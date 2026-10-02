---
name: cg-string-normalization-and-equalfoldany
description: Autonomously audits, refactors, and standardizes string normalization, case-insensitive comparisons, and chained equality into reusable strutil helpers (EqualFoldAny, EqualFoldAnyTrim) across Go, TypeScript, and Rust.
---

# Skill: String Normalization, `EqualFoldAny` & Centralized Utility Discovery (`cg-string-normalization-and-equalfoldany`)

This skill governs autonomous execution for auditing, refactoring, and standardizing repetitive string normalization (`ToLower`, `TrimSpace`) and chained equality checks (`EqualFold` OR-chains) across Go, TypeScript, and Rust codebases.

## Trigger Keywords & Aliases

- `cg-string-normalization`
- `cg-equalfoldany`
- `strutil-normalization`
- `audit-string-comparisons`
- `cg-execute equalfoldany`

---

## 1. Search First Protocol for `<repo>/pkg/strutil/strutil.go`

Before modifying call sites, creating local helpers, or introducing duplicate logic, agents must strictly follow the **Search First Protocol**:

### 1.1 Canonical Package Locations

| Language | Primary Canonical Path | Permitted Internal Path | Prohibited Ad-Hoc Paths |
| :--- | :--- | :--- | :--- |
| **Go** | `<repo>/pkg/strutil/strutil.go` | `<repo>/internal/strutil/strutil.go` | `cmd/*.go`, `<repo>/pkg/util/helpers.go`, inline file helpers |
| **TypeScript** | `<repo>/src/lib/strutil.ts` | `<repo>/src/utils/strutil.ts` | `src/components/*.tsx`, `src/pages/*.ts`, local helpers |
| **Rust** | `<repo>/src/util/strutil.rs` | `<repo>/src/strutil/mod.rs` | `<repo>/src/main.rs`, inline sub-modules in `<repo>/src/cmd/` |
| **Python** | `<repo>/pkg/strutil/strutil.py` | `<repo>/src/strutil.py` | `scripts/*.py`, inline functions in individual script files |

### 1.2 Protocol Execution Flow

1. **Locate Canonical String Utility:**
   - Scan for string utilities using `gitmap find-files-any "strutil"`.
   - In Go: Check `<repo>/pkg/strutil/strutil.go` or `<repo>/04-code/golang/pkg/strutil/strutil.go`.
   - In TypeScript: Check `<repo>/src/lib/strutil.ts` or `<repo>/src/utils/strutil.ts`.
   - In Rust: Check `<repo>/src/util/strutil.rs` or `crate::strutil`.
2. **Re-use Mandate:**
   - All refactored call sites MUST import and use the canonical utility functions (`strutil.EqualFoldAny`, `strutil.EqualFoldAnyTrim`, `strutil.NormalizeLowerTrim`).
   - Total ban on authoring duplicate unshared functions (`isYes`, `checkConfirm`, `matchesAny`, `cleanEqual`) inside consumer packages.
3. **Absence Protocol:**
   - If the canonical utility is missing in the target repository, agents must create or register it following `02-spec/02-coding-guidelines/01-cross-language/34-string-normalization-and-equalfoldany.md` BEFORE modifying call sites.

---

## 2. Step-by-Step Execution Workflow

Execute audits and refactoring runs in four bounded stages:

### Stage 1: Discovery & Inventory (Steps 1 .. N/4)
- Scan the repository for chained `EqualFold`, `ToLower`, and manual trimming patterns using GitMap streaming searches.
- Identify all call sites violating DRY principles or allocating intermediate strings unnecessarily.
- Populate the Violation Ledger in `.ai-memory/plans/pending/` with exact files, line numbers, and planned transformations.

### Stage 2: Planning & Subtask Generation (Steps N/4+1 .. N/2)
- Group target files into micro-batches of 5–8 files each.
- Generate modular subtask files in `.ai-memory/plans/subtasks/NN-<slug>/01-<title>.md`, `02-<title>.md`, etc.
- Verify each subtask assigns disjoint file boundaries to prevent merge collisions.

