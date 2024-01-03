package repository

import (
	"context"
	"time"

	"github.com/allisson/pgxutil/v2"
	"github.com/allisson/sqlquery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/allisson/psqlqueue/domain"
)

// Queue is an implementation of domain.QueueRepository.
type Queue struct {
	pool      *pgxpool.Pool
	tableName string
}

func (q *Queue) Create(ctx context.Context, queue *domain.Queue) error {
	return parseError(pgxutil.Insert(ctx, q.pool, "", q.tableName, queue), domain.ErrQueueNotFound, domain.ErrQueueAlreadyExists)
}

func (q *Queue) Update(ctx context.Context, queue *domain.Queue) error {
	return parseError(pgxutil.Update(ctx, q.pool, "", q.tableName, queue.ID, queue), domain.ErrQueueNotFound, domain.ErrQueueAlreadyExists)
}

func (q *Queue) Get(ctx context.Context, id string) (*domain.Queue, error) {
	queue := domain.Queue{}
	options := pgxutil.NewFindOptions().WithFilter("id", id)
	err := pgxutil.Get(ctx, q.pool, q.tableName, options, &queue)
	return &queue, parseError(err, domain.ErrQueueNotFound, domain.ErrQueueAlreadyExists)
}

func (q *Queue) List(ctx context.Context, offset, limit uint) ([]*domain.Queue, error) {
	queues := []*domain.Queue{}
	options := pgxutil.NewFindAllOptions().WithOffset(int(offset)).WithLimit(int(limit)).WithOrderBy("id asc")
	err := pgxutil.Select(ctx, q.pool, q.tableName, options, &queues)
	return queues, parseError(err, domain.ErrQueueNotFound, domain.ErrQueueAlreadyExists)
}

func (q *Queue) Delete(ctx context.Context, id string) error {
	return parseError(pgxutil.Delete(ctx, q.pool, q.tableName, id), domain.ErrQueueNotFound, domain.ErrQueueAlreadyExists)
}

func (q *Queue) Stats(ctx context.Context, id string) (*domain.QueueStats, error) {
	stats := &domain.QueueStats{}
	now := time.Now().UTC()
	options := pgxutil.NewFindAllOptions().
		WithFields([]string{"COUNT(1)"}).
		WithFilter("queue_id", id).
		WithFilter("expired_at.gte", now).
		WithFilter("scheduled_at.lte", now).
		WithLimit(1)
	sqlQuery, args := sqlquery.FindAllQuery("messages", options)
	err := q.pool.QueryRow(ctx, sqlQuery, args...).Scan(&stats.NumUndeliveredMessages)
	if err != nil {
		return stats, err
	}

	if stats.NumUndeliveredMessages == 0 {
		return stats, nil
	}

	var createdAt time.Time
	options = options.WithFields([]string{"created_at"}).WithOrderBy("created_at asc")
	sqlQuery, args = sqlquery.FindAllQuery("messages", options)
	err = q.pool.QueryRow(ctx, sqlQuery, args...).Scan(&createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return stats, nil
		}
		return stats, err
	}
	if !createdAt.IsZero() {
		stats.OldestUnackedMessageAgeSeconds = uint(now.Sub(createdAt).Seconds())
	}

	return stats, nil
}

func (q *Queue) Purge(ctx context.Context, id string) error {
	sqlQuery := `DELETE FROM messages WHERE queue_id = $1`
	_, err := q.pool.Exec(ctx, sqlQuery, id)
	return err
}

func (q *Queue) Cleanup(ctx context.Context, id string) error {
	now := time.Now().UTC()
	options := pgxutil.NewDeleteOptions().WithFilter("queue_id", id).WithFilter("expired_at.lte", now)
	return pgxutil.DeleteWithOptions(ctx, q.pool, "messages", options)
}

// NewQueue returns an implementation of domain.QueueRepository.
func NewQueue(pool *pgxpool.Pool) *Queue {
	return &Queue{pool: pool, tableName: "queues"}
}
