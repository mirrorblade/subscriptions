package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mirrorblade/subscriptions/internal/domain"
	"github.com/mirrorblade/subscriptions/internal/repository"
)

type Subscriptions struct {
	pool *pgxpool.Pool

	tableName string
}

func NewSubscriptions(pool *pgxpool.Pool, tableName string) *Subscriptions {
	return &Subscriptions{
		pool:      pool,
		tableName: tableName,
	}
}

func (s *Subscriptions) GetByID(context context.Context, id uuid.UUID) (domain.Subscription, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", s.tableName)

	rows, err := s.pool.Query(context, query, id)
	if err != nil {
		return domain.Subscription{}, err
	}

	subscription, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Subscription])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Subscription{}, domain.ErrSubscriptionNotFound
		}

		return domain.Subscription{}, err
	}

	return subscription, nil
}

func (s *Subscriptions) GetListByUserID(context context.Context, userID uuid.UUID) ([]domain.Subscription, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", s.tableName)

	rows, err := s.pool.Query(context, query, userID)
	if err != nil {
		return []domain.Subscription{}, err
	}

	subscriptions, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Subscription])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.Subscription{}, domain.ErrUserNotFound
		}

		return []domain.Subscription{}, err
	}

	return subscriptions, nil
}

func (s *Subscriptions) GetPriceSumByUserID(context context.Context, userID uuid.UUID, parameters repository.GetSumParameters) (int64, error) {
	query := fmt.Sprintf("SELECT SUM(price) FROM %s WHERE user_id = $1", s.tableName)

	args := []any{any(userID)}
	idx := 2

	if parameters.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", idx)
		args = append(args, *parameters.ServiceName)
		idx++
	}
	if parameters.FromDate != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", idx)
		args = append(args, *parameters.FromDate)
		idx++
	}
	if parameters.ToDate != nil {
		query += fmt.Sprintf(" AND end_date <= $%d", idx)
		args = append(args, *parameters.ToDate)
		idx++
	}

	rows, err := s.pool.Query(context, query, args...)
	if err != nil {
		return 0, err
	}

	sum, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[int64])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrUserNotFound
		}

		return 0, err
	}

	return sum, nil
}

func (s *Subscriptions) Create(context context.Context, subscription domain.Subscription) error {
	query := fmt.Sprintf("INSERT INTO %s (id, service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5, $6)", s.tableName)

	endDate := pgtype.Timestamp{}
	if subscription.EndDate == nil {
		endDate.Valid = false
	} else {
		endDate.Valid = true
		endDate.Time = *subscription.EndDate
	}

	if _, err := s.pool.Query(context, query, subscription.ID, subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, endDate); err != nil {
		return err
	}

	return nil
}

func (s *Subscriptions) UpdateByID(context context.Context, id uuid.UUID, parameters repository.UpdateParameters) error {
	query := fmt.Sprintf("UPDATE %s SET ", s.tableName)
	args := []any{}
	idx := 1

	if parameters.Price != nil {
		query += fmt.Sprintf("price = $%d, ", idx)
		args = append(args, *parameters.Price)
		idx++
	}
	if parameters.EndDate != nil {
		query += fmt.Sprintf("end_date = $%d, ", idx)
		args = append(args, *parameters.EndDate)
		idx++
	}

	if len(args) == 0 {
		return domain.ErrNoUpdateParameters
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", idx)
	args = append(args, id)

	commandTag, err := s.pool.Exec(context, query, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

func (s *Subscriptions) DeleteByID(context context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.tableName)

	commandTag, err := s.pool.Exec(context, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}
