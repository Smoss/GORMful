package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	migrationModels "github.com/GORMful/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	createMigrationFlags := flag.NewFlagSet("create-migration", flag.ExitOnError)
	createMigrationDescription := createMigrationFlags.String("description", "", "Description for the migration")
	migrationOutputDirectory := createMigrationFlags.String("output-directory", "migrations/migrations", "Output directory for the migration")

	dbHost := flag.String("db-host", os.Getenv("DB_HOST"), "Host for the database")
	dbUser := flag.String("db-user", os.Getenv("DB_USER"), "User for the database")
	dbPassword := flag.String("db-password", os.Getenv("DB_PASSWORD"), "Password for the database")
	dbPort := flag.String("db-port", os.Getenv("DB_PORT"), "Port for the database")
	dbName := flag.String("db-name", os.Getenv("DB_NAME"), "Name for the database")
	sslMode := flag.String("ssl-mode", "disable", "SSL mode for the database")
	flag.Parse()

	db, err := getDB(*dbHost, *dbUser, *dbPassword, *dbName, *dbPort, *sslMode)
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}

	switch os.Args[1] {
	case "create-migration":
		CreateMigrationFile(createMigrationDescription, db, migrationOutputDirectory)
	default:
		log.Fatalf("Invalid command: %s", os.Args[1])
	}
}

func CreateMigrationFile(description *string, db *gorm.DB, migrationOutputDirectory *string) {
	// Get the current date in the format YYYYMMDDHHMMSS
	now := time.Now()
	migrationName := now.Format("20060102150405")
	fileName := fmt.Sprintf("%s/migrate_%s.go", *migrationOutputDirectory, migrationName)

	previousMigrationId := ""
	ctx := context.Background()
	previousMigration, err := gorm.G[migrationModels.MigrationModel](db, gorm.WithResult()).Last(ctx)
	if err != nil {
		previousMigrationId = ""
	} else {
		previousMigrationId = previousMigration.ID
	}

	template := string(migrationModels.MigrationTemplate)
	template = strings.Replace(template, "<migrationName>", migrationName, -1)
	template = strings.Replace(template, "<migrationId>", uuid.New().String(), -1)
	template = strings.Replace(template, "<previousMigrationId>", previousMigrationId, -1)
	template = strings.Replace(template, "<description>", *description, -1)
	os.WriteFile(fileName, []byte(template), 0644)
}

func getDB(dbHost string, dbUser string, dbPassword string, dbName string, dbPort string, sslMode string) (*gorm.DB, error) {
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" || sslMode == "" {
		log.Fatal("DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, SSL_MODE are required")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPassword, dbName, dbPort, sslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
