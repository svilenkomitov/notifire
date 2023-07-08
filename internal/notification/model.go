package notification

type Channel string

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
	Subject   string
	Body      string
	Sender    string
	Recipient string
}

type Notification struct {
	Message Message
	Channel Channel
	Status  Status
}
