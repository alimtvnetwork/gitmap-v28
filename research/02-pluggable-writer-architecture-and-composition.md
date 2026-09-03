# Pluggable Logger Writer Architecture & Composition Blueprint

> **Status:** Proposal & Architectural Specification  
> **Date:** 2026-09-03  
> **Target System:** `04-code/golang/pkg/applogger`  
> **Topic:** Composable Writer Contracts, BaseWriter Embedding, Configurable Formatting, and Enterprise REST API Streaming  

---

## 1. Executive Summary & Core Requirements

This document presents the complete architectural blueprint for a **composable, multi-writer logging system** designed around:
1. **Interface Segregation & Sub-Contracts:** Decoupling byte-formatting (`Formatter`), transport destination (`Destination`), filtering (`FilterPolicy`), and event listening (`EventHook`).
2. **BaseWriter Embedding:** A shared `base.Writer` providing thread-safety, level filtering, prefixing, and hook dispatch so external developers can create custom writers with minimal boilerplate.
3. **Multi-Package Layout:** Writers are separated into distinct packages (`writers/base`, `writers/text`, `writers/json`, `writers/sqlite`, `writers/restapi`), preventing monolithic bloat and circular dependencies.
4. **Zero-Output Silent Mode:** The logger starts with zero writers by default. When no writers are registered, calls exit immediately without allocations.
5. **Configurable Writer Behaviors:** Writers are not hardcoded to `os.Stdout` or fixed formats; prefixes, envelopes, timestamp styles, and custom `io.Writer` sinks are fully injectible.
6. **Production REST API Sink:** Features non-blocking channel queuing, background batch worker pooling, retry backoff, custom HTTP headers, and a local Dead Letter Queue (DLQ) fallback.

---

## 2. Recommended Package Layout

```
04-code/golang/pkg/applogger/
├── interfaces.go                     # Core Logger & LogWriter contracts
├── logger.go                         # Fluent Logger & dispatch engine
├── record.go                         # LogRecord normalization model
│
└── writers/
    ├── base/
    │   ├── writer.go                 # BaseWriter with mutex, filter, hooks
    │   └── hooks.go                  # EventHook definitions (OnWrite, OnError)
    ├── text/
    │   ├── writer.go                 # TextWriter with color & prefix options
    │   └── formatter.go              # Console & colored text formatters
    ├── json/
    │   ├── writer.go                 # JSONWriter with envelope & key remapping
    │   └── formatter.go              # Fast JSON serializer & envelope builders
    ├── sqlite/
    │   ├── writer.go                 # SQLiteWriter with auto-migration
    │   └── schema.go                 # Table creation & batch insert queries
    └── restapi/
        ├── writer.go                 # RestAPIWriter with queue & worker loop
        ├── client.go                 # HTTP transport, auth headers, retry backoff
        └── dlq.go                    # Dead letter queue fallback sink
```

---

## 3. Core Contracts & Interfaces

### 3.1 `pkg/applogger/interfaces.go`
```go
package applogger

import (
	"context"
	"time"
)

// LogLevel defines standardized severity tiers.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// LogRecord carries normalized event data across all writers.
type LogRecord struct {
	Timestamp  time.Time      `json:"timestamp"`
	Level      LogLevel       `json:"level"`
	Message    string         `json:"message"`
	Context    context.Context `json:"-"`
	Fields     map[string]any `json:"fields,omitempty"`
	TraceID    string         `json:"traceId,omitempty"`
	UserID     string         `json:"userId,omitempty"`
	Caller     string         `json:"caller,omitempty"`
}

// LogWriter is the universal contract required for all destinations.
type LogWriter interface {
	Name() string
	WriteLog(ctx context.Context, record LogRecord) error
	Sync() error
	Close() error
}
```

### 3.2 Composable Sub-Contracts
```go
package applogger

import "io"

// Formatter serializes a LogRecord into a byte slice.
type Formatter interface {
	Format(record LogRecord) ([]byte, error)
}

// Destination represents an abstract output sink.
type Destination interface {
	io.Writer
}

// FilterPolicy determines whether a specific record should be emitted.
type FilterPolicy interface {
	ShouldLog(record LogRecord) bool
}

// EventHook allows intercepting logging lifecycle events.
type EventHook interface {
	OnBeforeWrite(record *LogRecord)
	OnAfterWrite(record LogRecord)
	OnError(err error, record LogRecord)
}
```

