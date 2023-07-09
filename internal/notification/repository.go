package notification

import (
	"context"
	"github.com/svilenkomitov/notifire/internal/domain"
	"github.com/svilenkomitov/notifire/internal/storage"
)

type Repository interface {
	Create(notification domain.Notification) (string, error)
	UpdateStatus(id string, status domain.Status) error
}

type repository struct {
	db *storage.Database
}

func New(db *storage.Database) Repository {
	return &repository{
		db: db,
	}
}

func (r repository) Create(notification domain.Notification) (string, error) {
	query := `INSERT INTO notifications (channel, status, subject, body, sender, recipient) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	var id string
	err := r.db.QueryRowContext(context.Background(), query,
		notification.Channel, notification.Status, notification.Subject,
		notification.Body, notification.Sender, notification.Recipient).Scan(&id)
	return id, err
}

func (r repository) UpdateStatus(id string, status domain.Status) error {
	query := `UPDATE notifications SET status = $1 WHERE id = $2;`
	_, err := r.db.ExecContext(context.Background(), query, status, id)
	return err
}
