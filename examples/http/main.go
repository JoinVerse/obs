package main

import (
	"fmt"
	"net/http"

	"github.com/JoinVerse/obs"
	"github.com/JoinVerse/obs/hlog"
	"github.com/JoinVerse/obs/log"
)

func main() {
	conf := obs.Config{
		NOGCloudEnabled: true,
	}

	observer := obs.New(conf)
	defer observer.Close()

	err := fmt.Errorf("main: ups, that was an error")
	observer.ErrorTags("ouch", map[string]string{"key": "value"}, err)

	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})

	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := fmt.Errorf("main: ups, something wrong happend with the request")
		observer.HttpError(r, err) // Report error to provider

		w.WriteHeader(500)
		fmt.Fprintf(w, "ups\n")
	})

	logger := hlog.New()

	http.Handle("/", logger.Handler(okHandler))
	// Use always logger.Handler hlog.Logger is deprecated keep it here for testing backward compatibility
	http.Handle("/error", hlog.Logger(errorHandler))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Startup failed", err)
	}
}
