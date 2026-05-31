package server

import (
	"net/http"

	_ "github.com/Asilbeek1/Subscription-Service/docs"
	"github.com/Asilbeek1/Subscription-Service/internal/service"
	handler "github.com/Asilbeek1/Subscription-Service/internal/transport/http/handler/Subscriptions"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(service *service.SubscriptionService) *chi.Mux {
	handler := handler.NewHandler(service)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", handler.CreateSubscriptionHandler)
		r.Get("/", handler.ListHandler)
		r.Get("/{id}", handler.ReadHandler)
		r.Delete("/{id}", handler.DeleteHandler)
		r.Get("/total", handler.CalculateTotalHandler)
	})
	return r
}
