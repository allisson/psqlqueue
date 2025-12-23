package service

import (
	"context"
	"time"

	"github.com/allisson/psqlqueue/domain"
)

// Topic is an implementation of domain.TopicService.
type Topic struct {
	topicRepository        domain.TopicRepository
	subscriptionRepository domain.SubscriptionRepository
	queueRepository        domain.QueueRepository
	messageRepository      domain.MessageRepository
}

func (t *Topic) Create(ctx context.Context, topic *domain.Topic) error {
	if err := topic.Validate(); err != nil {
		return err
	}

	topic.CreatedAt = time.Now().UTC()

	return t.topicRepository.Create(ctx, topic)
}

func (t *Topic) Get(ctx context.Context, id string) (*domain.Topic, error) {
	return t.topicRepository.Get(ctx, id)
}

func (t *Topic) List(ctx context.Context, offset, limit uint) ([]*domain.Topic, error) {
	return t.topicRepository.List(ctx, offset, limit)
}

func (t *Topic) Delete(ctx context.Context, id string) error {
	return t.topicRepository.Delete(ctx, id)
}

func (t *Topic) CreateMessage(ctx context.Context, topicID string, message *domain.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	topic, err := t.topicRepository.Get(ctx, topicID)
	if err != nil {
		return err
	}

	messages := []*domain.Message{}
	offset := 0
	limit := 50
	now := time.Now().UTC()

	for {
		subscriptions, err := t.subscriptionRepository.ListByTopic(ctx, topic.ID, uint(offset), uint(limit))
		if err != nil {
			return err
		}

		if len(subscriptions) == 0 {
			break
		}

		// Collect unique queue IDs to fetch in batch
		queueIDs := make([]string, 0, len(subscriptions))
		seenQueues := make(map[string]bool)
		for i := range subscriptions {
			subscription := subscriptions[i]
			if !subscription.ShouldCreateMessage(message) {
				continue
			}
			if !seenQueues[subscription.QueueID] {
				queueIDs = append(queueIDs, subscription.QueueID)
				seenQueues[subscription.QueueID] = true
			}
		}

		// Fetch all queues in a single query
		queues, err := t.queueRepository.GetMany(ctx, queueIDs)
		if err != nil {
			return err
		}

		for i := range subscriptions {
			subscription := subscriptions[i]
			if !subscription.ShouldCreateMessage(message) {
				continue
			}

			queue, ok := queues[subscription.QueueID]
			if !ok {
				return domain.ErrQueueNotFound
			}

			newMessage := &domain.Message{
				Label:      message.Label,
				Body:       message.Body,
				Attributes: message.Attributes,
			}
			newMessage.Enqueue(queue, now)
			messages = append(messages, newMessage)
		}

		offset += limit
	}

	return t.messageRepository.CreateMany(ctx, messages)
}

// NewTopic returns an implementation of domain.TopicService.
func NewTopic(topicRepository domain.TopicRepository, subscriptionRepository domain.SubscriptionRepository, queueRepository domain.QueueRepository, messageRepository domain.MessageRepository) *Topic {
	return &Topic{
		topicRepository:        topicRepository,
		subscriptionRepository: subscriptionRepository,
		queueRepository:        queueRepository,
		messageRepository:      messageRepository,
	}
}
