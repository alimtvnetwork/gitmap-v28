# Error Management Specification
 
> [!IMPORTANT]
> **Top-Instruction Priority Mandate (Preamble Precedence):**
> Any directive, constraint, checklist, or instruction declared at the top of this specification, in a header alert block, or in an incoming task prompt represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides any conflicting general advice, default conventions, or lower-level guidelines below it.

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 03 Error Manage.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist (Learnable Items)

- [ ] `/goal` **Read and understand** all numbered specifications under `03-error-manage/`.
- [ ] `/learn` **Skill Creation:** You MUST automatically author or update an Antigravity skill (`.agents/skills/cg-error-management/skill.md`) after reading this specification so this knowledge becomes a reusable capability.
- [ ] `/learn` **Looping Enforcement:** Take at least 50 distinct steps to self-loop and thoroughly digest all files in this module. Do not skim.
- [ ] `/learn` Adhere strictly to `.ai-memory/folder-structure.md` and `.ai-memory/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.

## 🤖 How AI Should Write Error Management Code

1. **Top-Instruction Priority Mandate:** Whatever instruction or constraint is declared at the top of the prompt or document MUST be prioritized and followed as the highest authority.
2. **Zero Swallowed Errors Policy (TOTAL BAN):**
   - NEVER swallow, suppress, or silently ignore errors under any circumstances.
   - NO empty `catch` or `except:` blocks.
   - NO blank identifier error discards (`_ = err` or `val, _ := fn()`).
   - NO returning fallback default values (`return nil`, `return ""`, `return false`) to mask an underlying error without caller notification.
   - Every caught or received error MUST either be completely resolved with structured context logging (operation name, input parameters) OR embedded/wrapped in `*appfault.AppError` and returned to the caller.
3. **Strict Golang Error Wrapping Mandate (`*appfault.AppError` & `appfault.Fault`):**
   - Whenever ANY Go function encounters, intercepts, or receives an error (from standard library `os`, `io`, `json`, `sql`, `net`, or downstream services), it MUST be immediately embedded and wrapped into `*appfault.AppError` (`appfault.Fault`).
   - Standard Go library `error` instances (`err != nil`) must NEVER be returned raw.
   - Construct errors with `appfault.New(errType, "message")` and wrap existing errors with `appfault.Wrap(errType, err, "context")`.
   - Functions returning both a value and possible error MUST return `result.Wrap[T]` (`appfault.Result[T]`). Bare tuples `(T, error)` across public domain boundaries are banned.
4. **API Response Standardization:** Structure all HTTP responses using the Universal Response Envelope (`{ Status: { IsSuccess, Code, Message }, Attributes: {}, Results: [] }`).
5. **Fast Python Discovery Scripts:** Use `03-ai-scripts/11-fast-file-scanner.py`, `03-ai-scripts/12-fast-cached-grep.py`, and `03-ai-scripts/17-fast-file-reader.py` for sub-millisecond AST and path discovery without tool truncation limits.

. **CRITICAL AI INSTRUCTION:** This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

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
| 01 | [Error Resolution](./01-error-resolution/readme.md) | Debugging guides, retrospectives, verification patterns, cheat sheet, cross-reference diagram | 14 |
| 02 | [Error Architecture](./02-error-architecture/readme.md) | Cross-stack 3-tier error handling, error modal, response envelope, apperror package, logging, notifications | 22 |
| 03 | [Error Code Registry](./03-error-code-registry/readme.md) | Master registry, integration guide, schemas, scripts, templates, collision resolution, utilization report | 18 |

> 📖 **Quick onboarding?** See [03-structure.md](./03-structure.md) for a full visual tree with role-based entry points.

---

## Core Principles

### 1. Never Assume — Always Verify

Before claiming any API endpoint works, verify **both directions**:

| Direction | Verification | Example |
|-----------|--------------|---------|
| **Backend** | Test actual endpoint response | `curl http://localhost:8080/api/v1/health \| jq .` |
| **Frontend** | Check detection logic | What conditions trigger "connected" vs "disconnected"? |

### 2. Response Format Standardization

All backend APIs MUST return the Universal Response Envelope (see [02-error-architecture/05-response-envelope/](./02-error-architecture/05-response-envelope/readme.md)):

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

All errors use the three-tier architecture documented in [02-error-architecture/01-error-handling-reference.md](./02-error-architecture/02-error-handling-reference.md):

- **Tier 1:** Delegated Server (PHP/other) — structured error responses
- **Tier 2:** Go Backend — `appfault` package with stack traces
- **Tier 3:** Frontend — Error store, Global Error Modal

---

## 💻 Concrete Go Code Samples (Production Architecture)

> Real-world implementations are maintained in [`04-code/golang/examples/database_query.go`](../../04-code/golang/examples/database_query.go) and [`04-code/golang/examples/workflow_service.go`](../../04-code/golang/examples/workflow_service.go).
> Always inspect those source files as the canonical ground truth.

