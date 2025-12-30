package migrationModels

import (
	"context"

	"gorm.io/gorm"
)

type Migration struct {
	MigrationId         string
	PreviousMigrationID string
	MigrateFunc         func(context.Context, *gorm.DB) error
	Description         string
}

func (m Migration) Apply(ctx context.Context, db *gorm.DB) error {
	if m.shouldRunMigration(db) {
		err := m.MigrateFunc(ctx, db)
		if err != nil {
			return err
		}
		m.createMigrationIdIfNotExists(db)
	}
	return nil
}

func getCurrentMigration(db *gorm.DB) (MigrationModel, error) {
	ctx := context.Background()
	return gorm.G[MigrationModel](db, gorm.WithResult()).Last(ctx)
}

func (m Migration) shouldRunMigration(db *gorm.DB) bool {
	previousMigration, err := getCurrentMigration(db)

	// If there is no previous migration and the previous migration ID is empty, run the migration
	if err != nil && m.PreviousMigrationID == "" {
		return true
		// If there is a previous migration and the previous migration ID is the same as the current migration ID, run the migration
	} else if err == nil && m.PreviousMigrationID == previousMigration.ID && m.PreviousMigrationID != "" {
		return true
	}
	return false
}

// createMigrationIdIfNotExists creates a migration ID if it does not exist
// It creates a migration ID if it does not exist
func (m Migration) createMigrationIdIfNotExists(db *gorm.DB) {
	ctx := context.Background()
	db.WithContext(ctx).Delete(&MigrationModel{ID: m.PreviousMigrationID})
	db.WithContext(ctx).Create(&MigrationModel{ID: m.MigrationId})
}
