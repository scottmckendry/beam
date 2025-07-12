package handlers

import (
	"github.com/scottmckendry/beam/ui/views"
	"net/http"
)

// HandleLogin processes GET requests to the login page and redirects authenticated users to the root.
func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("user_name")
	if err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	views.Login().Render(r.Context(), w)
}

// HandleLogout clears the user session and redirects to the login page.
func (h *Handlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(
		w,
		&http.Cookie{Name: "user_name", Value: "", Path: "/", HttpOnly: true, MaxAge: -1},
	)
	http.Redirect(w, r, "/login", http.StatusFound)
}
