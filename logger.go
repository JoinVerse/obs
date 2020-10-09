package obs

import (
	"github.com/rs/zerolog"
	"os"
)

// Logger contains logger interface.
type Logger struct {
	zl zerolog.Logger
}

type ILogger interface {
	Info(msg string)
	Infof(format string, v ...interface{})
	Error(err error)
	Errorf(msg string, err error)
	Fatalf(msg string, err error)
}

func (l *Logger) Info(msg string) {
	l.zl.Info().Msg(msg)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.zl.Info().Msgf(format, v...)
}

func (l *Logger) Error(err error) {
	l.zl.Err(err).Msg("")
}

func (l *Logger) Errorf(msg string, err error) {
	l.zl.Err(err).Msg(msg)
}

func (l *Logger) Fatalf(msg string, err error) {
	l.zl.Fatal().Err(err).Msg(msg)
}

// NewLogger returns a new logger.
func NewLogger() ILogger {
	return &Logger{zerolog.New(os.Stderr).With().Timestamp().Logger()}
}

// NewNopLogger returns a disabled logger for which all operation are no-op.
func NewNopLogger() ILogger {
	return &Logger{zerolog.Nop()}
}

