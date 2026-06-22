package logger

import "log"

// Logger defines logging behavior expected by commands.
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

// StdLogger is a minimal production logger.
type StdLogger struct {
	level string
}

func New(level string) *StdLogger {
	if level == "" {
		level = "info"
	}
	return &StdLogger{level: level}
}

func (l *StdLogger) Infof(format string, args ...any) {
	log.Printf("INFO: "+format, args...)
}

func (l *StdLogger) Errorf(format string, args ...any) {
	log.Printf("ERROR: "+format, args...)
}
