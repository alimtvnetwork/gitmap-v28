# Error Report Generation

> **Parent:** [Error Modal Reference](./00-overview.md)  
> **Version:** 2.2.0  
> **Updated:** 2026-03-31  
> **See also:** [Copy Format Samples](../01-copy-formats/00-overview.md) for complete output examples.

---

The `errorReportGenerator.ts` produces two report formats from any `CapturedError`:

## Compact Report (Default)

The **Compact Error Report** is the default copy format. Designed for high-signal diagnostic sharing — stripped of noise, built entirely from client-side memory (no API calls). The main "Copy" button uses this format.

```typescript
export function generateCompactReport(error: CapturedError, app?: AppInfo): string {
  // Sections: App metadata, code+level, page, user interaction, trigger context,
  // message, request, frontend execution chain (stripped), context JSON,
  // frontend stack trace, backend error.log.txt (from CapturedError memory)
  return markdownReport;
}
```

**Stripping rules for execution chain:**
- Timestamps (e.g., `[12:58:22 AM] ⬡`) are removed
- Base API URLs (e.g., `http://localhost:8080/api/v1`) are stripped to relative paths
- Result: clean, scannable list like `GET /sites`, `POST /error-history`

## Full Report

```typescript
export function generateErrorReport(error: CapturedError, app?: AppInfo): string {
  // 17 sections including full URLs, backend logs, PHP traces, etc.
  return markdownReport;
}
```

## ErrorDetailModal Adapter

The standalone `ErrorDetailModal` receives `ErrorLog` objects from the backend API, not `CapturedError`. The `errorLogAdapter.ts` bridges this gap:

```typescript
// src/components/errors/errorLogAdapter.ts
export function errorLogToCapturedError(error: ErrorLog): CapturedError;
```

---

## Download/Copy Options

**Copy (Split Button):**

| Action | Trigger | Content |
|--------|---------|---------|
| **Copy (main click)** | Primary button area | Compact Report (instant, from memory) |
| **Copy Compact Report** | Chevron dropdown | Same as main click |
| **Copy Full Report** | Chevron dropdown | Full 17-section Markdown report |
| **Copy with Backend Logs** | Chevron dropdown | Full report + error.log.txt (API fetch) |
| **Copy error.log.txt** | Chevron dropdown | Raw error log content (API fetch) |
| **Copy log.txt** | Chevron dropdown | Raw full log content (API fetch) |

**Download:**

| Action | Content |
|--------|---------|
| **Full Bundle (ZIP)** | Markdown report + logs, sent to `/api/v1/errors/bundle` |
| **error.log.txt** | Backend error log file |
| **log.txt** | Backend full diagnostic log |
| **Report (.md)** | Markdown error report |

---

*Report generation — updated: 2026-03-31*
