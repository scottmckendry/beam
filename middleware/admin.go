// Package middleware provides HTTP middleware for authentication, authorization, and logging.
package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/scottmckendry/beam/oauth"
)

// AdminKey is the context key used to indicate admin status in the request context.
var AdminKey ContextKey = "admin"

// Admin is middleware that restricts access to admin users only.
// It checks if the user is authenticated and an admin, otherwise redirects to /login or /no-access.
func Admin(oauthEnv *oauth.OAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserKey).(string)
			if !ok || user == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			isAdmin, err := oauthEnv.DB.IsUserAdmin(r.Context(), user)
			if err != nil || !isAdmin {
				slog.Error("User is not an admin or error occurred", "user", user, "err", err)
				http.Redirect(w, r, "/no-access", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), AdminKey, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
