package utils

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"strings"
	"time"
)

func TestDecodeBase64Image(t *testing.T) {
	plain := base64.StdEncoding.EncodeToString([]byte("hello"))
	withPrefix := "data:image/png;base64," + plain
	got, err := DecodeBase64Image(plain)
	if err != nil || string(got) != "hello" {
		t.Errorf("plain: got %q, err %v", got, err)
	}
	got, err = DecodeBase64Image(withPrefix)
	if err != nil || string(got) != "hello" {
		t.Errorf("withPrefix: got %q, err %v", got, err)
	}
}

func TestGetImageExtension(t *testing.T) {
	if ext := GetImageExtension([]string{"image/jpeg"}, nil, ".png"); ext != ".jpg" {
		t.Errorf("got %q, want .jpg", ext)
	}
	if ext := GetImageExtension(nil, []string{"foo.gif"}, ".png"); ext != ".gif" {
		t.Errorf("got %q, want .gif", ext)
	}
}

func TestPluralise(t *testing.T) {
	if got := Pluralise(1, "cat", "cats"); got != "cat" {
		t.Errorf("got %q, want cat", got)
	}
	if got := Pluralise(2, "cat", "cats"); got != "cats" {
		t.Errorf("got %q, want cats", got)
	}
}

type testFormStruct struct {
	Name  string
	Email sql.NullString
	ID    uuid.UUID
}

type allTypesFormStruct struct {
	Str      string
	Int      int
	Int64    int64
	Float32  float32
	Float64  float64
	Bool     bool
	UUID     uuid.UUID
	NullStr  sql.NullString
	NullBool sql.NullBool
	NullInt  sql.NullInt64
	NullFlt  sql.NullFloat64
	Date     time.Time
}

func TestMapFormToStruct_EdgeCases(t *testing.T) {
	form := url.Values{}
	r, _ := http.NewRequest("POST", "/", nil)
	r.Form = form
	// Nil pointer
	var nilPtr *testFormStruct = nil
	err := MapFormToStruct(r, nilPtr)
	if err == nil || !strings.Contains(err.Error(), "non-nil pointer") {
		t.Errorf("nil pointer: got err %v, want non-nil pointer error", err)
	}
	// Non-pointer
	var notPtr testFormStruct
	err = MapFormToStruct(r, notPtr)
	if err == nil || !strings.Contains(err.Error(), "non-nil pointer") {
		t.Errorf("non-pointer: got err %v, want non-nil pointer error", err)
	}
}

