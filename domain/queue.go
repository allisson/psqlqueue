package domain

import (
	"context"
	"regexp"
	"time"

	"github.com/jellydator/validation"
)

var (
	idRegex = regexp.MustCompile(`^[a-zA-Z0-9-._]+$`)
)

// Queue entity.
type Queue struct {
	ID                      string    `json:"id" db:"id" form:"id"`
	AckDeadlineSeconds      uint      `json:"ack_deadline_seconds" db:"ack_deadline_seconds" form:"ack_deadline_seconds"`
	MessageRetentionSeconds uint      `json:"message_retention_seconds" db:"message_retention_seconds" form:"message_retention_seconds"`
	DeliveryDelaySeconds    uint      `json:"delivery_delay_seconds" db:"delivery_delay_seconds" form:"delivery_delay_seconds"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}

func (q Queue) Validate() error {
	return validation.ValidateStruct(&q,
		validation.Field(&q.ID, validation.Required, validation.Match(idRegex)),
		validation.Field(&q.AckDeadlineSeconds, validation.Required),
		validation.Field(&q.MessageRetentionSeconds, validation.Required),
	)
}

// QueueStats entity.
type QueueStats struct {
	NumUndeliveredMessages         int `json:"num_undelivered_messages"`
	OldestUnackedMessageAgeSeconds int `json:"oldest_unacked_message_age_seconds"`
}

// QueueRepository is the repository interface for the Queue entity.
type QueueRepository interface {
	Create(ctx context.Context, queue *Queue) error
	Update(ctx context.Context, queue *Queue) error
	Get(ctx context.Context, id string) (*Queue, error)
	List(ctx context.Context, offset, limit int) ([]*Queue, error)
	Delete(ctx context.Context, id string) error
	Stats(ctx context.Context, id string) (*QueueStats, error)
	Purge(ctx context.Context, id string) error
	Cleanup(ctx context.Context, id string) error
}

// QueueService is the service interface for the Queue entity.
type QueueService interface {
	Create(ctx context.Context, queue *Queue) error
	Update(ctx context.Context, queue *Queue) error
	Get(ctx context.Context, id string) (*Queue, error)
	List(ctx context.Context, offset, limit int) ([]*Queue, error)
	Delete(ctx context.Context, id string) error
	Stats(ctx context.Context, id string) (*QueueStats, error)
	Purge(ctx context.Context, id string) error
	Cleanup(ctx context.Context, id string) error
}
