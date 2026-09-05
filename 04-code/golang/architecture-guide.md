# Golang Architecture Guide: Universal Writers, Result Wrappers, Enum-Driven File I/O, & Logging Engine

## 1. Executive Summary & Design Principles

This document provides the authoritative technical reference for the enterprise Go components implemented in `04-code/golang/pkg/`. The ecosystem is built on four core pillars:

1. **Structured Failure Standard (`*appfault.AppError` & `result.Wrap[T]`):**
   - Strict standard error return type: `*appfault.AppError`.
   - Error classification is purely driven by the `errtype.Variation` enum (`uint16`), with **NO** redundant `ErrID` or `ErrCode` string fields.
   - Caller tracking is encapsulated in a dedicated structured object: `CallerInfo` (containing `Function`, `File`, `Line`), **NEVER** a raw string.
   - Universal monadic container `result.Wrap[T]` (alias for `appfault.Result[T]`) eliminates package-stutter (`result.Wrap[T]` and `appfault.AppError`).
2. **Deterministic 16-Bit / UTF-16 Enum Taxonomy:**
   - Error types are represented as a 16-bit unsigned integer enum: `type Variation uint16` (`errtype.None = 0`, `errtype.Validation = 2`, `errtype.NotFound = 3`, etc.).
   - File permissions: `fileutil.FilePermType` (`uint32` POSIX bitmask).
   - File open modes: `fileutil.FileOpenModeType` (`byte`).
3. **Dedicated Printing, Display, & Formatting Subsystems:**
   - Both the error (`*appfault.AppError`) and the result wrapper (`result.Wrap[T]`) provide default printing and formatting sections (`.Print()`, `.PrintWith()`, `.Format()`, `.DisplayError()`, `.FullString()`, `.ToClipboard()`) that can be customized at runtime with a `FaultFormatter` or `ResultFormatter[T]`.
4. **Pluggable & Swappable Write Pipelines (`pkg/streamwriter` & `pkg/appwriter`):**
   - Universal `Writer[T]` and `Streamer[T]` interfaces satisfying `sync.Locker`.
   - Runtime swappable formatters (`SetFormatMethod`), write handlers (`SetWriteMethod`), destinations (`SetDestination`), and streamers (`SetStreamer`).
   - Standard writers: Filesystem (append/truncate with permissions), JSON serializer, REST API client, Console/Text.
   - Fluent composite logger (`streamwriter.Logger[T]`) with zero-allocation silent mode.

---

## 2. Error Taxonomy & Standard Return Types (`pkg/errtype`, `pkg/appfault`)

### 2.1 The Error Variation Enum: `type Variation uint16` (`pkg/errtype`)

Error classification is **NOT** a string. It is an extensible 16-bit unsigned integer enum (`uint16`) representing UTF-16 code points:

```go
package errtype

// Variation represents an extensible 16-bit unsigned integer error type code.
type Variation uint16

const (
    None         Variation = 0   // No error occurred (successful state)
    NoError      Variation = None
    Generic      Variation = 1   // Unspecified standard error
    Validation   Variation = 2   // Input, schema, or invariant validation failure (HTTP 400)
    NotFound     Variation = 3   // Resource or record not found (HTTP 404)
    Precondition Variation = 4   // State prerequisites unsatisfied (HTTP 412)
    Execution    Variation = 5   // General runtime execution failure
    Database     Variation = 6   // Database query, execution, or connection failure
    Network      Variation = 7   // Network transport or connectivity failure
    Timeout      Variation = 8   // Operation exceeded deadline (HTTP 504)
    IO           Variation = 9   // Disk I/O or filesystem failure (HTTP 500)
    Unauthorized Variation = 10  // Authentication required (HTTP 401)
    Forbidden    Variation = 11  // Authenticated but unauthorized (HTTP 403)
    Internal     Variation = 12  // Unexpected internal server fault (HTTP 500)
    Unknown      Variation = 13  // Unclassified error state
)
```

#### Variation Built-in Methods
`errtype.Variation` provides methods mapping the 16-bit enum to names, codes, and HTTP statuses:

```go
func (v Variation) Code() uint16          // Returns raw integer (e.g. 2 for Validation)
func (v Variation) Name() string          // Returns PascalCase name (e.g. "Validation")
func (v Variation) Description() string   // Returns human-readable explanation
func (v Variation) HttpStatus() int       // Maps to HTTP code (e.g. 400 for Validation)
func (v Variation) HasError() bool        // Returns true if v != None
func (v Variation) IsNone() bool          // Returns true if v == None
```

### 2.2 Standard Error Carrier: `*appfault.AppError` (`pkg/appfault`)

