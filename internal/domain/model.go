package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Channel string

func (c Channel) ToLower() Channel {
	return Channel(strings.ToLower(string(c)))
}

const (
	EMAIL Channel = "email"
	SMS   Channel = "sms"
	SLACK Channel = "slack"
)

type Status string

const (
	Pending Status = "pending"
	Sent    Status = "sent"
	Failed  Status = "failed"
)

type Message struct {
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type Notification struct {
	Message Message `json:"message"`
	Channel Channel `json:"channel"`
	Status  Status  `json:"status"`
}

func (n Notification) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
}

func Unmarshal(data []byte) (Notification, error) {
	var notification Notification
	if err := json.Unmarshal(data, &notification); err != nil {
		return Notification{}, errors.New(fmt.Sprintf("failed to unmarshal: %v", err))
	}
	return notification, nil
}
