package errortracking

import "net/http"

// User entity used for sentry messages.
type User struct {
	// All fields are optional
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	IP    string `json:"ip_address,omitempty"`
}

// ErrorTracker defines the interface to report errors
type ErrorTracker interface {
	Panic(rval interface{}, r *http.Request, tags map[string]string, user User)
	Error(err error, r *http.Request, tags map[string]string, user User)
	Warning(err error, r *http.Request, tags map[string]string, user User)
	InternalError(err error, tags map[string]string)
	Close()
}
