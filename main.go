package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/scottmckendry/beam/db"
	"github.com/scottmckendry/beam/ui/views"
)

func main() {
	if err := db.InitialiseDb(); err != nil {
		log.Fatalf("Failed to initialise database: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", handleRoot)

	// static content
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	if err := http.ListenAndServe(":1337", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	views.Login().Render(r.Context(), w)
}
