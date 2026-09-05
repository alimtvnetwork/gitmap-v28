# `FileWrapper` Architecture, `FileModeType` Enums, and Error Object Integration

> **Document:** `research/12-fileutil-wrapper-and-filemode-enums.md`  
> **Status:** Implemented & Verified  
> **Package Reference:** `04-code/golang/pkg/fileutil` & `04-code/golang/pkg/appfault`  
> **Date:** 2026-09-05  

---

## 1. Executive Summary & Context

File I/O operations in high-reliability Go applications frequently suffer from string-based error handling, loose integer mode constants (e.g., `0644`), and unstructured error wrapping. 

Drawing inspiration from `03-aukgo/core/filemode` and `03-aukgo/core/chmodhelper`, this architecture formalizes:
1. **Type-Safe Enum Architecture**:
   - `FileActionType`: Defines file operations (`Create`, `Delete`, `ReadOnly`, `Read`, `Write`, `Append`, `ReadWrite`, `Truncate`) with explicit OS flag translation.
   - `FileModeType`: Typed octal file permission models (`DefaultFile` `0644`, `DefaultDir` `0755`, `ReadOnly` `0444`, `WriteOnly` `0222`, `ReadWrite` `0666`, `OwnerOnly` `0700`, `AllPermission` `0777`) with permission inspection helpers (`IsExecutable`, `IsWritable`, `IsReadOnly`).
2. **Error Object Standard (`*appfault.AppError`)**:
   - Replaces primitive error wrappers with first-class `AppError` objects carrying an explicit `errorId`.
   - Dedicated wrap constructors (`WrapFailure`, `WrapReaderFailure`, `WrapWriterFailure`) that bind domain failure codes, error IDs, causes, and descriptive messages.
3. **`FileWrapper` Encapsulation**:
   - A composable file utility struct encapsulating file paths, default modes, and actions.
   - Methods for constructing failures with error IDs directly from the wrapper definition.
   - High-level, safety-first file operations (`Read`, `ReadString`, `Write`, `WriteString`, `Append`, `Delete`, `Create`, `Execute`).

---

## 2. Type-Safe Enum Specifications

### 2.1 `FileActionType`
All enum names strictly follow the `Type` suffix convention per repository guidelines.

```go
type FileActionType string

const (
    ActionCreate    FileActionType = "create"
    ActionDelete    FileActionType = "delete"
    ActionReadOnly  FileActionType = "read_only"
    ActionRead      FileActionType = "read"
    ActionWrite     FileActionType = "write"
    ActionAppend    FileActionType = "append"
    ActionReadWrite FileActionType = "read_write"
    ActionTruncate  FileActionType = "truncate"
)
```

#### OS Flags Mapping
The `ToOSFlags()` method translates each `FileActionType` to canonical `os.O_*` bitmasks:
- `ActionReadOnly` / `ActionRead` $\rightarrow$ `os.O_RDONLY`
- `ActionWrite` $\rightarrow$ `os.O_WRONLY|os.O_CREATE|os.O_TRUNC`
- `ActionAppend` $\rightarrow$ `os.O_WRONLY|os.O_CREATE|os.O_APPEND`
- `ActionCreate` $\rightarrow$ `os.O_RDWR|os.O_CREATE|os.O_EXCL`
- `ActionReadWrite` $\rightarrow$ `os.O_RDWR|os.O_CREATE`
- `ActionTruncate` $\rightarrow$ `os.O_RDWR|os.O_TRUNC`

### 2.2 `FileModeType`
Octal permissions are typed and represented via `FileModeType`:

```go
type FileModeType uint32

const (
    ModeDefaultFile   FileModeType = 0644
    ModeDefaultDir    FileModeType = 0755
    ModeReadOnly      FileModeType = 0444
    ModeWriteOnly     FileModeType = 0222
    ModeReadWrite     FileModeType = 0666
    ModeOwnerOnly     FileModeType = 0700
    ModeAllPermission FileModeType = 0777
)
```

Methods provide positive boolean inspection:
- `ToFileMode() os.FileMode`: Direct cast to standard `os.FileMode`.
- `IsExecutable() bool`: Returns `true` if any execute bit (owner, group, other) is set.
- `IsWritable() bool`: Returns `true` if any write bit is set.
- `IsReadOnly() bool`: Returns `true` if write bits are absent and read bits are present.

---

## 3. Error Object Architecture with `errorId`

Rather than returning plain `fmt.Errorf` string wrappers, failures produce a structured `*appfault.AppError` instance containing:
- `errType`: Enumerated error domain (e.g., `errtype.StorageIOWriteFailure`, `errtype.StorageIOReadFailure`).
- `errorId`: Unique tracking or correlation identifier (e.g., `"ERR_FILE_NOT_FOUND"`, `"ERR_FILE_WRITE"`).
- `cause`: Underlying root-cause `error`.
- `message`: Contextual human-readable description.

### 3.1 Constructor Signatures in `appfault`
```go
func WrapFailure(errType errtype.ErrorType, errorId string, cause error, message string) *AppError
func WrapWriterFailure(errorId string, cause error, message string) *AppError
func WrapReaderFailure(errorId string, cause error, message string) *AppError
func NewWithId(errType errtype.ErrorType, errorId string, message string) *AppError
```

### 3.2 Wrapper-Bound Error Methods
Inside `FileWrapper`, failure construction methods bind the current target file path into the generated `AppError`:
```go
func (w *FileWrapper) WrapFailure(errType errtype.ErrorType, errorId string, cause error, msg string) *appfault.AppError
func (w *FileWrapper) WrapReaderFailure(errorId string, cause error, msg string) *appfault.AppError
func (w *FileWrapper) WrapWriterFailure(errorId string, cause error, msg string) *appfault.AppError
```

---

## 4. `FileWrapper` API Reference

### 4.1 Struct Definition
```go
type FileWrapper struct {
    filePath string
    mode     FileModeType
    action   FileActionType
}
```

### 4.2 Factory Functions
- `NewFileWrapper(filePath string) *FileWrapper`: Initializes with `ModeDefaultFile` (0644) and `ActionRead`.
- `NewFileWrapperWithOptions(filePath string, mode FileModeType, action FileActionType) *FileWrapper`: Configures custom permissions and action.
- `NewDefault() *FileWrapper`: Creates an unanchored wrapper with default settings.

### 4.3 Execution Interface
- `Read() ([]byte, *appfault.AppError)`
- `ReadString() (string, *appfault.AppError)`
- `Write(data []byte) *appfault.AppError`
- `WriteString(content string) *appfault.AppError`
- `Append(data []byte) *appfault.AppError`
- `Delete() *appfault.AppError`
- `Create() *appfault.AppError`
- `Execute(data []byte) ([]byte, *appfault.AppError)`: Dynamically executes the configured `FileActionType`.

---

## 5. Verification & Quality Compliance

1. **Unit Test Coverage**: `04-code/golang/pkg/fileutil/fileutil_test.go` verifies all enums, read/write/append/delete operations, dynamic `Execute`, and structured `errorId` preservation.
2. **Coding Guideline Compliance**:
   - Every function strictly adheres to the $\le 15$ lines rule.
   - All returns are preceded by a blank line.
   - Positive boolean prefixes (`isSuccess`, `isExist`, `isExec`).
   - Zero swallowed errors.
