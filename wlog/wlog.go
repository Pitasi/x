package wlog

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

// Set sets the given attributes to the widelog for this context.
// Arguments are converted to attributes as if by [slog.Logger.Log].
//
// If the context doesn't contain a wlog, this is a no-op.
func Set(ctx context.Context, args ...any) {
	logger, ok := ctx.Value(evtKey).(*Event)
	if ok {
		logger.set(args...)
	}
}

// Group collects several attributes under a single key on a log line.
//
// If the context doesn't contain a wlog, this is a no-op.
func Group(ctx context.Context, name string, args ...any) {
	logger, ok := ctx.Value(evtKey).(*Event)
	if ok {
		logger.group(name, args...)
	}
}

// MustLog sets a flag on the underlaying [Event] so that the [SamplerFunc]
// function is aware that this instance should be logged.
func MustLog(ctx context.Context) {
	logger, ok := ctx.Value(evtKey).(*Event)
	if ok {
		logger.setMustLog()
	}
}

// Prepare adds a new [Event] to the context.
func Prepare(ctx context.Context) context.Context {
	logger := new()
	return context.WithValue(ctx, evtKey, logger)
}

// GetEvent returns the [Event] attached to this context, or nil.
func GetEvent(ctx context.Context) *Event {
	evt := ctx.Value(evtKey)
	if evt != nil {
		return evt.(*Event)
	}
	return nil
}

// Enrich adds the arguments of the current [Event] to the [*slog.Logger].
func Enrich(ctx context.Context, l *slog.Logger) *slog.Logger {
	wlog, ok := ctx.Value(evtKey).(*Event)
	if !ok {
		return l
	}
	for _, args := range wlog.args {
		l = l.With(args...)
	}
	return l
}

type Event struct {
	args    [][]any
	mustLog bool
}

func new() *Event { return &Event{} }

func (l *Event) setMustLog() { l.mustLog = true }

func (l *Event) set(args ...any) { l.args = append(l.args, args) }

func (l *Event) group(name string, args ...any) {
	l.args = append(l.args, []any{slog.Group(name, args...)})
}

type key string

var evtKey = key("wlog")

// Middleware is a reference implementation of a [http.Handler] that will
// emit a single widelog at the end of some of the incoming requests.
//
// Use the [SamplerFunc] argument to decide which requests should be logged. A
// reference implementation can be found at [NewSampler].
func Middleware(h http.Handler, baseLog *slog.Logger, s SamplerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := Prepare(r.Context())
		start := time.Now()

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)

		level := s(ctx, GetEvent(ctx))
		Enrich(ctx, baseLog).Log(r.Context(), level, "request", "took", time.Since(start))
	})
}

// LevelIgnore is a [slog.Level] that shouldn't be logged.
var LevelIgnore = slog.Level(-999)

// SamplerFunc decides whether an [Event] should be logged or not.
type SamplerFunc func(context.Context, *Event) slog.Level

// NewSampler returns a [Sampler] that will log:
// - all events with a "err" or "error" field, with an error level
// - all events that were [MustLog]ged
// - a random sample of events based on the given [0-1] rate
func NewSampler(rate float64) SamplerFunc {
	return func(ctx context.Context, evt *Event) slog.Level {
		// Always log events that have "err" or "error" fields.
		for _, g := range evt.args {
			for _, key := range g {
				if key == "err" || key == "error" {
					return slog.LevelError
				}
			}
		}

		// Sample all the rest, unless they have mustLog set.
		if !evt.mustLog && rand.Float64() > rate {
			return LevelIgnore
		}

		return slog.LevelInfo
	}
}
