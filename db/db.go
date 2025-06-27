package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"sort"

	_ "github.com/tursodatabase/go-libsql"

	"github.com/scottmckendry/beam/db/sqlc"
)

//go:embed migrations
var migrationsFS embed.FS

const dbName = "file:data/beam.db"

func InitialiseDb() error {
	ctx := context.Background()

	store, err := sql.Open("libsql", dbName)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return err
	}
	defer store.Close()

	queries := db.New(store)
	err = queries.CreateMigrationsTable(ctx)
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
		return err
	}

	if err := applyMigrations(ctx, store, queries); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
		return err
	}

	return nil
}

func applyMigrations(ctx context.Context, db *sql.DB, queries *db.Queries) error {
	files, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	// Convert to slice for sorting
	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			fileNames = append(fileNames, file.Name())
		}
	}
	sort.Strings(fileNames)

	for _, fileName := range fileNames {
		// skip any already applied migrations
		_, err := queries.GetMigration(ctx, fileName)
		if err == nil {
			continue
		}
		if err != sql.ErrNoRows {
			return fmt.Errorf("error checking migration status %s: %v", fileName, err)
		}

		log.Printf("Applying migration: %s", fileName)

		content, err := migrationsFS.ReadFile(filepath.Join("migrations", fileName))
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", fileName, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("error beginning transaction: %v", err)
		}
		defer tx.Rollback()

		qtx := queries.WithTx(tx)

		// Execute the migration file
		_, err = tx.ExecContext(ctx, string(content))
		if err != nil {
			return fmt.Errorf("error executing migration %s: %v", fileName, err)
		}

		// Record the migration
		err = qtx.ApplyMigration(ctx, fileName)
		if err != nil {
			return fmt.Errorf("error recording migration %s: %v", fileName, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration %s: %v", fileName, err)
		}
	}

	return nil
}
