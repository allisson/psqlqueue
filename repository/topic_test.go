package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/allisson/psqlqueue/domain"
)

func makeTopic(topicID string) *domain.Topic {
	return &domain.Topic{ID: topicID, CreatedAt: time.Now().UTC()}
}

func TestTopic(t *testing.T) {
	cfg := domain.NewConfig()
	ctx := context.Background()
	pool, _ := pgxpool.New(ctx, cfg.TestDatabaseURL)
	defer pool.Close()

	t.Run("Create", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topic := makeTopic("my-topic")
		topicRepo := NewTopic(pool)

		err := topicRepo.Create(ctx, topic)
		assert.Nil(t, err)

		err = topicRepo.Create(ctx, topic)
		assert.ErrorIs(t, err, domain.ErrTopicAlreadyExists)
	})

	t.Run("Get", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topicRepo := NewTopic(pool)

		err := topicRepo.Create(ctx, makeTopic("my-topic"))
		assert.Nil(t, err)

		topic, err := topicRepo.Get(ctx, "my-topic")
		assert.Nil(t, err)
		assert.Equal(t, "my-topic", topic.ID)

		_, err = topicRepo.Get(ctx, "not-found-topic")
		assert.ErrorIs(t, err, domain.ErrTopicNotFound)
	})

	t.Run("List", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topicRepo := NewTopic(pool)

		err := topicRepo.Create(ctx, makeTopic("my-topic-1"))
		assert.Nil(t, err)
		err = topicRepo.Create(ctx, makeTopic("my-topic-2"))
		assert.Nil(t, err)

		topics, err := topicRepo.List(ctx, uint(0), uint(10))
		assert.Nil(t, err)
		assert.Len(t, topics, 2)
		assert.Equal(t, "my-topic-1", topics[0].ID)
		assert.Equal(t, "my-topic-2", topics[1].ID)
	})

	t.Run("Delete", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		topicRepo := NewTopic(pool)

		err := topicRepo.Create(ctx, makeTopic("my-topic"))
		assert.Nil(t, err)

		err = topicRepo.Delete(ctx, "my-topic")
		assert.Nil(t, err)

		_, err = topicRepo.Get(ctx, "my-topic")
		assert.ErrorIs(t, err, domain.ErrTopicNotFound)
	})
}
