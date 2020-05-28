package api

import (
	"net/http"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/httptrace"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

func traceMW(log zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tracer := global.Tracer("service")
			attrs, _, _ := httptrace.Extract(r.Context(), r)
			ctx, span := tracer.Start(
				r.Context(),
				"api.endpoint:"+r.URL.EscapedPath(),
				trace.WithAttributes(attrs...),
			)
			defer span.End()

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

	version := newVersionAPI(&versionSvc, logger.With().Str("api", "auth").Logger())
	version.Routes(r)

	return r
}
