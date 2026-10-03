# File & Folder Naming — Acceptance Criteria Registry (AI Execution Prompt)

> **/goal** Provide a consolidated, traceable registry of testable acceptance criteria across all file and folder naming guideline specifications in `08-file-folder-naming/`.
> **/learn** Enforce the canonical criteria taxonomy (`AC-CG-FILE-[NUM]`), cross-language casing rules, framework idioms, strict lowercase repository hygiene, and verify compliance using targeted linters.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify each criterion maps 1:1 to an authoritative specification file in `08-file-folder-naming/`.
- [ ] `/learn` Ensure all criteria follow the structured `Given / When / Then` verification contract.
- [ ] `/goal` Validate that all verification commands execute cleanly with `Expected: exit 0`.
- [ ] `/learn` Verify 100% relative paths and zero absolute filesystem paths across all guideline links.
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 4.0.0
**Last Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. File and Folder Naming Criteria Inventory (`AC-CG-FILE-`)

| ID | Title | Authoritative Specification | Verification Command |
|:---|:---|:---|:---|
| `AC-CG-FILE-000` | File & Folder Naming Conventions Overview | [`readme.md`](readme.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |
| `AC-CG-FILE-002` | Cross-Language Universal File and Folder Naming Rules | [`02-cross-language.md`](02-cross-language.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |
| `AC-CG-FILE-003` | PHP and WordPress File and Folder Naming Standards | [`03-php-wordpress.md`](03-php-wordpress.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |
| `AC-CG-FILE-004` | Go File and Package Naming Conventions | [`04-golang.md`](04-golang.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |
| `AC-CG-FILE-005` | TypeScript and JavaScript File and Folder Naming Standards | [`05-typescript-javascript.md`](05-typescript-javascript.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |
| `AC-CG-FILE-006` | Rust and C# File and Directory Naming Conventions | [`06-rust-csharp.md`](06-rust-csharp.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |
| `AC-CG-FILE-REG-001` | File & Folder Naming Criteria Registry Conformance | [`97-acceptance-criteria.md`](97-acceptance-criteria.md) | `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only` |

---

## 2. Detailed Acceptance Criteria Specifications

### AC-CG-FILE-000: File & Folder Naming Conventions Overview

- [ ] Root `readme.md` provides navigation, scoring table, categories table, and cross-references.
- [ ] All files in `08-file-folder-naming/` adhere to active prompt anatomy, checklist headers, and acceptance criteria blocks.
- [ ] Internal relative markdown links resolve accurately with zero broken links.

**Given** Repository files, folders, and scripts.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase characters, invalid separators, or non-canonical folder paths exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-FILE-002: Cross-Language Universal File and Folder Naming Rules

- [ ] File names across all languages contain zero spaces or illegal special characters.
- [ ] Folder names are strictly lowercase across all languages except C# PascalCase directories.
- [ ] File extensions accurately reflect the implementation language.
- [ ] Test files follow source naming conventions (`*_test.go`, `*.test.ts`, `*Test.php`, `*Tests.cs`).
- [ ] PowerShell files use lowercase kebab-case for filenames and PascalCase `Verb-Noun` for internal functions.

**Given** Repository files, folders, and scripts.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase characters, invalid separators, or non-canonical folder paths exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-FILE-003: PHP and WordPress File and Folder Naming Standards

- [ ] General PHP files use `kebab-case.php`.
- [ ] WordPress class files strictly use `class-{name}.php` prefix with lowercase hyphenation.
- [ ] WordPress interface and trait files use `interface-{name}.php` and `trait-{name}.php`.
- [ ] Template files adhere to standard WordPress template hierarchy naming.
- [ ] Plugin folders use kebab-case matching the plugin slug, and internal directories use lowercase.

**Given** Repository files, folders, and scripts.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase characters, invalid separators, or non-canonical folder paths exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-FILE-004: Go File and Package Naming Conventions

- [ ] Go files use `snake_case.go` and unit tests use `*_test.go`.
- [ ] Package directories are strictly lowercase single words without hyphens or underscores.
- [ ] Enum packages strictly follow the `type` suffix convention (`providertype/`, `enginetype/`).
- [ ] Go project follows canonical layouts (`cmd/`, `internal/`, `pkg/`).

**Given** Repository files, folders, and scripts.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase characters, invalid separators, or non-canonical folder paths exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-FILE-005: TypeScript and JavaScript File and Folder Naming Standards

- [ ] React component files strictly use `PascalCase.tsx`.
- [ ] General utility, service, and config files use `kebab-case.ts`.
- [ ] Custom hooks strictly use `use-{name}.ts`.
- [ ] Type files use `kebab-case.types.ts`.
- [ ] Directory names are strictly lowercase `kebab-case`.

**Given** Repository files, folders, and scripts.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase characters, invalid separators, or non-canonical folder paths exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.

---

### AC-CG-FILE-006: Rust and C# File and Directory Naming Conventions

- [ ] Rust files and directories use `snake_case.rs` and `snake_case/`.
- [ ] C# files use `PascalCase.cs` matching class names with `I` prefix for interfaces.
- [ ] C# directories use `PascalCase/`.
- [ ] Rust crates use kebab-case crate names and snake_case in code.

**Given** Repository files, folders, and scripts.
**When** Audited against this file and folder naming specification.
**Then** Zero uppercase characters, invalid separators, or non-canonical folder paths exist and exit code is 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.

---

## 3. Cross-References & Traceability

- [`readme.md`](readme.md) — File & Folder Naming Conventions Overview
- [`02-cross-language.md`](02-cross-language.md) — Universal cross-language rules
- [`03-php-wordpress.md`](03-php-wordpress.md) — PHP and WordPress naming standards
- [`04-golang.md`](04-golang.md) — Go file and package naming rules
- [`05-typescript-javascript.md`](05-typescript-javascript.md) — TypeScript and JavaScript conventions
- [`06-rust-csharp.md`](06-rust-csharp.md) — Rust and C# conventions

---

## Verification & Acceptance Criteria

### AC-CG-FILE-REG-001: File & Folder Naming Criteria Registry Conformance

**Given** The file and folder naming criteria registry in `02-spec/02-coding-guidelines/08-file-folder-naming/97-acceptance-criteria.md`.
**When** Audited by the automated guideline validator.
**Then** All file and folder naming specification criteria correctly map to active guideline specifications with runnable verification commands and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/08-file-folder-naming --check-only
```
**Expected:** exit 0. Zero violations.
