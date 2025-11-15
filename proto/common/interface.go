// Package common defines the common plugin interface.
// All plugins must implement this interface for generic invocation.
package common

// PluginInterface is the common interface that all plugins must implement
type PluginInterface interface {
	// GetMetadata returns plugin metadata
	GetMetadata() (*MetadataResponse, error)

	// Execute executes a plugin method with generic parameters
	Execute(method string, params map[string]string) (*ExecuteResponse, error)
}
