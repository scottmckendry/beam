// Package oauth provides GitHub OAuth2 authentication and secure cookie utilities for Beam.
package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOauthConfig *oauth2.Config
	hashKey           []byte
	blockKey          []byte
	sc                *securecookie.SecureCookie
)

// InitOAuth initializes OAuth2 config and secure cookie keys from environment variables.
func InitOAuth() {
	githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
		RedirectURL:  os.Getenv("GITHUB_CALLBACK_URL"),
		Scopes:       []string{"read:user", "user:email"},
	}
	hashKey = []byte(os.Getenv("COOKIE_HASH_KEY"))
	blockKey = []byte(os.Getenv("COOKIE_BLOCK_KEY"))
	sc = securecookie.New(hashKey, blockKey)
}

// SetSignedCookie encodes and sets a signed, HTTP-only cookie.
func SetSignedCookie(w http.ResponseWriter, name, value string) {
	encoded, err := sc.Encode(name, value)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{Name: name, Value: encoded, Path: "/", HttpOnly: true})
}

// GetSignedCookie retrieves and decodes a signed cookie value.
func GetSignedCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	var value string
	if err = sc.Decode(name, cookie.Value, &value); err != nil {
		return "", err
	}
	return value, nil
}

// RegisterRoutes registers OAuth login and callback routes on the given router.
func RegisterRoutes(r chi.Router) {
	r.Get(
		"/login/github",
		func(w http.ResponseWriter, r *http.Request) { githubLoginHandler(w, r) },
	)
	r.Get(
		"/auth/github/callback",
		func(w http.ResponseWriter, r *http.Request) { githubCallbackHandler(w, r) },
	)
}

// githubLoginHandler starts the GitHub OAuth2 login flow.
func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	http.SetCookie(w, &http.Cookie{Name: "oauthstate", Value: state, Path: "/", HttpOnly: true})
	url := githubOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// githubCallbackHandler handles the GitHub OAuth2 callback and user authentication.
func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("oauthstate")
	if err != nil || r.URL.Query().Get("state") != state.Value {
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}
	code := r.URL.Query().Get("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "OAuth token exchange failed", http.StatusInternalServerError)
		return
	}
	client := githubOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}
	SetSignedCookie(w, "user_name", user.Login)
	http.Redirect(w, r, "/", http.StatusFound)
}

// generateState generates a random state string for OAuth2 CSRF protection.
func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
