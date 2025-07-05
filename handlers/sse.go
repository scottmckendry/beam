package handlers

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"time"

	datastar "github.com/starfederation/datastar/sdk/go"

	"github.com/scottmckendry/beam/ui/views"
)

const simulateSlowEvents = false

// getRandomDelay returns a random delay between 500ms and 1500ms
func getRandomDelay() time.Duration {
	return time.Duration(rand.Float64()*1000+500) * time.Millisecond
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

	sse := datastar.NewSSE(w, r)
	buf := &bytes.Buffer{}
	views.DashboardStats(stats).Render(r.Context(), buf)
	sse.MergeFragments(
		`<div id="dashboard-stats-section">`+buf.String()+`</div>`,
		datastar.WithUseViewTransitions(true),
	)
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

	sse := datastar.NewSSE(w, r)
	buf := &bytes.Buffer{}
	views.DashboardActivity(activities).Render(r.Context(), buf)
	sse.MergeFragments(
		`<div id="dashboard-activity-section">`+buf.String()+`</div>`,
		datastar.WithUseViewTransitions(true),
	)
}

// HandleSSECustomerNav streams the rendered CustomerNavigation component via SSE for Datastar
func (h *Handlers) HandleSSECustomerNav(w http.ResponseWriter, r *http.Request) {
	if simulateSlowEvents {
		time.Sleep(getRandomDelay())
	}
	sse := datastar.NewSSE(w, r)
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		http.Error(w, "Failed to load customers", http.StatusInternalServerError)
		return
	}
	currentPage := r.URL.Query().Get("page")
	buf := &bytes.Buffer{}
	views.CustomerNavigation(customers, currentPage).Render(r.Context(), buf)
	sse.MergeFragments(
		`<div id="customer-nav-section">`+buf.String()+`</div>`,
		datastar.WithUseViewTransitions(true),
	)
}
