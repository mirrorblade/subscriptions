package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mirrorblade/subscriptions/internal/domain"
	"github.com/mirrorblade/subscriptions/internal/repository"
)

type SubscriptionsService struct {
	subscriptions repository.Subscriptions
}

func NewSubscriptionsService(subscriptions Subscriptions) *SubscriptionsService {
	return &SubscriptionsService{
		subscriptions: subscriptions,
	}
}

func (s *SubscriptionsService) GetByID(context context.Context, id uuid.UUID) (domain.Subscription, error) {
	return s.subscriptions.GetByID(context, id)
}

func (s *SubscriptionsService) GetListByUserID(context context.Context, userID uuid.UUID) ([]domain.Subscription, error) {
	return s.subscriptions.GetListByUserID(context, userID)
}

func (s *SubscriptionsService) GetPriceSumByUserID(context context.Context, userID uuid.UUID, parameters repository.GetSumParameters) (int64, error) {
	return s.subscriptions.GetPriceSumByUserID(context, userID, parameters)
}

func (s *SubscriptionsService) Create(context context.Context, subscription domain.Subscription) error {
	return s.subscriptions.Create(context, subscription)
}

func (s *SubscriptionsService) UpdateByID(context context.Context, id uuid.UUID, parameters repository.UpdateParameters) error {
	return s.subscriptions.UpdateByID(context, id, parameters)
}

func (s *SubscriptionsService) DeleteByID(context context.Context, id uuid.UUID) error {
	return s.subscriptions.DeleteByID(context, id)
}
