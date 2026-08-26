# Go Utility Packages: `pathutil` and `fileutil`

> **Version:** 1.0.0  
> **Updated:** 2026-03-31  
> **Applies to:** All Go backend code  
> **Cross-refs:** [no-negatives](../01-cross-language/12-no-negatives.md), [golang-standards-reference](./04-golang-standards-reference/00-overview.md), [error-code-registry](../../03-error-manage/03-error-code-registry/01-registry.md#gen-700-file-system)

---

## Purpose

Two utility packages eliminate raw `os`/`io` calls from application code:

| Package | Role | Returns |
|---------|------|---------|
| `pathutil` | Boolean guard functions for file/directory existence checks | `bool` — positive-named, no `!` needed |
| `fileutil` | File I/O wrappers that return `apperror.Result[T]` | `apperror.Result[T]` — no raw `error` escaping |

**Rule:** Application code MUST NOT call `os.Open`, `os.Stat`, `os.ReadFile`, `os.WriteFile`, `os.MkdirAll`, or `os.Remove` directly. Use the corresponding `pathutil` or `fileutil` wrapper.

---

## `pathutil` — Boolean Guards

All functions return `bool`. Every negative check has a positive-named counterpart so callers never use `!`.

### Function Signatures

```go
package pathutil

// File existence
func IsFileExists(path string) bool   // wraps os.Stat — true if file exists and is not a directory
func IsFileMissing(path string) bool  // !IsFileExists(path)
func IsFileValid(path string) bool    // exists, not dir, size > 0

// Directory existence
func IsDir(path string) bool          // wraps os.Stat — true if path exists and is a directory
func IsDirMissing(path string) bool   // !IsDir(path)
func IsDirWritable(path string) bool  // IsDir(path) && write permission test
func IsDirReadonly(path string) bool  // !IsDirWritable(path)

// Path utilities
func IsPathAbsolute(path string) bool // filepath.IsAbs wrapper
func IsPathRelative(path string) bool // !IsPathAbsolute(path)

// Cleanup
func Remove(path string) error        // os.Remove — only util that returns raw error (stdlib boundary)
func RemoveAll(path string) error      // os.RemoveAll — stdlib boundary
```

### Usage Examples

```go
// ❌ FORBIDDEN: Raw negation on stdlib
if !pathutil.IsDir(gitDir) {
    return apperror.FailNew[StatusResult](apperror.ErrGitNotRepo, "not a git repo")
}

// ✅ REQUIRED: Positive-named guard
if pathutil.IsDirMissing(gitDir) {
    return apperror.FailNew[StatusResult](apperror.ErrGitNotRepo, "not a git repo")
}
```

```go
// ❌ FORBIDDEN: Mixed polarity
if config.IsEnabled() && !pathutil.IsDir(exportDir) {
    createDir(exportDir)
}

// ✅ REQUIRED: Named boolean, no negation
isEnabledWithMissingDir := config.IsEnabled() && pathutil.IsDirMissing(exportDir)
if isEnabledWithMissingDir {
    createDir(exportDir)
}
```

```go
// ✅ Zip validation after creation
isValid := pathutil.IsFileValid(absZipPath)

if !isValid {
    pathutil.Remove(absZipPath)

    return apperror.FailNew[ExportResult](apperror.ErrZipInvalid, "zip file invalid after creation")
}
```

### Naming Convention

| Pattern | Positive | Negative |
|---------|----------|----------|
| File exists | `IsFileExists` | `IsFileMissing` |
| Dir exists | `IsDir` | `IsDirMissing` |
| Dir writable | `IsDirWritable` | `IsDirReadonly` |
| Path absolute | `IsPathAbsolute` | `IsPathRelative` |

**Rule:** Every `Is<Thing>` MUST have its negation counterpart so callers never write `!pathutil.Is<Thing>()`.

---

## `fileutil` — `apperror.Result[T]` Wrappers

All functions return `apperror.Result[T]`, wrapping stdlib errors at the boundary. Application code receives typed results with error codes — no raw `error` propagation.

### Function Signatures

```go
package fileutil

// File open/read
func Open(path string) apperror.Result[*os.File]       // os.Open → ErrFileNotFound / ErrFileReadFailed
func ReadAll(path string) apperror.Result[[]byte]       // os.ReadFile → ErrFileNotFound / ErrFileReadFailed
func ReadString(path string) apperror.Result[string]    // ReadAll + string conversion

// File write
func WriteFile(path string, data []byte) apperror.Result[bool]          // os.WriteFile → ErrFileWriteFailed
func WriteFileWithPerm(path string, data []byte, perm os.FileMode) apperror.Result[bool]

// Directory
func MkdirAll(path string, perm os.FileMode) apperror.Result[bool]     // os.MkdirAll → ErrDirCreateFailed
func EnsureDir(path string) apperror.Result[bool]                       // MkdirAll with 0755 default

// File info
func Stat(path string) apperror.Result[os.FileInfo]     // os.Stat → ErrFileNotFound
func FileSize(path string) apperror.Result[int64]        // Stat + Size()
```

### Error Code Mapping

All error codes are from the [GEN-700: File System](../../03-error-manage/03-error-code-registry/01-registry.md#gen-700-file-system) range:

| Function | Success | Error Code | Error Constant |
|----------|---------|------------|----------------|
| `Open` | `*os.File` | GEN-700-01 | `ErrFileNotFound` |
| `ReadAll` | `[]byte` | GEN-700-02 | `ErrFileReadFailed` |
| `WriteFile` | `true` | GEN-700-03 | `ErrFileWriteFailed` |
| `MkdirAll` | `true` | GEN-700-04 | `ErrDirCreateFailed` |
| `Stat` | `os.FileInfo` | GEN-700-01 | `ErrFileNotFound` |
| `EnsureDir` | `true` | GEN-700-04 | `ErrDirCreateFailed` |

Permission errors (any function): GEN-700-05 / `ErrFsPermissionDenied`

### Internal Implementation Pattern

```go
// Every fileutil function follows this pattern:
func Open(path string) apperror.Result[*os.File] {
    file, err := os.Open(path)

    if err == nil {
        return apperror.Ok(file)
    }

    if os.IsNotExist(err) {
        return apperror.FailWrap[*os.File](err, apperror.ErrFileNotFound, "file not found: "+path)
    }

    if os.IsPermission(err) {
        return apperror.FailWrap[*os.File](err, apperror.ErrFsPermissionDenied, "permission denied: "+path)
    }

    return apperror.FailWrap[*os.File](err, apperror.ErrFileReadFailed, "failed to open: "+path)
}
```

### Usage Examples

```go
// ❌ FORBIDDEN: Raw os.Open in application code
func (s *PluginService) Upload(ctx context.Context, req UploadRequest) apperror.Result[UploadResult] {
    file, err := os.Open(req.Path)
    if err != nil {
        return apperror.FailWrap[UploadResult](err, "E4001", "failed to open file: "+req.Path)
    }
    defer file.Close()
    // ...
}

// ✅ REQUIRED: Use fileutil — error codes are consistent, wrapping is automatic
func (s *PluginService) Upload(ctx context.Context, req UploadRequest) apperror.Result[UploadResult] {
    result := fileutil.Open(req.Path)
    if result.IsErr() {
        return apperror.Fail[UploadResult](result.AppError())
    }
    file := result.Value
    defer file.Close()
    // ...
}
```

```go
// ✅ Reading config file
func LoadConfig(path string) apperror.Result[AppConfig] {
    data := fileutil.ReadAll(path)

    if data.IsErr() {
        return apperror.Fail[AppConfig](data.AppError())
    }

    var config AppConfig

    if err := json.Unmarshal(data.Value, &config); err != nil {
        return apperror.FailWrap[AppConfig](err, apperror.ErrJsonDecode, "invalid config JSON")
    }

    return apperror.Ok(config)
}
```

```go
// ✅ Ensuring output directory exists before writing
func (s *ExportService) SaveExport(dir, name string, data []byte) apperror.Result[string] {
    dirResult := fileutil.EnsureDir(dir)

    if dirResult.IsErr() {
        return apperror.Fail[string](dirResult.AppError())
    }

    outPath := filepath.Join(dir, name)
    writeResult := fileutil.WriteFile(outPath, data)

    if writeResult.IsErr() {
        return apperror.Fail[string](writeResult.AppError())
    }

    return apperror.Ok(outPath)
}
```

---

## Decision Matrix: `pathutil` vs `fileutil`

| Need | Use | Why |
|------|-----|-----|
| "Does this path exist?" (boolean check) | `pathutil.IsFileExists` / `pathutil.IsDir` | No error propagation needed |
| "Open/read/write a file" (I/O operation) | `fileutil.Open` / `fileutil.ReadAll` / `fileutil.WriteFile` | Returns `apperror.Result[T]` with proper error codes |
| "Guard clause before I/O" | `pathutil` check → early return → `fileutil` operation | Boolean guard first, then typed I/O |
| "Clean up a temp file" | `pathutil.Remove` | Fire-and-forget cleanup at stdlib boundary |

---

## PHP Equivalent: `PathHelper`

PHP uses a static helper class with the same positive-naming principle:

| Go (`pathutil`) | PHP (`PathHelper`) |
|-----------------|-------------------|
| `pathutil.IsFileExists(p)` | `PathHelper::isFileExists($p)` |
| `pathutil.IsFileMissing(p)` | `PathHelper::isFileMissing($p)` |
| `pathutil.IsDir(p)` | `PathHelper::isDirExists($p)` |
| `pathutil.IsDirMissing(p)` | `PathHelper::isDirMissing($p)` |
| `pathutil.IsDirWritable(p)` | `PathHelper::isDirWritable($p)` |
| `pathutil.IsDirReadonly(p)` | `PathHelper::isDirReadonly($p)` |

PHP namespace: `RiseupAsia\Helpers\PathHelper`

See [no-negatives spec](../01-cross-language/12-no-negatives.md) for the full PHP guard method table.

---

## AI Checklist

```
[ ] No raw os.Open / os.Stat / os.ReadFile / os.WriteFile in application code — use fileutil
[ ] No !pathutil.Is<X>() — use the named negative counterpart
[ ] Every pathutil boolean has both positive and negative versions
[ ] fileutil functions return apperror.Result[T], never raw error
[ ] Error codes match GEN-700 range from error-code-registry
[ ] PHP PathHelper methods mirror Go pathutil naming
```
