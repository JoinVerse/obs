package obs

import (
	"cloud.google.com/go/profiler"
	"github.com/JoinVerse/obs/errtrack"
)

// Observer provides observer object
type Observer struct {
	Log *Logger
	*errtrack.ErrorTracker
}

// New returns a new observer.
func New(errConfig errtrack.Config) Observer {
	obs :=  Observer{Log: NewLogger(), ErrorTracker: errtrack.New(errConfig)}

	if err := profiler.Start(profiler.Config{
		Service:        errConfig.ServiceName,
		ServiceVersion: errConfig.ServiceVersion,
	}); err != nil {
		obs.Log.Error("obs: cannot start profiling", err)
		obs.CaptureError(err, nil)
	}
	return obs
}
