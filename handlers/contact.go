package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	al "github.com/scottmckendry/beam/activitylog"
	"github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/ui/views"
)

// RegisterContactRoutes registers all contact-related routes on the given router.
func (h *Handlers) RegisterContactRoutes(r chi.Router) {
	r.Get("/sse/customer/{customerID}/add-contact", h.AddContactFormSSE)
	r.Get("/sse/customer/{customerID}/add-contact-submit", h.AddContactSubmitSSE)
	r.Get("/sse/customer/{customerID}/edit-contact/{contactID}", h.EditContactFormSSE)
	r.Get("/sse/customer/{customerID}/edit-contact-submit/{contactID}", h.EditContactSubmitSSE)
	r.Get("/sse/customer/{customerID}/delete-contact/{contactID}", h.DeleteContactSSE)
}

// AddContactFormSSE renders the form to add a new contact for a customer via SSE.
func (h *Handlers) AddContactFormSSE(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerID")
	if customerID == "" {
		slog.Error("No customerID provided in URL param")
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Missing Customer ID", "No customer ID provided.", w, r)
		return
	}

	h.renderSSE(w, r, SSEOpts{
		Views: []templ.Component{
			views.AddContact(customerID),
		},
	})
}

// AddContactSubmitSSE handles the submission of the add contact form, creates the contact, and refreshes the contact list via SSE.
func (h *Handlers) AddContactSubmitSSE(w http.ResponseWriter, r *http.Request) {
	customerID := chi.URLParam(r, "customerID")
	if err := r.ParseForm(); err != nil {
		slog.Error("Error parsing form", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid form", "The submitted form is invalid.", w, r)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	role := r.FormValue("role")
	isPrimary := r.FormValue("is_primary") == "on"
	notes := r.FormValue("notes")

	cid, err := uuid.Parse(customerID)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}

	newContact, err := h.Queries.CreateContact(r.Context(), db.CreateContactParams{
		CustomerID: cid,
		Name:       name,
		Role:       sql.NullString{String: role, Valid: role != ""},
		Email:      sql.NullString{String: email, Valid: email != ""},
		Phone:      sql.NullString{String: phone, Valid: phone != ""},
		IsPrimary:  sql.NullBool{Bool: isPrimary, Valid: true},
		Notes:      sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		slog.Error("Error adding contact", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to add contact", "An error occurred while adding the contact. Please try again.", w, r)
		return
	}

	// If this contact is primary, unset all others for this customer
	if isPrimary {
		err = h.Queries.UnsetOtherPrimaryContacts(r.Context(), db.UnsetOtherPrimaryContactsParams{
			CustomerID: cid,
			ID:         newContact.ID,
		})
		if err != nil {
			slog.Error("Failed to unset other primary contacts", "err", err)
		}
	}

	h.Notify(NotifySuccess, "Contact added", "The contact has been successfully added.", w, r)
	al.LogContactAdded(r.Context(), h.Queries, newContact.CustomerID, newContact.Name)

	// Refresh the contact list for the customer
	contacts, err := h.Queries.ListContactsByCustomer(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to list contacts for customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the updated contacts. Please try again.", w, r)
		return
	}
	customer, err := h.Queries.GetCustomer(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}

	h.renderSSE(w, r, SSEOpts{
		Views: []templ.Component{
			views.CustomerContacts(customer, contacts),
		},
	})
}

// DeleteContactSSE handles deleting a contact and refreshing the contact list via SSE.
func (h *Handlers) DeleteContactSSE(w http.ResponseWriter, r *http.Request) {
	contactID := chi.URLParam(r, "contactID")
	cid, err := uuid.Parse(contactID)
	if err != nil {
		slog.Error("Invalid contact ID", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid contact ID", "The provided contact ID is invalid.", w, r)
		return
	}
	// Get the contact to find the customer ID (for refreshing list)
	contact, err := h.Queries.GetContact(r.Context(), cid)
	if err != nil {
		slog.Error("Contact not found for delete", "err", err)
		h.Notify(NotifyError, "Contact not found", "Could not find the contact to delete.", w, r)
		return
	}
	// Delete the contact (soft delete)
	_, err = h.Queries.DeleteContact(r.Context(), cid)
	if err != nil {
		slog.Error("Error deleting contact", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to delete contact", "An error occurred while deleting the contact. Please try again.", w, r)
		return
	}
	h.Notify(NotifySuccess, "Contact deleted", "The contact has been successfully deleted.", w, r)
	al.LogContactDeleted(r.Context(), h.Queries, contact.CustomerID, contact.Name)

	// Refresh the contact list for the customer
	contacts, err := h.Queries.ListContactsByCustomer(r.Context(), contact.CustomerID)
	if err != nil {
		slog.Error("Failed to list contacts for customer", "customerID", contact.CustomerID.String(), "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the updated contacts. Please try again.", w, r)
		return
	}
	customer, err := h.Queries.GetCustomer(r.Context(), contact.CustomerID)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", contact.CustomerID.String(), "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}

	h.renderSSE(w, r, SSEOpts{
		Views: []templ.Component{
			views.CustomerContacts(customer, contacts),
		},
	})
}

// EditContactFormSSE renders the form to edit an existing contact via SSE.
func (h *Handlers) EditContactFormSSE(w http.ResponseWriter, r *http.Request) {
	contactID := chi.URLParam(r, "contactID")
	cid, err := uuid.Parse(contactID)
	if err != nil {
		slog.Error("Invalid contact ID", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid contact ID", "The provided contact ID is invalid.", w, r)
		return
	}
	contact, err := h.Queries.GetContact(r.Context(), cid)
	if err != nil {
		slog.Error("Failed to get contact", "err", err)
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}
	views.EditContact(contact).Render(r.Context(), w)
}

// EditContactSubmitSSE handles the submission of the edit contact form, updates the contact, and refreshes the contact list via SSE.
func (h *Handlers) EditContactSubmitSSE(w http.ResponseWriter, r *http.Request) {
	contactID := chi.URLParam(r, "contactID")
	customerID := chi.URLParam(r, "customerID")
	if err := r.ParseForm(); err != nil {
		slog.Error("Error parsing form", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid form", "The submitted form is invalid.", w, r)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	role := r.FormValue("role")
	isPrimary := r.FormValue("is_primary") == "on"
	notes := r.FormValue("notes")

	cid, err := uuid.Parse(contactID)
	if err != nil {
		slog.Error("Invalid contact ID", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		h.Notify(NotifyError, "Invalid contact ID", "The provided contact ID is invalid.", w, r)
		return
	}
	parsedCustID, err := uuid.Parse(customerID)

	// If this contact is being set as primary, unset all others for this customer
	if isPrimary {
		if err == nil {
			err = h.Queries.UnsetOtherPrimaryContacts(r.Context(), db.UnsetOtherPrimaryContactsParams{
				CustomerID: parsedCustID,
				ID:         cid,
			})
			if err != nil {
				slog.Error("Failed to unset other primary contacts", "err", err)
			}
		}
	}
	err = h.Queries.UpdateContact(r.Context(), db.UpdateContactParams{
		ID:        cid,
		Name:      name,
		Role:      sql.NullString{String: role, Valid: role != ""},
		Email:     sql.NullString{String: email, Valid: email != ""},
		Phone:     sql.NullString{String: phone, Valid: phone != ""},
		IsPrimary: sql.NullBool{Bool: isPrimary, Valid: true},
		Notes:     sql.NullString{String: notes, Valid: notes != ""},
	})
	if err != nil {
		slog.Error("Failed to update contact", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to update contact", "An error occurred while updating the contact. Please try again.", w, r)
		return
	}

	h.Notify(NotifySuccess, "Contact updated", "The contact has been successfully updated.", w, r)
	al.LogContactUpdated(r.Context(), h.Queries, parsedCustID, name)

	// Fetch the updated contacts for the customer
	contacts, err := h.Queries.ListContactsByCustomer(r.Context(), parsedCustID)
	if err != nil {
		slog.Error("Failed to list contacts for customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the updated contacts. Please try again.", w, r)
		return
	}

	customer, err := h.Queries.GetCustomer(r.Context(), parsedCustID)
	if err != nil {
		slog.Error("Failed to get customer", "customerID", customerID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the customer details. Please try again.", w, r)
		return
	}

	h.renderSSE(w, r, SSEOpts{
		Views: []templ.Component{
			views.CustomerContacts(customer, contacts),
		},
	})
}
