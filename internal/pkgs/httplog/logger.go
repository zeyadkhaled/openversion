// Package httplog rovides utilities for logging http request with zerolog
package httplog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"

	"gitlab.innology.com.tr/zabuamer/open-telemetry-go-integration/auth"
)

func NewEvent(l zerolog.Logger, r *http.Request) *zerolog.Event {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	u := auth.UserFromContext(r.Context())

	return l.Info().
		Str("scheme", scheme).
		Str("remote_addr", r.RemoteAddr).
		Str("method", r.Method).
		Str("user_agent", r.UserAgent()).
		Str("uri", r.RequestURI).Str("request_id", middleware.GetReqID(r.Context())).
		Str("user_source", u.Source).
		Str("user_id", u.ID)
}

func NewLogFormatter(logger zerolog.Logger) middleware.LogFormatter {
	return structuredLogger{logger: logger}
}

type structuredLogger struct {
	logger zerolog.Logger
}

func (l structuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	return structuredLoggerEntry{event: NewEvent(l.logger, r)}
}

type structuredLoggerEntry struct {
	event *zerolog.Event
}

func (entry structuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	entry.event.Int("status", status).Int("length", bytes).Dur("dur", elapsed).Msg("")
}

func (entry structuredLoggerEntry) Panic(v interface{}, stack []byte) {
	entry.event.Err(nil).Str("stack", string(stack)).Str("panic", fmt.Sprintf("%+v", v)).Msg("")
}
