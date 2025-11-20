package plugin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
)

type Manager struct {
	registry         *Registry
	clients          map[uint]*plugin.Client
	clientInterfaces map[uint]any
	mu               sync.RWMutex
	downloadTimeout  time.Duration
	startupTimeout   time.Duration
}

type ManagerConfig struct {
	DownloadTimeout time.Duration
	StartupTimeout  time.Duration
}

func NewManager(registry *Registry, config *ManagerConfig) *Manager {
	if config == nil {
		config = &ManagerConfig{
			DownloadTimeout: 5 * time.Minute,
			StartupTimeout:  30 * time.Second,
		}
	}

	m := &Manager{
		registry:         registry,
		clients:          make(map[uint]*plugin.Client),
		clientInterfaces: make(map[uint]any),
		downloadTimeout:  config.DownloadTimeout,
		startupTimeout:   config.StartupTimeout,
	}

	return m
}

func (m *Manager) DownloadPlugin(url, destPath string) error {
	client := &http.Client{
		Timeout: m.downloadTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download plugin: HTTP %d", resp.StatusCode)
	}

	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	tempFile := destPath + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to save plugin: %w", err)
	}

	if err := os.Chmod(tempFile, 0755); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to make plugin executable: %w", err)
	}

	if err := os.Rename(tempFile, destPath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to move plugin to destination: %w", err)
	}

	return nil
}

func (m *Manager) LoadPlugin(pluginID uint, pluginPath string, pluginName string) error {
	m.mu.RLock()
	if _, exists := m.clients[pluginID]; exists {
		m.mu.RUnlock()
		return nil
	}
	m.mu.RUnlock()

	pluginInterfaceInfo, err := m.registry.GetPluginInterfaceInfo(pluginName)
	if err != nil {
		return fmt.Errorf("failed to get plugin interface info for '%s': %w", pluginName, err)
	}

	if err := m.validatePluginBinary(pluginPath); err != nil {
		return fmt.Errorf("invalid plugin binary at '%s': %w", pluginPath, err)
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: pluginInterfaceInfo.HandshakeConfig,
		Plugins:         pluginInterfaceInfo.PluginMap,
		Cmd:             exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
		StartTimeout: m.startupTimeout,
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to connect to plugin: %w", err)
	}

	raw, err := rpcClient.Dispense(pluginInterfaceInfo.PluginName)
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense plugin interface '%s': %w", pluginInterfaceInfo.PluginName, err)
	}

	if err := m.verifyProtocolVersion(raw); err != nil {
		client.Kill()
		return fmt.Errorf("protocol version incompatible: %w", err)
	}

	m.mu.Lock()
	m.clients[pluginID] = client
	m.clientInterfaces[pluginID] = raw
	m.mu.Unlock()

	return nil
}

func (m *Manager) UnloadPlugin(pluginID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[pluginID]
	if !exists {
		return nil
	}

	client.Kill()

	delete(m.clients, pluginID)
	delete(m.clientInterfaces, pluginID)

	return nil
}

func (m *Manager) GetPluginClient(pluginID uint) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clientInterface, exists := m.clientInterfaces[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin not loaded")
	}

	return clientInterface, nil
}

func (m *Manager) UnloadAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, client := range m.clients {
		client.Kill()
	}

	m.clients = make(map[uint]*plugin.Client)
	m.clientInterfaces = make(map[uint]any)
}

// validatePluginBinary validates that the plugin binary exists and is executable
func (m *Manager) validatePluginBinary(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("plugin binary not found: %w", err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("plugin binary is not a regular file")
	}

	if info.Mode()&0111 == 0 {
		return fmt.Errorf("plugin binary is not executable")
	}

	return nil
}

// verifyProtocolVersion verifies that the plugin's protocol version is compatible
// with the host. This follows Terraform's approach of checking version ranges.
func (m *Manager) verifyProtocolVersion(pluginInterface any) error {
	type MetadataGetter interface {
		GetMetadata() (*struct {
			Name            string
			Version         string
			Description     string
			Methods         []string
			Capabilities    map[string]string
			ProtocolVersion int32
		}, error)
	}

	getter, ok := pluginInterface.(MetadataGetter)
	if !ok {
		// If plugin doesn't implement GetMetadata, assume it's compatible
		// (for backward compatibility with plugins that don't report version)
		return nil
	}

	metadata, err := getter.GetMetadata()
	if err != nil {
		return fmt.Errorf("failed to get plugin metadata: %w", err)
	}

	pluginVersion := int(metadata.ProtocolVersion)

	// If plugin doesn't report version (0), assume version 1 for backward compatibility
	if pluginVersion == 0 {
		pluginVersion = 1
	}

	// Check if version is in supported range
	// We need to import the common package constants
	// For now, hardcode the check (will be replaced with proper import)
	minVersion := 1
	maxVersion := 1

	if pluginVersion < minVersion || pluginVersion > maxVersion {
		return fmt.Errorf(
			"plugin protocol version %d is not compatible (host supports versions %d-%d)",
			pluginVersion,
			minVersion,
			maxVersion,
		)
	}

	return nil
}
