package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Asilbeek1/Subscription-Service/internal/domain"
	"github.com/Asilbeek1/Subscription-Service/pkg/ctxutil"
)

type SubscriptionService struct {
	log  *slog.Logger
	repo domain.Repository
}

func NewSubscriptionService(repo domain.Repository, log *slog.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
		log:  log,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, input domain.CreateSubscription) (*domain.Subscription, error) {
	requestID := ctxutil.GetRequestID(ctx)

	result, err := s.repo.Create(ctx, input)
	if err != nil {
		s.log.Error("Error Creating Subscription",
			"request_id", requestID,
			"error", err)
		return nil, err
	}
	return result, nil
}

func (s *SubscriptionService) GetById(ctx context.Context, id int) (*domain.Subscription, error) {
	requestID := ctxutil.GetRequestID(ctx)

	result, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			s.log.Error("error record not found",
				"requestID", requestID,
				"error", err)
			s.log.Error("get subscription failed",
				"request_ID", requestID,
				"error", err)
			return nil, err
		}
		return nil, err
	}

	return result, nil
}

func (s *SubscriptionService) List(ctx context.Context) ([]domain.Subscription, error) {

	requestID := ctxutil.GetRequestID(ctx)

	result, err := s.repo.List(ctx)
	if err != nil {
		s.log.Error("list subscriptions failed",
			"request_id", requestID,
			"error", err,
		)
		return nil, err
	}

	return result, nil
}
func (s *SubscriptionService) Delete(ctx context.Context, id int) error {
	requestID := ctxutil.GetRequestID(ctx)

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			s.log.Info("subscription not found",
				"request_id", requestID,
				"id", id,
			)
			return err
		}

		s.log.Error("delete failed",
			"request_id", requestID,
			"error", err,
		)
		return err
	}

	return nil
}

func (s *SubscriptionService) CalculateTotal(ctx context.Context, filter domain.FilteredSum) (int64, error) {
	var list []domain.Subscription
	requestID := ctxutil.GetRequestID(ctx)
	list, err := s.repo.ListByFilter(ctx, filter)
	if err != nil {
		s.log.Error("Error Filtering From DB",
			"request_id", requestID,
			"error", err)
		return 0, err
	}

	total := 0
	for _, sub := range list {
		end := sub.EndDate
		if end == nil {
			continue
		}

		months, err := calcMonth(sub.StartDate, *end, *filter.From, *filter.To, s.log)
		if err != nil {
			return 0, err
		}
		total += int(sub.Price) * months
	}
	return int64(total), nil
}
func calcMonth(start_date, end_date, from, to time.Time, log *slog.Logger) (int, error) {
	start := maxTime(start_date, from)
	end := minTime(end_date, to)

	if start.After(end) {
		log.Error("Error start date is bigger than end date")
		return 0, domain.ErrWrongDate
	}
	startIdx := from.Year()*12 + int(from.Month())
	endIdx := to.Year()*12 + int(to.Month())
	return endIdx - startIdx, nil
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
