# Error Management Specification

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Purpose

Consolidated error management specification covering error resolution/debugging, cross-stack error architecture, and the error code registry. This folder is the **single canonical location** for all error management documentation.

---

## Keywords

`error-management` · `error-resolution` · `debugging` · `error-handling` · `error-codes` · `registry` · `apperror` · `response-envelope` · `error-modal` · `diagnostics` · `stack-trace`

---

## Scoring

| Metric | Value |
|--------|-------|
| AI Confidence | Production-Ready |
| Ambiguity | None |
| Health Score | 100/100 (A+) |

---

## Categories

| # | Category | Description | Files |
|---|----------|-------------|-------|
| 01 | [Error Resolution](./01-error-resolution/00-overview.md) | Debugging guides, retrospectives, verification patterns, cheat sheet, cross-reference diagram | 14 |
| 02 | [Error Architecture](./02-error-architecture/00-overview.md) | Cross-stack 3-tier error handling, error modal, response envelope, apperror package, logging, notifications | 22 |
| 03 | [Error Code Registry](./03-error-code-registry/00-overview.md) | Master registry, integration guide, schemas, scripts, templates, collision resolution, utilization report | 18 |

> 📖 **Quick onboarding?** See [structure.md](./structure.md) for a full visual tree with role-based entry points.

---

## Core Principles

### 1. Never Assume — Always Verify

Before claiming any API endpoint works, verify **both directions**:

| Direction | Verification | Example |
|-----------|--------------|---------|
| **Backend** | Test actual endpoint response | `curl http://localhost:8080/api/v1/health \| jq .` |
| **Frontend** | Check detection logic | What conditions trigger "connected" vs "disconnected"? |

### 2. Response Format Standardization

All backend APIs MUST return the Universal Response Envelope (see [02-error-architecture/05-response-envelope/](./02-error-architecture/05-response-envelope/00-overview.md)):

```json
{
  "Status": { "IsSuccess": true, "Code": 200, "Message": "OK" },
  "Attributes": { "RequestedAt": "..." },
  "Results": [{ "..." }]
}
```

### 3. HTTP Status as Primary Indicator

Frontend detection logic MUST use HTTP status codes (2xx) as the primary indicator, NOT response body fields.

### 4. Structured Error Architecture

All errors use the three-tier architecture documented in [02-error-architecture/01-error-handling-reference.md](./02-error-architecture/01-error-handling-reference.md):
- **Tier 1:** Delegated Server (PHP/other) — structured error responses
- **Tier 2:** Go Backend — `apperror` package with stack traces
- **Tier 3:** Frontend — Error store, Global Error Modal

---

## Quick Reference: Common Pitfalls

| Symptom | Likely Cause | Check |
|---------|--------------|-------|
| "Backend disconnected" but backend running | Response format mismatch | Compare handler output to frontend detection logic |
| 404 on API base URL | No index route registered | Check router for `GET /api/v1` handler |
| VITE_API_URL shows wrong value | Resolved vs raw env confusion | Distinguish raw env var from resolved origin |
| HTML instead of JSON | SPA fallback serving index.html | Check if route exists in backend router |
| CORS errors | Missing CORS headers | Check backend CORS middleware configuration |
| 401/403 on protected routes | Token not sent or expired | Check Authorization header, token validity |

---

## Migration Note

This folder consolidates content previously located at:

| Old Location | Status |
|-------------|--------|

---

## Document Inventory

| File |
|------|
| 97-acceptance-criteria.md |
| 98-changelog.md |
| 99-consistency-report.md |


## Cross-References

| Reference | Location |
|-----------|----------|
| Coding Guidelines | [../02-coding-guidelines/00-overview.md](../02-coding-guidelines/00-overview.md) |
| Rust Error Handling | [../02-coding-guidelines/05-rust/02-error-handling.md](../02-coding-guidelines/05-rust/02-error-handling.md) |
| Cross-Language Guidelines | [../02-coding-guidelines/01-cross-language/00-overview.md](../02-coding-guidelines/01-cross-language/00-overview.md) |
| Database Conventions | [../04-database-conventions/00-overview.md](../04-database-conventions/00-overview.md) |
| [structure.md](./structure.md) | Full visual tree |

---

*This specification is mandatory for all projects and is the **highest priority** — error handling must be implemented from the very first line of code. Violations result in debugging time waste.*

---

## Verification

_Auto-generated section — see `spec/03-error-manage/97-acceptance-criteria.md` for the full criteria index._

### AC-ERR-000: Error-management conformance: Overview

**Given** Audit error-handling sites for use of the `apperror` package, error codes, and explicit file/path logging.  
**When** Run the verification command shown below.  
**Then** Every error site uses `apperror.Wrap`/`apperror.New` with a registered code; no bare `errors.New` or swallowed errors remain.

**Verification command:**

```bash
python3 linter-scripts/check-forbidden-strings.py && go run linter-scripts/validate-guidelines.go --path spec --max-lines 15
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
