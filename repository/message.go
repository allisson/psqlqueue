package repository

import (
	"context"
	"time"

	"github.com/allisson/pgxutil/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/allisson/psqlqueue/domain"
)

// Message is an implementation of domain.MessageRepository.
type Message struct {
	pool      *pgxpool.Pool
	tableName string
}

func (m *Message) Create(ctx context.Context, message *domain.Message) error {
	return parseError(pgxutil.Insert(ctx, m.pool, "", m.tableName, message), domain.ErrMessageNotFound, domain.ErrMessageAlreadyExists)
}

func (m *Message) Get(ctx context.Context, id string) (*domain.Message, error) {
	message := domain.Message{}
	options := pgxutil.NewFindOptions().WithFilter("id", id)
	err := pgxutil.Get(ctx, m.pool, m.tableName, options, &message)
	return &message, err
}

func (m *Message) List(ctx context.Context, queue *domain.Queue, label *string, limit uint) ([]*domain.Message, error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	messages := []*domain.Message{}
	now := time.Now().UTC()
	options := pgxutil.NewFindAllOptions().
		WithFilter("queue_id", queue.ID).
		WithFilter("expired_at.gte", now).
		WithFilter("scheduled_at.lte", now).
		WithLimit(int(limit)).
		WithForUpdate("SKIP LOCKED").
		WithOrderBy("scheduled_at asc")
	if label != nil {
		options = options.WithFilter("label", label)
	}
	if err := parseError(pgxutil.Select(ctx, tx, m.tableName, options, &messages), domain.ErrMessageNotFound, domain.ErrMessageAlreadyExists); err != nil {
		executeRollback(ctx, tx)
		return nil, err
	}

	for i := range messages {
		message := messages[i]

		message.DeliverySetup(queue, now)
		if err := pgxutil.Update(ctx, tx, "", m.tableName, message.ID, &message); err != nil {
			executeRollback(ctx, tx)
			return nil, err
		}
	}

	return messages, tx.Commit(ctx)
}

func (m *Message) Ack(ctx context.Context, id string) error {
	message, err := m.Get(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	message.Ack(now)

	return pgxutil.Update(ctx, m.pool, "", m.tableName, message.ID, &message)
}

func (m *Message) Nack(ctx context.Context, id string, visibilityTimeoutSeconds uint) error {
	message, err := m.Get(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	message.Nack(now, visibilityTimeoutSeconds)

	return pgxutil.Update(ctx, m.pool, "", m.tableName, message.ID, &message)
}

// NewMessage returns an implementation of domain.MessageRepository.
func NewMessage(pool *pgxpool.Pool) *Message {
	return &Message{pool: pool, tableName: "messages"}
}
