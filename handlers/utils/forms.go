package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MapFormToStruct maps form values to a struct using field tags or field names.
func MapFormToStruct(r *http.Request, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}
	v = v.Elem()
	t := v.Type()
	var errs []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		formKey := field.Tag.Get("form")
		if formKey == "" {
			formKey = strings.ToLower(field.Name)
		}
		formValue := r.FormValue(formKey)
		if err := setFieldValue(field, fieldValue, formValue); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", field.Name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("form mapping errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func setFieldValue(field reflect.StructField, fieldValue reflect.Value, formValue string) error {
	// Handle basic types
	switch field.Type.Kind() {
	case reflect.String:
		fieldValue.SetString(formValue)
		return nil
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return setIntField(fieldValue, formValue)
	case reflect.Float64, reflect.Float32:
		return setFloatField(fieldValue, formValue)
	case reflect.Bool:
		fieldValue.SetBool(formValue == "on" || formValue == "true" || formValue == "1")
		return nil
	}

	// Handle special types
	switch field.Type {
	case reflect.TypeOf(uuid.UUID{}):
		return setUUIDField(fieldValue, formValue)
	case reflect.TypeOf(sql.NullString{}):
		fieldValue.Set(reflect.ValueOf(sql.NullString{String: formValue, Valid: formValue != ""}))
	case reflect.TypeOf(sql.NullBool{}):
		fieldValue.Set(reflect.ValueOf(sql.NullBool{
			Bool:  formValue == "on" || formValue == "true" || formValue == "1",
			Valid: true,
		}))
	case reflect.TypeOf(sql.NullInt64{}):
		return setNullInt64Field(fieldValue, formValue)
	case reflect.TypeOf(sql.NullFloat64{}):
		return setNullFloat64Field(fieldValue, formValue)
	case reflect.TypeOf(time.Time{}):
		return setTimeField(fieldValue, formValue)
	}
	return nil
}

func setIntField(fieldValue reflect.Value, formValue string) error {
	if formValue == "" {
		fieldValue.SetInt(0)
		return nil
	}
	i, err := strconv.ParseInt(formValue, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer value: %v", err)
	}
	fieldValue.SetInt(i)
	return nil
}

func setFloatField(fieldValue reflect.Value, formValue string) error {
	if formValue == "" {
		fieldValue.SetFloat(0)
		return nil
	}
	f, err := strconv.ParseFloat(formValue, 64)
	if err != nil {
		return fmt.Errorf("invalid float value: %v", err)
	}
	fieldValue.SetFloat(f)
	return nil
}

func setUUIDField(fieldValue reflect.Value, formValue string) error {
	if formValue == "" {
		fieldValue.Set(reflect.Zero(fieldValue.Type()))
		return nil
	}
	parsedUUID, err := uuid.Parse(formValue)
	if err != nil {
		return fmt.Errorf("invalid UUID: %v", err)
	}
	fieldValue.Set(reflect.ValueOf(parsedUUID))
	return nil
}

func setTimeField(fieldValue reflect.Value, formValue string) error {
	if formValue == "" {
		fieldValue.Set(reflect.Zero(fieldValue.Type()))
		return nil
	}
	t, err := time.Parse("2006-01-02", formValue)
	if err != nil {
		return fmt.Errorf("invalid date format: %v", err)
	}
	fieldValue.Set(reflect.ValueOf(t))
	return nil
}

func setNullInt64Field(fieldValue reflect.Value, formValue string) error {
	if formValue == "" {
		fieldValue.Set(reflect.ValueOf(sql.NullInt64{Valid: false}))
		return nil
	}
	i, err := strconv.ParseInt(formValue, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer value: %v", err)
	}
	fieldValue.Set(reflect.ValueOf(sql.NullInt64{Int64: i, Valid: true}))
	return nil
}

func setNullFloat64Field(fieldValue reflect.Value, formValue string) error {
	if formValue == "" {
		fieldValue.Set(reflect.ValueOf(sql.NullFloat64{Valid: false}))
		return nil
	}
	f, err := strconv.ParseFloat(formValue, 64)
	if err != nil {
		return fmt.Errorf("invalid float value: %v", err)
	}
	fieldValue.Set(reflect.ValueOf(sql.NullFloat64{Float64: f, Valid: true}))
	return nil
}
