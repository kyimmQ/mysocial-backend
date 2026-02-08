// Package logger provides a standardized logging wrapper.
//
// Wraps a structured logger (zerolog, zap, or slog) behind
// a clean interface so the rest of the app doesn't depend
// on a specific logging library.
//
// Example:
//
//	type Logger interface {
//	    Info(msg string, fields ...Field)
//	    Error(msg string, err error, fields ...Field)
//	    Debug(msg string, fields ...Field)
//	    Warn(msg string, fields ...Field)
//	    With(fields ...Field) Logger
//	}
package logger
