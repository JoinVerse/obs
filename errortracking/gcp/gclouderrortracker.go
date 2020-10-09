package gcp

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/errorreporting"
	"github.com/JoinVerse/obs/errortracking"
)

// ErrorTracker implements sending reports to google cloud.
type ErrorTracker struct {
	errorClient *errorreporting.Client
	ctx         context.Context
}

// NewErrorTracker creates new error tracker that sends reports to google cloud.
func NewErrorTracker(ctx context.Context, projectID string, serviceName string, serviceVersion string) (*ErrorTracker, error) {
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
	return &ErrorTracker{errorClient: errorClient, ctx: ctx}, nil
}

// Close shutdowns the google cloud error tracker.
func (et *ErrorTracker) Close() {
	_ = et.errorClient.Close()
}

// Panic send panic to Google Cloud's Stack Driver.
func (et *ErrorTracker) Panic(rval interface{}, r *http.Request, tags map[string]string, user errortracking.User) {
	err := fmt.Errorf("panic: %#v", rval)
	et.errorClient.Report(errorreporting.Entry{
		Error: err,
		Req:   r,
		User:  buildUser(user),
	})
	et.errorClient.Flush()
}

// Error send error to Google Cloud's Stack Driver.
func (et *ErrorTracker) Error(err error, r *http.Request, tags map[string]string, user errortracking.User) {
	et.errorClient.Report(errorreporting.Entry{
		Error: err,
		Req:   r,
		User:  buildUser(user),
	})
	et.errorClient.Flush()
}

// Warning send warning to Google Cloud's Stack Driver.
func (et *ErrorTracker) Warning(err error, r *http.Request, tags map[string]string, user errortracking.User) {
	et.errorClient.Report(errorreporting.Entry{
		Error: err,
		Req:   r,
		User:  buildUser(user),
	})
	et.errorClient.Flush()
}

// InternalError send error to Google Cloud's Stack Driver.
func (et *ErrorTracker) InternalError(err error, tags map[string]string) {
	et.errorClient.Report(errorreporting.Entry{
		Error: err,
	})
	et.errorClient.Flush()
}

func buildUser(user errortracking.User) string {
	return fmt.Sprintf("%s:%s", user.ID, user)
}
