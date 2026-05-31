package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Asilbeek1/Subscription-Service/internal/domain"
	"github.com/Asilbeek1/Subscription-Service/internal/service"
	"github.com/Asilbeek1/Subscription-Service/internal/transport/http/dto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *service.SubscriptionService
}

func NewHandler(service *service.SubscriptionService) *Handler {
	return &Handler{
		service: service,
	}
}

// @Summary      Create subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription body dto.CreateSubscriptionRequest true "Subscription"
// @Success      201 {object} domain.Subscription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [post]
func (h *Handler) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		writeError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		writeError(w, "invalid start_date, ", http.StatusBadRequest)
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		t, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			writeError(w, "invalid end_date, expected {M-Y}", http.StatusBadRequest)
			return
		}
		endDate = &t
	}

	sub := domain.CreateSubscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	res, err := h.service.CreateSubscription(r.Context(), sub)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, res)
}

// @Summary      Get subscription by ID
// @Tags         subscriptions
// @Produce      json
// @Param        id path int true "Subscription ID"
// @Success      200 {object} domain.Subscription
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /subscriptions/{id} [get]
func (h *Handler) ReadHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "id")
	id, err := strconv.Atoi(uid)
	if err != nil {
		writeError(w, "invalid id", http.StatusBadRequest)
		return
	}

	res, err := h.service.GetById(r.Context(), id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// @Summary      List all subscriptions
// @Tags         subscriptions
// @Produce      json
// @Success      200 {array} domain.Subscription
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [get]
func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.List(r.Context())
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// @Summary      Delete subscription
// @Tags         subscriptions
// @Param        id path int true "Subscription ID"
// @Success      204
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [delete]
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Calculate total cost
// @Tags         subscriptions
// @Param        user_id     query string false "User UUID"
// @Param        service_name query string false "Service name"
// @Param        from        query string false "From month MM-YYYY"
// @Param        to          query string false "To month MM-YYYY"
// @Success      200 {object} map[string]int64
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/total [get]
func (h *Handler) CalculateTotalHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := domain.FilteredSum{
		ServiceName: ptrString(q.Get("service_name")),
	}

	if v := q.Get("user_id"); v != "" {
		uuidVal, err := uuid.Parse(v)
		if err != nil {
			writeError(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		filter.UserID = &uuidVal
	}

	if v := q.Get("from"); v != "" {
		t, err := time.Parse("01-2006", v)
		if err != nil {
			writeError(w, "invalid from date, expected MM-YYYY", http.StatusBadRequest)
			return
		}
		filter.From = &t
	}

	if v := q.Get("to"); v != "" {
		t, err := time.Parse("01-2006", v)
		if err != nil {
			writeError(w, "invalid to date, expected MM-YYYY", http.StatusBadRequest)
			return
		}
		filter.To = &t
	}

	res, err := h.service.CalculateTotal(r.Context(), filter)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWrongDate):
			writeError(w, err.Error(), http.StatusBadRequest)
		default:
			writeError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]int64{
		"total": res,
	})
}