### Stage 3: Micro-Batched Surgical Refactoring (Steps N/2+1 .. 3N/4)
- Replace verbose OR-chains with canonical `strutil` helper calls.
- Hoist repeated string normalizations outside tight loops.
- Enforce positive boolean variables (`isMatch`, `hasTargetRole`) and implicit boolean checks (`if isMatch { ... }`).
- Maintain function lengths <= 8–15 lines by decomposing complex blocks.

### Stage 4: Verification & Consolidation (Steps 3N/4+1 .. N)
- Run targeted file-level linters and verify `exit 0` across modified files.
- Consolidate completed subtasks into `.ai-memory/plans/completed/`.
- Stage all modified files and create a single grouped atomic commit via GitMap (`gitmap cpf "<module> - <summary>"`).

---

## 3. Polyglot Code Transformation Catalog

### 3.1 Go (`<repo>/pkg/strutil/strutil.go`)

#### ❌ Anti-Pattern: Chained `strings.EqualFold` & Manual Trimming

```go
// ❌ WRONG — manual trim, intermediate variable, chained EqualFold calls
trimmedRole := strings.TrimSpace(userRole)
if strings.EqualFold(trimmedRole, "admin") || strings.EqualFold(trimmedRole, "owner") || strings.EqualFold(trimmedRole, "moderator") {
    grantPermission()
}
```

#### ✅ Canonical Pattern: `strutil.EqualFoldAnyTrim`

```go
// ✅ REQUIRED — single-line atomic check via canonical strutil helper
isElevated := strutil.EqualFoldAnyTrim(userRole, "admin", "owner", "moderator")
if isElevated {
    grantPermission()
}
```

#### ❌ Anti-Pattern: Allocating `ToLower` in Loops

```go
// ❌ WRONG — allocates heap copy for every iteration
for _, item := range rawItems {
    if strings.ToLower(item) == "active" || strings.ToLower(item) == "pending" {
        results = append(results, item)
    }
}
```

#### ✅ Canonical Pattern: Zero-Allocation `EqualFoldAny`

```go
// ✅ REQUIRED — zero-allocation Unicode case-folding match
for _, item := range rawItems {
    isTargetStatus := strutil.EqualFoldAny(item, "active", "pending")
    if isTargetStatus {
        results = append(results, item)
    }
}
```

#### Real-World CLI Confirmation Example (`cli/cmd/releaseundo.go`)

```go
// ❌ BAD:
func confirmUndoRelease(tag string) bool {
    fmt.Printf("Delete %s locally and on origin? [y/N]: ", tag)
    var reply string
    _, _ = fmt.Scanln(&reply)
    trimmed := strings.TrimSpace(reply)

    return strings.EqualFold(trimmed, "y") || strings.EqualFold(trimmed, "yes")
}

// ✅ GOOD:
func confirmUndoRelease(tag string) bool {
    fmt.Printf("Delete %s locally and on origin? [y/N]: ", tag)
    var reply string
    _, _ = fmt.Scanln(&reply)

    return strutil.EqualFoldAnyTrim(reply, "y", "yes")
}
```

---

### 3.2 TypeScript (`<repo>/src/lib/strutil.ts`)

#### ❌ Anti-Pattern: Verbose Chained `.toLowerCase()` and `.trim()`

```typescript
// ❌ WRONG — repetitive method chaining and fragile equality checks
const clean = input ? input.trim().toLowerCase() : '';
if (clean === 'yes' || clean === 'y' || clean === 'true' || clean === '1') {
  enableFeature();
}
```

#### ✅ Canonical Pattern: `strutil.equalFoldAnyTrim`

