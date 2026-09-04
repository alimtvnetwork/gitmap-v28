package applogger

import "sync"

// CompositeSink broadcasts log entries to multiple sinks.
type CompositeSink struct {
	mu    sync.RWMutex
	sinks []LogSink
}

// NewCompositeSink creates a multi-destination broadcaster sink.
func NewCompositeSink(sinks ...LogSink) *CompositeSink {
	return &CompositeSink{
		sinks: sinks,
	}
}

// AddSink appends an additional sink destination.
func (cs *CompositeSink) AddSink(sink LogSink) *CompositeSink {
	if sink == nil {
		return cs
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.sinks = append(cs.sinks, sink)

	return cs
}

// WriteEntry forwards the entry to all configured sinks.
func (cs *CompositeSink) WriteEntry(e LogEntry) error {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	for _, sink := range cs.sinks {
		_ = sink.WriteEntry(e)
	}

	return nil
}

// Sync flushes all inner sinks.
func (cs *CompositeSink) Sync() error {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	for _, sink := range cs.sinks {
		_ = sink.Sync()
	}

	return nil
}

// Close closes all inner sinks.
func (cs *CompositeSink) Close() error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for _, sink := range cs.sinks {
		_ = sink.Close()
	}

	return nil
}