func TestMapFormToStruct_MixedTypes(t *testing.T) {
	type mixedFormStruct struct {
		Str     string
		Int     int
		Bool    bool
		UUID    uuid.UUID
		NullStr sql.NullString
		Date    time.Time
	}
	id := uuid.New()
	form := url.Values{}
	form.Set("str", "mixed")
	form.Set("int", "7")
	form.Set("bool", "1")
	form.Set("uuid", id.String())
	form.Set("nullstr", "bar")
	form.Set("date", "2025-07-30")
	r, _ := http.NewRequest("POST", "/", nil)
	r.Form = form
	var dest mixedFormStruct
	if err := MapFormToStruct(r, &dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.Str != "mixed" || dest.Int != 7 || !dest.Bool || dest.UUID != id || !dest.NullStr.Valid || dest.NullStr.String != "bar" || dest.Date.Format("2006-01-02") != "2025-07-30" {
		t.Errorf("got %+v, want correct values", dest)
	}
}

func TestMapFormToStruct_FieldTags(t *testing.T) {
	type tagFormStruct struct {
		FirstName string `form:"fname"`
		LastName  string // no tag, should use field name
	}
	form := url.Values{}
	form.Set("fname", "John")
	form.Set("lastname", "Doe")
	r, _ := http.NewRequest("POST", "/", nil)
	r.Form = form
	var dest tagFormStruct
	if err := MapFormToStruct(r, &dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.FirstName != "John" || dest.LastName != "Doe" {
		t.Errorf("got %+v, want FirstName=John, LastName=Doe", dest)
	}
}

func TestMapFormToStruct_EmptyValues(t *testing.T) {
	type emptyFormStruct struct {
		Str      string
		Int      int
		Float    float64
		Bool     bool
		UUID     uuid.UUID
		NullStr  sql.NullString
		NullBool sql.NullBool
		NullInt  sql.NullInt64
		NullFlt  sql.NullFloat64
		Date     time.Time
	}
	form := url.Values{} // no values set
	r, _ := http.NewRequest("POST", "/", nil)
	r.Form = form
	var dest emptyFormStruct
	if err := MapFormToStruct(r, &dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.Str != "" ||
		dest.Int != 0 ||
		dest.Float != 0 ||
		dest.Bool != false ||
		dest.UUID != (uuid.UUID{}) ||
		dest.NullStr.Valid ||
		dest.NullBool.Valid != true ||
		dest.NullBool.Bool != false ||
		dest.NullInt.Valid ||
		dest.NullFlt.Valid ||
		!dest.Date.IsZero() {
		t.Errorf("got %+v, want zero/invalid values", dest)
	}
}

func TestMapFormToStruct_Errors(t *testing.T) {
	type errorFormStruct struct {
		Int   int
		Float float64
		UUID  uuid.UUID
		Date  time.Time
	}
	cases := []struct {
		field, value, wantErr string
	}{
		{"int", "notanint", "invalid integer value"},
		{"float", "notafloat", "invalid float value"},
		{"uuid", "notauuid", "invalid UUID"},
		{"date", "notadate", "invalid date format"},
	}
	for _, tc := range cases {
		form := url.Values{}
		form.Set(tc.field, tc.value)
		r, _ := http.NewRequest("POST", "/", nil)
		r.Form = form
		var dest errorFormStruct
		err := MapFormToStruct(r, &dest)
		if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
			t.Errorf("field %s: got err %v, want %q", tc.field, err, tc.wantErr)
		}
	}
}

func TestMapFormToStruct_AllTypes(t *testing.T) {
	id := uuid.New()
	form := url.Values{}
	form.Set("str", "hello")
	form.Set("int", "42")
	form.Set("int64", "42000")
	form.Set("float32", "3.14")
	form.Set("float64", "2.718")
	form.Set("bool", "true")
	form.Set("uuid", id.String())
	form.Set("nullstr", "foo")
	form.Set("nullbool", "on")
	form.Set("nullint", "123")
	form.Set("nullflt", "1.23")
	form.Set("date", "2023-07-30")
	r, _ := http.NewRequest("POST", "/", nil)
	r.Form = form
	var dest allTypesFormStruct
	if err := MapFormToStruct(r, &dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.Str != "hello" || dest.Int != 42 || dest.Int64 != 42000 || dest.Float32 != float32(3.14) || dest.Float64 != 2.718 || !dest.Bool || dest.UUID != id || !dest.NullStr.Valid || dest.NullStr.String != "foo" || !dest.NullBool.Valid || !dest.NullBool.Bool || !dest.NullInt.Valid || dest.NullInt.Int64 != 123 || !dest.NullFlt.Valid || dest.NullFlt.Float64 != 1.23 || dest.Date.Format("2006-01-02") != "2023-07-30" {
		t.Errorf("got %+v, want correct values", dest)
	}
}

func TestMapFormToStruct(t *testing.T) {
	id := uuid.New()
	form := url.Values{}
	form.Set("name", "Alice")
	form.Set("email", "alice@example.com")
	form.Set("id", id.String())
	r, _ := http.NewRequest("POST", "/", nil)
	r.Form = form
	var dest testFormStruct
	if err := MapFormToStruct(r, &dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.Name != "Alice" || !dest.Email.Valid || dest.Email.String != "alice@example.com" || dest.ID != id {
		t.Errorf("got %+v, want correct values", dest)
	}
}
