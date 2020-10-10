package main

import (
	"fmt"
	"github.com/JoinVerse/obs/errtrack"
	"github.com/JoinVerse/obs/errtrack/noop"
)

func main() {
	conf := errtrack.Config{
		ServiceName:     "",
		ServiceVersion:  "",
		SentryDSN:       "",
		SentryOnGetUser: nil,
		GCloudEnabled:   false,
		GCloudProjectID: "",
		GCloudOnGetUser: nil,
	}

	errorTracker := errtrack.New(conf)
	defer errorTracker.Close()

	err := fmt.Errorf("main: ups, that was an error")
	errorTracker.CaptureError(err, map[string]string{"key": "value"})

	// You can use the noopExporter exporter for testing purposes, it does nothing
	noopExporter := noop.New()
	defer noopExporter.Close()
	noopExporter.CaptureError(err, map[string]string{"os":"Darwin"})

}
