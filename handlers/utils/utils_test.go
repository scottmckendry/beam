package utils

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
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
