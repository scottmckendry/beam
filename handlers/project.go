package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/scottmckendry/beam/handlers/utils"
	"github.com/scottmckendry/beam/ui/views"
)

// RegisterProjectRoutes registers all customer-related routes to the given router.
func (h *Handlers) RegisterProjectRoutes(r chi.Router) {
	r.Get("/sse/customer/{customerID}/projects", h.GetCustomerProjectsSSE)
}

// GetCustomerProjectsSSE retrieves a customer's projects by ID and renders them via SSE
func (h *Handlers) GetCustomerProjectsSSE(w http.ResponseWriter, r *http.Request) {
	c, ok := h.getCustomerByID(w, r, "customerID")
	if !ok {
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.CustomerProjects(c),
			views.HeaderIcon("customer"),
		},
	})
}
