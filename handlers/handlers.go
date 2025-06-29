// Package handlers provides HTTP handler functions for user authentication, session management, and main application routes.
package handlers

import (
	"net/http"

	"github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/github"
	"github.com/scottmckendry/beam/oauth"
	"github.com/scottmckendry/beam/ui/views"
)

type Handlers struct {
	Queries *db.Queries
}

// New creates a new Handlers instance with the provided database queries.
func New(queries *db.Queries) *Handlers {
	return &Handlers{Queries: queries}
}

// HandleLogin processes GET requests to the login page and redirects authenticated users to the root.
func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("user_name")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	views.Login().Render(r.Context(), w)
}

// HandleLogout clears the user session and redirects to the login page.
func (h *Handlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(
		w,
		&http.Cookie{Name: "user_name", Value: "", Path: "/", HttpOnly: true, MaxAge: -1},
	)
	http.Redirect(w, r, "/login", http.StatusFound)
}

// HandleRoot serves the main application page, handling user authentication and repository info display.
func (h *Handlers) HandleRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := oauth.GetSignedCookie(r, "user_name")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	tokenCookie, err := r.Cookie("oauth_token")
	if err != nil {
		views.Root(user, false, "", "", 0, 0).Render(ctx, w)
		return
	}
	ghClient := github.NewClient(tokenCookie.Value)
	repo, err := ghClient.GetRepo(user, "beam")
	if err != nil {
		views.Root(user, false, "", "", 0, 0).Render(ctx, w)
		return
	}
	isAdmin, err := h.Queries.IsUserAdmin(ctx, user)
	if err != nil {
		views.Root(user, false, repo.FullName, repo.Description, repo.StargazersCount, repo.ForksCount).
			Render(ctx, w)
		return
	}
	views.Root(user, isAdmin, repo.FullName, repo.Description, repo.StargazersCount, repo.ForksCount).
		Render(ctx, w)
}
