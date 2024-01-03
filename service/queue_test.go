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
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		queueRepository := mocks.NewQueueRepository(t)
		queueService := NewQueue(queueRepository)
		queue := makeQueue("my-queue")

		queueRepository.On("Create", ctx, queue).Return(nil)

		err := queueService.Create(ctx, queue)
		assert.Nil(t, err)
	})

	t.Run("Create with invalid queue", func(t *testing.T) {
		expectedErrorPayload := `{"id":"must be in a valid format"}`
		queueRepository := mocks.NewQueueRepository(t)
		queueService := NewQueue(queueRepository)
		queue := makeQueue("my@queue")

		err := queueService.Create(ctx, queue)
		assert.NotNil(t, err)
		errorPayload, err := json.Marshal(err)
		assert.Nil(t, err)
		assert.Equal(t, expectedErrorPayload, string(errorPayload))
	})

	t.Run("Update", func(t *testing.T) {
		queueRepository := mocks.NewQueueRepository(t)
		queueService := NewQueue(queueRepository)
		queue := makeQueue("my-queue")

		queueRepository.On("Get", ctx, queue.ID).Return(queue, nil)
		queueRepository.On("Update", ctx, queue).Return(nil)

		err := queueService.Update(ctx, queue)
		assert.Nil(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		queueRepository := mocks.NewQueueRepository(t)
		queueService := NewQueue(queueRepository)
		queue := makeQueue("my-queue")

		queueRepository.On("Get", ctx, queue.ID).Return(queue, nil)

		_, err := queueService.Get(ctx, queue.ID)
		assert.Nil(t, err)
	})

	t.Run("List", func(t *testing.T) {
		queueRepository := mocks.NewQueueRepository(t)
		queueService := NewQueue(queueRepository)
		queue1 := makeQueue("my-queue-1")
		queue2 := makeQueue("my-queue-2")

		queueRepository.On("List", ctx, uint(0), uint(10)).Return([]*domain.Queue{queue1, queue2}, nil)

		queues, err := queueService.List(ctx, 0, 10)
		assert.Nil(t, err)
		assert.Len(t, queues, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		queueRepository := mocks.NewQueueRepository(t)
		queueService := NewQueue(queueRepository)
		queue := makeQueue("my-queue")

		queueRepository.On("Get", ctx, queue.ID).Return(queue, nil)
		queueRepository.On("Delete", ctx, queue.ID).Return(nil)

		err := queueService.Delete(ctx, queue.ID)
		assert.Nil(t, err)
	})

	t.Run("Stats", func(t *testing.T) {
		queueRepository := mocks.NewQueueRepository(t)
		queueService := NewQueue(queueRepository)
		queue := makeQueue("my-queue")

		queueRepository.On("Get", ctx, queue.ID).Return(queue, nil)
		queueRepository.On("Stats", ctx, queue.ID).Return(&domain.QueueStats{}, nil)

		_, err := queueService.Stats(ctx, queue.ID)
		assert.Nil(t, err)
	})

	t.Run("Purge", func(t *testing.T) {
		queueRepository := mocks.NewQueueRepository(t)
		queueService := NewQueue(queueRepository)
		queue := makeQueue("my-queue")

		queueRepository.On("Get", ctx, queue.ID).Return(queue, nil)
		queueRepository.On("Purge", ctx, queue.ID).Return(nil)

		err := queueService.Purge(ctx, queue.ID)
		assert.Nil(t, err)
	})
}