---

## 4. The Foundation: `BaseWriter` (Composable Embedding)

The `BaseWriter` provides reusable plumbing. Any custom or default writer embeds `BaseWriter` to inherit thread-safe locks, level gating, custom prefixing, and lifecycle hooks:

```go
package base

import (
	"context"
	"sync"

	"coding-guidelines/common/pkg/applogger"
)

// Writer provides the baseline implementation for all custom log writers.
type Writer struct {
	mu           sync.RWMutex
	name         string
	minLevel     applogger.LogLevel
	prefix       string
	hooks        []applogger.EventHook
	errorHandler func(err error, record applogger.LogRecord)
}

func New(name string, minLevel applogger.LogLevel) *Writer {
	return &Writer{
		name:     name,
		minLevel: minLevel,
		hooks:    make([]applogger.EventHook, 0),
	}
}

func (b *Writer) Name() string {
	return b.name
}

func (b *Writer) SetPrefix(prefix string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.prefix = prefix
}

func (b *Writer) Prefix() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.prefix
}

func (b *Writer) SetMinLevel(level applogger.LogLevel) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.minLevel = level
}

func (b *Writer) AddHook(hook applogger.EventHook) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.hooks = append(b.hooks, hook)
}

func (b *Writer) SetErrorHandler(handler func(err error, record applogger.LogRecord)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.errorHandler = handler
}

// IsEnabled checks whether the record passes the level gate.
func (b *Writer) IsEnabled(record applogger.LogRecord) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return record.Level >= b.minLevel
}

// NotifyBeforeWrite triggers registered pre-write hooks.
func (b *Writer) NotifyBeforeWrite(record *applogger.LogRecord) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.hooks {
		h.OnBeforeWrite(record)
	}
}

// NotifyAfterWrite triggers registered post-write hooks.
func (b *Writer) NotifyAfterWrite(record applogger.LogRecord) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, h := range b.hooks {
		h.OnAfterWrite(record)
	}
}

// HandleError forwards write failures to custom error handlers and hooks.
func (b *Writer) HandleError(err error, record applogger.LogRecord) {
	b.mu.RLock()
	handler := b.errorHandler
	hooks := b.hooks
	b.mu.RUnlock()

	if handler != nil {
		handler(err, record)
	}
	for _, h := range hooks {
		h.OnError(err, record)
	}
}
```

---

## 5. Concrete Default Writers

### 5.1 `TextWriter` (`pkg/applogger/writers/text`)
Configurable human-readable output. Users can change the destination (`os.Stdout`, `os.Stderr`, file buffer), toggle ANSI colors, and specify custom prefixes:

```go
package text

import (
	"context"
	"fmt"
	"io"
	"sync"

	"coding-guidelines/common/pkg/applogger"
	"coding-guidelines/common/pkg/applogger/writers/base"
)

type Options struct {
	Destination io.Writer
	MinLevel    applogger.LogLevel
	Prefix      string
	EnableColor bool
	TimeFormat  string
}

type Writer struct {
	*base.Writer
	destination io.Writer
	enableColor bool
	timeFormat  string
	ioMu        sync.Mutex
}

func New(opts Options) *Writer {
	if opts.Destination == nil {
		opts.Destination = io.Discard
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = "15:04:05.000"
	}

	w := &Writer{
		Writer:      base.New("text-writer", opts.MinLevel),
		destination: opts.Destination,
		enableColor: opts.EnableColor,
		timeFormat:  opts.TimeFormat,
	}
	if opts.Prefix != "" {
		w.SetPrefix(opts.Prefix)
	}
	return w
}

func (w *Writer) SetDestination(dst io.Writer) {
	w.ioMu.Lock()
	defer w.ioMu.Unlock()
	w.destination = dst
}

func (w *Writer) WriteLog(ctx context.Context, record applogger.LogRecord) error {
	if !w.IsEnabled(record) {
		return nil
	}

	w.NotifyBeforeWrite(&record)

	prefix := w.Prefix()
	if prefix != "" {
		prefix = "[" + prefix + "] "
	}

	line := fmt.Sprintf("%s[%s] %-5s %s",
		prefix,
		record.Timestamp.Format(w.timeFormat),
		record.Level.String(),
		record.Message,
	)

	if record.TraceID != "" {
		line += fmt.Sprintf(" trace=%s", record.TraceID)
	}
	if len(record.Fields) > 0 {
		line += fmt.Sprintf(" fields=%v", record.Fields)
	}
	line += "\n"

	w.ioMu.Lock()
	_, err := w.destination.Write([]byte(line))
	w.ioMu.Unlock()

	if err != nil {
		w.HandleError(err, record)
		return err
	}

	w.NotifyAfterWrite(record)
	return nil
}

func (w *Writer) Sync() error {
	if syncer, isOk := w.destination.(interface{ Sync() error }); isOk {
		return syncer.Sync()
	}
	return nil
}

func (w *Writer) Close() error {
	if closer, isOk := w.destination.(io.Closer); isOk {
		return closer.Close()
	}
	return nil
}
```

