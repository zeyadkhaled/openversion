package main

// TODO add config and dev
// TODO add migrations file
import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/api"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/filterenc"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version/store/versioninmem"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version/store/versionpostgre"
	versionredisstore "gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version/store/versionredis"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
)

func initTracer() {
	collectorAddr, ok := os.LookupEnv("OTEL_RECIEVER_ENDPOINT")
	if !ok {
		collectorAddr = otlp.DefaultCollectorHost + ":" + string(otlp.DefaultCollectorHost)
	}
	// collectorAddr = "localhost:55678"
	exporter, err := otlp.NewExporter(otlp.WithAddress(collectorAddr), otlp.WithInsecure())

	if err != nil {
		log.Fatal(err)
	}

	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}

func main() {
	initTracer()
	logger := zerolog.New(zerolog.NewConsoleWriter()).Level(zerolog.DebugLevel)
	ctx := context.Background()

	var versionStore version.Store
	switch "postgres" {
	case "inmem":
		versionStore = versioninmem.New()
	case "postgres":
		versionStore, _ = versionpostgre.New(
			ctx,
			"postgres://postgres:roottoor@db:5432/backend",
			logger.With().Str("package", "versionpostgre").Logger(),
		)
	}

	filterKey, err := hex.DecodeString("ce4f34331feab353c0a6c5f27f98097c8e81c65b1f0dac259074d0063e27eddd")
	if err != nil {
		fmt.Printf("Failed to decode filter key: %v", err)
		os.Exit(1)
	}
	filterEncoder := filterenc.New(filterKey)

	versionCacheStore, err := versionredisstore.New("redis:6379", "",
		0, "versionredis", time.Duration(time.Minute*30),
		logger.With().Str("package", "version").Logger(), versionStore)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create version cache store")
		os.Exit(1)
	}

	versionSvc := version.New(versionCacheStore, filterEncoder)

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer cancel()
		server := &http.Server{
			Addr: ":8088",
			Handler: api.Handler(
				*versionSvc,
				logger.With().Str("api", "root").Logger(),
			),
		}
		logger.Info().Msg("rest server started")
		err := server.ListenAndServe()
		logger.Err(err).Msg("rest server end")
		return err
	})

	err = g.Wait()
	if err != nil {
		fmt.Println("server ended:", err)
	}
}
