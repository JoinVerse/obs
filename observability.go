package obs

import (
	"context"
	"log"

	"github.com/JoinVerse/obs/errortracking"
	"github.com/JoinVerse/obs/errortracking/gcp"
	"github.com/JoinVerse/obs/errortracking/sentry"
)

// Observer provides observer objects
type Observer struct {
	Log      ILogger
	ErrTrack *errortracking.Dispatcher
}

// NewObserver is the construct to create a new observer.
func NewObserver(config Config) (Observer, error) {
	return Observer{Log: NewLogger(), ErrTrack: NewErrorTracker(config)}, nil
}

// Close close the err trackers connections.
func (obs *Observer) Close() {
	obs.ErrTrack.Close()
}

// Config handles observer configuration
type Config struct {
	ServiceName       string
	ServiceVersion    string
	SentryDSN         string
	GCloudEnabled     bool
	GCloudProjectID   string
	Environment       string
}

// NewErrorTracker create the errorTracker
func NewErrorTracker(config Config) *errortracking.Dispatcher {
	var errorTrackers []errortracking.ErrorTracker
	if config.SentryDSN != "" {
		sentryErrorTracker, err := sentry.NewErrorTracker(config.SentryDSN, config.ServiceVersion)
		if err != nil {
			log.Printf("Error initializing Sentry Error Tracker %v\n", err)
		} else {
			errorTrackers = append(errorTrackers, sentryErrorTracker)
		}
	}

	if config.GCloudEnabled {
		gcloudErrTracker, err := gcp.NewErrorTracker(
			context.Background(),
			config.GCloudProjectID,
			config.ServiceName,
			config.ServiceVersion)
		if err != nil {
			log.Printf("Error initializing Google Cloud Error Tracker: %v\n", err)
		} else {
			errorTrackers = append(errorTrackers, gcloudErrTracker)
		}
	}

	return errortracking.NewDispatcher(errorTrackers...)
}