`AppError` encapsulates failure diagnostics. All internal fields are unexported for immutability and encapsulation:

```go
package appfault

type AppError struct {
    errType    errtype.Variation // 16-bit error type enum (NO string ErrorCode or ErrID!)
    message    string            // Human-readable diagnostic message
    caller     CallerInfo        // Value-based caller site object (zero heap allocation, NO pointer!)
    stack      StackTrace        // Structured slice of StackFrame objects
    ctx        ContextMap        // Thread-safe structured metadata map
    cause      error             // Underlying root cause error
    statusCode int               // Explicit HTTP status override (or derived from errType)
}
```

### 2.3 Value-Based Caller Object: `CallerInfo`

The caller site is **never** a loose string or pointer indirection. It is a strongly-typed value object (`CallerInfo`) embedded directly into `AppError` without heap allocation:

```go
package appfault

type CallerInfo struct {
    Function string `json:"Function,omitempty" yaml:"Function,omitempty"`
    File     string `json:"File,omitempty" yaml:"File,omitempty"`
    Line     int    `json:"Line,omitempty" yaml:"Line,omitempty"`
}

// String formats as "file:line (function)" or "file:line"
func (c CallerInfo) String() string

// IsEmpty returns true if caller metadata is unset
func (c CallerInfo) IsEmpty() bool
```

`CaptureCallerInfo(skip int)` automatically populates `CallerInfo` by value using `runtime.Caller` when an error is constructed.

### 2.4 Strict Immutability & The AppBuilder Pattern

`*appfault.AppError` is **strictly immutable**. Once constructed, an `AppError` cannot be mutated in place, guaranteeing complete thread safety and preventing race conditions during concurrent logging or telemetry propagation.

#### 1. Copy-on-Write Immutability
All `With*` methods on `*AppError` return a fresh cloned instance without modifying the receiver:

```go
baseErr := appfault.New(errtype.Validation, "validation failed")

// Each derivation produces a new, independent immutable instance:
err400 := baseErr.WithStatusCode(400)
err422 := baseErr.WithStatusCode(422)

// baseErr remains 100% unmodified!
```

#### 2. Mutable Staging via `AppBuilder` (`AppErrorBuilder`)
For assembling errors with multiple diagnostic attributes before freezing, use `AppBuilder`:

```go
// 1. Accumulate diagnostic attributes mutably in staging
builder := appfault.NewAppBuilder(errtype.Validation, "input invalid").
    WithStatusCode(400).
    WithCaller(appfault.CallerInfo{File: "auth/service.go", Line: 55, Function: "HandleRegistration"}).
    WithContext("traceId", "tx-9901")

// 2. Freeze into 100% immutable *AppError
appErr := builder.Build()

// 3. Round-trip serialization & staging modification
jsonBytes, _ := builder.MarshalJSON()

var restoredBuilder appfault.AppBuilder
_ = restoredBuilder.UnmarshalJSON(jsonBytes)
restoredErr := restoredBuilder.Build()

// 4. Staging modifications on an existing immutable error via ToBuilder()
updatedErr := restoredErr.ToBuilder().
    SetStatusCode(422).
    Build()
```

### 2.5 Error Merging & Multi-Loop Stack Trace Tracking

When retrying operations or iterating through loops, multiple errors can occur sequentially. `appfault.Merge(prev, next)` combines two errors into one immutable `*AppError` while preserving complete diagnostic lineage:

1. **Existence Validation:** Checks if errors actually exist (`HasError()`). If either is nil or none, the active error is returned.
2. **First Error Stack Trace Preservation:** Automatically stores the very first failure's stack trace under context key `"FirstErrorStackTrace"`, caller under `"FirstErrorCaller"`, and message under `"FirstErrorMessage"`.
3. **Prior Error Tracking:** Records the immediately preceding error's stack trace under `"PriorStackTrace"` and caller under `"PriorCaller"`.
4. **Loop & Attempt Counter:** Automatically manages `"LoopCount"`, incrementing each time an error is merged.
5. **Full History Accumulation:** Appends each attempt into `"StackTraceHistory"` (`[]string`), allowing engineers to trace exactly where and how failures occurred across every loop iteration.

```go
var accumulatedErr *appfault.AppError

// Retrying an operation across a loop:
for attempt := 1; attempt <= 3; attempt++ {
    currentErr := executeNetworkCall()
    if currentErr.HasError() {
        // Merges prior failure into current failure, accumulating stack traces
        accumulatedErr = appfault.Merge(accumulatedErr, currentErr)
    }
}

// Inspecting accumulated diagnostics:
fmt.Println("Total failure attempts:", accumulatedErr.LoopCount())
fmt.Println("Initial failure stack trace:\n", accumulatedErr.FirstErrorStackTrace())
fmt.Println("All attempt traces:\n", accumulatedErr.StackTraceHistory())
```

