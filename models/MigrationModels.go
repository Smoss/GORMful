package models

// MigrationModel is a model for the migration table
// It is used to track the migrations that have been applied to the database
// The ID is a UUID and is used to track the migrations that have been applied to the database
type MigrationModel struct {
	ID string `gorm:"primaryKey"`
}
