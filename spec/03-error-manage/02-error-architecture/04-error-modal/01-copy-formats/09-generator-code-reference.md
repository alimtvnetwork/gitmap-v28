# Generator Code Reference

> **Parent:** [Copy Formats Index](./00-overview.md)  
> **Version:** 3.2.0  
> **Updated:** 2026-03-31  
> **Purpose:** Source files, function signatures, and replication guide for the error report generation system.

---

> Full architectural context: [03-error-modal-reference.md § 11 – Error Report Generation](../03-error-modal-reference.md#11-error-report-generation)

---

## Source Files

| File | Purpose |
|------|---------|
| `src/components/errors/errorReportGenerator.ts` | Pure-function report generators (no React deps) |
| `src/components/errors/errorLogAdapter.ts` | Maps backend `ErrorLog` → `CapturedError` for `ErrorDetailModal` |

---

## Function Signatures

```typescript
// Full report — all sections, verbose
function generateErrorReport(
  error: CapturedError,
  app?: { appName: string; appVersion: string; gitCommit?: string; buildTime?: string }
): string

// Compact report — essential sections, stripped execution chain, backend from CapturedError
function generateCompactReport(
  error: CapturedError,
  app?: { appName: string; appVersion: string; gitCommit?: string; buildTime?: string }
): string

// Adapter — ErrorDetailModal uses this before calling either generator
function errorLogToCapturedError(error: ErrorLog): CapturedError
```

---

## Dependencies

| Import | Purpose |
|--------|---------|
| `CapturedError` from `@/stores/errorStore` | Error data type |
| `ErrorLog` from `@/lib/api` | Backend API error shape (used by adapter) |
| `formatDateTimeUtc` from `@/lib/logText` | Timestamp formatting |
| `toClipboardText` from `@/lib/logText` | Clipboard text sanitization |
| `unescapeEmbeddedNewlines` from `@/lib/logText` | Newline unescaping in log entries |

---

## Replication Guide

To replicate this error reporting system in another project:

1. **Copy `errorReportGenerator.ts`** — Pure function, no side effects
2. **Copy `errorLogAdapter.ts`** — If you need to feed backend-stored errors into the generators
3. **Implement `CapturedError` interface** — See Section 2 of [03-error-modal-reference.md](../03-error-modal-reference.md)
4. **Implement `errorStore`** — Zustand store with `captureError()`, `buildCapturedError()`
5. **Implement `parseEnvelope()`** — Extract `Errors`, `MethodsStack`, `Attributes` from API responses
6. **Implement click path tracking** — DOM click listener that records interactive element clicks
7. **Implement execution logger** — Optional debug-mode call chain tracker
8. **Wire Copy/Download menus** — Use the menu structures documented in [00-overview.md](./00-overview.md)

---

*Generator code reference — updated: 2026-03-31*