---

## 3. Dedicated Error Display & Printing Subsystem (`pkg/appfault`)

`AppError` provides built-in methods for clean console display, customizable formatting, and AI/clipboard reporting.

### 3.1 Display & Printing Methods on `*AppError`

```go
// 1. Default Print to standard output
err.Print()
// Output: ❌ [Validation:2] username is required (at user.go:45 (CreateUser))

// 2. Custom Formatting with FaultFormatter
err.PrintWith(func(e *appfault.AppError) string {
    return fmt.Sprintf("🚨 ALERT [%s] -> %s", e.Type().Name(), e.Message())
})

// 3. String Formatting
formattedLine := err.Format(appfault.DefaultFaultFormatter)

// 4. Comprehensive Full Diagnostic Dump
diagnosticReport := err.FullString()
// Output includes: ERROR line, CALLER site, CAUSE, CONTEXT map, and STACK TRACE

// 5. Markdown Report for AI Analysis / Issue Tracking
clipboardReport := err.ToClipboard()
// Output: Markdown header, bulleted metadata, and fenced stack trace codeblock

// 6. Direct Terminal Banner
err.DisplayError()
```

### 3.2 Multi-Destination Presentation Subsystem

`*appfault.AppError` and `result.Wrap[T]` provide 3 first-class presentation formats out of the box:

```go
// 1. Terminal / Stdout Banner (with status icon, caller info, and HTTP status)
stdoutBanner := err.FormatStdout()
err.PrintStdout()
// Output:
// ❌ ERROR [Validation:2] username is required (HTTP 400)
//    Caller:  services/user/validator.go:42 (ValidateUser)
//    Cause:   string length is 0
//    Context: email="bad@"

// 2. Structured JSON (RFC-compatible machine-readable payload)
jsonString := err.FormatJson()
err.PrintJson()
// Output:
// {
//   "Type": 2,
//   "Message": "username is required",
//   "Caller": { "Function": "ValidateUser", "File": "services/user/validator.go", "Line": 42 },
//   "StatusCode": 400
// }

// 3. Structured Text Log (Single-line log aggregator format for Loki/ELK)
logLine := err.FormatTextLog()
err.PrintLog()
// Output:
// [ERROR] [Validation:2] status=400 caller="services/user/validator.go:42" msg="username is required"
```

### 3.3 Custom Formatter Function Signature
```go
type FaultFormatter func(e *AppError) string
```

Users can define custom formatters for custom syslog formats, Slack notifications, or telemetry pipelines.

---

## 4. Result & Wrap Envelopes (`pkg/result`, `pkg/appwriter`)

### 4.1 Generic Container: `result.Wrap[T]`

`result.Wrap[T]` (an alias for `appfault.Result[T]`) replaces `(T, error)` tuples:

```go
package appfault

type Result[T any] struct {
    Value    T         `json:"Value,omitempty" yaml:"Value,omitempty"`
    AppError *AppError `json:"AppError,omitempty" yaml:"AppError,omitempty"`
}
```

```go
package result

type Wrap[T any] = appfault.Result[T]
```

#### Inspection Methods
```go
wrap.IsSuccess()     // bool: true if AppError == nil
wrap.IsFailed()      // bool: true if AppError != nil
wrap.Data()          // returns Value (type T)
wrap.Value()         // returns Value (type T)
wrap.Fault()         // returns *appfault.AppError
wrap.Unwrap()        // returns (T, *appfault.AppError)
wrap.UnwrapOr(def)   // returns T if success, or def if failed
```

### 4.2 Result Display & Printing Subsystem

`Result[T]` / `Wrap[T]` also has built-in printing and formatting methods:

```go
// 1. Default Print
res.Print()
// If success: ✅ [OK] <value>
// If failed:  ❌ [Validation:2] <message> (at caller)

// 2. Print fault only if failed
res.PrintFault()

// 3. Format with custom ResultFormatter[T]
formatted := res.Format(func(r result.Wrap[MyData]) string {
    if r.IsFailed() {
        return "FAILED: " + r.Fault().Message()
    }
    return fmt.Sprintf("SUCCESS: %+v", r.Data())
})
```

### 4.3 Wrap Constructors & Error Propagation

