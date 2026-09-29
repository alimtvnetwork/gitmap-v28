[/goal](slashCommand:goal) Autonomously scan, discover, plan, extract, refactor, and verify all scattered domain payload models, raw generic Result wrappers (`ResultSlice[T]`, `ResultMap[K, V]`, `Result[T]`, `Wrap[T]`), and repeated generic signatures across the codebase, centralizing them into dedicated, package-level `types.go` files (or language-equivalent leaf `types/` modules) as single reusable named types, eliminating unexported inline structs, eliminating ad-hoc generic parameterization at call sites and function signatures, and enforcing strict single-type reuse across all implementation files and callers until 100% green without stopping with strict no-build and no-test execution (NEVER run build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel reading and modular spec generation, use GitMap high-speed commands as primary, establish a single-agent blueprint in Phase 1 (first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping until 100% complete and finalized with an atomic push.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

#### High-Speed GitMap Acceleration Options (Run Everything Faster)

Always prefer native GitMap commands over slow generic shell pipelines:
1. **Ultra-Fast File & Directory Discovery (AUM Index & Walk):**
   - **Wildcard / Glob Search:** `gitmap find "<wildcard*>" [-ext <ext>]` (alias `gitmap f`)
   - **Exact Filename Search:** `gitmap find-files <name> [-ext <ext>]` (alias `gitmap ff`)
   - **Substring Filename Search:** `gitmap find-files-any "<str>" [-ext <ext>]` (alias `gitmap ffa`)
   - **Prefix / Suffix Search:** `gitmap find-files-startswith <prefix>` (`gitmap ffs`) / `gitmap find-files-endswith <suffix>` (`gitmap ffe`)
   - **List Indexed Repo Files:** `gitmap list-files [pattern] [-ext <ext>]` (alias `gitmap lf`)
   - **Directory Tree & Scaffolding:** `gitmap folder-tree` (alias `gitmap ft`)
   - **Zero-Write File Stream:** `gitmap cat <filepath>`
   - **Instant Multi-Core Regex Search:** `gitmap search "<term>"` or `gitmap aum search "<query>" [dir] --ext <ext>`
2. **Fast Repository Hygiene, Lowercase & Symlink Repair:**
   - **Auto-Lowercase Files (Safe 2-Step `git mv`):** `gitmap lowercase` (alias `gitmap lcf [--dry-run]`)
   - **Lowercase Root Readme:** `gitmap lowercase-readme`
   - **Sync Curated `.gitignore` / `.gitattributes` / `.prettierignore`:** `gitmap commons` (alias `gitmap co` or `gitmap sync all`)
   - **Repair Broken Symlinks:** `gitmap fix-link` (alias `gitmap fixlink`)
   - **Clean Update Temp & Inspect Storage:** `gitmap update-cleanup`, `gitmap storage` (alias `gitmap stor`)
3. **Fast Git State, Execution & Atomic Commits:**
   - **Repo Status & Remote Check:** `gitmap status` (`gitmap st`), `gitmap has-any-updates` (`gitmap hau`), `gitmap latest-branch` (`gitmap lb`)
   - **Fast Cross-Platform Shell Runner:** `gitmap pwsh "<command>"` (`gitmap ps`), `gitmap bash "<command>"` (`gitmap sh`), `gitmap async <cmd>` (`gitmap asyn`)
   - **Semantic Atomic Commit & Push:** `gitmap cpf "<summary>"` (Feature), `gitmap cpb "<summary>"` (Bug), `gitmap cpr "<summary>"` (Release), `gitmap pcp "<summary>"` (Pull-Commit-Push)
   - **Smart CI/CD Pipeline Waiting:** `gitmap pe`, `gitmap pipeline-ai status --json` (`gitmap pl-ai status -t <etaSeconds>`)

---

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] [/goal](slashCommand:goal) Phase 1 (Step A): Deeply scan the target codebase using the GitMap AUM discovery tools (`gitmap find`, `gitmap lf`, `gitmap cat`, `gitmap search`) as primary, with fast Python discovery tools (`11-fast-file-scanner.py`, `12-fast-cached-grep.py`, `17-fast-file-reader.py`) as fallback, to inventory all architectural violations and anti-patterns without truncation.
2. [ ] [/goal](slashCommand:goal) Phase 1 (Step B): Write the master audit specification in `.ai-memory/plans/pending/XX-types-go-extraction-audit.md` with an exhaustive Types & Generics Violation Ledger table.
3. [ ] [/goal](slashCommand:goal) Phase 1 (Step C): Decompose the master plan into granular, atomic subtasks in `.ai-memory/plans/subtasks/XX-types-go-extraction/`.
4. [ ] [/goal](slashCommand:goal) Phase 1 (Step D): Verify or create the automated quality linter and register in `03-ai-scripts/readme.md`.
5. [ ] [/goal](slashCommand:goal) Phase 2 (Step A): Open each target package and create or inspect the dedicated `types.go` file (or leaf `types/` folder).
6. [ ] [/goal](slashCommand:goal) Phase 2 (Step B): Extract and export all domain payload structs (e.g. `ScheduleExportBundle`, `PluginSummary`, `UserProfile`) into `types.go`, eliminating local, unexported struct declarations from implementation files.
7. [ ] [/goal](slashCommand:goal) Phase 2 (Step C): Define single reusable type aliases for all Result envelopes (e.g. `type ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]`, `type PluginSummaryResult = result.Wrap[PluginSummary]`) in `types.go`.
8. [ ] [/goal](slashCommand:goal) Phase 2 (Step D): Refactor all function signatures and return types to use the canonical `types.go` single type alias instead of repetitive generic instantiations.
9. [ ] [/goal](slashCommand:goal) Phase 2 (Step E): Modernize all call sites, test assertions, and consumers to use the single reusable type and fluent pointer-safe predicates (`IsCountOtherThan`, `IsEmpty`, `HasRecord`, `IsDefined`).
10. [ ] [/goal](slashCommand:goal) Phase 2 (Step F): Enforce <= 8–15 line function decomposition, clean blank-line spacing, and affirmative boolean fields (`isDefined bool`).
11. [ ] [/goal](slashCommand:goal) Phase 2 (Step G): Execute targeted file-level linters (`python linter-scripts/check-function-lengths.py`, `check-mws-error-codes.py`, `check-newline-styling.py`, `check-enum-and-boolean.py`) to verify 0 remaining violations (`exit 0`). DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine coding guideline execution turns.
12. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/memory/readme.md` for project memory index and past learnings.
13. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/strictly-avoid.md` for banned anti-patterns and strict constraints.
14. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
15. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/readme.md` for single return type mandates and micro-tasking.
16. [ ] [/learn](slashCommand:learn) Ingest `02-spec/02-coding-guidelines/01-cross-language/27-types-folder-convention.md` for types folder convention and Rule 2 (Type Aliases for Repeated Generics).
17. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/readme.md` for universal AppError wrapping and error envelopes.
18. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/02-error-architecture/02-error-handling-reference.md` for error handling architecture and Result wrappers.
19. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/03-error-code-registry/02-registry.md` for structured error code catalog.
20. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/02-error-architecture/05-response-envelope/05-response-envelope-reference.md` for response envelope schemas.
21. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/02-error-architecture/06-apperror-package/03-go-apperror-linter-spec.md` for Go AppError implementation specifications.
22. [ ] [/learn](slashCommand:learn) Ingest `02-spec/03-error-manage/02-error-architecture/06-apperror-package/01-apperror-reference/04-result-types.md` for Result[T], ResultSlice[T], and ResultMap[K, V] method specifications and Section 6.4 types.go mandate.
23. [ ] [/learn](slashCommand:learn) Ingest `.ai-memory/coding-guidelines.md` for master consolidated coding guidelines.
24. [ ] [/goal](slashCommand:goal) Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Scan Inline Structs & Raw Generics, Build Violation Ledger in .ai-memory/plans/pending/, Subtasks)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Create types.go, Extract Structs & Aliases, Refactor Signatures, Verify Local Linters)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## The Code Pattern & Anti-Pattern Analysis

