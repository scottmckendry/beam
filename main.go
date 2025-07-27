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
			admin.Get("/sse/invoice", h.InvoicesSSE)
			admin.Get("/sse/dashboard", h.DashboardSSE)
			admin.Get("/sse/customer/{id}", h.GetCustomerSSE)
			admin.Get("/sse/customer/overview/{id}", h.GetCustomerOverviewSSE)
			admin.Get("/sse/customer/contacts/{id}", h.GetCustomerContactsSSE)
			admin.Get("/sse/customer/subscriptions/{id}", h.GetCustomerSubscriptionsSSE)
			admin.Get("/sse/customer/projects/{id}", h.GetCustomerProjectsSSE)
			admin.Get("/sse/customer/add", h.AddCustomerSSE)
			admin.Get("/sse/customer/add-submit", h.SubmitAddCustomerSSE)
			admin.Get("/sse/customer/delete/{id}", h.DeleteCustomerSSE)
			admin.Get("/sse/customer/edit/{id}", h.EditCustomerFormSSE)
			admin.Get("/sse/customer/edit-submit/{id}", h.EditCustomerSubmitSSE)
			admin.Post("/sse/customer/upload-logo/{id}", h.UploadCustomerLogoSSE)
			admin.Get("/sse/customer/delete-logo/{id}", h.DeleteCustomerLogoSSE)
			admin.Get("/sse/customer/{customerID}/add-contact", h.AddContactFormSSE)
			admin.Get("/sse/customer/{customerID}/add-contact-submit", h.AddContactSubmitSSE)
			admin.Get("/sse/customer/{customerID}/edit-contact/{contactID}", h.EditContactFormSSE)
			admin.Get("/sse/customer/{customerID}/edit-contact-submit/{contactID}", h.EditContactSubmitSSE)
			admin.Get("/sse/customer/{customerID}/delete-contact/{contactID}", h.DeleteContactSSE)
			admin.Get("/sse/dashboard/stats", h.DashboardStatsSSE)
			admin.Get("/sse/dashboard/activity", h.DashboardActivitySSE)
		})

		// Final catch-all for authenticated routes
		protected.Get("/*", h.HandleNotFound)
	})

	if err := http.ListenAndServe(":1337", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
