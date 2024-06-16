package authentication

import (
	"context"
	"net/http"
)

type Authentication struct {
	Email string `json:"-"` // Email represents the user's email address as set by the "sub" jwt-claims structure.
}

type Implementation interface {
	Value(ctx context.Context) *Authentication
	Configuration(options ...Variadic) Implementation
	Middleware(next http.Handler) http.Handler
}

func New() Implementation {
	return &generic{
		options: settings(),
	}
}