| Constructor | Description |
| :--- | :--- |
| `result.WrapSuccess[T](data)` | Creates a successful `Wrap[T]` containing data. |
| `result.WrapFailure[T](err)` | Wraps an existing `*appfault.AppError` into a failed `Wrap[T]`. |
| `result.WrapFailureFromError[T](err)` | Creates a failed `Wrap[T]` directly from `*appfault.AppError`. |
| `result.WrapFailureWithId[T](variation, msg)` | Constructs an `AppError` from `errtype.Variation` and returns failed `Wrap[T]`. |
| `result.WrapFailureWithCause[T](variation, cause, msg)` | Constructs an `AppError` with root cause and error enum. |
| `result.WrapFailureFromWrap[T, U](failedWrap)` | Propagates failure across different types (`Wrap[U]` to `Wrap[T]`). |

```go
// Clean propagation across distinct payload types:
func FetchUser(id string) result.Wrap[*User] {
    fileRes := fileutil.ReadAll("users/" + id + ".json")
    if fileRes.IsFailed() {
        // Propagate failure from Wrap[[]byte] to Wrap[*User]
        return result.WrapFailureFromWrap[*User](fileRes)
    }

    var u User
    if err := json.Unmarshal(fileRes.Data(), &u); err != nil {
        return result.WrapFailureWithCause[*User](errtype.Validation, err, "malformed JSON")
    }

    return result.WrapSuccess(&u)
}
```

### 4.4 The Writer Wrap Pattern: `BaseWriterWrap` (`pkg/appwriter`)

`appwriter.BaseWriterWrap` aliases `result.Wrap[*BaseWriter]`:

```go
package appwriter

type BaseWriterWrap = result.Wrap[*BaseWriter]

// Dedicated constructors defined directly inside the appwriter wrapper definition:
func WrapWriterSuccess(w *BaseWriter) BaseWriterWrap
func WrapWriterFailure(err *appfault.AppError) BaseWriterWrap
func WrapWriterFailureFromError(err *appfault.AppError) BaseWriterWrap
func WrapWriterFailureWithId(errType errtype.Variation, msg string) BaseWriterWrap
func WrapWriterFailureWithCause(errType errtype.Variation, cause error, msg string) BaseWriterWrap
func WrapWriterFailureFromWrap[U any](failed result.Wrap[U]) BaseWriterWrap
```

---

## 5. Enum-Driven File I/O Engine (`pkg/fileutil`)

All file reading and writing utilities rely on strongly-typed enums instead of raw integer flags or octal literals.

### 5.1 Permission Bitmask Enum: `FilePermType` (`uint32`)

```go
type FilePermType uint32

const (
    FilePermStandard        FilePermType = 0644  // rw-r--r-- (Standard files)
    FilePermPrivate         FilePermType = 0600  // rw------- (Private keys, secrets)
    FilePermExecutable      FilePermType = 0755  // rwxr-xr-x (Scripts, binaries, dirs)
    FilePermReadOnly        FilePermType = 0444  // r--r--r-- (Protected read-only)
    FilePermOwnerExec       FilePermType = 0700  // rwx------ (Owner-only dir)
    FilePermGroupReadWrite  FilePermType = 0660  // rw-rw---- (Group collaboration)
    FilePermPublicAll       FilePermType = 0777  // rwxrwxrwx (Full public access)
)

// Mode converts the enum to standard os.FileMode
mode := FilePermStandard.Mode() // os.FileMode(0644)
```

### 5.2 Open Mode Enum: `FileOpenModeType` (`byte`)

```go
type FileOpenModeType byte

const (
    FileOpenReadOnly       FileOpenModeType = iota // os.O_RDONLY
    FileOpenWriteOnly                              // os.O_WRONLY
    FileOpenReadWrite                              // os.O_RDWR
    FileOpenAppend                                 // os.O_WRONLY | os.O_APPEND
    FileOpenCreateAppend                           // os.O_CREATE | os.O_WRONLY | os.O_APPEND
    FileOpenCreateTruncate                         // os.O_CREATE | os.O_WRONLY | os.O_TRUNC
    FileOpenCreateNew                              // os.O_CREATE | os.O_EXCL | os.O_WRONLY
)

// Flags maps the enum to os.OpenFile integer flags
flags := FileOpenCreateAppend.Flags()
```

### 5.3 Utility Functions

