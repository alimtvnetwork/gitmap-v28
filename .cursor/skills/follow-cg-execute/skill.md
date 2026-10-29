---
name: follow-cg-execute
description: Master 39-step orchestrator for executing all coding guideline compliance prompts in order. Steps 1-38 run in fast-forwarding mode (zero-build, zero-test, bounded micro-batches), and step 39 conducts terminal full builds, E2E test verification to 100% pass, and release ceremony.
---

# Master 39-Step Coding Guideline Orchestrator (`follow-cg-execute`)

> **Invocation:** `@[follow-cg-execute](file;.cursor/skills/follow-cg-execute)`

## 1. Overview & Operational Directive

`follow-cg-execute` is the master sequencing skill designed to execute all **39 Coding Guideline (`cg-*`) compliance prompts** across the target codebase in strict sequential order.

Rather than running isolated, ad-hoc passes, `follow-cg-execute` enforces a systematic, deterministic 39-step pipeline:
- **Sequential Flow via `next`:** The engine tracks progress from step 01 to step 39. Each iteration executes the active step and fast-forwards to the next.
- **Fast-Forwarding Micro-Batches (Steps 1–38):**
  * Strictly bounded to 5–8 files per pass.
  * Zero builds and zero test suite executions (`NO-BUILD`, `NO-TEST`).
  * Fast surgical refactoring adhering to leaf error sovereignty, parameter structs, positive booleans, and vertical newline spacing.
  * Intermediate checkpoints recorded locally via `gitmap commit "wip - step XX: <scope>"`.
- **Terminal Verification & Release Gate (Step 39):**
  * Full codebase build verification (`go build`, `npm run build`, `cargo build`).
  * End-to-end and integration test suite execution to 100% green pass.
  * Regression auditing and cognitive heatmap closure verification (`gitmap heatmap status`).
  * Automated minor bump release ceremony (`gitmap cpr "vX.Y.Z - coding guideline compliance"`).

---

## 2. The 39 Ordered Coding Guideline Steps

