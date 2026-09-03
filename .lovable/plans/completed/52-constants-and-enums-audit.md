# 31 - Constants & Enums Architecture Audit Specification

## 1. Verbatim Acceptance Criteria Echo (from spec/02-coding-guidelines/01-cross-language/97-acceptance-criteria.md)

### AC-01: Guideline Coverage

- [ ] Boolean principles define naming, evaluation, and composition patterns
- [ ] Casting elimination patterns cover type-safe alternatives to type assertions
- [ ] Code style defines formatting, naming, and structural conventions

### AC-02: Enforcement

- [ ] All guidelines include ❌ (forbidden) and ✅ (compliant) code examples
- [ ] ESLint/linter rules are documented for automated enforcement
- [ ] Master guidelines document consolidates all standards for AI reference

---

## 2. Task-Specific Rule Set (Domain Rules)

1. **Rule CE-1 (Mandatory `*Type` Suffix):** Every enum definition across Go, TypeScript, PHP, and Python MUST end with `Type` (e.g. `CompressionModeType`, `OutputModeType`, `MatchKindType`, `OutputFormatType`, `ExitCodeType`).
2. **Rule CE-2 (Zero Magic Numbers, Runes, and Delimiters):** Raw rune conversions (`rune(10)`), inline magic status codes, and delimiter strings must be extracted into dedicated constants (`constants/` or `enums/`).
3. **Rule CE-3 (Logging & Test Exemption):** Informational log messages, format strings, and test assertion failure descriptions are exempt from constant extraction.
4. **Rule CE-4 (Dedicated Definition Files):** Enums and constants must live in dedicated packages or modules (`gitmap/constants`, `gitmap/enums`, `src/enums`, `src/types`).
5. **Rule CE-5 (Typed Enum Safety):** Call sites and function parameters must accept typed enum values rather than raw primitive strings/integers.

---

## 3. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
| :---: | :--- | :---: | :--- | :--- | :---: |
| V-01 | gitmap/cmd/find_files.go | 15 | `type MatchKind int` | Define `MatchKindType` with alias `MatchKind = MatchKindType` | FIXED |
| V-02 | gitmap/cliexit/report.go | 50 | `type OutputMode int` | Define `OutputModeType` with alias `OutputMode = OutputModeType` | FIXED |
| V-03 | gitmap/archive/create.go | 36 | `type CompressionMode string` | Define `CompressionModeType` with alias `CompressionMode = CompressionModeType` | FIXED |
| V-04 | gitmap/cmd/folder/folder.go | 14 | `type OutputFormat string` | Define `OutputFormatType` with alias `OutputFormat = OutputFormatType` | FIXED |
| V-05 | 03-ai-scripts/06-cicd-local-runner.py | 32 | Missing Enum Guidelines Linter | Added `Enum Guidelines Linter` to Batch 1 | FIXED |

---

## 4. Subtasks Breakdown

1. **Subtask 01 (`.lovable/plans/subtasks/31-enums/01-enum-type-suffix.md`)**:
   Refactor Go enum types (`MatchKind`, `OutputMode`, `CompressionMode`, `OutputFormat`) to `*Type` definitions.
2. **Subtask 02 (`.lovable/plans/subtasks/31-enums/02-constants-and-runes-verification.md`)**:
   Verify zero raw rune casts (`rune(10)`), zero magic string delimiters, and enforce dedicated constant packages.
3. **Subtask 03 (`.lovable/plans/subtasks/31-enums/03-ci-linter-verification.md`)**:
   Execute `python linter-scripts/check-enum-guidelines.py` and run full CI quality gates via `03-ai-scripts/06-cicd-local-runner.py`.
