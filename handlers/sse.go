package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/scottmckendry/beam/ui/views"
)

type PageSignals struct {
	HeaderTitle       string `json:"headerTitle"`
	HeaderDescription string `json:"headerDescription"`
	CurrentPage       string `json:"currentPage,omitempty"`
}

func serveSSEFragment(w http.ResponseWriter, r *http.Request, fragment string) {
	sse := datastar.NewSSE(w, r)
	sse.MergeFragments(
		fragment,
		datastar.WithUseViewTransitions(true),
	)
}

func pluralise(count int64, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// HandleSSEDashboardStats streams the dashboard stats cards via SSE for Datastar
func (h *Handlers) HandleSSEDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Queries.GetDashboardStats(r.Context())
	if err != nil {
		log.Printf("Failed to load dashboard stats: %v", err)
		http.Error(w, "Failed to load dashboard stats", http.StatusInternalServerError)
		return
	}

	buf := &bytes.Buffer{}
	views.DashboardStats(stats).Render(r.Context(), buf)
	serveSSEFragment(w, r, buf.String())
}

// HandleSSEDashboardActivity streams the dashboard recent activity via SSE for Datastar
func (h *Handlers) HandleSSEDashboardActivity(w http.ResponseWriter, r *http.Request) {
	activities, err := h.Queries.GetRecentActivity(r.Context())
	if err != nil {
		log.Printf("Failed to load recent activity: %v", err)
		http.Error(w, "Failed to load recent activity", http.StatusInternalServerError)
		return
	}

	buf := &bytes.Buffer{}
	views.DashboardActivity(activities).Render(r.Context(), buf)
	serveSSEFragment(w, r, buf.String())
}

// HandleSSECustomerNav streams the rendered CustomerNavigation component via SSE for Datastar
func (h *Handlers) HandleSSECustomerNav(w http.ResponseWriter, r *http.Request) {
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		log.Printf("Failed to load customers: %v", err)
		http.Error(w, "Failed to load customers", http.StatusInternalServerError)
		return
	}
	currentPage := r.URL.Query().Get("page")
	buf := &bytes.Buffer{}

	views.CustomerNavigation(customers, currentPage).Render(r.Context(), buf)
	serveSSEFragment(w, r, buf.String())
}

// HandleSSEGetAddCustomer streams the rendered AddCustomer component via SSE for Datastar, including header signals.
func (h *Handlers) HandleSSEGetAddCustomer(w http.ResponseWriter, r *http.Request) {
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
	sse.MergeSignals(encodedSignals)
	sse.MergeFragments(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}

func (h *Handlers) HandleSSEGetInvoices(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	views.Invoices().Render(r.Context(), buf)
	views.HeaderIcon("invoices").Render(r.Context(), buf)

	pageSignals := PageSignals{
		HeaderTitle:       "Invoices",
		HeaderDescription: "Manage invoices for all customers",
		CurrentPage:       "invoices",
	}

	encodedSignals, _ := json.Marshal(pageSignals)

	sse := datastar.NewSSE(w, r)
	sse.MergeSignals(encodedSignals)
	sse.MergeFragments(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}

func (h *Handlers) HandleSSEGetDashboard(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	views.Dashboard().Render(r.Context(), buf)
	views.HeaderIcon("dashboard").Render(r.Context(), buf)

	pageSignals := PageSignals{
		HeaderTitle:       "Dashboard",
		HeaderDescription: "Overview of your business metrics",
		CurrentPage:       "dashboard",
	}

	encodedSignals, _ := json.Marshal(pageSignals)

	sse := datastar.NewSSE(w, r)
	sse.MergeSignals(encodedSignals)
	sse.MergeFragments(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}

func (h *Handlers) HandleSSEGetCustomer(w http.ResponseWriter, r *http.Request) {
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
			pluralise(c.ContactCount, "contact", "contacts"),
			c.SubscriptionCount,
			pluralise(c.SubscriptionCount, "subscription", "subscriptions"),
			c.ProjectCount,
			pluralise(c.ProjectCount, "project", "projects"),
		),
		CurrentPage: c.ID.String(),
	}

	encodedSignals, _ := json.Marshal(pageSignals)

	sse := datastar.NewSSE(w, r)
	sse.MergeSignals(encodedSignals)
	sse.MergeFragments(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}

// HandleSSEAddCustomer adds a new customer via SSE for Datastar.
func (h *Handlers) HandleSSEAddCustomer(w http.ResponseWriter, r *http.Request) {
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
