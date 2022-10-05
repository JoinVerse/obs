package obs

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger implements interface.
type Logger struct {
	zl zerolog.Logger
}

func (l *Logger) Info(msg string) {
	l.zl.Info().Msg(msg)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.zl.Info().Msgf(format, v...)
}

func (l *Logger) Error(msg string, err error) {
	l.zl.Err(err).Msg(msg)
}

func (l *Logger) Errorf(format string, err error, v ...interface{}) {
	l.zl.Err(err).Msgf(format, v...)
}

func (l *Logger) Fatal(msg string, err error) {
	l.zl.Fatal().Err(err).Msg(msg)
}

// NewLogger returns a new Logger.
func NewLogger() *Logger {
	host, _ := os.Hostname()
	return &Logger{zerolog.New(os.Stderr).With().Timestamp().Str("host", host).Logger()}
}

// NewNopLogger returns a disabled Logger for which all operation are no-op.
func NewNopLogger() *Logger {
	return &Logger{zerolog.Nop()}
}
