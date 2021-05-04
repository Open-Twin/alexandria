// Package hclog provides a Logur adapter for hclog.
package plugins
import (
	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	stdlog "log"
)

type Logger struct {
	Lusi zerolog.Logger
}

// Args are alternating key, val pairs
// keys must be strings
// vals can be any type, but display is implementation specific
// Emit a message and key/value pairs at a provided log level
func (Logger) Log(level hclog.Level, msg string, args ...interface {}) {
	switch level {
		case hclog.Error:
			log.Error().Msg(msg)
		case hclog.Warn:
			log.Warn().Msg(msg)
		case hclog.Info:
			log.Info().Msg(msg)
		case hclog.Debug:
			log.Debug().Msg(msg)
		case hclog.Trace:
			log.Trace().Msg(msg)
	}
}

// Emit a message and key/value pairs at the TRACE level
func (logger Logger) Trace(msg string, args ...interface{}){
	log.Trace().Msg(msg)
}

// Emit a message and key/value pairs at the DEBUG level
func (logger Logger) Debug(msg string, args ...interface{}) {
	log.Debug().Msg(msg)
}

// Emit a message and key/value pairs at the INFO level
func (logger Logger) Info(msg string, args ...interface {}) {
	log.Info().Msg(msg)
}

// Emit a message and key/value pairs at the WARN level
func (logger Logger) Warn(msg string, args ...interface {}){
	log.Warn().Msg(msg)
}

// Emit a message and key/value pairs at the ERROR level
func (logger Logger) Error(msg string, args ...interface {}){
	log.Error().Msg(msg)
}

// Indicate if TRACE logs would be emitted. This and the other Is* guards
// are used to elide expensive logging code based on the current level.
func (logger Logger) IsTrace() bool {
	return log.Trace().Enabled()
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func (Logger) IsDebug() bool {
	return log.Debug().Enabled()
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func (Logger) IsInfo() bool {
	return log.Info().Enabled()
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func (Logger) IsWarn() bool {
	return log.Info().Enabled()
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func (Logger) IsError() bool {
	return log.Info().Enabled()
}

// ImpliedArgs returns With key/value pairs
func (Logger) ImpliedArgs() []interface{} {
	return nil
}

// Creates a sublogger that will always have the given key/value pairs
func (Logger) With(args ...interface {}) hclog.Logger {
	//return Logger(log.With().Logger())
	return Logger{}
}

// Returns the Name of the logger
func (Logger) Name() string {
	return "Zerolog"
}

// Create a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (Logger) Named(name string) hclog.Logger {
	//return Logger(log.With().
	//	Str(name, "foo").Logger())
	return Logger{}
}

// Create a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (Logger) ResetNamed(name string) hclog.Logger {
	/*return log.With().
		Str(name, "nothing").
		Logger()*/
	return Logger{}
}

// Updates the level. This should affect all related loggers as well,
// unless they were created with IndependentLevels. If an
// implementation cannot update the level on the fly, it should no-op.
func (Logger) SetLevel(level hclog.Level) {
	zerolog.SetGlobalLevel(zerolog.Level(level))
}

// Return a value that conforms to the stdlib log.Logger interface
func (Logger) StandardLogger(opts *hclog.StandardLoggerOptions) *stdlog.Logger {
	/*logi := zerolog.New(os.Stdout).With().Str("foo","bar").Logger()
	lag := stdlog.New(logi,"",0)*/
	return nil
}

// Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (Logger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return nil
}
