package sentry

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

type User sentry.User

// Exporter implements sending reports to sentry.
type Exporter struct {
	getUserFn func(r *http.Request) User
}

// New creates new error tracker that sends reports to sentry.
func New(DSN string, release string, getUserFn func(r *http.Request) User) (*Exporter, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              DSN,
		Release:          release,
		AttachStacktrace: true,
	})

	if err != nil {
		return nil, err
	}

	return &Exporter{getUserFn: getUserFn}, nil
}

// Close shutdowns the sentry error tracker.
func (e *Exporter) Close() {
	defer sentry.Flush(2 * time.Second)
}

// CaptureError send error to Sentry.
func (e *Exporter) CaptureError(err error, tags map[string]string, context map[string]string) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
		scope.SetContext("context", context)
		sentry.CaptureException(err)
	})
}

// CaptureHttpError send error to Sentry.
func (e *Exporter) CaptureHttpError(err error, r *http.Request, tags map[string]string, context map[string]string) {
	user := e.getUser(r)
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetRequest(r)
		// Add r.Body explicitly because setRequest only set it at same time is read
		// so you MUST calls SetRequest before read the body
		if r.Body != nil || r.Body != http.NoBody {
			if bodyBytes, err := ioutil.ReadAll(r.Body); err == nil{
				scope.SetRequestBody(bodyBytes)
			}
		}
		scope.SetTags(tags)
		scope.SetUser(sentry.User(user))
		scope.SetContext("context", context)
		sentry.CaptureException(err)
	})
}

func (e *Exporter) getUser(r *http.Request) User {
	if e.getUserFn != nil {
		return e.getUserFn(r)
	}
	var user User
	if r != nil {
		user.ID = r.Header.Get("X-User-Id")
		user.IPAddress = r.Header.Get("X-Forwarder-For")
	}
	return user
}
