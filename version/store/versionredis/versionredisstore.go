package versionredisstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog"

	tracer "go.opentelemetry.io/otel/api/trace"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
)

type Store struct {
	c       *redis.Client
	prefix  string
	timeout time.Duration
	logger  zerolog.Logger

	base   version.Store
	tracer tracer.Tracer
}

func New(addr, pass string, db int, keyPrefix string, timeout time.Duration, logger zerolog.Logger, base version.Store, tracer tracer.Tracer) (*Store, error) {
	if keyPrefix != "" && !strings.HasSuffix(keyPrefix, ":") {
		keyPrefix += ":"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	s := client.Ping()
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Store{
		c:       client,
		prefix:  keyPrefix,
		timeout: timeout,
		logger:  logger,

		base:   base,
		tracer: tracer,
	}, nil
}

func (store *Store) upsertRedis(ctx context.Context, id string, a version.Application) error {
	ctx, span := store.tracer.Start(ctx, "store.redis.upsertRedis")
	defer span.End()

	appByte, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("failed to marshal application: %v", err)
	}
	_, err = store.c.Set(store.prefix+id, string(appByte), store.timeout).Result()
	if err != nil {
		return fmt.Errorf("failed to insert application: %v", err)
	}

	return nil
}

func (store *Store) deleteRedis(ctx context.Context, id string) {
	ctx, span := store.tracer.Start(ctx, "store.redis.deleteRedis")
	defer span.End()

	_, err := store.c.Del(store.prefix + id).Result()
	if err != nil {
		store.logger.Info().Err(err).Msg("failed to delete from redis while getting")
	}
}

func (store *Store) get(ctx context.Context, id string) (version.Application, error) {
	ctx, span := store.tracer.Start(ctx, "store.redis.get")
	defer span.End()

	a, err := store.base.Get(ctx, id)
	if err != nil {
		return version.Application{}, fmt.Errorf("failed to get from base store: %v", err)
	}

	err = store.upsertRedis(ctx, id, a)
	if err != nil {
		store.logger.Info().Err(err).Msg("failed to insert to redis")
	}
	return a, nil
}

func (store *Store) Get(ctx context.Context, id string) (version.Application, error) {
	ctx, span := store.tracer.Start(ctx, "store.redis.Get")
	defer span.End()

	r, err := store.c.Get(store.prefix + id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return store.get(ctx, id)
		}
		return version.Application{}, fmt.Errorf("failed to get: %v", err)
	}

	var a version.Application
	err = json.Unmarshal([]byte(r), &a)
	if err != nil {
		store.deleteRedis(ctx, id)
		return version.Application{}, fmt.Errorf("failed to unmarshal get response: %v", err)
	}
	return a, nil
}

func (store *Store) Upsert(ctx context.Context, a version.Application) error {
	ctx, span := store.tracer.Start(ctx, "store.redis.Upsert")
	defer span.End()

	err := store.base.Upsert(ctx, a)
	store.deleteRedis(ctx, a.ID)
	return err
}

func (store *Store) List(ctx context.Context, limit int) ([]version.Application, error) {
	ctx, span := store.tracer.Start(ctx, "store.redis.List")
	span.End()

	return store.base.List(ctx, limit)
}
