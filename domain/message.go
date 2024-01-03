package domain

import (
	"context"
	"time"

	"github.com/jellydator/validation"
	"github.com/oklog/ulid/v2"
)

// Message entity.
type Message struct {
	ID               string            `json:"id" db:"id"`
	QueueID          string            `json:"queue_id" db:"queue_id"`
	Label            *string           `json:"label" db:"label" form:"label"`
	Body             string            `json:"body" db:"body" form:"body"`
	Attributes       map[string]string `json:"attributes" db:"attributes" form:"attributes"`
	DeliveryAttempts uint              `json:"delivery_attempts" db:"delivery_attempts"`
	ExpiredAt        time.Time         `json:"-" db:"expired_at"`
	ScheduledAt      time.Time         `json:"-" db:"scheduled_at"`
	CreatedAt        time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time         `json:"-" db:"updated_at"`
}

func (m Message) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Body, validation.Required),
	)
}

func (m *Message) Enqueue(queue *Queue, now time.Time) {
	scheduledAt := now
	if queue.DeliveryDelaySeconds > 0 {
		scheduledAt = scheduledAt.Add(time.Duration(queue.DeliveryDelaySeconds) * time.Second)
	}

	m.ID = ulid.Make().String()
	m.QueueID = queue.ID
	m.DeliveryAttempts = 0
	m.ExpiredAt = now.Add(time.Duration(queue.MessageRetentionSeconds) * time.Second)
	m.ScheduledAt = scheduledAt
	m.CreatedAt = now
	m.UpdatedAt = now
}

func (m *Message) DeliverySetup(queue *Queue, now time.Time) {
	m.DeliveryAttempts = m.DeliveryAttempts + 1
	m.ScheduledAt = now.Add(time.Duration(queue.AckDeadlineSeconds) * time.Second)
	m.UpdatedAt = now
}

func (m *Message) Ack(now time.Time) {
	m.ExpiredAt = now
	m.UpdatedAt = now
}

func (m *Message) Nack(now time.Time, visibilityTimeoutSeconds uint) {
	m.ScheduledAt = now.Add(time.Duration(visibilityTimeoutSeconds) * time.Second)
	m.UpdatedAt = now
}

// MessageRepository is the repository interface for the Message entity.
type MessageRepository interface {
	CreateMany(ctx context.Context, messages []*Message) error
	Create(ctx context.Context, message *Message) error
	Get(ctx context.Context, id string) (*Message, error)
	List(ctx context.Context, queue *Queue, label *string, limit uint) ([]*Message, error)
	Ack(ctx context.Context, id string) error
	Nack(ctx context.Context, id string, visibilityTimeoutSeconds uint) error
}

// MessageService is the service interface for the Message entity.
type MessageService interface {
	Create(ctx context.Context, message *Message) error
	List(ctx context.Context, queueID string, label *string, limit uint) ([]*Message, error)
	Ack(ctx context.Context, id string) error
	Nack(ctx context.Context, id string, visibilityTimeoutSeconds uint) error
}
