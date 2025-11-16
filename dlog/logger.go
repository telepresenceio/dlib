// Package dlog implements a generic logger facade.
//
// There are three first-class things of value in this package:
//
// First: The Logger interface.  This is a simple structured logging
// interface that is mostly trivial to implement on top of most
// logging backends, and allows library code to not need to care about
// what specific logging system the calling program uses.
//
// Second: The WithLogger and WithField functions for affecting the
// logging for a context.
//
// Third: The actual logging functions.  If you are writing library
// code and want to log things, then you should take a context.Context
// as an argument, and then call dlog.{Level}{,f,ln}(ctx, args) to
// log.
package dlog

import (
	"fmt"
	"log"
)

// PlainLogger is the most basic logger. that all loggers must implement,
// so that consumers don't need to care about the actual log implementation.
type PlainLogger interface {
	Helper()

	LogMessage(level LogLevel, message string)

	// StdLogger returns a stdlib *log.Logger that writes to this
	// Logger at the specified loglevel; for use with external
	// libraries that demand a stdlib *log.Logger.	Since
	StdLogger(LogLevel) *log.Logger
}

// GenericLogger provides the ability to first check the log-level to determine if
// logging will be made, and then format the message only when that is the case.
type GenericLogger interface {
	PlainLogger

	// Log formats then logs a message if the logger's MaxLevel is >= the given level.
	// The message is formatted using the default formats for its operands and adds
	// spaces between operands when neither is a string; in the manner of fmt.Print().
	Log(level LogLevel, args ...any)

	// Logf formats then logs a message if the logger's MaxLevel is >= the given level.
	// The message is formatted according to the format specifier; in the manner of
	// fmt.Printf().
	Logf(level LogLevel, fmt string, args ...any)

	// Logln formats then logs a message if the logger's MaxLevel is >= the given level.
	// The message is formatted using the default formats for its operands and always
	// adds spaces between operands; in the manner of fmt.Println() but without appending
	// a newline.
	Logln(level LogLevel, args ...any)

	WithField(key string, value any) Logger
}

// LoggerWithMaxLevel can be implemented by loggers that define a maximum
// level that will be logged, e.g. if a logger defines a max-level of
// LogLevelInfo, then only LogLevelError, LogLevelWarn, and LogLevelInfo will
// be logged; while LogLevelDebug and LogLevelTrace will be discarded.
//
// This interface can be used for examining what the loggers max level is
// so that resource consuming string formatting can be avoided if its known
// that the resulting message will be discarded anyway.
//
// The MaxLevel method is provided in an extra interface so that expected
// implementations of Logger that don't need a MaxLevel, such as a wrapper for
// log.Logger don't need to implement it.
type LoggerWithMaxLevel interface {
	// MaxLevel return the maximum loglevel that will be logged
	MaxLevel() LogLevel
}

// LogLevel is an abstracted common log-level type for Logger.StdLogger().
type LogLevel uint32

const (
	// LogLevelError is for errors that should definitely be noted.
	LogLevelError LogLevel = iota
	// LogLevelWarn is for non-critical entries that deserve eyes.
	LogLevelWarn
	// LogLevelInfo is for general operational entries about what's
	// going on inside the application.
	LogLevelInfo
	// LogLevelDebug is for debugging.  Very verbose logging.
	LogLevelDebug
	// LogLevelTrace is for extreme debugging.  Even finer-grained
	// informational events than the Debug.
	LogLevelTrace
)

// GenericImpl implements all level-specific functions by calling the generic function with
// the level of the specific function. It's intended to be used as the base for implementations
// that lack the level-specific functions, such as the standard logger.
type GenericImpl struct {
	PlainLogger
}

func (l GenericImpl) WithField(key string, value any) Logger {
	return &BaseLogger{GenericLogger: l}
}

func (l GenericImpl) Log(level LogLevel, args ...any) {
	l.Helper()
	l.LogMessage(level, fmt.Sprint(args...))
}

func (l GenericImpl) Logf(level LogLevel, format string, args ...any) {
	l.Helper()
	l.LogMessage(level, fmt.Sprintf(format, args...))
}

func (l GenericImpl) Logln(level LogLevel, args ...any) {
	l.Helper()
	// Trim the trailing newline; what we care about is that spaces are added in between
	// arguments, not that there's a trailing newline.
	// See also: logrus.Entry.sprintlnn
	msg := fmt.Sprintln(args...)
	l.LogMessage(level, msg[:len(msg)-1])
}

type BaseLogger struct {
	GenericLogger
}
