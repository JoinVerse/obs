# obs
Observability module for GO applications

## Logs
It provides a simple interface for logging based on [zerolog](https://github.com/rs/zerolog):

```go
package main

import (
    "fmt"
    "github.com/JoinVerse/obs/log"
)

func main() {

    log.Info("hello world")
    log.Infof("hello %s", "world")
    
    err := fmt.Errorf("be water my friend")
    log.Error("main: He said", err)

    log.Fatal("main: cannot start service", err)
}

// Output: {"level":"info","time":"2020-10-10T13:21:37+02:00","message":"hello world"}
// Output: {"level":"info","time":"2020-10-10T13:21:37+02:00","message":"hello world"}
// Output: {"level":"error","error":"be water my friend","time":"2020-10-10T13:21:37+02:00","message":"main: He said"}
// Output: {"level":"fatal","error":"be water my friend","time":"2020-10-10T13:21:37+02:00","message":"main: cannot start service"}
// Output: exit status 1
```

## Error tracking

Error tracking provides an interface to send your errors to different providers, it supports [sentry](sentry.io) and [GCP Error Reporting](https://cloud.google.com/error-reporting)

```go
package main

import (
	"fmt"
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

	errorTracker := errtrack.New(conf)
	defer errorTracker.Close()

	err := fmt.Errorf("main: ups, that was an error")
	errorTracker.CaptureError(err, map[string]string{"key":"value"})
}
``` 


## Observer

Observer is a wrapper of Logger and Error tracker, it also enables [Cloud Profiler](https://cloud.google.com/profiler/) 

```go
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
	observer.Log.Error("main: ouch", err)
}
```


## net/http

This module also provides functionality to be used with `net/http`. See how to use it [here](github.com/JoinVerse/obs/examples/http/main.go)

- `errtrack.CaptureHttpError` capture requests information along with the user id if `X-User-ID` header has being set.
- `htop.Logger` is a middleware that logs end of each request, along with some useful data about what was requested, what the response status was, and how long it took to return.