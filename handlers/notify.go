package handlers

import (
	"github.com/scottmckendry/beam/handlers/utils"
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

func (h *Handlers) Notify(notifyType NotificationType, title string, description string, w http.ResponseWriter, r *http.Request) {
	notification := notifyType(title, description)
	utils.RenderSSE(w, r, utils.SSEOpts{Views: []templ.Component{notification}})
}
