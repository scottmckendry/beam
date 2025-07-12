package handlers

import (
	"net/http"

	"github.com/starfederation/datastar/sdk/go/datastar"
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
