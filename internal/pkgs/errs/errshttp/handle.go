package errshttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
)

var httpStatus = map[errs.Kind]int{
	errs.KindInternal:         http.StatusInternalServerError,
	errs.KindUnauthorized:     http.StatusUnauthorized,
	errs.KindLoginPassword:    http.StatusUnauthorized,
	errs.KindForbidden:        http.StatusForbidden,
	errs.KindNotFound:         http.StatusNotFound,
	errs.KindConflict:         http.StatusConflict,
	errs.KindDuplicate:        http.StatusConflict,
	errs.KindParameterErr:     http.StatusBadRequest,
	errs.KindTooManyRequests:  http.StatusTooManyRequests,
	errs.KindDependentService: http.StatusServiceUnavailable,
	errs.KindBlocked:          http.StatusForbidden,
}

func Handle(l zerolog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	var e errs.E
	errors.As(err, &e)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus[e.Kind])
	_ = json.NewEncoder(w).Encode(errResp{
		Parameters: e.Parameters,
	})
}

type errResp struct {
	Parameters []string `json:"params,omitempty"`
}
