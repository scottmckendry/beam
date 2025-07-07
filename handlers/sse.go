package handlers

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/scottmckendry/beam/ui/views"
)

const simulateSlowEvents = false

// getRandomDelay returns a random delay between 100ms and 500ms
// good for testing lazy-loading and view transitions - wehre view transitions are not supported (*cough, cough* Firefox...), this will make things look a bit nicer and less "flickery"
func getRandomDelay() time.Duration {
	return time.Duration(rand.Float64()*400+100) * time.Millisecond
}

func serveSSEFragment(w http.ResponseWriter, r *http.Request, fragment string) {
	sse := datastar.NewSSE(w, r)
	sse.MergeFragments(
		fragment,
		datastar.WithUseViewTransitions(true),
	)
}

// HandleSSEDashboardStats streams the dashboard stats cards via SSE for Datastar
func (h *Handlers) HandleSSEDashboardStats(w http.ResponseWriter, r *http.Request) {
	if simulateSlowEvents {
		time.Sleep(getRandomDelay())
	}

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
	if simulateSlowEvents {
		time.Sleep(getRandomDelay())
	}

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
	if simulateSlowEvents {
		time.Sleep(getRandomDelay())
	}

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

// HandleSSECustomerOverview streams the rendered CustomerOverview component via SSE for Datastar
func (h *Handlers) HandleSSECustomerOverview(w http.ResponseWriter, r *http.Request) {
	if simulateSlowEvents {
		time.Sleep(getRandomDelay())
	}
	customerID := chi.URLParam(r, "id")
	if customerID == "" {
		log.Println("Id is required")
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	cid, err := uuid.Parse(customerID)
	if err != nil {
		log.Printf("Invalid customer_id: %v", err)
		http.Error(w, "Invalid customer_id", http.StatusBadRequest)
		return
	}

	customer, err := h.Queries.GetCustomer(r.Context(), cid)
	if err != nil {
		log.Printf("Failed to load customer: %v", err)
		http.Error(w, "Failed to load customer", http.StatusInternalServerError)
		return
	}

	buf := &bytes.Buffer{}
	views.CustomerOverview(customer).Render(r.Context(), buf)
	serveSSEFragment(w, r, buf.String())
}

// HandleSSEGetAddCustomer streams the rendered AddCustomer component via SSE for Datastar, including header signals.
// This is an on-click element swap, so 'lazy-loading' is not required.
func (h *Handlers) HandleSSEGetAddCustomer(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	views.AddCustomer().Render(r.Context(), buf)

	headerSignal := []byte(
		`{"headerTitle":"Add Customer","headerDescription":"Woohoo! Let's add a new customer 🚀"}`,
	)

	sse := datastar.NewSSE(w, r)
	sse.MergeSignals(headerSignal)
	sse.MergeFragments(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}
