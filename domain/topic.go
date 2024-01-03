package domain

import (
	"context"
	"time"

	"github.com/jellydator/validation"
)

// Topic entity.
type Topic struct {
	ID        string    `json:"id" db:"id" form:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (t Topic) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.ID, validation.Required, validation.Match(idRegex)),
	)
}

// TopicRepository is the repository interface for the Topic entity.
type TopicRepository interface {
	Create(ctx context.Context, topic *Topic) error
	Get(ctx context.Context, id string) (*Topic, error)
	List(ctx context.Context, offset, limit uint) ([]*Topic, error)
	Delete(ctx context.Context, id string) error
}

// TopicService is the service interface for the Topic entity.
type TopicService interface {
	Create(ctx context.Context, topic *Topic) error
	Get(ctx context.Context, id string) (*Topic, error)
	List(ctx context.Context, offset, limit uint) ([]*Topic, error)
	Delete(ctx context.Context, id string) error
	CreateMessage(ctx context.Context, topicID string, message *Message) error
}