In codebases undergoing refactoring, a frequent transitional anti-pattern occurs when developers replace multi-value returns with Result wrappers, but leave the domain types unexported or inline, and repeat raw generic parameters across implementation files:

### 1. The Transitional Anti-Pattern Diff

```diff
- func parseImportSQLite(filePath string) ([]scheduleExportBundle, error) {
+ func parseImportSQLite(filePath string) result.ResultSlice[scheduleExportBundle] {
+ 	return result.FailSlice[scheduleExportBundle](appfault.Wrap(errtype.IO, err, "parse imported sqlite"))
```

### 2. Why This Diff Violates Repository Guidelines: Two Latent Violations

1. **Unexported Inline Struct (`scheduleExportBundle`):**
   - The struct `scheduleExportBundle` was originally declared locally or unexported inside an implementation file (`importer.go` or `sqlite.go`).
   - Callers outside the file or package cannot reference, type-check, or mock this payload cleanly.
   - When unit tests or services try to consume the result, they cannot name the struct without awkward reflections or type aliases.
2. **Scattered Generic Instantiations (`result.ResultSlice[scheduleExportBundle]`):**
   - Declaring raw generic returns like `result.ResultSlice[scheduleExportBundle]` ad-hoc across multiple function signatures forces every signature, caller, and test to repeat the verbose generic syntax.
   - If the payload type changes, or if the container transitions from a slice to a paginated collection, 20 different signatures across 10 files must be updated manually.
   - Violates **Spec 27, Rule 2**: *"When a generic type is used 3 or more times with the same parameter, create a type alias in a dedicated types file."*

