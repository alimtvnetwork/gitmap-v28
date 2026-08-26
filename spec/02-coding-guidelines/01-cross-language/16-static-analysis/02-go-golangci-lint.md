# Go — golangci-lint Enforcement Rule Mapping

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Keywords

`go` · `golang` · `golangci-lint` · `static-analysis` · `revive` · `cyclop` · `gocognit`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present | ✅ |
| AI Confidence assigned | ✅ |
| Ambiguity assigned | ✅ |
| Keywords present | ✅ |
| Scoring table present | ✅ |

---

## Purpose

Maps every Go coding guideline to its **golangci-lint** linter rule. Teams adopting these guidelines can copy the `.golangci.yml` config and get instant enforcement.

---

## 1. Linter Selection

golangci-lint bundles 100+ linters. We enable only those that enforce our coding guidelines:

| Linter | Purpose | Our Guidelines |
|--------|---------|----------------|
| `revive` | Style, naming, complexity | Nesting, naming, function size |
| `govet` | Correctness, suspicious constructs | Type safety |
| `staticcheck` | Bug detection, simplification | Dead code, simplifications |
| `gosimple` | Code simplification | DRY, flatten logic |
| `unused` | Unused code detection | Dead code removal |
| `errcheck` | Unchecked error returns | Error handling |
| `gocritic` | Opinionated style checks | Nesting, boolean patterns |
| `cyclop` | Cyclomatic complexity | Max function complexity |
| `funlen` | Function length | Max 15 lines |
| `gocognit` | Cognitive complexity | Zero nesting / complexity |
| `dupl` | Duplicate code detection | DRY principles |
| `goconst` | Repeated string literals | No magic strings |
| `nestif` | Nested if detection | Zero nesting rule |
| `predeclared` | Shadowing predeclared identifiers | Type safety |
| `goimports` | Import ordering | 3-group import convention |
| `misspell` | Typo detection | Code quality |
| `nolintlint` | Enforce nolint justifications | No silent suppressions |

---

## 2. Guideline → golangci-lint Rule Mapping

### 2.1 Code Style & Structure

| Guideline | Spec Source | Linter | Rule / Setting |
|-----------|-------------|--------|----------------|
| Zero nested `if` | [Code Style §R2](../04-code-style/01-braces-and-nesting.md) | `nestif` | `min-complexity: 1` |
| Zero nested `if` (alt) | [Forbidden Patterns](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) | `revive` | `nested-structs` |
| No `else` after return | [Code Style §R7](../04-code-style/01-braces-and-nesting.md) | `revive` | `superfluous-else` |
| Max 15-line functions | [Code Style §R6](../04-code-style/04-function-and-type-size.md) | `funlen` | `lines: 15, statements: 10` |
| Max 400-line files (target 300) | [File Rules](../../03-golang/04-golang-standards-reference/01-file-and-function-rules.md) | `revive` | `file-header` + custom |
| Blank line before return | [Code Style §R4](../04-code-style/03-blank-lines-and-spacing.md) | `whitespace` | `multi-func: true` |
| Import 3-group ordering | [Concurrency & Patterns](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) | `goimports` | `local-prefixes: project/` |
| No dead code | [Code Style §R5](../04-code-style/00-overview.md) | `unused` | (enabled by default) |

### 2.2 Naming Conventions

| Guideline | Spec Source | Linter | Rule / Setting |
|-----------|-------------|--------|----------------|
| Boolean naming (`is/has/can/should/was/will`) | [Boolean Principles](../02-boolean-principles/00-overview.md) | `revive` | `var-naming` + custom regex |
| PascalCase keys (API/DB) | [Key Naming](../11-key-naming-pascalcase.md) | `revive` | `var-naming` |
| No raw negation (`!fn()`) | [No Negatives](../12-no-negatives.md) | `gocritic` | `unslice` / code review |
| No boolean flag params | [Function Naming](../10-function-naming.md) | `revive` | `flag-parameter` |

### 2.3 Type Safety & Error Handling

| Guideline | Spec Source | Linter | Rule / Setting |
|-----------|-------------|--------|----------------|
| No `interface{}` / `any` in exported APIs | [Forbidden Patterns](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) | `revive` | `use-any` (inverted) |
| No `fmt.Errorf` for service errors | [Type Safety §2](../../03-golang/04-golang-standards-reference/02-type-safety-and-errors.md) | `revive` | `error-strings` + code review |
| Unchecked error returns | [Type Safety §2](../../03-golang/04-golang-standards-reference/02-type-safety-and-errors.md) | `errcheck` | `check-blank: true` |
| No `init()` functions | [Forbidden Patterns](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) | `revive` | `no-init` (custom) |
| No global mutable state | [Forbidden Patterns](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) | `gocritic` | code review |
| No panic in handlers | [Forbidden Patterns](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) | `revive` | `deep-exit` |

