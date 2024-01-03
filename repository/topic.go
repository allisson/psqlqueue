package repository

import (
	"context"

	"github.com/allisson/pgxutil/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/allisson/psqlqueue/domain"
)

// Topic is an implementation of domain.TopicRepository.
type Topic struct {
	pool      *pgxpool.Pool
	tableName string
}

func (t *Topic) Create(ctx context.Context, topic *domain.Topic) error {
	return parseError(pgxutil.Insert(ctx, t.pool, "", t.tableName, topic), domain.ErrTopicNotFound, domain.ErrTopicAlreadyExists)
}

func (t *Topic) Get(ctx context.Context, id string) (*domain.Topic, error) {
	topic := domain.Topic{}
	options := pgxutil.NewFindOptions().WithFilter("id", id)
	err := pgxutil.Get(ctx, t.pool, t.tableName, options, &topic)
	return &topic, parseError(err, domain.ErrTopicNotFound, domain.ErrTopicAlreadyExists)
}

func (t *Topic) List(ctx context.Context, offset, limit uint) ([]*domain.Topic, error) {
	topics := []*domain.Topic{}
	options := pgxutil.NewFindAllOptions().WithOffset(int(offset)).WithLimit(int(limit)).WithOrderBy("id asc")
	err := pgxutil.Select(ctx, t.pool, t.tableName, options, &topics)
	return topics, parseError(err, domain.ErrTopicNotFound, domain.ErrTopicAlreadyExists)
}

func (t *Topic) Delete(ctx context.Context, id string) error {
	return parseError(pgxutil.Delete(ctx, t.pool, t.tableName, id), domain.ErrTopicNotFound, domain.ErrTopicAlreadyExists)
}

// NewTopic returns an implementation of domain.TopicRepository.
func NewTopic(pool *pgxpool.Pool) *Topic {
	return &Topic{pool: pool, tableName: "topics"}
}
