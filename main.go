package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/scottmckendry/beam/db"
	"github.com/scottmckendry/beam/handlers"
	"github.com/scottmckendry/beam/oauth"
)

func main() {
	_ = godotenv.Load()

	dbConn, queries, err := db.InitialiseDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	auth := oauth.New(queries)
	h := handlers.New(queries, auth)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	auth.RegisterRoutes(r)

	// Public routes
	handlers.RegisterRoutes(r, []handlers.Route{
		{Method: "GET", Pattern: "/login", Handler: h.HandleLogin},
		{Method: "GET", Pattern: "/logout", Handler: h.HandleLogout},
	})

	// Authenticated routes
	r.Group(func(protected chi.Router) {
		protected.Use(handlers.AuthMiddleware(auth))
		handlers.RegisterRoutes(protected, []handlers.Route{
			{Method: "GET", Pattern: "/no-access", Handler: h.HandleNoAccess},
		})

		// Admin routes
		protected.Group(func(admin chi.Router) {
			admin.Use(handlers.AdminMiddleware(auth))
			admin.Get("/", h.HandleRoot)
			admin.Get("/dashboard", h.HandleDashboard)
			admin.Get("/invoices", h.HandleInvoices)
			admin.Get("/customer/{id}", h.HandleCustomer)
		})
	})

	// Static file server for public assets
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	if err := http.ListenAndServe(":1337", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
