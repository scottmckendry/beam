package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"

	"github.com/scottmckendry/beam/ui/views"
)

func (h *Handlers) InvoicesSSE(w http.ResponseWriter, r *http.Request) {
	pageSignals := PageSignals{
		HeaderTitle:       "Invoices",
		HeaderDescription: "Manage invoices for all customers",
		CurrentPage:       "invoices",
	}
	encodedSignals, _ := json.Marshal(pageSignals)

	h.renderSSE(w, r, SSEOpts{
		Signals: encodedSignals,
		Views: []templ.Component{
			views.Invoices(),
			views.HeaderIcon("invoices"),
		},
	})
}
