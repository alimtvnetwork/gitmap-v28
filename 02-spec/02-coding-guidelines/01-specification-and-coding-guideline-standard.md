# Specification and Coding Guideline Authoring Standard (AI Execution Prompt)

> **/goal** Establish the mandatory 4-part architectural anatomy, testable acceptance criteria, zero-violation lint standards, and strict authoring discipline across all coding guidelines in `02-spec/02-coding-guidelines/`.
> **/learn** Master the authoring structure for AI execution prompts, actionable CI/CD checklists, numbered normative rules (R1–R7), ❌ FORBIDDEN vs ✅ REQUIRED code comparisons, and Given/When/Then verification contracts.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify every coding guideline document begins with the canonical AI Execution Prompt header and `> **/goal**` / `> **/learn**` blockquotes.
- [ ] `/learn` Ensure every guideline contains the `## 🎯 Actionable CI/CD & Agent Checklist` with interactive checkboxes and the `. **CRITICAL AI INSTRUCTION:**` directive.
- [ ] `/goal` Validate that all specifications enforce numbered normative rules (`R1`–`R7`) with explicit ❌ FORBIDDEN vs ✅ REQUIRED code comparisons across Style, Booleans, Naming, and Polyglot standards.
- [ ] `/learn` Confirm that every specification ends with testable Acceptance Criteria (`AC-CG-*`) specifying automated verification commands returning exit code 0.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified and all coding guideline documents must strictly adhere to the standards outlined in this specification.

**Version:** 1.0.0
**Updated:** 2026-10-03
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Architectural Scope & Purpose

This specification governs the authoring, formatting, structuring, and validation of all coding guidelines in `02-spec/02-coding-guidelines/` and across connected repositories. Every coding guideline document serves a dual role:

1. **Human Engineering Manual:** Authoritative technical specification establishing clean code patterns, parameter boundaries, and architectural invariants for engineers.
2. **Autonomous AI Execution Protocol:** Grounded, prompt-engineered directives instructing AI agents, linters, and CI/CD validation runners how to refactor, verify, and enforce compliance deterministically without hallucination.

---

## 2. Mandatory 4-Part Anatomical Structure

Every guideline specification in `02-spec/02-coding-guidelines/` MUST strictly adhere to the following 4-part architectural anatomy:

```
┌────────────────────────────────────────────────────────┐
│ 1. AI Execution Prompt Header                          │
│    - Title with (AI Execution Prompt) suffix           │
│    - > **/goal** [Concrete operational goal]           │
│    - > **/learn** [Key mental models & anti-patterns]  │
├────────────────────────────────────────────────────────┤
│ 2. Actionable CI/CD & Agent Checklist                  │
│    - ## 🎯 Actionable CI/CD & Agent Checklist          │
│    - Checkboxes with /goal and /learn directives       │
│    - . CRITICAL AI INSTRUCTION directive block         │
├────────────────────────────────────────────────────────┤
│ 3. Core Specification Body                             │
│    - Metadata block (Version, Updated, Status, etc.)   │
│    - Numbered rules (R1, R2, ...)                      │
│    - ❌ FORBIDDEN vs ✅ REQUIRED code examples        │
│    - Style, Boolean, Naming, and Polyglot standards    │
├────────────────────────────────────────────────────────┤
│ 4. Verification & Acceptance Criteria                  │
│    - ## Verification & Acceptance Criteria             │
│    - AC-CG-[CATEGORY]-[NUM] with Given / When / Then   │
│    - Verification command and expected exit code 0     │
└────────────────────────────────────────────────────────┘
```

### 2.1 Part 1: AI Execution Prompt Header

Every file MUST begin with a top-level `# [Title] (AI Execution Prompt)` heading, immediately followed by `> **/goal**` and `> **/learn**` blockquotes:

- `> **/goal**`: Defines the concrete operational objective, files to modify, and standards to enforce.
- `> **/learn**`: Directs the AI agent to internalize anti-patterns, precedence rules, and structural boundaries before modifying code.

