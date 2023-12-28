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
)
