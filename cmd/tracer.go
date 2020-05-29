package main

import (
	"log"
	"os"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() {
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
}