```typescript
// ✅ REQUIRED — concise variadic helper invocation
const isAffirmative = strutil.equalFoldAnyTrim(input, 'yes', 'y', 'true', '1');
if (isAffirmative) {
  enableFeature();
}
```

---

### 3.3 Rust (`<repo>/src/util/strutil.rs`)

#### ❌ Anti-Pattern: Chained `eq_ignore_ascii_case`

```rust
// ❌ WRONG — repetitive trimming and chained ASCII case comparison
let val = input.trim();
if val.eq_ignore_ascii_case("draft") || val.eq_ignore_ascii_case("review") {
    process_document();
}
```

#### ✅ Canonical Pattern: `strutil::equal_fold_any_trim`

```rust
// ✅ REQUIRED — slice-based case-folding helper
let is_pending = strutil::equal_fold_any_trim(input, &["draft", "review"]);
if is_pending {
    process_document();
}
```

---

## 4. High-Speed AST Search & Discovery via GitMap

Use GitMap multi-core live streaming search commands to identify candidates across languages. Generic shell search tools (`Select-String`, `git grep`, `grep`, `findstr`) are strictly banned:

```bash
# Go: Find chained EqualFold checks
gitmap aum search -r "EqualFold\(.*\|\|.*EqualFold\(" [dir] -e .go

# Go: Find TrimSpace followed by EqualFold
gitmap aum search -r "TrimSpace\(.*EqualFold" [dir] -e .go

# Go: Find ToLower equality OR-chains
gitmap aum search -r "strings\.ToLower\(.*==.*\|\|" [dir] -e .go

# Go: Find all EqualFold call sites
gitmap aum search "strings.EqualFold" [dir] -e .go

# TypeScript: Find chained .toLowerCase() checks
gitmap aum search -r "\.toLowerCase\(\)\s*===\s*.*\|\|" [dir] -e .ts

# TypeScript: Find .trim().toLowerCase() calls
gitmap aum search -r "\.trim\(\)\.toLowerCase\(\)" [dir] -e .ts

# Rust: Find chained eq_ignore_ascii_case checks
gitmap aum search -r "eq_ignore_ascii_case\(.*\|\|" [dir] -e .rs
```

---

## 5. Verification Checklist & Quality Gates

- [ ] Canonical `<repo>/pkg/strutil/strutil.go` verified before refactoring call sites.
- [ ] Chained `EqualFold` and `ToLower` OR-chains converted to `strutil.EqualFoldAny` or `strutil.EqualFoldAnyTrim`.
- [ ] Unshared ad-hoc helper functions (`isYes`, `checkConfirm`, `matchesAny`) eliminated.
- [ ] Affirmative boolean variable naming enforced (`is`, `has` only).
- [ ] Implicit boolean checks enforced (`if isMatch { ... }`, never `if isMatch == true`).
- [ ] No mixed polarity conditions (`if isA && !isB` is banned).
- [ ] Vertical line gaps enforced (blank line before `if`, after `}`, before `return`).
- [ ] Functions adhere to canonical size limits (<= 8–15 lines).
- [ ] Targeted linters exit 0:
  - `python linter-scripts/check-relative-paths.py`
  - `python 03-ai-scripts/21-sequence-integrity-linter.py 01-prompts/15-cg-execute`
  - `python linter-scripts/check-forbidden-strings.py`
- [ ] ZERO test runs or build commands executed during routine execution (R1).

---

## 6. Cross-References

- [34-string-normalization-and-equalfoldany.md](02-spec/02-coding-guidelines/01-cross-language/34-string-normalization-and-equalfoldany.md)
- [01-architecture-spec.md](02-spec/21-app/02-string-normalization-and-equalfoldany/01-architecture-spec.md)
- [02-prompt-and-skill-spec.md](02-spec/21-app/02-string-normalization-and-equalfoldany/02-prompt-and-skill-spec.md)
- [37-string-normalization-and-equalfoldany.md](01-prompts/15-cg-execute/37-string-normalization-and-equalfoldany.md)
