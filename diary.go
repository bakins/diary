package diary

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	DefaultTimeKey    = "ts"
	DefaultLevelKey   = "lvl"
	DefaultMessageKey = "message"
)

type Level int

const (
	LevelFatal = iota
	LevelError
	LevelInfo
	LevelDebug
)

type (
	Context map[string]interface{}

	Logger struct {
		level      Level
		context    Context
		writer     io.Writer
		timeKey    string
		levelKey   string
		messageKey string
	}
)

func SetLevel(lvl Level) func(*Logger) error {
	return func(l *Logger) error {
		l.level = lvl
		return nil
	}
}

func SetContext(ctx Context) func(*Logger) error {
	return func(l *Logger) error {
		l.context = ctx
		return nil
	}
}

func SetWriter(w io.Writer) func(*Logger) error {
	return func(l *Logger) error {
		l.writer = w
		return nil
	}
}

func TimeKey(key string) func(*Logger) error {
	return func(l *Logger) error {
		l.timeKey = key
		return nil
	}
}

func LevelKey(key string) func(*Logger) error {
	return func(l *Logger) error {
		l.levelKey = key
		return nil
	}
}

func MessageKey(key string) func(*Logger) error {
	return func(l *Logger) error {
		l.messageKey = key
		return nil
	}
}

func (l *Logger) doOptions(options []func(*Logger) error) error {
	for _, f := range options {
		if err := f(l); err != nil {
			return err
		}
	}
	return nil
}

func New(context Context, options ...func(*Logger) error) (*Logger, error) {
	l := &Logger{
		level:      LevelInfo,
		context:    context,
		writer:     os.Stdout,
		timeKey:    DefaultTimeKey,
		levelKey:   DefaultLevelKey,
		messageKey: DefaultMessageKey,
	}

	if err := l.doOptions(options); err != nil {
		return nil, err
	}

	return l, nil
}

func (l *Logger) New(context Context, options ...func(*Logger) error) (*Logger, error) {
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

	if err := n.doOptions(options); err != nil {
		return nil, err
	}

	return n, nil
}

func (l *Logger) Fatal(msg string, context ...Context) {
	l.write(LevelFatal, msg, context)
	os.Exit(-1)
}

func (l *Logger) Error(msg string, context ...Context) {
	l.write(LevelError, msg, context)
}

func (l *Logger) Info(msg string, context ...Context) {
	l.write(LevelInfo, msg, context)
}

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

// Returns the name of a Level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "dbug"
	case LevelInfo:
		return "info"
		return "warn"
	case LevelError:
		return "eror"
	case LevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// Returns the appropriate Level from a string name.
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
