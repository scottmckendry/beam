package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

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
		log.Printf("No customerID provided in URL param")
		h.Notify(NotifyError, "Missing Customer ID", "No customer ID provided.", w, r)
		w.WriteHeader(http.StatusBadRequest)
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
		log.Printf("Error parsing form: %v", err)
		h.Notify(NotifyError, "Invalid form", "The submitted form is invalid.", w, r)
		w.WriteHeader(http.StatusBadRequest)
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
		log.Printf("Failed to get customer %s: %v", customerID, err)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the customer details. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Printf("Error adding contact: %v", err)
		h.Notify(NotifyError, "Failed to add contact", "An error occurred while adding the contact. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If this contact is primary, unset all others for this customer
	if isPrimary {
		err = h.Queries.UnsetOtherPrimaryContacts(r.Context(), db.UnsetOtherPrimaryContactsParams{
			CustomerID: cid,
			ID:         newContact.ID,
		})
		if err != nil {
			log.Printf("Failed to unset other primary contacts: %v", err)
		}
	}
	h.Notify(NotifySuccess, "Contact added", "The contact has been successfully added.", w, r)

	// Refresh the contact list for the customer
	contacts, err := h.Queries.ListContactsByCustomer(r.Context(), cid)
	if err != nil {
		log.Printf("Failed to list contacts for customer %s: %v", customerID, err)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the updated contacts. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	customer, err := h.Queries.GetCustomer(r.Context(), cid)
	if err != nil {
		log.Printf("Failed to get customer %s: %v", customerID, err)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the customer details. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Printf("Invalid contact ID: %v", err)
		h.Notify(NotifyError, "Invalid contact ID", "The provided contact ID is invalid.", w, r)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get the contact to find the customer ID (for refreshing list)
	contact, err := h.Queries.GetContact(r.Context(), cid)
	if err != nil {
		log.Printf("Contact not found for delete: %v", err)
		h.Notify(NotifyError, "Contact not found", "Could not find the contact to delete.", w, r)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Delete the contact (soft delete)
	_, err = h.Queries.DeleteContact(r.Context(), cid)
	if err != nil {
		log.Printf("Error deleting contact: %v", err)
		h.Notify(NotifyError, "Failed to delete contact", "An error occurred while deleting the contact. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Notify(NotifySuccess, "Contact deleted", "The contact has been successfully deleted.", w, r)
	// Refresh the contact list for the customer
	contacts, err := h.Queries.ListContactsByCustomer(r.Context(), contact.CustomerID)
	if err != nil {
		log.Printf("Failed to list contacts for customer %s: %v", contact.CustomerID.String(), err)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the updated contacts. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	customer, err := h.Queries.GetCustomer(r.Context(), contact.CustomerID)
	if err != nil {
		log.Printf("Failed to get customer %s: %v", contact.CustomerID.String(), err)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the customer details. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
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
		log.Printf("Invalid contact ID: %v", err)
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	contact, err := h.Queries.GetContact(r.Context(), cid)
	if err != nil {
		log.Printf("Failed to get contact: %v", err)
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
		log.Printf("Error parsing form: %v", err)
		h.Notify(NotifyError, "Invalid form", "The submitted form is invalid.", w, r)
		w.WriteHeader(http.StatusBadRequest)
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
		log.Printf("Invalid contact ID: %v", err)
		h.Notify(NotifyError, "Invalid contact ID", "The provided contact ID is invalid.", w, r)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// If this contact is being set as primary, unset all others for this customer
	if isPrimary {
		parsedCustID, err := uuid.Parse(customerID)
		if err == nil {
			err = h.Queries.UnsetOtherPrimaryContacts(r.Context(), db.UnsetOtherPrimaryContactsParams{
				CustomerID: parsedCustID,
				ID:         cid,
			})
			if err != nil {
				log.Printf("Failed to unset other primary contacts: %v", err)
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
		log.Printf("Failed to update contact: %v", err)
		h.Notify(NotifyError, "Failed to update contact", "An error occurred while updating the contact. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Notify(NotifySuccess, "Contact updated", "The contact has been successfully updated.", w, r)

	// Refresh the contact list for the customer
	parsedCustID, err := uuid.Parse(customerID)
	if err != nil {
		log.Printf("Invalid customer ID: %v", err)
		h.Notify(NotifyError, "Could not refresh contacts", "The provided customer ID is invalid.", w, r)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch the updated contacts for the customer
	contacts, err := h.Queries.ListContactsByCustomer(r.Context(), parsedCustID)
	if err != nil {
		log.Printf("Failed to list contacts for customer %s: %v", customerID, err)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the updated contacts. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	customer, err := h.Queries.GetCustomer(r.Context(), parsedCustID)
	if err != nil {
		log.Printf("Failed to get customer %s: %v", customerID, err)
		h.Notify(NotifyError, "Failed to refresh contacts", "An error occurred while fetching the customer details. Please try again.", w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.renderSSE(w, r, SSEOpts{
		Views: []templ.Component{
			views.CustomerContacts(customer, contacts),
		},
	})
}
