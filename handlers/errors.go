package handlers

import (
	"net/http"

	"github.com/scottmckendry/beam/ui/views"
)

// HandleNotFound serves a 404 Not Found page when a requested resource is not found.
func (h *Handlers) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	views.NotFound().Render(r.Context(), w)
}

// HandleNoAccess serves a page indicating that the user does not have access to the requested resource.
func (h *Handlers) HandleNoAccess(w http.ResponseWriter, r *http.Request) {
	views.NonAdmin().Render(r.Context(), w)
}
