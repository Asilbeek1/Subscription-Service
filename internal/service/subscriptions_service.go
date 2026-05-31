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

func (s *SubscriptionService) CreateSubscription(ctx context.Context, sub domain.CreateSubscription) (*domain.Subscription, error) {
	requestID := ctxutil.GetRequestID(ctx)

	if sub.Price <= 0 {
		s.log.Warn("invalid price price must be bigger than 0",
			"request_id", requestID)
		return nil, domain.ErrWrongPrice
	}
	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		s.log.Warn("invalid date range: end before start", "request_id", requestID,
			"start_date", sub.StartDate, "end_date", sub.EndDate)
		return nil, domain.ErrWrongDate
	}

	result, err := s.repo.Create(ctx, sub)
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
			return nil, err
		}
		s.log.Error("get subscription failed",
			"request_ID", requestID,
			"error", err)
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
	requestID := ctxutil.GetRequestID(ctx)

	from := time.Time{}
	to := time.Now()
	if filter.From != nil {
		from = *filter.From
	}
	if filter.To != nil {
		to = *filter.To
	}

	list, err := s.repo.ListByFilter(ctx, filter)
	if err != nil {
		s.log.Error("Error to filter from db", "request_id", requestID, "error", err)
		return 0, err
	}

	var total int64
	for _, sub := range list {
		start := maxTime(sub.StartDate, from)

		var end time.Time
		if sub.EndDate == nil {
			end = to
		} else {
			end = minTime(*sub.EndDate, to)
		}

		if start.After(end) {
			continue
		}

		total += int64(sub.Price) * int64(calcMonth(start, end))
	}

	return total, nil
}
func calcMonth(start, end time.Time) int {
	startIdx := start.Year()*12 + int(start.Month())
	endIdx := end.Year()*12 + int(end.Month())
	return endIdx - startIdx + 1
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
