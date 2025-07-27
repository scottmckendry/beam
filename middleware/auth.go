// Package middleware provides HTTP middleware for authentication, authorization, and logging.
package middleware

import (
	"context"
	"net/http"

	"github.com/scottmckendry/beam/oauth"
)

// ContextKey is a custom type for context keys used in middleware.
type ContextKey string

// UserKey is the context key used to store the authenticated user in the request context.
var UserKey ContextKey = "user"

// Auth is middleware that checks for a signed user cookie and injects the user into the request context.
// If the user is not authenticated, it redirects to /login.
func Auth(oauthEnv *oauth.OAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := oauthEnv.GetSignedCookie(r, "user_name")
			if err != nil || user == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