```go
// Open file with mode and permission enums (automatically creates parent directory if create mode)
fileWrap := fileutil.OpenFile("logs/app.log", fileutil.FileOpenCreateAppend, fileutil.FilePermStandard)
if fileWrap.IsFailed() {
    fileWrap.Print() // Nicely prints formatted fault
    return
}
f := fileWrap.Data()
defer f.Close()

// Read entire file content into byte slice
dataWrap := fileutil.ReadAll("config.json")

// Read entire file content as string
strWrap := fileutil.ReadString("config.json")

// Write entire file replacing content
writeWrap := fileutil.WriteFile("build/status.txt", []byte("SUCCESS"), fileutil.FilePermStandard)

// Append string to file
appendWrap := fileutil.AppendString("logs/audit.log", "Event occurred\n", fileutil.FilePermStandard)
```

### 5.4 Advanced File Utilities (Atomic Writes & Chunked Streaming)

```go
// 1. Atomic Write (Writes to temp file, syncs to physical disk, and atomically swaps via os.Rename)
// Prevents partial writes, corruption, or zero-byte files on crash/power loss
atomWrap := fileutil.WriteAtomic("config/active.json", []byte(`{"version":2}`), fileutil.FilePermStandard)

// 2. Chunked File Streaming (Reads fixed-size chunks using memory pool sync.Pool)
// Prevents high memory consumption on large multi-gigabyte files
readWrap := fileutil.ReadChunked("data/large_dataset.bin", 64*1024, func(chunk []byte) *appfault.AppError {
    // Process chunk without buffering entire file in memory
    return nil
})

// 3. Chunked Streaming Writer (Streams from io.Reader to disk with sync)
writeChunkWrap := fileutil.WriteChunked("data/downloaded.tar.gz", fileutil.FilePermStandard, resp.Body, 64*1024)

// 4. File Writer Bridge (Attaches streamwriter.PluggableWriter directly to fileutil.OpenFile)
writerWrap := fileutil.NewFileWriter("logs/pipeline.log", fileutil.FileOpenCreateAppend, fileutil.FilePermStandard)
if writerWrap.IsSuccess() {
    writer := writerWrap.Data()
    _ = writer.Write(ctx, "pipeline initialized")
}
```

---

## 6. JSON Result Wrapper & Data Flow (`pkg/streamwriter`)

`streamwriter.JsonResult` provides a structured envelope around JSON data, integrating with `*appfault.AppError`.

### 6.1 Multi-Source Creation via `JsonSource`
```go
// From struct or map
res1 := streamwriter.JsonSource.FromPayload(map[string]any{"status": "active", "count": 10})

// From raw bytes with JSON validation
res2 := streamwriter.JsonSource.FromBytes([]byte(`{"metric":"cpu","value":42}`))

// From string
res3 := streamwriter.JsonSource.FromString(`{"enabled":true}`)

// From io.Reader
res4 := streamwriter.JsonSource.FromReader(resp.Body)
```

### 6.2 Transformation & Formatting
```go
prettyJSON  := res1.Pretty()   // Indented with 2 spaces
compactJSON := res1.Compact()  // Minified single-line string
rawBytes    := res1.Raw()      // Underlying byte slice
bytesObj    := res1.ToBytes()  // Converts to strongly-typed Bytes[T] envelope

// Safe unmarshal with *appfault.AppError return
var target Config
if appErr := res1.Unmarshal(&target); appErr != nil {
    appErr.Print()
}
```

### 6.3 Deterministic Object Compiler (`streamwriter.Compile`)
Ensures reproducible string serialization with deterministically sorted keys across nested maps, slices, and structs:

```go
out := streamwriter.Compile(map[string]any{"z": 1, "a": 2, "m": 3})
// Always formatted as: {a: 2, m: 3, z: 1}
```

---

## 7. Pluggable & Swappable Writers (`pkg/streamwriter`)

### 7.1 Writer Interfaces & Locker Synchronization

```go
type Writer[T any] interface {
    Name() string
    Write(ctx context.Context, payload T) *appfault.AppError
    AsWriter() Writer[T]
    Lock()
    Unlock()
    Sync() *appfault.AppError
    Close() *appfault.AppError
}

type Streamer[T any] interface {
    Name() string
    Stream(ctx context.Context, payload T) *appfault.AppError
    AsStreamer() Streamer[T]
    AsWriter() Writer[T]
    IsLocked() bool
    Lock()
    Unlock()
    Destination() io.Writer
    Sync() *appfault.AppError
    Close() *appfault.AppError
}
```

### 7.2 Locked vs. Lockless Streamers
- **`LockedStreamer[T]`**: Wraps destination `io.Writer` with an internal re-entrant mutex. Thread-safe for concurrent writes across goroutines.
- **`LocklessStreamer[T]`**: Direct un-synchronized streaming for memory buffers (`bytes.Buffer`) or single-threaded workers.

### 7.3 `PluggableWriter[T]` & Hot-Swapping

