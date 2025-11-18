package database

import (
	"fmt"
	"log"

	"github.com/wylu1037/polyglot-plugin-host-server/app/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SeedData defines the structure for seed data
type SeedData struct {
	Plugins []models.Plugin
}

// GetSeedData returns the seed data for the database
func GetSeedData() *SeedData {
	return &SeedData{
		Plugins: []models.Plugin{
			{
				Name:            "desensitization",
				Version:         "1.0.0",
				Type:            models.PluginTypeDesensitization,
				Description:     "A plugin for data desensitization. Supports masking, hashing, and tokenization.",
				Status:          models.PluginStatusInactive,
				BinaryPath:      "bin/plugins/desensitization/desensitization_v1.0.0",
				Protocol:        models.PluginProtocolGRPC,
				ProtocolVersion: 1,
				Config: models.JSONMap{
					"default_strategy": "mask",
					"mask_char":        "*",
				},
				Metadata: models.JSONMap{
					"author":      "Polyglot Team",
					"license":     "MIT",
					"repository":  "https://github.com/wylu1037/polyglot-plugin-showcase",
					"tags":        []string{"security", "privacy", "data-protection"},
					"min_version": "1.0.0",
				},
			},
		},
	}
}

func Seed(db *gorm.DB) error {
	log.Println("🌱 Starting database seeding...")

	seedData := GetSeedData()

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := seedPlugins(tx, seedData.Plugins); err != nil {
			return fmt.Errorf("failed to seed plugins: %w", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("❌ Seeding failed: %v\n", err)
		return err
	}

	log.Println("✅ Database seeding completed successfully")
	return nil
}

func seedPlugins(tx *gorm.DB, plugins []models.Plugin) error {
	if len(plugins) == 0 {
		log.Println("⏭️  No plugins to seed, skipping...")
		return nil
	}

	log.Printf("📦 Seeding %d plugin(s)...\n", len(plugins))

	for _, plugin := range plugins {
		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "name"}, {Name: "version"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"type",
				"description",
				"binary_path",
				"protocol",
				"protocol_version",
				"config",
				"metadata",
			}),
		}).Create(&plugin)

		if result.Error != nil {
			return fmt.Errorf("failed to upsert plugin '%s@%s': %w", plugin.Name, plugin.Version, result.Error)
		}

		if result.RowsAffected > 0 {
			log.Printf("✓ Plugin upserted: %s@%s\n", plugin.Name, plugin.Version)
		}
	}

	return nil
}
