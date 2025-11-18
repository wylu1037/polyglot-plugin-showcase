package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/repository"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/request"
	"github.com/wylu1037/polyglot-plugin-host-server/internal/errors"
	"github.com/wylu1037/polyglot-plugin-host-server/internal/plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

type PluginService interface {
	InstallPlugin(req *request.InstallPluginRequest) (*models.Plugin, error)
	ActivatePlugin(id uint) error
	DeactivatePlugin(id uint) error
	UninstallPlugin(id uint) error
	ListPlugins(filters map[string]any) ([]*models.Plugin, error)
	GetPluginInfo(id uint) (*models.Plugin, error)
	CallPlugin(id uint, req *request.CallPluginRequest) (any, error)
}

type pluginService struct {
	repo      repository.PluginRepository
	manager   *plugin.Manager
	pluginDir string
}

func NewPluginService(
	repo repository.PluginRepository,
	manager *plugin.Manager,
	pluginDir string,
) PluginService {
	return &pluginService{
		repo:      repo,
		manager:   manager,
		pluginDir: pluginDir,
	}
}

func (s *pluginService) InstallPlugin(req *request.InstallPluginRequest) (*models.Plugin, error) {
	// Check if plugin already exists
	existing, err := s.repo.FindByNameAndVersion(req.Name, req.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing plugin: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("plugin %s version %s already exists", req.Name, req.Version)
	}

	// Construct binary path: {pluginDir}/{type}/{name}_{version}
	binaryPath := filepath.Join(s.pluginDir, string(req.Type), fmt.Sprintf("%s_%s", req.Name, req.Version))

	// Create plugin record with installing status
	pluginRecord := &models.Plugin{
		Name:            req.Name,
		Version:         req.Version,
		Type:            req.Type,
		Description:     req.Description,
		Status:          models.PluginStatusInstalling,
		BinaryPath:      binaryPath,
		DownloadURL:     req.DownloadURL,
		Protocol:        models.PluginProtocolGRPC,
		ProtocolVersion: 1,
		Config:          req.Config,
		Metadata:        req.Metadata,
	}

	if err := s.repo.Create(pluginRecord); err != nil {
		return nil, fmt.Errorf("failed to create plugin record: %w", err)
	}

	// Download plugin
	if err := s.manager.DownloadPlugin(req.DownloadURL, binaryPath); err != nil {
		s.repo.UpdateStatus(pluginRecord.ID, models.PluginStatusError)
		return nil, fmt.Errorf("failed to download plugin: %w", err)
	}

	// Update status to inactive (installed but not activated)
	if err := s.repo.UpdateStatus(pluginRecord.ID, models.PluginStatusInactive); err != nil {
		return nil, fmt.Errorf("failed to update plugin status: %w", err)
	}

	// Reload plugin record
	return s.repo.FindByID(pluginRecord.ID)
}

func (s *pluginService) ActivatePlugin(id uint) error {
	// Get plugin record
	pluginRecord, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find plugin: %w", err)
	}

	// Check if already active
	if pluginRecord.Status == models.PluginStatusActive {
		return nil
	}

	// Check if plugin is in a loadable state
	if pluginRecord.Status != models.PluginStatusInactive && pluginRecord.Status != models.PluginStatusDisabled {
		return fmt.Errorf("plugin cannot be activated from status: %s", pluginRecord.Status)
	}

	// Load plugin
	if err := s.manager.LoadPlugin(id, pluginRecord.BinaryPath, pluginRecord.Type); err != nil {
		s.repo.UpdateStatus(id, models.PluginStatusError)
		return fmt.Errorf("failed to load plugin: %w", err)
	}

	// Update status to active
	if err := s.repo.UpdateStatus(id, models.PluginStatusActive); err != nil {
		s.manager.UnloadPlugin(id)
		return fmt.Errorf("failed to update plugin status: %w", err)
	}

	return nil
}

func (s *pluginService) DeactivatePlugin(id uint) error {
	// Get plugin record
	pluginRecord, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find plugin: %w", err)
	}

	// Check if already inactive
	if pluginRecord.Status == models.PluginStatusInactive {
		return nil
	}

	// Unload plugin
	if err := s.manager.UnloadPlugin(id); err != nil {
		return fmt.Errorf("failed to unload plugin: %w", err)
	}

	// Update status to inactive
	if err := s.repo.UpdateStatus(id, models.PluginStatusInactive); err != nil {
		return fmt.Errorf("failed to update plugin status: %w", err)
	}

	return nil
}

func (s *pluginService) UninstallPlugin(id uint) error {
	// Get plugin record
	pluginRecord, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find plugin: %w", err)
	}

	// Deactivate if active
	if pluginRecord.Status == models.PluginStatusActive {
		if err := s.DeactivatePlugin(id); err != nil {
			return fmt.Errorf("failed to deactivate plugin: %w", err)
		}
	}

	// Remove binary file
	if err := os.Remove(pluginRecord.BinaryPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plugin binary: %w", err)
	}

	// Delete database record
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete plugin record: %w", err)
	}

	return nil
}

func (s *pluginService) ListPlugins(filters map[string]any) ([]*models.Plugin, error) {
	return s.repo.FindAll(filters)
}

func (s *pluginService) GetPluginInfo(id uint) (*models.Plugin, error) {
	return s.repo.FindByID(id)
}

func (s *pluginService) CallPlugin(id uint, req *request.CallPluginRequest) (any, error) {
	pluginRecord, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.ErrPluginNotFound.WithInternal(err)
	}

	if pluginRecord.Status != models.PluginStatusActive {
		return nil, fmt.Errorf("plugin is not active")
	}

	now := time.Now().Unix()
	s.repo.UpdateLastUsedAt(id, now)

	clientInterface, err := s.manager.GetPluginClient(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin client: %w", err)
	}

	// Type assert to common plugin interface
	pluginClient, ok := clientInterface.(common.PluginInterface)
	if !ok {
		return nil, fmt.Errorf("plugin does not implement common.PluginInterface")
	}

	// Convert params from map[string]any to map[string]string
	stringParams := make(map[string]string)
	for key, value := range req.Params {
		if strValue, ok := value.(string); ok {
			stringParams[key] = strValue
		} else {
			stringParams[key] = fmt.Sprintf("%v", value)
		}
	}

	// Execute plugin method using generic interface
	result, err := pluginClient.Execute(req.Method, stringParams)
	if err != nil {
		return nil, fmt.Errorf("plugin execution failed: %w", err)
	}

	if !result.Success {
		errMsg := "unknown error"
		if result.Error != nil {
			errMsg = *result.Error
		}
		return nil, fmt.Errorf("plugin returned error: %s", errMsg)
	}

	if result.Result == nil {
		return "", nil // Success but no result
	}

	return *result.Result, nil
}
