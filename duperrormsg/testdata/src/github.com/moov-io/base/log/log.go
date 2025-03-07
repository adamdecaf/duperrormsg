package log

import "fmt"

// Logger defines the behavior of a logging system
type Logger interface {
	Set(key string, value Valuer) Logger
	With(ctxs ...Context) Logger
	Details() map[string]interface{}

	Debug() Logger
	Info() Logger
	Warn() Logger
	Error() Logger
	Fatal() Logger

	Log(message string)
	Logf(format string, args ...interface{})
	Send()

	LogError(error error) LoggedError
	LogErrorf(format string, args ...interface{}) LoggedError
}

// Context interface for log contexts
type Context interface {
	Context() map[string]Valuer
}

// Fields implements a simple context
type Fields map[string]Valuer

func (f Fields) Context() map[string]Valuer {
	return f
}

// Level is a string wrapper for log levels
type Level string

// Log level constants
const Debug = Level("debug")
const Info = Level("info")
const Warn = Level("warn")
const Error = Level("error")
const Fatal = Level("fatal")

func (l Level) Context() map[string]Valuer {
	return map[string]Valuer{
		"level": String(string(l)),
	}
}

// Valuer interface for typed values
type Valuer interface {
	getValue() interface{}
}

// Mock implementation of Valuer
type mockValuer struct {
	value interface{}
}

func (m *mockValuer) getValue() interface{} {
	return m.value
}

// String creates a string Valuer
func String(s string) Valuer {
	return &mockValuer{s}
}

// Basic value creator functions - all just return a mockValuer
func Int(i int) Valuer             { return &mockValuer{i} }
func Int64(i int64) Valuer         { return &mockValuer{i} }
func Int64OrNil(i *int64) Valuer   { return &mockValuer{i} }
func Uint32(i uint32) Valuer       { return &mockValuer{i} }
func Uint64(i uint64) Valuer       { return &mockValuer{i} }
func Float32(f float32) Valuer     { return &mockValuer{f} }
func Float64(f float64) Valuer     { return &mockValuer{f} }
func Bool(b bool) Valuer           { return &mockValuer{b} }
func StringOrNil(s *string) Valuer { return &mockValuer{s} }
func ByteString(b []byte) Valuer   { return &mockValuer{b} }
func ByteBase64(b []byte) Valuer   { return &mockValuer{b} }
func Strings(vals []string) Valuer { return &mockValuer{vals} }

// Mock implementation of Logger that does nothing
type mockLogger struct{}

// All methods return the logger itself to allow chaining
func (l *mockLogger) Set(key string, value Valuer) Logger { return l }
func (l *mockLogger) With(ctxs ...Context) Logger         { return l }
func (l *mockLogger) Debug() Logger                       { return l }
func (l *mockLogger) Info() Logger                        { return l }
func (l *mockLogger) Warn() Logger                        { return l }
func (l *mockLogger) Error() Logger                       { return l }
func (l *mockLogger) Fatal() Logger                       { return l }

// These methods do nothing
func (l *mockLogger) Log(message string)                      {}
func (l *mockLogger) Logf(format string, args ...interface{}) {}
func (l *mockLogger) Send()                                   {}

// Details returns an empty map
func (l *mockLogger) Details() map[string]interface{} {
	return make(map[string]interface{})
}

// Error logging just returns a wrapped error
func (l *mockLogger) LogError(err error) LoggedError {
	return LoggedError{err: err}
}

func (l *mockLogger) LogErrorf(format string, args ...interface{}) LoggedError {
	return LoggedError{err: fmt.Errorf(format, args...)}
}

// LoggedError wraps an error
type LoggedError struct {
	err error
}

func (l LoggedError) Err() error {
	return l.err
}

func (l LoggedError) Nil() error {
	return nil
}

// Constructor functions
func NewDefaultLogger() Logger {
	return &mockLogger{}
}

func NewNopLogger() Logger {
	return &mockLogger{}
}

func NewLogFmtLogger() Logger {
	return &mockLogger{}
}

func NewJSONLogger() Logger {
	return &mockLogger{}
}

func NewTestLogger() Logger {
	return &mockLogger{}
}

func NewLogger(writer interface{}) Logger {
	return &mockLogger{}
}

// BufferedLogger type for compatibility
type BufferedLogger struct{}

func (b *BufferedLogger) Write(p []byte) (int, error) {
	return len(p), nil
}

func (b *BufferedLogger) Reset() {}

func (b *BufferedLogger) String() string {
	return ""
}

func NewBufferLogger() (*BufferedLogger, Logger) {
	return &BufferedLogger{}, &mockLogger{}
}
