package telemetrics

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/x-ethr/middleware/types"
)

type generic struct {
	types.Valuer[*Telemetry]

	options *Settings
}

func (g *generic) Configuration(options ...Variadic) Implementation {
	var o = settings()
	for _, option := range options {
		option(o)
	}

	g.options = o

	return g
}

func (*generic) Value(ctx context.Context) Telemetry {
	return ctx.Value(key).(Telemetry)
}

func (g *generic) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		value := Telemetry{
			Headers: make(map[string]string),
		}

		{ // --> headers
			for header, values := range r.Header {
				if len(values) > 0 {
					for index := range g.options.Headers {
						v := http.CanonicalHeaderKey(g.options.Headers[index])
						if http.CanonicalHeaderKey(header) == v {
							slog.Log(ctx, g.options.Level.Level(), "Found Telemetry Header in Request", slog.String("header", header), slog.String("value", values[0]))

							value.Headers[v] = values[0]
						}
					}
				}
			}

			ctx = context.WithValue(ctx, key, value)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