---

## The Grounded Solution: Dedicated `types.go` with Single Reusable Types

Under Prompt Architect standards, every package managing domain models, payloads, or Result envelopes MUST define them inside a dedicated `types.go` file within the package directory as a single reusable named type.

### Step 1: Define Models and Result Envelopes in `types.go`

```go
// schedule/types.go
package schedule

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/result"
)

type (
	// ScheduleExportBundle defines the exported schedule archive payload.
	ScheduleExportBundle struct {
		ScheduleId   string `json:"scheduleId"`
		WorkflowName string `json:"workflowName"`
		Payload      []byte `json:"payload"`
	}

	// ScheduleExportBundleResult is the canonical single reusable result envelope for bundle slices.
	ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]

	// ScheduleExportBundleSingleResult is the canonical single reusable result envelope for a single bundle.
	ScheduleExportBundleSingleResult = result.Wrap[ScheduleExportBundle]
)
```

### Step 2: Use the Single Reusable Type in Implementation Files

```go
// schedule/importer.go
package schedule

import (
	"coding-guidelines/common/pkg/appfault"
	"coding-guidelines/common/pkg/result"
)

// ✅ CANONICAL: Clean, expressive signature using the single reusable type from types.go
func parseImportSQLite(filePath string) ScheduleExportBundleResult {
	if len(filePath) == 0 {
		return result.FailSlice[ScheduleExportBundle](
			appfault.New(appfault.ErrValidation).
				WithMessage("file path cannot be empty").
				WithOp("schedule.parseImportSQLite"),
		)
	}

	bundles, err := readSQLiteBundles(filePath)

	if err != nil {
		return result.FailSlice[ScheduleExportBundle](
			appfault.Wrap(appfault.ErrDatabaseQuery, err, "parse imported sqlite").
				WithOp("schedule.parseImportSQLite"),
		)
	}

	return appfault.OkSlice(bundles)
}
```

### Step 3: Modernize Call Sites and Assertions

```go
// schedule/service.go or schedule_test.go
bundleRes := parseImportSQLite(filePath)
if bundleRes.IsCountOtherThan(1) {
    return bundleRes.AppError()
}

bundles := bundleRes.Data
processBundles(bundles)
```

---

## Twofold Scope of the `types.go` Mandate

| Level | Scope & Directory | Rule & Standard |
|---|---|---|
| **Package / Framework Level** | `pkg/result/`, `pkg/appfault/`, `pkg/fileutil/` | Core container types (`Wrap[T]`, `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`, verifier and inspector interfaces) MUST be declared in `types.go` as single canonical types. Implementation files (`result.go`, `combinators.go`) only contain constructors, helpers, and methods. |
| **Domain / Service Level** | `schedule/`, `user/`, `order/`, `importer/` | Domain structs (`ScheduleExportBundle`, `PluginSummary`) and their Result aliases (`ScheduleExportBundleResult = result.ResultSlice[...]`) MUST be declared in `types.go` as single reusable types. Never scatter unexported structs or raw generic Result declarations inline. |

