package api

import (
	"net/http"
	"net/http/httputil"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/httptrace"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

func traceMW(log zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			path := "api.endpoint:" + r.URL.EscapedPath() + ":" + r.Method
			tracer := global.Tracer("service")
			attrs, _, _ := httptrace.Extract(r.Context(), r)
			ctx, span := tracer.Start(
				r.Context(),
				path,
				trace.WithAttributes(attrs...),
			)
			defer span.End()

			meter := global.Meter("service")
			counter := metric.Must(meter).NewInt64Counter("api_hit").Bind(kv.String("endpoint", path))
			defer counter.Unbind()
			counter.Add(ctx, 1)

			dump, _ := httputil.DumpRequestOut(r, true)
			recorder := metric.Must(meter).NewInt64ValueRecorder("bytes_recieved").Bind(kv.String("endpoint", path))
			defer recorder.Unbind()
			recorder.Record(ctx, int64(len(dump)))

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
		traceMW(logger.With().Str("api", "tracemw").Logger()),
		cors.Handler,
	)

	version := newVersionAPI(&versionSvc, logger.With().Str("api", "version").Logger())
	version.Routes(r)

	return r
}
