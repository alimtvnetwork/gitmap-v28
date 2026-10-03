# String Normalization & EqualFoldAny Specification (AI Execution Prompt)

> **/goal** Eliminate ad-hoc chained string comparisons, repeated trimming, and case-transformation allocations by centralizing all candidate matching behind canonical `strutil.EqualFoldAnyTrim` and `strutil.EqualFoldAny` utilities across Go, TypeScript, Rust, and Python.
> **/learn** Master the Search First protocol to locate existing `strutil` helpers, understand Unicode case folding vs lowercase allocations, short-circuiting candidate iteration, and variadic candidate signatures.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Replace ad-hoc chained equality or `strings.EqualFold()` OR-chains with canonical `strutil.EqualFoldAnyTrim(target, ...candidates)`.
- [ ] `/learn` Enforce Search First protocol: scan `pkg/strutil` for existing helpers before authoring one-off string comparison functions.
- [ ] `/goal` Guarantee zero-allocation upfront trimming and short-circuit evaluation for multi-candidate string matching.
- [ ] `/learn` Verify strict target-first variadic candidate parameter signatures and 100% relative repository paths.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 1.0.0
**Updated:** 2026-10-02
**Applies to:** All languages (Go, TypeScript, Rust, Python)
**Status:** Active
**AI Confidence:** High
**Ambiguity:** None

---

## 1. Overview & Core Motivation

Across polyglot systems, CLI utilities, and HTTP services, verifying whether a user input, query parameter, or runtime token matches one of several candidate string values is a common operation:

```go
// Common confirmation prompt in CLI tools:
reply := getUserInput()
trimmed := strings.TrimSpace(reply)
if strings.EqualFold(trimmed, "y") || strings.EqualFold(trimmed, "yes") {
    // proceed
}
```

While functional on the surface, this ubiquitous idiom introduces severe architectural, performance, and maintenance friction when repeated ad-hoc across codebases:

### 1.1 Redundant Allocations & Computational Waste

In garbage-collected environments such as Go and Python, functions like `strings.ToLower()` or `str.lower()` allocate brand-new string objects on the heap. Even when utilizing case-insensitive comparisons such as Go's `strings.EqualFold()`, call sites frequently combine it with repeated `strings.TrimSpace()` calls. When matching against multiple candidates (`"y"`, `"yes"`, `"true"`, `"1"`), callers either allocate intermediate transformed strings or repeatedly execute transformation routines across chained conditions.

### 1.2 Cognitive Bloat & Call-Site Clutter

Chained boolean expressions using logical OR (`||`) bury core business logic underneath layers of mechanical string manipulation boilerplate:
- **Go:** `strings.EqualFold(trimmed, "y") || strings.EqualFold(trimmed, "yes") || strings.EqualFold(trimmed, "true")`
- **TypeScript:** `input.trim().toLowerCase() === 'y' || input.trim().toLowerCase() === 'yes' || input.trim().toLowerCase() === 'true'`
- **Rust:** `s.trim().eq_ignore_ascii_case("y") || s.trim().eq_ignore_ascii_case("yes")`
- **Python:** `s.strip().lower() == "y" or s.strip().lower() == "yes" or s.strip().lower() == "true"`

Every additional candidate broadens horizontal complexity and cyclomatic branching, expanding the surface area for logic errors and cluttering code reviews.

### 1.3 Asymmetry & Inconsistent Edge-Case Handling

Because individual developers implement string checks ad-hoc at each call site, edge-case handling fractures:
- Caller A trims whitespace but performs case-sensitive matching (`trimmed == "y"`).
- Caller B performs case-insensitive matching but neglects whitespace trimming (`strings.EqualFold(raw, "y")`).
- Caller C trims whitespace on the first candidate check but accidentally evaluates the raw string on the second candidate (`EqualFold(trimmed, "y") || EqualFold(raw, "yes")`).
- Caller D applies lowercase transformation instead of Unicode case folding.

### 1.4 Code Duplication & Re-Invention Fatigue

