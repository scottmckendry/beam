package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/starfederation/datastar/sdk/go/datastar"

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
	buf := &bytes.Buffer{}
	views.DashboardStats(stats).Render(r.Context(), buf)
	ServeSSEElement(w, r, buf.String())
}

func (h *Handlers) DashboardActivitySSE(w http.ResponseWriter, r *http.Request) {
	activities, err := h.Queries.GetRecentActivity(r.Context())
	if err != nil {
		log.Printf("Failed to load recent activity: %v", err)
		h.Notify(NotifyError, "Activity Error", "Failed to load recent activity.", w, r)
		http.Error(w, "Failed to load recent activity", http.StatusInternalServerError)
		return
	}
	buf := &bytes.Buffer{}
	views.DashboardActivity(activities).Render(r.Context(), buf)
	ServeSSEElement(w, r, buf.String())
}

func (h *Handlers) DashboardSSE(w http.ResponseWriter, r *http.Request) {
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
	sse.PatchSignals(encodedSignals)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}
