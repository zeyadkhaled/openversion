package versionpostgre

import (
	"context"
	"os"
	"testing"
	"time"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rs/zerolog"
)

var connStr string

func init() {
	connStr = os.Getenv("POSTGRE_CONNSTR")
}

func TestStore_RoundTrip(t *testing.T) {
	if connStr == "" {
		t.SkipNow()
	}

	approxTime := cmpopts.EquateApproxTime(time.Millisecond)
	ctx := context.Background()

	store, err := New(ctx, connStr, zerolog.Nop())
	if err != nil {
		t.Fatal("cannot create store", err)
	}

	app := version.Application{
		ID:         "rider_iOS",
		MinVersion: "25",
		Package:    "id1487640704",
	}

	// insert test
	err = store.Upsert(ctx, app)
	if err != nil {
		t.Fatal("couldn't insert application", err)
	}

	got, err := store.Get(ctx, app.ID)
	if err != nil {
		t.Fatal("couldn't get application", err)
	}
	if diff := cmp.Diff(got, app, approxTime); diff != "" {
		t.Fatal("applications not matching:", diff)
	}

	// update test
	app.MinVersion = "26"
	err = store.Upsert(ctx, app)
	if err != nil {
		t.Fatal("couldn't update application", err)
	}
	got, err = store.Get(ctx, app.ID)
	if err != nil {
		t.Fatal("couldn't get application", err)
	}
	if diff := cmp.Diff(got, app, approxTime); diff != "" {
		t.Fatal("applications not matching:", diff)
	}

	// list test
	gotList, err := store.List(ctx, version.Filter{}, 1)
	if err != nil {
		t.Fatal("couldn't list applications", err)
	}
	if len(gotList) != 1 {
		t.Fatal("fewer elements than expected", err)
	}
}
