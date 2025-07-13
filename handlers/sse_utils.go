package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

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

// MapFormToStruct maps form values to a struct using field names as keys
// It automatically handles sql.NullString conversions
func MapFormToStruct(r *http.Request, dest any) error {
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
