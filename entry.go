package log

import (
	"fmt"
	stdlog "log"
	"os"
	"sort"
	"strings"
	"time"
)

// Fields represents a map of entry level data used for structured logging.
type Fields map[string]interface{}

// Names returns field names sorted.
// map is not
func (f Fields) Names() (v []string) {
	for k := range f {
		v = append(v, k)
	}

	sort.Strings(v)
	return
}

// Get field value by name.
func (f Fields) Get(name string) interface{} {
	return f[name]
}

// Entry defines a single log entry
type Entry struct {
	logger *logger
	start  time.Time
	fields []Fields // private used; store all fields when withFields is called.  improve performance.

	Level     Level     `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Fields    Fields    `json:"fields"` // single map; easy to use for handlers
}

func newEntry(l *logger) Entry {
	e := Entry{}
	e.logger = l
	e.fields = l.defaultFields
	return e
}

// Debug level message.
func (e Entry) Debug(msg string) {
	e.Level = DebugLevel
	e.Message = msg
	handler(e)
}

// Debugf level message.
func (e Entry) Debugf(msg string, v ...interface{}) {
	e.Level = DebugLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Info level message.
func (e Entry) Info(msg string) {
	e.Level = InfoLevel
	e.Message = msg
	handler(e)
}

// Infof level message.
func (e Entry) Infof(msg string, v ...interface{}) {
	e.Level = InfoLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Warn level message.
func (e Entry) Warn(msg string) {
	e.Level = WarnLevel
	e.Message = msg
	handler(e)
}

// Warnf level message.
func (e Entry) Warnf(msg string, v ...interface{}) {
	e.Level = WarnLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Error level message.
func (e Entry) Error(msg string) {
	e.Level = ErrorLevel
	e.Message = msg
	handler(e)
}

// Errorf level message.
func (e Entry) Errorf(msg string, v ...interface{}) {
	e.Level = ErrorLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
}

// Panic level message.
func (e Entry) Panic(msg string) {
	e.Level = PanicLevel
	e.Message = msg
	handler(e)
	os.Exit(1)
}

// Panicf level message.
func (e Entry) Panicf(msg string, v ...interface{}) {
	e.Level = PanicLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
	os.Exit(1)
}

// Fatal level message.
func (e Entry) Fatal(msg string) {
	e.Level = FatalLevel
	e.Message = msg
	handler(e)
	os.Exit(1)
}

// Fatalf level message.
func (e Entry) Fatalf(msg string, v ...interface{}) {
	e.Level = FatalLevel
	e.Message = fmt.Sprintf(msg, v...)
	handler(e)
	os.Exit(1)
}

// Str add string field to current entry
func (e Entry) Str(key string, val string) Entry {
	return e.WithFields(Fields{key: val})
}

// Bool add bool field to current entry
func (e Entry) Bool(key string, val bool) Entry {
	return e.WithFields(Fields{key: val})
}

// Int add Int field to current entry
func (e Entry) Int(key string, val int) Entry {
	return e.WithFields(Fields{key: val})
}

// Int8 add Int8 field to current entry
func (e Entry) Int8(key string, val int8) Entry {
	return e.WithFields(Fields{key: val})
}

// Int16 add Int16 field to current entry
func (e Entry) Int16(key string, val int16) Entry {
	return e.WithFields(Fields{key: val})
}

// Int32 add Int32 field to current entry
func (e Entry) Int32(key string, val int32) Entry {
	return e.WithFields(Fields{key: val})
}

// Int64 add Int64 field to current entry
func (e Entry) Int64(key string, val int64) Entry {
	return e.WithFields(Fields{key: val})
}

// Uint add Uint field to current entry
func (e Entry) Uint(key string, val uint) Entry {
	return e.WithFields(Fields{key: val})
}

// Uint8 add Uint8 field to current entry
func (e Entry) Uint8(key string, val uint8) Entry {
	return e.WithFields(Fields{key: val})
}

// Uint16 add Uint16 field to current entry
func (e Entry) Uint16(key string, val uint16) Entry {
	return e.WithFields(Fields{key: val})
}

// Uint32 add Uint32 field to current entry
func (e Entry) Uint32(key string, val uint32) Entry {
	return e.WithFields(Fields{key: val})
}

// Uint64 add Uint64 field to current entry
func (e Entry) Uint64(key string, val uint64) Entry {
	return e.WithFields(Fields{key: val})
}

// Float32 add Float32 field to current entry
func (e Entry) Float32(key string, val float32) Entry {
	return e.WithFields(Fields{key: val})
}

// Float64 add Float64 field to current entry
func (e Entry) Float64(key string, val float64) Entry {
	return e.WithFields(Fields{key: val})
}

// WithField returns a new entry with the `key` and `value` set.
func (e Entry) WithField(key string, value interface{}) Entry {
	return e.WithFields(Fields{key: value})
}

// WithFields adds the provided fields to the current entry
func (e Entry) WithFields(fields Fields) Entry {
	f := make([]Fields, 0, len(e.fields)+len(fields))
	f = append(f, e.fields...)
	f = append(f, fields)

	e.fields = f
	return e
}

// WithError returns a new entry with the "error" set to `err`.
func (e Entry) WithError(err error) Entry {
	if err == nil {
		return e
	}
	return e.WithField("error", fmt.Sprintf("%+v", err))
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (e Entry) Trace(msg string) Entry {
	e.Message = msg
	e.start = time.Now().UTC()
	return e
}

// mergedFields returns the fields list collapsed into a single map.
func (e Entry) mergedFields() Fields {
	f := Fields{}

	for _, fields := range e.fields {
		for k, v := range fields {
			f[k] = v
		}
	}

	return f
}

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

func duration(d time.Duration) string {
	if d < day {
		return d.String()
	}

	var b strings.Builder

	if d >= year {
		years := d / year
		fmt.Fprintf(&b, "%dy", years)
		d -= years * year
	}

	days := d / day
	d -= days * day
	fmt.Fprintf(&b, "%dd%s", days, d)

	return b.String()
}

// Stop should be used with Trace, to fire off the completion message. When
// an `err` is passed the "error" field is set, and the log level is error.
func (e Entry) Stop() {
	e.WithField("duration", duration(time.Since(e.start))).Info(e.Message)
}

func handler(e Entry) {
	// I guess we don't need to lock here and the performance can be improved
	// e.logger.rwMutex.RLock()
	// defer e.logger.rwMutex.RUnlock()

	for _, h := range e.logger.cacheLeveledHandlers(e.Level) {
		e.Timestamp = time.Now().UTC()
		e.Fields = e.mergedFields()
		err := h.Log(e)
		if err != nil {
			stdlog.Printf("log: log failed: %v", err)
		}
	}
}
