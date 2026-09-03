# Cross-Language Coding Guidelines

**Version:** 3.2.0
**Status:** Active
**Updated:** 2026-04-16
**AI Confidence:** High
**Ambiguity:** None

---

## Keywords

`coding-standards` · `cross-language` · `code-style` · `naming-conventions` · `variable-naming` · `singular-plural` · `boolean-patterns` · `dry` · `strict-typing` · `function-naming` · `newline-condition` · `dead-code` · `comment-spacing` · `method-documentation` · `nesting` · `complexity` · `lazy-evaluation` · `regex` · `mutation-avoidance` · `null-safety` · `casting` · `database-naming` · `pascalcase` · `test-naming` · `solid-principles` · `interface-segregation` · `dependency-inversion` · `single-responsibility` · `braces` · `early-return` · `blank-lines` · `method-chaining` · `function-size` · `struct-size` · `else-after-return` · `loop-variable` · `map-naming` · `abbreviations`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `01-index.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |


## Purpose

Cross-language coding standards and conventions that apply to **all languages** in the project (TypeScript, Go, PHP, Rust). Language-specific specs reference these as the single source of truth.

---

## Categories

### 🎨 Code Style & Formatting

Rules for braces, spacing, blank lines, function size, struct size, chaining, and newlines.

| # | File | Description |
|---|------|-------------|
| 04 | `04-code-style/` | Braces, nesting ban, blank lines, function/struct size limits, dead code, comment spacing, docs (subfolder — 8 files) |
| 21 | `21-newline-styling-examples.md` | Before/after newline examples, Go newline constants |

### 📛 Naming Conventions

Variable, function, boolean, key, and database naming rules.

| # | File | Description |
|---|------|-------------|
| 02 | `02-boolean-principles/01-index.md` | P1–P6 boolean naming and logic patterns |
| 10 | `10-function-naming.md` | Verb-led naming, no boolean flag parameters |
| 11 | `11-key-naming-pascalcase.md` | PascalCase for API/database keys |
| 12 | `12-no-negatives.md` | Positive guard functions instead of `!` |
| 22 | `22-variable-naming-conventions.md` | Singular/plural, maps, loop variables, abbreviations |

### 🏗️ Architecture & Design Principles

SOLID, DRY, complexity limits, and structural patterns.

| # | File | Description |
|---|------|-------------|
| 06 | `06-cyclomatic-complexity.md` | Complexity limits and measurement |
| 08 | `08-dry-principles.md` | Don't Repeat Yourself — deduplication rules |
| 09 | `09-dry-refactoring-summary.md` | DRY refactoring examples and techniques |
| 23 | `23-solid-principles.md` | SOLID principles with cross-language examples |

### 🛡️ Type Safety & Data Handling

Strict typing, casting elimination, null safety, mutation avoidance.

| # | File | Description |
|---|------|-------------|
| 03 | `03-casting-elimination-patterns.md` | Type casting elimination patterns |
| 07 | `07-database-naming.md` | Database table/column naming conventions |
| 13 | `13-strict-typing.md` | Strict typing rules, max parameters |
| 18 | `18-code-mutation-avoidance.md` | Immutability patterns, avoid side effects |
| 19 | `19-null-pointer-safety.md` | Null/nil safety guards |

### 🔧 Patterns & Techniques

Nesting resolution, lazy evaluation, regex usage.

| # | File | Description |
|---|------|-------------|
| 16 | `16-lazy-evaluation-patterns.md` | Defer computation until needed |
| 17 | `17-regex-usage-guidelines.md` | Compile once, named groups, readability |
| 20 | `20-nesting-resolution-patterns.md` | Flatten nested code, avoid else after return |

### 🧪 Testing

Test naming and structure conventions.

| # | File | Description |
|---|------|-------------|
| 14 | `14-test-naming-and-03-structure.md` | Test naming patterns and file organization |
| 28 | `28-slug-conventions.md` | Slug format rules, REST API examples, generation code |

### 📋 Reference & Meta

Consolidated reference, audit logs, contradiction checks.

| # | File | Description |
|---|------|-------------|
| 01 | `01-issues-and-fixes-log.md` | Historical issues and fixes |
| 05 | `05-cross-spec-contradiction-checks.md` | Contradiction audit across specs |
| 15 | `15-master-coding-guidelines/01-index.md` | Condensed master guidelines (AI quick-ref) |

---

## Full Document Inventory

| # | File | Category |
|---|------|----------|
| 01 | `01-issues-and-fixes-log.md` | Reference |
| 02 | `02-boolean-principles/01-index.md` | Naming |
| 03 | `03-casting-elimination-patterns.md` | Type Safety |
| 04 | `04-code-style/` | Style (subfolder) |
| 05 | `05-cross-spec-contradiction-checks.md` | Reference |
| 06 | `06-cyclomatic-complexity.md` | Architecture |
| 07 | `07-database-naming.md` | Type Safety |
| 08 | `08-dry-principles.md` | Architecture |
| 09 | `09-dry-refactoring-summary.md` | Architecture |
| 10 | `10-function-naming.md` | Naming |
| 11 | `11-key-naming-pascalcase.md` | Naming |
| 12 | `12-no-negatives.md` | Naming |
| 13 | `13-strict-typing.md` | Type Safety |
| 14 | `14-test-naming-and-03-structure.md` | Testing |
| 15 | `15-master-coding-guidelines/01-index.md` | Reference |
| 16 | `16-lazy-evaluation-patterns.md` | Patterns |
| 17 | `17-regex-usage-guidelines.md` | Patterns |
| 18 | `18-code-mutation-avoidance.md` | Type Safety |
| 19 | `19-null-pointer-safety.md` | Type Safety |
| 20 | `20-nesting-resolution-patterns.md` | Patterns |
| 21 | `21-newline-styling-examples.md` | Style |
| 22 | `22-variable-naming-conventions.md` | Naming |
| 23 | `23-solid-principles.md` | Architecture |
| 16a | `16-static-analysis/01-index.md` | Enforcement |
| 28 | `28-slug-conventions.md` | Naming |
| 97 | `97-acceptance-criteria.md` | Meta |
| 98 | `98-changelog.md` | Meta |
| 99 | `99-consistency-report.md` | Meta |

**Total:** 29 files (25 spec files + 1 overview + 3 meta)

---

## Document Inventory

| File |
|------|
| 01-issues-and-fixes-log.md |
| 02-boolean-principles.md |
| 03-casting-elimination-patterns.md |
| 05-cross-spec-contradiction-checks.md |
| 06-cyclomatic-complexity.md |
| 07-database-naming.md |
| 08-dry-principles.md |
| 09-dry-refactoring-summary.md |
| 10-function-naming.md |
| 11-key-naming-pascalcase.md |
| 12-no-negatives.md |
| 13-strict-typing.md |
| 14-test-naming-and-03-structure.md |
| 15-master-coding-guidelines.md |
| 16-lazy-evaluation-patterns.md |
| 17-regex-usage-guidelines.md |
| 18-code-mutation-avoidance.md |
| 19-null-pointer-safety.md |
| 20-nesting-resolution-patterns.md |
| 21-newline-styling-examples.md |
| 22-variable-naming-conventions.md |
| 23-solid-principles.md |
| 24-boolean-flag-methods.md |
| 25-generic-return-types.md |
| 26-magic-values-and-immutability.md |
| 27-types-folder-convention.md |
| 97-acceptance-criteria.md |
| 98-changelog.md |
| 99-consistency-report.md |


## Cross-References

- [Parent Overview](../01-index.md) — Broader coding guidelines context
- [AI Quick-Reference Checklist](../06-ai-optimization/03-ai-quick-reference-checklist.md) — Condensed checklist for AI code generation
- [Condensed Master Guidelines](../06-ai-optimization/05-condensed-master-guidelines.md) — AI-optimized single-file reference
