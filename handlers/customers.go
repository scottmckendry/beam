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

	"github.com/scottmckendry/beam/ui/views"
)

func (h *Handlers) CustomerNavSSE(w http.ResponseWriter, r *http.Request) {
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		log.Printf("Failed to load customers: %v", err)
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
	c, err := h.Queries.GetCustomer(r.Context(), parsedID)
	if err != nil {
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
	type DatastarPayload struct {
		Customer struct {
			Name   string `json:"name"`
			Email  string `json:"email"`
			Status string `json:"status"`
			Notes  string `json:"notes"`
		} `json:"customer"`
	}
	var payload DatastarPayload
	datastarParam := r.URL.Query().Get("datastar")
	if err := json.Unmarshal([]byte(datastarParam), &payload); err != nil {
		log.Printf("Error parsing customer data: %v", err)
		http.Error(w, "Invalid customer data", http.StatusBadRequest)
		return
	}
	customer := payload.Customer
	log.Printf("Extracted customer data: %+v", customer)
	// Now you can use customer.Name, customer.Email, etc.
}
