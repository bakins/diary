// Package diary provides a simple JSON logger.
package diary

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
)

// Default keys for log output
const (
	DefaultTimeKey    = "ts"
	DefaultLevelKey   = "lvl"
	DefaultMessageKey = "message"
	DefaultCallerKey  = "caller"
	DefaultCallerSkip = 2
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

	// OptionsFunc is a function passed to new for setting options on a new logger.
	OptionsFunc func(*Logger) error

	// Context is a map of key/value pairs. These are Marshalled and included in the log output.
	Context map[string]interface{}

	// Logger is the actual logger. The default log level is debug and the default writer is STDOUT.
	Logger struct {
		level      Level
		context    Context
		writer     io.Writer
		timeKey    string
		levelKey   string
		messageKey string
		callerKey  string
		callerSkip int
	}

	// A Value generates a log value. It represents a dynamic value which is re-evaluated with each log event.
	Value struct {
		Func interface{}
	}
)

var (
	defaultLogger *Logger
)

func init() {
	defaultLogger, _ = New(nil)
}

// GetDefaultLogger returns a Logger with the default settings
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// SetLevel creates a function that sets the log level. Generally, used when create a new logger.
func SetLevel(lvl Level) func(*Logger) error {
	return func(l *Logger) error {
		l.level = lvl
		return nil
	}
}

// SetContext creates a function that sets the context. Generally, used when create a new logger.
func SetContext(ctx Context) func(*Logger) error {
	return func(l *Logger) error {
		l.context = ctx
		return nil
	}
}

// SetWriter creates a function that will set the writer. Generally, used when create a new logger.
func SetWriter(w io.Writer) func(*Logger) error {
	return func(l *Logger) error {
		l.writer = w
		return nil
	}
}

// SetTimeKey creates a funtion that sets the time key. Generally, used when create a new logger.
func SetTimeKey(key string) func(*Logger) error {
	return func(l *Logger) error {
		l.timeKey = key
		return nil
	}
}

// SetLevelKey creates a funtion that sets the level key. Generally, used when create a new logger.
func SetLevelKey(key string) func(*Logger) error {
	return func(l *Logger) error {
		l.levelKey = key
		return nil
	}
}

// SetMessageKey creates a funtion that sets the message key. Generally, used when create a new logger.
func SetMessageKey(key string) func(*Logger) error {
	return func(l *Logger) error {
		l.messageKey = key
		return nil
	}
}

// SetCallerKey creates a function that will set the caller key. Generally, used when create a new logger.
func SetCallerKey(key string) func(*Logger) error {
	return func(l *Logger) error {
		l.callerKey = key
		return nil
	}
}

// SetCallerSkip creates a function that will set the caller stack skip. Generally, used when create a new logger.
// This can be used if you call the logger fromyour own utility function but want the caller info for the caller
// of your own utility function rather than the utility function.
func SetCallerSkip(i int) func(*Logger) error {
	return func(l *Logger) error {
		if i < DefaultCallerSkip {
			return fmt.Errorf("caller ship must be >= %d", DefaultCallerSkip)
		}
		l.callerSkip = i
		return nil
	}
}

func (l *Logger) doOptions(options []OptionsFunc) error {
	for _, f := range options {
		if err := f(l); err != nil {
			return err
		}
	}
	return nil
}

// New creates a logger.
func New(context Context, options ...OptionsFunc) (*Logger, error) {
	l := &Logger{
		level:      LevelDebug,
		context:    context,
		writer:     os.Stdout,
		timeKey:    DefaultTimeKey,
		levelKey:   DefaultLevelKey,
		messageKey: DefaultMessageKey,
		callerKey:  DefaultCallerKey,
		callerSkip: DefaultCallerSkip,
	}

	if err := l.doOptions(options); err != nil {
		return nil, err
	}

	return l, nil
}

// New creates a copy of the logger with additional options.  Initial options are inherited from the original.
// The two loggers are independent.
func (l *Logger) New(context Context, options ...OptionsFunc) (*Logger, error) {
	n := &Logger{
		callerKey:  l.callerKey,
		level:      l.level,
		writer:     l.writer,
		timeKey:    l.timeKey,
		levelKey:   l.levelKey,
		messageKey: l.messageKey,
		callerSkip: l.callerSkip,
	}

	ctx := make(Context)

	for k, v := range l.context {
		ctx[k] = v
	}

	for k, v := range context {
		ctx[k] = v
	}

	n.context = ctx

	if err := n.doOptions(options); err != nil {
		return nil, err
	}

	return n, nil
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

	record[l.timeKey] = time.Now().Format(time.RFC3339Nano)
	record[l.messageKey] = msg
	record[l.levelKey] = l.level.String()
	record[l.callerKey] = caller(l.callerSkip)

	if data, err := json.Marshal(record); err == nil {
		data = append(data, '\n')
		l.writer.Write(data)
	} else {
		fmt.Println(err)
	}

}

var levelsMap = map[Level]string{
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelError: "error",
	LevelFatal: "fatal",
}

// String returns the name of a Level.
func (l Level) String() string {
	if v, ok := levelsMap[l]; ok {
		return v
	}
	return "unknown"
}

// LevelFromString returns the appropriate Level from a string name.
// Useful for parsing command line args and configuration files.
func LevelFromString(levelString string) (Level, error) {
	switch levelString {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "error", "eror", "err":
		return LevelError, nil
	case "fatal":
		return LevelFatal, nil
	default:
		return LevelDebug, fmt.Errorf("Unknown level: %v", levelString)
	}
}

// Fatal uses the default logger to log a message at the "fatal" log level. It then calls os.Exit
func Fatal(msg string, context ...Context) {
	defaultLogger.Fatal(msg, context...)
}

// Error uses the default logger to log a message at the "error" log level.
func Error(msg string, context ...Context) {
	defaultLogger.Error(msg, context...)
}

// Info uses the default logger to log a message at the "info" log level.
func Info(msg string, context ...Context) {
	defaultLogger.Info(msg, context...)
}

// Debug uses the default logger to log a message at the "debug" log level.
func Debug(msg string, context ...Context) {
	defaultLogger.Debug(msg, context...)
}

func (v Value) MarshalJSON() ([]byte, error) {
	// copied from evaluateLazy in log15

	t := reflect.TypeOf(v.Func)
	// we do not currently add errors to the log entries, so these may go unnoticed
	// TODO: handled json errors better
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("not func: %+v", v.Func)
	}

	if t.NumIn() > 0 {
		return nil, fmt.Errorf("func takes args: %+v", v.Func)
	}

	if t.NumOut() == 0 {
		return nil, fmt.Errorf("no func return val: %+v", v.Func)
	}
	value := reflect.ValueOf(v.Func)
	results := value.Call([]reflect.Value{})
	if len(results) == 1 {
		return json.Marshal(results[0].Interface())
	} else {
		values := make([]interface{}, len(results))
		for i, v := range results {
			values[i] = v.Interface()
		}
		return json.Marshal(values)
	}
}
