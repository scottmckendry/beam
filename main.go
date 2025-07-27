package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"

	"github.com/scottmckendry/beam/db"
	"github.com/scottmckendry/beam/handlers"
	middlewares "github.com/scottmckendry/beam/middleware"
	"github.com/scottmckendry/beam/oauth"
)

func main() {
	_ = godotenv.Load()
	initLogger()

	dbConn, queries, err := db.InitialiseDB()
	if err != nil {
		slog.Error("Failed to initialize database", "err", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	auth := oauth.New(queries)
	h := handlers.New(queries, auth)

	r := chi.NewRouter()

	// Use a custom logger for middleware with AddSource: false
	middlewareLogger := newLogger(false)
	r.Use(middlewares.Slog(middlewareLogger))
	r.Use(middleware.Recoverer)

	auth.RegisterRoutes(r)

	// Public routes
	r.Get("/login", h.HandleLogin)
	r.Get("/logout", h.HandleLogout)

	// Static file server for public assets
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Authenticated routes
	r.Group(func(protected chi.Router) {
		protected.Use(middlewares.Auth(auth))
		protected.Get("/no-access", h.HandleNoAccess)

		// Admin routes
		protected.Group(func(admin chi.Router) {
			admin.Use(middlewares.Admin(auth))
			h.RegisterRootRoutes(admin)
			h.RegisterInvoiceRoutes(admin)
			h.RegisterDashboardRoutes(admin)
			h.RegisterCustomerRoutes(admin)
			h.RegisterContactRoutes(admin)
		})

		// Final catch-all for authenticated routes
		protected.Get("/*", h.HandleNotFound)
	})

	if err := http.ListenAndServe(":1337", r); err != nil {
		slog.Error("Failed to start server", "err", err)
		os.Exit(1)
	}
}

// newLogger returns a slog.Logger with the given AddSource option and LOG_FORMAT env.
func newLogger(addSource bool) *slog.Logger {
	logFormat := os.Getenv("LOG_FORMAT")
	var handler slog.Handler
	if logFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: addSource})
	} else {
		handler = tint.NewHandler(os.Stdout, &tint.Options{AddSource: addSource})
	}
	return slog.New(handler)
}

// initLogger configures slog with either tint (pretty) or JSON handler based on LOG_FORMAT.
func initLogger() {
	slog.SetDefault(newLogger(true))
	slog.Info("Logger initialized", "format", os.Getenv("LOG_FORMAT"))
}
