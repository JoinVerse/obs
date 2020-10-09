package errortracking

import (
	"net/http"
)

// Dispatcher implements the error tracker interface
type Dispatcher struct {
	errTracks []ErrorTracker
}

// NewDispatcher creates an error tracker that dispatch the errors
// a list of other errorTrackers.
func NewDispatcher(errTrackers ...ErrorTracker) *Dispatcher {
	var d Dispatcher
	for _, et := range errTrackers {
		if et != nil {
			d.errTracks = append(d.errTracks, et)
		}
	}
	return &d
}

// Panic send panic to all other error trackers.
func (d *Dispatcher) Panic(rval interface{}, r *http.Request, tags map[string]string) {
	user := d.buildErrorTrackingUser(r)
	for _, et := range d.errTracks {
		et.Panic(rval, r, tags, user)
	}
}

// Error send error to all other error trackers.
func (d *Dispatcher) Error(err error, r *http.Request, tags map[string]string) {
	user := d.buildErrorTrackingUser(r)
	for _, et := range d.errTracks {
		et.Error(err, r, tags, user)
	}
}

// Warning send warning to all other error trackers.
func (d *Dispatcher) Warning(err error, r *http.Request, tags map[string]string) {
	user := d.buildErrorTrackingUser(r)
	for _, et := range d.errTracks {
		et.Warning(err, r, tags, user)
	}
}

// InternalError send error to all other error trackers.
func (d *Dispatcher) InternalError(err error, tags map[string]string) {
	for _, et := range d.errTracks {
		et.InternalError(err, tags)
	}
}

// Close calls each children Close.
func (d *Dispatcher) Close() {
	for _, et := range d.errTracks {
		et.Close()
	}
}

func (d *Dispatcher) buildErrorTrackingUser(r *http.Request) User {
	var user User
	if r != nil {
		user.ID = r.Header.Get("X-User-Id")
	}
	return user
}
