package main

// [START opentelemetry_trace_import]
import (
	"context"
	"log"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() {
	projectID := "digital-waters-276111"
	exporter, err := texporter.NewExporter(texporter.WithProjectID(projectID))
	if err != nil {
		log.Fatalf("texporter.NewExporter: %v", err)
	}
	tp, err := sdktrace.NewProvider(sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)

	// Create custom span.
	tracer := global.TraceProvider().Tracer("example.com/trace")
	tracer.WithSpan(context.Background(), "foo",
		func(_ context.Context) error {
			// Do some work.
			return nil
		})
}
