package plugin

import (
	"context"
	"fmt"

	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
	"github.com/wylu1037/polyglot-plugin-host-server/app/modules/plugins/repository"
	"github.com/wylu1037/polyglot-plugin-host-server/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func ProvideRegistry() *Registry {
	return NewRegistry()
}

func ProvideManager(cfg *config.Config, registry *Registry) *Manager {
	return NewManager(registry, &ManagerConfig{
		DownloadTimeout: cfg.Plugin.DownloadTimeout,
		StartupTimeout:  cfg.Plugin.StartupTimeout,
	})
}

func ProvidePluginDir(cfg *config.Config) string {
	return cfg.Plugin.Dir
}

type AutoLoadPluginsParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Manager   *Manager
	Repo      repository.PluginRepository
	DB        *gorm.DB
	Config    *config.Config
}

func AutoLoadPlugins(p AutoLoadPluginsParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Find all active plugins
			plugins, err := p.Repo.FindAll(map[string]any{
				"status": models.PluginStatusActive,
			})
			if err != nil {
				return fmt.Errorf("failed to find active plugins: %w", err)
			}

			// Load each plugin
			for _, plugin := range plugins {
				if err := p.Manager.LoadPlugin(plugin.ID, plugin.BinaryPath, plugin.Type); err != nil {
					// Log error but don't fail startup
					fmt.Printf("Failed to load plugin %s (ID: %d): %v\n", plugin.Name, plugin.ID, err)
					// Update plugin status to error
					p.Repo.UpdateStatus(plugin.ID, models.PluginStatusError)
				} else {
					fmt.Printf("Successfully loaded plugin %s (ID: %d)\n", plugin.Name, plugin.ID)
				}
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Manager.UnloadAll()
			return nil
		},
	})
}
