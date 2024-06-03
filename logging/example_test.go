package logging_test

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/x-ethr/middleware"
)

func Example() {
	middlewares := middleware.Middleware()
	middlewares.Add(middleware.New().Logging().Middleware)

	mux := http.NewServeMux()

	handler := middlewares.Handler(mux)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger := middleware.New().Logging().Value(ctx)

		var response = map[string]interface{}{
			"key": "value",
		}

		logger = logger.With(slog.String("key", "value"))

		logger.InfoContext(ctx, "Response", slog.Any("response", response))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":8080", handler)
}