```go
writer := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
    Name:        "app-writer",
    Destination: os.Stdout,
})
```

#### 1. Swapping Formatters (`SetFormatMethod`)
Transforms generic payload `T` into `Bytes[T]`:
```go
writer.SetFormatMethod(func(payload any) streamwriter.Bytes[any] {
    line := fmt.Sprintf(">>> %v <<<\n", payload)
    return streamwriter.NewBytes([]byte(line), payload)
})
```

#### 2. Swapping Write Functions (`SetWriteMethod`)
Receives the attached `streamer Streamer[T]`, `context.Context`, writer instance pointer `*PluggableWriter[T]`, and payload `T`:
```go
writer.SetWriteMethod(func(s streamwriter.Streamer[any], ctx context.Context, w *streamwriter.PluggableWriter[any], payload any) *appfault.AppError {
    dest := w.Destination()
    if dest == nil {
        return appfault.New(errtype.IO, "destination writer is nil")
    }

    record := fmt.Sprintf("[%s] %s\n", w.Name(), streamwriter.Compile(payload))
    _, err := dest.Write([]byte(record))
    if err != nil {
        return appfault.Wrap(errtype.IO, err, "write failed")
    }
    return nil
})
```

#### 3. Swapping Destination (`SetDestination`) & Streamer (`SetStreamer`)
```go
// Switch target from stdout to disk file at runtime
fileWrap := fileutil.OpenFile("logs/runtime.log", fileutil.FileOpenCreateAppend, fileutil.FilePermStandard)
if fileWrap.IsSuccess() {
    writer.SetDestination(fileWrap.Data())
}
```

### 7.4 Payload Intelligence & `AnyWriter` (`pkg/streamwriter`)

For heterogeneous write pipelines where events, strings, errors, and raw bytes pass through the same writer, use `AnyWriter`:

```go
// Non-generic alias for PluggableWriter[any]
type AnyWriter = PluggableWriter[any]

// Constructor with smart payload inspection
writer := streamwriter.NewAnyWriter(streamwriter.WriterOptions[any]{
    Name:        "universal-logger",
    Destination: os.Stdout,
})
```

#### Reflect Converter & Byte Dispatch:
In accordance with core engine guidelines (`aukgo/core`), `ExtractBytes(payload)` and `InspectPayload(payload)` enforce strict serialization rules:
1. **Raw `[]byte` Preservation:** If payload is already `[]byte`, it is streamed directly to the destination. Calling `json.Marshal([]byte)` is strictly avoided because it transforms raw bytes into a **Base64-encoded string**.
2. **String Bypass:** Plain strings are streamed directly as UTF-8 bytes (`[]byte(s)`), eliminating JSON quotation overhead.
3. **Structured Errors:** `*appfault.AppError` objects are extracted via their structured methods rather than generic reflection.
4. **Maps & Structs:** Serialized into deterministic JSON or compiled strings.

---

## 8. Default Destination Writers & Integration

### 8.1 Filesystem Disk Writer
```go
fileWrap := fileutil.OpenFile("logs/audit.log", fileutil.FileOpenCreateAppend, fileutil.FilePermStandard)
diskStreamer := streamwriter.NewLockedStreamer[any](streamwriter.StreamerOptions[any]{
    Name:        "disk-streamer",
    Destination: fileWrap.Data(),
})
diskWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
    Name:     "disk-writer",
    Streamer: diskStreamer,
})
```

### 8.2 JSON Serializing Writer
```go
jsonWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
    Name:        "json-writer",
    Destination: os.Stdout,
    WriteMethod: func(s streamwriter.Streamer[any], ctx context.Context, w *streamwriter.PluggableWriter[any], payload any) *appfault.AppError {
        res := streamwriter.JsonSource.FromPayload(payload)
        if res.HasError() {
            return res.Fault()
        }
        _, err := w.Destination().Write([]byte(res.Compact() + "\n"))
        if err != nil {
            return appfault.Wrap(errtype.IO, err, "json write error")
        }
        return nil
    },
})
```

### 8.3 REST API / HTTP Remote Sink
```go
apiWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
    Name: "http-api-writer",
    WriteMethod: func(s streamwriter.Streamer[any], ctx context.Context, w *streamwriter.PluggableWriter[any], payload any) *appfault.AppError {
        jsonRes := streamwriter.JsonSource.FromPayload(payload)
        if jsonRes.HasError() {
            return jsonRes.Fault()
        }

        req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.internal/events", bytes.NewReader(jsonRes.Raw()))
        if err != nil {
            return appfault.Wrap(errtype.Internal, err, "failed to construct request")
        }
        req.Header.Set("Content-Type", "application/json")

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            return appfault.Wrap(errtype.IO, err, "HTTP request failed")
        }
        defer resp.Body.Close()

        if resp.StatusCode >= 400 {
            return appfault.New(errtype.IO, fmt.Sprintf("endpoint returned HTTP %d", resp.StatusCode))
        }
        return nil
    },
})
```

