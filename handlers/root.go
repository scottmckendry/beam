package handlers

import (
	"log"
	"net/http"

	"github.com/scottmckendry/beam/ui/views"
)

// HandleRoot serves the main application page.
func (h *Handlers) HandleRoot(w http.ResponseWriter, r *http.Request) {
	customers, err := h.Queries.ListCustomers(r.Context())
	if err != nil {
		log.Printf("Failed to load customers: %v", err)
		http.Error(w, "Failed to load customers", http.StatusInternalServerError)
		return
	}
	// Admin middleware handles the admin check, so we can assume the user is authenticated here and has admin privileges.
	views.Root(true, customers).Render(r.Context(), w)
}
