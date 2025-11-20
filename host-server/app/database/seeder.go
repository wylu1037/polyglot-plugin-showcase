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
				Namespace:       "builtin",
				Name:            "converter",
				Version:         "1.0.0",
				Type:            "data-processing",
				Description:     "A plugin for data format conversion. Supports converting JSON to CSV, TXT, and HTML formats.",
				Status:          models.PluginStatusActive,
				BinaryPath:      "bin/plugins/builtin/data-processing/converter/v1.0.0/darwin_arm64/plugin",
				Protocol:        models.PluginProtocolGRPC,
				ProtocolVersion: 1,
				OS:              "darwin",
				Arch:            "arm64",
				Config: models.JSONMap{
					"default_format": "csv",
					"csv_delimiter":  ",",
					"html_styled":    true,
					"txt_format":     "key-value",
				},
				Metadata: models.JSONMap{
					"author":      "Polyglot Team",
					"license":     "MIT",
					"repository":  "https://github.com/wylu1037/polyglot-plugin-showcase",
					"tags":        []string{"conversion", "data-format", "csv", "html", "text"},
					"min_version": "1.0.0",
				},
			},
			{
				Namespace:       "builtin",
				Name:            "desensitization",
				Version:         "1.0.0",
				Type:            "data-processing",
				Description:     "A plugin for data desensitization. Supports masking, hashing, and tokenization.",
				Status:          models.PluginStatusActive,
				BinaryPath:      "bin/plugins/builtin/data-processing/desensitization/v1.0.0/darwin_arm64/plugin",
				Protocol:        models.PluginProtocolGRPC,
				ProtocolVersion: 1,
				OS:              "darwin",
				Arch:            "arm64",
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
			{
				Namespace:       "builtin",
				Name:            "dpanonymizer",
				Version:         "1.0.0",
				Type:            "data-processing",
				Description:     "A plugin for differential privacy anonymization. Supports Laplace/Gaussian noise addition, DP count, sum, mean, and variance calculations.",
				Status:          models.PluginStatusActive,
				BinaryPath:      "bin/plugins/builtin/data-processing/dpanonymizer/v1.0.0/darwin_arm64/plugin",
				Protocol:        models.PluginProtocolGRPC,
				ProtocolVersion: 1,
				OS:              "darwin",
				Arch:            "arm64",
				Config: models.JSONMap{
					"default_epsilon":            0.1,
					"default_delta":              1e-5,
					"default_sensitivity":        1.0,
					"max_partitions_contributed": 1,
				},
				Metadata: models.JSONMap{
					"author":       "Polyglot Team",
					"license":      "MIT",
					"repository":   "https://github.com/wylu1037/polyglot-plugin-showcase",
					"tags":         []string{"differential-privacy", "privacy", "anonymization", "statistics"},
					"min_version":  "1.0.0",
					"dependencies": []string{"github.com/google/differential-privacy/go/v3"},
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
