package main

// TODO find why there is overhead time

import (
	"context"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/api"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/filterenc"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version/store/versionpostgre"
	versionredisstore "gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version/store/versionredis"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	postgreConnStr string = "postgres://postgres:roottoor@db:5432/backend"
	redisConnStr   string = "redis:6379"
	filterString   string = "ce4f34331feab353c0a6c5f27f98097c8e81c65b1f0dac259074d0063e27eddd"
	serverAddr     string = ":8088"
)

func initService(ctx context.Context, logger zerolog.Logger) version.Service {
	versionStore, err := versionpostgre.New(
		ctx,
		postgreConnStr,
		logger.With().Str("package", "versionpostgre").Logger(),
	)

	if err != nil {
		log.Fatal("Failed to create postgre store", err)
	}

	filterKey, err := hex.DecodeString(filterString)
	if err != nil {
		log.Fatal("Failed to decode filter key", err)
	}
	filterEncoder := filterenc.New(filterKey)

	versionCacheStore, err := versionredisstore.New(redisConnStr, "",
		0, "versionredis", time.Duration(time.Minute*30),
		logger.With().Str("package", "version").Logger(), versionStore)

	if err != nil {
		logger.Error().Err(err).Msg("failed to create version cache store")
		log.Fatal(err)
	}

	return *version.New(versionCacheStore, filterEncoder)
}

func main() {
	initTracer()

	ctx, cancel := context.WithCancel(context.Background())
	logger := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.DebugLevel)
	versionSvc := initService(ctx, logger)

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
