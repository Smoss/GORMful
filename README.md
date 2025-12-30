# GORMful

GORMful is a migration management tool for GORM (Go Object-Relational Mapping) that helps you manage database migrations in a structured and trackable way.

## Installation

```bash
go get github.com/smoss/GORMful
```

## Using the Command Line Tool

The GORMful CLI tool helps you generate migration files with proper structure and tracking.

### Creating a Migration

To create a new migration file, use the `create-migration` command:

```bash
gormful create-migration --description "Add users table"
```

### Required Flags and Environment Variables

The CLI requires database connection information. You can provide this via command-line flags or environment variables:

**Command-line flags:**
- `--db-host`: Database host (defaults to `DB_HOST` environment variable)
- `--db-user`: Database user (defaults to `DB_USER` environment variable)
- `--db-password`: Database password (defaults to `DB_PASSWORD` environment variable)
- `--db-name`: Database name (defaults to `DB_NAME` environment variable)
- `--db-port`: Database port (defaults to `DB_PORT` environment variable)
- `--ssl-mode`: SSL mode for the database connection (defaults to `disable`)

**Migration-specific flags:**
- `--description`: Description for the migration (required)
- `--output-directory`: Output directory for the migration file (defaults to `migrations/migrations`)

### Example Usage

**Using environment variables:**
```bash
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=mypassword
export DB_NAME=mydb
export DB_PORT=5432

gormful create-migration --description "Add users table"
```

**Using command-line flags:**
```bash
gormful create-migration \
  --db-host localhost \
  --db-user postgres \
  --db-password mypassword \
  --db-name mydb \
  --db-port 5432 \
  --ssl-mode disable \
  --description "Add users table" \
  --output-directory migrations/migrations
```

### Generated Migration File

The CLI generates a migration file with the following structure:

```go
package db_migrations

import (
	"context"

	migration_models "github.com/smoss/GORMful/models"
	"gorm.io/gorm"
)

func migrate20240101120000(ctx context.Context, db *gorm.DB) error {
	// Your migration code here
	return nil
}

var Migration20240101120000 = migration_models.Migration{
	MigrationId:         "uuid-here",
	PreviousMigrationID: "previous-uuid-or-empty",
	MigrateFunc:         migrate20240101120000,
	Description:         "Add users table",
}
```

## Using the Library

You can also use GORMful programmatically in your Go application to manage migrations.

### Basic Setup

First, ensure you have a GORM database connection:

```go
import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

dsn := "host=localhost user=postgres password=mypassword dbname=mydb port=5432 sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
	log.Fatal("Failed to connect to database")
}
```

### Creating a Migration

Create a `Migration` struct with your migration logic:

```go
import (
	"context"
	"github.com/smoss/GORMful/models"
	"gorm.io/gorm"
)

func migrateAddUsers(ctx context.Context, db *gorm.DB) error {
	// Your migration code here
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

var MyMigration = models.Migration{
	MigrationId:         "550e8400-e29b-41d4-a716-446655440000",
	PreviousMigrationID: "", // Empty for the first migration
	MigrateFunc:         migrateAddUsers,
	Description:         "Add users table",
}
```

### Applying Migrations

Apply a migration using the `Apply` method:

```go
ctx := context.Background()

err := MyMigration.Apply(ctx, db)
if err != nil {
	log.Fatalf("Migration failed: %v", err)
}
```

The `Apply` method:
- Checks if the migration should run based on previous migrations
- Executes the migration function if needed
- Records the migration in the database to prevent duplicate runs

### Chaining Migrations

For subsequent migrations, set the `PreviousMigrationID` to the `MigrationId` of the previous migration:

```go
func migrateAddPosts(ctx context.Context, db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			title VARCHAR(255) NOT NULL,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

var SecondMigration = models.Migration{
	MigrationId:         "660e8400-e29b-41d4-a716-446655440001",
	PreviousMigrationID: "550e8400-e29b-41d4-a716-446655440000", // ID of previous migration
	MigrateFunc:         migrateAddPosts,
	Description:         "Add posts table",
}
```

### Running Multiple Migrations

You can run multiple migrations in sequence:

```go
migrations := []models.Migration{
	MyMigration,
	SecondMigration,
	// Add more migrations...
}

ctx := context.Background()
for _, migration := range migrations {
	err := migration.Apply(ctx, db)
	if err != nil {
		log.Fatalf("Migration %s failed: %v", migration.Description, err)
	}
}
```

### Migration Tracking

GORMful uses a `MigrationModel` table to track which migrations have been applied. The migration system ensures:
- Migrations run only once
- Migrations run in the correct order based on `PreviousMigrationID`
- The database state matches the migration history

The `MigrationModel` struct is automatically managed by GORMful - you don't need to create it manually, but you can reference it if needed:

```go
type MigrationModel struct {
	ID string `gorm:"primaryKey"`
}
```

## License

See [LICENSE](LICENSE) file for details.