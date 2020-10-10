package errtrack

import (
	"context"
	"github.com/JoinVerse/obs/errtrack/gcp"
	"github.com/JoinVerse/obs/errtrack/sentry"
	"log"
	"net/http"
)

// Config handles error tracker exporters configuration
type Config struct {
	ServiceName     string
	ServiceVersion  string
	SentryDSN       string
	SentryOnGetUser func(r *http.Request) sentry.User
	GCloudEnabled   bool
	GCloudProjectID string
	GCloudOnGetUser func(r *http.Request) string
}

type ErrorTracker struct {
	errorExporters []errorExporter
}

// errorExporter defines the interface to export errors to providers
type errorExporter interface {
	CaptureError(err error, tags map[string]string)
	CaptureHttpError(err error, r *http.Request, tags map[string]string)
	Close()
}

// New creates a new ErrorTracker
func New(config Config) *ErrorTracker {
	var errorExporters []errorExporter
	if config.SentryDSN != "" {
		sentryExporter, err := sentry.New(config.SentryDSN, config.ServiceVersion, config.SentryOnGetUser)
		if err != nil {
			log.Printf("CaptureError initializing Sentry CaptureError Tracker %v\n", err)
		} else {
			errorExporters = append(errorExporters, sentryExporter)
		}
	}

	if config.GCloudEnabled {
		gcloudExporter, err := gcp.New(
			context.Background(),
			config.GCloudProjectID,
			config.ServiceName,
			config.ServiceVersion,
			config.GCloudOnGetUser)
		if err != nil {
			log.Printf("CaptureError initializing Google Cloud CaptureError Tracker: %v\n", err)
		} else {
			errorExporters = append(errorExporters, gcloudExporter)
		}
	}

	return &ErrorTracker{errorExporters}
}

// CaptureError send error to all other error trackers.
func (et *ErrorTracker) CaptureError(err error, tags map[string]string) {
	//user := et.buildErrorTrackingUser(r)
	for _, et := range et.errorExporters {
		et.CaptureError(err, tags)
	}
}

// Close calls each children Close.
func (et *ErrorTracker) Close() {
	for _, et := range et.errorExporters {
		et.Close()
	}
}
