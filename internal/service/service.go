package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mirrorblade/subscriptions/internal/domain"
	"github.com/mirrorblade/subscriptions/internal/repository"
)

type Subscriptions interface {
	GetByID(context context.Context, id uuid.UUID) (domain.Subscription, error)
	GetListByUserID(context context.Context, userID uuid.UUID) ([]domain.Subscription, error)
	GetPriceSumByUserID(context context.Context, userID uuid.UUID, parameters repository.GetSumParameters) (int64, error)
	Create(context context.Context, subscription domain.Subscription) error
	UpdateByID(context context.Context, id uuid.UUID, parameters repository.UpdateParameters) error
	DeleteByID(context context.Context, id uuid.UUID) error
}

type Service struct {
	subscriptions Subscriptions
}

func New(subscriptions Subscriptions) *Service {
	return &Service{
		subscriptions: subscriptions,
	}
}
