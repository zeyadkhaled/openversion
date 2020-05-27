package versioninmem

import (
	"context"
	"sort"
	"strings"
	"sync"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
)

type Store struct {
	applications map[string]version.Application
	lock         sync.Mutex
}

func New() *Store {
	return &Store{
		applications: make(map[string]version.Application),
	}
}

func (s *Store) Upsert(ctx context.Context, a version.Application) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.applications[a.ID] = a
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (version.Application, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	a, ok := s.applications[id]
	if !ok {
		return a, errs.E{Kind: errs.KindNotFound}
	}

	return a, nil

}

func (s *Store) List(ctx context.Context, f version.Filter, limit int) ([]version.Application, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	applications := []version.Application{}

	for _, v := range s.applications {
		if (f.Older != v.CreatedAt.Before(f.LastAt)) ||
			(f.LastAt.Equal(v.CreatedAt) &&
				(strings.Compare(f.LastID, v.ID) < 0) != f.Older) {

			continue
		}

		applications = append(applications, v)
	}

	sort.Slice(applications, func(i, j int) bool {
		if applications[i].CreatedAt.Before(applications[j].CreatedAt) {
			return f.Older
		} else if applications[i].CreatedAt.After(applications[j].CreatedAt) {
			return !f.Older
		} else {
			return f.Older == (applications[i].ID < applications[j].ID)
		}
	})

	return applications, nil
}
