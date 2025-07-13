package handlers

import (
	"bytes"
	"net/http"

	"github.com/a-h/templ"

	"github.com/scottmckendry/beam/ui/views"
)

type NotificationType func(title string, description string) templ.Component

var (
	NotifyInfo    = NotificationType(views.InfoToast)
	NotifyError   = NotificationType(views.ErrorToast)
	NotifySuccess = NotificationType(views.SuccessToast)
)

func (h *Handlers) Notify(
	notifyType NotificationType,
	title string,
	description string,
	w http.ResponseWriter,
	r *http.Request,
) {
	buf := &bytes.Buffer{}
	notification := notifyType(title, description)
	if err := notification.Render(r.Context(), buf); err != nil {
		http.Error(w, "Failed to render notification", http.StatusInternalServerError)
		return
	}
	ServeSSEElement(w, r, buf.String())
}