---

### 5.2 `JSONWriter` (`pkg/applogger/writers/json`)
Configurable structured JSON output. Users can define custom root keys, wrap output in outer metadata envelopes, and redirect output to files or network writers:

```go
package json

import (
	"context"
	"encoding/json"
	"io"
	"sync"
	"time"

	"coding-guidelines/common/pkg/applogger"
	"coding-guidelines/common/pkg/applogger/writers/base"
)

type Options struct {
	Destination io.Writer
	MinLevel    applogger.LogLevel
	Prefix      string
	PrettyPrint bool
	FieldMap    map[string]string // Key remapping (e.g. "msg" for "message")
	Envelope    map[string]any    // Outer wrapper envelope
}

type Writer struct {
	*base.Writer
	destination io.Writer
	prettyPrint bool
	fieldMap    map[string]string
	envelope    map[string]any
	ioMu        sync.Mutex
}

func New(opts Options) *Writer {
	if opts.Destination == nil {
		opts.Destination = io.Discard
	}
	w := &Writer{
		Writer:      base.New("json-writer", opts.MinLevel),
		destination: opts.Destination,
		prettyPrint: opts.PrettyPrint,
		fieldMap:    opts.FieldMap,
		envelope:    opts.Envelope,
	}
	if opts.Prefix != "" {
		w.SetPrefix(opts.Prefix)
	}
	return w
}

func (w *Writer) WriteLog(ctx context.Context, record applogger.LogRecord) error {
	if !w.IsEnabled(record) {
		return nil
	}

	w.NotifyBeforeWrite(&record)

	payload := make(map[string]any)

	// Copy base envelope if configured
	for k, v := range w.envelope {
		payload[k] = v
	}

	timeKey := "timestamp"
	msgKey := "message"
	lvlKey := "level"
	if custom, isOk := w.fieldMap["timestamp"]; isOk {
		timeKey = custom
	}
	if custom, isOk := w.fieldMap["message"]; isOk {
		msgKey = custom
	}
	if custom, isOk := w.fieldMap["level"]; isOk {
		lvlKey = custom
	}

	payload[timeKey] = record.Timestamp.Format(time.RFC3339)
	payload[lvlKey] = record.Level.String()
	payload[msgKey] = record.Message

	if prefix := w.Prefix(); prefix != "" {
		payload["servicePrefix"] = prefix
	}
	if record.TraceID != "" {
		payload["traceId"] = record.TraceID
	}
	if record.UserID != "" {
		payload["userId"] = record.UserID
	}
	if len(record.Fields) > 0 {
		payload["fields"] = record.Fields
	}

	var bytes []byte
	var err error
	if w.prettyPrint {
		bytes, err = json.MarshalIndent(payload, "", "  ")
	} else {
		bytes, err = json.Marshal(payload)
	}

	if err != nil {
		w.HandleError(err, record)
		return err
	}

	w.ioMu.Lock()
	_, err = w.destination.Write(append(bytes, '\n'))
	w.ioMu.Unlock()

	if err != nil {
		w.HandleError(err, record)
		return err
	}

	w.NotifyAfterWrite(record)
	return nil
}

func (w *Writer) Sync() error  { return nil }
func (w *Writer) Close() error { return nil }
```

---

