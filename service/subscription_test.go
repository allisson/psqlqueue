package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/allisson/psqlqueue/domain"
	"github.com/allisson/psqlqueue/mocks"
)

func makeSubscription(id, topicID, queueID string) *domain.Subscription {
	return &domain.Subscription{
		ID:        id,
		TopicID:   topicID,
		QueueID:   queueID,
		CreatedAt: time.Now().UTC(),
	}
}

func TestSubscription(t *testing.T) {
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		subscriptionService := NewSubscription(subscriptionRepository)
		subscription := makeSubscription("my-subscription", "my-topic", "my-queue")

		subscriptionRepository.On("Create", ctx, subscription).Return(nil)

		err := subscriptionService.Create(ctx, subscription)
		assert.Nil(t, err)
	})

	t.Run("Create with invalid id", func(t *testing.T) {
		expectedErrorPayload := `{"id":"must be in a valid format"}`
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		subscriptionService := NewSubscription(subscriptionRepository)
		subscription := makeSubscription("my@subscription", "my-topic", "my-queue")

		err := subscriptionService.Create(ctx, subscription)
		assert.NotNil(t, err)
		errorPayload, err := json.Marshal(err)
		assert.Nil(t, err)
		assert.Equal(t, expectedErrorPayload, string(errorPayload))
	})

	t.Run("Get", func(t *testing.T) {
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		subscriptionService := NewSubscription(subscriptionRepository)
		subscription := makeSubscription("my-subscription", "my-topic", "my-queue")

		subscriptionRepository.On("Get", ctx, subscription.ID).Return(subscription, nil)

		_, err := subscriptionService.Get(ctx, subscription.ID)
		assert.Nil(t, err)
	})

	t.Run("List", func(t *testing.T) {
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		subscriptionService := NewSubscription(subscriptionRepository)
		subscription1 := makeSubscription("my-subscription-1", "my-topic-1", "my-queue-1")
		subscription2 := makeSubscription("my-subscription-1", "my-topic-1", "my-queue-2")

		subscriptionRepository.On("List", ctx, uint(0), uint(10)).Return([]*domain.Subscription{subscription1, subscription2}, nil)

		subscriptions, err := subscriptionService.List(ctx, uint(0), uint(10))
		assert.Nil(t, err)
		assert.Len(t, subscriptions, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		subscriptionService := NewSubscription(subscriptionRepository)
		subscription := makeSubscription("my-subscription", "my-topic", "my-queue")

		subscriptionRepository.On("Get", ctx, subscription.ID).Return(subscription, nil)
		subscriptionRepository.On("Delete", ctx, subscription.ID).Return(nil)

		err := subscriptionService.Delete(ctx, subscription.ID)
		assert.Nil(t, err)
	})
}
