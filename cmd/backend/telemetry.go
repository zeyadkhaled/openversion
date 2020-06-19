package main

import (
	"log"
	"os"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initProviders() *push.Controller {
	collectorAddr, ok := os.LookupEnv("OTEL_RECIEVER_ENDPOINT")
	if !ok {
		collectorAddr = otlp.DefaultCollectorHost + ":" + string(otlp.DefaultCollectorHost)
	}
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

	// This sets OTLP as the global metric exporter
	// Comment out if another metric exporter is going to be used i.e prometheus
	// pusher := push.New(
	// 	simple.NewWithExactDistribution(),
	// 	exporter,
	// 	push.WithStateful(true),
	// 	push.WithPeriod(2*time.Second),
	// )
	// global.SetMeterProvider(pusher.Provider())
	// pusher.Start()

	return nil
}
