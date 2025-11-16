package dlog

import (
	"log"
)

type splitLogger struct {
	outLog GenericLogger
	errLog GenericLogger
}

func (l splitLogger) WithField(key string, value any) Logger {
	return &BaseLogger{GenericLogger: splitLogger{outLog: l.outLog.WithField(key, value), errLog: l.errLog.WithField(key, value)}}
}

var _ GenericLogger = splitLogger{}

func (l splitLogger) levelLogger(level LogLevel) (lg GenericLogger) {
	if level == LogLevelError {
		lg = l.errLog
	} else {
		lg = l.outLog
	}
	return lg
}

func (l splitLogger) Helper() {
	l.outLog.Helper()
}

func (l splitLogger) LogMessage(level LogLevel, message string) {
	l.levelLogger(level).Log(level, message)
}

func (l splitLogger) StdLogger(level LogLevel) *log.Logger {
	return l.levelLogger(level).StdLogger(level)
}

func (l splitLogger) Log(level LogLevel, args ...any) {
	l.levelLogger(level).Log(level, args...)
}

func (l splitLogger) Logf(level LogLevel, format string, args ...any) {
	l.levelLogger(level).Logf(level, format, args...)
}

func (l splitLogger) Logln(level LogLevel, args ...any) {
	l.levelLogger(level).Logln(level, args...)
}

// NewSplitLogger creates a logger that logs to two different loggers, depending on the log level.
// This type of logger is useful for separating logs into stdout and stderr when running as a
// service so that the service logs can be easily filtered.
//
//   - Error level: errLog
//   - Other levels: outLog
func NewSplitLogger(outLog GenericLogger, errLog GenericLogger) Logger {
	return &BaseLogger{GenericLogger: &splitLogger{outLog: outLog, errLog: errLog}}
}
