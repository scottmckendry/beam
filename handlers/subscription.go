package handlers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	db "github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/handlers/utils"
	"github.com/scottmckendry/beam/ui/views"
)

// RegisterSubscriptionRoutes registers all subscription-related routes to the given router.
func (h *Handlers) RegisterSubscriptionRoutes(r chi.Router) {
	r.Get("/sse/customer/{customerID}/subscriptions", h.GetCustomerSubscriptionsSSE)
	r.Get("/sse/customer/{customerID}/add-subscription", h.AddSubscriptionFormSSE)
	r.Get("/sse/customer/{customerID}/add-subscription-submit", h.AddSubscriptionSubmitSSE)
	r.Get("/sse/customer/{customerID}/edit-subscription/{subscriptionID}", h.EditSubscriptionFormSSE)
	r.Get("/sse/customer/{customerID}/edit-subscription-submit/{subscriptionID}", h.EditSubscriptionSubmitSSE)
	r.Get("/sse/customer/{customerID}/delete-subscription/{subscriptionID}", h.DeleteSubscriptionSSE)
}

// AddSubscriptionFormSSE renders the form to add a new subscription for a customer via SSE.
func (h *Handlers) AddSubscriptionFormSSE(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerID")
	if customerID == "" {
		slog.Error("No customerID provided in URL param")
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Missing Customer ID", "No customer ID provided.", w, r)
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.AddSubscription(customerID),
		},
	})
}

// AddSubscriptionSubmitSSE handles the submission of the add subscription form, creates the subscription, and refreshes the subscription list via SSE.
func (h *Handlers) AddSubscriptionSubmitSSE(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerID")
	cid, err := uuid.Parse(customerID)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh subscriptions", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}

	var params db.CreateSubscriptionParams
	if err := utils.MapFormToStruct(r, &params); err != nil {
		slog.Error("Error parsing/mapping form", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Form Error", "An error occurred while processing the form.", w, r)
		return
	}

	// ensure the customer ID is set in the params
	params.CustomerID = cid

	_, err = h.Queries.CreateSubscription(r.Context(), params)
	if err != nil {
		slog.Error("Error adding subscription", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to add subscription", "An error occurred while adding the subscription. Please try again.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Subscription added", "The subscription has been successfully added.", w, r)

	// Refresh the subscription list for the customer
	subscriptions, err := h.Queries.ListSubscriptionsByCustomer(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to list subscriptions for customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh subscriptions", "An error occurred while fetching the updated subscriptions. Please try again.", w, r)
		return
	}
	customer, err := h.Queries.GetCustomer(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh subscriptions", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.CustomerSubscriptions(customer, subscriptions),
		},
	})
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

// EditSubscriptionFormSSE renders the form to edit a subscription for a customer via SSE.
func (h *Handlers) EditSubscriptionFormSSE(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerID")
	subscriptionID := chi.URLParam(r, "subscriptionID")
	if customerID == "" || subscriptionID == "" {
		slog.Error("Missing customerID or subscriptionID in URL param")
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Missing ID", "Customer or subscription ID missing.", w, r)
		return
	}

	sid, err := uuid.Parse(subscriptionID)
	if err != nil {
		slog.Error("Invalid subscriptionID", "subscriptionID", subscriptionID, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid subscription ID", "The subscription ID is invalid.", w, r)
		return
	}

	sub, err := h.Queries.GetSubscription(r.Context(), sid)
	if err != nil {
		slog.Error("Failed to get subscription", "subscriptionID", subscriptionID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to get subscription", "Could not fetch subscription details.", w, r)
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.EditSubscription(customerID, sub),
		},
	})
}

// EditSubscriptionSubmitSSE handles the submission of the edit subscription form, updates the subscription, and refreshes the subscription list via SSE.
func (h *Handlers) EditSubscriptionSubmitSSE(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerID")
	subscriptionID := chi.URLParam(r, "subscriptionID")
	sid, err := uuid.Parse(subscriptionID)
	if err != nil {
		slog.Error("Invalid subscriptionID", "subscriptionID", subscriptionID, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid subscription ID", "The subscription ID is invalid.", w, r)
		return
	}

	var params db.UpdateSubscriptionParams
	if err := utils.MapFormToStruct(r, &params); err != nil {
		slog.Error("Error mapping form to struct", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Form Error", "An error occurred while processing the form.", w, r)
		return
	}
	params.ID = sid

	_, err = h.Queries.UpdateSubscription(r.Context(), params)
	if err != nil {
		slog.Error("Error updating subscription", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to update subscription", "An error occurred while updating the subscription. Please try again.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Subscription updated", "The subscription has been successfully updated.", w, r)

	cid, err := uuid.Parse(customerID)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh subscriptions", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}

	subscriptions, err := h.Queries.ListSubscriptionsByCustomer(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to list subscriptions for customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh subscriptions", "An error occurred while fetching the updated subscriptions. Please try again.", w, r)
		return
	}
	customer, err := h.Queries.GetCustomer(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh subscriptions", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.CustomerSubscriptions(customer, subscriptions),
		},
	})
}

// DeleteSubscriptionSSE handles deleting a subscription and refreshing the subscription list via SSE.
func (h *Handlers) DeleteSubscriptionSSE(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerID")
	subscriptionID := chi.URLParam(r, "subscriptionID")
	if customerID == "" || subscriptionID == "" {
		slog.Error("Missing customerID or subscriptionID in URL param")
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Missing ID", "Customer or subscription ID missing.", w, r)
		return
	}

	sid, err := uuid.Parse(subscriptionID)
	if err != nil {
		slog.Error("Invalid subscriptionID", "subscriptionID", subscriptionID, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid subscription ID", "The subscription ID is invalid.", w, r)
		return
	}

	cid, err := uuid.Parse(customerID)
	if err != nil {
		slog.Error("Invalid customerID", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid customer ID", "The customer ID is invalid.", w, r)
		return
	}

	_, err = h.Queries.DeleteSubscription(r.Context(), sid)
	if err != nil {
		slog.Error("Error deleting subscription", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to delete subscription", "An error occurred while deleting the subscription. Please try again.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Subscription deleted", "The subscription has been successfully deleted.", w, r)

	subscriptions, err := h.Queries.ListSubscriptionsByCustomer(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to list subscriptions for customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh subscriptions", "An error occurred while fetching the updated subscriptions. Please try again.", w, r)
		return
	}
	customer, err := h.Queries.GetCustomer(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh subscriptions", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}
	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.CustomerSubscriptions(customer, subscriptions),
		},
	})
}
