package obs

import (
	"net/http"

	"cloud.google.com/go/profiler"
	"github.com/JoinVerse/obs/errtrack"
)

// Config ...
type Config struct {
	GCloudConfig errtrack.GoogleCloudErrorReportingConfig
	//When true, GCP integration is disabled.
	NOGCloudEnabled bool
	SentryConfig    errtrack.SentryConfig
}

// Observer provides observer object
type Observer struct {
	log      *Logger
	errTrack *errtrack.ErrorTracker
}

// New returns a new observer.
func New(config Config) Observer {
	log := NewLogger()
	errTrack := errtrack.New()
	if err := errTrack.InitSentry(config.SentryConfig); err != nil {
		log.Error("obs: cannot init Sentry", err)
	}

	if !config.NOGCloudEnabled {
		if err := errTrack.InitGoogleCloudErrorReporting(config.GCloudConfig); err != nil {
			log.Error("obs: cannot init GoogleCloudErrorReporting", err)
		}
		if err := profiler.Start(profiler.Config{
			Service:        config.GCloudConfig.ServiceName,
			ServiceVersion: config.GCloudConfig.ServiceVersion,
		}); err != nil {
			log.Error("obs: cannot start GoogleCloudProfiler", err)
			errTrack.CaptureError(err, nil, nil)
		}
	}
	return Observer{log: log, errTrack: errTrack}
}

// Close calls Flush, then closes any resources held by the client.
// Close should be called when the client is no longer needed.
func (o *Observer) Close() {
	o.errTrack.Close()
}

// Info logs an info message to Stderr.
func (o *Observer) Info(msg string) {
	o.log.Info(msg)
}

// Infof formats and logs an info message to Stderr.
func (o *Observer) Infof(format string, v ...interface{}) {
	o.log.Infof(format, v...)
}

// Error logs an error message to Stderr and send the error to configured trackers.
func (o *Observer) Error(msg string, err error) {
	o.ErrorTags(msg, nil, err)
}

// ErrorTags logs an error message to Stderr and send the error among the tags, to configured trackers.
func (o *Observer) ErrorTags(msg string, tags map[string]string, err error) {
	o.errTrack.CaptureError(err, tags, nil)
	o.log.Error(msg, err)
}

// ErrorTagsAndContext logs an error message to Stderr and send the error among the tags and context, to configured trackers.
func (o *Observer) ErrorTagsAndContext(msg string, tags map[string]string, context map[string]interface{}, err error) {
	o.errTrack.CaptureError(err, tags, context)
	o.log.Error(msg, err)
}

// HTTPError logs an error message to Stderr and send the error to configured trackers.
func (o *Observer) HttpError(r *http.Request, err error) {
	o.HttpErrorTags(r, nil, err)
}

// HttpErrorTags logs an error message to Stderr and send the error among the tags, to configured trackers.
func (o *Observer) HttpErrorTags(r *http.Request, tags map[string]string, err error) {
	o.errTrack.CaptureHttpError(err, r, tags, nil)
	o.log.Error("", err)
}

// HttpErrorTagsAndContext logs an error message to Stderr and send the error among the tags, to configured trackers.
func (o *Observer) HttpErrorTagsAndContext(r *http.Request, tags map[string]string, context map[string]interface{}, err error) {
	o.errTrack.CaptureHttpError(err, r, tags, context)
	o.log.Error("", err)
}

// Fatal logs a fatal message to Stderr and send the error to configured trackers.
// The os.Exit(1) function is called, which terminates the program immediately.
func (o *Observer) Fatal(msg string, err error) {
	o.errTrack.CaptureError(err, nil, nil)
	o.log.Fatal(msg, err)
}
