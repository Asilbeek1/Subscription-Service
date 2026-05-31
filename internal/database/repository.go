package database

import (
	"context"
	"errors"
	"time"

	"github.com/Asilbeek1/Subscription-Service/pkg/ctxutil"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("Not Found")
)

type Subscription struct {
	ServiceName string
	Price       int32
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

func (db *DB) Create(ctx context.Context, sub Subscription) error {
	requestID := ctxutil.GetRequestID(ctx)
	query := `
		INSERT INTO subscriptions
		(service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		`

	_, err := db.pool.Exec(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)
	if err != nil {
		db.log.Error("failed to create subscription",
			"request_id", requestID,
			"error", err,
		)
		return err
	}
	return nil
}

func (db *DB) Update(ctx context.Context, id int, sub Subscription) error {
	requestID := ctxutil.GetRequestID(ctx)
	query := `
		UPDATE subscriptions
		SET service_name = $1,
    		price = $2,
    		user_id = $3,
    		start_date = $4,
    		end_date = $5
		WHERE id = $6;
		`

	result, err := db.pool.Exec(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		id,
	)

	if err != nil {
		db.log.Error("Error updating the record",
			"request_id", requestID,
			"error", err,
		)
		return err
	}
	if result.RowsAffected() == 0 {
		db.log.Info("subscription not found")
		return ErrNotFound
	}
	return nil
}
