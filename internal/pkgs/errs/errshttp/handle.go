// Package errshttp provides function to convert err to http response
package errshttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	"golang.org/x/text/language"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/errs"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/httplog"
	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/internal/pkgs/langutil"
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

var msgi18n = map[language.Tag]map[errs.Kind]string{
	language.English: {
		errs.KindInternal:         "Internal server error.",
		errs.KindUnauthorized:     "Unauthorized. Please login again.",
		errs.KindLoginPassword:    "Wrong login credentials.",
		errs.KindForbidden:        "Forbidden.",
		errs.KindNotFound:         "Not found.",
		errs.KindConflict:         "Conflict.",
		errs.KindDuplicate:        "Duplicate.",
		errs.KindParameterErr:     "Incorrect field.",
		errs.KindTooManyRequests:  "Too many requests.",
		errs.KindDependentService: "Error in upstream service.",
		errs.KindBlocked:          "Your account has been blocked. You can contact us for more information.",
	},
	language.Turkish: {
		errs.KindInternal:         "Sunucu hatası.",
		errs.KindUnauthorized:     "İzinsiz giriş. Lütfen tekrar giriş yapınız.",
		errs.KindLoginPassword:    "Hatalı giriş.",
		errs.KindForbidden:        "Yasak işlem.",
		errs.KindNotFound:         "Bulunamadı.",
		errs.KindConflict:         "Çelişen istek.",
		errs.KindDuplicate:        "Çift istek.",
		errs.KindParameterErr:     "Hatalı alan.",
		errs.KindTooManyRequests:  "İstek limit aşımı.",
		errs.KindDependentService: "Bağlı servislerde hata.",
		errs.KindBlocked:          "Hesabınız engellenmiştir. Detaylı bilgi için bizimle iletişime geçebilirsiniz.",
	},
}

func Handle(l zerolog.Logger, w http.ResponseWriter, r *http.Request, err error) {

	// Default E is internal server error which is what unknown errors are so no
	// need to check if conversion succeed.
	var e errs.E
	errors.As(err, &e)

	var m string
	if e.PublicMsg != "" {
		m = e.PublicMsg
	} else {
		msgs, ok := msgi18n[langutil.Language(r.Context())]
		if !ok {
			msgs = msgi18n[language.English]
		}
		m = msgs[e.Kind]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus[e.Kind])
	_ = json.NewEncoder(w).Encode(errResp{
		Msg:        m,
		Parameters: e.Parameters,
	})

	event := httplog.NewEvent(l, r).Err(err)
	if e.Kind == errs.KindParameterErr {
		event = event.Strs("failed_params", e.Parameters)
	}
	event.Msg("request not successful")
}

type errResp struct {
	Msg        string   `json:"msg"`
	Parameters []string `json:"params,omitempty"`
}
