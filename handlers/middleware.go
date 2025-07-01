package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/scottmckendry/beam/oauth"
)

type contextKey string

const (
	userKey  contextKey = "user"
	adminKey contextKey = "admin"
)

func AuthMiddleware(oauthEnv *oauth.OAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := oauthEnv.GetSignedCookie(r, "user_name")
			if err != nil || user == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminMiddleware(oauthEnv *oauth.OAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(userKey).(string)
			if !ok || user == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			isAdmin, err := oauthEnv.DB.IsUserAdmin(r.Context(), user)
			if err != nil || !isAdmin {
				log.Printf("User %s is not an admin or an error occurred: %v", user, err)
				http.Redirect(w, r, "/no-access", http.StatusFound)
				return
			}

			ctx := context.WithValue(r.Context(), adminKey, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
