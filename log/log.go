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

func Error(msg string, err error) {
	Logger.Error(msg, err)
}

func Fatal(msg string, err error) {
	Logger.Fatal(msg, err)
}
