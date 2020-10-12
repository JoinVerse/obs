package errtrack

import (
	"context"
	"fmt"
	"github.com/JoinVerse/obs/errtrack/gcp"
	"github.com/JoinVerse/obs/errtrack/sentry"
	"net/http"
)

// SentryConfig handles Sentry exporter configuration.
type SentryConfig struct {
	SentryDSN      string
	ServiceVersion string
	OnGetUser      func(r *http.Request) sentry.User
}

// GoogleCloudErrorReportingConfig handles GoogleCloudErrorReporting configuration.
type GoogleCloudErrorReportingConfig struct {
	ServiceName     string
	ServiceVersion  string
	GCloudProjectID string
	OnGetUser       func(r *http.Request) string
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
func New() *ErrorTracker {
	return &ErrorTracker{}
}

// InitSentry initializes Sentry error tracker
func (e *ErrorTracker) InitSentry(config SentryConfig) error {
	sentryExporter, err := sentry.New(config.SentryDSN, config.ServiceVersion, config.OnGetUser)
	if err != nil {
		return fmt.Errorf("errtrack: cannot start Sentry error tracker %w", err)
	}
	e.errorExporters = append(e.errorExporters, sentryExporter)

	return nil
}

// InitGoogleCloudErrorReporting initializes Google Cloud Error Reporting
func (e *ErrorTracker) InitGoogleCloudErrorReporting(config GoogleCloudErrorReportingConfig) error {
	gcloudExporter, err := gcp.New(
		context.Background(),
		config.GCloudProjectID,
		config.ServiceName,
		config.ServiceVersion,
		config.OnGetUser)
	if err != nil {
		return fmt.Errorf("errtrack: cannot start Google Cloud Error Reporting %w", err)
	}
	e.errorExporters = append(e.errorExporters, gcloudExporter)
	return nil
}

// CaptureError sends error to all other error trackers.
func (e *ErrorTracker) CaptureError(err error, tags map[string]string) {
	for _, e := range e.errorExporters {
		e.CaptureError(err, tags)
	}
}

// CaptureHttpError sends error to all other error trackers.
func (e *ErrorTracker) CaptureHttpError(err error, r *http.Request, tags map[string]string) {
	for _, e := range e.errorExporters {
		e.CaptureHttpError(err, r, tags)
	}
}

// Close calls each children Close.
func (e *ErrorTracker) Close() {
	for _, e := range e.errorExporters {
		e.Close()
	}
}
