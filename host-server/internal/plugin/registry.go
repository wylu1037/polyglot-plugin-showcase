package plugin

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

type PluginClientConfig struct {
	PluginName      string // Plugin name (e.g., "converter", "desensitization")
	HandshakeConfig plugin.HandshakeConfig
	PluginMap       map[string]plugin.Plugin
}

type Registry struct {
	configs map[string]*PluginClientConfig // key: plugin name (e.g., "converter")
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		configs: make(map[string]*PluginClientConfig),
	}
}

func (r *Registry) Register(config *PluginClientConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if config.PluginName == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	r.configs[config.PluginName] = config
	return nil
}

func createPluginMap(pluginName string) map[string]plugin.Plugin {
	return map[string]plugin.Plugin{
		pluginName: &common.PluginGRPCPlugin{},
	}
}

func (r *Registry) GetPluginConfig(pluginName string) (*PluginClientConfig, error) {
	r.mu.RLock()
	config, ok := r.configs[pluginName]
	r.mu.RUnlock()

	if !ok {
		return r.autoRegisterPlugin(pluginName)
	}
	return config, nil
}

func (r *Registry) autoRegisterPlugin(pluginName string) (*PluginClientConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock
	if config, ok := r.configs[pluginName]; ok {
		return config, nil
	}

	config := &PluginClientConfig{
		PluginName:      pluginName,
		HandshakeConfig: common.Handshake,
		PluginMap:       createPluginMap(pluginName),
	}

	r.configs[pluginName] = config
	fmt.Printf("✨ Auto-registered plugin: %s\n", pluginName)
	return config, nil
}
