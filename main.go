package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/scottmckendry/beam/db"
	"github.com/scottmckendry/beam/github"
	"github.com/scottmckendry/beam/oauth"
	"github.com/scottmckendry/beam/ui/views"
)

func main() {
	_ = godotenv.Load()
	oauth.InitOAuth()
	if err := db.InitialiseDb(); err != nil {
		log.Fatalf("Failed to initialise database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	oauth.RegisterRoutes(r)

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
	user, err := oauth.GetSignedCookie(r, "user_name")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	tokenCookie, err := r.Cookie("oauth_token")
	if err != nil {
		w.Write([]byte("Logged in as: " + user + " (no token, cannot fetch repo)"))
		return
	}

	ghClient := github.NewClient(tokenCookie.Value)
	repo, err := ghClient.GetRepo(user, "beam")
	if err != nil {
		w.Write([]byte("Logged in as: " + user + " (could not fetch repo)"))
		return
	}

	w.Write(
		[]byte(
			"Logged in as: " + user + "<br>Repo: " + repo.FullName + "<br>Description: " + repo.Description + "<br>Stars: " + fmt.Sprint(
				repo.StargazersCount,
			) + "<br>Forks: " + fmt.Sprint(
				repo.ForksCount,
			),
		),
	)
}
