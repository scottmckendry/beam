package handlers

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/starfederation/datastar/sdk/go/datastar"

	"github.com/scottmckendry/beam/db/sqlc"
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

// renderSSE renders a collection of templ.Components to a Server-Sent Events (SSE) response.
// Optionally accepts signals to patch into the SSE stream.
func (h *Handlers) renderSSE(w http.ResponseWriter, r *http.Request, opts SSEOpts) error {
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

	// TODO: remove the replace mode - depending on the outcome of #999
	sse.PatchElements(buf.String(), datastar.WithUseViewTransitions(true), datastar.WithModeReplace())
	return nil
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

// mapFormToStruct maps form values to a struct using field names as keys
// handles uuid.UUID and sql.NullString types specifically
func mapFormToStruct(r *http.Request, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Get form field name from struct field name (converting to lower case)
		formKey := strings.ToLower(field.Name)
		formValue := r.FormValue(formKey)

		// Handle uuid.UUID fields
		if field.Type == reflect.TypeOf(uuid.UUID{}) {
			if formValue == "" {
				fieldValue.Set(reflect.Zero(field.Type))
				continue
			}
			parsedUUID, err := uuid.Parse(formValue)
			if err != nil {
				return fmt.Errorf("invalid UUID for field %s: %w", field.Name, err)
			}
			fieldValue.Set(reflect.ValueOf(parsedUUID))
			continue
		}

		// Handle sql.NullString fields
		if field.Type == reflect.TypeOf(sql.NullString{}) {
			fieldValue.Set(reflect.ValueOf(sql.NullString{
				String: formValue,
				Valid:  formValue != "",
			}))
		} else {
			// Handle regular string fields
			fieldValue.SetString(formValue)
		}
	}

	return nil
}

// pluralize returns the plural form of a word based on the count.
func pluralise(count int64, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}
