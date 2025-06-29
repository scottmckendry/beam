package db

import (
	"context"
	"os"
	"testing"
)

func TestInitialiseDB(t *testing.T) {
	os.Remove("data/beam.db")
	os.MkdirAll("data", 0755) // Ensure the data directory exists
	db, queries, err := InitialiseDB()
	if err != nil {
		t.Fatalf("InitialiseDB failed: %v", err)
	}
	defer db.Close()
	if queries == nil {
		t.Fatal("queries should not be nil")
	}
}

func TestApplyMigrations_Idempotent(t *testing.T) {
	db, queries, err := InitialiseDB()
	if err != nil {
		t.Fatalf("InitialiseDB failed: %v", err)
	}
	defer db.Close()
	err = applyMigrations(context.Background(), db, queries)
	if err != nil {
		t.Fatalf("applyMigrations should be idempotent: %v", err)
	}
}

func TestSplitSQLStatements(t *testing.T) {
	sql := "CREATE TABLE test (id INTEGER); INSERT INTO test VALUES ('a;bc');"
	stmts := splitSQLStatements(sql)
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
}
