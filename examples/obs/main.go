package main

import (
	"fmt"
	"github.com/JoinVerse/obs"
	"github.com/JoinVerse/obs/errtrack"
)

func main() {
	conf := errtrack.Config{
		ServiceName:     "",
		ServiceVersion:  "",
		SentryDSN:       "",
		GCloudEnabled:   false,
		GCloudProjectID: "",
	}

	observer := obs.New(conf)
	defer observer.Close()

	err := fmt.Errorf("main: ups, that was an error")
	observer.CaptureError(err, map[string]string{"key": "value"})
	observer.Log.Error("ouch", err)
}