| Step # | Prompt Slug | Canonical Domain Focus | Fast-Forward Rules |
| :---: | :--- | :--- | :--- |
| **01** | `cg-boolean-and-naming` | Affirmative boolean prefixes (`is`, `has`), implicit evaluation, total ban on `== true`. | Zero build/test. Bounded 5–8 files. |
| **02** | `cg-database-conventions` | Split-DB SQLite three-tier architecture, PascalCase tables, zero raw SQL strings. | Zero build/test. Bounded 5–8 files. |
| **03** | `cg-enum-standards` | PascalCase enums with `Type` suffixes across all languages, exhaustive switch dispatch. | Zero build/test. Bounded 5–8 files. |
| **04** | `cg-error-management` | Structured error return type `*appfault.AppError`, leaf error sovereignty. | Zero build/test. Bounded 5–8 files. |
| **05** | `cg-spec-authoring` | Specification authoring standard, deterministic testability, zero vague guidelines. | Zero build/test. Spec files only. |
| **06** | `go-guidelines` | Go idiomatic standards, error wrapping, package naming without stutter. | Zero build/test. Bounded 5–8 files. |
| **07** | `cg-audit-gap-n-steps` | Deep gap analysis, scoring guideline compliance from 0 to 100 without code mutations. | Analysis only. Zero code edits. |
| **08** | `style-guidelines` | Mandatory vertical spacing (blank line before `if`, after `}`, before `return`). | Zero build/test. Bounded 5–8 files. |
| **09** | `code-hygiene` | LF line endings, trailing whitespace elimination, redundant `.gitkeep` removal. | Zero build/test. Bounded 5–8 files. |
| **10** | `relative-paths` | Strict relative Git paths, total ban on absolute paths or `file:///` URIs. | Zero build/test. Bounded 5–8 files. |
| **11** | `testing-and-coverage` | Test structure, table-driven tests, isolated unit fixtures. | Zero test runs. Structure edits only. |
| **12** | `terminal-ui-and-cli-styling` | Lipgloss color tokens, unified CLI help menus, and status formatters. | Zero build/test. Bounded 5–8 files. |
| **13** | `cli-commands-and-help` | Complete CLI subcommand registration, flag help text, and routing parity. | Zero build/test. Bounded 5–8 files. |
| **14** | `react-frontend-guidelines` | Decomposing oversized React components, strict prop types, hooks hygiene. | Zero build/test. Bounded 5–8 files. |
| **15** | `function-signatures-and-return-types` | Signature simplification, parameter structs for >2-3 arguments. | Zero build/test. Bounded 5–8 files. |
| **16** | `multi-language-enums-and-traits` | String-backed enums, MIME and media type enums, `HasEnumHelpers` traits. | Zero build/test. Bounded 5–8 files. |
| **17** | `nested-if-and-guard-clauses` | Early returns, guard clauses, flattening nested `if` statements. | Zero build/test. Bounded 5–8 files. |
| **18** | `lazy-regex-and-pattern-matching` | Pre-compiled lazy regular expressions, elimination of inline regex compilation. | Zero build/test. Bounded 5–8 files. |
| **19** | `function-argument-reduction-and-params` | Parameter struct DTOs for signatures with >2-3 parameters. | Zero build/test. Bounded 5–8 files. |
| **20** | `extract-types` | Centralizing domain types, structs, and interfaces into dedicated `types.go` files. | Zero build/test. Bounded 5–8 files. |
| **21** | `typescript-guidelines-and-types` | Discriminated unions, no `any`, strongly-typed `Result<T, AppError>` envelopes. | Zero build/test. Bounded 5–8 files. |
| **22** | `file-size-and-function-reduction` | 8–15 line functions max, 150–300 line files max (React UI 100 max). | Zero build/test. Bounded 5–8 files. |
| **23** | `string-operations-and-efficiency` | Zero-allocation string comparisons, short-circuit lazy evaluation. | Zero build/test. Bounded 5–8 files. |
| **24** | `isolate-destructive-os-and-heavy-unit-tests` | Strict mocking of system calls, ensuring tests never trigger OS shutdown/reboots. | Zero build/test. Bounded 5–8 files. |
| **25** | `nuclear-package-modularization-and-unit-test-optimization` | Modularizing monolithic packages into discrete bounded contexts. | Zero build/test. Bounded 5–8 files. |
| **26** | `python-boolean-and-conditional-enhancement` | Pythonic affirmative booleans, guard clauses, and early returns. | Zero build/test. Bounded 5–8 files. |
| **27** | `python-constants-and-magic-number-elimination` | Extracting magic numbers and raw strings into uppercase constants/enums. | Zero build/test. Bounded 5–8 files. |
| **28** | `python-array-constants-and-dynamic-enums` | Immutable tuples and Enum classes for fixed array collections. | Zero build/test. Bounded 5–8 files. |
| **29** | `code-dryness-and-library-extraction` | Eliminating repetitive duplication into reusable packages and libraries. | Zero build/test. Bounded 5–8 files. |
| **30** | `clean-repo-build-and-caches` | Purging ephemeral build caches, OS temp files, and sync `.gitignore`. | Cache purge only. Zero builds. |
| **31** | `cg-execute-in-below-steps` | Surgical refactoring fixes for instructions appended below using 5–8 file batches. | Zero build/test. Bounded 5–8 files. |
| **32** | `cg-follow-other-prompts` | Ingesting and executing referenced external prompts across the codebase. | Zero build/test. Bounded 5–8 files. |
| **33** | `branch-immutability-and-clean-construction` | Copy-on-write branch immutability, pure constructor returns. | Zero build/test. Bounded 5–8 files. |
| **34** | `golang-pointer-reduction-and-value-semantics` | Value semantics for small structs, reducing `*T` indirection. | Zero build/test. Bounded 5–8 files. |
| **35** | `clean-work-artifacts-and-os-caches` | Purging ephemeral build outputs and runtime caches. | Cache purge only. Zero builds. |
| **36** | `variadic-and-spread-parameters` | Enforcing variadic parameters over artificial slice wrapping. | Zero build/test. Bounded 5–8 files. |
| **37** | `string-normalization-and-equalfoldany` | Standardizing `strutil.EqualFoldAny` and zero-allocation case folding. | Zero build/test. Bounded 5–8 files. |
| **38** | `map-dispatch-and-handler-lookup` | Replacing long `switch` and `if-else` chains with map dispatch tables. | Zero build/test. Bounded 5–8 files. |
| **39** | `execute-with-e2e` | Terminal full builds, end-to-end verification, 100% test pass gate, release ceremony. | **FULL BUILD & E2E TESTS ENABLED.** |

---

## 3. Operational Execution Workflow

### Step 1: Initialize Task & Enqueue 39-Step Ledger
When triggered, initialize tracking in SQLite and Heatmap:
```bash
gitmap task init --name "Coding Guideline 39-Step Pass" --budget 400 --spec-file "02-spec/21-app/XX-task/01-architecture-spec.md" --spec-folder "02-spec/21-app/XX-task/"
gitmap heatmap log --task "<slug>" --action ANALYZE --file "all" --state EXPLORATION --expected "Initialize 39-step guideline pass"
```

### Step 2: The `next` Execution Loop (Steps 1–38)
On each `next` iteration:
1. Identify the current active step from the task plan.
2. Search candidate files using GitMap: `gitmap aum search "<pattern>" <dir> -e <ext>`
3. Refactor strictly 5–8 files applying the specific step's rule.
4. Log cognitive intent and closure:
   ```bash
   gitmap heatmap log --task "<slug>" --action EDIT --file "<path>" --state CODE_EDITING --expected "Apply step XX rule"
   gitmap heatmap update --id <id> --outcome "Refactored cleanly" --met true
   ```
5. Checkpoint commit: `gitmap commit "wip - step XX: <scope>"`
6. Advance pointer to step XX+1. Do NOT run compiler builds or test runners.

### Step 3: Terminal Verification & Release Gate (Step 39)
Upon reaching Step 39:
1. Execute complete clean build:
   ```bash
   gitmap pwsh "go build ./...; npm run build"
   ```
2. Run full test suite:
   ```bash
   gitmap pwsh "go test -v ./...; npm test"
   ```
3. Audit git status to verify zero uncommitted changes and clean tree.
4. Perform minor version bump release ceremony:
   ```bash
   gitmap cpr "vX.Y.Z - complete 39-step coding guideline compliance pass"
   git tag vX.Y.Z
   git push origin main --tags
   ```
