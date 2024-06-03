package tracing_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/x-ethr/telemetry"

	"github.com/x-ethr/middleware"
	"github.com/x-ethr/middleware/tracing"
)

func Test(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug - 4)

	ctx := context.Background()
	service := "test-service"
	version := "development"

	var tracer = otel.Tracer("service-example")

	t.Run("Tracing-Middleware", func(t *testing.T) {
		middlewares := middleware.Middleware()
		middlewares.Add(middleware.New().Tracer().Configuration(func(options *tracing.Settings) { options.Tracer = tracer }).Middleware)

		mux := http.NewServeMux()

		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tracer := middleware.New().Tracer().Value(ctx)
			ctx, span := tracer.Start(ctx, "example-handler", trace.WithAttributes(attribute.String("handler", "example-handler")))

			defer span.End()

			span.SetStatus(http.StatusOK, "N/A")

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status": "ok",
			})
		})

		// --> Telemetry Setup
		shutdown, e := telemetry.Setup(ctx, service, version, func(options *telemetry.Settings) {
			if version == "development" && os.Getenv("CI") == "" {
				options.Zipkin.Enabled = false

				options.Tracer.Local = true
				options.Metrics.Local = true
				options.Logs.Local = true
			}
		})

		if e != nil {
			panic(e)
		}

		handler := otelhttp.NewHandler(middlewares.Handler(mux), "server", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))

		server := httptest.NewServer(handler)

		defer server.Close()
		defer shutdown(ctx)

		t.Run("Successful-API-Request", func(t *testing.T) {
			client := server.Client()
			request, e := http.NewRequest(http.MethodGet, server.URL, nil)
			if e != nil {
				t.Fatal(e)
			}

			response, exception := client.Do(request)
			if exception != nil {
				t.Fatal(exception)
			}

			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				t.Fatalf("Unexpected Status Code: %d", response.StatusCode)
			}

			content, e := io.ReadAll(response.Body)
			if e != nil {
				t.Fatalf("Couldn't Read Response Body: %v", e)
			}

			t.Logf("Response: %s", string(content))
		})
	})
}
