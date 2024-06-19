package authentication

import (
	"context"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"

	"github.com/x-ethr/middleware/types"
)

type Claims[Generic struct{}] struct {
	data *Generic // data represents the generic's custom claim data.

	jwt.RegisteredClaims
}

// Data represents a getter for the generic's custom claim data.
func (c Claims[Generic]) Data() *Generic {
	return c.data
}

// Set represents a setter for the generic's custom claim data.
func (c Claims[Generic]) Set(input *Generic) {
	c.data = input
}

type Claim[Generic struct{}] interface {
	Set(input *Generic)
	Data() *Generic
}

type Settings struct {
	Verification func(ctx context.Context, token string) (*jwt.Token, error) // Verification is a user-provided jwt-verification function.

	Level slog.Leveler // Level represents a [log/slog] log level - defaults to (slog.LevelDebug - 4) (trace)
}

type Variadic types.Variadic[Settings]

func settings() *Settings {
	return &Settings{
		Level: (slog.LevelDebug - 4),
	}
}
