package notification

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJson = "application/json"

	SendNotificationEndpoint = "/notifications"
)

var supportedChannels = []Channel{EMAIL, SMS, SLACK}

type MessageSendRequest struct {
	Subject   string
	Body      string
	Sender    string
	Recipient string
}

type SendRequest struct {
	Message MessageSendRequest `json:"message"`
	Channel Channel            `json:"channel"`
}

type Handler struct {
}

func (h *Handler) Routes(router *chi.Mux) {
	router.MethodFunc(http.MethodPost, SendNotificationEndpoint, h.Send)
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	var requestBody SendRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Errorf("invalid data: %v", err)
		respond(w, http.StatusBadRequest, nil)
		return
	}

	if !isSupportedChannel(requestBody.Channel) {
		errMsg := fmt.Sprintf("channel [%s] is not supported", requestBody.Channel)
		log.Errorf(errMsg)
		respond(w, http.StatusBadRequest, errMsg)
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

func isSupportedChannel(channel Channel) bool {
	for _, ch := range supportedChannels {
		if ch == channel.ToLower() {
			return true
		}
	}
	return false
}