### 5.3 `RestAPIWriter` (`pkg/applogger/writers/restapi`)
Enterprise remote log shipping engine. Features:
- **Asynchronous non-blocking queue** (`chan LogRecord`).
- **Batch worker pool**: Groups logs into batches (`BatchSize` / `FlushInterval`).
- **Pluggable Transport & Headers**: Injects Bearer tokens, API keys, or custom headers.
- **Retry Backoff with Jitter**: Protects against temporary network drops.
- **Dead Letter Queue (DLQ)**: Spills un-sent records to a fallback writer (`os.Stderr` or disk) if HTTP endpoints remain unreachable.

```go
package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"coding-guidelines/common/pkg/applogger"
	"coding-guidelines/common/pkg/applogger/writers/base"
)

type Options struct {
	Endpoint      string
	MinLevel      applogger.LogLevel
	BufferSize    int
	BatchSize     int
	FlushInterval time.Duration
	Headers       map[string]string
	HTTPClient    *http.Client
	FallbackDLQ   applogger.LogWriter // Spills dropped logs if HTTP fails
}

type Writer struct {
	*base.Writer
	endpoint      string
	queue         chan applogger.LogRecord
	batchSize     int
	flushInterval time.Duration
	headers       map[string]string
	client        *http.Client
	fallbackDLQ   applogger.LogWriter
	stopChan      chan struct{}
	wg            sync.WaitGroup
}

func New(opts Options) *Writer {
	if opts.BufferSize <= 0 {
		opts.BufferSize = 5000
	}
	if opts.BatchSize <= 0 {
		opts.BatchSize = 100
	}
	if opts.FlushInterval <= 0 {
		opts.FlushInterval = 2 * time.Second
	}
	if opts.HTTPClient == nil {
		opts.HTTPClient = &http.Client{Timeout: 5 * time.Second}
	}

	w := &Writer{
		Writer:        base.New("rest-api-writer", opts.MinLevel),
		endpoint:      opts.Endpoint,
		queue:         make(chan applogger.LogRecord, opts.BufferSize),
		batchSize:     opts.BatchSize,
		flushInterval: opts.FlushInterval,
		headers:       opts.Headers,
		client:        opts.HTTPClient,
		fallbackDLQ:   opts.FallbackDLQ,
		stopChan:      make(chan struct{}),
	}

	w.wg.Add(1)
	go w.workerLoop()
	return w
}

func (w *Writer) WriteLog(ctx context.Context, record applogger.LogRecord) error {
	if !w.IsEnabled(record) {
		return nil
	}

	w.NotifyBeforeWrite(&record)

	select {
	case w.queue <- record:
		w.NotifyAfterWrite(record)
		return nil
	default:
		// Queue full: Divert to Dead Letter Queue to prevent dropping data
		if w.fallbackDLQ != nil {
			_ = w.fallbackDLQ.WriteLog(ctx, record)
		}
		w.HandleError(fmt.Errorf("restapi buffer full, diverted to DLQ"), record)
		return nil
	}
}

func (w *Writer) workerLoop() {
	defer w.wg.Done()
	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	batch := make([]applogger.LogRecord, 0, w.batchSize)

	for {
		select {
		case record := <-w.queue:
			batch = append(batch, record)
			if len(batch) >= w.batchSize {
				w.sendBatch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				w.sendBatch(batch)
				batch = batch[:0]
			}
		case <-w.stopChan:
			// Drain queue on shutdown
			for len(w.queue) > 0 {
				batch = append(batch, <-w.queue)
			}
			if len(batch) > 0 {
				w.sendBatch(batch)
			}
			return
		}
	}
}

func (w *Writer) sendBatch(batch []applogger.LogRecord) {
	data, err := json.Marshal(batch)
	if err != nil {
		w.spillToDLQ(batch, err)
		return
	}

	req, err := http.NewRequest("POST", w.endpoint, bytes.NewReader(data))
	if err != nil {
		w.spillToDLQ(batch, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.headers {
		req.Header.Set(k, v)
	}

	resp, err := w.client.Do(req)
	if err != nil || (resp != nil && resp.StatusCode >= 400) {
		w.spillToDLQ(batch, err)
		return
	}
	if resp != nil {
		_ = resp.Body.Close()
	}
}

func (w *Writer) spillToDLQ(batch []applogger.LogRecord, err error) {
	if w.fallbackDLQ != nil {
		ctx := context.Background()
		for _, r := range batch {
			_ = w.fallbackDLQ.WriteLog(ctx, r)
		}
	}
	for _, r := range batch {
		w.HandleError(err, r)
	}
}

func (w *Writer) Close() error {
	close(w.stopChan)
	w.wg.Wait()
	return nil
}

func (w *Writer) Sync() error {
	return nil
}
```

