package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/api/kv"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs/errshttp"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/version"
)

type versionAPI struct {
	svc    *version.Service
	logger zerolog.Logger
}

func newVersionAPI(svc *version.Service, logger zerolog.Logger) versionAPI {
	return versionAPI{
		svc:    svc,
		logger: logger,
	}
}

func (api versionAPI) Routes(r chi.Router) {
	r.Post("/v2/version", api.newVersion)
	r.Get("/v2/version/{appid}", api.getVersion)
	r.Put("/v2/version/{appid}", api.updateVersion)
	r.Get("/v2/version", api.listVersions)
}

func (api versionAPI) newVersion(w http.ResponseWriter, r *http.Request) {
	defer processDuration(r.Context(), time.Now(), "newVersion", api.svc.Meterics)

	var app version.Application
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		errshttp.Handle(api.logger, w, r, errs.E{
			Kind:    errs.KindParameterErr,
			Wrapped: err,
		})
		return
	}

	err = api.svc.Add(r.Context(), &app)
	if err != nil {
		errshttp.Handle(api.logger, w, r, err)
		return
	}

	jsonWrite(api.logger, w, http.StatusCreated, "{}")
}

func (api versionAPI) getVersion(w http.ResponseWriter, r *http.Request) {
	defer processDuration(r.Context(), time.Now(), "getVersion", api.svc.Meterics)

	app, err := api.svc.Get(r.Context(), chi.URLParam(r, "appid"))
	if err != nil {
		errshttp.Handle(api.logger, w, r, err)
		return
	}
	jsonWrite(api.logger, w, http.StatusOK, app)
}

func (api versionAPI) updateVersion(w http.ResponseWriter, r *http.Request) {
	defer processDuration(r.Context(), time.Now(), "updateVersion", api.svc.Meterics)

	id := chi.URLParam(r, "appid")

	var app version.Application
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		errshttp.Handle(api.logger, w, r, errs.E{
			Wrapped: err,
		})
		return
	}

	if id != "" && id != app.ID {
		errshttp.Handle(api.logger, w, r, errs.E{
			Kind:       errs.KindParameterErr,
			Parameters: []string{"id"},
			Wrapped:    err,
		})
		return
	}

	err = api.svc.UpdateVersion(r.Context(), app)
	if err != nil {
		errshttp.Handle(api.logger, w, r, err)
		return
	}

	jsonWrite(api.logger, w, http.StatusOK, "{}")
}

func (api versionAPI) listVersions(w http.ResponseWriter, r *http.Request) {
	defer processDuration(r.Context(), time.Now(), "listVersions", api.svc.Meterics)

	lim := 100
	if l, err := strconv.Atoi(r.FormValue("limit")); err == nil {
		lim = l
	}

	resp, err := api.svc.List(r.Context(), lim)
	if err != nil {
		errshttp.Handle(api.logger, w, r, err)
		return
	}

	jsonWrite(api.logger, w, http.StatusOK, resp)
}

func processDuration(ctx context.Context, start time.Time, endpoint string, metric version.Metric) {
	elapsed := time.Since(start)
	metric.Meter.RecordBatch(ctx, []kv.KeyValue{kv.String("api.handler", endpoint)}, metric.Instruments.ProcessDuration.Measurement(elapsed.Seconds()))
}

func jsonWrite(logger zerolog.Logger, w http.ResponseWriter, status int, t interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(t)
	if err != nil {
		logger.Info().Err(err).Msg("cannot write response")
	}
}
