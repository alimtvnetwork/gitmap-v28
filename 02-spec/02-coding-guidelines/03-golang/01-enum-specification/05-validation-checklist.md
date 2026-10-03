# Validation Checklist (AI Execution Prompt)

> **/goal** Provide a comprehensive 50-point compliance audit checklist and grep tooling for verifying Go enum implementations.
> **/learn** Master enum audit scoring, identify legacy string comparisons in switches, detect missing predicate methods, and eliminate hardcoded string constants.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each enum meets the 50-point compliance criteria across structure, declaration, methods, lookup tables, and string elimination.
- [ ] `/learn` Audit switch statements and condition blocks to ensure zero string literals are used for enum variants.
- [ ] `/goal` Verify compile-time validation rules and grep commands detect all legacy enum anti-patterns.
- [ ] `/learn` Ensure zero tolerance for legacy `Unknown` zero values; enforce replacement with `Invalid Variant = iota`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Status:** Complete
**Updated:** 2026-04-16

---

## Audit Process

Use this checklist when auditing a CLI for enum compliance.

---

## Pre-Audit: Find String Literals

Search for hardcoded strings that should be enums:

```bash

# Find string comparisons in switch statements

grep -rn 'case "' --include="*.go" .

# Find string equality checks

grep -rn '== "' --include="*.go" .
grep -rn '!= "' --include="*.go" .

# Find string assignments to type fields

grep -rn 'Type.*=' --include="*.go" .
grep -rn 'Mode.*=' --include="*.go" .
grep -rn 'Status.*=' --include="*.go" .
grep -rn 'Provider.*=' --include="*.go" .
```

---

## Checklist

### 1. Structure (10 points)

| # | Check | Points | Pass |
|---|-------|--------|------|
| 1.1 | Enums in `internal/enums/` directory | 2 | ☐ |
| 1.2 | Each enum in own package folder | 2 | ☐ |
| 1.3 | Main file named `variant.go` | 2 | ☐ |
| 1.4 | Package name is `snake_case` | 2 | ☐ |
| 1.5 | Type named `Variant` | 2 | ☐ |

### 2. Declaration (10 points)

| # | Check | Points | Pass |
|---|-------|--------|------|
| 2.1 | Uses `byte` as underlying type | 2 | ☐ |
| 2.2 | First constant is `Invalid` with `iota` | 2 | ☐ |
| 2.3 | Uses `iota` for all values | 2 | ☐ |
| 2.4 | Constants are PascalCase | 2 | ☐ |
| 2.5 | Each constant has doc comment | 2 | ☐ |

### 3. Required Methods (14 points)

| # | Check | Points | Pass |
|---|-------|--------|------|
| 3.1 | `String() string` implemented | 2 | ☐ |
| 3.2 | `Label() string` implemented | 2 | ☐ |
| 3.3 | `IsValid() bool` implemented | 2 | ☐ |
| 3.4 | `Is{Value}()` for each variant | 2 | ☐ |
| 3.5 | `All() []Variant` implemented | 2 | ☐ |
| 3.6 | `ByIndex(int) Variant` implemented | 2 | ☐ |
| 3.7 | `Parse(string) (Variant, error)` implemented | 2 | ☐ |

### 4. Lookup Tables (6 points)

| # | Check | Points | Pass |
|---|-------|--------|------|
| 4.1 | `variantLabels` array exists | 4 | ☐ |
| 4.2 | Array uses array literal (not slice) | 2 | ☐ |

### 5. No Hardcoded Strings (10 points)

| # | Check | Points | Pass |
|---|-------|--------|------|
| 5.1 | No string literals in switch cases | 3 | ☐ |
| 5.2 | No string comparisons for types | 3 | ☐ |
| 5.3 | Uses `Is{Value}()` for conditions | 2 | ☐ |
| 5.4 | Uses `Parse()` for string input | 2 | ☐ |

---

## Scoring

| Score | Rating | Action |
|-------|--------|--------|
| 45-50 | ✅ Compliant | No action needed |
| 35-44 | ⚠️ Partial | Fix missing methods |
| 25-34 | 🔶 Needs Work | Refactor required |
| 0-24 | ❌ Non-Compliant | Full rewrite needed |

---

## Audit Report Template

```markdown

# Enum Compliance Audit: {CLI Name}

**Date:** YYYY-MM-DD
**Auditor:** AI/Human
**Version:** X.X.X

## Summary

| Category | Score | Max |
|----------|-------|-----|
| Structure | X | 10 |
| Declaration | X | 10 |
| Required Methods | X | 14 |
| Lookup Tables | X | 6 |
| No Hardcoded Strings | X | 10 |
| **Total** | **X** | **50** |

## Enums Found

| Enum | Package | Status |
|------|---------|--------|
| Provider | `internal/enums/provider` | ✅/⚠️/❌ |
| Platform | `internal/enums/platform` | ✅/⚠️/❌ |

## Issues Found

### Critical

1. [Issue description]
   - File: `path/to/file.go`
   - Line: XX
   - Fix: [Description]

### Warnings

1. [Warning description]

## Recommendations

1. [Recommendation]

## Files to Create/Modify

| File | Action | Priority |
|------|--------|----------|
| `internal/enums/providertype/variant.go` | Create | High |
```

---

## Common Issues

### Issue: String-Based Type Field

**Before:**
```go
type SearchRequest struct {
    Provider string `json:"provider"`
}
```

**After:**
```go
import "myapp/internal/enums/providertype"

type SearchRequest struct {
    Provider providertype.Variant
}
```

### Issue: String Switch

**Before:**
```go
switch req.Provider {
case "serpapi":
    // ...
case "colly":
    // ...
}
```

**After:**
```go
switch req.Provider {
case provider.SerpApi:
    // ...
case provider.Colly:
    // ...
default:
    return fmt.Errorf("invalid provider: %s", req.Provider)
}
```

### Issue: String Comparison

**Before:**
```go

if config.Mode == "parallel" {
    // ...
}
```

**After:**
```go

if config.Mode.IsParallel() {
    // ...
}
```

### Issue: Hardcoded Default

**Before:**
```go
mode := "sequential"
```

**After:**
```go
mode := search_mode.Sequential
```

### Issue: Using "Unknown" as Zero Value

**Before:**
```go
const (
    Unknown Variant = iota  // ❌ Ambiguous
)
```

**After:**
```go
const (
    Invalid Variant = iota  // ✅ Clearly invalid
)
```

---

## Grep Commands for Audit

```bash

# Find potential enum candidates

grep -rn 'type.*string' --include="*.go" internal/

# Find switch on strings

grep -rn 'switch.*{' -A5 --include="*.go" . | grep 'case "'

# Find string fields that could be enums

grep -rn 'Type\|Mode\|Status\|Provider\|Format' --include="*.go" internal/models/

# Find missing Is* methods

grep -rn 'if.*==' --include="*.go" . | grep -v 'err\|nil\|0\|false\|true'

# Find legacy "Unknown" zero values (should be "Invalid")

grep -rn 'Unknown.*Variant.*=.*iota' --include="*.go" internal/enums/
```

---

*Validation checklist for enum compliance audits.*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/03-golang/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-GO-ENUM-005: Go Enum Validation and Compile-Time Verification

**Given** Go source code defining domain enums.
**When** Enum implementations are audited against Go enum specifications.
**Then** Full compliance across the 50-point audit checklist, zero string comparisons in switch cases, and exhaustive verification are strictly satisfied with zero violations and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/03-golang/01-enum-specification --check-only
```
**Expected:** exit 0. Zero violations.
