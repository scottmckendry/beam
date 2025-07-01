package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

func RegisterRoutes(r chi.Router, routes []Route) {
	for _, route := range routes {
		switch route.Method {
		case "GET":
			r.Get(route.Pattern, route.Handler)
		case "POST":
			r.Post(route.Pattern, route.Handler)
		case "PUT":
			r.Put(route.Pattern, route.Handler)
		case "DELETE":
			r.Delete(route.Pattern, route.Handler)
		}
	}
}
