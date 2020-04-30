package log

import (
	"context"
	stdlog "log"
	"sync"
)

// Logger is the default instance of the log package
var (
	_logger = new()
)

// Handler is an interface that log handlers need to be implemented
type Handler interface {
	Log(Entry) error
}

// Flusher is an interface that allow handles have the ability to clear buffer and close connection
type Flusher interface {
	Flush() error
}

type logger struct {
	handles              []Handler
	leveledHandlers      map[Level][]Handler
	cacheLeveledHandlers func(level Level) []Handler
	defaultFields        []Fields
	rwMutex              sync.RWMutex
}

func new() *logger {
	logger := logger{
		leveledHandlers: map[Level][]Handler{},
		defaultFields:   []Fields{},
	}

	logger.cacheLeveledHandlers = logger.getLeveledHandlers()
	return &logger
}

func (l *logger) getLeveledHandlers() func(level Level) []Handler {
	debugHandlers := l.leveledHandlers[DebugLevel]
	infoHandlers := l.leveledHandlers[InfoLevel]
	warnHandlers := l.leveledHandlers[WarnLevel]
	errorHandlers := l.leveledHandlers[ErrorLevel]
	panicHandlers := l.leveledHandlers[PanicLevel]
	fatalHandlers := l.leveledHandlers[FatalLevel]

	return func(level Level) []Handler {
		switch level {
		case DebugLevel:
			return debugHandlers
		case InfoLevel:
			return infoHandlers
		case WarnLevel:
			return warnHandlers
		case ErrorLevel:
			return errorHandlers
		case PanicLevel:
			return panicHandlers
		case FatalLevel:
			return fatalHandlers
		}

		return []Handler{}
	}
}

// RegisterHandler adds a new Log Handler and specifies what log levels
// the handler will be passed log entries for
func RegisterHandler(handler Handler, levels ...Level) {
	_logger.rwMutex.Lock()
	defer _logger.rwMutex.Unlock()

	for _, level := range levels {
		_logger.leveledHandlers[level] = append(_logger.leveledHandlers[level], handler)
	}

	_logger.handles = append(_logger.handles, handler)
	_logger.cacheLeveledHandlers = _logger.getLeveledHandlers()
}

// Debug level formatted message.
func Debug(msg string) {
	e := newEntry(_logger)
	e.Debug(msg)
}

// Debugf level formatted message.
func Debugf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Debugf(msg, v...)
}

// Info level formatted message.
func Info(msg string) {
	e := newEntry(_logger)
	e.Info(msg)
}

// Infof level formatted message.
func Infof(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Infof(msg, v...)
}

// Warn level formatted message.
func Warn(msg string) {
	e := newEntry(_logger)
	e.Warn(msg)
}

// Warnf level formatted message.
func Warnf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Warnf(msg, v...)
}

// Error level formatted message
func Error(msg string) {
	e := newEntry(_logger)
	e.Error(msg)
}

// Errorf level formatted message
func Errorf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Errorf(msg, v...)
}

// Panic level formatted message
func Panic(msg string) {
	e := newEntry(_logger)
	e.Panic(msg)
}

// Panicf level formatted message
func Panicf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Panicf(msg, v...)
}

// Fatal level formatted message, followed by an exit.
func Fatal(msg string) {
	e := newEntry(_logger)
	e.Fatal(msg)
}

// Fatalf level formatted message, followed by an exit.
func Fatalf(msg string, v ...interface{}) {
	e := newEntry(_logger)
	e.Fatalf(msg, v...)
}

// Str add string field to current entry
func Str(key string, val string) Entry {
	e := newEntry(_logger)
	return e.Str(key, val)
}

// Bool add bool field to current entry
func Bool(key string, val bool) Entry {
	e := newEntry(_logger)
	return e.Bool(key, val)
}

// Int add Int field to current entry
func Int(key string, val int) Entry {
	e := newEntry(_logger)
	return e.Int(key, val)
}

// Int8 add Int8 field to current entry
func Int8(key string, val int8) Entry {
	e := newEntry(_logger)
	return e.Int8(key, val)
}

// Int16 add Int16 field to current entry
func Int16(key string, val int16) Entry {
	e := newEntry(_logger)
	return e.Int16(key, val)
}

// Int32 add Int32 field to current entry
func Int32(key string, val int32) Entry {
	e := newEntry(_logger)
	return e.Int32(key, val)
}

// Int64 add Int64 field to current entry
func Int64(key string, val int64) Entry {
	e := newEntry(_logger)
	return e.Int64(key, val)
}

// Uint add Uint field to current entry
func Uint(key string, val uint) Entry {
	e := newEntry(_logger)
	return e.Uint(key, val)
}

// Uint8 add Uint8 field to current entry
func Uint8(key string, val uint8) Entry {
	e := newEntry(_logger)
	return e.Uint8(key, val)
}

// Uint16 add Uint16 field to current entry
func Uint16(key string, val uint16) Entry {
	e := newEntry(_logger)
	return e.Uint16(key, val)
}

// Uint32 add Uint32 field to current entry
func Uint32(key string, val uint32) Entry {
	e := newEntry(_logger)
	return e.Uint32(key, val)
}

// Uint64 add Uint64 field to current entry
func Uint64(key string, val uint64) Entry {
	e := newEntry(_logger)
	return e.Uint64(key, val)
}

// Float32 add float32 field to current entry
func Float32(key string, val float32) Entry {
	e := newEntry(_logger)
	return e.Float32(key, val)
}

// Float64 add Float64 field to current entry
func Float64(key string, val float64) Entry {
	e := newEntry(_logger)
	return e.Float64(key, val)
}

// Flush clear all handler's buffer
func Flush() {
	for _, h := range _logger.handles {
		flusher, ok := h.(Flusher)
		if ok {
			err := flusher.Flush()
			if err != nil {
				stdlog.Printf("log: flush log handler: %v", err)
			}
		}
	}
}

// WithFields returns a log Entry with fields set
func WithFields(fields Fields) Entry {
	e := newEntry(_logger)
	return e.WithFields(fields)
}

// WithField returns a new entry with the `key` and `value` set.
func WithField(key string, value interface{}) Entry {
	e := newEntry(_logger)
	return e.WithField(key, value)
}

// WithDefaultFields adds fields to every entry instance
func WithDefaultFields(fields Fields) {
	f := make([]Fields, 0, len(_logger.defaultFields)+len(fields))
	f = append(f, _logger.defaultFields...)
	f = append(f, fields)

	_logger.defaultFields = f
}

// WithError returns a new entry with the "error" set to `err`.
func WithError(err error) Entry {
	e := newEntry(_logger)
	return e.WithError((err))
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func Trace(msg string) Entry {
	e := newEntry(_logger)
	return e.Trace(msg)
}

var (
	ctxKey = &struct {
		name string
	}{
		name: "log",
	}
)

// NewContext return a new context with a logger value
func NewContext(ctx context.Context, e Entry) context.Context {
	return context.WithValue(ctx, ctxKey, e)
}

// FromContext return a logger from the context
func FromContext(ctx context.Context) Entry {
	v := ctx.Value(ctxKey)
	if v == nil {
		return newEntry(_logger)
	}

	return v.(Entry)
}
