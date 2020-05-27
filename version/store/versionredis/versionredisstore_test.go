package versionredisstore

import (
	"context"
	"os"
	"testing"
	"time"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version/store/versionpostgre"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rs/zerolog"
)

var connStr string

func init() {
	connStr = os.Getenv("POSTGRE_CONNSTR")
}

func TestStore_Roundtrip(t *testing.T) {

	if connStr == "" {
		t.SkipNow()
	}

	approxTime := cmpopts.EquateApproxTime(time.Millisecond)
	ctx := context.Background()
	postGreStore, err := versionpostgre.New(ctx, connStr, zerolog.Nop())
	if err != nil {
		t.Fatal("cannot create postgre store", err)
	}

	app := version.Application{
		ID:         "rider_test",
		MinVersion: "25",
		Package:    "id1487640704",
	}

	redisstore, err := New("localhost:6379", "", 0, "versionstore", time.Second*3, zerolog.Nop(), postGreStore)
	if err != nil {
		t.Fatal("cannot create redis store", err)
	}

	// Insert new record
	err = redisstore.Upsert(ctx, app)
	if err != nil {
		t.Fatal("couldn't insert application", err)
	}

	// Get record from base to insert to redis with 5 seconds timeout
	_, err = redisstore.Get(ctx, app.ID)
	if err != nil {
		t.Fatal("couldn't get application", err)
	}

	// Update base store
	app.MinVersion = "26"
	err = postGreStore.Upsert(ctx, app)
	if err != nil {
		t.Fatal("couldn't update base store", err)
	}

	// Get old version from redis
	gotRedis, err := redisstore.Get(ctx, app.ID)
	if err != nil {
		t.Fatal("couldn't get application", err)
	}

	if app.MinVersion == gotRedis.MinVersion {
		t.Fatal("gotRedis shouldn't have new version")
	}

	// Wait for cache to be invalidated
	time.Sleep(time.Second * 3)

	// Get validated version from redis
	gotRedis, err = redisstore.Get(ctx, app.ID)
	if err != nil {
		t.Fatal("couldn't get application", err)
	}

	if diff := cmp.Diff(gotRedis, app, approxTime); diff != "" {
		t.Fatal("applications not matching:", diff)
	}

	// Invalidate redis store by upserting
	app.MinVersion = "27"
	err = redisstore.Upsert(ctx, app)
	if err != nil {
		t.Fatal("couldn't update application", err)
	}

	// Get from redis again which should get from base latest version
	gotRedis, err = redisstore.Get(ctx, app.ID)
	if err != nil {
		t.Fatal("couldn't get application", err)
	}

	if diff := cmp.Diff(gotRedis, app, approxTime); diff != "" {
		t.Fatal("applications not matching:", diff)
	}
}
