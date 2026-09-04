package applogger

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
)

// enrichErrorContext copies all fields from AppError to Logger context.
func (l *appLogger) enrichErrorContext(err *appfault.AppError) Logger {
	enriched := l.WithContext("ErrorType", err.Type().Name()).WithContext("ErrorCode", err.Type().Code())
	ctx := err.Context()
	for k, v := range ctx {
		enriched = enriched.WithContext(k, v)
	}

	return enriched
}

// LogError logs a structured AppError.
func (l *appLogger) LogError(err *appfault.AppError) {
	if err != nil {
		enriched := l.enrichErrorContext(err)
		enriched.(*appLogger).write(LevelError, err.Message(), err.StackTrace().String())
	}
}

// LogFaults logs an entire collection of AppErrors.
func (l *appLogger) LogFaults(faults *appfaults.Collection) {
	if faults != nil && faults.HasError() {
		for _, item := range faults.Items() {
			l.LogError(item)
		}
	}
}

// WithContext returns a child logger with key-value attached.
func (l *appLogger) WithContext(key string, val any) Logger {
	return &appLogger{
		minLevel: l.minLevel,
		sink:     l.sink,
		fields:   l.fields.Clone().Set(key, val),
	}
}

// WithFields returns a child logger with multiple fields attached.
func (l *appLogger) WithFields(fields map[string]any) Logger {
	cloned := l.fields.Clone()
	for k, v := range fields {
		cloned.Set(k, v)
	}

	return &appLogger{
		minLevel: l.minLevel,
		sink:     l.sink,
		fields:   cloned,
	}
}

// Sync flushes the underlying sink.
func (l *appLogger) Sync() error {
	return l.sink.Sync()
}

// Close flushes and releases resources held by sink.
func (l *appLogger) Close() error {
	return l.sink.Close()
}
