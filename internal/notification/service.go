package notification

type Service interface {
	Send(notification Notification) error
	Validate(notification Notification) (bool, error)
}
