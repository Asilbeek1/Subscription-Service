package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrWrongDate = errors.New("Error wrong date")
	ErrNotFound  = errors.New("Not Found")
)

type Repository interface {
	Create(ctx context.Context, sub CreateSubscription) (*Subscription, error)
	GetById(ctx context.Context, id int) (*Subscription, error)
	Update(ctx context.Context, id int, sub UpdateSubscription) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]Subscription, error)
	ListByFilter(ctx context.Context, filter FilteredSum) ([]Subscription, error)
}

type Subscription struct {
	ID          int        `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int32      `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type CreateSubscription struct {
	ServiceName string     `json:"service_name"`
	Price       int32      `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type UpdateSubscription struct {
	ServiceName string     `json:"service_name"`
	Price       int32      `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type FilteredSum struct {
	UserID      *uuid.UUID `json:"user_id"`
	ServiceName *string    `json:"service_name"`
	From        *time.Time `json:"from"`
	To          *time.Time `json:"to"`
}
