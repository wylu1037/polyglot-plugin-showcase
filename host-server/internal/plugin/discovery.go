package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
)

// PluginInfo contains information about a discovered plugin
type PluginInfo struct {
	Name       string
	Version    string
	Type       models.PluginType
	BinaryPath string
}

// Discovery handles plugin discovery from filesystem
type Discovery struct {
	registry *Registry
}

// NewDiscovery creates a new plugin discovery instance
func NewDiscovery(registry *Registry) *Discovery {
	return &Discovery{
		registry: registry,
	}
}

// DiscoverPlugins scans a directory for plugin binaries
// Expected directory structure: {pluginDir}/{type}/{name}_{version}
// Example: ./bin/plugins/desensitization/my-plugin_v1.0.0
func (d *Discovery) DiscoverPlugins(pluginDir string) ([]PluginInfo, error) {
	var plugins []PluginInfo

	// Check if plugin directory exists
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return plugins, nil
	}

	// Walk through plugin directory
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file is executable
		if info.Mode()&0111 == 0 {
			return nil
		}

		// Parse plugin info from path
		pluginInfo, err := d.parsePluginPath(pluginDir, path)
		if err != nil {
			// Log error but continue discovery
			return nil
		}

		plugins = append(plugins, *pluginInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover plugins: %w", err)
	}

	return plugins, nil
}

// parsePluginPath extracts plugin information from file path
// Expected format: {pluginDir}/{type}/{name}_{version}
func (d *Discovery) parsePluginPath(pluginDir, fullPath string) (*PluginInfo, error) {
	// Get relative path from plugin directory
	relPath, err := filepath.Rel(pluginDir, fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	// Split path into components
	parts := strings.Split(relPath, string(os.PathSeparator))
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid plugin path structure: %s", relPath)
	}

	// Extract type from directory name
	typeStr := parts[0]
	pluginType := models.PluginType(typeStr)

	// Verify plugin type is supported
	if !d.registry.IsSupported(pluginType) {
		return nil, fmt.Errorf("unsupported plugin type: %s", typeStr)
	}

	// Extract name and version from filename
	filename := parts[len(parts)-1]
	name, version, err := parsePluginFilename(filename)
	if err != nil {
		return nil, err
	}

	return &PluginInfo{
		Name:       name,
		Version:    version,
		Type:       pluginType,
		BinaryPath: fullPath,
	}, nil
}

// parsePluginFilename extracts name and version from filename
// Expected format: {name}_{version}
// Example: my-plugin_v1.0.0
func parsePluginFilename(filename string) (name, version string, err error) {
	// Remove extension if present
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))

	// Split by underscore
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid plugin filename format: %s (expected: name_version)", filename)
	}

	// Last part is version, everything else is name
	version = parts[len(parts)-1]
	name = strings.Join(parts[:len(parts)-1], "_")

	if name == "" || version == "" {
		return "", "", fmt.Errorf("invalid plugin filename: name or version is empty")
	}

	return name, version, nil
}

// ValidatePluginBinary checks if a file is a valid plugin binary
func (d *Discovery) ValidatePluginBinary(path string) error {
	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("plugin binary not found: %w", err)
	}

	// Check if it's a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("plugin binary is not a regular file")
	}

	// Check if it's executable
	if info.Mode()&0111 == 0 {
		return fmt.Errorf("plugin binary is not executable")
	}

	return nil
}
