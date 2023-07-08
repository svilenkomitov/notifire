package notification

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"github.com/svilenkomitov/notifire/internal/domain"
	"github.com/svilenkomitov/notifire/internal/mq"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJson = "application/json"

	SendNotificationEndpoint = "/notifications"
)

var supportedChannels = []domain.Channel{domain.EMAIL, domain.SMS, domain.SLACK}

type Handler struct {
	MQService mq.Service
}

func (h *Handler) Routes(router *chi.Mux) {
	router.MethodFunc(http.MethodPost, SendNotificationEndpoint, h.Send)
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	var notification domain.Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		log.Errorf("invalid data: %v", err)
		respond(w, http.StatusBadRequest, nil)
		return
	}

	if !isSupportedChannel(notification.Channel) {
		errMsg := fmt.Sprintf("channel [%s] is not supported", notification.Channel)
		log.Errorf(errMsg)
		respond(w, http.StatusBadRequest, errMsg)
		return
	}

	err = h.MQService.Publish(notification)
	if err != nil {
		log.Errorf("failed to send notification: %v", err)
		respond(w, http.StatusAccepted, nil)
		return
	}

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

func isSupportedChannel(channel domain.Channel) bool {
	for _, ch := range supportedChannels {
		if ch == channel.ToLower() {
			return true
		}
	}
	return false
}
