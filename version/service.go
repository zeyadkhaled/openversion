package version

import (
	"context"
	"fmt"
	"time"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/auth"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/filterenc"
	"go.opentelemetry.io/otel/api/global"
)

type Service struct {
	store     Store
	filterEnc filterenc.Encer
}

func New(store Store, filterEnc filterenc.Encer) *Service {
	return &Service{
		store:     store,
		filterEnc: filterEnc,
	}
}

type Store interface {
	Upsert(ctx context.Context, a Application) error
	Get(ctx context.Context, id string) (Application, error)
	List(ctx context.Context, filter Filter, limit int) ([]Application, error)
}

func (svc *Service) Add(ctx context.Context, a *Application) error {
	u := auth.UserFromContext(ctx)
	if !u.HasRole(auth.RoleVersionsMake) {
		return errs.E{
			Kind:    errs.KindForbidden,
			Wrapped: fmt.Errorf("user not have version-make role"),
		}
	}
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

	err := svc.store.Upsert(context.Background(), *a)
	if err == nil {
		return nil
	}

	return err
}

func (svc *Service) Get(ctx context.Context, id string) (Application, error) {
	u := auth.UserFromContext(ctx)
	if !u.HasRole(auth.RoleVersionsList) {
		return Application{}, errs.E{
			Kind:    errs.KindForbidden,
			Wrapped: fmt.Errorf("user not have versions-list role"),
		}
	}

	return svc.store.Get(ctx, id)
}

func (svc *Service) UpdateVersion(ctx context.Context, a Application) error {
	u := auth.UserFromContext(ctx)
	if !u.HasRole(auth.RoleVersionsMake) {
		return errs.E{Kind: errs.KindForbidden}
	}

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

type Filter struct {
	// Rest is for sorting and pagination. LastAt is used as page selection so
	// it is also used in filtering. If LastID is empty it is not used in filtering.
	LastAt time.Time `json:"last_at,omitempty"` // createdAt
	LastID string    `json:"last_id,omitempty"`
	Older  bool      `json:"older,omitempty"`
}

type PaginatedApplications struct {
	After    string        `json:"after"`
	Current  string        `json:"current"`
	Before   string        `json:"before"`
	Elements []Application `json:"elements"`
}

func (svc *Service) List(ctx context.Context, _ Filter, cursor string, limit int) (PaginatedApplications, error) {
	tracer := global.Tracer("v2/list_versions")
	tracer.WithSpan(ctx, "Service Hit",
		func(ctx context.Context) error {
			return nil
		},
	)

	f := Filter{}
	applications, err := svc.store.List(ctx, f, limit)
	if err != nil {
		return PaginatedApplications{}, err
	}

	r := PaginatedApplications{
		Elements: applications,
	}

	c, err := svc.filterEnc.CursorFromFilter(f)
	if err != nil {
		return PaginatedApplications{}, err
	}
	r.Current = c

	if len(applications) > 0 {
		returnFilter := Filter{}

		c, err := svc.filterEnc.CursorFromFilter(returnFilter)
		if err != nil {
			return PaginatedApplications{}, err
		}
		r.Before = c

		returnFilter.LastAt = applications[len(applications)-1].CreatedAt
		returnFilter.LastID = applications[len(applications)-1].ID
		returnFilter.Older = true
		c, err = svc.filterEnc.CursorFromFilter(returnFilter)
		if err != nil {
			return PaginatedApplications{}, err
		}
		r.After = c
	}

	return r, nil
}