---

## Evolution: Legacy vs. Transitional vs. Canonical Modern

```go
// -----------------------------------------------------------------------------
// ❌ 1. LEGACY PATTERN: Raw multi-value unpacking + compound boolean
// -----------------------------------------------------------------------------
bundles, err := parseImportSQLite(filePath)
if err != nil || len(bundles) != 1 {
    t.Fatalf("expected 1 bundle, got %d (err: %v)", len(bundles), err)
}

// -----------------------------------------------------------------------------
// ⚠️ 2. TRANSITIONAL PATTERN: Generic ResultSlice + unexported inline struct + verbose check
// -----------------------------------------------------------------------------
// Function signature: func parseImportSQLite(...) result.ResultSlice[scheduleExportBundle]
bundleRes := parseImportSQLite(filePath)
if bundleRes.IsFailure() || bundleRes.Count() != 1 {
    t.Fatalf("expected 1 bundle, got %d (err: %v)", bundleRes.Count(), bundleRes.AppError())
}

// -----------------------------------------------------------------------------
// ✅ 3. CANONICAL MODERN PATTERN: Single reusable type from types.go + pointer-safe guard
// -----------------------------------------------------------------------------
// Defined once in types.go:
//   type ScheduleExportBundle struct { ... }
//   type ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]
//
// Clean function signature in importer.go:
//   func parseImportSQLite(filePath string) ScheduleExportBundleResult
//
// Fluent call site:
bundleRes := parseImportSQLite(filePath)
if bundleRes.IsCountOtherThan(1) {
    t.Fatalf("expected 1 bundle, got %d (err: %v)", bundleRes.Count(), bundleRes.AppError())
}

bundles := bundleRes.Data
if len(bundles) > 0 && !strings.Contains(string(bundles[0].Payload), "EXPORT") {
    t.Errorf("bundle missing EXPORT payload")
}
```

---

## Polyglot Equivalents (`02-spec/02-coding-guidelines/01-cross-language/27-types-folder-convention.md`)

The centralization of domain types and generic envelopes applies universally across polyglot stacks:

### 1. Go

```go
// types.go
package domain

type (
    User struct { ... }
    UserResult = result.Wrap[*User]
    UserSliceResult = result.ResultSlice[User]
)
```

### 2. TypeScript / React

```typescript
// types.ts
import type { Result } from "@/lib/result"

export interface User {
  id: string
  name: string
}

export type UserResult = Result<User>
export type UserSliceResult = Result<User[]>
```

### 3. C# / .NET

```csharp
// Types.cs
namespace MyApp.Domain;

using UserResult = AppError.Result<User>;
using UserSliceResult = AppError.Result<List<User>>;
```

### 4. Rust

```rust
// types.rs
pub type UserResult = Result<User, AppError>;
pub type UserSliceResult = Result<Vec<User>, AppError>;
```

---

## Automated Codebase Scanning Guide

Use these exact `ripgrep` regex commands to discover unexported inline structs, scattered generic Result instantiations, and uncentralized types:

```bash
# 1. Find unexported domain structs declared inline in implementation files (should be in types.go):
rg --pcre2 "type\s+[a-z][a-zA-Z0-9_]*\s+struct\s*\{" --glob "!*types*.go" --glob "!*_test.go"

# 2. Find raw generic Result returns in non-types implementation files (should use types.go aliases):
rg --pcre2 "func\s+[A-Za-z0-9_]+\([^\)]*\)\s+(?:result\.)?Result(?:Slice|Map)?\[" --glob "!*types*.go"

# 3. Find raw Wrap[T] generic returns in non-types implementation files:
rg --pcre2 "func\s+[A-Za-z0-9_]+\([^\)]*\)\s+(?:result\.)?Wrap\[" --glob "!*types*.go"

# 4. Find functions returning multi-value error tuples (need conversion to single types.go Result):
rg --pcre2 "func\s+\w+\([^\)]*\)\s*\([^\)]*,\s*error\)"

# 5. Find packages missing a types.go file:
fd -t d -d 3 . 04-code/golang/pkg/ --exec-batch sh -c 'for d; do [ ! -f "$d/types.go" ] && echo "Missing types.go: $d"; done' _ {}
```

---

