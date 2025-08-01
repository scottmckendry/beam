// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID           uuid.UUID
	CustomerID   uuid.UUID
	ActivityType string
	Action       string
	Description  string
	CreatedAt    sql.NullTime
}

type Contact struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Name       string
	Role       sql.NullString
	Email      sql.NullString
	Phone      sql.NullString
	Avatar     sql.NullString
	IsPrimary  sql.NullBool
	Notes      sql.NullString
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
	DeletedAt  sql.NullTime
}

type Customer struct {
	ID        uuid.UUID
	Name      string
	Logo      sql.NullString
	Status    string
	Email     sql.NullString
	Phone     sql.NullString
	Address   sql.NullString
	Website   sql.NullString
	Notes     sql.NullString
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type Migration struct {
	ID      uuid.UUID
	Name    string
	Applied time.Time
}

type Subscription struct {
	ID             uuid.UUID
	CustomerID     uuid.UUID
	Description    string
	Amount         float64
	Term           string
	BillingCadence string
	StartDate      time.Time
	EndDate        sql.NullTime
	Status         string
	Notes          sql.NullString
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
	DeletedAt      sql.NullTime
}

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	GithubID string
	IsAdmin  bool
}