### 2.2 Part 2: Actionable CI/CD & Agent Checklist

Immediately following the header (or preceding the metadata block), every guideline MUST declare an actionable checklist containing at least 4 checkboxes prefixed with `/goal` and `/learn`, concluded by the critical instruction marker:

- Each checkbox represents a non-negotiable operational checkpoint.
- The `. **CRITICAL AI INSTRUCTION:**` directive seals the checklist and binds autonomous agents.

### 2.3 Part 3: Core Specification Body

The specification body details the normative requirements:
- **Metadata Block:** Version, Updated date, Status, AI Confidence, and Ambiguity rating.
- **Numbered Rules (`R1`, `R2`, ...):** Unambiguous normative requirements.
- **❌ FORBIDDEN vs ✅ REQUIRED Code Examples:** Side-by-side or paired comparisons demonstrating violations and compliant refactorings across supported languages.

### 2.4 Part 4: Verification & Acceptance Criteria

Every guideline file MUST conclude with formal, testable acceptance criteria using the `AC-CG-[CATEGORY]-[NUM]` taxonomy, structured in `Given / When / Then` Gherkin format, paired with an automated verification command returning exit code 0.

---

## 3. Numbered Normative Rules

### R1: Mandatory 4-Part Anatomical Structure
All specification files in `02-spec/02-coding-guidelines/` must implement all four architectural sections without omitting headers, prompt directives, actionable checklists, or testable acceptance criteria.

### R2: Strict Relative Git Paths Mandate (TOTAL BAN on Absolute Paths / `file:///` URIs)
All file references, documentation links, script invocations, subtask paths, and markdown anchors MUST use relative paths starting from the repository root or relative to the current file (e.g., `02-spec/02-coding-guidelines/03-coding-style-checklist.md`).
- ❌ **TOTAL BAN:** Never write absolute paths (`C:\...`, `/d/...`, `/home/...`, `/Users/...`) or URI protocols (`file:///...`).
- Cross-references must link cleanly within the repository structure.

### R3: Strict Lowercase File and Folder Naming
All files, folders, scripts, and documentation MUST use strict lowercase kebab-case naming.
- Two-digit zero-padded prefixes (e.g., `01-*.md`, `02-*.md`) are required for sequential ordering.
- Root or subfolder entry points MUST be named `readme.md` (unadorned, no `00-` or `01-` prefixes).
- Overview and index variants (`00-overview.md`, `01-index.md`, `index.md`) are strictly prohibited.

### R4: Code Style Guidelines & Canonical Sizing Limits
All coding guidelines must enforce the canonical sizing limits defined in `02-canonical-size-tier.md` and parameter constraints in `03-coding-style-checklist.md`:
- **Function Body:** ≤ 8 lines preferred, ≤ 15 lines hard cap (CODE-RED-004).
- **File Length:** ≤ 300 lines hard cap (React TSX components: ≤ 100 lines).
- **Struct / Class Length:** ≤ 120 lines.
- **Parameters per Function:** Maximum 3 parameters. Functions requiring 4+ parameters or exceeding 100 characters in signature length MUST use an options struct or class (`*Options` or `*Params`).
- **Return Values:** Exactly 1 value only (use structured result wrappers like `Result[T]` or `*appfault.AppError`).
- **Zero Magic Values:** All literal strings and magic numbers must be extracted into named constants or enum types.
- **Vertical Spacing:** Mandatory blank lines before `if`, after `}`, before `return`, and around multiline struct calls.
- **File Encoding & Line Endings:** UTF-8 without BOM, Unix LF (`\n`) line endings, and a trailing EOF newline.

