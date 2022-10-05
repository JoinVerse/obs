package log

import (
	"context"

	"github.com/JoinVerse/obs"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Logger is the global logger.
var Logger = obs.NewLogger()

func Info(msg string) {
	Logger.Info(msg)
}

func Infof(format string, v ...interface{}) {
	Logger.Infof(format, v...)
}

func Error(msg string, err error) {
	Logger.Error(msg, err)
}

func Errorf(err error, format string, v ...interface{}) {
	Logger.Errorf(format, err, v...)
}

// ErrorWithSpan gets the span from the provided context.
// Used for Datadog.
func ErrorWithSpan(ctx context.Context, msg string, err error) {
	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		Error(msg, err)
		return
	}

	Errorf(err, msg+", span: %v", span)
}

func ErrorfWithSpan(ctx context.Context, format string, err error, v ...interface{}) {
	span, ok := tracer.SpanFromContext(ctx)
	if !ok {
		Errorf(err, format, v...)
		return
	}

	v = append(v, []interface{}{span})
	Errorf(err, format+", span: %v", v...)
}

func Fatal(msg string, err error) {
	Logger.Fatal(msg, err)
}
