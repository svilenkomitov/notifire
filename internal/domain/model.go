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

type Notification struct {
	ID        int     `json:"id" db:"id"`
	Subject   string  `json:"subject" db:"subject"`
	Body      string  `json:"body" db:"body"`
	Sender    string  `json:"sender" db:"sender"`
	Recipient string  `json:"recipient" db:"recipient"`
	Channel   Channel `json:"channel" db:"channel"`
	Status    Status  `json:"status" db:"status"`
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
