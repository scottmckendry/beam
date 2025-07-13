package handlers

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/starfederation/datastar/sdk/go/datastar"

	"github.com/scottmckendry/beam/db/sqlc"
)

type PageSignals struct {
	HeaderTitle       string `json:"_headerTitle"`
	HeaderDescription string `json:"_headerDescription"`
	CurrentPage       string `json:"_currentPage,omitempty"`
}

func ServeSSEElement(w http.ResponseWriter, r *http.Request, elements string) {
	sse := datastar.NewSSE(w, r)
	sse.PatchElements(
		elements,
		datastar.WithUseViewTransitions(true),
	)
}

func Pluralise(count int64, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

// logActivity inserts a new activity log entry for a customer
// TODO: make activity types type-safe with constants or an enum
// TODO: pass a struct for activity details instead of multiple strings
func (h *Handlers) logActivity(
	r *http.Request,
	customerID uuid.UUID,
	activityType,
	action, description string,
) {
	activity := db.LogActivityParams{
		CustomerID:   customerID,
		ActivityType: activityType,
		Action:       action,
		Description:  description,
	}
	_, err := h.Queries.LogActivity(r.Context(), activity)
	if err != nil {
		log.Printf("Failed to log customer activity: %v", err)
	}
}
