package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/allisson/psqlqueue/domain"
)

func clearDatabase(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	sqlQuery := `
	DELETE FROM queues;
	DELETE FROM topics;
	`
	_, err := pool.Exec(ctx, sqlQuery)
	assert.Nil(t, err)
}

func makeQueue(queueID string) *domain.Queue {
	return &domain.Queue{
		ID:                      queueID,
		AckDeadlineSeconds:      60,
		MessageRetentionSeconds: 3600,
		DeliveryDelaySeconds:    0,
		CreatedAt:               time.Now().UTC(),
		UpdatedAt:               time.Now().UTC(),
	}
}

func TestQueue(t *testing.T) {
	cfg := domain.NewConfig()
	ctx := context.Background()
	pool, _ := pgxpool.New(ctx, cfg.TestDatabaseURL)
	defer pool.Close()

	t.Run("Create", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		queue := makeQueue("my-queue")
		queueRepo := NewQueue(pool)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = queueRepo.Create(ctx, queue)
		assert.ErrorIs(t, err, domain.ErrQueueAlreadyExists)
	})

	t.Run("Update", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		queue := makeQueue("my-queue")
		queueRepo := NewQueue(pool)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		queue.AckDeadlineSeconds = 30
		queue.MessageRetentionSeconds = 600
		queue.DeliveryDelaySeconds = 10

		err = queueRepo.Update(ctx, queue)
		assert.Nil(t, err)

		queueFromDB, err := queueRepo.Get(ctx, queue.ID)
		assert.Nil(t, err)
		assert.Equal(t, uint(30), queueFromDB.AckDeadlineSeconds)
		assert.Equal(t, uint(600), queueFromDB.MessageRetentionSeconds)
		assert.Equal(t, uint(10), queueFromDB.DeliveryDelaySeconds)
	})

	t.Run("Get", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		queueRepo := NewQueue(pool)

		err := queueRepo.Create(ctx, makeQueue("my-queue"))
		assert.Nil(t, err)

		queue, err := queueRepo.Get(ctx, "my-queue")
		assert.Nil(t, err)
		assert.Equal(t, "my-queue", queue.ID)

		_, err = queueRepo.Get(ctx, "not-found-queue")
		assert.ErrorIs(t, err, domain.ErrQueueNotFound)
	})

	t.Run("List", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		queueRepo := NewQueue(pool)

		err := queueRepo.Create(ctx, makeQueue("my-queue-1"))
		assert.Nil(t, err)
		err = queueRepo.Create(ctx, makeQueue("my-queue-2"))
		assert.Nil(t, err)

		queues, err := queueRepo.List(ctx, 0, 10)
		assert.Nil(t, err)
		assert.Len(t, queues, 2)
		assert.Equal(t, "my-queue-1", queues[0].ID)
		assert.Equal(t, "my-queue-2", queues[1].ID)
	})

	t.Run("Delete", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		queueRepo := NewQueue(pool)

		err := queueRepo.Create(ctx, makeQueue("my-queue"))
		assert.Nil(t, err)

		err = queueRepo.Delete(ctx, "my-queue")
		assert.Nil(t, err)

		_, err = queueRepo.Get(ctx, "my-queue")
		assert.ErrorIs(t, err, domain.ErrQueueNotFound)
	})

	t.Run("Stats", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)
		message := makeMessage(queue.ID)
		message.Enqueue(queue, now)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message)
		assert.Nil(t, err)

		time.Sleep(1 * time.Second)

		stats, err := queueRepo.Stats(ctx, queue.ID)
		assert.Nil(t, err)
		assert.Equal(t, uint(1), stats.NumUndeliveredMessages)
		assert.Equal(t, uint(1), stats.OldestUnackedMessageAgeSeconds)
	})

	t.Run("Purge", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)
		message := makeMessage(queue.ID)
		message.Enqueue(queue, now)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message)
		assert.Nil(t, err)

		err = queueRepo.Purge(ctx, queue.ID)
		assert.Nil(t, err)

		stats, err := queueRepo.Stats(ctx, queue.ID)
		assert.Nil(t, err)
		assert.Equal(t, uint(0), stats.NumUndeliveredMessages)
		assert.Equal(t, uint(0), stats.OldestUnackedMessageAgeSeconds)
	})

	t.Run("Cleanup", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)
		message := makeMessage(queue.ID)
		message.Enqueue(queue, now)
		message.ExpiredAt = now

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message)
		assert.Nil(t, err)

		err = queueRepo.Cleanup(ctx, queue.ID)
		assert.Nil(t, err)
	})
}
