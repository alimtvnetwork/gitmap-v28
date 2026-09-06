package applogger

import "fmt"

type (
	ZapLoggerInterface interface {
		Debug(args ...any)
		Info(args ...any)
		Warn(args ...any)
		Error(args ...any)
		Fatal(args ...any)
		Sync() error
	}

	ZapAdapter struct {
		zapLogger ZapLoggerInterface
	}
)

// NewZapAdapter wraps a Zap logger instance into a LogSink.
func NewZapAdapter(zapLogger ZapLoggerInterface) *ZapAdapter {
	return &ZapAdapter{zapLogger: zapLogger}
}

// dispatchLevel routes the message to the corresponding Zap level.
func (za *ZapAdapter) dispatchLevel(lvl LogLevel, msg string) {
	switch lvl {
	case LevelDebug:
		za.zapLogger.Debug(msg)
	case LevelWarn:
		za.zapLogger.Warn(msg)
	case LevelError:
		za.zapLogger.Error(msg)
	case LevelFatal:
		za.zapLogger.Fatal(msg)
	default:
		za.zapLogger.Info(msg)
	}
}

// WriteEntry forwards log entry to the Zap logger.
func (za *ZapAdapter) WriteEntry(e LogEntry) error {
	if za.zapLogger == nil {
		return nil
	}

	msg := e.Message
	if len(e.Fields) > 0 {
		msg += fmt.Sprintf(" fields=%s", e.Fields.Format())
	}

	za.dispatchLevel(e.Level, msg)

	return nil
}

// Sync delegates to Zap logger sync.
func (za *ZapAdapter) Sync() error {
	if za.zapLogger != nil {
		return za.zapLogger.Sync()
	}

	return nil
}

// Close closes the zap adapter.
func (za *ZapAdapter) Close() error {
	return za.Sync()
}
