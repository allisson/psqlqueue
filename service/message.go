package service

import (
	"context"
	"time"

	"github.com/allisson/psqlqueue/domain"
)

// Message is an implementation of domain.MessageService
type Message struct {
	messageRepository domain.MessageRepository
	queueRepository   domain.QueueRepository
}

func (m *Message) Create(ctx context.Context, message *domain.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	queue, err := m.queueRepository.Get(ctx, message.QueueID)
	if err != nil {
		return err
	}

	message.Enqueue(queue, time.Now().UTC())

	return m.messageRepository.Create(ctx, message)
}

func (m *Message) List(ctx context.Context, queueID string, label *string, limit uint) ([]*domain.Message, error) {
	queue, err := m.queueRepository.Get(ctx, queueID)
	if err != nil {
		return nil, err
	}

	return m.messageRepository.List(ctx, queue, label, limit)
}

func (m *Message) Ack(ctx context.Context, id string) error {
	return m.messageRepository.Ack(ctx, id)
}

func (m *Message) Nack(ctx context.Context, id string, visibilityTimeoutSeconds uint) error {
	return m.messageRepository.Nack(ctx, id, visibilityTimeoutSeconds)
}

// NewMessage returns an implementation of domain.MessageService.
func NewMessage(messageRepository domain.MessageRepository, queueRepository domain.QueueRepository) *Message {
	return &Message{
		messageRepository: messageRepository,
		queueRepository:   queueRepository,
	}
}
