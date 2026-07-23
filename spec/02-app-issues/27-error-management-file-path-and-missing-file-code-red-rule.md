# Error Management — File Path & Missing File (Code Red Rule)

## Priority

**CODE RED** — This rule is mandatory and must never be bypassed.

---

## Rule

Every error related to a file or path MUST include:

1. **Exact file path** — the full resolved path that was attempted.
2. **Failure reason** — why the file was not found or the operation failed.

Generic messages like "file not found" without a path are **prohibited**.

---

## Mandatory Log Fields

| Field              | Required | Description                                      |
|--------------------|----------|--------------------------------------------------|
| Exact file path    | ✅        | Full resolved path that was attempted             |
| Operation          | ✅        | read, write, copy, move, extract, load, resolve   |
| Failure reason     | ✅        | Why the file was not found or operation failed     |
| Module / component | ✅        | Which package or function encountered the error    |
| Recovery action    | Optional | Fallback taken, if any                            |
| Timestamp          | Auto     | Included by verbose logger when enabled            |

---

## Failure Reason Categories

| Reason                         | When to Use                                      |
|--------------------------------|--------------------------------------------------|
| File does not exist            | `os.Stat` or `os.Open` returns `ErrNotExist`     |
| Path is invalid                | Malformed or empty path string                    |
| Path is inaccessible           | Permission denied or OS restriction               |
| Permission denied              | Explicit permission error from OS                 |
| File name mismatch             | Expected name differs from actual                 |
| Extension mismatch             | Wrong file extension for the operation             |
| File was never created         | Expected from a prior step but absent              |
| File was removed or renamed    | Previously existed but no longer present           |
| Archive entry missing          | Expected entry not found inside zip/tar            |
| Environment variable not set   | Path depends on env var that is empty              |

---

## Error Format Standard

All file/path error constants MUST follow this pattern:

```go
// ✅ Correct — includes path and reason
ErrSettingsZipNotFound = "Error: settings zip not found at %s: %v\n"
ErrExtractEntryFailed  = "Error: failed to extract %s to %s: %v (operation: extract)\n"
ErrSettingsSourceDir   = "Error: settings source directory not found at %s: %v\n"

// ❌ Wrong — missing path or reason
"Settings zip not found"
"Failed to extract settings"
"Source not found"
```

---

## Affected Areas

This rule applies to every module that handles files or paths:

| Area                          | Examples                                           |
|-------------------------------|----------------------------------------------------|
| Install / NPP settings       | Zip extraction, settings copy, exe verification     |
| Clone pipeline                | Source file loading, target directory creation       |
| Release pipeline              | Asset upload, compression, checksums, metadata       |
| Configuration loading         | Config files, task files, gitignore patterns         |
| Database operations           | SQLite file open, migration files                    |
| SSH key management            | Key file read/write, config file update              |
| Export / Import               | File read, file write, format detection              |
| GoMod rename                  | go.mod read/write, source file scanning              |
| Verbose logging               | Log file creation                                    |
| Settings / Fallback           | Loose file copy when zip is missing                  |

---

## Implementation Pattern

```go
// ✅ Correct — file error with path and reason
data, err := os.ReadFile(path)
if err != nil {
    fmt.Fprintf(os.Stderr, constants.ErrFileReadFailed, path, err)
    return
}

// ✅ Correct — zip extraction with entry path
src, err := file.Open()
if err != nil {
    fmt.Fprintf(os.Stderr, constants.ErrExtractEntryFailed, file.Name, destPath, err)
    return
}

// ✅ Correct — fallback with source path
entries, err := os.ReadDir(source)
if err != nil {
    fmt.Fprintf(os.Stderr, constants.ErrSettingsSourceDir, source, err)
    return
}
```

---

## Validation Checklist

Before merging any code that handles files:

- [ ] Every `os.Open`, `os.ReadFile`, `os.Stat`, `os.Create` error includes the path
- [ ] Every `zip.OpenReader` error includes the zip path
- [ ] Every `os.ReadDir` error includes the directory path
- [ ] Every fallback logs what was attempted and why the primary path failed
- [ ] No error message says "file not found" without the exact path
- [ ] No error is silently swallowed (empty `if err != nil { return }`)

---

## Ambiguities Noted

1. **Warnings vs errors**: This rule currently applies to errors only. Whether warnings
   should also include exact paths is not yet defined.
2. **Sensitive paths**: Whether paths containing usernames require masking in shared logs
   is not yet defined.
3. **UI surfaces**: Whether this rule extends to browser DevTools output or only CLI
   stderr is not yet defined.

---

## See Also

- [Error Handling](../04-generic-cli/07-error-handling.md) — CLI error handling patterns
- [Error Handling Patterns](../05-coding-guidelines/04-error-handling.md) — Cross-language error patterns
- [Verbose Logging](../04-generic-cli/16-verbose-logging.md) — Debug logging system

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC](https://riseup-asia.com) (2026)
