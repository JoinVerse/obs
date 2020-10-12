package main

import (
	"fmt"
	"github.com/JoinVerse/obs"
	"github.com/JoinVerse/obs/errtrack"
	"github.com/JoinVerse/obs/hlog"
	"github.com/JoinVerse/obs/log"
	"net/http"
)

func main() {
	conf := errtrack.Config{
		ServiceName:     "example",
		ServiceVersion:  "",
		SentryDSN:       "",
	}

	observer := obs.New(conf)
	defer observer.Close()

	err := fmt.Errorf("main: ups, that was an error")
	observer.CaptureError(err, map[string]string{"key": "value"})
	observer.Log.Error("ouch", err)

	okHandler :=http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})

	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err:= fmt.Errorf("main: ups, something wrong happend with the request")
		observer.CaptureHttpError(err, r, nil) // Report error to provider
		observer.Log.Error("mamma mia", err) // Write log to Stderr

		w.WriteHeader(500)
		fmt.Fprintf(w, "ups\n")
	})


	http.Handle("/", hlog.Logger(okHandler))
	http.Handle("/error", hlog.Logger(errorHandler))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Startup failed", err)
	}
}
