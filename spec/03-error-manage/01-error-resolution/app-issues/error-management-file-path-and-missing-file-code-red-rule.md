# 🔴 CODE RED: File Path and Missing File Error Logging Rule

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Priority:** 🔴 Code Red — Zero Tolerance  
**AI Confidence:** High  
**Ambiguity:** Low (see § 9 for edge cases)

---

## 1. Purpose

Every file-related or path-related error **must** include the exact file path and a clear reason for failure in the error log. Generic "file not found" messages without path details are a **🔴 Code Red violation**.

This rule applies system-wide — not just in one module.

---

## 2. Mandatory Logging Fields

Every file/path error log entry **must** include:

| Field | Required | Description |
|-------|----------|-------------|
| Error level | ✅ | Code Red or equivalent severity |
| Exact file path | ✅ | Full resolved path that was attempted |
| Operation | ✅ | What was being done (see § 2.1) |
| Failure reason | ✅ | Why the operation failed (see § 3) |
| Module / component | ✅ | Which module or component triggered the error |
| Timestamp | ✅ | Standard logging timestamp |
| Recovery action | ⚠️ | Fallback action taken, if any |

### 2.1 Operations

The operation field must use one of these values:

```
Read | Write | Copy | Move | Inject | Load | Extract | Resolve
```

---

## 3. Acceptable Failure Reasons

The failure reason **must** be specific. Acceptable values include:

| Reason | When to Use |
|--------|-------------|
| `FileDoesNotExist` | File is confirmed absent at the expected path |
| `PathInvalid` | Path string is malformed or unresolvable |
| `PathInaccessible` | Path exists but cannot be accessed (network, mount) |
| `FileNameMismatch` | Expected filename differs from actual |
| `ExtensionMismatch` | Expected file extension differs from actual |
| `PermissionDenied` | OS-level permission prevents access |
| `FileNeverCreated` | File was expected from a prior pipeline step but was never generated |
| `FileMovedOrRenamed` | File existed previously but was moved, renamed, or deleted |

---

## 4. File and Path Error Categories

### 4.1 Missing File

```
🔴 [Code Red] File not found
  Path: /app/uploads/theme-v2.zip
  Reason: FileDoesNotExist — directory listing confirms no matching file
  Operation: Extract
  Module: ThemeInstaller
```

### 4.2 Invalid Path

```
🔴 [Code Red] Path resolution failed
  Path: ../../../etc/passwd (resolved: /etc/passwd)
  Reason: PathInvalid — path escapes allowed directory
  Operation: Resolve
  Module: AssetLoader
```

### 4.3 Missing Generated File

```
🔴 [Code Red] Expected output file not found
  Path: /tmp/build/output/bundle.js
  Reason: FileNeverCreated — build step "compile" exited with error
  Operation: Read
  Module: DeploymentPipeline
```

### 4.4 Missing Uploaded or Extracted File

```
🔴 [Code Red] Extracted file not found
  Path: /tmp/extract/wp-content/plugins/my-plugin/main.php
  Reason: FileDoesNotExist — archive did not contain expected entry
  Operation: Extract
  Module: PluginInstaller
```

### 4.5 Injection or Asset Load Failure

```
🔴 [Code Red] Asset load failed
  Path: /app/assets/custom.css
  Reason: PermissionDenied — file mode 0000
  Operation: Load
  Module: CSSInjector
```

---

## 5. Implementation Patterns

### 5.1 Go

```go
// ✅ Correct — exact path + reason
apperror.New(apperrtype.FileNotFound,
    "File does not exist",
    apperror.WithField("path", filePath),
    apperror.WithField("operation", "Read"),
    apperror.WithField("reason", "FileDoesNotExist"),
    apperror.WithField("module", "ThemeInstaller"),
)

// 🔴 VIOLATION — generic message, no path, no reason
apperror.New(apperrtype.FileNotFound, "file not found")
```

### 5.2 TypeScript

```typescript
// ✅ Correct
captureError({
  code: "E5010",
  level: "error",
  message: `File not found: ${filePath}`,
  context: {
    path: filePath,
    operation: "Load",
    reason: "FileDoesNotExist",
    module: "AssetLoader",
  },
});

// 🔴 VIOLATION
console.error("File not found");
```

### 5.3 PHP

```php
// ✅ Correct
throw new AppError(
    code: 'E5010',
    message: "File not found: {$filePath}",
    context: [
        'path' => $filePath,
        'operation' => 'Read',
        'reason' => 'FileDoesNotExist',
        'module' => 'PluginLoader',
    ]
);

// 🔴 VIOLATION
throw new \Exception("File not found");
```

---

## 6. System-Wide Applicability

This rule applies to **all** areas where file/path errors can occur:

| Area | Examples |
|------|----------|
| Injection pipeline | Script, CSS, HTML injection failures |
| Cache rebuild / load | Cache file missing after rebuild |
| Extraction / unzip | Archive extraction target missing |
| Script resolution | Dependency file not found |
| Spec file references | Cross-reference target missing |
| Asset / CSS loading | Static asset load failures |
| Deployment cleanup | Path cleanup target missing |
| Centralized logger | Any error logger module |

---

## 7. Validation Rule

A linter or code review check **must** reject any file/path error log that:

1. Does not include the exact file path
2. Does not include a failure reason from the approved list (§ 3)
3. Uses a generic message like "file not found" without context

---

## 8. Error Template

Standard template for all file/path errors:

```
🔴 [Code Red] {summary}
  Path: {exact_file_path}
  Reason: {reason_code} — {human_readable_detail}
  Operation: {operation}
  Module: {module_name}
  Recovery: {fallback_action_or_none}
```

---

## 9. Ambiguities and Clarifications

| Question | Decision |
|----------|----------|
| Do warnings also require exact paths? | **Yes** — all file-related warnings must include exact paths |
| Do sensitive paths require masking? | **Not defined** — document separately if masking is needed for production logs |
| Does this apply to UI error surfaces? | **Yes** for DevTools/debug output; user-facing messages may use friendly text but the underlying log must have the full path |

---

## 10. Monitoring Recommendations

1. Add a validation rule preventing generic "file not found" logs without path detail
2. Monitor for repeated failures on the same path (indicates systemic issue)
3. Distinguish between missing source files and missing generated files in dashboards

---

## Cross-References

- [Code Severity Taxonomy](../../../02-coding-guidelines/03-golang/07-code-severity-taxonomy.md) — Code Red classification
- [Error Documentation Guideline](../00-error-documentation-guideline.md) — Mandatory documentation process
- [Data Model: CapturedError](../../02-error-architecture/04-error-modal/03-error-modal-reference/01-data-model.md) — `file`, `context` fields
- [Error Code Registry](../../03-error-code-registry/02-integration-guide.md) — Error code integration

---

## Do you understand?

```
If you are an AI reading this specification:

1. Every file-related or path-related error MUST include the exact file path in the log.
2. Every such error MUST include the reason why the file was not found or the path failed.
3. This is 🔴 Code Red priority — zero tolerance for generic file-not-found messages.
4. Apply this rule everywhere file/path errors can occur — not just one module.
5. When creating or reviewing code, reject any file/path error that omits the path or reason.

Do you understand? Can you please do that?
```

---

*File path and missing file error logging rule — created: 2026-04-07*
