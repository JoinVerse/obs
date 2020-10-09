package noop

import (
	"net/http"

	"github.com/JoinVerse/obs/errortracking"
)

// ErrorTracker does nothing but implements the ErrorTracker interface
type ErrorTracker struct {
}

// NewErrorTracker creates new NoopErrorTracker.
func NewErrorTracker() *ErrorTracker {
	return &ErrorTracker{}
}

// Panic send panic to nowhere.
func (*ErrorTracker) Panic(rval interface{}, r *http.Request, tags map[string]string, user errortracking.User) {
}

// Error send error to nowhere.
func (*ErrorTracker) Error(err error, r *http.Request, tags map[string]string, user errortracking.User) {
}

// Warning send warning to nowhere.
func (*ErrorTracker) Warning(err error, r *http.Request, tags map[string]string, user errortracking.User) {
}

// InternalError send error to nowhere.
func (*ErrorTracker) InternalError(err error, tags map[string]string) {
}

// Close does nothing.
func (*ErrorTracker) Close() {
}
