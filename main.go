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
	r.Get("/login", h.HandleLogin)
	r.Get("/logout", h.HandleLogout)

	// Static file server for public assets
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Authenticated routes
	r.Group(func(protected chi.Router) {
		protected.Use(handlers.AuthMiddleware(auth))
		protected.Get("/no-access", h.HandleNoAccess)

		// Admin routes
		protected.Group(func(admin chi.Router) {
			admin.Use(handlers.AdminMiddleware(auth))
			admin.Get("/", h.HandleRoot)

			// Server-Sent Events (Powered by Datastar 🚀🚀)
			admin.Get("/sse/invoices", h.HandleSSEGetInvoices)
			admin.Get("/sse/dashboard", h.HandleSSEGetDashboard)
			admin.Get("/sse/customers", h.HandleSSECustomerNav)
			admin.Get("/sse/customer/{id}", h.HandleSSEGetCustomer)
			admin.Get("/sse/customers/add", h.HandleSSEGetAddCustomer)
			admin.Get("/sse/customer/overview/{id}", h.HandleSSECustomerOverview)
			admin.Get("/sse/dashboard/stats", h.HandleSSEDashboardStats)
			admin.Get("/sse/dashboard/activity", h.HandleSSEDashboardActivity)
		})

		// Final catch-all for authenticated routes
		protected.Get("/*", h.HandleNotFound)
	})

	if err := http.ListenAndServe(":1337", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
