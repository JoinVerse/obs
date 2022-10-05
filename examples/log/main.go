package main

import (
	"context"
	"fmt"

	"github.com/JoinVerse/obs/log"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {

	log.Info("hello world")
	log.Infof("hello %s", "world")

	err := fmt.Errorf("be water my friend")
	log.Error("He said", err)

	log.Errorf(err, "the problem are the %s", "logs")

	trace := mocktracer.Start()
	defer trace.Stop()

	ctx := context.Background()
	span, ctx := tracer.StartSpanFromContext(ctx, "obs log", tracer.ResourceName("obs"))

	log.ErrorWithSpan(ctx, "error while tracing", err)

	log.ErrorfWithSpan(ctx, "this obs log has a %s", err, "string")

	span.Finish()

	log.Fatal("hello world", err)
}
