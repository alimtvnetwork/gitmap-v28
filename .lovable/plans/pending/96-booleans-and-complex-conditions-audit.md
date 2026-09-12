# Plan 96: Boolean Principles, Negatives & Complex Conditions — Master Architectural Specification

Trigger Keywords & Aliases: `cg-boolean`, `cg-execute boolean`, `audit boolean`, `fix boolean negatives`, `fix complex conditions`, `affirmative booleans`

> **Prompt Version:** 2.1.0  
> **Synchronization:** Main Meta-Repo & Connected Workspaces  
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)

---

## 1. Executive Summary & Objective

Autonomously scan, plan, refactor, and fix all boolean naming, double negatives, mixed polarity, and complex condition violations across the codebase, modifying source files directly to enforce affirmative prefixes (`is` and `has` only; `can`, `should`, `was`, etc. are banned), implicit evaluation (no `== true`), positive framing (no `!isSuccess`), and discrete condition decomposition until 100% green without stopping.

Automated AST and regex scanning identified:
1. **17 Inverted Success Violations:** Expressions using `!isSuccess` / `!r.IsSuccess` instead of positive failure methods (`r.IsFailed()`).
2. **28 Banned Boolean Function Prefixes:** Functions returning `bool` prefixed with `should*` (26) and `can*` (2) instead of mandatory `is*` or `has*`.
3. **Mixed Polarity Chains:** Conditions combining positive and negated boolean terms in `cmd/pullparallel.go` and `cmd/pushparallel.go`.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md`: Every boolean identifier (variable, property, parameter, method) must begin with `is` or `has` only. Prefixes such as `can`, `should`, `was`, `will`, `did` are strictly banned.
- `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/03-parameters-and-conditions.md`: Never mix positive and negative booleans in a single condition (`isA && !isB`). Extract discrete positive booleans.
- `spec/02-coding-guidelines/01-cross-language/12-no-negatives.md`: Total ban on `!isSuccess` inverted checks. Use `.IsFailed()` or `.IsFail`.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: Functions <= 8 lines preferred, <= 15 lines max. Files <= 100 coding lines.
- `spec/02-coding-guidelines/01-cross-language/04-code-style/04-blank-lines-and-spacing.md`: Required blank lines before `return` and after closing `}` (R13-R16).

---

## 3. Domain-Specific Rules for Boolean & Condition Hardening

1. **Rule 1 (Strict Affirmative Prefixes `is*` / `has*`):** Functions and methods that return a boolean MUST begin with `is` or `has` (e.g. `isStackTraceEnabled`, `isMultiCloneEnabled`, `isPRCreationEnabled`, `hasMatchingPattern`). Never use `can` or `should`.
2. **Rule 2 (Positive Failure Over Inverted Success):** Replace all occurrences of `!r.IsSuccess` with `r.IsFailed()`. Models returning results (`CloneResult`, `PromptRunStatusRecord`, `RunRecord`) must provide explicit affirmative failure methods.
3. **Rule 3 (Zero Mixed Polarity in Condition Chains):** Never combine positive checks with negated checks in one condition (`if a && !b`). Extract the negation into an affirmative boolean (`isStopRequested := result.IsFailed() && stopOnFail`).
4. **Rule 4 (Implicit Evaluation & Anti-Compression):** Never evaluate booleans explicitly against `true` or `false`. Never collapse `if` blocks into a single line.

---

## 4. Exhaustive Violation Ledger

### Group A: Inverted Success Checks (`!.*IsSuccess`)

| Id | File | Line | Snippet | Planned Fix | Status |
|:---|:---|:---:|:---|:---|:---:|
| V-S01 | `gitmap/cmd/install_audit_runner_test.go` | 13 | `if !res.IsSuccess {` | Replace with `if res.IsFailed() {` | Pending |
| V-S02 | `gitmap/cmd/install_logs.go` | 75 | `if !r.IsSuccess && len(filtered) < limit {` | Replace with `if r.IsFailed() && len(filtered) < limit {` | Pending |
| V-S03 | `gitmap/cmd/prompt_failure_reporter.go` | 13 | `if !r.IsSuccess {` | Replace with `if r.IsFailed() {` | Pending |
| V-S04 | `gitmap/cmd/pull.go` | 517 | `if !result.IsSuccess {` | Replace with `if result.IsFailed() {` | Pending |
| V-S05 | `gitmap/cmd/pullparallel.go` | 108 | `if !result.IsSuccess {` | Replace with `if result.IsFailed() {` | Pending |
| V-S06 | `gitmap/cmd/pullparallel.go` | 112 | `shouldStop := !result.IsSuccess && stopOnFail` | Replace with `isStopRequested := result.IsFailed() && stopOnFail` | Pending |
| V-S07 | `gitmap/cmd/push.go` | 223 | `if !result.IsSuccess {` | Replace with `if result.IsFailed() {` | Pending |
| V-S08 | `gitmap/cmd/pushparallel.go` | 99 | `if !result.IsSuccess {` | Replace with `if result.IsFailed() {` | Pending |
| V-S09 | `gitmap/cmd/pushparallel.go` | 102 | `if !result.IsSuccess && stopOnFail {` | Replace with `if result.IsFailed() && stopOnFail {` | Pending |
| V-S10 | `gitmap/cmd/pushparallel.go` | 105 | `if !result.IsSuccess {` | Replace with `if result.IsFailed() {` | Pending |
| V-S11 | `gitmap/cmd/schedule_cmd.go` | 408 | `if !r.IsSuccess {` | Replace with `if r.IsFailed() {` | Pending |
| V-S12 | `gitmap/lazyregex/compile_result.go` | 41 | `return !it.IsSuccess()` in `IsFailed()` | Replace with `return it.appError != nil || it.re == nil` | Pending |
| V-S13 | `gitmap/lazyregex/lazyregex_test.go` | 212 | `if !validRes.IsSuccess() {` | Replace with `if validRes.IsFailed() {` | Pending |
| V-S14 | `gitmap/lazyregex/lazyregex_test.go` | 359 | `if !res.IsSuccess() \|\| res.Value != re1 {` | Replace with `if res.IsFailed() \|\| res.Value != re1 {` | Pending |
| V-S15 | `gitmap/pipelinedb/pipeline_repository_test.go` | 90 | `if !runRecord.IsSuccess {` | Replace with `if runRecord.IsFailed() {` | Pending |
| V-S16 | `gitmap/result/result_test.go` | 13 | `if !res.IsSuccess() {` | Replace with `if res.IsFailed() {` | Pending |
| V-S17 | `gitmap/tests/cmd_test/prompt_install_e2e_test.go` | 15 | `if !res.IsSuccess {` | Replace with `if res.IsFailed() {` | Pending |

### Group B: Banned Boolean Function Prefixes (`can*` / `should*`)

| Id | File | Line | Function Name | Replacement Name | Status |
|:---|:---|:---:|:---|:---|:---:|
| V-P01 | `gitmap/cliexit/handle.go` | 120 | `shouldPrintStackTrace` | `isStackTraceEnabled` | Pending |
| V-P02 | `gitmap/clonefrom/execute_dest.go` | 53 | `shouldSkip` | `isSkippable` | Pending |
| V-P03 | `gitmap/cmd/agy_scan.go` | 83 | `shouldSkipScanDir` | `isSkippableScanDir` | Pending |
| V-P04 | `gitmap/cmd/audit.go` | 39 | `shouldAuditCommand` | `isAuditableCommand` | Pending |
| V-P05 | `gitmap/cmd/clone.go` | 127 | `shouldUseMultiClone` | `isMultiCloneEnabled` | Pending |
| V-P06 | `gitmap/cmd/clonenextbatchdispatch.go` | 31 | `shouldRunBatch` | `isBatchRunnable` | Pending |
| V-P07 | `gitmap/cmd/inject_idempotency.go` | 74 | `shouldRunDesktop` | `isDesktopRunnable` | Pending |
| V-P08 | `gitmap/cmd/inject_idempotency.go` | 89 | `shouldRunVSCode` | `isVSCodeRunnable` | Pending |
| V-P09 | `gitmap/cmd/installer_ls.go` | 45 | `shouldSkipInstallerRow` | `isSkippableInstallerRow` | Pending |
| V-P10 | `gitmap/cmd/install_add.go` | 64 | `shouldPromptInteractively` | `isInteractivePromptNeeded` | Pending |
| V-P11 | `gitmap/cmd/install_prompt.go` | 9 | `shouldProceedInstall` | `isInstallApproved` | Pending |
| V-P12 | `gitmap/cmd/prune.go` | 48 | `shouldProceedWithPrune` | `isPruneConfirmed` | Pending |
| V-P13 | `gitmap/cmd/pull.go` | 147 | `shouldPullCWD` | `isPullCWDEnabled` | Pending |
| V-P14 | `gitmap/cmd/push.go` | 95 | `shouldPushCWD` | `isPushCWDEnabled` | Pending |
| V-P15 | `gitmap/cmd/releaseautobump.go` | 67 | `shouldAutoBumpMinor` | `isAutoBumpEligible` | Pending |
| V-P16 | `gitmap/cmd/root.go` | 368 | `shouldRewriteToClone` | `isCloneRewriteRequired` | Pending |
| V-P17 | `gitmap/cmd/scanresolve.go` | 81 | `shouldAnnounceResolve` | `isAnnounceResolveNeeded` | Pending |
| V-P18 | `gitmap/cmd/schedule_export.go` | 123 | `shouldSkipSchedule` | `isSkippableSchedule` | Pending |
| V-P19 | `gitmap/cmd/schemaregistry_assert_test.go` | 129 | `shouldUpdateSchema` | `isSchemaUpdateNeeded` | Pending |
| V-P20 | `gitmap/cmd/selfuninstall.go` | 121 | `shouldHandoffSelfUninstall` | `isSelfUninstallHandoffNeeded` | Pending |
| V-P21 | `gitmap/cmd/seowritegit.go` | 155 | `shouldStop` | `isStopRequested` | Pending |
| V-P22 | `gitmap/cmd/updaterepo.go` | 111 | `canPromptForRepoPath` | `isRepoPathPromptAllowed` | Pending |
| V-P23 | `gitmap/committransfer/pr_engine.go` | 41 | `shouldCreatePR` | `isPRCreationEnabled` | Pending |
| V-P24 | `gitmap/committransfer/replay.go` | 144 | `shouldSkipPath` | `isSkippablePath` | Pending |
| V-P25 | `gitmap/detector/detector.go` | 127 | `shouldExcludeDir` | `isExcludedDir` | Pending |
| V-P26 | `gitmap/diff/tree.go` | 87 | `shouldIgnore` | `isIgnoredPath` | Pending |
| V-P27 | `gitmap/release/selfrelease_resolve.go` | 93 | `canPromptForPath` | `isPathPromptAllowed` | Pending |
| V-P28 | `gitmap/vscodepm/optimize.go` | 37 | `shouldWrite` | `isWriteEnabled` | Pending |

---

## 5. Subtask Decomposition

- **Subtask 01:** `.lovable/plans/subtasks/96-booleans-and-complex-conditions-audit/01-task-inverted-success-and-failure-methods.md`
  - Target files: `gitmap/model/record.go`, `gitmap/model/prompt_status.go`, `gitmap/pipelinedb/pipeline_repository_test.go`, `gitmap/lazyregex/compile_result.go`, `gitmap/lazyregex/lazyregex_test.go`, `gitmap/result/result_test.go`, `gitmap/cmd/install_audit_runner_test.go`, `gitmap/cmd/prompt_failure_reporter.go`, `gitmap/cmd/schedule_cmd.go`, `gitmap/cmd/pull.go`, `gitmap/cmd/pullparallel.go`, `gitmap/cmd/push.go`, `gitmap/cmd/pushparallel.go`, `gitmap/cmd/install_logs.go`, `gitmap/tests/cmd_test/prompt_install_e2e_test.go`.
  - Actions: Add `IsFailed() bool` / `IsFail() bool` methods, replace all 17 `!.*IsSuccess` instances, decompose mixed polarity conditions.
- **Subtask 02:** `.lovable/plans/subtasks/96-booleans-and-complex-conditions-audit/02-task-banned-prefix-function-refactoring-and-linter.md`
  - Target files: `gitmap/cliexit/`, `gitmap/clonefrom/`, `gitmap/cmd/`, `gitmap/committransfer/`, `gitmap/detector/`, `gitmap/diff/`, `gitmap/release/`, `gitmap/vscodepm/`, `linter-scripts/check-boolean-guidelines.py`, `03-ai-scripts/01-index.md`.
  - Actions: Rename all 28 `can*` / `should*` functions and callers to `is*` / `has*`. Update `check-boolean-guidelines.py` to AST-enforce affirmative prefixes (`is`/`has` only) and case-insensitive inverted success checks.

---

## 6. Verification Plan

1. **Boolean Linter:** `python linter-scripts/check-boolean-guidelines.py` exits 0 with zero violations.
2. **Go Test Suite:** `go test -v -short ./...` passes with zero failures.
3. **CI/CD Quality Gate:** `python 03-ai-scripts/06-cicd-local-runner.py --filter "Boolean"` passes 100% green.
