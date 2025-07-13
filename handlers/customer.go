package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/starfederation/datastar/sdk/go/datastar"

	"github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/ui/views"
)

func (h *Handlers) CustomerNavSSE(w http.ResponseWriter, r *http.Request) {
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		log.Printf("Failed to load customers: %v", err)
		h.Notify(
			NotifyError,
			"View Error",
			"An error occurred while loading the customer navigation.",
			w,
			r,
		)
		http.Error(w, "Failed to load customers", http.StatusInternalServerError)
		return
	}
	currentPage := r.URL.Query().Get("page")
	buf := &bytes.Buffer{}
	views.CustomerNavigation(customers, currentPage).Render(r.Context(), buf)
	ServeSSEElement(w, r, buf.String())
}

func (h *Handlers) AddCustomerSSE(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	views.AddCustomer().Render(r.Context(), buf)
	views.HeaderIcon("customer").Render(r.Context(), buf)
	pageSignals := PageSignals{
		HeaderTitle:       "Add Customer",
		HeaderDescription: "Woohoo! Let's add a new customer 🚀",
		CurrentPage:       "none",
	}
	encodedSignals, _ := json.Marshal(pageSignals)
	sse := datastar.NewSSE(w, r)
	sse.PatchSignals(encodedSignals)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}

func (h *Handlers) GetCustomerSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		views.NotFound().Render(r.Context(), w)
		return
	}
	h.renderCustomerOverviewSSE(w, r, parsedID)
}

// Helper to render customer overview via SSE
func (h *Handlers) renderCustomerOverviewSSE(
	w http.ResponseWriter,
	r *http.Request,
	customerID uuid.UUID,
) {
	c, err := h.Queries.GetCustomer(r.Context(), customerID)
	if err != nil {
		log.Printf("GetCustomer failed for ID=%v: %v", customerID, err)
		h.Notify(NotifyError, "Customer Not Found", "No customer found for the provided ID.", w, r)
		w.WriteHeader(http.StatusNotFound)
		views.NotFound().Render(r.Context(), w)
		return
	}
	buf := &bytes.Buffer{}
	views.Customer(c).Render(r.Context(), buf)
	views.HeaderIcon("customer").Render(r.Context(), buf)
	pageSignals := PageSignals{
		HeaderTitle: c.Name,
		HeaderDescription: fmt.Sprintf(
			"%d %s • %d %s • %d %s",
			c.ContactCount,
			Pluralise(c.ContactCount, "contact", "contacts"),
			c.SubscriptionCount,
			Pluralise(c.SubscriptionCount, "subscription", "subscriptions"),
			c.ProjectCount,
			Pluralise(c.ProjectCount, "project", "projects"),
		),
		CurrentPage: c.ID.String(),
	}
	encodedSignals, _ := json.Marshal(pageSignals)
	sse := datastar.NewSSE(w, r)
	sse.PatchSignals(encodedSignals)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}

func (h *Handlers) SubmitAddCustomerSSE(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	status := r.URL.Query().Get("status")
	notes := r.URL.Query().Get("notes")

	if name == "" || email == "" {
		log.Printf("Missing required fields: name or email")
		h.Notify(NotifyError, "Missing Fields", "Name and email are required.", w, r)
		http.Error(w, "Missing required fields: name and email", http.StatusBadRequest)
		return
	}

	params := db.CreateCustomerParams{
		Name:    name,
		Logo:    sql.NullString{},
		Status:  status,
		Email:   sql.NullString{String: email, Valid: email != ""},
		Phone:   sql.NullString{},
		Address: sql.NullString{},
		Website: sql.NullString{},
		Notes:   sql.NullString{String: notes, Valid: notes != ""},
	}

	customer, err := h.Queries.CreateCustomer(r.Context(), params)
	if err != nil {
		log.Printf("Error adding customer: %v", err)
		h.Notify(NotifyError, "Add Failed", "An error occurred while adding the customer.", w, r)
		http.Error(w, "Failed to add customer", http.StatusInternalServerError)
		return
	}

	h.Notify(NotifySuccess, "Customer Added", "Customer has been successfully added.", w, r)
	h.renderCustomerOverviewSSE(w, r, customer.ID)

	// refresh the customer navigation
	h.CustomerNavSSE(w, r)
}

func (h *Handlers) EditCustomerFormSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	c, err := h.Queries.GetCustomer(r.Context(), parsedID)
	if err != nil {
		log.Printf("GetCustomer failed for ID=%v: %v", parsedID, err)
		h.Notify(NotifyError, "Customer Not Found", "No customer found for the provided ID.", w, r)
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	buf := &bytes.Buffer{}
	views.EditCustomer(c).Render(r.Context(), buf)
	views.HeaderIcon("customer").Render(r.Context(), buf)

	pageSignals := PageSignals{
		HeaderTitle:       "Edit Customer",
		HeaderDescription: fmt.Sprintf("Editing %s", c.Name),
		CurrentPage:       c.ID.String(),
	}
	encodedSignals, _ := json.Marshal(pageSignals)

	sse := datastar.NewSSE(w, r)
	sse.PatchSignals(encodedSignals)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}

func (h *Handlers) EditCustomerSubmitSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	status := r.URL.Query().Get("status")
	notes := r.URL.Query().Get("notes")

	if name == "" || email == "" {
		log.Printf("Missing required fields: name or email")
		h.Notify(NotifyError, "Missing Fields", "Name and email are required.", w, r)
		http.Error(w, "Missing required fields: name and email", http.StatusBadRequest)
		return
	}

	params := db.UpdateCustomerParams{
		ID:     parsedID,
		Name:   name,
		Status: status,
		Email:  sql.NullString{String: email, Valid: email != ""},
		Notes:  sql.NullString{String: notes, Valid: notes != ""},
	}

	_, err = h.Queries.UpdateCustomer(r.Context(), params)
	if err != nil {
		log.Printf("Error updating customer: %v", err)
		h.Notify(
			NotifyError,
			"Update Failed",
			"An error occurred while updating the customer.",
			w,
			r,
		)
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	h.Notify(
		NotifySuccess,
		"Customer Updated",
		fmt.Sprintf("%s has been successfully updated.", name),
		w,
		r,
	)
	h.renderCustomerOverviewSSE(w, r, parsedID)
	// refresh the customer navigation
	h.CustomerNavSSE(w, r)
}

func (h *Handlers) DeleteCustomerSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid Customer ID", "The customer ID provided is not valid.", w, r)
		return
	}

	c, err := h.Queries.DeleteCustomer(r.Context(), parsedID)
	if err != nil {
		log.Printf("Error deleting customer: %v", err)
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		h.Notify(
			NotifyError,
			"Delete Failed",
			"An error occurred while trying to delete the customer.",
			w,
			r,
		)
		return
	}

	// render dashboard, refresh customer navigation
	h.Notify(
		NotifySuccess,
		"Customer Deleted",
		fmt.Sprintf("Customer %s has been successfully deleted.", c.Name),
		w,
		r,
	)
	h.DashboardSSE(w, r)
	h.CustomerNavSSE(w, r)
}
