package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	al "github.com/scottmckendry/beam/activitylog"
	"github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/handlers/utils"
	"github.com/scottmckendry/beam/ui/views"
)

// RegisterCustomerRoutes registers all customer-related routes to the given router.
func (h *Handlers) RegisterCustomerRoutes(r chi.Router) {
	r.Get("/sse/customer/{id}", h.GetCustomerSSE)
	r.Get("/sse/customer/{customerID}/overview", h.GetCustomerOverviewSSE)
	r.Get("/sse/customer/add", h.AddCustomerSSE)
	r.Get("/sse/customer/add-submit", h.SubmitAddCustomerSSE)
	r.Get("/sse/customer/delete/{id}", h.DeleteCustomerSSE)
	r.Get("/sse/customer/edit/{id}", h.EditCustomerFormSSE)
	r.Get("/sse/customer/edit-submit/{id}", h.EditCustomerSubmitSSE)
	r.Post("/sse/customer/upload-logo/{id}", h.UploadCustomerLogoSSE)
	r.Get("/sse/customer/delete-logo/{id}", h.DeleteCustomerLogoSSE)
}

// AddCustomerSSE renders the form to add a new customer via SSE
func (h *Handlers) AddCustomerSSE(w http.ResponseWriter, r *http.Request) {
	pageSignals := utils.PageSignals{
		HeaderTitle:       "Add Customer",
		HeaderDescription: "Woohoo! Let's add a new customer 🚀",
		CurrentPage:       "none",
	}
	encodedSignals, _ := json.Marshal(pageSignals)

	utils.RenderSSE(w, r, utils.SSEOpts{
		Signals: encodedSignals,
		Views: []templ.Component{
			views.AddCustomer(),
			views.HeaderIcon("customer"),
		},
	})
}

// GetCustomerSSE retrieves a customer by ID and renders the overview page via SSE
func (h *Handlers) GetCustomerSSE(w http.ResponseWriter, r *http.Request) {
	c, ok := h.getCustomerByID(w, r, "id")
	if !ok {
		return
	}
	pageSignals := buildCustomerPageSignals(c)

	utils.RenderSSE(w, r, utils.SSEOpts{
		Signals: pageSignals,
		Views: []templ.Component{
			views.Customer(c),
			views.HeaderIcon("customer"),
		},
	})
}

