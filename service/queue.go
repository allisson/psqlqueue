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
	queue, err := q.queueRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	return q.queueRepository.Delete(ctx, queue.ID)

}

func (q *Queue) Stats(ctx context.Context, id string) (*domain.QueueStats, error) {
	queue, err := q.queueRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return q.queueRepository.Stats(ctx, queue.ID)

}

func (q *Queue) Purge(ctx context.Context, id string) error {
	queue, err := q.queueRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	return q.queueRepository.Purge(ctx, queue.ID)
}

func (q *Queue) Cleanup(ctx context.Context, id string) error {
	queue, err := q.queueRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	return q.queueRepository.Cleanup(ctx, queue.ID)
}

// NewQueue returns an implementation of domain.QueueService.
func NewQueue(queueRepository domain.QueueRepository) *Queue {
	return &Queue{queueRepository: queueRepository}
}
