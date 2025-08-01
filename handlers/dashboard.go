package handlers

import (
	"encoding/json"
	"github.com/scottmckendry/beam/handlers/utils"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/scottmckendry/beam/ui/views"
)

// RegisterDashboardRoutes registers all dashboard-related routes on the given router.
func (h *Handlers) RegisterDashboardRoutes(r chi.Router) {
	r.Get("/sse/dashboard", h.DashboardSSE)
	r.Get("/sse/dashboard/stats", h.DashboardStatsSSE)
	r.Get("/sse/dashboard/activity", h.DashboardActivitySSE)
}

func (h *Handlers) DashboardStatsSSE(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Queries.GetDashboardStats(r.Context())
	if err != nil {
		slog.Error("Failed to load dashboard stats", "err", err)
		h.Notify(NotifyError, "Dashboard Error", "Failed to load dashboard stats.", w, r)
		http.Error(w, "Failed to load dashboard stats", http.StatusInternalServerError)
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{Views: []templ.Component{views.DashboardStats(stats)}})
}

func (h *Handlers) DashboardActivitySSE(w http.ResponseWriter, r *http.Request) {
	activities, err := h.Queries.GetRecentActivity(r.Context())
	if err != nil {
		slog.Error("Failed to load recent activity", "err", err)
		h.Notify(NotifyError, "Activity Error", "Failed to load recent activity.", w, r)
		http.Error(w, "Failed to load recent activity", http.StatusInternalServerError)
		return
	}

	utils.RenderSSE(w, r, utils.SSEOpts{Views: []templ.Component{views.DashboardActivity(activities)}})
}

func (h *Handlers) DashboardSSE(w http.ResponseWriter, r *http.Request) {
	pageSignals := utils.PageSignals{
		HeaderTitle:       "Dashboard",
		HeaderDescription: "Overview of your business metrics",
		CurrentPage:       "dashboard",
	}
	encodedSignals, _ := json.Marshal(pageSignals)

	utils.RenderSSE(w, r, utils.SSEOpts{
		Signals: encodedSignals,
		Views: []templ.Component{
			views.Dashboard(),
			views.HeaderIcon("dashboard"),
		},
	})
}
