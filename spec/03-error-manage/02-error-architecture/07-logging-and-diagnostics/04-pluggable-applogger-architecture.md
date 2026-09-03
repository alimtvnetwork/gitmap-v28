# Specification: `applogger` Pluggable Logger Architecture

**Version:** 3.3.0  
**Status:** Draft Specification (Pending Review)  
**Package:** `04-code/golang/pkg/applogger`  
**Reference Implementations:** `D:\work\03-aukgo\core\coreinterface\loggerinf\` (`Logger`, `StandardLogger`, `BasePersistentLogger`)

---

## 1. Executive Summary & Design Goals

Logging requirements differ across development, test, and production environments:
- **Local Development:** Colored terminal stdout logs with instant stack traces.
- **Enterprise Services:** High-throughput JSON logs delegated to **Uber Zap** (`go.uber.org/zap`).
- **CLI & Embedded Agents:** Local file system logs (rotating `.log` files) or local **SQLite database** audit trails (`logs` table).

The **`applogger`** package defines a clean, decoupled logger interface that enables swapping logging backends (Console, File, SQLite, Uber Zap, or Composite) through configuration without changing application code.

---

## 2. Core Interface Contracts

```go
package applogger

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/appfaults"
)

// Logger is the unified logging contract.
type Logger interface {
	// Level-based logging (PascalCase names: Debug, Info, Warn, Error, Fatal)
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)

	// Formatted logging
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)

	// Structured Error & Fault Collection Logging
	LogError(err *appfault.AppError)
	LogFaults(faults *appfaults.Collection)

	// Context & Field Enrichment
	WithContext(key string, val any) Logger
	WithFields(fields map[string]any) Logger

	// Sink Lifecycle
	Sync() error
	Close() error
}
```

---

## 3. Log Sink / Driver Architecture

The logger delegates log dispatching to one or more pluggable **Sinks** (`LogSink`):

```go
// LogEntry represents a structured log event payload.
type LogEntry struct {
	Timestamp string              `json:"Timestamp" yaml:"Timestamp"`
	Level     LogLevel            `json:"Level" yaml:"Level"`
	Message   string              `json:"Message" yaml:"Message"`
	Fields    appfault.ContextMap `json:"Fields,omitempty" yaml:"Fields,omitempty"`
	Caller    string              `json:"Caller,omitempty" yaml:"Caller,omitempty"`
	Stack     string              `json:"Stack,omitempty" yaml:"Stack,omitempty"`
}

// LogSink defines the driver interface for persistent destinations.
type LogSink interface {
	WriteEntry(entry LogEntry) error
	Sync() error
	Close() error
}
```

---

## 4. Built-in Sinks & Adapters

```
┌─────────────────────────────────────────────────────────────┐
│                      applogger.Logger                       │
└──────────────────────────────┬──────────────────────────────┘
                               │
            ┌──────────────────┼──────────────────┐
            ▼                  ▼                  ▼
   ┌─────────────────┐ ┌───────────────┐ ┌─────────────────┐
   │   ConsoleSink   │ │   FileSink    │ │   SQLiteSink    │
   │ (Stdout/JSON)   │ │  (Rotated)    │ │ (.db / table)   │
   └─────────────────┘ └───────────────┘ └─────────────────┘
            │                                     │
            ▼                                     ▼
   ┌─────────────────┐                   ┌─────────────────┐
   │   ZapAdapter    │                   │  CompositeSink  │
   │  (Uber Zap API) │                   │(Multi-Broadcast)│
   └─────────────────┘                   └─────────────────┘
```

### 4.1 ConsoleSink (Terminal / Stdout)
- Prints human-readable colored output during CLI operations or structured JSON when `UseJSON = true`.

### 4.2 FileSink (Filesystem Log Files)
- Writes log entries sequentially to a specified log file path (e.g. `tmp/logs/app.log`).
- Supports automatic size-based rotation and header flushing.

### 4.3 SQLiteSink (SQLite Database Audit Trail)
- Writes structured entries to a local or embedded SQLite database table (`app_logs`):
  ```sql
  CREATE TABLE IF NOT EXISTS app_logs (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      timestamp TEXT NOT NULL,
      level TEXT NOT NULL,
      message TEXT NOT NULL,
      caller TEXT,
      fields_json TEXT,
      stack_trace TEXT
  );
  ```

### 4.4 ZapAdapter (Uber Zap Integration)
- Adapts any existing `*zap.Logger` or `*zap.SugaredLogger` instance so enterprise teams using Zap can use `applogger.Logger` without rewriting their logging pipeline.

### 4.5 CompositeSink (Multi-Destination Broadcaster)
- Fan-outs every log entry to multiple configured sinks simultaneously (e.g., Console for dev output + SQLite for auditing + File for backup).

---

## 5. Logger Factory & Configuration

```go
// DriverType selects the backend implementation.
type DriverType byte

const (
	DriverConsole DriverType = iota
	DriverFile
	DriverSQLite
	DriverZap
	DriverComposite
)

// Config configures the logger instance.
type Config struct {
	MinLevel      LogLevel
	Driver        DriverType
	FilePath      string            // Required for DriverFile
	SQLitePath    string            // Required for DriverSQLite (e.g. "tmp/logs.db")
	ZapLogger     any               // Optional *zap.Logger for DriverZap
	Sinks         []LogSink         // For DriverComposite
	IsStackTrace  bool
	UseJSON       bool
}

// New constructs a Logger using the requested configuration and sink driver.
func New(cfg Config) (Logger, error)

// Default returns a standard Console logger at Info level.
func Default() Logger
```

---

## 6. Integration with `appfault` and `appfaults`

1. **Logging Single AppError:**
   ```go
   log := applogger.Default()
   err := appfault.New(errtype.Database, "connection timeout").WithOp("db.connect")
   log.LogError(err)
   ```
2. **Logging Error Collections:**
   ```go
   faults := appfaults.New()
   faults.Add(err1).Add(err2)
   log.LogFaults(faults)
   ```
3. **Fluent Fault Dispatch:**
   ```go
   // AppFaults can also dispatch directly to a Logger
   faults.LogTo(logger)
   ```
