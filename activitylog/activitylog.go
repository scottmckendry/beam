// Package activitylog centralizes all activity log logic.
package activitylog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/scottmckendry/beam/db/sqlc"
)

type ActivityType string

const (
	ActivityTypeCustomer ActivityType = "customer"
	ActivityTypeContact  ActivityType = "contact"
)

// LogCustomerCreated logs a customer creation event.
func LogCustomerCreated(ctx context.Context, queries *db.Queries, customer db.Customer) {
	logActivity(ctx, queries, customer.ID, ActivityTypeCustomer, "customer_created", fmt.Sprintf("Customer %s created", customer.Name))
}

// LogCustomerUpdated logs a customer update event.
func LogCustomerUpdated(ctx context.Context, queries *db.Queries, customer db.GetCustomerRow) {
	logActivity(ctx, queries, customer.ID, ActivityTypeCustomer, "customer_updated", fmt.Sprintf("Customer %s updated", customer.Name))
}

// LogCustomerDeleted logs a customer deletion event.
func LogCustomerDeleted(ctx context.Context, queries *db.Queries, customer db.Customer) {
	logActivity(ctx, queries, customer.ID, ActivityTypeCustomer, "customer_deleted", fmt.Sprintf("Customer %s deleted", customer.Name))
}

// LogContactAdded logs a contact creation event.
func LogContactAdded(ctx context.Context, queries *db.Queries, customerID uuid.UUID, contactName string) {
	logActivity(ctx, queries, customerID, ActivityTypeContact, "contact_added", fmt.Sprintf("Contact %s added", contactName))
}

// LogContactUpdated logs a contact update event.
func LogContactUpdated(ctx context.Context, queries *db.Queries, customerID uuid.UUID, contactName string) {
	logActivity(ctx, queries, customerID, ActivityTypeContact, "contact_updated", fmt.Sprintf("Contact %s updated", contactName))
}

// LogContactDeleted logs a contact deletion event.
func LogContactDeleted(ctx context.Context, queries *db.Queries, customerID uuid.UUID, contactName string) {
	logActivity(ctx, queries, customerID, ActivityTypeContact, "contact_deleted", fmt.Sprintf("Contact %s deleted", contactName))
}

// logActivity inserts a new activity log entry for a customer or contact
func logActivity(ctx context.Context, queries *db.Queries, customerID uuid.UUID, activityType ActivityType, action, description string) {
	activity := db.LogActivityParams{
		CustomerID:   customerID,
		ActivityType: string(activityType),
		Action:       action,
		Description:  description,
	}
	_, err := queries.LogActivity(ctx, activity)
	if err != nil {
		slog.Error("Failed to log activity", "err", err)
	}
}
