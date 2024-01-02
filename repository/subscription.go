package repository

import (
	"context"

	"github.com/allisson/pgxutil/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/allisson/psqlqueue/domain"
)

// Subscription is an implementation of domain.SubscriptionRepository.
type Subscription struct {
	pool      *pgxpool.Pool
	tableName string
}

func (s *Subscription) Create(ctx context.Context, subscription *domain.Subscription) error {
	return parseError(pgxutil.Insert(ctx, s.pool, "", s.tableName, subscription), domain.ErrSubscriptionNotFound, domain.ErrSubscriptionAlreadyExists)
}

func (s *Subscription) Get(ctx context.Context, id string) (*domain.Subscription, error) {
	subscription := domain.Subscription{}
	options := pgxutil.NewFindOptions().WithFilter("id", id)
	err := pgxutil.Get(ctx, s.pool, s.tableName, options, &subscription)
	return &subscription, parseError(err, domain.ErrSubscriptionNotFound, domain.ErrSubscriptionAlreadyExists)
}

func (s *Subscription) List(ctx context.Context, offset, limit uint) ([]*domain.Subscription, error) {
	subscriptions := []*domain.Subscription{}
	options := pgxutil.NewFindAllOptions().WithOffset(int(offset)).WithLimit(int(limit)).WithOrderBy("id asc")
	err := pgxutil.Select(ctx, s.pool, s.tableName, options, &subscriptions)
	return subscriptions, parseError(err, domain.ErrSubscriptionNotFound, domain.ErrSubscriptionAlreadyExists)
}

func (s *Subscription) ListByTopic(ctx context.Context, topicID string, offset, limit uint) ([]*domain.Subscription, error) {
	subscriptions := []*domain.Subscription{}
	options := pgxutil.NewFindAllOptions().WithFilter("topic_id", topicID).WithOffset(int(offset)).WithLimit(int(limit)).WithOrderBy("id asc")
	err := pgxutil.Select(ctx, s.pool, s.tableName, options, &subscriptions)
	return subscriptions, parseError(err, domain.ErrSubscriptionNotFound, domain.ErrSubscriptionAlreadyExists)
}

func (s *Subscription) Delete(ctx context.Context, id string) error {
	return parseError(pgxutil.Delete(ctx, s.pool, s.tableName, id), domain.ErrSubscriptionNotFound, domain.ErrSubscriptionAlreadyExists)
}

// NewSubscription returns an implementation of domain.SubscriptionRepository.
func NewSubscription(pool *pgxpool.Pool) *Subscription {
	return &Subscription{pool: pool, tableName: "subscriptions"}
}
