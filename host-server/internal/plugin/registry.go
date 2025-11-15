package plugin

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

// PluginTypeInfo contains metadata about a plugin type
type PluginTypeInfo struct {
	Type            models.PluginType
	HandshakeConfig plugin.HandshakeConfig
	PluginMap       map[string]plugin.Plugin
	InterfaceName   string
}

// Registry manages plugin type registrations
type Registry struct {
	types map[models.PluginType]*PluginTypeInfo
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	r := &Registry{
		types: make(map[models.PluginType]*PluginTypeInfo),
	}
	r.registerBuiltinTypes()
	return r
}

// registerBuiltinTypes registers all built-in plugin types
// Each plugin type uses its type name as the interface name for better flexibility
func (r *Registry) registerBuiltinTypes() {
	// All plugin types use the same common handshake
	commonHandshake := common.Handshake

	// Register desensitization plugin
	r.Register(&PluginTypeInfo{
		Type:            models.PluginTypeDesensitization,
		HandshakeConfig: commonHandshake,
		PluginMap:       createPluginMap(string(models.PluginTypeDesensitization)),
		InterfaceName:   string(models.PluginTypeDesensitization),
	})

	// Register encryption plugin
	r.Register(&PluginTypeInfo{
		Type:            models.PluginTypeEncryption,
		HandshakeConfig: commonHandshake,
		PluginMap:       createPluginMap(string(models.PluginTypeEncryption)),
		InterfaceName:   string(models.PluginTypeEncryption),
	})

	// Register validation plugin
	r.Register(&PluginTypeInfo{
		Type:            models.PluginTypeValidation,
		HandshakeConfig: commonHandshake,
		PluginMap:       createPluginMap(string(models.PluginTypeValidation)),
		InterfaceName:   string(models.PluginTypeValidation),
	})

	// Register transform plugin
	r.Register(&PluginTypeInfo{
		Type:            models.PluginTypeTransform,
		HandshakeConfig: commonHandshake,
		PluginMap:       createPluginMap(string(models.PluginTypeTransform)),
		InterfaceName:   string(models.PluginTypeTransform),
	})

	// Register custom plugin
	r.Register(&PluginTypeInfo{
		Type:            models.PluginTypeCustom,
		HandshakeConfig: commonHandshake,
		PluginMap:       createPluginMap(string(models.PluginTypeCustom)),
		InterfaceName:   string(models.PluginTypeCustom),
	})
}

// createPluginMap creates a plugin map with the given interface name
func createPluginMap(interfaceName string) map[string]plugin.Plugin {
	return map[string]plugin.Plugin{
		interfaceName: &common.PluginGRPCPlugin{},
	}
}

// Register registers a new plugin type
func (r *Registry) Register(info *PluginTypeInfo) {
	r.types[info.Type] = info
}

// GetTypeInfo returns the plugin type info for a given type
func (r *Registry) GetTypeInfo(pluginType models.PluginType) (*PluginTypeInfo, error) {
	info, ok := r.types[pluginType]
	if !ok {
		return nil, fmt.Errorf("plugin type %s not registered", pluginType)
	}
	return info, nil
}

// GetHandshake returns the handshake config for a plugin type
func (r *Registry) GetHandshake(pluginType models.PluginType) (plugin.HandshakeConfig, error) {
	info, err := r.GetTypeInfo(pluginType)
	if err != nil {
		return plugin.HandshakeConfig{}, err
	}
	return info.HandshakeConfig, nil
}

// GetPluginMap returns the plugin map for a plugin type
func (r *Registry) GetPluginMap(pluginType models.PluginType) (map[string]plugin.Plugin, error) {
	info, err := r.GetTypeInfo(pluginType)
	if err != nil {
		return nil, err
	}
	return info.PluginMap, nil
}

// IsSupported checks if a plugin type is supported
func (r *Registry) IsSupported(pluginType models.PluginType) bool {
	_, ok := r.types[pluginType]
	return ok
}

// GetSupportedTypes returns all supported plugin types
func (r *Registry) GetSupportedTypes() []models.PluginType {
	types := make([]models.PluginType, 0, len(r.types))
	for t := range r.types {
		types = append(types, t)
	}
	return types
}
