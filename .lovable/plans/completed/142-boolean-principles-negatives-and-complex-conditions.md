# Plan 142: Boolean Principles, Negatives & Complex Conditions Coding Guideline Audit

Trigger Keywords & Aliases: `cg-boolean`, `cg-execute boolean`, `audit boolean`, `fix boolean negatives`, `fix complex conditions`, `affirmative booleans`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)

---

## 1. Executive Summary & Objective

Autonomously scan, plan, refactor, and fix all boolean naming, double negatives, mixed polarity, and complex condition violations across the codebase, modifying source files directly to enforce affirmative prefixes (`is` and `has` only (can, should, was, etc. are banned)), implicit evaluation (no `== true`), positive framing (no `!isSuccess`), and discrete condition decomposition until 100% green without stopping.

Plan 142 executes surgical refactoring across three major clusters of boolean anti-patterns:
1. **Banned Function Prefixes & Negative Variables:** Eliminating `wasFlagPassed`, `noCmdErr`, `noRemote`, `noVSCode`, `noDesktop`, `noVSCodeSync`, and `noAutoTags`, replacing them with affirmative `is*` and `has*` names.
2. **Negative CLI Options & Struct Fields:** Eliminating `NoCommit`, `NoPush`, `noMerge`, `noTidy`, `noFetch`, `shouldFetch`, `shouldSwitch`, `IsNoPush`, and `IsNoCommit`, replacing them with semantic positive equivalents (`IsSkipCommit`, `IsSkipPush`, `isSkipMerge`, `isSkipTidy`, `isSkipFetch`, `isFetchEnabled`, `isSwitchEnabled`, `IsSkipPush`, `IsSkipCommit`).
3. **Mixed Polarity Condition Decomposition:** Eliminating `&& !` and `! &&` mixed polarity expressions across `cli/cmd/` (in `backup.go`, `backup_cloud_ops.go`, `cdops.go`, `changelog.go`, `changeloggen.go`, `clustersubcmd.go`, `commit_push.go`, `create_cmd.go`, `dopendingretry.go`, `find_duplicates.go`, `find_files_match.go`, `fixauth.go`, `historyrewrite_push.go`, `ipchange_cmd.go`) using guard clauses and extracted affirmative booleans.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/01-index.md`: Hallucination prevention, micro-tasking, and strict relative paths.
- `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-index.md`: Boolean principles index, implicit checks, and inverse patterns.
- `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/02-naming-prefixes.md`: Principle 1 (`is`/`has` prefixes only) and Principle 2 (total ban on negative words `no`, `not`, `non`).
- `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-parameters-and-conditions.md`: Principle 6 (ban on mixed polarity in single if conditions).
- `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`: Positive guard functions and semantic inverse pairs.
- `spec/02-coding-guidelines/01-cross-language/22-variable-naming-conventions.md`: Affirmative boolean variable naming rules.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: File and function size limits (functions <= 8–15 lines, files <= 100 lines).
- `.agents/skills/cg-boolean-and-naming/skill.md`: Canonical skill for boolean principles and naming conventions.

---

## 3. Domain-Specific Rules for Boolean Refactoring

1. **Rule 1 (Affirmative Prefixes Only):** Every boolean variable, function, parameter, and struct field must begin with `is` or `has`. Prefixes such as `can`, `should`, `was`, `will`, `did`, and `must` are strictly forbidden.
2. **Rule 2 (Total Ban on Negative Words):** The words `no`, `not`, and `non` are strictly forbidden in boolean identifiers (e.g., `noPush` -> `isSkipPush`, `noMerge` -> `isSkipMerge`, `noCmdErr` -> `hasZeroCmdErr` / `isCmdSuccess`).
3. **Rule 3 (Zero Mixed Polarity in If Conditions):** Never combine positive checks with negative checks in a single `if` condition (`a && !b` or `!a && b`). Extract the negative check into a positive semantic boolean first, or decompose into discrete linear guard clauses.
4. **Rule 4 (Guard Clause Flattening & Anti-Compression):** Keep maximum conditional nesting depth <= 1. Never compress `if/else` into single lines or omit mandatory braces. Maintain clean blank lines before control statements and returns.

---

## 4. Exhaustive Violation Ledger

| Id | File | Line | Snippet | Planned Fix | Status |
|:---|:---|:---:|:---|:---|:---:|
| V-01 | `cli/cmdscan/flags.go` | 162 | `func wasFlagPassed(...) bool` | Rename to `isFlagPassed` and update call sites | Completed |
| V-02 | `cli/cluster/exec_install.go` | 88 | `noCmdErr := cmdErr == nil` | Rename to `isCmdSuccess := cmdErr == nil` | Completed |
| V-03 | `cli/cmd/create_ops.go` | 22 | `NoRemote bool` | Rename to `IsSkipRemote bool` | Completed |
| V-04 | `cli/cmd/create_ops.go` | 50 | `noRemote := hasArgFlag(args, "--no-remote")` | Rename to `isSkipRemote := hasArgFlag(...)` | Completed |
| V-05 | `cli/cmd/mv_flags.go` | 6 | `noVSCode bool` | Rename to `isSkipVSCode bool` | Completed |
| V-06 | `cli/cmd/mv_flags.go` | 7 | `noDesktop bool` | Rename to `isSkipDesktop bool` | Completed |
| V-07 | `cli/cmd/mv_relocate.go` | 31 | `if !opts.noVSCode` | Update to `if !opts.isSkipVSCode` | Completed |
| V-08 | `cli/cmd/mv_relocate.go` | 35 | `if !opts.noDesktop` | Update to `if !opts.isSkipDesktop` | Completed |
| V-09 | `cli/cmdvscode/vscodepmscan.go` | 12 | `func syncRecordsToVSCodePM(..., noVSCodeSync, noAutoTags bool)` | Rename to `isSkipVSCodeSync`, `isSkipAutoTags` | Completed |
| V-10 | `cli/cmdvscode/exports.go` | 69 | `func SyncRecordsToVSCodePM(..., noVSCodeSync, noAutoTags bool)` | Rename params to affirmative `isSkip*` | Completed |
| V-11 | `cli/cmdscan/exports.go` | 14 | `SyncRecordsToVSCodePMFn func(..., noVSCodeSync, noAutoTags bool)` | Rename params to affirmative `isSkip*` | Completed |
| V-12 | `cli/cmd/clihelpers.go` | 91 | `func syncRecordsToVSCodePM(..., noVSCodeSync, noAutoTags bool)` | Rename params to affirmative `isSkip*` | Completed |
| V-13 | `cli/cmdclone/helpers.go` | 123 | `type CGCommitOpts struct { NoCommit, NoPush bool }` | Rename to `IsSkipCommit`, `IsSkipPush` | Completed |
| V-14 | `cli/cmd/codingguidelines_commit.go` | 32 | `type CGCommitOpts struct { NoCommit, NoPush bool }` | Rename to `IsSkipCommit`, `IsSkipPush` | Completed |
| V-15 | `cli/cmd/codingguidelines_commit.go` | 171 | `func emitCGSkipNotes(..., noCommit, noPush bool)` | Rename to `isSkipCommit`, `isSkipPush` | Completed |
| V-16 | `cli/cmd/gomod.go` | 18 | `noMerge bool, noTidy bool` | Rename to `isSkipMerge bool, isSkipTidy bool` | Completed |
| V-17 | `cli/cmd/gomod.go` | 70 | `noMerge := fs.Bool(...)` | Rename to `isSkipMerge` | Completed |
| V-18 | `cli/cmd/gomod.go` | 149 | `func runGoModTidy(noTidy bool)` | Rename to `func runGoModTidy(isSkipTidy bool)` | Completed |
| V-19 | `cli/cmd/gomod_test.go` | 79 | `opts.dryRun \|\| opts.noMerge \|\| opts.noTidy` | Update to affirmative `opts.isSkipMerge`, etc. | Completed |
| V-20 | `cli/cmd/latestbranch.go` | 24 | `shouldFetch bool, shouldSwitch bool` | Rename to `isFetchEnabled bool, isSwitchEnabled bool` | Completed |
| V-21 | `cli/cmd/latestbranch.go` | 147 | `var allRemotes, noFetch, jsonOut, switchLong, switchShort bool` | Rename to affirmative `isSkipFetch`, `isAllRemotes`, etc. | Completed |
| V-22 | `cli/movemerge/types.go` | 65 | `IsNoPush bool, IsNoCommit bool` | Rename to `IsSkipPush bool, IsSkipCommit bool` | Completed |
| V-23 | `cli/cmd/movemergeflags.go` | 14 | `noPush, noCommit, forceFold, pullFold, dryRun bool` | Rename to affirmative `isSkipPush`, `isSkipCommit`, etc. | Completed |
| V-24 | `cli/cmd/historyrewrite_flags.go` | 16 | `noPush bool` | Rename to `isSkipPush bool` | Completed |
| V-25 | `cli/cmd/historyrewrite_push.go` | 18 | `if opts.noPush` | Rename to `if opts.isSkipPush` | Completed |
| V-26 | `cli/cmd/historyrewrite_push.go` | 24 | `if !opts.yes && !confirmHistoryPush(...)` | Decompose mixed/double negative into discrete affirmative flags | Completed |
| V-27 | `cli/cmd/backup.go` | 204 | `if err == nil && !info.IsDir()` | Flatten to discrete guard clauses | Completed |
| V-28 | `cli/cmd/backup_cloud_ops.go` | 222 | `if len(args) > 0 && !strings.HasPrefix(args[0], "-")` | Decompose using affirmative `isPositional` helper | Completed |
| V-29 | `cli/cmd/cdops.go` | 114 | `if len(dflt) > 0 && !pick` | Restructure to guard clause `if pick` then load default | Completed |
| V-30 | `cli/cmd/changelog.go` | 29 | `if !latest && len(version) == 0 && openFile` | Extract affirmative `isDefaultVersion` and compose | Completed |
| V-31 | `cli/cmd/changeloggen.go` | 76 | `if err != nil && !os.IsNotExist(err)` | Normalize `if errors.Is(err, os.ErrNotExist) { err = nil }` | Completed |
| V-32 | `cli/cmd/clustersubcmd.go` | 34 | `if isGit && !hasMinTokens` | Extract positive `hasTokensMissing` and compose `isGitMissing` | Completed |
| V-33 | `cli/cmd/clustersubcmd.go` | 38 | `if isProj && !hasMinTokens` | Extract positive `hasTokensMissing` and compose `isProjMissing` | Completed |
| V-34 | `cli/cmd/clustersubcmd.go` | 122 | `if hasCommaSuffix && !isStrippedEmpty` | Replace with positive `hasStrippedContent := stripped != ""` | Completed |
| V-35 | `cli/cmd/commit_push.go` | 154 | `if targetSha == "" && !strings.HasPrefix(arg, "-")` | Restructure with guard clause `if strings.HasPrefix(...) continue` | Completed |
| V-36 | `cli/cmd/create_cmd.go` | 15 | `if len(args) == 0 && !isInteractiveStdin()` | Extract affirmative `isHeadlessError` composite | Completed |
| V-37 | `cli/cmd/dopendingretry.go` | 75 | `if workDir != "" && !pathExists(workDir)` | Extract `isDirMissing := !pathExists(workDir)` and compose | Completed |
| V-38 | `cli/cmd/find_duplicates.go` | 36 | `if len(args) > 0 && !strings.HasPrefix(args[0], "-")` | Restructure into discrete linear guard clauses | Completed |
| V-39 | `cli/cmd/find_files_match.go` | 24 | `if len(opts.Exts) > 0 && !hasAllowedExtension(...)` | Extract positive `isExtDisallowed` and compose `isFilterRejected` | Completed |
| V-40 | `cli/cmd/fixauth.go` | 113 | `if keyExists && !force` | Extract `isSafeMode := !force` and compose `isReuseExisting` | Completed |
| V-41 | `cli/cmd/fixauth.go` | 119 | `if keyExists && !assumeYes && !confirmOverwrite(keyPath)` | Extract affirmative composed flags and eliminate mixed polarity | Completed |
| V-42 | `cli/cmd/ipchange_cmd.go` | 60 | `if doPing && !validatePing(ctx, "8.8.8.8", 3)` | Extract `isPingFailed := !validatePing(...)` and compose `isRollbackNeeded` | Completed |

---

## 5. Subtask Decomposition

- **Subtask 01:** `.lovable/plans/subtasks/142-booleans/01-banned-function-prefixes-and-negative-vars.md`
  - Targets: `cli/cmdscan/flags.go`, `cli/cluster/exec_install.go`, `cli/cmd/create_ops.go`, `cli/cmd/mv_flags.go`, `cli/cmd/mv_relocate.go`, `cli/cmdvscode/vscodepmscan.go`, `cli/cmdvscode/exports.go`, `cli/cmdscan/exports.go`, `cli/cmd/clihelpers.go`.
  - Scope: Refactor `wasFlagPassed` to `isFlagPassed`, eliminate `noCmdErr`, `noRemote`, `noVSCode`, `noDesktop`, `noVSCodeSync`, and `noAutoTags`.
- **Subtask 02:** `.lovable/plans/subtasks/142-booleans/02-negative-options-and-struct-fields.md`
  - Targets: `cli/cmdclone/helpers.go`, `cli/cmd/codingguidelines_commit.go`, `cli/cmd/gomod.go`, `cli/cmd/gomod_test.go`, `cli/cmd/latestbranch.go`, `cli/movemerge/types.go`, `cli/cmd/movemergeflags.go`, `cli/cmd/historyrewrite_flags.go`, `cli/cmd/historyrewrite_push.go`.
  - Scope: Refactor `NoCommit`, `NoPush`, `noMerge`, `noTidy`, `noFetch`, `shouldFetch`, `shouldSwitch`, `IsNoPush`, and `IsNoCommit` to affirmative positive identifiers.
- **Subtask 03:** `.lovable/plans/subtasks/142-booleans/03-eliminate-mixed-polarity-in-cli-cmd.md`
  - Targets: `cli/cmd/backup.go`, `cli/cmd/backup_cloud_ops.go`, `cli/cmd/cdops.go`, `cli/cmd/changelog.go`, `cli/cmd/changeloggen.go`, `cli/cmd/clustersubcmd.go`, `cli/cmd/commit_push.go`, `cli/cmd/create_cmd.go`, `cli/cmd/dopendingretry.go`, `cli/cmd/find_duplicates.go`, `cli/cmd/find_files_match.go`, `cli/cmd/fixauth.go`, `cli/cmd/ipchange_cmd.go`.
  - Scope: Eliminate all mixed polarity condition joins (`&& !` and `! &&`) using discrete guard clauses and affirmative composite booleans.
- **Subtask 04:** `.lovable/plans/subtasks/142-booleans/04-upgrade-linter-and-verify-quality-gates.md`
  - Targets: `linter-scripts/check-boolean-guidelines.py`, quality linters, test inventory lock, milestone consolidation.
  - Scope: Upgrade boolean linter, run all repository linters, record modified files under lock, atomic commit & push.

---

## 6. Verification Plan

1. **Go Compilation & Vet:** `go vet ./...` in `cli/` must exit 0 with zero errors.
2. **Boolean Guidelines Linter:** `python linter-scripts/check-boolean-guidelines.py` must pass with 0 violations.
3. **Nested If Linter:** `python linter-scripts/check-nested-ifs.py` must pass.
4. **Relative Paths Linter:** `python linter-scripts/check-relative-paths.py` must pass.
5. **Enum Guidelines Linter:** `python linter-scripts/check-enum-guidelines.py` must pass.
6. **Result Wrapper Auditor:** `python 03-ai-scripts/35-result-wrapper-auditor.py` must pass.
7. **Test Inventory Lock:** `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>` must register all modified files.

---

## 7. Execution Outcome & Verification Results

All 42 identified violations across Subtasks 01–04 were refactored, verified, and locked:
1. **Subtask 01 (Banned Function Prefixes & Negative Variables):**
   - Refactored `wasFlagPassed` to `isFlagPassed` in `cli/cmdscan/flags.go`.
   - Modernized `noCmdErr` to `isCmdSuccess` in `cli/cluster/exec_install.go`.
   - Renamed `noRemote` / `NoRemote` to `isSkipRemote` / `IsSkipRemote` in `cli/cmd/create_ops.go`.
   - Renamed `noVSCode` / `noDesktop` to `isSkipVSCode` / `isSkipDesktop` in `cli/cmd/mv_flags.go` and `cli/cmd/mv_relocate.go`.
   - Converted `noVSCodeSync` and `noAutoTags` across `cli/cmdvscode/`, `cli/cmdscan/`, and `cli/cmd/` to affirmative `isSkipVSCodeSync` and `isSkipAutoTags`.
2. **Subtask 02 (Negative Options & Struct Fields):**
   - Modernized `CGCommitOpts` in `cli/cmdclone/helpers.go`, `cli/cmdclone/clonefixrepo.go`, `cli/cmdclone/clonefixrepo_modifiers.go`, and `cli/cmd/codingguidelines_commit.go` to affirmative `IsSkipCommit` and `IsSkipPush`.
   - Modernized `goModOpts` in `cli/cmd/gomod.go` to affirmative `isSkipMerge` and `isSkipTidy`.
   - Modernized `latestBranchConfig` in `cli/cmd/latestbranch.go` and `cli/cmd/latestbranchresolve.go` to affirmative `isFetchEnabled`, `isSwitchEnabled`, `hasContainsFallback`, and `isRemoteFiltered`.
   - Modernized `cli/movemerge/types.go`, `cli/movemerge/finalize.go`, and `cli/cmd/movemergeflags.go` to `IsSkipPush` and `IsSkipCommit`.
   - Modernized `cli/cmd/historyrewrite_flags.go` and `cli/cmd/historyrewrite_push.go` to `isSkipPush`.
3. **Subtask 03 (Mixed Polarity Condition Elimination):**
   - Refactored all 14 mixed polarity conditions in `cli/cmd/` (`backup.go`, `backup_cloud_ops.go`, `cdops.go`, `changelog.go`, `changeloggen.go`, `clustersubcmd.go`, `commit_push.go`, `create_cmd.go`, `dopendingretry.go`, `find_duplicates.go`, `find_files_match.go`, `fixauth.go`, `historyrewrite_push.go`, `ipchange_cmd.go`).
   - Replaced mixed polarity with discrete linear guard clauses and affirmative composite booleans while maintaining maximum nesting depth <= 1.
4. **Subtask 04 (Linter Upgrades & Quality Gates):**
   - Upgraded `linter-scripts/check-boolean-guidelines.py` regexes for banned prefixes (`can`, `should`, `was`, `will`, `did`, `must`).
   - Verified `go vet ./...` in `cli/` exits 0 with zero errors.
   - Verified `python linter-scripts/check-boolean-guidelines.py` passes 100% across 2,863 files.
   - Verified `python linter-scripts/check-nested-ifs.py` passes 100% across 2,863 files and 205 TS files.
   - Verified `python linter-scripts/check-relative-paths.py` passes 100% across 6,837 files.
   - Verified `python linter-scripts/check-enum-guidelines.py` passes 100%.
   - Verified `python 03-ai-scripts/35-result-wrapper-auditor.py` passes 100%.
   - Recorded 41 modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record`.