### R5: Boolean Logic & Conditioning Standards
Boolean variables, properties, methods, and expressions must enforce strict affirmative logic:
- **Prefix Standard:** Booleans MUST begin with `is` or `has` prefixes (e.g., `isValid`, `hasPermission`). `can`, `should`, and `was` are banned for state representations.
- **Strict Implicit Evaluation (TOTAL BAN on Explicit True):** NEVER compare a boolean explicitly against `true` (e.g., `if isReady == true` or `if isReady === true`). Positive booleans must always be evaluated implicitly: `if isReady { ... }`.
- **No Mixed Polarity:** NEVER combine a positive check and a negative check in the same `if` condition (e.g., `if isA && !isB`). Extract the combined intent into a single affirmative boolean.
- **Zero Nested Ifs:** Flatten conditional branches using early returns and guard clauses.
- **Maximum 2 Conditions per `if`:** Never chain more than 2 conditions without extracting into descriptive intermediate booleans. Never mix `&&` and `||` in the same condition.
- **Semantic Inverse Naming:** Ban negative tokens (`not`, `no`, `non`, `un`). Use positive opposites: `isInactive` instead of `isNotActive`, `isMissing` instead of `isNotFound`, `isDisabled` instead of `isNotEnabled`.

### R6: Naming Conventions & Acronym Casing
Identifiers must follow strict semantic and typographic rules across all languages:
- **PascalCase Acronyms:** Treat acronyms as standard words with initial capitalization only (`UserId` not `UserID`, `HttpServer` not `HTTPServer`, `ApiUrl` not `APIURL`, `JsonData` not `JSONData`).
- **Enum Standards:** Enums must use PascalCase with a mandatory `Type` suffix (e.g., `StatusType`, `LogLevelType`, `PaymentMethodType`).
- **Language-Specific Idioms:**
  - TypeScript / JavaScript: `camelCase` for functions, methods, and variables; `PascalCase` for classes, interfaces, and types.
  - Go: `PascalCase` for exported identifiers; `camelCase` for unexported identifiers.
  - Python: `snake_case` for functions, methods, and variables; `PascalCase` for classes; `UPPER_SNAKE_CASE` for constants.
  - Rust: `snake_case` for functions and variables; `PascalCase` for structs, enums, and traits.
- **Ban on Generic Names:** Identifiers such as `data`, `obj`, `temp`, `val`, `res`, `foo` are banned. Every variable must clearly convey domain meaning.

