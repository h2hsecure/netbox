package server

import (
	"io"
	"log"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
)

type internalLogger struct {
	log *zerolog.Logger
}

// Debug implements hclog.Logger.
func (i *internalLogger) Debug(msg string, args ...interface{}) {
	i.log.Debug().Msgf(msg, args...)
}

// Error implements hclog.Logger.
func (i *internalLogger) Error(msg string, args ...interface{}) {
	i.log.Error().Msgf(msg, args...)
}

// GetLevel implements hclog.Logger.
func (i *internalLogger) GetLevel() hclog.Level {
	panic("unimplemented")
}

// ImpliedArgs implements hclog.Logger.
func (i *internalLogger) ImpliedArgs() []interface{} {
	panic("unimplemented")
}

// Info implements hclog.Logger.
func (i *internalLogger) Info(msg string, args ...interface{}) {
	i.log.Info().Msgf(msg, args...)
}

// IsDebug implements hclog.Logger.
func (i *internalLogger) IsDebug() bool {
	return i.log.GetLevel() >= zerolog.DebugLevel
}

// IsError implements hclog.Logger.
func (i *internalLogger) IsError() bool {
	return i.log.GetLevel() >= zerolog.ErrorLevel
}

// IsInfo implements hclog.Logger.
func (i *internalLogger) IsInfo() bool {
	return i.log.GetLevel() >= zerolog.InfoLevel
}

// IsTrace implements hclog.Logger.
func (i *internalLogger) IsTrace() bool {
	return i.log.GetLevel() >= zerolog.TraceLevel
}

// IsWarn implements hclog.Logger.
func (i *internalLogger) IsWarn() bool {
	return i.log.GetLevel() >= zerolog.WarnLevel
}

// Log implements hclog.Logger.
func (i *internalLogger) Log(level hclog.Level, msg string, args ...interface{}) {
	i.log.GetLevel()
}

// Name implements hclog.Logger.
func (i *internalLogger) Name() string {
	return "zerolog"
}

// Named implements hclog.Logger.
func (i *internalLogger) Named(name string) hclog.Logger {
	return i
}

// ResetNamed implements hclog.Logger.
func (i *internalLogger) ResetNamed(name string) hclog.Logger {
	return i
}

// SetLevel implements hclog.Logger.
func (i *internalLogger) SetLevel(level hclog.Level) {

}

// StandardLogger implements hclog.Logger.
func (i *internalLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return nil
}

// StandardWriter implements hclog.Logger.
func (i *internalLogger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return i.log.With().Logger()
}

// Trace implements hclog.Logger.
func (i *internalLogger) Trace(msg string, args ...interface{}) {
	i.log.Trace().Msgf(msg, args...)
}

// Warn implements hclog.Logger.
func (i *internalLogger) Warn(msg string, args ...interface{}) {
	i.log.Warn().Msgf(msg, args...)
}

// With implements hclog.Logger.
func (i *internalLogger) With(args ...interface{}) hclog.Logger {
	return i
}
