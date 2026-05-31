package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/Asilbeek1/Subscription-Service/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (db *DB) Create(ctx context.Context, sub domain.CreateSubscription) (*domain.Subscription, error) {
	query := `
		INSERT INTO subscriptions
		(service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err := db.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &domain.Subscription{
		ID:          id,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}, nil
}

func (db *DB) GetById(ctx context.Context, id int) (*domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`
	var response domain.Subscription
	err := db.pool.QueryRow(ctx, query, id).Scan(
		&response.ID,
		&response.ServiceName,
		&response.Price,
		&response.UserID,
		&response.StartDate,
		&response.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &response, nil
}

func (db *DB) Update(ctx context.Context, id int, sub domain.UpdateSubscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1,
			price = $2,
			user_id = $3,
			start_date = $4,
			end_date = $5
		WHERE id = $6
	`

	result, err := db.pool.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	result, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (db *DB) List(ctx context.Context) ([]domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		ORDER BY id ASC
	`

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Subscription

	for rows.Next() {
		var sub domain.Subscription

		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return nil, err
		}

		result = append(result, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *DB) FilterSum(ctx context.Context, filter domain.FilteredSum) ([]domain.Subscription, error) {
	query, args := buildQueryFilter(filter)

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Subscription

	for rows.Next() {
		var sub domain.Subscription

		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return nil, err
		}

		result = append(result, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func buildQueryFilter(filter domain.FilteredSum) (string, []interface{}) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *filter.UserID)
		argPos++
	}

	if filter.ServiceName != nil && *filter.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argPos)
		args = append(args, *filter.ServiceName)
		argPos++
	}

	if filter.From != nil {
		query += fmt.Sprintf(" AND start_date >= $%d", argPos)
		args = append(args, *filter.From)
		argPos++
	}

	if filter.To != nil {
		query += fmt.Sprintf(" AND start_date <= $%d", argPos)
		args = append(args, *filter.To)
		argPos++
	}

	return query, args
}
