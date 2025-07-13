package handlers

import (
	"bytes"
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

// AddCustomerSSE renders the form to add a new customer via SSE
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

// GetCustomerSSE retrieves a customer by ID and renders the overview page via SSE
func (h *Handlers) GetCustomerSSE(w http.ResponseWriter, r *http.Request) {
	c, ok := h.getCustomerByID(w, r, "id")
	if !ok {
		return
	}
	buf := &bytes.Buffer{}
	views.Customer(c).Render(r.Context(), buf)
	views.HeaderIcon("customer").Render(r.Context(), buf)
	pageSignals := buildCustomerPageSignals(c)
	sse := datastar.NewSSE(w, r)
	sse.PatchSignals(pageSignals)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}

// SubmitAddCustomerSSE handles the submission of the add customer form, upon success it will render the customer overview page and refresh the customer navigation
func (h *Handlers) SubmitAddCustomerSSE(w http.ResponseWriter, r *http.Request) {
	var params db.CreateCustomerParams
	if err := MapFormToStruct(r, &params); err != nil {
		log.Printf("Error mapping form to struct: %v", err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		h.Notify(NotifyError, "Form Error", "An error occurred while processing the form.", w, r)
		return
	}

	customer, err := h.Queries.CreateCustomer(r.Context(), params)
	if err != nil {
		log.Printf("Error adding customer: %v", err)
		h.Notify(NotifyError, "Add Failed", "An error occurred while adding the customer.", w, r)
		http.Error(w, "Failed to add customer", http.StatusInternalServerError)
		return
	}

	h.Notify(NotifySuccess, "Customer Added", "Customer has been successfully added.", w, r)
	h.logActivity(
		r,
		customer.ID,
		"customer",
		"customer_created",
		fmt.Sprintf("New customer \"%s\" created", customer.Name),
	)

	c, err := h.Queries.GetCustomer(r.Context(), customer.ID)
	if err != nil {
		log.Printf("GetCustomer failed for ID=%v: %v", customer.ID, err)
		h.Notify(NotifyError, "Customer Not Found", "No customer found for the provided ID.", w, r)
		w.WriteHeader(http.StatusNotFound)
		views.NotFound().Render(r.Context(), w)
		return
	}

	buf := &bytes.Buffer{}
	views.Customer(c).Render(r.Context(), buf)
	views.HeaderIcon("customer").Render(r.Context(), buf)

	// update the navigation
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		log.Printf("Failed to load customers for navigation: %v", err)
		h.Notify(NotifyError, "Navigation Error", "An error occurred while loading the customer navigation.", w, r)
	}
	views.CustomerNavigation(customers).Render(r.Context(), buf)

	sse := datastar.NewSSE(w, r)

	pageSignals := buildCustomerPageSignals(c)
	sse.PatchSignals(pageSignals)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
		// TODO: Morph should be working here, but for whatever reason it isn't.
		// Hence, we use replace to ensure that the customer navigation is updated correctly
		// revisit this later once the following issue is resolved:
		// https://github.com/starfederation/datastar/issues/999
		datastar.WithModeReplace(),
	)
}

// EditCustomerFormSSE renders the form to edit an existing customer via SSE
func (h *Handlers) EditCustomerFormSSE(w http.ResponseWriter, r *http.Request) {
	c, ok := h.getCustomerByID(w, r, "id")
	if !ok {
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

// EditCustomerSubmitSSE handles the submission of the edit customer, upon success it will render the customer overview page and refresh the customer navigation
func (h *Handlers) EditCustomerSubmitSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	var params db.UpdateCustomerParams
	if err := MapFormToStruct(r, &params); err != nil {
		log.Printf("Error mapping form to struct: %v", err)
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		h.Notify(NotifyError, "Form Error", "An error occurred while processing the form.", w, r)
		return
	}

	_, err = h.Queries.UpdateCustomer(r.Context(), params)
	if err != nil {
		log.Printf("Error updating customer: %v", err)
		h.Notify(NotifyError, "Update Failed", "An error occurred while updating the customer.", w, r)
		http.Error(w, "Failed to update customer", http.StatusInternalServerError)
		return
	}

	h.Notify(NotifySuccess, "Customer Updated", fmt.Sprintf("%s has been successfully updated.", params.Name), w, r)
	h.logActivity(r, parsedID, "customer", "customer_updated", fmt.Sprintf("Customer %s updated", params.Name))

	buf := &bytes.Buffer{}
	c, _ := h.getCustomerByID(w, r, "id")
	views.Customer(c).Render(r.Context(), buf)
	views.HeaderIcon("customer").Render(r.Context(), buf)

	// update the navigation
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		log.Printf("Failed to load customers for navigation: %v", err)
		h.Notify(NotifyError, "Navigation Error", "An error occurred while loading the customer navigation.", w, r)
	}
	views.CustomerNavigation(customers).Render(r.Context(), buf)

	sse := datastar.NewSSE(w, r)
	pageSignals := buildCustomerPageSignals(c)
	sse.PatchSignals(pageSignals)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
		datastar.WithModeReplace(),
	)
}

// DeleteCustomerSSE handles the deletion of a customer - if successful, it will render the dashboard and refresh the customer navigation
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
		h.Notify(NotifyError, "Delete Failed", "An error occurred while trying to delete the customer.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Customer Deleted", fmt.Sprintf("Customer %s has been successfully deleted.", c.Name), w, r)
	// INFO: this will fail while we still have delete cascade constraints in place - see TODO in the the init migration
	h.logActivity(r, c.ID, "customer", "customer_deleted", fmt.Sprintf("Customer %s deleted", c.Name))

	// render dashboard, refresh customer navigation
	h.DashboardSSE(w, r)

	buf := &bytes.Buffer{}
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		log.Printf("Failed to load customers for navigation: %v", err)
		h.Notify(NotifyError, "Navigation Error", "An error occurred while loading the customer navigation.", w, r)
	}
	views.CustomerNavigation(customers).Render(r.Context(), buf)
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
		datastar.WithModeReplace(),
	)
}

// getCustomerByID fetches a customer by ID from the URL param and handles errors consistently
func (h *Handlers) getCustomerByID(
	w http.ResponseWriter,
	r *http.Request,
	idParam string,
) (db.GetCustomerRow, bool) {
	id := chi.URLParam(r, idParam)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		h.Notify(NotifyError, "Invalid Customer ID", "The customer ID provided is not valid.", w, r)
		w.WriteHeader(http.StatusBadRequest)
		views.NotFound().Render(r.Context(), w)
		return db.GetCustomerRow{}, false
	}
	c, err := h.Queries.GetCustomer(r.Context(), parsedID)
	if err != nil {
		log.Printf("GetCustomer failed for ID=%v: %v", parsedID, err)
		h.Notify(NotifyError, "Customer Not Found", "No customer found for the provided ID.", w, r)
		w.WriteHeader(http.StatusNotFound)
		views.NotFound().Render(r.Context(), w)
		return db.GetCustomerRow{}, false
	}
	return c, true
}

// buildCustomerPageSignals constructs the page signals for a customer overview
func buildCustomerPageSignals(c db.GetCustomerRow) []byte {
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
	return encodedSignals
}
