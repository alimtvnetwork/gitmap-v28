# Golang Coding Standards — Typed constants, enums, DRY enforcement (AI Execution Prompt)

> **/goal** Eliminate all magic strings and numbers in Go codebases by enforcing byte-based enums, typed const blocks, and DRY abstractions across errors, schemas, and API clients.
> **/learn** Internalize the Go enum standard (byte representation, iota, variantLabels, Type suffix), Result[T] client returns, and reusable validation methods.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace all magic strings, HTTP methods, status codes, and config keys with typed constants and byte-based enums.
- [ ] `/learn` Never return tuple `(*T, error)` from API client helper functions; return `apperror.Result[T]` with structured error context.
- [ ] `/goal` Enforce DRY principles across duplicate error handling, validation, JSON access, and database query patterns.
- [ ] `/learn` Ensure all enum variants define PascalCase `variantLabels` and conform to the canonical Go enum specification.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

> **Parent:** [Golang Coding Standards](./readme.md)
> **Version:** 3.7.0
> **Updated:** 2026-03-31

---

## Typed Constants & Enums

> **Canonical source:** [Go Enum Specification](../01-enum-specification/readme.md) — core pattern, required methods, folder structure
> **Cross-language reference:** [Enum Naming Quick Reference](../../06-ai-optimization/07-enum-naming-quick-reference.md) — Go, TypeScript, PHP comparison

All Go enum rules (byte type, `Invalid` zero value, `iota`, `variantLabels`, required methods, folder structure) are defined in the [Go Enum Specification](../01-enum-specification/readme.md). Do not duplicate here.

### Zero Magic Strings/Numbers

- All HTTP status codes → typed constants
- All error codes → `apperror` code constants
- All config keys → typed const block
- All status/event strings → typed byte-based enum constants
- All HTTP methods → `httpmethodtype.Variant` enum (see [HTTP Method Enum](../03-httpmethod-enum.md))
- All API operation names → per-domain operation enum

### API Client — No Tuple Returns, No Magic Strings

API client helper functions MUST return `apperror.Result[T]`, not `(*T, error)`. All struct literal fields for method and operation MUST use enum constants.

```go
// ❌ FORBIDDEN: magic strings + tuple return
func (c *Client) CleanupSnapshots(opts SnapshotCleanupOptions) (*SnapshotCleanupResult, error) {
    callInput := apiCallInput{
        Method:    "POST",
        Operation: "snapshot cleanup",
    }

    return doAPICall[SnapshotCleanupResult](c, callInput)
}

// ✅ REQUIRED: enum constants + Result[T] return
func (c *Client) CleanupSnapshots(opts SnapshotCleanupOptions) apperror.Result[SnapshotCleanupResult] {
    callInput := apiCallInput{
        Method:    httpmethod.Post,
        Operation: snapshotoperationtype.Cleanup,
    }

    return apiCallTo[SnapshotCleanupResult](c, callInput)
}
```

---

## DRY Enforcement

| Pattern | Solution |
|---------|----------|
| Repeated error handling | `apperror.Result[T]` or helper functions |
| Repeated JSON key access | Typed response structs |
| Repeated validation | `Validate()` method on input structs |
| Repeated DB patterns | `dbutil` generic wrappers |
| Repeated string constants | Typed const blocks with `Type` suffix |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/03-golang/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-REF-006: Go Enum Architecture, Constants, and Code Reuse

**Given** Go source code under review or development.
**When** Codebases are audited against Go coding standards.
**Then** All magic literals are replaced with typed enums/constants, API helpers return Result types, and repeated patterns are consolidated with exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/04-golang-standards-reference --check-only
```
**Expected:** exit 0. Zero violations.
