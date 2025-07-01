// Package handlers provides HTTP handler functions for user authentication, session management, and main application routes.
package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/oauth"
	"github.com/scottmckendry/beam/ui/views"
)

type Handlers struct {
	Queries *db.Queries
	OAuth   *oauth.OAuth
}

// HandleDashboard serves the dashboard page.
func (h *Handlers) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	views.Dashboard().Render(r.Context(), w)
}

// HandleInvoices serves the invoices page.
func (h *Handlers) HandleInvoices(w http.ResponseWriter, r *http.Request) {
	views.Invoices().Render(r.Context(), w)
}

// HandleNoAccess serves a page indicating that the user does not have access to the requested resource.
func (h *Handlers) HandleNoAccess(w http.ResponseWriter, r *http.Request) {
	views.NonAdmin().Render(r.Context(), w)
}

// New creates a new Handlers instance with the provided database queries.
func New(queries *db.Queries, env *oauth.OAuth) *Handlers {
	return &Handlers{Queries: queries, OAuth: env}
}

// HandleCustomers serves the customers page for a specific customer.
func (h *Handlers) HandleCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// TODO: Fetch customer data from DB using id
	var name string
	switch id {
	case "1":
		name = "Acme Corporation"
	case "2":
		name = "Tech Solutions Inc"
	case "3":
		name = "Startup Co"
	default:
		name = "Unknown Customer"
	}
	views.Customer(id, name).Render(r.Context(), w)
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

// HandleRoot serves the main application page.
func (h *Handlers) HandleRoot(w http.ResponseWriter, r *http.Request) {
	// Admin middleware handles the admin check, so we can assume the user is authenticated here and has admin privileges.
	views.Root(true).Render(r.Context(), w)
}
