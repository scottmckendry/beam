package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
)

type mockDB struct{}

func (m *mockDB) InsertUser(ctx context.Context, p any) error { return nil }

func newTestEnv() *OAuth {
	os.Setenv("COOKIE_HASH_KEY", "testhashkeytesthashkeytesthashkeytesth")
	os.Setenv("COOKIE_BLOCK_KEY", "testblockkeytestblockkeytestblockkeytes")
	cfg := &oauth2.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		Endpoint:     oauth2.Endpoint{AuthURL: "http://auth", TokenURL: "http://token"},
		RedirectURL:  "http://cb",
		Scopes:       []string{"read:user"},
	}
	sc := newSecureCookieFromEnv()
	return &OAuth{OauthConfig: cfg, SecureCookie: sc, DB: nil}
}

func TestRegisterRoutes(t *testing.T) {
	env := newTestEnv()
	r := chi.NewRouter()
	env.RegisterRoutes(r)
	req := httptest.NewRequest("GET", "/login/github", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Errorf("expected redirect, got %d", rec.Code)
	}
}

func TestLoginHandlerSetsStateCookie(t *testing.T) {
	env := newTestEnv()
	r := chi.NewRouter()
	env.RegisterRoutes(r)
	req := httptest.NewRequest("GET", "/login/github", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	cookies := rec.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "oauthstate" {
			found = true
		}
	}
	if !found {
		t.Error("oauthstate cookie not set")
	}
}

func TestCallbackHandlerInvalidState(t *testing.T) {
	env := newTestEnv()
	r := chi.NewRouter()
	env.RegisterRoutes(r)
	req := httptest.NewRequest("GET", "/auth/github/callback?state=bad", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid state, got %d", rec.Code)
	}
}

func newSecureCookieFromEnv() *securecookie.SecureCookie {
	hashKey := []byte(os.Getenv("COOKIE_HASH_KEY"))
	blockKey := []byte(os.Getenv("COOKIE_BLOCK_KEY"))
	return securecookie.New(hashKey, blockKey)
}

// Add more tests for token exchange, user info, and DB errors by mocking OauthConfig and DB as needed.
