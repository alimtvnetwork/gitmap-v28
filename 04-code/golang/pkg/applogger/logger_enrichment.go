package applogger

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
)

func (l *appLogger) enrichErrorContext(err *appfault.AppError) Logger {
	enriched := l.WithContext("ErrorType", err.Type().Name()).WithContext("ErrorCode", err.Type().Code())
	ctx := err.Context()
	for k, v := range ctx {
		enriched = enriched.WithContext(k, v)
	}

	return enriched
}

func (l *appLogger) LogError(err *appfault.AppError) Logger {
	if err != nil {
		enriched := l.enrichErrorContext(err)
		enriched.(*appLogger).write(LevelError, err.Message(), err.StackTrace().String())
	}

	return l
}

func (l *appLogger) LogFaults(faults *appfaults.Collection) Logger {
	if faults != nil && faults.HasError() {
		for _, item := range faults.Items() {
			l.LogError(item)
		}
	}

	return l
}

func (l *appLogger) WithContext(key string, val any) Logger {
	return &appLogger{
		minLevel: l.minLevel,
		sink:     l.sink,
		fields:   l.fields.Clone().Set(key, val),
	}
}

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

func (l *appLogger) Sync() error {
	return l.sink.Sync()
}

func (l *appLogger) Close() error {
	return l.sink.Close()
}