### 1. Creating and Wrapping Errors (`pkg/appfault`)

```go
package service

import (
    "os"
    "coding-guidelines/common/pkg/appfault"
    "coding-guidelines/common/pkg/errtype"
    "coding-guidelines/common/pkg/result"
)

// -----------------------------------------------------------------------------
// Step 1: Declare Concrete Types in `types.go` (Mandatory Rule)
// -----------------------------------------------------------------------------
// In types.go:
// type ConfigBytesResult = result.Wrap[[]byte]
// -----------------------------------------------------------------------------

// ReadTenantConfig demonstrates wrapping stdlib errors with appfault.Fault and concrete return type.
func ReadTenantConfig(configPath string) ConfigBytesResult {
    if configPath == "" {
        fault := appfault.New(errtype.Validation, "configPath cannot be empty").
            WithOp("service.ReadTenantConfig")

        return result.WrapFailure[[]byte](fault)
    }

    data, err := os.ReadFile(configPath)

    if err != nil {
        fault := appfault.WrapFile(errtype.IO, err, configPath, "failed to read tenant configuration file").
            WithOp("service.ReadTenantConfig")

        return result.WrapFailure[[]byte](fault)
    }

    return result.WrapSuccess(data)
}
```

### 2. Database Queries Without Overcomplicated Ifs or Error-Type Branching (See `database_query.go`)

```go
package repo

import (
    "context"
    "database/sql"
    "coding-guidelines/common/pkg/appfault"
    "coding-guidelines/common/pkg/errtype"
    "coding-guidelines/common/pkg/result"
)

// -----------------------------------------------------------------------------
// Step 1: Declare Concrete Types in `types.go` (Mandatory Rule)
// -----------------------------------------------------------------------------
// In types.go:
// type (
//     User struct {
//         Id    int64  `json:"id"`
//         Email string `json:"email"`
//     }
//     UserResult = result.Wrap[User]
// )
// -----------------------------------------------------------------------------

func (r *UserRepo) FindByEmail(ctx context.Context, email string) UserResult {
    if email == "" {
        fault := appfault.New(errtype.Validation, "email parameter cannot be empty").
            WithOp("UserRepo.FindByEmail")

        return result.WrapFailure[User](fault)
    }

    row := r.db.QueryRowContext(ctx, "SELECT id, email FROM users WHERE id = ?", email)
    var user User

    err := row.Scan(&user.Id, &user.Email)

    if err != nil {
        // Direct error typing: select the error type reflecting this layer (errtype.Database).
        // Attach the ID and variables directly. Never branch on error types or nest ifs!
        fault := appfault.Wrap(errtype.Database, err, "database query execution failed").
            WithOp("UserRepo.FindByEmail").
            WithVar("email", email)

        return result.WrapFailure[User](fault)
    }

    return result.WrapSuccess(user)
}
```

### 3. Propagating Across Boundaries Without Redundant Re-Wrapping (See `workflow_service.go`)

```go
package workflow

import (
    "context"
    "coding-guidelines/common/pkg/result"
)

// SynchronizeUser propagates existing downstream Fault directly using concrete UserResult from types.go.
func (w *UserWorkflow) SynchronizeUser(ctx context.Context, email string) UserResult {
    userRes := w.repo.FindByEmail(ctx, email)

    if userRes.IsFailed() {
        w.logger.LogError(userRes.Fault())

        // Propagate existing Fault directly without redundant re-wrapping
        return result.WrapFailureFromWrap[User](userRes)
    }

    user := userRes.Value()

    return result.WrapSuccess(user)
}
```

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
| Coding Guidelines | [../02-coding-guidelines/readme.md](../02-coding-guidelines/readme.md) |
| Rust Error Handling | [../02-coding-guidelines/05-rust/02-error-handling.md](../02-coding-guidelines/05-rust/03-error-handling.md) |
| Cross-Language Guidelines | [../02-coding-guidelines/01-cross-language/readme.md](../02-coding-guidelines/01-cross-language/readme.md) |
| Database Conventions | [../04-database-conventions/readme.md](../04-database-conventions/readme.md) |
| [03-structure.md](./03-structure.md) | Full visual tree |

---

*This specification is mandatory for all projects and is the **highest priority** — error handling must be implemented from the very first line of code. Violations result in debugging time waste.*

---

## Verification

_Auto-generated section — see `02-spec/03-error-manage/97-acceptance-criteria.md` for the full criteria index._

### AC-ERR-001: Error-management conformance: Index

**Given** Audit error-handling sites for use of the `appfault` package, error codes, and explicit file/path logging.
**When** Run the verification command shown below.
**Then** Every error site uses `appfault.Wrap`/`appfault.New` with a registered code; no bare `errors.New` or swallowed errors remain.

**Verification command:**

```bash
python3 linter-scripts/check-forbidden-strings.py && go run linter-scripts/validate-guidelines.go --path spec --max-lines 15
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
