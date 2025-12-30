package migrationModels

var MigrationTemplate = `
package db_migrations

import (
	"context"

	migration_models "github.com/GORMful/models"
	"gorm.io/gorm"
)

func migrate<migrationName>(ctx context.Context, db *gorm.DB) error {
	// Your migration code here
	return nil
}

var Migration<migrationName> = migration_models.Migration{
	MigrationId:         "<migrationId>",
	PreviousMigrationID: "<previousMigrationId>",
	MigrateFunc:         migrate<migrationName>,
	Description:         "<description>",
}
`
