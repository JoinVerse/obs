package obs

import (
	"cloud.google.com/go/profiler"
	"github.com/JoinVerse/obs/errtrack"
)

type Config struct {
	GCloudConfig errtrack.GoogleCloudErrorReportingConfig
	//When true, GCP integration is disabled.
	NOGCloudEnabled bool
	SentryConfig    errtrack.SentryConfig
}

// Observer provides observer object
type Observer struct {
	Log *Logger
	*errtrack.ErrorTracker
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
			errTrack.CaptureError(err, nil)
		}
	}
	return Observer{Log: log, ErrorTracker: errTrack}
}
