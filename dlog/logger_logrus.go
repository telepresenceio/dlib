package dlog

import (
	"io"
	"log"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// loggerOrEntry is an interface that lists all methods that are common for the
// logrus.Logger and logrus.Entry interfaces.
type loggerOrEntry interface {
	logrus.Ext1FieldLogger
	Log(level logrus.Level, args ...any)
	Logf(level logrus.Level, format string, args ...any)
	Logln(level logrus.Level, args ...any)
	WriterLevel(level logrus.Level) *io.PipeWriter
}

type logrusWrapper struct {
	loggerOrEntry
}

var _ Logger = logrusWrapper{}

// Helper does nothing--we use a Logrus Hook instead (see below).
func (l logrusWrapper) Helper() {}

func (l logrusWrapper) WithField(key string, value any) Logger {
	return logrusWrapper{l.loggerOrEntry.WithField(key, value)}
}

//nolint:gochecknoglobals // constant
var dlogLevel2logrusLevel = [5]logrus.Level{
	logrus.ErrorLevel,
	logrus.WarnLevel,
	logrus.InfoLevel,
	logrus.DebugLevel,
	logrus.TraceLevel,
}

func logrusLevel(level LogLevel) logrus.Level {
	if level > LogLevelTrace {
		panic(errors.Errorf("invalid LogLevel: %d", level))
	}
	return dlogLevel2logrusLevel[level]
}

func (l logrusWrapper) LogMessage(level LogLevel, message string) {
	l.Log(level, message)
}

func (l logrusWrapper) StdLogger(level LogLevel) *log.Logger {
	return log.New(l.WriterLevel(logrusLevel(level)), "", 0)
}

func (l logrusWrapper) Log(level LogLevel, args ...any) {
	l.loggerOrEntry.Log(logrusLevel(level), args...)
}

func (l logrusWrapper) Logf(level LogLevel, format string, args ...any) {
	l.loggerOrEntry.Logf(logrusLevel(level), format, args...)
}

func (l logrusWrapper) Logln(level LogLevel, args ...any) {
	l.loggerOrEntry.Logln(logrusLevel(level), args...)
}

func (l logrusWrapper) MaxLevel() LogLevel {
	var ll *logrus.Logger
	switch l := l.loggerOrEntry.(type) {
	case *logrus.Logger:
		ll = l
	case *logrus.Entry:
		ll = l.Logger
	default:
		l.Panic("unable to get logrus.Level from a %T", l)
	}
	lrv := ll.GetLevel()
	for i, l := range dlogLevel2logrusLevel {
		if l == lrv {
			return LogLevel(i)
		}
	}
	panic(errors.Errorf("invalid logrus LogLevel: %d", lrv))
}

// WrapLogrus converts a logrus *Logger into a generic Logger.
//
// You should only really ever call WrapLogrus from the initial
// process set up (i.e. directly inside your 'main()' function), and
// you should pass the result directly to WithLogger.
func WrapLogrus(in *logrus.Logger) Logger {
	in.AddHook(logrusFixCallerHook{})
	return logrusWrapper{in}
}

type logrusFixCallerHook struct{}

func (logrusFixCallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (logrusFixCallerHook) Fire(entry *logrus.Entry) error {
	if entry.Caller != nil && strings.HasPrefix(entry.Caller.Function, dlogPackageDot) {
		entry.Caller = getCaller()
	}
	return nil
}

const (
	dlogPackageDot         = "github.com/telepresenceio/dlib/v2/dlog."
	logrusPackageDot       = "github.com/sirupsen/logrus."
	maximumCallerDepth int = 25
	minimumCallerDepth int = 2 // runtime.Callers + getCaller
)

// Duplicate of logrus.getCaller() because Logrus doesn't have the
// kind if skip/.Helper() functionality that testing.TB has.
//
// https://github.com/sirupsen/logrus/issues/972
func getCaller() *runtime.Frame {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		// If the caller isn't part of this package, we're done
		if strings.HasPrefix(f.Function, logrusPackageDot) {
			continue
		}
		if strings.HasPrefix(f.Function, dlogPackageDot) {
			continue
		}
		return &f
	}

	// if we got here, we failed to find the caller's context
	return nil
}
