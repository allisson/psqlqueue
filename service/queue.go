package service

import (
	"context"
	"time"

	"github.com/allisson/psqlqueue/domain"
)

// Queue is an implementation of domain.QueueService.
type Queue struct {
	queueRepository domain.QueueRepository
}

func (q *Queue) Create(ctx context.Context, queue *domain.Queue) error {
	if err := queue.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()
	queue.CreatedAt = now
	queue.UpdatedAt = now

	return q.queueRepository.Create(ctx, queue)
}

func (q *Queue) Update(ctx context.Context, queue *domain.Queue) error {
	if err := queue.Validate(); err != nil {
		return err
	}

	queueFromDB, err := q.queueRepository.Get(ctx, queue.ID)
	if err != nil {
		return err
	}

	queue.CreatedAt = queueFromDB.CreatedAt
	queue.UpdatedAt = time.Now().UTC()

	return q.queueRepository.Update(ctx, queue)
}

func (q *Queue) Get(ctx context.Context, id string) (*domain.Queue, error) {
	return q.queueRepository.Get(ctx, id)
}

func (q *Queue) List(ctx context.Context, offset, limit uint) ([]*domain.Queue, error) {
	return q.queueRepository.List(ctx, offset, limit)
}

func (q *Queue) Delete(ctx context.Context, id string) error {
	return q.queueRepository.Delete(ctx, id)
}

func (q *Queue) Stats(ctx context.Context, id string) (*domain.QueueStats, error) {
	return q.queueRepository.Stats(ctx, id)
}

func (q *Queue) Purge(ctx context.Context, id string) error {
	return q.queueRepository.Purge(ctx, id)
}

func (q *Queue) Cleanup(ctx context.Context, id string) error {
	return q.queueRepository.Cleanup(ctx, id)
}

// NewQueue returns an implementation of domain.QueueService.
func NewQueue(queueRepository domain.QueueRepository) *Queue {
	return &Queue{queueRepository: queueRepository}
}
