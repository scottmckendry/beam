// Package handlers provides HTTP handler functions for user authentication, session management, and main application routes.
package handlers

import (
	"bytes"
	"log"
	"net/http"

	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/oauth"
	"github.com/scottmckendry/beam/ui/views"
)

type Handlers struct {
	Queries *db.Queries
	OAuth   *oauth.OAuth
}

// HandleHelloSSE streams a message and a signal using Datastar SSE.
func (h *Handlers) HandleHelloSSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	sse.MergeFragments(`<div id="hello-message">Hello from the backend via SSE!</div>`)
	sse.MergeSignals([]byte(`{"hello": "world"}`))
}

// HandleNav streams page content and a signal for navigation.
func (h *Handlers) HandleNav(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "dashboard"
	}
	var buf bytes.Buffer
	switch page {
	case "dashboard":
		views.Dashboard().Render(r.Context(), &buf)
	case "invoices":
		views.Invoices().Render(r.Context(), &buf)
	default:
		buf.WriteString(
			`<div id="page-content" data-fragment="main-content"><p>Page not found.</p></div>`,
		)
	}
	sse.MergeFragments(buf.String())
	sse.MergeSignals([]byte(`{"page": "` + page + `"}`))
}

// New creates a new Handlers instance with the provided database queries.
func New(queries *db.Queries, env *oauth.OAuth) *Handlers {
	return &Handlers{Queries: queries, OAuth: env}
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
	user, err := h.OAuth.GetSignedCookie(r, "user_name")
	if err != nil || user == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	isAdmin, err := h.Queries.IsUserAdmin(ctx, user)
	if err != nil {
		log.Printf("Error checking admin status for user %s: %v", user, err)
		views.Root(false).Render(r.Context(), w)
		return
	}

	views.Root(isAdmin).Render(r.Context(), w)
}
