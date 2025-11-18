package database

import (
	"fmt"
	"log"

	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
	"gorm.io/gorm"
)

// AutoMigrate runs database schema migrations and seeds initial data.
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	if err := db.AutoMigrate(&models.Plugin{}); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully.")

	// Seed initial data
	if err := Seed(db); err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

	return nil
}
