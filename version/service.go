package version

import (
	"context"
	"fmt"
	"time"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
)

type Instruments struct {
	ErrCounter      metric.Int64Counter
	ProcessDuration metric.Int64ValueRecorder
}

type Metric struct {
	Meter       metric.Meter
	Instruments Instruments
}

type Service struct {
	store    Store
	Tracer   trace.Tracer
	Meterics Metric
}

func New(store Store, tracer trace.Tracer, meterics Metric) *Service {
	return &Service{
		store:    store,
		Tracer:   tracer,
		Meterics: meterics,
	}
}

type Store interface {
	Upsert(ctx context.Context, a Application) error
	Get(ctx context.Context, id string) (Application, error)
	List(ctx context.Context, limit int) ([]Application, error)
}

func (svc *Service) Add(ctx context.Context, a *Application) error {
	ctx, span := svc.Tracer.Start(ctx, "service.Add")
	defer span.End()

	var errParams []string

	if a.ID == "" {
		errParams = append(errParams, "id")
	}
	if a.MinVersion == "" {
		errParams = append(errParams, "min-version")
	}
	if a.Package == "" {
		errParams = append(errParams, "package")
	}
	if len(errParams) > 0 {
		svc.Meterics.Meter.RecordBatch(ctx, []kv.KeyValue{kv.String("service.method", "Add")}, svc.Meterics.Instruments.ErrCounter.Measurement(1))
		return errs.E{Kind: errs.KindParameterErr, Parameters: errParams}
	}

	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	err := svc.store.Upsert(ctx, *a)
	if err == nil {
		return nil
	}

	svc.Meterics.Meter.RecordBatch(ctx, []kv.KeyValue{kv.String("service.method", "Add")}, svc.Meterics.Instruments.ErrCounter.Measurement(1))
	return err
}

func (svc *Service) Get(ctx context.Context, id string) (Application, error) {
	ctx, span := svc.Tracer.Start(ctx, "service.Get")
	defer span.End()

	a, err := svc.store.Get(ctx, id)
	if err != nil {
		svc.Meterics.Meter.RecordBatch(ctx, []kv.KeyValue{kv.String("service.method", "Get")}, svc.Meterics.Instruments.ErrCounter.Measurement(1))
	}
	return a, err
}

func (svc *Service) UpdateVersion(ctx context.Context, a Application) error {
	ctx, span := svc.Tracer.Start(ctx, "service.UpdateVersion")
	defer span.End()

	if a.ID == "" {
		svc.Meterics.Meter.RecordBatch(ctx, []kv.KeyValue{kv.String("service.method", "UpdateVersion")}, svc.Meterics.Instruments.ErrCounter.Measurement(1))
		return errs.E{
			Kind:    errs.KindInternal,
			Wrapped: fmt.Errorf("given struct is missing ID"),
		}
	}

	aFrom, err := svc.store.Get(ctx, a.ID)
	if err != nil {
		svc.Meterics.Meter.RecordBatch(ctx, []kv.KeyValue{kv.String("service.method", "UpdateVersion")}, svc.Meterics.Instruments.ErrCounter.Measurement(1))
		return err
	}

	a.UpdatedAt = time.Now()
	a.CreatedAt = aFrom.CreatedAt

	return svc.store.Upsert(ctx, a)
}

func (svc *Service) List(ctx context.Context, limit int) ([]Application, error) {
	ctx, span := svc.Tracer.Start(ctx, "service.List")
	defer span.End()

	applications, err := svc.store.List(ctx, limit)
	if err != nil {
		svc.Meterics.Meter.RecordBatch(ctx, []kv.KeyValue{kv.String("service.method", "List")}, svc.Meterics.Instruments.ErrCounter.Measurement(1))
		return []Application{}, err
	}

	return applications, nil
}
