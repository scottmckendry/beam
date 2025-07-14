package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/templ"

	"github.com/scottmckendry/beam/ui/views"
)

func (h *Handlers) DashboardStatsSSE(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Queries.GetDashboardStats(r.Context())
	if err != nil {
		log.Printf("Failed to load dashboard stats: %v", err)
		h.Notify(NotifyError, "Dashboard Error", "Failed to load dashboard stats.", w, r)
		http.Error(w, "Failed to load dashboard stats", http.StatusInternalServerError)
		return
	}

	h.renderSSE(w, r, SSEOpts{Views: []templ.Component{views.DashboardStats(stats)}})
}

func (h *Handlers) DashboardActivitySSE(w http.ResponseWriter, r *http.Request) {
	activities, err := h.Queries.GetRecentActivity(r.Context())
	if err != nil {
		log.Printf("Failed to load recent activity: %v", err)
		h.Notify(NotifyError, "Activity Error", "Failed to load recent activity.", w, r)
		http.Error(w, "Failed to load recent activity", http.StatusInternalServerError)
		return
	}

	h.renderSSE(w, r, SSEOpts{Views: []templ.Component{views.DashboardActivity(activities)}})
}

func (h *Handlers) DashboardSSE(w http.ResponseWriter, r *http.Request) {
	pageSignals := PageSignals{
		HeaderTitle:       "Dashboard",
		HeaderDescription: "Overview of your business metrics",
		CurrentPage:       "dashboard",
	}
	encodedSignals, _ := json.Marshal(pageSignals)

	h.renderSSE(w, r, SSEOpts{
		Signals: encodedSignals,
		Views: []templ.Component{
			views.Dashboard(),
			views.HeaderIcon("dashboard"),
		},
	})
}