// SubmitAddCustomerSSE handles the submission of the add customer form, upon success it will render the customer overview page and refresh the customer navigation
func (h *Handlers) SubmitAddCustomerSSE(w http.ResponseWriter, r *http.Request) {
	var params db.CreateCustomerParams
	if err := utils.MapFormToStruct(r, &params); err != nil {
		slog.Error("Error mapping form to struct", "err", err)
		h.Notify(NotifyError, "Form Error", "An error occurred while processing the form.", w, r)
		return
	}

	customer, err := h.Queries.CreateCustomer(r.Context(), params)
	if err != nil {
		slog.Error("Error adding customer", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Add Failed", "An error occurred while adding the customer.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Customer Added", "Customer has been successfully added.", w, r)
	al.LogCustomerCreated(r.Context(), h.Queries, customer)

	c, err := h.Queries.GetCustomer(r.Context(), customer.ID)
	if err != nil {
		slog.Error("GetCustomer failed", "customer_id", customer.ID, "err", err)
		h.Notify(NotifyError, "Customer Not Found", "No customer found for the provided ID.", w, r)
		views.NotFound().Render(r.Context(), w)
		return
	}

	// get the updated customer list for navigation
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		slog.Error("Failed to load customers for navigation", "err", err)
		h.Notify(NotifyError, "Navigation Error", "An error occurred while loading the customer navigation.", w, r)
	}

	// render the customer overview page with the new customer, along with a navigation refresh
	utils.RenderSSE(w, r, utils.SSEOpts{
		Signals: buildCustomerPageSignals(c),
		Views: []templ.Component{
			views.Customer(c),
			views.HeaderIcon("customer"),
			views.CustomerNavigation(customers),
		},
	})
}

// EditCustomerFormSSE renders the form to edit an existing customer via SSE
func (h *Handlers) EditCustomerFormSSE(w http.ResponseWriter, r *http.Request) {
	c, ok := h.getCustomerByID(w, r, "id")
	if !ok {
		return
	}

	pageSignals := utils.PageSignals{
		HeaderTitle:       "Edit Customer",
		HeaderDescription: fmt.Sprintf("Editing %s", c.Name),
		CurrentPage:       c.ID.String(),
	}
	encodedSignals, _ := json.Marshal(pageSignals)

	utils.RenderSSE(w, r, utils.SSEOpts{
		Signals: encodedSignals,
		Views: []templ.Component{
			views.EditCustomer(c),
			views.HeaderIcon("customer"),
		},
	})
}

// EditCustomerSubmitSSE handles the submission of the edit customer, upon success it will render the customer overview page and refresh the customer navigation
func (h *Handlers) EditCustomerSubmitSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("Invalid customer ID", "err", err)
		h.Notify(NotifyError, "Invalid Customer ID", "The customer ID provided is not valid.", w, r)
		return
	}

	// prevent non-form fields from being overwritten
	var params db.UpdateCustomerParams
	c, _ := h.getCustomerByID(w, r, "id")
	params.ID = c.ID
	params.Logo = c.Logo

	_, err = h.Queries.UpdateCustomer(r.Context(), params)
	if err != nil {
		slog.Error("Error updating customer", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Update Failed", "An error occurred while updating the customer.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Customer Updated", fmt.Sprintf("%s has been successfully updated.", params.Name), w, r)
	c, _ = h.Queries.GetCustomer(r.Context(), parsedID)
	al.LogCustomerUpdated(r.Context(), h.Queries, c)

	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		slog.Error("Failed to load customers for navigation", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Navigation Error", "An error occurred while loading the customer navigation.", w, r)
	}

	utils.RenderSSE(w, r, utils.SSEOpts{
		Signals: buildCustomerPageSignals(c),
		Views: []templ.Component{
			views.Customer(c),
			views.HeaderIcon("customer"),
			views.CustomerNavigation(customers),
		},
	})
}

// DeleteCustomerSSE handles the deletion of a customer - if successful, it will render the dashboard and refresh the customer navigation
func (h *Handlers) DeleteCustomerSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("Invalid customer ID", "err", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid Customer ID", "The customer ID provided is not valid.", w, r)
		return
	}

	c, err := h.Queries.DeleteCustomer(r.Context(), parsedID)
	if err != nil {
		slog.Error("Error deleting customer", "err", err)
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		h.Notify(NotifyError, "Delete Failed", "An error occurred while trying to delete the customer.", w, r)
		return
	}

	// Soft delete all contacts for this customer
	err = h.Queries.DeleteContactsByCustomer(r.Context(), parsedID)
	if err != nil {
		slog.Error("Error soft deleting contacts for customer", "customer_id", parsedID, "err", err)
	}

	h.Notify(NotifySuccess, "Customer Deleted", fmt.Sprintf("Customer %s has been deleted and removed from active lists.", c.Name), w, r)
	al.LogCustomerDeleted(r.Context(), h.Queries, c)

	// render dashboard, refresh customer navigation
	h.DashboardSSE(w, r)

	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		slog.Error("Failed to load customers for navigation", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Navigation Error", "An error occurred while loading the customer navigation.", w, r)
	}
	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.CustomerNavigation(customers),
		},
	})
}

// GetCustomerOverviewSSE retrieves a customer overview by ID and renders it via SSE
func (h *Handlers) GetCustomerOverviewSSE(w http.ResponseWriter, r *http.Request) {
	c, ok := h.getCustomerByID(w, r, "customerID")
	if !ok {
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.CustomerOverview(c),
		},
	})
}

// getCustomerByID fetches a customer by ID from the URL param and handles errors consistently
func (h *Handlers) getCustomerByID(w http.ResponseWriter, r *http.Request, idParam string) (db.GetCustomerRow, bool) {
	id := chi.URLParam(r, idParam)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("Invalid customer ID", "err", err)
		h.Notify(NotifyError, "Invalid Customer ID", "The customer ID provided is not valid.", w, r)
		views.NotFound().Render(r.Context(), w)
		return db.GetCustomerRow{}, false
	}
	c, err := h.Queries.GetCustomer(r.Context(), parsedID)
	if err != nil {
		slog.Error("GetCustomer failed", "customer_id", parsedID, "err", err)
		h.Notify(NotifyError, "Customer Not Found", "No customer found for the provided ID.", w, r)
		views.NotFound().Render(r.Context(), w)
		return db.GetCustomerRow{}, false
	}
	return c, true
}

