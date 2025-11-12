package database

import (
	"fmt"
	"log"

	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
	"gorm.io/gorm"
)

// AutoMigrate runs automatic migration for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	models := []any{
		&models.PluginStore{},
		// Add other models here
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
