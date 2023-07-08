package notification

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJson = "application/json"

	SendNotificationEndpoint = "/notifications"
)

type Handler struct {
}

func (h *Handler) Routes(router *chi.Mux) {
	router.MethodFunc(http.MethodPost, SendNotificationEndpoint, h.Send)
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusAccepted, nil)
}

func respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set(HeaderContentType, MimeApplicationJson)
	w.WriteHeader(code)
	if data != nil {
		resp, _ := json.Marshal(data)
		_, _ = w.Write(resp)
	}
}