### 8.4 Console / Terminal Writer
```go
consoleWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
    Name:        "console-writer",
    Destination: os.Stdout,
    FormatMethod: func(payload any) streamwriter.Bytes[any] {
        timestamp := time.Now().Format("15:04:05.000")
        line := fmt.Sprintf("[%s] %s\n", timestamp, streamwriter.Compile(payload))
        return streamwriter.NewBytes([]byte(line), payload)
    },
})
```

---

## 9. Fluent Composite Logging Engine (`streamwriter.Logger[T]`)

### 9.1 Multi-Writer Chaining & Fan-Out
```go
logger := streamwriter.NewLogger[any]().
    AddWriter(consoleWriter).
    AddWriter(diskWriter).
    AddWriter(apiWriter)

count := logger.WriterCount() // Returns 3
```

### 9.2 Zero-Allocation Silent Mode
When `logger.WriterCount() == 0`, calls to `logger.Emit(...)` return immediately with `nil`, generating 0 allocations during testing or disabled logging states.

### 9.3 Context & Trace ID Enrichment
```go
type traceKeyType string
const traceKey traceKeyType = "traceId"

ctx := context.WithValue(context.Background(), traceKey, "trace-xyz-789")

// Emits to all destinations in parallel / sequence
appErr := logger.Info(ctx, "Order processed", map[string]any{
    "orderId": "ord-883",
    "total":   149.50,
})
if appErr != nil {
    appErr.Print() // Nicely prints formatted fault
}
```

---

## 10. Complete End-to-End Walkthrough

```go
package main

import (
    "bytes"
    "context"
    "fmt"
    "os"

    "coding-guidelines/common/pkg/appfault"
    "coding-guidelines/common/pkg/errtype"
    "coding-guidelines/common/pkg/fileutil"
    "coding-guidelines/common/pkg/streamwriter"
)

func main() {
    ctx := context.Background()

    // 1. OPEN DISK FILE VIA TYPE-SAFE ENUMS
    fileWrap := fileutil.OpenFile(
        "tmp/demo-audit.log",
        fileutil.FileOpenCreateAppend,
        fileutil.FilePermStandard,
    )
    if fileWrap.IsFailed() {
        fileWrap.Print() // Nicely prints formatted fault
        return
    }
    logFile := fileWrap.Data()
    defer logFile.Close()

    // 2. CONSTRUCT DESTINATIONS
    buf := &bytes.Buffer{}
    lockedDiskStreamer := streamwriter.NewLockedStreamer[any](streamwriter.StreamerOptions[any]{
        Name:        "disk-streamer",
        Destination: logFile,
    })

    // 3. CONSTRUCT PLUGGABLE WRITERS
    consoleWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
        Name:        "console",
        Destination: os.Stdout,
    })

    memoryWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
        Name:        "memory-buffer",
        Destination: buf,
    })

    diskWriter := streamwriter.NewPluggableWriter[any](streamwriter.WriterOptions[any]{
        Name:     "disk",
        Streamer: lockedDiskStreamer,
    })

    // 4. ASSEMBLE COMPOSITE LOGGER
    logger := streamwriter.NewLogger[any]().
        AddWriters(consoleWriter, memoryWriter, diskWriter)

    fmt.Printf("Registered %d active writers\n\n", logger.WriterCount())

    // 5. EMIT STRUCTURED EVENT ACROSS ALL 3 DESTINATIONS
    _ = logger.Info(ctx, "System initialization complete", map[string]any{
        "version":     "2.4.0",
        "environment": "production",
    })

    // 6. HOT-SWAP FORMATTER AT RUNTIME (Custom JSON Formatting)
    consoleWriter.SetFormatMethod(func(payload any) streamwriter.Bytes[any] {
        jsonRes := streamwriter.JsonSource.FromPayload(payload)
        line := fmt.Sprintf("🔥 [CUSTOM-FORMATTER] %s\n", jsonRes.Pretty())
        return streamwriter.NewBytes([]byte(line), payload)
    })

    // 7. HOT-SWAP WRITE METHOD AT RUNTIME (Custom Watermarking)
    memoryWriter.SetWriteMethod(func(s streamwriter.Streamer[any], c context.Context, w *streamwriter.PluggableWriter[any], p any) *appfault.AppError {
        record := fmt.Sprintf("[WATERMARKED][%s] %s\n", w.Name(), streamwriter.Compile(p))
        _, err := w.Destination().Write([]byte(record))
        if err != nil {
            return appfault.Wrap(errtype.IO, err, "custom memory write failed")
        }
        return nil
    })

    // 8. EMIT SECOND EVENT WITH UPDATED HANDLERS
    _ = logger.Warn(ctx, "High memory utilization detected", map[string]any{
        "usagePct": 87.4,
    })

    // 9. VERIFY IN-MEMORY BUFFER
    fmt.Printf("\n--- In-Memory Buffer Contents ---\n%s", buf.String())
}
```

