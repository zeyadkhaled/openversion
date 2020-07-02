package api

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/instrumentation/httptrace"
	"opentelemetry.version.service/version"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

func telemetryMW(log zerolog.Logger, tracer trace.Tracer, meter metric.Meter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := "api.endpoint:" + r.URL.EscapedPath() + ":" + r.Method
			// Skips tracing and metric collection for prometheus exposed endpoint
			if strings.Contains(path, "/metrics") {
				next.ServeHTTP(w, r)
				return
			}

			// Tracing start
			attrs, _, _ := httptrace.Extract(r.Context(), r)
			ctx, span := tracer.Start(
				r.Context(),
				path,
				trace.WithAttributes(attrs...),
			)
			defer span.End()

			// Metrics start
			labels := []kv.KeyValue{kv.String("endpoint", path)}
			counter := metric.Must(meter).NewInt64Counter("api.hit.count")
			recorder := metric.Must(meter).NewInt64ValueRecorder("bytes.recieved")
			meter.RecordBatch(ctx,
				labels,
				counter.Measurement(1),
				recorder.Measurement(r.ContentLength))

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func Handler(versionSvc version.Service, logger zerolog.Logger) *chi.Mux {

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(
		telemetryMW(logger.With().Str("api", "tracemw").Logger(), versionSvc.Tracer, versionSvc.Meterics.Meter),
		cors.Handler,
	)

	version := newVersionAPI(&versionSvc, logger.With().Str("api", "version").Logger())
	version.Routes(r)
	return r
}
