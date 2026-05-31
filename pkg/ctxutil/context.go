package ctxutil

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
)

type key string

const requestIDKey key = "request_id"

func WithRequestId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}
