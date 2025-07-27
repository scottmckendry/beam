package handlers

import (
	"encoding/json"
	"github.com/scottmckendry/beam/handlers/utils"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/scottmckendry/beam/ui/views"
)

// RegisterInvoiceRoutes registers all invoice-related routes on the given router.
func (h *Handlers) RegisterInvoiceRoutes(r chi.Router) {
	r.Get("/sse/invoice", h.InvoicesSSE)
}

func (h *Handlers) InvoicesSSE(w http.ResponseWriter, r *http.Request) {
	pageSignals := utils.PageSignals{
		HeaderTitle:       "Invoices",
		HeaderDescription: "Manage invoices for all customers",
		CurrentPage:       "invoices",
	}
	encodedSignals, _ := json.Marshal(pageSignals)

	utils.RenderSSE(w, r, utils.SSEOpts{
		Signals: encodedSignals,
		Views: []templ.Component{
			views.Invoices(),
			views.HeaderIcon("invoices"),
		},
	})
}
