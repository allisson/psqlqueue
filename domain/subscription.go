package domain

import (
	"context"
	"slices"
	"time"

	"github.com/jellydator/validation"
	"golang.org/x/exp/maps"
)

// Subscription entity.
type Subscription struct {
	ID             string              `json:"id" db:"id" form:"id"`
	TopicID        string              `json:"topic_id" db:"topic_id" form:"topic_id"`
	QueueID        string              `json:"queue_id" db:"queue_id" form:"queue_id"`
	MessageFilters map[string][]string `json:"message_filters" db:"message_filters" form:"message_filters"`
	CreatedAt      time.Time           `json:"created_at" db:"created_at"`
}

func (s Subscription) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required, validation.Match(idRegex)),
		validation.Field(&s.TopicID, validation.Required, validation.Match(idRegex)),
		validation.Field(&s.QueueID, validation.Required, validation.Match(idRegex)),
	)
}

func (s *Subscription) ShouldCreateMessage(message *Message) bool {
	if len(s.MessageFilters) == 0 {
		return true
	}

	messageAttributesKeys := maps.Keys(message.Attributes)
	messageFiltersKeys := maps.Keys(s.MessageFilters)
	for _, key := range messageFiltersKeys {
		if !slices.Contains(messageAttributesKeys, key) || !slices.Contains(s.MessageFilters[key], message.Attributes[key]) {
			return false
		}
	}

	return true
}

// SubscriptionRepository is the repository interface for the Subscription entity.
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *Subscription) error
	Get(ctx context.Context, id string) (*Subscription, error)
	List(ctx context.Context, offset, limit uint) ([]*Subscription, error)
	ListByTopic(ctx context.Context, topicID string, offset, limit uint) ([]*Subscription, error)
	Delete(ctx context.Context, id string) error
}

// SubscriptionService is the service interface for the Subscription entity.
type SubscriptionService interface {
	Create(ctx context.Context, subscription *Subscription) error
	Get(ctx context.Context, id string) (*Subscription, error)
	List(ctx context.Context, offset, limit uint) ([]*Subscription, error)
	Delete(ctx context.Context, id string) error
}
