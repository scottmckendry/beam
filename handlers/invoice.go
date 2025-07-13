package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/starfederation/datastar/sdk/go/datastar"

	"github.com/scottmckendry/beam/ui/views"
)

func (h *Handlers) InvoicesSSE(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}
	if err := views.Invoices().Render(r.Context(), buf); err != nil {
		h.Notify(NotifyError, "Invoice Error", "Failed to load invoices.", w, r)
		http.Error(w, "Failed to load invoices", http.StatusInternalServerError)
		return
	}
	views.HeaderIcon("invoices").Render(r.Context(), buf)
	pageSignals := PageSignals{
		HeaderTitle:       "Invoices",
		HeaderDescription: "Manage invoices for all customers",
		CurrentPage:       "invoices",
	}
	encodedSignals, _ := json.Marshal(pageSignals)
	sse := datastar.NewSSE(w, r)
	sse.PatchSignals(encodedSignals)
	sse.PatchElements(
		buf.String(),
		datastar.WithUseViewTransitions(true),
	)
}
