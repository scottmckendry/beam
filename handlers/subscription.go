package handlers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/scottmckendry/beam/handlers/utils"
	"github.com/scottmckendry/beam/ui/views"
)

// RegisterSubscriptionRoutes registers all subscription-related routes to the given router.
func (h *Handlers) RegisterSubscriptionRoutes(r chi.Router) {
	r.Get("/sse/customer/{customerID}/subscriptions", h.GetCustomerSubscriptionsSSE)
}

// GetCustomerSubscriptionsSSE retrieves a customer's subscriptions by ID and renders them via SSE
func (h *Handlers) GetCustomerSubscriptionsSSE(w http.ResponseWriter, r *http.Request) {
	c, ok := h.getCustomerByID(w, r, "customerID")
	if !ok {
		return
	}

	subscriptions, err := h.Queries.ListSubscriptionsByCustomer(r.Context(), c.ID)
	if err != nil {
		slog.Error("ListSubscriptionsByCustomer failed", "customer_id", c.ID, "err", err)
		h.Notify(NotifyError, "Subscriptions Not Found", "No subscriptions found for the provided customer ID.", w, r)
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.CustomerSubscriptions(c, subscriptions),
			views.HeaderIcon("customer"),
		},
	})
}
