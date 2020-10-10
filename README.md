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
	conf := errortracking.Config{
		ServiceName:     "",
		ServiceVersion:  "",
		SentryDSN:       "",
		GCloudEnabled:   false,
		GCloudProjectID: "",
	}

	errorTracker := errortracking.New(conf)
	defer errorTracker.Close()

	err := fmt.Errorf("main: ups, that was an error")
	errorTracker.CaptureError(err, map[string]string{"key":"value"})
}
``` 


## Observer

Observer is a wrapper of Logger and Error tracker

```go
package main

import (
    "fmt"
	"github.com/JoinVerse/obs"
	"github.com/JoinVerse/obs/errtrack"
)

func main() {
	conf := errortracking.Config{
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