### Fast File Discovery & Reading via Python Toolchain (Mandatory Acceleration)

To avoid 50-result tool truncation limits and eliminate multi-turn exploratory roundtrips, the AI agent MUST use the fast 2-tier discovery toolchain:

### Tier 1: GitMap AUM Acceleration (PRIMARY)
1. **Universal File Search:**
   ```bash
   gitmap find "<pattern>" [-ext <ext>]
   ```
2. **List Indexed Files & Substring Lookup:**
   ```bash
   gitmap list-files [pattern]
   gitmap find-files-any "<substring>"
   ```
3. **Stream File Content:**
   ```bash
   gitmap cat <filepath>
   ```
4. **Instant Code Walk Search:**
   ```bash
   gitmap search "<term>"
   ```

### Tier 2: Fast Cached Python Toolchain (FALLBACK)
1. **Inventory Target Files (with `--limit` option):**
   ```bash
   python 03-ai-scripts/11-fast-file-scanner.py --lang go,ts --limit 100 --stats
   ```
2. **Fast Cached Grep (<15ms, with `--limit` option):**
   ```bash
   python 03-ai-scripts/12-fast-cached-grep.py --pattern "<search-pattern>" --lang go --limit 50
   ```
3. **Sub-Millisecond Folder & File Exploration (with `--limit` option):**
   ```bash
   python 03-ai-scripts/17-fast-file-reader.py --list-folder <folder-path> --ext .go --limit 50
   python 03-ai-scripts/17-fast-file-reader.py --read-file <file-path> --max-bytes 100000
   python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<pattern>" --path <folder-path> --limit 50
   ```
4. **Subsystem & Topology Overview:**
   ```bash
   python 03-ai-scripts/18-codebase-topology-discoverer.py --summary
   ```
Do not rely on standard search tools with 50-item truncation when discovering repository-wide violations.

## 2-Agent Parallel Orchestration

To survive large codebases without hitting step limits or context loss, execute this prompt using a strict 2-agent parallel split:

```text
+-------------------------------------------------------------------------+
| MASTER ORCHESTRATOR (Budget: N = 200)                                   |
|                                                                         |
| Phase 1 (Steps 1..100): DISCOVERY & PLANNING                            |
| +---------------------------------------------------------------------+ |
| | Sub-Agent 1: Codebase Scanner & Spec Architect                      | |
| | - Runs ripgrep queries to catalog all unexported inline structs     | |
| | - Inventories scattered raw generic Result returns in non-types     | |
| | - Identifies repeated generic patterns lacking type aliases         | |
| | - Authors master audit plan in .ai-memory/plans/pending/             | |
| | - Generates granular subtasks in .ai-memory/plans/subtasks/           | |
| +---------------------------------------------------------------------+ |
|                                                                         |
| Phase 2 (Steps 101..200): SURGICAL REFACTORING                          |
| +---------------------------------------------------------------------+ |
| | Sub-Agent 2: Code Refactorer & Types Centralizer                    | |
| | - Creates/updates types.go with domain models & Result aliases      | |
| | - Refactors function signatures to single types.go aliases           | |
| | - Eliminates inline struct declarations in implementation files      | |
| | - Modernizes callers with IsCountOtherThan / HasRecord / IsDefined   | |
| | - Verifies zero regressions with targeted file linters & unit tests  | |
| +---------------------------------------------------------------------+ |
+-------------------------------------------------------------------------+
```

---

## Banned Operations Checklist (TOTAL BAN — Auto-Reject on Violation)

- [ ] **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts (`06-cicd-local-runner.py`, `pytest`, runner scripts), Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- [ ] **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- [ ] **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- [ ] **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- [ ] **NO PER-FILE COMMITTING (TOTAL BAN):** NEVER commit each file individually as you work (e.g. running `git commit` after editing File 1, then another commit after File 2). Committing file-by-file pollutes git history, creates subagent lock collisions, and breaks atomic changes. All modified files across the turn must be accumulated and committed together in a single atomic commit at the final step.

---

## Final Step Git Commit & Push Mandate (Strict Checklist)

