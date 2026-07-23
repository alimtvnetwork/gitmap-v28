# 07 — Logging & Observability

Standards for structured logging, log levels, and diagnostics.

## Log Levels

| Level | Purpose |
|-------|---------|
| Error | Operation failed, action required |
| Warn | Degraded but recoverable |
| Info | Key lifecycle events |
| Debug | Detailed internals (verbose only) |

Default production level: Info. Debug requires `--verbose`.

## Structured Format

Every log entry must include context: component/stage, entity name, count/progress, duration, error detail.

CLI tools: `[HH:MM:SS] [component] message key=value` to log files.

## Correlation IDs

Every multi-stage operation carries a correlation ID. Generate at the outermost entry point. Short IDs (8 chars) for CLI.

## Sensitive Data Redaction

Secrets, credentials, and PII must never appear in logs. Redact before the value reaches any log call.

## Verbose Mode

Verbose logs go to files, not stdout. Stderr shows colored summaries. Stdout reserved for machine-readable output.

## Progress Reporting

Use `[current/total]` counters on stderr for batch operations. Print summary at completion. Support `--quiet`.

## Error Logging

Always wrap errors with operation context. User-facing (stderr): plain and actionable. Log file: detailed and technical.

---

Source: `spec/05-coding-guidelines/07-logging-observability.md`
