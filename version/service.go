package version

import (
	"context"
	"fmt"
	"time"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
	"go.opentelemetry.io/otel/api/global"
)

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{
		store: store,
	}
}

type Store interface {
	Upsert(ctx context.Context, a Application) error
	Get(ctx context.Context, id string) (Application, error)
	List(ctx context.Context, limit int) ([]Application, error)
}

func (svc *Service) Add(ctx context.Context, a *Application) error {
	ctx, span := global.Tracer("service").Start(ctx, "service.Add")
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
		return errs.E{Kind: errs.KindParameterErr, Parameters: errParams}
	}

	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	err := svc.store.Upsert(ctx, *a)
	if err == nil {
		return nil
	}

	return err
}

func (svc *Service) Get(ctx context.Context, id string) (Application, error) {
	ctx, span := global.Tracer("service").Start(ctx, "service.Get")
	defer span.End()

	return svc.store.Get(ctx, id)
}

func (svc *Service) UpdateVersion(ctx context.Context, a Application) error {
	ctx, span := global.Tracer("service").Start(ctx, "service.UpdateVersion")
	defer span.End()

	if a.ID == "" {
		return errs.E{
			Kind:    errs.KindInternal,
			Wrapped: fmt.Errorf("given struct is missing ID"),
		}
	}

	aFrom, err := svc.store.Get(ctx, a.ID)
	if err != nil {
		return err
	}

	a.UpdatedAt = time.Now()
	a.CreatedAt = aFrom.CreatedAt

	return svc.store.Upsert(ctx, a)
}

func (svc *Service) List(ctx context.Context, limit int) ([]Application, error) {
	ctx, span := global.Tracer("service").Start(ctx, "service.List")
	defer span.End()

	applications, err := svc.store.List(ctx, limit)
	if err != nil {
		return []Application{}, err
	}

	return applications, nil
}
