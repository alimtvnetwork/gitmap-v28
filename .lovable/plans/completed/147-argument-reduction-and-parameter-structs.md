# Plan 147: Argument Reduction, Parameter Structs & Return Architecture Audit

Trigger Keywords & Aliases: `cg-argument-reduction`, `cg-params`, `cg-struct-params`, `cg-execute params`, `audit function arguments`, `reduce arguments`, `struct parameters`, `mandatory appfault return`, `parameter objects`, `no void functions`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)
> **Status:** COMPLETED

---

## 1. Executive Summary & Objective

Autonomously scan, discover, plan, refactor, and format all function signatures across the codebase, enforcing argument reduction via dedicated value-based parameter Structs/DTOs for signatures with >2–3 parameters, affirmative boolean prefixing (`is*` and `has*` only (can, should, was, etc. are banned)) on all struct fields and parameters, mandatory `*apperror.AppError` returns (eliminating bare "void" functions in Go domain/service logic), wrapping external framework errors into `*apperror.AppError`, and single `Result[T]` return envelopes until 100% green without stopping.

Plan 147 executes surgical parameter and return refactoring across four key subsystems:
1. **Cluster Parameter Structs & Affirmative Booleans (`cli/cluster/`):** Encapsulate loose multi-parameter signatures (`ExecRestart`, `ExecShutdown`, `ExecLogoff`, `finishClusterPool`, `reportNodeExecutionResult`, `PrintPreflight`) into dedicated value-based structs (`LifecycleExecParams`, `FinishClusterPoolParams`, `NodeExecutionReportParams`, `PreflightParams`), rename non-affirmative booleans (`forceLifecycle` -> `isForceLifecycle`, `showRole` -> `isShowRole`, `allOk` -> `isAllOk`, `autoConfirm` -> `isAutoConfirm`), and convert `PrintPreflight` error returns to `*apperror.AppError`.
2. **Clonenow Parameter Structs & Idempotent Engine (`cli/clonenow/`):** Create `cli/clonenow/types.go` declaring `CloneIdempotentParams`, `ConcurrentDispatchParams`, `ConcurrentWorkerParams`, `GitCloneParams`, `ProgressWriteParams`. Encapsulate high-arity signatures in `execute_idempotent.go`, `execute_concurrent.go`, and `execute.go`.
3. **Clonefrom Parameter Structs & Hooks (`cli/clonefrom/`):** Create `cli/clonefrom/types.go` declaring `CloneFromDispatchParams`, `CloneFromWorkerParams`, `BeforeRowInvokeParams`, `ProgressWriteParams`. Encapsulate high-arity signatures in `execute_concurrent.go` and `execute.go`.
4. **Visibility & GoldenGuard Parameter Refactoring (`cli/visibility/`, `cli/goldenguard/`):** Introduce `ExclusionTokenParams` in `cli/visibility/exclude.go` and return `*apperror.AppError`. Rename non-affirmative booleans `trigger` -> `isTrigger` in `cli/goldenguard/goldenguard.go` and `determinism.go`.
5. **Automated Auditor Tooling & Quality Gates:** Author and register `03-ai-scripts/36-param-struct-auditor.py` in `03-ai-scripts/01-index.md`. Run quality linters, record modified files under lock in test inventory, consolidate Plan 147, and push atomically via SSH.

---

## 2. Authoritative Spec Citations

- `spec/02-coding-guidelines/01-cross-language/01-index.md`: Hallucination prevention, micro-tasking, strict relative paths, and Rule 9a/9b multi-line parameter formatting.
- `spec/02-coding-guidelines/01-cross-language/02-boolean-principles/02-naming-prefixes.md`: Principle 1 (`is`/`has` prefixes only) and Principle 2 (total ban on negative words).
- `spec/02-coding-guidelines/01-cross-language/10-function-naming.md`: Semantic verb and predicate prefix standards.
- `spec/02-coding-guidelines/02-canonical-size-tier.md`: File and function size limits (functions <= 8–15 lines, files <= 100–300 lines).
- `spec/03-error-manage/01-index.md`: Universal `*AppError` wrapping, `Result[T]` envelopes, and zero swallowed errors.
- `.lovable/coding-guidelines.md`: Master consolidated coding guidelines.

---

## 3. Exhaustive Parameter & Return Ledger

