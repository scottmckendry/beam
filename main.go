package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	beamDb "github.com/scottmckendry/beam/db"
	"github.com/scottmckendry/beam/db/sqlc"
	"github.com/scottmckendry/beam/github"
	"github.com/scottmckendry/beam/oauth"
	"github.com/scottmckendry/beam/ui/views"
)

var queries *db.Queries

func main() {
	_ = godotenv.Load()
	oauth.InitOAuth()
	dbConn, queries, err := beamDb.InitialiseDB()
	if err != nil {
		log.Fatalf("Failed to initialise database: %v", err)
	}
	defer dbConn.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	oauth.RegisterRoutes(r, queries)

	r.Get("/login", handleLogin)
	r.Get("/logout", handleLogout)
	r.Get("/", handleRoot)

	// static content
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	if err := http.ListenAndServe(":1337", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// If already authenticated, redirect to /
	_, err := r.Cookie("user_name")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	views.Login().Render(r.Context(), w)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(
		w,
		&http.Cookie{Name: "user_name", Value: "", Path: "/", HttpOnly: true, MaxAge: -1},
	)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := oauth.GetSignedCookie(r, "user_name")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tokenCookie, err := r.Cookie("oauth_token")
	if err != nil {
		views.Root(user, false, "", "", 0, 0).Render(ctx, w)
		return
	}

	ghClient := github.NewClient(tokenCookie.Value)
	repo, err := ghClient.GetRepo(user, "beam")
	if err != nil {
		views.Root(user, false, "", "", 0, 0).Render(ctx, w)
		return
	}

	isAdmin, err := queries.IsUserAdmin(ctx, user)
	if err != nil {
		views.Root(user, false, repo.FullName, repo.Description, repo.StargazersCount, repo.ForksCount).
			Render(ctx, w)
		return
	}

	views.Root(user, isAdmin, repo.FullName, repo.Description, repo.StargazersCount, repo.ForksCount).
		Render(ctx, w)
}
