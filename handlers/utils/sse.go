package utils

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/starfederation/datastar-go/datastar"
)

type PageSignals struct {
	HeaderTitle       string `json:"_headerTitle"`
	HeaderDescription string `json:"_headerDescription"`
	CurrentPage       string `json:"_currentPage,omitempty"`
}

type SSEOpts struct {
	Signals []byte
	Views   []templ.Component
}

// RenderSSE renders a collection of templ.Components to a Server-Sent Events (SSE) response.
// Optionally accepts signals to patch into the SSE stream.
func RenderSSE(w http.ResponseWriter, r *http.Request, opts SSEOpts) error {
	buf := &bytes.Buffer{}
	for _, view := range opts.Views {
		if err := view.Render(r.Context(), buf); err != nil {
			return fmt.Errorf("failed to render view: %w", err)
		}
	}

	sse := datastar.NewSSE(w, r)
	if opts.Signals != nil {
		sse.PatchSignals(opts.Signals)
	}

	sse.PatchElements(buf.String(), datastar.WithUseViewTransitions(true))
	return nil
}