| Id | File | Line | Param Count | Current Signature | Violation | Target Refactoring | Status |
|:---|:---|:---:|:---:|---|---|---|:---:|
| V-01 | `cli/cluster/exec_lifecycle.go` | 19 | 4 | `ExecRestart(ctx, node, forceLifecycle, providedPassword)` | 4 loose params, unprefixed bool `forceLifecycle` | Encapsulate into `LifecycleExecParams`, return `*apperror.AppError` | Completed |
| V-02 | `cli/cluster/exec_lifecycle.go` | 42 | 4 | `ExecShutdown(ctx, node, forceLifecycle, providedPassword)` | 4 loose params, unprefixed bool `forceLifecycle` | Encapsulate into `LifecycleExecParams`, return `*apperror.AppError` | Completed |
| V-03 | `cli/cluster/exec_lifecycle.go` | 65 | 4 | `ExecLogoff(ctx, node, forceLifecycle, providedPassword)` | 4 loose params, unprefixed bool `forceLifecycle` | Encapsulate into `LifecycleExecParams`, return `*apperror.AppError` | Completed |
| V-04 | `cli/cluster/exec_lifecycle.go` | 86 | 3 | `checkLifecycleGuards(node, forceLifecycle, providedPassword)` | Unprefixed bool `forceLifecycle`, returns raw `error` | Rename `isForceLifecycle bool`, return `*apperror.AppError` | Completed |
| V-05 | `cli/cluster/pool.go` | 219 | 7 | `reportNodeExecutionResult(spinner, nodeLabel, displayCmd, durMs, exitCode, ctxErr, allOk)` | 7 loose params, unprefixed bool `allOk` | Encapsulate into `NodeExecutionReportParams` (`IsAllOk bool`) | Completed |
| V-06 | `cli/cluster/pool.go` | 275 | 8 | `finishClusterPool(multi, isMultiActive, updateCounts, runId, totalNodes, succeeded, failed, skipped)` | 8 loose params | Encapsulate into `FinishClusterPoolParams` | Completed |
| V-07 | `cli/cluster/preflight.go` | 14 | 5 | `PrintPreflight(selector, effective, command, runRef, autoConfirm)` | 5 loose params, unprefixed bool `autoConfirm`, raw `error` | Encapsulate into `PreflightParams`, return `*apperror.AppError` | Completed |
| V-08 | `cli/cluster/lsrender.go` | 9 | 2 | `RenderNodeTable(nodes, showRole)` | Unprefixed bool `showRole` | Rename to `isShowRole bool` | Completed |
| V-09 | `cli/cmd/clustercommand.go` | 172 | 5 | `performPreflight(flags, selector, effective, cmdStr, runRef)` | 5 loose params, calls `PrintPreflight` | Modernize to pass `PreflightParams` and check `*apperror.AppError` | Completed |
| V-10 | `cli/clonenow/types.go` | 1 | N/A | Missing `types.go` | Parameter structs scattered | Create `cli/clonenow/types.go` declaring parameter structs | Completed |
| V-11 | `cli/clonenow/execute_idempotent.go` | 105 | 6 | `dispatchOnExists(r, url, absDest, cwd, policy, state)` | 6 loose params | Encapsulate into `CloneIdempotentParams` | Completed |
| V-12 | `cli/clonenow/execute_idempotent.go` | 140 | 4 | `cloneFresh(r, url, absDest, cwd)` | 4 loose params | Encapsulate into `CloneIdempotentParams` | Completed |
| V-13 | `cli/clonenow/execute_idempotent.go` | 275 | 4 | `forceReclone(r, url, absDest, cwd)` | 4 loose params | Encapsulate into `CloneIdempotentParams` | Completed |
| V-14 | `cli/clonenow/execute_concurrent.go` | 66 | 5 | `dispatchConcurrent(plan, cwd, beforeRow, workers, out)` | 5 loose params | Encapsulate into `ConcurrentDispatchParams` | Completed |
| V-15 | `cli/clonenow/execute_concurrent.go` | 82 | 5 | `runConcurrentWorker(jobs, plan, cwd, out, wg)` | 5 loose params | Encapsulate into `ConcurrentWorkerParams` | Completed |
| V-16 | `cli/clonenow/execute.go` | 133 | 4 | `runGitClone(r, url, dest, cwd)` | 4 loose params | Encapsulate into `GitCloneParams` | Completed |
| V-17 | `cli/clonenow/execute.go` | 199 | 4 | `writeProgress(w, n, total, res)` | 4 loose params | Encapsulate into `ProgressWriteParams` | Completed |
| V-18 | `cli/clonefrom/types.go` | 1 | N/A | Missing `types.go` | Parameter structs scattered | Create `cli/clonefrom/types.go` declaring parameter structs | Completed |
| V-19 | `cli/clonefrom/execute_concurrent.go` | 59 | 5 | `dispatchConcurrent(plan, cwd, beforeRow, workers, out)` | 5 loose params | Encapsulate into `CloneFromDispatchParams` | Completed |
| V-20 | `cli/clonefrom/execute_concurrent.go` | 74 | 4 | `runConcurrentWorker(jobs, cwd, out, wg)` | 4 loose params | Encapsulate into `CloneFromWorkerParams` | Completed |
| V-21 | `cli/clonefrom/execute_concurrent.go` | 96 | 4 | `invokeBeforeRow(hook, i, total, r)` | 4 loose params | Encapsulate into `BeforeRowInvokeParams` | Completed |
| V-22 | `cli/clonefrom/execute.go` | 305 | 4 | `writeProgress(w, n, total, res)` | 4 loose params | Encapsulate into `ProgressWriteParams` | Completed |
| V-23 | `cli/visibility/exclude.go` | 62 | 4 | `absorbExclusionToken(tok, tokIdx, totalCount, out)` | 4 loose params, raw stdlib `error` | Encapsulate into `ExclusionTokenParams`, return `*apperror.AppError` | Completed |
| V-24 | `cli/visibility/exclude.go` | 86 | 4 | `absorbExclusionRange(tok, tokIdx, totalCount, out)` | 4 loose params, raw stdlib `error` | Encapsulate into `ExclusionTokenParams`, return `*apperror.AppError` | Completed |
| V-25 | `cli/goldenguard/goldenguard.go` | 64 | 2 | `AllowUpdate(t, trigger)` | Unprefixed bool `trigger` | Rename to `isTrigger bool` | Completed |
| V-26 | `cli/goldenguard/determinism.go` | 149 | 4 | `AllowUpdateAfterDeterminism(t, trigger, label, writer)` | Unprefixed bool `trigger` | Rename to `isTrigger bool` | Completed |
| V-27 | `03-ai-scripts/36-param-struct-auditor.py` | 1 | N/A | Missing auditor script | No dedicated parameter struct auditor | Author and register `36-param-struct-auditor.py` | Completed |
