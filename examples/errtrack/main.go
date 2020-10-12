package main

import (
	"fmt"
	"github.com/JoinVerse/obs/errtrack"
	"github.com/JoinVerse/obs/errtrack/noop"
)

func main() {
	errorTracker := errtrack.New()
	_ = errorTracker.InitGoogleCloudErrorReporting(errtrack.GoogleCloudErrorReportingConfig{})
	_ = errorTracker.InitSentry(errtrack.SentryConfig{})
	defer errorTracker.Close()

	err := fmt.Errorf("main: ups, that was an error")
	errorTracker.CaptureError(err, map[string]string{"key": "value"})

	// You can use the noopExporter exporter for testing purposes, it does nothing
	noopExporter := noop.New()
	defer noopExporter.Close()
	noopExporter.CaptureError(err, map[string]string{"os": "Darwin"})

}
