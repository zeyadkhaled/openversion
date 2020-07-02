package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"opentelemetry.version.service/api"
	"opentelemetry.version.service/version"
	"opentelemetry.version.service/version/store/versionpostgre"
	versionredisstore "opentelemetry.version.service/version/store/versionredis"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	postgreConnStr string = "postgres://postgres:roottoor@db:5432/backend"
	redisConnStr   string = "redis:6379"
	serverAddr     string = ":8088"
)

func initService(ctx context.Context, logger zerolog.Logger) version.Service {
	tracer := global.Tracer("service")
	meter := global.Meter("service")
	versionStore, err := versionpostgre.New(
		ctx,
		postgreConnStr,
		logger.With().Str("package", "versionpostgre").Logger(), tracer,
	)
	if err != nil {
		log.Fatal("Failed to create postgre store", err)
	}

	versionCacheStore, err := versionredisstore.New(redisConnStr, "",
		0, "versionredis", time.Duration(time.Minute*30),
		logger.With().Str("package", "version").Logger(), versionStore, tracer)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create version cache store")
		log.Fatal(err)
	}

	instruments := version.Instruments{
		ErrCounter:      metric.Must(meter).NewInt64Counter("errors.counter"),
		ProcessDuration: metric.Must(meter).NewFloat64ValueRecorder("process.duration"),
	}

	metrics := version.Metric{
		Meter:       meter,
		Instruments: instruments,
	}

	return *version.New(versionCacheStore, tracer, metrics)
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	logger := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.DebugLevel)
	versionSvc := initService(ctx, logger)
	defer initProviders().Stop()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer cancel()
		server := &http.Server{
			Addr: serverAddr,
			Handler: api.Handler(
				versionSvc,
				logger.With().Str("api", "root").Logger(),
			),
		}
		logger.Info().Msg("rest server started")
		err := server.ListenAndServe()
		logger.Err(err).Msg("rest server end")
		return err
	})

	g.Wait()
}
