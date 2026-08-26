# Error Modal — Copy & Export Formats (Index)

> **Parent:** [Error Modal Spec](../00-overview.md)  
> **Version:** 3.2.0  
> **Updated:** 2026-03-31  
> **Status:** Active  
> **AI Confidence:** 95%  
> **Ambiguity Score:** 5%  
> **Purpose:** Complete, copy-pasteable samples of every error report format produced by the Global Error Modal. Each format lives in its own file for focused AI consumption.

---

## File Index

| # | File | Format | Description |
|---|------|--------|-------------|
| 01 | [01-compact-report.md](./01-compact-report.md) | Compact Report (Markdown) ⭐ | **DEFAULT** — stripped-down, instant copy (no API call). Includes delegated server info built from CapturedError |
| 02 | [02-full-report.md](./02-full-report.md) | Full Report (Markdown) | All frontend + backend diagnostics, verbose |
| 03 | [03-full-report-with-backend-logs.md](./03-full-report-with-backend-logs.md) | Full Report + Backend Logs | Full Report with error.log.txt appended (async, fetches API) |
| 04 | [04-error-log-txt.md](./04-error-log-txt.md) | error.log.txt | Raw backend error log from `GET /api/v1/logs/error` |
| 05 | [05-full-log-txt.md](./05-full-log-txt.md) | log.txt | Raw backend full log from `GET /api/v1/logs/full` |
| 06 | [06-error-log-with-delegated-info.md](./06-error-log-with-delegated-info.md) | error.log.txt + Delegated Server | Enhanced error log with downstream server diagnostics |
| 07 | [07-envelope-error-response.md](./07-envelope-error-response.md) | Envelope Error (JSON) | Raw Go backend JSON error response |
| 08 | [08-session-diagnostics.md](./08-session-diagnostics.md) | Session Diagnostics (JSON) | Session-linked request/response data |
| 09 | [09-generator-code-reference.md](./09-generator-code-reference.md) | Generator Code | Source files, function signatures, replication guide |

---

## Format Overview

| Format | Trigger | Content |
|--------|---------|---------|
| **Compact Report** ⭐ | Copy button (main click) | Stripped-down Markdown — key diagnostics + backend error.log built from CapturedError |
| **Full Report** | Copy ▸ Copy Full Report | Markdown with all frontend + backend diagnostics |
| **Report + Backend Logs** | Copy ▸ Copy with Backend Logs | Full Report + error.log.txt appended (fetched from API) |
| **error.log.txt** | Copy ▸ Copy error.log.txt | Raw backend error log file |
| **log.txt** | Copy ▸ Copy log.txt | Raw backend full log file |
| **Full Bundle (ZIP)** | Download ▸ Full Bundle | ZIP containing report.md + error.log.txt + log.txt |

---

## Copy Button — Split Button Pattern

The Copy button uses a **Split Button** pattern:
- **Main area** (left): Copies the **Compact Report** instantly (no API call)
- **Arrow/Chevron** (right): Opens a dropdown with all copy options

```
┌──────────┬───┐
│ 📋 Copy  │ ▼ │   ← Main click = Compact Report (instant, no API call)
└──────────┴───┘
               │
               ├── Copy Compact Report    → generateCompactReport() output
               ├── Copy Full Report       → generateErrorReport() output
               ├── Copy with Backend Logs → generateErrorReport() + error.log.txt (async, fetches API)
               ├── ─────────────────
               ├── Copy error.log.txt     → Raw error log from GET /api/v1/logs/error
               └── Copy log.txt           → Raw full log from GET /api/v1/logs/full
```

## Download Menu Structure

```
[▼ Download]
├── Full Bundle (ZIP)         → POST /api/v1/errors/bundle
├── ─────────────────
├── error.log.txt             → GET /api/v1/logs/error
├── log.txt (Full)            → GET /api/v1/logs/full
├── ─────────────────
└── Report (.md)              → generateErrorReport() as .md file
```

---

## Critical Note: Delegated Server Info in Compact Report

The **Compact Report** (the default copy format) **MUST** include the `Delegated Server Info` section when available. This section is built from `CapturedError.envelopeErrors.DelegatedRequestServer` — no API call is needed. It includes:

- Delegated endpoint URL
- HTTP method and status code
- PHP/downstream stack trace
- Request body (if POST/PUT/PATCH)
- Additional error messages from the downstream server

This is essential for debugging proxy-chain errors (React → Go → WordPress/PHP). Without it, the recipient only sees the Go backend error and cannot diagnose the root cause on the delegated server.

See [01-compact-report.md § Backend error.log.txt Section](./01-compact-report.md#backend-errorlogtxt-section-built-from-capturederror) for field mapping.

---

## Document Inventory

| File |
|------|
| 99-consistency-report.md |


## Cross-References

- [Error Modal Spec](../03-error-modal-reference.md) — Full modal structure and component hierarchy
- [React Components Reference](../02-react-components.md) — Portable React code and component props
- [Envelope Schema](../../05-response-envelope/envelope.schema.json) — JSON Schema source of truth
- [Envelope Error Sample](../../05-response-envelope/envelope-error.json) — Canonical error response
- [Error Handling Spec](../../01-error-handling-reference.md) — Cross-stack error architecture
- [Session Logging Spec](../../07-logging-and-diagnostics/02-session-based-logging.md) — Session data model
- [React Execution Logger](../../07-logging-and-diagnostics/01-react-execution-logger.md) — Frontend call chain

---

*Copy format index — updated: 2026-03-31*