### R7: Polyglot Standards & Structured Error Wrapping
All supported languages (Go, TypeScript, Python, Rust, PHP, C#) must adhere to uniform enterprise-grade error handling and type safety:
- **Zero Error Swallowing (CODE RED 🔴):** Never ignore errors, empty catch blocks, or swallow exceptions. Every error must be wrapped, logged, or returned.
- **Structured Error Handling:**
  - Go: Functions returning failure metadata must use `*appfault.AppError` and wrap third-party errors at first contact.
  - TypeScript: Use `AppError` or structured Result wrappers; never use loose `any`.
  - Python: Raise structured custom domain exceptions inheriting from base application faults.
- **Monadic Result Wrappers:** Prefer monadic containers (`Result[T]`) for pipeline operations, providing explicit `.HasError()` and `.Value()` accessors.
- **Temporary Scripts Hygiene:** Scratchpads, prototypes, and diagnostic harnesses must reside strictly within `.ai-memory/temp-scripts/` and must never be committed to source control.

---

## 4. ❌ FORBIDDEN vs ✅ REQUIRED Code Examples

### 4.1 Code Style: Function Parameters & Sizing

❌ **FORBIDDEN (Loose Parameters Exceeding Limit):**
```go
// Forbidden: 5 parameters, exceeds 3-parameter limit, signature > 100 characters
func CreateUserAccount(name string, email string, role string, isActive bool, sendWelcomeEmail bool) (*User, error) {
    // implementation
}
```

✅ **REQUIRED (Options Struct Pattern):**
```go
// Required: Clean options struct with affirmative booleans and structured return
type CreateUserParams struct {
    Name                    string
    Email                   string
    Role                    string
    IsActive                bool
    IsWelcomeEmailRequested bool
}

func CreateUserAccount(params CreateUserParams) Result[*User] {
    // implementation
}
```

---

### 4.2 Boolean Logic: Explicit True Checks & Mixed Polarity

❌ **FORBIDDEN (Explicit Comparison to True & Mixed Polarity):**
```typescript
// Forbidden: explicit comparison to true and mixed positive/negative polarity in same condition
if (isAuthorized === true && !isAccountLocked) {
    grantAccess();
}
```

✅ **REQUIRED (Implicit Evaluation & Extracted Intent):**
```typescript
// Required: implicit evaluation and single affirmative intermediate boolean
const isEligible = isAuthorized && isAccountActive;

if (isEligible) {
    grantAccess();
}
```

---

### 4.3 Naming Conventions: Acronym Casing & Enums

❌ **FORBIDDEN (All-Caps Acronyms & Un-Suffixed Enum):**
```typescript
// Forbidden: all-caps acronyms and missing Type suffix
enum Status {
    ACTIVE = "active",
    INACTIVE = "inactive"
}

interface UserProfile {
    USER_ID: string;
    HTTP_URL: string;
    STATUS: Status;
}
```

✅ **REQUIRED (PascalCase Acronyms & Type Suffix):**
```typescript
// Required: PascalCase acronyms and Type suffix on enum
enum UserStatusType {
    Active = "active",
    Inactive = "inactive",
}

interface UserProfile {
    userId: string;
    httpUrl: string;
    status: UserStatusType;
}
```

---

### 4.4 Polyglot Standards: Structured Error Wrapping

❌ **FORBIDDEN (Swallowed Error & Bare Error Return):**
```go
// Forbidden: bare error return and lost error context
func FetchOrder(orderId string) (*Order, error) {
    orderRecord, err := db.Query(orderId)
    if err != nil {
        return nil, err
    }

    return orderRecord, nil
}
```

✅ **REQUIRED (Structured appfault.AppError Wrapping):**
```go
// Required: structured appfault.AppError with rich context and monadic Result wrapper
func FetchOrder(orderId string) Result[*Order] {
    orderRecord, err := db.Query(orderId)
    if err != nil {
        return ResultFail[*Order](appfault.Wrap(err, "E_ORDER_FETCH_FAILED", "failed to query order from db"))
    }

    return ResultOk(orderRecord)
}
```

---

## 5. Acceptance Criteria Taxonomy & Categorization

Acceptance criteria in `02-spec/02-coding-guidelines/` must use canonical taxonomy prefixes:

| Category Prefix | Domain / Scope | Target Rule Focus |
|---|---|---|
| `AC-CG-ROOT-` | Root coding guideline architecture | Sizing tier, parameter rules, authoring standards |
| `AC-CG-STYLE-` | Code style, spacing, sizing, braces | Blank line before return, after `}`, function length ≤15 lines, max 3 params |
| `AC-CG-BOOL-` | Boolean logic, polarity, naming | Affirmative `is`/`has` prefixes, zero `== true`, zero mixed polarity |
| `AC-CG-NAME-` | Naming conventions, semantic clarity | PascalCase types/schemas, semantic unit-tagged names, zero generic garbage |
| `AC-CG-TYPE-` | Type safety, immutability, param structs | Parameter structs, immutability, casting elimination, Result wrappers |
| `AC-CG-ARCH-` | Architecture, DRY, complexity | Cyclomatic complexity ≤ 10, DRY extraction, static analysis |

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-ROOT-001: Coding Guideline Authoring Standard Conformance

**Given** All specification and coding guideline markdown files in `02-spec/02-coding-guidelines/`.
**When** Guidelines are audited for anatomical structure, relative paths, lowercase naming, and coding style rules.
**Then** Every guideline file strictly satisfies the 4-part anatomy (AI Header, Actionable Checklist, Core Specification with R-rules & bad/good examples, and testable Acceptance Criteria) with zero violations, 100% relative paths, Unix LF line endings, and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines --check-only
```
**Expected:** exit 0. Zero violations.
