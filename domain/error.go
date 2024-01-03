package domain

import "errors"

var (
	// ErrQueueAlreadyExists is returned when the queue already exists.
	ErrQueueAlreadyExists = errors.New("queue already exists")
	// ErrQueueNotFound is returned when the queue is not found.
	ErrQueueNotFound = errors.New("queue not found")
	// ErrMessageAlreadyExists is returned when the message already exists.
	ErrMessageAlreadyExists = errors.New("message already exists")
	// ErrMessageNotFound is returned when the message is not found.
	ErrMessageNotFound = errors.New("message not found")
	// ErrTopicAlreadyExists is returned when the topic already exists.
	ErrTopicAlreadyExists = errors.New("topic already exists")
	// ErrTopicNotFound is returned when the topic is not found.
	ErrTopicNotFound = errors.New("topic not found")
	// ErrSubscriptionAlreadyExists is returned when the subscription already exists.
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	// ErrSubscriptionNotFound is returned when the subscription is not found.
	ErrSubscriptionNotFound = errors.New("subscription not found")
)