---

## 11. Summary Matrix

| Component | Package | Primary Type / Constructor | Key Capabilities |
| :--- | :--- | :--- | :--- |
| **Error Type** | `pkg/errtype` | `errtype.Variation` (`uint16`) | 16-bit / UTF-16 code enum (`Validation = 2`, `NotFound = 3`, `IO = 9`). Has `.Code()`, `.Name()`, `.HttpStatus()`. |
| **Error Carrier** | `pkg/appfault` | `*appfault.AppError` | Standard error return. Carries `errType Variation` (enum only, no string ErrorCode/ID) and `caller CallerInfo` (structured object, no string). |
| **Caller Object** | `pkg/appfault` | `appfault.CallerInfo` | Structured value object (embedded directly by value, zero heap allocation) containing `Function string`, `File string`, `Line int`. |
| **Display Subsystem** | `pkg/appfault` | `.Print()`, `.PrintWith()`, `.Format()`, `.ToClipboard()` | Default formatted printing with icon (`❌ [Validation:2] msg`), full diagnostic dumps, and markdown clipboard reports. |
| **Multi-Destination Display** | `pkg/appfault`, `pkg/result` | `.FormatStdout()`, `.FormatJSON()`, `.FormatTextLog()` | Dedicated presentation formatters for terminal ANSI banners, RFC JSON documents, and single-line log aggregator streams. |
| **Result Envelope** | `pkg/result` | `result.Wrap[T]` | Monadic container replacing `(T, error)`. Has `.Print()`, `.PrintFault()`, `.Format()`, and type propagation (`WrapFailureFromWrap`). |
| **Writer Envelope** | `pkg/appwriter` | `appwriter.BaseWriterWrap` | Dedicated wrapper for `*BaseWriter` with error ID and root cause constructors. |
| **File Permissions** | `pkg/fileutil` | `fileutil.FilePermType` (`uint32`) | POSIX octal bitmask enum with `.Mode()` conversion (`FilePermStandard`, `FilePermPrivate`). |
| **File Open Modes** | `pkg/fileutil` | `fileutil.FileOpenModeType` (`byte`) | Open flag enum with `.Flags()` conversion (`FileOpenReadOnly`, `FileOpenCreateAppend`, `FileOpenCreateTruncate`). |
| **File I/O Engine** | `pkg/fileutil` | `fileutil.OpenFile`, `WriteFile` | Safe filesystem utilities with automatic parent directory creation and AppError returns. |
| **Advanced File I/O** | `pkg/fileutil` | `fileutil.WriteAtomic`, `ReadChunked`, `NewFileWriter` | Atomic file swapping (crash-proof), pooled chunk streaming for large files, and PluggableWriter bridge. |
| **JSON Envelope** | `pkg/streamwriter` | `streamwriter.JsonResult` | Rich JSON wrapper with `.Pretty()`, `.Compact()`, and multi-source factory (`JsonSource`). |
| **Deterministic Compiler** | `pkg/streamwriter` | `streamwriter.Compile` | Recursive object serializer guaranteeing deterministic key ordering. |
| **Payload Intelligence** | `pkg/streamwriter` | `InspectPayload`, `ExtractBytes`, `NewAnyWriter` | Reflect converter preventing base64 double-marshaling of raw `[]byte`, bypassing string quotes, and AnyWriter non-generic ergonomics. |
| **Pluggable Writer** | `pkg/streamwriter` | `streamwriter.PluggableWriter[T]`, `AnyWriter` | Composable writer with runtime swappable formatting, write handlers (`WriteFunc`), and destinations. |
| **Streamers** | `pkg/streamwriter` | `LockedStreamer`, `LocklessStreamer` | Thread-safe (`sync.Locker`) or zero-lock high-throughput streaming abstractions. |
| **Composite Logger** | `pkg/streamwriter` | `streamwriter.Logger[T]` | Fluent multi-destination fan-out coordinator with zero-allocation silent mode. |
