package log

import (
	"github.com/JoinVerse/obs"
)

// Logger is the global logger.
var Logger = obs.NewLogger()

func Info(msg string) {
	Logger.Info(msg)
}

func Infof(format string, v ...interface{}) {
	Logger.Infof(format, v...)
}

func Error(err error) {
	Logger.Error(err)
}

func Errorf(msg string, err error) {
	Logger.Errorf(msg, err)
}

func Fatalf(msg string, err error) {
	Logger.Fatalf(msg, err)
}
