package gcp

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/errorreporting"
)

// Exporter implements sending reports to google cloud.
type Exporter struct {
	errorClient *errorreporting.Client
	ctx         context.Context

	getUserFn func(r *http.Request) string
}

// New creates new error exporter that sends reports to google cloud.
func New(ctx context.Context, projectID string, serviceName string, serviceVersion string, getUserFn func(r *http.Request) string) (*Exporter, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	errorClient, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		OnError: func(err error) {
			log.Printf("Could not log error: %v", err)
		},
	})
	if err != nil {
		return nil, err
	}
	return &Exporter{errorClient: errorClient, ctx: ctx, getUserFn: getUserFn}, nil
}

// Close shutdowns the google cloud error tracker.
func (e *Exporter) Close() {
	_ = e.errorClient.Close()
}

// CaptureError send error to Google Cloud's Stack Driver.
func (e *Exporter) CaptureError(err error, tags map[string]string) {
	e.errorClient.Report(errorreporting.Entry{
		Error: err,
	})
	e.errorClient.Flush()
}

// CaptureHttpError send error to Google Cloud's Stack Driver.
func (e *Exporter) CaptureHttpError(err error, r *http.Request, tags map[string]string) {
	e.errorClient.Report(errorreporting.Entry{
		Error: err,
		Req:   r,
		User:  e.getUser(r),
	})
	e.errorClient.Flush()
}

func (e *Exporter) getUser(r *http.Request) string {
	if e.getUserFn != nil {
		return e.getUserFn(r)
	}
	var user string
	if r != nil {
		user = r.Header.Get("X-User-Id")
	}
	return user
}