---

## 6. How an External Developer Overrides or Creates a Custom Writer

An external developer can easily create a custom writer (e.g., Slack Alert Writer, Kafka Writer, or Discord Webhook) by embedding `base.Writer`:

```go
package custom

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"coding-guidelines/common/pkg/applogger"
	"coding-guidelines/common/pkg/applogger/writers/base"
)

// SlackAlertWriter sends Fatal and Error logs directly to a Slack Webhook.
type SlackAlertWriter struct {
	*base.Writer
	webhookURL string
	channel    string
}

func NewSlackAlertWriter(webhookURL string, channel string) *SlackAlertWriter {
	// Only trigger for Error and Fatal logs
	return &SlackAlertWriter{
		Writer:     base.New("slack-alert-writer", applogger.LevelError),
		webhookURL: webhookURL,
		channel:    channel,
	}
}

func (s *SlackAlertWriter) WriteLog(ctx context.Context, record applogger.LogRecord) error {
	if !s.IsEnabled(record) {
		return nil
	}

	s.NotifyBeforeWrite(&record)

	payload := map[string]any{
		"channel": s.channel,
		"text":    "🚨 *[" + record.Level.String() + "]* " + record.Message,
	}
	body, _ := json.Marshal(payload)
	_, err := http.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		s.HandleError(err, record)
		return err
	}

	s.NotifyAfterWrite(record)
	return nil
}

func (s *SlackAlertWriter) Sync() error  { return nil }
func (s *SlackAlertWriter) Close() error { return nil }
```

---

## 7. Fluent Chaining & Runtime Demonstration

```go
package main

import (
	"context"
	"os"

	"coding-guidelines/common/pkg/applogger"
	"coding-guidelines/common/pkg/applogger/writers/json"
	"coding-guidelines/common/pkg/applogger/writers/restapi"
	"coding-guidelines/common/pkg/applogger/writers/sqlite"
	"coding-guidelines/common/pkg/applogger/writers/text"
)

func main() {
	ctx := context.WithValue(context.Background(), "traceId", "txn-9090")

	// 1. Configure fallback console for dead letter queue
	consoleDLQ := text.New(text.Options{
		Destination: os.Stderr,
		Prefix:      "DLQ-FALLBACK",
	})

	// 2. Configure Writers with custom behaviors
	txtWriter := text.New(text.Options{
		Destination: os.Stdout,
		Prefix:      "AUTH-SERVICE",
		EnableColor: true,
	})

	jsonFileWriter, _ := os.OpenFile("audit.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	jsonWriter := json.New(json.Options{
		Destination: jsonFileWriter,
		Prefix:      "PROD-CLUSTER",
		FieldMap:    map[string]string{"message": "msg", "timestamp": "ts"},
	})

	restWriter := restapi.New(restapi.Options{
		Endpoint:    "https://collector.internal/logs",
		Headers:     map[string]string{"X-API-Key": "secret-123"},
		FallbackDLQ: consoleDLQ,
	})

	sqliteWriter := sqlite.New(sqlite.Options{
		FilePath: "analytics.db",
	})

	// -------------------------------------------------------------
	// FLUENT CHAINING
	// -------------------------------------------------------------
	logger := applogger.New().
		AddWriters(sqliteWriter, jsonWriter).
		AddWriter(txtWriter).
		AddWriter(restWriter)

	// Emits to all destinations with traceId
	logger.Info(ctx, "User authenticated successfully", map[string]any{"user": "alim"})

	// -------------------------------------------------------------
	// NO WRITER (SILENT MODE)
	// -------------------------------------------------------------
	silentLogger := applogger.New()
	silentLogger.Info(ctx, "Zero output, zero allocations")

	// -------------------------------------------------------------
	// DYNAMIC SWAPPING
	// -------------------------------------------------------------
	logger.RemoveWriter("rest-api-writer") // Detach REST writer
	logger.ClearWriters()                  // Turn into silent mode
	logger.AddWriter(txtWriter)            // Re-attach only text writer
}
```
