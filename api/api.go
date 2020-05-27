package api

import (
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

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
		// middleware.Recoverer,
		// middleware.NewCompressor(5).Handler,
		// middleware.RequestID,
		// middleware.RealIP,
		cors.Handler,
	)

	version := newVersionAPI(&versionSvc, logger.With().Str("api", "auth").Logger())
	version.Routes(r)

	return r
}
