package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	middlewares "github.com/scottmckendry/beam/middleware"
	"github.com/scottmckendry/beam/ui/views"
)

// RegisterRootRoutes registers the admin root route on the given router.
func (h *Handlers) RegisterRootRoutes(r chi.Router) {
	r.Get("/", h.HandleRoot)
}

// HandleRoot serves the main application page.
func (h *Handlers) HandleRoot(w http.ResponseWriter, r *http.Request) {
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		slog.Error("Failed to load customers", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		views.ServerError().Render(r.Context(), w)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserKey).(string)
	if !ok || userID == "" {
		slog.Error("No user in context", "userID", userID)
		w.WriteHeader(http.StatusInternalServerError)
		views.ServerError().Render(r.Context(), w)
		return
	}

	user, err := h.Queries.GetUserByGithubID(r.Context(), userID)
	if err != nil {
		slog.Error("Failed to get user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		views.ServerError().Render(r.Context(), w)
		return
	}
	// Admin middleware handles the admin check, so we can assume the user is authenticated here and has admin privileges.
	views.Root(true, customers, user).Render(r.Context(), w)
}
