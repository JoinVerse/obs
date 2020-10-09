package sentry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/JoinVerse/obs/errortracking"
	"github.com/getsentry/sentry-go"
)

// ErrorTracker implements sending reports to sentry.
type ErrorTracker struct {
}

// NewErrorTracker creates new error tracker that sends reports to sentry.
func NewErrorTracker(DSN string, release string) (*ErrorTracker, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              DSN,
		Release:          release,
		AttachStacktrace: true,
	})

	if err != nil {
		return nil, err
	}
	return &ErrorTracker{}, nil
}

// Close shutdowns the sentry error tracker.
func (et *ErrorTracker) Close() {
	defer sentry.Flush(2 * time.Second)
}

// Panic send panic to Sentry.
func (et *ErrorTracker) Panic(rval interface{}, r *http.Request, tags map[string]string, user errortracking.User) {
	err := fmt.Errorf("panic: %#v", rval)
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetRequest(r)
		scope.SetLevel(sentry.LevelFatal)
		scope.SetTags(tags)
		scope.SetUser(sentry.User{ID: user.ID, Email: user.Email, IPAddress: user.IP})
		sentry.CaptureException(err)
	})
}

// Error send error to Sentry.
func (et *ErrorTracker) Error(err error, r *http.Request, tags map[string]string, user errortracking.User) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetRequest(r)
		scope.SetTags(tags)
		scope.SetUser(sentry.User{ID: user.ID, Email: user.Email, IPAddress: user.IP})
		sentry.CaptureException(err)
	})
}

// Warning send warning to Sentry.
func (et *ErrorTracker) Warning(err error, r *http.Request, tags map[string]string, user errortracking.User) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetRequest(r)
		scope.SetLevel(sentry.LevelWarning)
		scope.SetTags(tags)
		scope.SetUser(sentry.User{ID: user.ID, Email: user.Email, IPAddress: user.IP})
		sentry.CaptureException(err)
	})
}

// InternalError send error to Sentry.
func (et *ErrorTracker) InternalError(err error, tags map[string]string) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
		sentry.CaptureException(err)
	})
}
