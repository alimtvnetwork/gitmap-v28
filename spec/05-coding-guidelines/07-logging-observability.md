# Logging & Observability — Diagnostic Output Standards

## Overview

Standards for structured logging, log levels, correlation IDs, and
sensitive data redaction across CLI and application code. Logs are for
operators — make them useful, consistent, and actionable.

---

## 1. Log Levels

Use a clear hierarchy. Never log everything at the same level.

| Level | Purpose | Example |
|-------|---------|---------|
| Error | Operation failed, action required | `"failed to write metadata: permission denied"` |
| Warn | Degraded but recoverable | `"legacy directory found, migrating"` |
| Info | Key lifecycle events | `"scan complete: 42 repos found"` |
| Debug | Detailed internals (verbose only) | `"checking path: /home/user/projects/api"` |
| Trace | Ultra-fine-grained (dev only) | `"entering parseRecord: line 47"` |

### Rules

- Default production level: **Info**.
- Debug/Trace require explicit opt-in (`--verbose`, `LOG_LEVEL=debug`).
- Never log at Error for expected conditions (e.g., "no results found").
- Warn must be actionable — if no one will act on it, it's Debug.
- A single request should produce ≤5 Info lines under normal flow.

---

## 2. Structured Log Format

Every log entry must include context. Never log bare strings.

### Go

```go
verbose.Log("scan", "discovered repo: %s (branch: %s)", repoName, branch)
verbose.Log("release", "version resolved: %s (source: %s)", version, source)
verbose.Log("clone", "cloned %d/%d repos in %s", completed, total, elapsed)
```

### TypeScript

```ts
logger.info({ component: "scan", repo: repoName, branch }, "discovered repo");
logger.error({ component: "release", version, statusCode, err }, "upload failed");
```

### Required Fields

| Field | When |
|-------|------|
| Component/stage | Always — identifies the subsystem |
| Entity name | When operating on a specific item |
| Count/progress | When processing a batch |
| Duration | When timing an operation |
| Error detail | On failures — include the original error |
| Correlation ID | On all entries within a request/operation |

### Output Format

- **CLI tools**: `[HH:MM:SS] [component] message key=value` to log files.
- **Services/APIs**: JSON to stdout — one object per line.
- **Never** mix formats within the same output stream.

```json
{
  "timestamp": "2025-01-15T14:30:22.123Z",
  "level": "info",
  "component": "scan",
  "correlation_id": "req-a1b2c3",
  "message": "scan complete",
  "repos_found": 42,
  "duration_ms": 1200
}
```

---

## 3. Correlation IDs

Every operation that spans multiple stages or services must carry a
correlation ID for end-to-end tracing.

### Generation

```go
// Go — generate at operation entry point
correlationID := fmt.Sprintf("op-%s", uuid.New().String()[:8])
ctx := context.WithValue(parentCtx, ctxKeyCorrelation, correlationID)
```

```ts
// TypeScript — generate at request boundary
const correlationId = `req-${crypto.randomUUID().slice(0, 8)}`;
req.headers["x-correlation-id"] = correlationId;
```

### Propagation Rules

| Boundary | Mechanism |
|----------|-----------|
| HTTP requests | `X-Correlation-ID` header |
| CLI pipeline stages | Context value / struct field |
| Background jobs | Job metadata field |
| Log entries | Always included as `correlation_id` |

### Rules

- Generate at the **outermost** entry point only — never mid-pipeline.
- Accept incoming correlation IDs from callers; generate only if absent.
- Use short IDs (8 chars) for CLI; full UUIDs for distributed services.
- Log the correlation ID in the **first** and **last** log line of every operation.
- Never reuse correlation IDs across independent operations.

### Go Context Pattern

```go
func scanRepos(ctx context.Context, root string) error {
    cid := correlationFrom(ctx)
    verbose.Log("scan", "[%s] starting directory walk: %s", cid, root)
    // ... work ...
    verbose.Log("scan", "[%s] complete: %d repos in %s", cid, count, elapsed)
    return nil
}
```

---

## 4. Sensitive Data Redaction

Secrets, credentials, and PII must never appear in logs, diagnostics,
or error messages — at any log level.

### Must Redact

| Category | Examples |
|----------|----------|
| Credentials | API keys, tokens, passwords, SSH keys |
| PII | Email addresses, names, IP addresses (in privacy contexts) |
| Financial | Card numbers, account IDs |
| Internal paths | Full server filesystem paths (use relative) |

### Redaction Patterns

```go
// Go — redact function
func redact(s string) string {
    if len(s) <= 4 {
        return "****"
    }
    return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}

// Usage
verbose.Log("auth", "token: %s (source: %s)", redact(token), source)
// Output: "token: gh**********3f (source: env)"
```

```ts
// TypeScript — redact utility
function redact(value: string): string {
  if (value.length <= 4) return "****";
  return value.slice(0, 2) + "*".repeat(value.length - 4) + value.slice(-2);
}

// Usage
logger.info({ token: redact(apiKey), source: "env" }, "authenticated");
```

### Automated Detection

Build a pre-log scanner that catches common secret patterns before they
reach the log output.

```go
var sensitivePatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)(ghp_|gho_|github_pat_)[a-zA-Z0-9_]+`),  // GitHub tokens
    regexp.MustCompile(`(?i)(sk-|pk_live_|pk_test_)[a-zA-Z0-9]+`),   // API keys
    regexp.MustCompile(`(?i)bearer\s+[a-zA-Z0-9\-._~+/]+=*`),       // Bearer tokens
    regexp.MustCompile(`(?i)password\s*[:=]\s*\S+`),                  // Password assignments
}

