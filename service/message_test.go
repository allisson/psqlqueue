package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/psqlqueue/domain"
	"github.com/allisson/psqlqueue/mocks"
)

func nilString() *string {
	var s *string = nil
	return s
}

func TestMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		messageRepository := mocks.NewMessageRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageService := NewMessage(messageRepository, queueRepository)
		queue := makeQueue("my-queue", nil)
		message := domain.Message{Body: `{"data": true}`, QueueID: queue.ID}

		queueRepository.On("Get", ctx, queue.ID).Return(queue, nil)
		messageRepository.On("Create", ctx, mock.Anything).Return(nil)

		err := messageService.Create(ctx, &message)
		assert.Nil(t, err)
	})

	t.Run("List", func(t *testing.T) {
		messageRepository := mocks.NewMessageRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageService := NewMessage(messageRepository, queueRepository)
		queue := makeQueue("my-queue", nil)
		message1 := domain.Message{Body: `{"data": true}`}
		message1.Enqueue(queue, time.Now().UTC())
		message2 := domain.Message{Body: `{"data": true}`}
		message2.Enqueue(queue, time.Now().UTC())

		queueRepository.On("Get", ctx, queue.ID).Return(queue, nil)
		messageRepository.On("List", ctx, queue, nilString(), uint(10)).Return([]*domain.Message{&message1, &message2}, nil)

		messages, err := messageService.List(ctx, queue.ID, nilString(), 10)
		assert.Nil(t, err)
		assert.Len(t, messages, 2)
	})

	t.Run("Ack", func(t *testing.T) {
		messageRepository := mocks.NewMessageRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageService := NewMessage(messageRepository, queueRepository)
		queue := makeQueue("my-queue", nil)
		message := domain.Message{Body: `{"data": true}`}
		message.Enqueue(queue, time.Now().UTC())

		messageRepository.On("Ack", ctx, message.ID).Return(nil)

		err := messageService.Ack(ctx, message.ID)
		assert.Nil(t, err)
	})

	t.Run("Nack", func(t *testing.T) {
		messageRepository := mocks.NewMessageRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageService := NewMessage(messageRepository, queueRepository)
		queue := makeQueue("my-queue", nil)
		message := domain.Message{Body: `{"data": true}`}
		message.Enqueue(queue, time.Now().UTC())

		messageRepository.On("Nack", ctx, message.ID, uint(0)).Return(nil)

		err := messageService.Nack(ctx, message.ID, uint(0))
		assert.Nil(t, err)
	})
}
