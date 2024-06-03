package versioning_test

import (
	"encoding/json"
	"net/http"

	"github.com/x-ethr/middleware"
	"github.com/x-ethr/middleware/versioning"
)

func Example() {
	middlewares := middleware.Middleware()

	middlewares.Add(middleware.New().Version().Configuration(func(options *versioning.Settings) { options.Version = versioning.Version{Service: "0.0.0"} }).Middleware)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		version := middleware.New().Version().Value(ctx)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(version)
	})

	http.ListenAndServe(":8080", middlewares.Handler(mux))
}