func containsSensitive(msg string) bool {
    for _, pat := range sensitivePatterns {
        if pat.MatchString(msg) {
            return true
        }
    }
    return false
}
```

```ts
const SENSITIVE_PATTERNS = [
  /(?:ghp_|gho_|github_pat_)[a-zA-Z0-9_]+/gi,
  /(?:sk-|pk_live_|pk_test_)[a-zA-Z0-9]+/gi,
  /bearer\s+[a-zA-Z0-9\-._~+/]+=*/gi,
  /password\s*[:=]\s*\S+/gi,
];

function containsSensitive(msg: string): boolean {
  return SENSITIVE_PATTERNS.some((pat) => pat.test(msg));
}
```

### Rules

- Redact **before** the value reaches any log call — not after.
- Never log raw HTTP request/response bodies containing auth headers.
- Environment variable **names** may be logged; **values** must not.
- In error messages, reference secrets by source (`"env:GITHUB_TOKEN"`) not value.
- Automated pattern detection must run in CI as a lint check.

---

## 5. Verbose Mode Pattern

Verbose output is opt-in via a global `--verbose` flag and writes
to timestamped files while showing summaries on stderr.

### Architecture

```
User runs: gitmap-v28 scan --verbose
  │
  ├─ stderr: colored summary lines (always visible)
  │    "✓ Scanned 42 repos in 1.2s"
  │
  └─ file: .gitmap/output/scan-2025-01-15-143022.log
       [14:30:22] [scan] [op-a1b2c3] starting directory walk: /home/user
       [14:30:22] [scan] [op-a1b2c3] found .git: /home/user/api/.git
       [14:30:23] [scan] [op-a1b2c3] remote: https://github.com/user/api.git
       ...
```

### Rules

- Verbose logs go to files, not stdout.
- Stderr shows colored summaries regardless of verbose flag.
- Stdout is reserved for machine-readable output (JSON, CSV).
- Log file names include timestamp to prevent overwrites.
- Each log line has a timestamp, stage prefix, and correlation ID.

---

## 6. Pipeline Stage Logging

For multi-stage operations, log entry and exit of each stage.

```go
verbose.Log("stage", "[%s] starting: %s", cid, stageName)
// ... work ...
verbose.Log("stage", "[%s] completed: %s (%d items, %s)", cid, stageName, count, elapsed)
```

### Standard Stages (Release Example)

| # | Stage | Key Log Fields |
|---|-------|----------------|
| 1 | Version Resolution | source, resolved version |
| 2 | Git Operations | branch, commit, tag |
| 3 | Asset Collection | file paths, sizes |
| 4 | Cross-Compilation | GOOS, GOARCH, output path |
| 5 | Compression | algorithm, ratio, bytes |
| 6 | Upload | URL, status code, retry count |
| 7 | Metadata Persistence | file path, JSON size |

---

## 7. Diagnostic Output (Doctor Pattern)

Health checks follow a consistent pass/fail format.

```
✓ Config file       OK — parsed 6 fields
✓ Database          OK — 42 repos, 12 releases
✓ Legacy dirs       OK — no migration needed
✗ Git binary        FAIL — git not found in PATH
✓ Deploy path       OK — E:\bin-run\gitmap-v28\
```

### Rules

- One line per check — icon, name, status, detail.
- Use `✓` for pass, `✗` for fail, `⚠` for warning.
- Include the `--fix-path` pattern for auto-remediation.
- Exit code reflects overall health (0 = all pass, 1 = any fail).
- Constants for all check names and messages — no magic strings.

---

## 8. Progress Reporting

Use `[current/total]` counters on stderr for batch operations.

```
[1/42] Scanning api-gateway...
[2/42] Scanning frontend-app...
...
[42/42] Scanning legacy-tool...
✓ Scanned 42 repos in 3.4s
```

### Rules

- Counter format: `[current/total]` — always padded.
- Show item name being processed.
- Print summary line at completion.
- Support `--quiet` to suppress progress (keep summary).
- Never mix progress output with stdout data.

---

## 9. Error Logging

### Context Wrapping

Always wrap errors with the operation that failed.

```go
// ✅ Correct
return fmt.Errorf("writing release metadata v%s: %w", version, err)

// ❌ Wrong
return err
return fmt.Errorf("error: %w", err)
```

### User-Facing vs Internal

| Audience | Format | Example |
|----------|--------|---------|
| User (stderr) | Plain, actionable | `"Error: config.json not found. Run 'gitmap-v28 doctor' to diagnose."` |
| Log file | Detailed, technical | `"[config] ReadFile failed: /home/user/data/config.json: ENOENT"` |

---

## 10. What NOT to Log

- Secrets, tokens, or credentials — never (see §4).
- Redundant success messages for trivial operations.
- Stack traces in user-facing output (log file only).
- Raw HTTP response bodies (log truncated previews).
- PII without explicit user consent and redaction.

---

## Constraints

- Every log call must include a component/stage prefix.
- Correlation IDs are mandatory for operations spanning 2+ stages.
- Sensitive data redaction must be enforced at the logger level.
- Log format must not change between verbose and quiet modes — only volume.
- CI must include a lint step that scans for leaked secret patterns.

---

## References

- Verbose Logging Spec: `spec/04-generic-cli/16-verbose-logging.md`
- Progress Tracking: `spec/04-generic-cli/17-progress-tracking.md`
- Error Handling: `spec/05-coding-guidelines/04-error-handling.md`
- Security & Secrets: `spec/05-coding-guidelines/08-security-secrets.md`
- Code Quality: `spec/05-coding-guidelines/01-code-quality-improvement.md`
