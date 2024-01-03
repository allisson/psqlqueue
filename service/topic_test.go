package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/psqlqueue/domain"
	"github.com/allisson/psqlqueue/mocks"
)

func makeTopic(topicID string) *domain.Topic {
	return &domain.Topic{
		ID:        topicID,
		CreatedAt: time.Now().UTC(),
	}
}

func TestTopic(t *testing.T) {
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		topicRepository := mocks.NewTopicRepository(t)
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageRepository := mocks.NewMessageRepository(t)
		topicService := NewTopic(topicRepository, subscriptionRepository, queueRepository, messageRepository)
		topic := makeTopic("my-topic")

		topicRepository.On("Create", ctx, topic).Return(nil)

		err := topicService.Create(ctx, topic)
		assert.Nil(t, err)
	})

	t.Run("Create with invalid id", func(t *testing.T) {
		expectedErrorPayload := `{"id":"must be in a valid format"}`
		topicRepository := mocks.NewTopicRepository(t)
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageRepository := mocks.NewMessageRepository(t)
		topicService := NewTopic(topicRepository, subscriptionRepository, queueRepository, messageRepository)
		topic := makeTopic("my@topic")

		err := topicService.Create(ctx, topic)
		assert.NotNil(t, err)
		errorPayload, err := json.Marshal(err)
		assert.Nil(t, err)
		assert.Equal(t, expectedErrorPayload, string(errorPayload))
	})

	t.Run("Get", func(t *testing.T) {
		topicRepository := mocks.NewTopicRepository(t)
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageRepository := mocks.NewMessageRepository(t)
		topicService := NewTopic(topicRepository, subscriptionRepository, queueRepository, messageRepository)
		topic := makeTopic("my-topic")

		topicRepository.On("Get", ctx, topic.ID).Return(topic, nil)

		_, err := topicService.Get(ctx, topic.ID)
		assert.Nil(t, err)
	})

	t.Run("List", func(t *testing.T) {
		topicRepository := mocks.NewTopicRepository(t)
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageRepository := mocks.NewMessageRepository(t)
		topicService := NewTopic(topicRepository, subscriptionRepository, queueRepository, messageRepository)
		topic1 := makeTopic("my-topic-1")
		topic2 := makeTopic("my-topic-2")

		topicRepository.On("List", ctx, uint(0), uint(10)).Return([]*domain.Topic{topic1, topic2}, nil)

		topics, err := topicService.List(ctx, 0, 10)
		assert.Nil(t, err)
		assert.Len(t, topics, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		topicRepository := mocks.NewTopicRepository(t)
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageRepository := mocks.NewMessageRepository(t)
		topicService := NewTopic(topicRepository, subscriptionRepository, queueRepository, messageRepository)
		topic := makeTopic("my-topic")

		topicRepository.On("Get", ctx, topic.ID).Return(topic, nil)
		topicRepository.On("Delete", ctx, topic.ID).Return(nil)

		err := topicService.Delete(ctx, topic.ID)
		assert.Nil(t, err)
	})

	t.Run("CreateMessage", func(t *testing.T) {
		topicRepository := mocks.NewTopicRepository(t)
		subscriptionRepository := mocks.NewSubscriptionRepository(t)
		queueRepository := mocks.NewQueueRepository(t)
		messageRepository := mocks.NewMessageRepository(t)
		topicService := NewTopic(topicRepository, subscriptionRepository, queueRepository, messageRepository)
		topic := makeTopic("my-topic")
		queue := makeQueue("my-queue")
		subscription := makeSubscription("my-subscription", topic.ID, queue.ID)
		message := &domain.Message{Body: "my-message-body"}

		topicRepository.On("Get", ctx, topic.ID).Return(topic, nil)
		subscriptionRepository.On("ListByTopic", ctx, topic.ID, uint(0), uint(50)).Return([]*domain.Subscription{subscription}, nil)
		subscriptionRepository.On("ListByTopic", ctx, topic.ID, uint(50), uint(50)).Return([]*domain.Subscription{}, nil)
		queueRepository.On("Get", ctx, queue.ID).Return(queue, nil)
		messageRepository.On("CreateMany", ctx, mock.Anything).Return(nil)

		err := topicService.CreateMessage(ctx, topic.ID, message)
		assert.Nil(t, err)
	})
}