Without an authoritative canonical utility, AI agents and engineers repeatedly author one-off helper functions inside individual command or handler files (such as `isYes(s string) bool`, `checkConfirm(str string) bool`, or `matchesOption(opt string) bool`). This fragments repositories into unshared micro-helpers that violate the DRY (Don't Repeat Yourself) principle.

---

## 2. The Solution: `EqualFoldAnyTrim` & `EqualFoldAny`

To eliminate these inefficiencies across all supported platforms, this specification establishes a centralized utility function family anchored around `EqualFoldAnyTrim` and `EqualFoldAny`.

### 2.1 Function Signature & Variadic Design

The signature strictly adheres to the variadic parameter architecture codified in [`02-spec/02-coding-guidelines/01-cross-language/33-variadic-and-spread-parameters.md`](33-variadic-and-spread-parameters.md):

- **Target-First:** The primary string being inspected is passed as the first parameter.
- **Variadic Candidates:** Candidate matching strings are passed as trailing variadic arguments (`...string` in Go, `...candidates: readonly string[]` in TypeScript, `&[&str]` in Rust, `*candidates: str` in Python).
- **Zero-Wrapper Single & Multi Invocations:** Callers can check one, two, or multiple candidates with clean comma separation, or unpack an existing slice using the native spread operator (`candidates...`, `...candidates`).

```go
// Go: Centralized definition in pkg/strutil
package strutil

// EqualFoldAny reports whether target matches any candidate under Unicode case-folding.
func EqualFoldAny(target string, candidates ...string) bool

// EqualFoldAnyTrim reports whether target, with leading and trailing whitespace
// removed, matches any candidate under Unicode case-folding.
func EqualFoldAnyTrim(target string, candidates ...string) bool
```

### 2.2 Internal Mechanics & Zero-Allocation Efficiency

The canonical implementation of `EqualFoldAnyTrim` executes with optimal performance:

1. **Single Upfront Normalization:** The target string is trimmed exactly once upfront (`strings.TrimSpace(target)`).
2. **Short-Circuiting Evaluation:** Candidate iteration terminates immediately upon the first match, returning `true` without evaluating remaining candidates.
3. **Empty Collection Safety:** When zero candidates are provided, the utility returns `false` safely without panicking.
4. **Zero Intermediate Slice Construction:** When candidates are supplied as literals, compiler optimizations eliminate temporary heap allocation of slice headers.

---

## 3. Mandatory Rules

### Rule 1 (R1): Centralized String Utility Primacy (Banned Ad-Hoc Chained Comparisons)

- Application code, CLI commands, HTTP handlers, and services MUST NOT evaluate multiple string candidate matches using inline chained OR (`||`) expressions with repeated equality checks or transformations (such as `strings.EqualFold(s, "a") || strings.EqualFold(s, "b")`).
- Callers MUST invoke the centralized string utility function (`strutil.EqualFoldAnyTrim` or `strutil.EqualFoldAny`).
- Banned patterns include chaining `EqualFold`, repeating `TrimSpace()`, or concatenating `.trim().toLowerCase() === ...` chains.

### Rule 2 (R2): Target-First Variadic Candidate Signature

- The string utility function MUST place the target string being tested as the first parameter.
- Candidate matching strings MUST be accepted as trailing variadic arguments:
  - **Go:** `EqualFoldAnyTrim(target string, candidates ...string) bool`
  - **TypeScript:** `equalFoldAnyTrim(target: string, ...candidates: readonly string[]): boolean`
  - **Rust:** `equal_fold_any_trim(target: &str, candidates: &[&str]) -> bool`
  - **Python:** `equal_fold_any_trim(target: str, *candidates: str) -> bool`
- Callers can pass individual string candidates as comma-separated values or unpack an existing slice/array using the language's native spread operator (`candidates...`, `...candidates`).

### Rule 3 (R3): Pre-Normalization & Short-Circuit Optimization

- `EqualFoldAnyTrim` MUST perform whitespace trimming on the target string exactly once upfront prior to evaluating candidates.
- Candidate iteration MUST short-circuit and return `true` immediately upon finding the first match, avoiding redundant checks or allocations.
- If zero candidates are provided, the utility MUST safely return `false` without panicking or throwing errors.

### Rule 4 (R4): "Search First" Canonical Location Mandate

- Before authoring or proposing any string comparison helper, AI agents and engineers MUST inspect the repository's canonical utility package:
  - **Go:** `pkg/strutil/strutil.go` (or `internal/strutil/strutil.go`)
  - **TypeScript:** `src/lib/strutil.ts` (or `src/utils/strutil.ts`)
  - **Rust:** `src/util/strutil.rs` (or `src/strutil/mod.rs`)
  - **Python:** `pkg/strutil/strutil.py` (or `src/strutil.py`)
- Authoring private, unshared helper functions inside individual feature or command files (such as `isYes()`, `checkConfirm()`, or `matchesOption()`) is **STRICTLY PROHIBITED**.
- If the utility does not yet exist in the repository, it MUST be added to the canonical `strutil` location rather than introduced as an isolated private helper.

#### Canonical Location Directory Table

| Language | Primary Canonical Path | Permitted Internal Path | Prohibited Ad-Hoc Paths |
|:---|:---|:---|:---|
| **Go** | `pkg/strutil/strutil.go` | `internal/strutil/strutil.go` | `cmd/*.go`, `pkg/util/helpers.go`, inline file helpers |
| **TypeScript** | `src/lib/strutil.ts` | `src/utils/strutil.ts` | `src/components/*.tsx`, `src/pages/*.ts`, local `helpers.ts` |
| **Rust** | `src/util/strutil.rs` | `src/strutil/mod.rs` | `src/main.rs`, inline sub-modules in `src/cmd/` |
| **Python** | `pkg/strutil/strutil.py` | `src/strutil.py` | `scripts/*.py`, inline functions in individual script files |

### Rule 5 (R5): Prompt Architect Boolean Hygiene

- Positive booleans MUST ALWAYS be evaluated implicitly: `if isMatch { ... }`.
- Evaluating booleans explicitly against `true` (e.g. `if isMatch == true`) is **TOTALLY BANNED**.
- Mixed polarity within a single condition (e.g. `if isReady && !isBlocked`) is **TOTALLY BANNED**; split into separate discrete guard clauses.
- Boolean variables and helper return values MUST use `is` or `has` prefixes (e.g. `isMatch`, `hasCandidate`).

### Rule 6 (R6): Strict Vertical Line Spacing

- Maintain mandatory blank line spacing across all code examples and implementations:
  - Blank line before every `if` statement.
  - Blank line after every closing brace `}`.
  - Blank line before every `return` statement.
  - Blank lines around multiline struct initializations and function parameters.

---

## 4. Real-World Case Study: `releaseundo.go`

A prime example of the chained comparison anti-pattern is found in production CLI implementations such as `cli/cmd/releaseundo.go`:

### 4.1 The Anti-Pattern (`cli/cmd/releaseundo.go`)

```go
// ❌ WRONG: Chained EqualFold with manual TrimSpace (releaseundo.go anti-pattern)
func confirmUndoRelease(tag string) bool {
    fmt.Printf("Delete %s locally and on origin? [y/N]: ", tag)
    var reply string
    _, _ = fmt.Scanln(&reply)
    trimmed := strings.TrimSpace(reply)

    return strings.EqualFold(trimmed, "y") || strings.EqualFold(trimmed, "yes")
}
```

#### Deficiencies in this implementation:

1. **Redundant intermediate variable:** `trimmed` is declared solely to feed two successive `strings.EqualFold()` calls.
2. **Horizontal expansion:** Supporting additional confirmations (`"true"`, `"1"`, or localized equivalents) multiplies the `||` chain linearly.
3. **Zero reusability:** Any other command requiring user confirmation must duplicate this identical logic or create a divergent variant.

### 4.2 The Canonical Solution

```go
// ✅ REQUIRED: Centralized strutil.EqualFoldAnyTrim invocation
package cmd

import (
    "fmt"

    "github.com/alimtvnetwork/gitmap-v28/cli/pkg/strutil"
)

func confirmUndoRelease(tag string) bool {
    fmt.Printf("Delete %s locally and on origin? [y/N]: ", tag)
    var reply string
    _, _ = fmt.Scanln(&reply)

    return strutil.EqualFoldAnyTrim(reply, "y", "yes")
}
```

---

## 5. Cross-Language Implementations & Patterns

### 5.1 Go (`pkg/strutil/strutil.go`)

```go
// ✅ Canonical package implementation in pkg/strutil/strutil.go:
package strutil

import "strings"

// EqualFoldAny reports whether target matches any candidate under Unicode case-folding.
func EqualFoldAny(target string, candidates ...string) bool {
    for _, candidate := range candidates {
        isMatch := strings.EqualFold(target, candidate)
        if isMatch {
            return true
        }
    }

    return false
}

// EqualFoldAnyTrim reports whether target (after trimming whitespace) matches any candidate.
func EqualFoldAnyTrim(target string, candidates ...string) bool {
    trimmed := strings.TrimSpace(target)

    return EqualFoldAny(trimmed, candidates...)
}
```

### 5.2 TypeScript (`src/lib/strutil.ts`)

#### ❌ Anti-Pattern: Inefficient Repeated Chaining

```typescript
// ❌ WRONG: Inefficient repeated chaining and array inclusion
function isAffirmative(input: string): boolean {
  const clean = input.trim().toLowerCase();
  return clean === 'y' || clean === 'yes' || clean === 'true';
}
```

#### ✅ Canonical Implementation & Usage

```typescript
// ✅ Canonical implementation in src/lib/strutil.ts:
export function equalFoldAny(
  target: string,
  ...candidates: readonly string[]
): boolean {
  const normalizedTarget = target.toLowerCase();

  for (const candidate of candidates) {
    const isMatch = normalizedTarget === candidate.toLowerCase();
    if (isMatch) {
      return true;
    }
  }

  return false;
}

export function equalFoldAnyTrim(
  target: string,
  ...candidates: readonly string[]
): boolean {
  const trimmed = target.trim();

  return equalFoldAny(trimmed, ...candidates);
}

// Call site:
import { equalFoldAnyTrim } from '@/lib/strutil';

const isConfirmed = equalFoldAnyTrim(userInput, 'y', 'yes', 'true');
```

### 5.3 Rust (`src/util/strutil.rs`)

#### ❌ Anti-Pattern: Manual Trimming and Chained Calls

```rust
// ❌ WRONG: Manual trimming and chained eq_ignore_ascii_case
fn is_positive_response(input: &str) -> bool {
    let trimmed = input.trim();
    trimmed.eq_ignore_ascii_case("y") || trimmed.eq_ignore_ascii_case("yes")
}
```

#### ✅ Canonical Implementation & Usage

```rust
// ✅ Canonical implementation in src/util/strutil.rs:
pub fn equal_fold_any(target: &str, candidates: &[&str]) -> bool {
    for candidate in candidates {
        let is_match = target.eq_ignore_ascii_case(candidate);
        if is_match {
            return true;
        }
    }

    false
}

pub fn equal_fold_any_trim(target: &str, candidates: &[&str]) -> bool {
    let trimmed = target.trim();

    equal_fold_any(trimmed, candidates)
}

// Call site:
use crate::util::strutil::equal_fold_any_trim;

let is_positive = equal_fold_any_trim(user_input, &["y", "yes"]);
```

### 5.4 Python (`pkg/strutil/strutil.py`)

#### ❌ Anti-Pattern: Chained Transforms Across Scripts

```python
# ❌ WRONG: Chained lower/strip checks scattered across scripts
if s.strip().lower() == "y" or s.strip().lower() == "yes":
    process()
```

#### ✅ Canonical Implementation & Usage

```python
# ✅ Canonical implementation in pkg/strutil/strutil.py:
def equal_fold_any(target: str, *candidates: str) -> bool:
    """Report whether target matches any candidate case-insensitively."""
    normalized_target = target.casefold()

    for candidate in candidates:
        is_match = normalized_target == candidate.casefold()
        if is_match:
            return True

    return False


def equal_fold_any_trim(target: str, *candidates: str) -> bool:
    """Report whether trimmed target matches any candidate case-insensitively."""
    trimmed = target.strip()

    return equal_fold_any(trimmed, *candidates)


# Call site:
from pkg.strutil import equal_fold_any_trim

if equal_fold_any_trim(user_input, "y", "yes"):
    process()
```

---

## 6. Comprehensive Good vs Bad Contrast Matrix

| Attribute | ❌ Ad-Hoc Inline Chaining (`EqualFold(s, "y") \|\| ...`) | ✅ Centralized `EqualFoldAnyTrim(s, ...)` | Architectural Benefit |
|:---|:---|:---|:---|
| **Call Site Readability** | High clutter; multiple `strings.EqualFold()` calls | Clean single-line function call: `strutil.EqualFoldAnyTrim(s, "y", "yes")` | Immediate clarity of intent; low cognitive load |
| **Maintenance & Extensibility** | Expanding candidates requires modifying boolean logic operators | Adding a candidate requires only adding a comma-separated argument | Minimal diff surface; zero risk of operator precedence errors |
| **Memory Allocation** | Often allocates new strings per candidate check | Trims once upfront; zero allocations during candidate comparison | Reduced GC pressure in high-throughput hot paths |
| **Whitespace Resilience** | Callers frequently forget `strings.TrimSpace()`, leading to subtle input bugs | Trimming is guaranteed and standardized internally | Immunity to trailing newlines, carriage returns, and spaces |
| **Reusability** | Logic re-implemented ad-hoc in every command, package, and handler | 100% centralized in canonical `strutil` package | Strict DRY adherence across all repository packages |
| **AI Agent Consistency** | AI models generate inconsistent helper variants (`isYes`, `checkMatch`, `validStr`) | AI models discover and invoke canonical `strutil` methods via Search First protocol | Deterministic code generation; zero codebase fragmentation |

---

## 7. Related Specifications

- [`02-spec/02-coding-guidelines/01-cross-language/02-boolean-principles/readme.md`](02-boolean-principles/readme.md) — Boolean naming and implicit condition standards
- [`02-spec/02-coding-guidelines/01-cross-language/08-dry-principles.md`](08-dry-principles.md) — Deduplication and DRY principles
- [`02-spec/02-coding-guidelines/01-cross-language/33-variadic-and-spread-parameters.md`](33-variadic-and-spread-parameters.md) — Variadic and spread parameter standard
- [`02-spec/21-app/02-string-normalization-and-equalfoldany/01-architecture-spec.md`](../../21-app/02-string-normalization-and-equalfoldany/01-architecture-spec.md) — Parent architecture specification

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TYPE-034: String Normalization & EqualFoldAny Conformance

**Given** String equality, case-insensitive comparison, and trimming operations across Go, TypeScript, Rust, and Python.
**When** Linters and CI suites scan codebase repositories for string matching patterns.
**Then** Multi-candidate comparisons invoke canonical `EqualFoldAnyTrim` / `EqualFoldAny` without inline OR chaining, preserving zero unnecessary heap allocations, positive booleans, and 100% relative paths.

- **AC-CG-034-A:** Spec file `02-spec/02-coding-guidelines/01-cross-language/34-string-normalization-and-equalfoldany.md` exists and contains 100% relative paths.
- **AC-CG-034-B:** Registry in `02-spec/02-coding-guidelines/01-cross-language/readme.md` contains sequence #34 without gap.
- **AC-CG-034-C:** All code examples pass guideline checks (positive booleans, zero `== true`, zero mixed polarity, vertical line gaps preserved).
- **AC-CG-034-D:** Includes verbatim documentation of the `releaseundo.go` case study.
- **AC-CG-034-E:** Function signatures across Go, TypeScript, Rust, and Python adhere to the target-first variadic candidate parameter architecture.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language/34-string-normalization-and-equalfoldany.md --check-only
```
**Expected:** exit 0. Zero violations detected.