// DeleteCustomerLogoSSE handles the deletion of a customer logo, sets it to NULL in the DB, and returns an SSE response.
func (h *Handlers) DeleteCustomerLogoSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	customerID, err := uuid.Parse(id)
	if err != nil {
		slog.Error("Invalid customer ID", "err", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid Customer ID", "The customer ID provided is not valid.", w, r)
		return
	}

	// Update DB to set logo to NULL
	params := db.UpdateCustomerLogoParams{
		ID:   customerID,
		Logo: sql.NullString{String: "", Valid: false}, // Set logo to NULL
	}

	if err := h.Queries.UpdateCustomerLogo(r.Context(), params); err != nil {
		slog.Error("Error updating customer logo", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Update Failed", "An error occurred while updating the customer logo.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Logo Deleted", "Customer logo has been successfully deleted.", w, r)

	// Refresh the customer overview page
	updated, _ := h.Queries.GetCustomer(r.Context(), customerID)
	customers, _ := h.Queries.ListCustomers(r.Context())
	utils.RenderSSE(w, r, utils.SSEOpts{
		Views: []templ.Component{
			views.Customer(updated),
			views.CustomerNavigation(customers),
		},
	})

	// Delete the logo file from the filesystem
	// TODO: create dir if it doesn't exist
	matches, err := filepath.Glob(fmt.Sprintf("public/uploads/logos/%s*", customerID.String()))
	if err != nil {
		slog.Error("Error finding logo files", "err", err)
		return
	}
	for _, match := range matches {
		if err := os.Remove(match); err != nil {
			slog.Error("Error deleting logo file", "file", match, "err", err)
		}
	}
}

// UploadCustomerLogoSSE handles the upload of a customer logo, saves it, updates the DB, and returns an SSE response.
// TODO:
// 1. Figure out how to handle cases where the filename stays the same but the content changes (image doesn't change after refresh)
// 2. Not a big issue for logos, but public assets are, by nature, public. Need to consider the implications of this.
func (h *Handlers) UploadCustomerLogoSSE(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	customerID, err := uuid.Parse(id)

	if err != nil {
		slog.Error("Invalid customer ID", "err", err)
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid Customer ID", "The customer ID provided is not valid.", w, r)
		return
	}

	// Parse JSON body
	var payload struct {
		Logo      []string `json:"logo"`
		LogoMimes []string `json:"logoMimes"`
		LogoNames []string `json:"logoNames"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		slog.Error("Error decoding JSON", "err", err)
		h.Notify(NotifyError, "Upload Failed", "An error occurred while decoding the JSON payload.", w, r)
		return
	}
	if len(payload.Logo) == 0 || payload.Logo[0] == "" {
		http.Error(w, "No logo data provided", http.StatusBadRequest)
		h.Notify(NotifyError, "Upload Failed", "No logo data provided in the request.", w, r)
		return
	}

	// Determine file extension
	ext := utils.GetImageExtension(payload.LogoMimes, payload.LogoNames, ".png")

	// Decode base64
	data, err := utils.DecodeBase64Image(payload.Logo[0]) // handlers/utils
	if err != nil {
		slog.Error("Error decoding base64", "err", err)
		http.Error(w, "Invalid image data", http.StatusBadRequest)
		h.Notify(NotifyError, "Upload Failed", "An error occurred while decoding the image data.", w, r)
		return
	}

	// Save file
	logoPath := fmt.Sprintf("public/uploads/logos/%s%s", customerID.String(), ext)
	if err := os.WriteFile(logoPath, data, 0644); err != nil {
		slog.Error("Error saving file", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Upload Failed", "An error occurred while uploading the logo.", w, r)
		return
	}

	// Update DB with relative path
	urlPath := fmt.Sprintf("public/uploads/logos/%s%s", customerID.String(), ext)
	params := db.UpdateCustomerLogoParams{
		ID:   customerID,
		Logo: sql.NullString{String: urlPath, Valid: true},
	}

	if err := h.Queries.UpdateCustomerLogo(r.Context(), params); err != nil {
		slog.Error("Error updating customer logo", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Update Failed", "An error occurred while updating the customer logo.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Logo Uploaded", "Customer logo has been successfully uploaded.", w, r)

	// Refresh the customer overview page
	updated, _ := h.Queries.GetCustomer(r.Context(), customerID)
	customers, _ := h.Queries.ListCustomers(r.Context())
	utils.RenderSSE(w, r, utils.SSEOpts{
		Signals: []byte(`{"logo": ""}`),
		Views: []templ.Component{
			views.Customer(updated),
			views.CustomerNavigation(customers),
		},
	})
}

// buildCustomerPageSignals constructs the page signals for a customer overview
func buildCustomerPageSignals(c db.GetCustomerRow) []byte {
	pageSignals := utils.PageSignals{
		HeaderTitle: c.Name,
		HeaderDescription: fmt.Sprintf(
			"%d %s • %d %s • %d %s",
			c.ContactCount,
			utils.Pluralise(c.ContactCount, "contact", "contacts"),
			c.SubscriptionCount,
			utils.Pluralise(c.SubscriptionCount, "subscription", "subscriptions"),
			c.ProjectCount,
			utils.Pluralise(c.ProjectCount, "project", "projects"),
		),
		CurrentPage: c.ID.String(),
	}

	encodedSignals, _ := json.Marshal(pageSignals)
	return encodedSignals
}