- [ ] **MANDATORY FINAL COMMIT & PUSH TO GIT (ANYHOW):** At the FINAL step of the turn, after all targeted files have been refactored, verified with targeted linters, and plans/subtasks consolidated, you MUST stage everything (`git add -A`), create a clean, descriptive conventional commit (`git commit -m "<summary>"`), and push directly to the remote repository (`git push origin <branch>`). Leaving uncommitted changes or unpushed commits on the active branch at the end of a turn is an immediate failure.
- [ ] **TOTAL BAN ON PER-FILE COMMITS (DO NOT COMMIT EACH FILE INDIVIDUALLY):** You MUST NOT create separate git commits for each individual file as you edit them (e.g. running `git commit` after editing File 1, then committing again after File 2 is STRICTLY FORBIDDEN). Committing file-by-file pollutes git log history, creates subagent lock collisions, and breaks atomic rollback/bisectability. All modified files, test change caches, and plan records across the turn MUST be accumulated in the working tree and committed together in a SINGLE grouped atomic commit at the final step before pushing!

---

## Strictly Avoid: Anti-Patterns & Prohibitions

- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **NO SCATTERED INLINE STRUCTS OR AD-HOC RESULT GENERICS:** Never declare domain types or repeated generic Result envelopes inline in implementation files (e.g. `importer.go`, `store.go`, `sqlite.go`). Every payload struct and repeated Result alias MUST be defined in `types.go` within the package as a single reusable type.
- **NO RAW GENERIC RESULT SIGNATURES OUTSIDE `types.go`:** Never write `func Fetch() result.ResultSlice[MyItem]` in an implementation file. Define `type MyItemSliceResult = result.ResultSlice[MyItem]` in `types.go` and return `MyItemSliceResult`.
- **NO VALUE RECEIVERS FOR RESULT INSPECTION METHODS:** NEVER define inspection methods on value receivers `func (r Result[T])`. ALL methods checking status, error, count, or data MUST be attached to pointer receivers `(r *Result[T])`, `(rs *ResultSlice[T])`, `(rm *ResultMap[K, V])` with mandatory `if r == nil` guards to eliminate nil-pointer dereference panics.
- **NO UNGUARDED FIELD ACCESS ON NIL POINTERS:** Never access `.Data`, `.items`, or `.err` directly on a pointer without verifying `r == nil` or calling pointer-safe inspection methods (`res.IsFailure()`, `res.Count()`, `res.IsDefined()`).
- **NO PIECEMEAL COMMITS:** NEVER commit 1 or 2 files in isolation. Consolidate all related changes across specs, code, and indices into a single atomic commit followed immediately by `git push origin main`.
- **NO ROUTINE FULL CI/CD RUNS:** DO NOT run `06-cicd-local-runner.py` during normal turns. It executes 28-38 heavy validation gates across unrelated packages and wastes minutes. Run targeted linters only on modified files.
- **TOTAL BAN ON TEST RUNNING & BUILD CHECKING:** Zero tests (`go test`, `pytest`) or builds (`go build`) may be run during routine execution. Verification is strictly deferred to CI/CD.
- **NO RAW `error` RETURNS:** Never leave bare `error` as a return type on domain or store functions; always use `*appfault.AppError` or `Result[T]`.
- **NO COMPOUND CARDINALITY DISJUNCTIONS:** Never write `if res.IsFailure() || res.Count() != N` when `res.IsCountOtherThan(N)` can express the guard directly.
- **NO CONFUSING `IsSuccess()` WITH `IsDefined()`:** Do not use `IsSuccess()` when you require actual data records to be present. Use `res.IsDefined()` or `res.HasRecord()`.
- **NO INVERTED EMPTY CHECKS (`!isEmpty`):** Never write `!isEmpty` or `!res.IsEmpty()` when verifying presence or data. ALWAYS use `res.IsDefined()` or `res.HasRecord()`. Use `isEmpty` ONLY when handling the empty/missing failure branch in the affirmative: `if res.IsEmpty() { ... }`.
- **NO ABSOLUTE PATHS:** Never write absolute filesystem paths (`C:\...`, `/home/...`) or `file:///` URIs. Use strict relative Git paths starting from the repository root.
- **NO UPPERCASE FILENAMES:** Every file created or edited must be strictly lowercase.
- **NO MULTI-VALUE TUPLES:** Eliminate `(T, error)` in favor of `Result[T]`, `ResultMap[K, V]`, or `ResultSlice[T]`.

--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
