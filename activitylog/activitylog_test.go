package activitylog

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/scottmckendry/beam/db"
	sqlc "github.com/scottmckendry/beam/db/sqlc"
)

func setupTestDB(t *testing.T) (*sqlc.Queries, func()) {
	os.MkdirAll("data", 0755)
	dbConn, queries, err := db.InitialiseDB()
	if err != nil {
		t.Fatalf("InitialiseDB failed: %v", err)
	}
	cleanup := func() { dbConn.Close() }
	return queries, cleanup
}

func createTestCustomer(t *testing.T, queries *sqlc.Queries) sqlc.Customer {
	params := sqlc.CreateCustomerParams{
		Name:   "Test Customer",
		Logo:   sql.NullString{},
		Status: "active",
	}
	customer, err := queries.CreateCustomer(context.Background(), params)
	if err != nil {
		t.Fatalf("CreateCustomer failed: %v", err)
	}
	return customer
}

func TestLogCustomerCreated_Integration(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()
	customer := createTestCustomer(t, queries)
	LogCustomerCreated(context.Background(), queries, customer)
	logs, err := queries.ListRecentActivity(context.Background())
	if err != nil {
		t.Fatalf("ListRecentActivity failed: %v", err)
	}
	found := false
	for _, log := range logs {
		if log.CustomerID == customer.ID && log.Action == "customer_created" {
			found = true
			break
		}
	}
	if !found {
		t.Error("customer_created activity not found in log")
	}
}