### 2.4 Complexity & DRY

| Guideline | Spec Source | Linter | Rule / Setting |
|-----------|-------------|--------|----------------|
| Cyclomatic complexity | [Complexity](../06-cyclomatic-complexity.md) | `cyclop` | `max-complexity: 10` |
| Cognitive complexity | [Complexity](../06-cyclomatic-complexity.md) | `gocognit` | `min-complexity: 10` |
| No magic strings / numbers | [Magic Strings](../15-master-coding-guidelines/05-magic-strings-and-organization.md) | `goconst` | `min-len: 3, min-occurrences: 2` |
| No duplicate code | [DRY Principles](../08-dry-principles.md) | `dupl` | `threshold: 100` |

### 2.5 Concurrency

| Guideline | Spec Source | Linter | Rule / Setting |
|-----------|-------------|--------|----------------|
| Context propagation | [Concurrency](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) | `revive` | `context-as-argument` |
| Context first parameter | [Concurrency](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) | `revive` | `context-keys-type` |

---

## 3. SonarQube Go Rules

SonarQube's `sonar-go` plugin provides additional cross-language coverage:

| SonarQube Rule | Description | Our Guideline |
|----------------|-------------|---------------|
| S3776 | Cognitive Complexity | Max function lines + zero nesting |
| S1871 | Identical branches | DRY |
| S1066 | Collapsible if | Zero nesting |
| S4144 | Identical functions | DRY |
| S1192 | Duplicated string literals | No magic strings |
| S107 | Too many parameters | Max 3 parameters |
| S138 | Function too long | Max 15 lines |
| S134 | Nesting depth | Zero nesting |
| S1126 | Return boolean directly | Boolean principles |

---

## 4. Reference `.golangci.yml`

```yaml
run:
  timeout: 5m

linters:
  enable:
    - revive
    - govet
    - staticcheck
    - gosimple
    - unused
    - errcheck
    - gocritic
    - cyclop
    - funlen
    - gocognit
    - dupl
    - goconst
    - nestif
    - predeclared
    - goimports
    - misspell
    - nolintlint
    - whitespace

linters-settings:
  funlen:
    lines: 15
    statements: 10

  cyclop:
    max-complexity: 10

  gocognit:
    min-complexity: 10

  nestif:
    min-complexity: 1

  goconst:
    min-len: 3
    min-occurrences: 2

  dupl:
    threshold: 100

  errcheck:
    check-blank: true

  goimports:
    local-prefixes: project/

  revive:
    rules:
      - name: superfluous-else
        severity: error
      - name: context-as-argument
        severity: error
      - name: deep-exit
        severity: error
      - name: flag-parameter
        severity: error
      - name: var-naming
        severity: error

  nolintlint:
    require-explanation: true
    require-specific: true

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
```

---

## 5. Integration Checklist

| # | Task | Status |
|---|------|--------|
| 1 | golangci-lint v1.59+ installed | 🔲 |
| 2 | `.golangci.yml` placed at repo root | 🔲 |
| 3 | All linters from §1 enabled | 🔲 |
| 4 | `nestif` min-complexity set to 1 (zero tolerance) | 🔲 |
| 5 | `funlen` lines set to 15 | 🔲 |
| 6 | CI runs `golangci-lint run --timeout 5m` on every PR | 🔲 |
| 7 | SonarQube sonar-go plugin enabled | 🔲 |
| 8 | `nolintlint` requires explanation on every suppression | 🔲 |
| 9 | Team reviewed and approved thresholds | 🔲 |

---

## Cross-References

- [Static Analysis Overview](./00-overview.md) — Cross-language analyzer guide
- [Go Standards Reference](../../03-golang/04-golang-standards-reference/00-overview.md) — Full Go rules
- [Forbidden Patterns](../../03-golang/04-golang-standards-reference/06-concurrency-and-patterns.md) — Go forbidden patterns
- [Cross-Language Code Style](../04-code-style/00-overview.md) — Formatting rules
- [Boolean Principles](../02-boolean-principles/00-overview.md) — Boolean naming rules
- [DRY Principles](../08-dry-principles.md) — Deduplication rules

---

*Go golangci-lint enforcement v1.0.0 — maps every Go coding guideline to a linter rule — 2026-04-01*
