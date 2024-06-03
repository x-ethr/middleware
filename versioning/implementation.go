package versioning

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/x-ethr/middleware/types"
)

type generic struct {
	types.Valuer[string]

	options *Settings
}

func (*generic) Value(ctx context.Context) Version {
	if v, ok := ctx.Value(key).(Version); ok {
		return v
	}

	return Version{Service: "development"}
}

func (g *generic) Configuration(options ...Variadic) Implementation {
	var o = settings()
	for _, option := range options {
		option(o)
	}

	g.options = o

	return g
}

func (g *generic) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		{
			value := g.options.Version

			service := value.Service

			slog.Log(ctx, (slog.LevelDebug - 4), "Middleware", slog.Group("context", slog.String("key", string(key)), slog.Any("value", map[string]string{"service": service})))

			ctx = context.WithValue(ctx, key, value)

			w.Header().Set("X-Service-Version", service)
			if v := r.Header.Get(http.CanonicalHeaderKey("X-API-Version")); v != "" {
				w.Header().Set("X-API-Version", v)
			}
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
