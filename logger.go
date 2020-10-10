package obs

import (
	"github.com/rs/zerolog"
	"os"
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

func (l *Logger) Fatal(msg string, err error) {
	l.zl.Fatal().Err(err).Msg(msg)
}

// NewLogger returns a new Logger.
func NewLogger() *Logger {
	return &Logger{zerolog.New(os.Stderr).With().Timestamp().Logger()}
}

// NewNopLogger returns a disabled Logger for which all operation are no-op.
func NewNopLogger() *Logger {
	return &Logger{zerolog.Nop()}
}
