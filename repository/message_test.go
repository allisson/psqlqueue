package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/allisson/psqlqueue/domain"
)

func pointString(x string) *string {
	return &x
}

func makeMessage(queueID string) *domain.Message {
	return &domain.Message{
		QueueID:    queueID,
		Body:       `{"data": true}`,
		Attributes: map[string]string{"attribute1": "attribute1"},
	}
}

func TestMessage(t *testing.T) {
	cfg := domain.NewConfig()
	ctx := context.Background()
	pool, _ := pgxpool.New(ctx, cfg.TestDatabaseURL)
	defer pool.Close()

	t.Run("Create", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		message := makeMessage(queue.ID)
		message.Enqueue(queue, now)
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message)
		assert.ErrorIs(t, err, domain.ErrMessageAlreadyExists)
	})

	t.Run("Get", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		message := makeMessage(queue.ID)
		message.Enqueue(queue, now)
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message)
		assert.Nil(t, err)

		messageFromDB, err := messageRepo.Get(ctx, message.ID)
		assert.Nil(t, err)
		assert.Equal(t, message.ID, messageFromDB.ID)
		assert.Equal(t, message.Body, messageFromDB.Body)
		assert.Equal(t, message.Attributes, messageFromDB.Attributes)
	})

	t.Run("List", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		message1 := makeMessage(queue.ID)
		message1.Enqueue(queue, now)
		message2 := makeMessage(queue.ID)
		message2.Enqueue(queue, now)
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message1)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message2)
		assert.Nil(t, err)

		messages, err := messageRepo.List(ctx, queue, nil, 10)
		assert.Nil(t, err)
		assert.Len(t, messages, 2)

		messages, err = messageRepo.List(ctx, queue, nil, 10)
		assert.Nil(t, err)
		assert.Len(t, messages, 0)
	})

	t.Run("List with label", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		message1 := makeMessage(queue.ID)
		message1.Label = pointString("label-1")
		message1.Enqueue(queue, now)
		message2 := makeMessage(queue.ID)
		message2.Label = pointString("label-2")
		message2.Enqueue(queue, now)
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message1)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message2)
		assert.Nil(t, err)

		messages, err := messageRepo.List(ctx, queue, pointString("label-1"), 10)
		assert.Nil(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, message1.ID, messages[0].ID)

		messages, err = messageRepo.List(ctx, queue, pointString("label-2"), 10)
		assert.Nil(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, message2.ID, messages[0].ID)
	})

	t.Run("Ack", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		message := makeMessage(queue.ID)
		message.Enqueue(queue, now)
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message)
		assert.Nil(t, err)

		err = messageRepo.Ack(ctx, message.ID)
		assert.Nil(t, err)
	})

	t.Run("Nack", func(t *testing.T) {
		defer clearDatabase(t, ctx, pool)

		now := time.Now().UTC()
		queue := makeQueue("my-queue")
		message := makeMessage(queue.ID)
		message.Enqueue(queue, now)
		queueRepo := NewQueue(pool)
		messageRepo := NewMessage(pool)

		err := queueRepo.Create(ctx, queue)
		assert.Nil(t, err)

		err = messageRepo.Create(ctx, message)
		assert.Nil(t, err)

		err = messageRepo.Nack(ctx, message.ID, uint(0))
		assert.Nil(t, err)
	})
}
