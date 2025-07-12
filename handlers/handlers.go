// Package handlers provides the Handlers struct and constructor for HTTP handler functions.
package handlers

import (
	"github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/oauth"
)

type Handlers struct {
	Queries *db.Queries
	OAuth   *oauth.OAuth
}

// New creates a new Handlers instance with the provided database queries and OAuth environment.
func New(queries *db.Queries, env *oauth.OAuth) *Handlers {
	return &Handlers{Queries: queries, OAuth: env}
}
