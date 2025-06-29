// Package oauth provides GitHub OAuth2 authentication and secure cookie utilities for Beam.
package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/scottmckendry/beam/db/sqlc"
)

type OAuth struct {
	OauthConfig  *oauth2.Config
	SecureCookie *securecookie.SecureCookie
	DB           *db.Queries
}

// New creates a new OAuthEnv instance using environment variables for configuration.
func New(db *db.Queries) *OAuth {
	const githubOAuthScopes = "read:user,user:email,repo"
	config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
		RedirectURL:  os.Getenv("GITHUB_CALLBACK_URL"),
		Scopes:       strings.Split(githubOAuthScopes, ","),
	}
	hashKey := []byte(os.Getenv("COOKIE_HASH_KEY"))
	blockKey := []byte(os.Getenv("COOKIE_BLOCK_KEY"))
	sc := securecookie.New(hashKey, blockKey)
	return &OAuth{
		OauthConfig:  config,
		SecureCookie: sc,
		DB:           db,
	}
}

// SetSignedCookie encodes and sets a signed, HTTP-only cookie.
func (env *OAuth) SetSignedCookie(w http.ResponseWriter, name, value string) {
	encoded, err := env.SecureCookie.Encode(name, value)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{Name: name, Value: encoded, Path: "/", HttpOnly: true})
}

// GetSignedCookie retrieves and decodes a signed cookie value.
func (env *OAuth) GetSignedCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	var value string
	if err = env.SecureCookie.Decode(name, cookie.Value, &value); err != nil {
		return "", err
	}
	return value, nil
}

// RegisterRoutes registers OAuth login and callback routes on the given router.
func (env *OAuth) RegisterRoutes(r chi.Router) {
	r.Get(
		"/login/github",
		env.githubLoginHandler(),
	)
	r.Get(
		"/auth/github/callback",
		env.githubCallbackHandler(),
	)
}

// githubLoginHandler starts the GitHub OAuth2 login flow.
func (env *OAuth) githubLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := generateState()
		http.SetCookie(w, &http.Cookie{Name: "oauthstate", Value: state, Path: "/", HttpOnly: true})
		url := env.OauthConfig.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// githubCallbackHandler handles the GitHub OAuth2 callback and user authentication.
func (env *OAuth) githubCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		state, err := r.Cookie("oauthstate")
		if err != nil || r.URL.Query().Get("state") != state.Value {
			http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		token, err := env.OauthConfig.Exchange(ctx, code)
		if err != nil {
			http.Error(w, "OAuth token exchange failed", http.StatusInternalServerError)
			return
		}
		client := env.OauthConfig.Client(ctx, token)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil || resp.StatusCode != 200 {
			http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		var user struct {
			ID    string `json:"login"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}
		_ = env.DB.InsertUser(ctx, db.InsertUserParams{
			Email:    user.Email,
			GithubID: fmt.Sprint(user.ID),
		})
		env.SetSignedCookie(w, "user_name", user.ID)
		http.SetCookie(
			w,
			&http.Cookie{Name: "oauth_token", Value: token.AccessToken, Path: "/", HttpOnly: true},
		)
		http.Redirect(w, r, "/", http.StatusFound)
	}
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
