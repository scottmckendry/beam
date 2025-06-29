// Package db handles database initialization, migrations, and utility functions for interacting with the application's persistent storage.
package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/tursodatabase/go-libsql"

	"github.com/scottmckendry/beam/db/sqlc"
)

//go:embed migrations
var migrationsFS embed.FS

const dbName = "file:data/beam.db"

// InitialiseDB sets up the database and runs migrations. Returns the DB and Queries for use/testing.
func InitialiseDB() (*sql.DB, *db.Queries, error) {
	ctx := context.Background()

	store, err := sql.Open("libsql", dbName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	queries := db.New(store)
	if err := queries.CreateMigrationsTable(ctx); err != nil {
		store.Close()
		return nil, nil, fmt.Errorf("failed to create migrations table: %w", err)
	}

	if err := applyMigrations(ctx, store, queries); err != nil {
		store.Close()
		return nil, nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return store, queries, nil
}

// applyMigrations applies all pending migrations in order.
func applyMigrations(ctx context.Context, db *sql.DB, queries *db.Queries) error {
	fileNames, err := getMigrationFileNames()
	if err != nil {
		return err
	}
	for _, fileName := range fileNames {
		applied, err := isMigrationApplied(ctx, queries, fileName)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := applySingleMigration(ctx, db, queries, fileName); err != nil {
			return err
		}
	}
	return nil
}

// getMigrationFileNames retrieves the list of migration files from the embedded filesystem.
func getMigrationFileNames() ([]string, error) {
	files, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("error reading migrations directory: %v", err)
	}
	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			fileNames = append(fileNames, file.Name())
		}
	}
	sort.Strings(fileNames)
	return fileNames, nil
}

// isMigrationApplied checks if a migration has already been applied by querying the migrations table.
func isMigrationApplied(ctx context.Context, queries *db.Queries, fileName string) (bool, error) {
	_, err := queries.GetMigration(ctx, fileName)
	if err == nil {
		return true, nil
	}
	if err != sql.ErrNoRows {
		return false, fmt.Errorf("error checking migration status %s: %v", fileName, err)
	}
	return false, nil
}

// applySingleMigration applies a single migration file to the database.
func applySingleMigration(
	ctx context.Context,
	db *sql.DB,
	queries *db.Queries,
	fileName string,
) error {
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
	stmts := splitSQLStatements(string(content))
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if len(stmt) == 0 {
			continue
		}
		_, err = tx.ExecContext(ctx, stmt)
		if err != nil {
			return fmt.Errorf("error executing migration %s: %v", fileName, err)
		}
	}
	err = qtx.ApplyMigration(ctx, fileName)
	if err != nil {
		return fmt.Errorf("error recording migration %s: %v", fileName, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing migration %s: %v", fileName, err)
	}
	return nil
}

// splitSQLStatements splits SQL migration files into individual statements by semicolon.
func splitSQLStatements(sql string) []string {
	stmts := []string{}
	curr := ""
	inString := false
	for _, r := range sql {
		if r == '\'' {
			inString = !inString
		}
		if r == ';' && !inString {
			stmts = append(stmts, curr)
			curr = ""
		} else {
			curr += string(r)
		}
	}
	if len(curr) > 0 {
		stmts = append(stmts, curr)
	}
	return stmts
}
