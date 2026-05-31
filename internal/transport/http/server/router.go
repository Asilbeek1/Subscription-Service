package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New() *http.Server {
	r := chi.NewRouter()
}
