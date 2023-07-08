package notification

type Service interface {
	Send(notification Notification) (Notification, error)
	Validate(notification Notification) *ValidationError
}
