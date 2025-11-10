package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// AutoMigrate runs automatic migration for all models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	models := []any{
		// Add models here
	}

	// Run migrations
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
