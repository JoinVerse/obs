# obs
Observability module for GO applications

[![CodeQL](https://github.com/JoinVerse/obs/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/JoinVerse/obs/actions/workflows/codeql-analysis.yml)

## Observer

Observer is a wrapper of Logger and Error tracker, it also enables [Cloud Profiler](https://cloud.google.com/profiler/).
This is the recommended way to use this module.

```go
package main

import (
    "fmt"
	"github.com/JoinVerse/obs"
)

func main() {
	conf := obs.Config{
		NOGCloudEnabled: true,
	}

	observer := obs.New(conf)
	defer observer.Close()

    observer.Info("Starting program")

	err := fmt.Errorf("main: ups, that was an error")
	observer.ErrorTags("message", map[string]string{"key": "value"}, map[string]string{"body": "{'hi':'bye'}"}, err)
}
```


## net/http

This module also provides functionality to be used with `net/http`. See how to use it [here](github.com/JoinVerse/obs/examples/http/main.go)

- `errtrack.CaptureHTTPError` capture requests information along with the user id if `X-User-ID` header has being set. Also `context` is used to send more context about the error there you can send until 8kb of data.
- `htop.Logger` is a middleware that logs end of each request, along with some useful data about what was requested, 
what the response status was, and how long it took to return.


## Logs
It provides a simple interface for logging based on [zerolog](https://github.com/rs/zerolog)

> Avoid package-global logger https://github.com/uber-go/zap/blob/master/FAQ.md#why-include-package-global-loggers

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

Error tracking provides an interface to send your errors to different providers, it supports [sentry](sentry.io) and 
[GCP Error Reporting](https://cloud.google.com/error-reporting).

> On GKE, you must add the cloud-platform access scope when creating the cluster, as the following example command shows:
> `gcloud container clusters create example-cluster-name --scopes https://www.googleapis.com/auth/cloud-platform`

```go
package main

import (
	"fmt"
	"github.com/JoinVerse/obs/errtrack"
)

func main() {
	errorTracker := errtrack.New()
    _ = errorTracker.InitGoogleCloudErrorReporting(errtrack.GoogleCloudErrorReportingConfig{})
    _ = errorTracker.InitSentry(errtrack.SentryConfig{})
	defer errorTracker.Close()

	err := fmt.Errorf("main: ups, that was an error")
	errorTracker.CaptureError(err, map[string]string{"key":"value"}, nil)
}
``` 
