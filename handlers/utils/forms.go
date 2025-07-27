package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"database/sql"
	"github.com/google/uuid"
)

// MapFormToStruct maps form values to a struct using field names as keys
// handles uuid.UUID and sql.NullString types specifically
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
