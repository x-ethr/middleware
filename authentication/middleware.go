package authentication

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Authentication struct {
	Token *jwt.Token
	Email string `json:"-"` // Email represents the user's email address as set by the "sub" jwt-claims structure.
	Raw   string // Raw represents the raw jwt token as submitted by the client.
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
