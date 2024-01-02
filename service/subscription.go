package service

import (
	"context"
	"time"

	"github.com/allisson/psqlqueue/domain"
)

// Subscription is an implementation of domain.SubscriptionService.
type Subscription struct {
	subscriptionRepository domain.SubscriptionRepository
}

func (s *Subscription) Create(ctx context.Context, subscription *domain.Subscription) error {
	if err := subscription.Validate(); err != nil {
		return err
	}

	subscription.CreatedAt = time.Now().UTC()

	return s.subscriptionRepository.Create(ctx, subscription)
}

func (s *Subscription) Get(ctx context.Context, id string) (*domain.Subscription, error) {
	return s.subscriptionRepository.Get(ctx, id)
}

func (s *Subscription) List(ctx context.Context, offset, limit uint) ([]*domain.Subscription, error) {
	return s.subscriptionRepository.List(ctx, offset, limit)
}

func (s *Subscription) Delete(ctx context.Context, id string) error {
	subscription, err := s.subscriptionRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	return s.subscriptionRepository.Delete(ctx, subscription.ID)
}

// NewSubscription returns an implementation of domain.SubscriptionService.
func NewSubscription(subscriptionRepository domain.SubscriptionRepository) *Subscription {
	return &Subscription{subscriptionRepository: subscriptionRepository}
}
