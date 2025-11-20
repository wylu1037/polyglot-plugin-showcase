package plugin

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

type PluginInterfaceInfo struct {
	PluginName      string // Plugin name (e.g., "converter", "desensitization")
	HandshakeConfig plugin.HandshakeConfig
	PluginMap       map[string]plugin.Plugin
}

type Registry struct {
	plugins map[string]*PluginInterfaceInfo // key: plugin name (e.g., "converter")
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]*PluginInterfaceInfo),
	}
}

func (r *Registry) Register(info *PluginInterfaceInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if info.PluginName == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	r.plugins[info.PluginName] = info
	return nil
}

func createPluginMap(pluginName string) map[string]plugin.Plugin {
	return map[string]plugin.Plugin{
		pluginName: &common.PluginGRPCPlugin{},
	}
}

func (r *Registry) GetPluginInterfaceInfo(pluginName string) (*PluginInterfaceInfo, error) {
	r.mu.RLock()
	info, ok := r.plugins[pluginName]
	r.mu.RUnlock()

	if !ok {
		return r.autoRegisterPluginInterface(pluginName)
	}
	return info, nil
}

func (r *Registry) autoRegisterPluginInterface(pluginName string) (*PluginInterfaceInfo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock
	if info, ok := r.plugins[pluginName]; ok {
		return info, nil
	}

	info := &PluginInterfaceInfo{
		PluginName:      pluginName,
		HandshakeConfig: common.Handshake,
		PluginMap:       createPluginMap(pluginName),
	}

	r.plugins[pluginName] = info
	fmt.Printf("✨ Auto-registered plugin interface: %s\n", pluginName)
	return info, nil
}
