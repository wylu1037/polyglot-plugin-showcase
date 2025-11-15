package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
)

type Manager struct {
	registry         *Registry
	discovery        *Discovery
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
		discovery:        NewDiscovery(registry),
		clients:          make(map[uint]*plugin.Client),
		clientInterfaces: make(map[uint]any),
		downloadTimeout:  config.DownloadTimeout,
		startupTimeout:   config.StartupTimeout,
	}

	return m
}

func (m *Manager) DownloadPlugin(url, destPath string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: m.downloadTimeout,
	}

	// Make request
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download plugin: HTTP %d", resp.StatusCode)
	}

	// Create destination directory if not exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create temporary file
	tempFile := destPath + ".tmp"
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer out.Close()

	// Copy content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to save plugin: %w", err)
	}

	// Make file executable
	if err := os.Chmod(tempFile, 0755); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to make plugin executable: %w", err)
	}

	// Rename to final destination
	if err := os.Rename(tempFile, destPath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to move plugin to destination: %w", err)
	}

	return nil
}

// VerifyChecksum verifies the checksum of a file
func (m *Manager) VerifyChecksum(filePath, expectedChecksum string) error {
	if expectedChecksum == "" {
		return nil // Skip verification if no checksum provided
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for checksum verification: %w", err)
	}
	defer file.Close()

	// Calculate SHA256
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Get hex string
	actualChecksum := hex.EncodeToString(hash.Sum(nil))

	// Compare
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// LoadPlugin loads a plugin and returns the client
func (m *Manager) LoadPlugin(pluginID uint, pluginPath string, pluginType models.PluginType) error {
	// Check if plugin is already loaded
	m.mu.RLock()
	if _, exists := m.clients[pluginID]; exists {
		m.mu.RUnlock()
		return nil // Already loaded
	}
	m.mu.RUnlock()

	// Get plugin type info from registry
	typeInfo, err := m.registry.GetTypeInfo(pluginType)
	if err != nil {
		return fmt.Errorf("failed to get plugin type info: %w", err)
	}

	// Validate plugin binary
	if err := m.discovery.ValidatePluginBinary(pluginPath); err != nil {
		return fmt.Errorf("invalid plugin binary: %w", err)
	}

	// Create plugin client
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: typeInfo.HandshakeConfig,
		Plugins:         typeInfo.PluginMap,
		Cmd:             exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
		StartTimeout: m.startupTimeout,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to connect to plugin: %w", err)
	}

	// Dispense plugin interface
	raw, err := rpcClient.Dispense(typeInfo.InterfaceName)
	if err != nil {
		client.Kill()
		return fmt.Errorf("failed to dispense plugin: %w", err)
	}

	// Verify protocol version compatibility (similar to Terraform's approach)
	if err := m.verifyProtocolVersion(raw); err != nil {
		client.Kill()
		return fmt.Errorf("protocol version incompatible: %w", err)
	}

	// Store client and interface
	m.mu.Lock()
	m.clients[pluginID] = client
	m.clientInterfaces[pluginID] = raw
	m.mu.Unlock()

	return nil
}

// UnloadPlugin unloads a plugin
func (m *Manager) UnloadPlugin(pluginID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clients[pluginID]
	if !exists {
		return nil // Already unloaded
	}

	// Kill plugin process
	client.Kill()

	// Remove from maps
	delete(m.clients, pluginID)
	delete(m.clientInterfaces, pluginID)

	return nil
}

// GetPluginClient returns the plugin client interface
func (m *Manager) GetPluginClient(pluginID uint) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clientInterface, exists := m.clientInterfaces[pluginID]
	if !exists {
		return nil, fmt.Errorf("plugin not loaded")
	}

	return clientInterface, nil
}

// IsPluginLoaded checks if a plugin is loaded
func (m *Manager) IsPluginLoaded(pluginID uint) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.clients[pluginID]
	return exists
}

// UnloadAll unloads all plugins
func (m *Manager) UnloadAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, client := range m.clients {
		client.Kill()
	}

	m.clients = make(map[uint]*plugin.Client)
	m.clientInterfaces = make(map[uint]any)
}

// GetLoadedPluginIDs returns all loaded plugin IDs
func (m *Manager) GetLoadedPluginIDs() []uint {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]uint, 0, len(m.clients))
	for id := range m.clients {
		ids = append(ids, id)
	}
	return ids
}

// verifyProtocolVersion verifies that the plugin's protocol version is compatible
// with the host. This follows Terraform's approach of checking version ranges.
func (m *Manager) verifyProtocolVersion(pluginInterface any) error {
	// Import common package to access version checking
	// Try to get metadata from the plugin
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
