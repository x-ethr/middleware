package cors

import (
	"context"
	"net/http"
)

type Implementation interface {
	Value(ctx context.Context) bool
	Middleware(next http.Handler) http.Handler
}

func New() Implementation {
	return &generic{}
}
