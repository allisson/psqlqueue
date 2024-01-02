package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/allisson/psqlqueue/domain"
)

func makeSubscription(id, topicID, queueID string, messageFilters map[string][]string) *domain.Subscription {
	return &domain.Subscription{
		ID:             id,
		TopicID:        topicID,
		QueueID:        queueID,
		MessageFilters: messageFilters,
		CreatedAt:      time.Now().UTC(),
	}
}

func TestSubscription(t *testing.T) {
	cfg := domain.NewConfig()
	ctx := context.Background()
	pool, _ := pgxpool.New(ctx, cfg.TestDatabaseURL)
	defer pool.Close()

	t.Run("Create", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topic := makeTopic("my-topic")
		topicRepo := NewTopic(pool)
		queue := makeQueue("my-queue")
		queueRepo := NewQueue(pool)
		subscription := makeSubscription("my-subscription", topic.ID, queue.ID, nil)
		subscriptionRepo := NewSubscription(pool)

		err := topicRepo.Create(ctx, topic)
		assert.Nil(t, err)

		err = queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = subscriptionRepo.Create(ctx, subscription)
		assert.Nil(t, err)

		err = subscriptionRepo.Create(ctx, subscription)
		assert.ErrorIs(t, err, domain.ErrSubscriptionAlreadyExists)

		subscription.ID = "another-subscription-with-same-topic-id-and-queue-id"
		err = subscriptionRepo.Create(ctx, subscription)
		assert.ErrorIs(t, err, domain.ErrSubscriptionAlreadyExists)
	})

	t.Run("Get", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topic := makeTopic("my-topic")
		topicRepo := NewTopic(pool)
		queue := makeQueue("my-queue")
		queueRepo := NewQueue(pool)
		subscriptionRepo := NewSubscription(pool)

		err := topicRepo.Create(ctx, topic)
		assert.Nil(t, err)

		err = queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = subscriptionRepo.Create(ctx, makeSubscription("my-subscription", topic.ID, queue.ID, nil))
		assert.Nil(t, err)

		subscription, err := subscriptionRepo.Get(ctx, "my-subscription")
		assert.Nil(t, err)
		assert.Equal(t, "my-subscription", subscription.ID)

		_, err = subscriptionRepo.Get(ctx, "not-found-subscription")
		assert.ErrorIs(t, err, domain.ErrSubscriptionNotFound)
	})

	t.Run("List", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topic1 := makeTopic("my-topic-1")
		topic2 := makeTopic("my-topic-2")
		topicRepo := NewTopic(pool)
		queue1 := makeQueue("my-queue-1")
		queue2 := makeQueue("my-queue-2")
		queueRepo := NewQueue(pool)
		subscriptionRepo := NewSubscription(pool)

		err := topicRepo.Create(ctx, topic1)
		assert.Nil(t, err)
		err = topicRepo.Create(ctx, topic2)
		assert.Nil(t, err)

		err = queueRepo.Create(ctx, queue1)
		assert.Nil(t, err)
		err = queueRepo.Create(ctx, queue2)
		assert.Nil(t, err)

		err = subscriptionRepo.Create(ctx, makeSubscription("my-subscription-1", topic1.ID, queue1.ID, nil))
		assert.Nil(t, err)
		err = subscriptionRepo.Create(ctx, makeSubscription("my-subscription-2", topic2.ID, queue2.ID, nil))
		assert.Nil(t, err)

		subscriptions, err := subscriptionRepo.List(ctx, uint(0), uint(10))
		assert.Nil(t, err)
		assert.Len(t, subscriptions, 2)
		assert.Equal(t, "my-subscription-1", subscriptions[0].ID)
		assert.Equal(t, "my-subscription-2", subscriptions[1].ID)
	})

	t.Run("ListByTopic", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topic1 := makeTopic("my-topic-1")
		topic2 := makeTopic("my-topic-2")
		topicRepo := NewTopic(pool)
		queue1 := makeQueue("my-queue-1")
		queue2 := makeQueue("my-queue-2")
		queueRepo := NewQueue(pool)
		subscriptionRepo := NewSubscription(pool)

		err := topicRepo.Create(ctx, topic1)
		assert.Nil(t, err)
		err = topicRepo.Create(ctx, topic2)
		assert.Nil(t, err)

		err = queueRepo.Create(ctx, queue1)
		assert.Nil(t, err)
		err = queueRepo.Create(ctx, queue2)
		assert.Nil(t, err)

		err = subscriptionRepo.Create(ctx, makeSubscription("my-subscription-1", topic1.ID, queue1.ID, nil))
		assert.Nil(t, err)
		err = subscriptionRepo.Create(ctx, makeSubscription("my-subscription-2", topic2.ID, queue2.ID, nil))
		assert.Nil(t, err)

		subscriptions, err := subscriptionRepo.ListByTopic(ctx, topic1.ID, uint(0), uint(10))
		assert.Nil(t, err)
		assert.Len(t, subscriptions, 1)
		assert.Equal(t, "my-subscription-1", subscriptions[0].ID)

		subscriptions, err = subscriptionRepo.ListByTopic(ctx, topic2.ID, uint(0), uint(10))
		assert.Nil(t, err)
		assert.Len(t, subscriptions, 1)
		assert.Equal(t, "my-subscription-2", subscriptions[0].ID)
	})

	t.Run("Delete", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topic := makeTopic("my-topic")
		topicRepo := NewTopic(pool)
		queue := makeQueue("my-queue")
		queueRepo := NewQueue(pool)
		subscriptionRepo := NewSubscription(pool)

		err := topicRepo.Create(ctx, topic)
		assert.Nil(t, err)

		err = queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = subscriptionRepo.Create(ctx, makeSubscription("my-subscription", topic.ID, queue.ID, nil))
		assert.Nil(t, err)

		err = subscriptionRepo.Delete(ctx, "my-subscription")
		assert.Nil(t, err)

		_, err = subscriptionRepo.Get(ctx, "my-subscription")
		assert.ErrorIs(t, err, domain.ErrSubscriptionNotFound)
	})
}
