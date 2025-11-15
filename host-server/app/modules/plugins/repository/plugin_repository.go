package repository

import (
	"fmt"

	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
	"gorm.io/gorm"
)

// PluginRepository defines the interface for plugin data access
type PluginRepository interface {
	Create(plugin *models.Plugin) error
	FindByID(id uint) (*models.Plugin, error)
	FindAll(filters map[string]any) ([]*models.Plugin, error)
	Update(plugin *models.Plugin) error
	Delete(id uint) error
	UpdateStatus(id uint, status models.PluginStatus) error
	FindByType(pluginType models.PluginType) ([]*models.Plugin, error)
	FindByNameAndVersion(name, version string) (*models.Plugin, error)
	UpdateLastUsedAt(id uint, timestamp int64) error
}

type pluginRepository struct {
	db *gorm.DB
}

// NewPluginRepository creates a new plugin repository
func NewPluginRepository(db *gorm.DB) PluginRepository {
	return &pluginRepository{
		db: db,
	}
}

// Create creates a new plugin record
func (r *pluginRepository) Create(plugin *models.Plugin) error {
	if err := r.db.Create(plugin).Error; err != nil {
		return fmt.Errorf("failed to create plugin: %w", err)
	}
	return nil
}

// FindByID finds a plugin by ID
func (r *pluginRepository) FindByID(id uint) (*models.Plugin, error) {
	var plugin models.Plugin
	if err := r.db.First(&plugin, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("plugin not found")
		}
		return nil, fmt.Errorf("failed to find plugin: %w", err)
	}
	return &plugin, nil
}

// FindAll finds all plugins with optional filters
func (r *pluginRepository) FindAll(filters map[string]any) ([]*models.Plugin, error) {
	var plugins []*models.Plugin
	query := r.db.Model(&models.Plugin{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&plugins).Error; err != nil {
		return nil, fmt.Errorf("failed to find plugins: %w", err)
	}

	return plugins, nil
}

// Update updates a plugin record
func (r *pluginRepository) Update(plugin *models.Plugin) error {
	if err := r.db.Save(plugin).Error; err != nil {
		return fmt.Errorf("failed to update plugin: %w", err)
	}
	return nil
}

// Delete deletes a plugin by ID
func (r *pluginRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Plugin{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete plugin: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of a plugin
func (r *pluginRepository) UpdateStatus(id uint, status models.PluginStatus) error {
	if err := r.db.Model(&models.Plugin{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update plugin status: %w", err)
	}
	return nil
}

// FindByType finds all plugins of a specific type
func (r *pluginRepository) FindByType(pluginType models.PluginType) ([]*models.Plugin, error) {
	var plugins []*models.Plugin
	if err := r.db.Where("type = ?", pluginType).Find(&plugins).Error; err != nil {
		return nil, fmt.Errorf("failed to find plugins by type: %w", err)
	}
	return plugins, nil
}

// FindByNameAndVersion finds a plugin by name and version
func (r *pluginRepository) FindByNameAndVersion(name, version string) (*models.Plugin, error) {
	var plugin models.Plugin
	if err := r.db.Where("name = ? AND version = ?", name, version).First(&plugin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found is not an error
		}
		return nil, fmt.Errorf("failed to find plugin: %w", err)
	}
	return &plugin, nil
}

// UpdateLastUsedAt updates the last used timestamp
func (r *pluginRepository) UpdateLastUsedAt(id uint, timestamp int64) error {
	if err := r.db.Model(&models.Plugin{}).Where("id = ?", id).Update("last_used_at", timestamp).Error; err != nil {
		return fmt.Errorf("failed to update last used at: %w", err)
	}
	return nil
}
