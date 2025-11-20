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

var Module = fx.Options(
	fx.Provide(provideRegistry),
	fx.Provide(provideManager),
	fx.Invoke(autoLoadPlugins),
)

func provideRegistry() *Registry {
	return NewRegistry()
}

func provideManager(cfg *config.Config, registry *Registry) *Manager {
	return NewManager(registry, &ManagerConfig{
		DownloadTimeout: cfg.Plugin.DownloadTimeout,
		StartupTimeout:  cfg.Plugin.StartupTimeout,
	})
}

type autoLoadPluginsParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Manager   *Manager
	Repo      repository.PluginRepository
	DB        *gorm.DB
	Config    *config.Config
}

func autoLoadPlugins(p autoLoadPluginsParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			plugins, err := p.Repo.FindAll(map[string]any{
				"status": models.PluginStatusActive,
			})
			if err != nil {
				return fmt.Errorf("failed to find active plugins: %w", err)
			}

			if len(plugins) == 0 {
				fmt.Println("📦 No active plugins to load")
				return nil
			}

			for _, plugin := range plugins {
				if err := p.Manager.LoadPlugin(plugin.ID, plugin.BinaryPath, plugin.Name); err != nil {
					fmt.Printf("❌ Failed to load plugin %s (ID: %d): %v\n", plugin.Name, plugin.ID, err)
					p.Repo.UpdateStatus(plugin.ID, models.PluginStatusError)
				} else {
					fmt.Printf("✅ Successfully loaded plugin %s (ID: %d)\n", plugin.Name, plugin.ID)
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
