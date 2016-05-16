// Package diary provides a simple JSON logger.
package diary

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Default keys for log output
const (
	DefaultTimeKey    = "ts"
	DefaultLevelKey   = "lvl"
	DefaultMessageKey = "message"
)

// Level is the level of the log entry
type Level int

// Log Levels
const (
	LevelFatal = iota
	LevelError
	LevelInfo
	LevelDebug
)

type (
	// Context is a map of key/value pairs. These are Marshalled and included in the log output.
	Context map[string]interface{}

	// Logger is the actual logger. The default log level is info and the default writer is STDOUT.
	Logger struct {
		level      Level
		context    Context
		writer     io.Writer
		timeKey    string
		levelKey   string
		messageKey string
	}
)

// SetLevel creates a function that sets the log level. Generally, used when create a new logger.
func SetLevel(lvl Level) func(*Logger) {
	return func(l *Logger) {
		l.level = lvl
		return
	}
}

// SetContext creates a function that sets the context. Generally, used when create a new logger.
func SetContext(ctx Context) func(*Logger) {
	return func(l *Logger) {
		l.context = ctx
		return
	}
}

// SetWriter creates a function that will set the writer. Generally, used when create a new logger.
func SetWriter(w io.Writer) func(*Logger) {
	return func(l *Logger) {
		l.writer = w
		return
	}
}

// SetTimeKey creates a funtion that sets the time key. Generally, used when create a new logger.
func SetTimeKey(key string) func(*Logger) {
	return func(l *Logger) {
		l.timeKey = key
		return
	}
}

// SetLevelKey creates a funtion that sets the level key. Generally, used when create a new logger.
func SetLevelKey(key string) func(*Logger) {
	return func(l *Logger) {
		l.levelKey = key
		return
	}
}

// SetMessageKey creates a funtion that sets the message key. Generally, used when create a new logger.
func SetMessageKey(key string) func(*Logger) {
	return func(l *Logger) {
		l.messageKey = key
		return
	}
}

func (l *Logger) doOptions(options []func(*Logger)) {
	for _, f := range options {
		f(l)
	}
	return
}

// New creates a logger.
func New(context Context, options ...func(*Logger)) *Logger {
	l := &Logger{
		level:      LevelInfo,
		context:    context,
		writer:     os.Stdout,
		timeKey:    DefaultTimeKey,
		levelKey:   DefaultLevelKey,
		messageKey: DefaultMessageKey,
	}

	l.doOptions(options)

	return l
}

// New creates a child logger.  Initial options are inherited from the parent.
func (l *Logger) New(context Context, options ...func(*Logger)) *Logger {
	n := &Logger{
		level:      l.level,
		writer:     l.writer,
		timeKey:    l.timeKey,
		levelKey:   l.levelKey,
		messageKey: l.messageKey,
	}

	ctx := make(Context)

	for k, v := range l.context {
		ctx[k] = v
	}

	for k, v := range context {
		ctx[k] = v
	}

	n.context = ctx

	n.doOptions(options)

	return n
}

// Fatal logs a message at the "fatal" log level. It then calls os.Exit
func (l *Logger) Fatal(msg string, context ...Context) {
	l.write(LevelFatal, msg, context)
	os.Exit(-1)
}

// Error logs a message at the "error" log level.
func (l *Logger) Error(msg string, context ...Context) {
	l.write(LevelError, msg, context)
}

// Info logs a message at the "info" log level.
func (l *Logger) Info(msg string, context ...Context) {
	l.write(LevelInfo, msg, context)
}

// Debug logs a message at the "debug" log level.
func (l *Logger) Debug(msg string, context ...Context) {
	l.write(LevelDebug, msg, context)
}

func (l *Logger) write(level Level, msg string, context []Context) {
	if level > l.level {
		return
	}

	record := make(map[string]interface{}, 8)

	for k, v := range l.context {
		record[k] = v
	}

	for _, ctx := range context {
		for k, v := range ctx {
			record[k] = v
		}
	}

	record[l.timeKey] = time.Now()
	record[l.messageKey] = msg
	record[l.levelKey] = l.level.String()

	if data, err := json.Marshal(record); err == nil {
		data = append(data, '\n')
		l.writer.Write(data)
	}
}

// String returns the name of a Level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// LevelFromString returns the appropriate Level from a string name.
// Useful for parsing command line args and configuration files.
func LevelFromString(levelString string) (Level, bool) {
	switch levelString {
	case "debug":
		return LevelDebug, true
	case "info":
		return LevelInfo, true
	case "error", "eror", "err":
		return LevelError, true
	case "fatal":
		return LevelFatal, true
	default:
		return LevelDebug, false
	}
}
