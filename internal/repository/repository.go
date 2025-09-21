package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mirrorblade/subscriptions/internal/domain"
)

type GetSumParameters struct {
	ServiceName *string
	FromDate    *time.Time
	ToDate      *time.Time
}

type UpdateParameters struct {
	Price   *int64
	EndDate *time.Time
}

type Subscriptions interface {
	GetByID(context context.Context, id uuid.UUID) (domain.Subscription, error)
	GetListByUserID(context context.Context, userID uuid.UUID) ([]domain.Subscription, error)
	GetPriceSumByUserID(context context.Context, userID uuid.UUID, parameters GetSumParameters) (int64, error)
	Create(context context.Context, subscription domain.Subscription) error
	UpdateByID(context context.Context, id uuid.UUID, parameters UpdateParameters) error
	DeleteByID(context context.Context, id uuid.UUID) error
}

type Respository struct {
	Subscriptions Subscriptions
}

func New(subscriptions Subscriptions) *Respository {
	return &Respository{
		Subscriptions: subscriptions,
	}
}
