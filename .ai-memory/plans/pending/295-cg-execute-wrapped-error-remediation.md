# Master Plan 295: Repository-Wide Monadic Wrapped Result Architecture (CG Execute)

## Status: COMPLETED
- **Plan ID:** 295
- **Slug:** `295-cg-execute-wrapped-error-remediation`
- **Architecture Spec:** [01-architecture-spec.md](../../../02-spec/21-app/295-cg-execute-wrapped-error-remediation/01-architecture-spec.md)
- **Component Spec:** [02-component-spec.md](../../../02-spec/21-app/295-cg-execute-wrapped-error-remediation/02-component-spec.md)
- **Migration & Verification Spec:** [03-migration-and-verification-spec.md](../../../02-spec/21-app/295-cg-execute-wrapped-error-remediation/03-migration-and-verification-spec.md)
- **Scope:** Refactoring `cli/store/repo_repository.go`, `cli/store/repo_query_filter.go`, `cli/store/repo.go`, and dependent pipeline renderers from multi-value tuples `(T, *apperror.AppError)` to monadic `result.Result[T]`.
- **Execution Mandate:** Zero-Build, Zero-Test, Strictly Relative Lowercase Git Paths, Bounded 5–8 File Micro-Batches

---

## 1. Subtask Registry & Traceability Matrix

- [x] **Subtask-01:** [subtask-01.md](../subtasks/295-cg-execute-wrapped-error-remediation/subtask-01.md)
  - **Targets:** `cli/store/repo_repository.go`
  - **Objective:** Refactor all `RepoRepository` methods (`ListByBranch`, `ListDirty`, `ListByTag`, `ListUnpushed`, `FindPathBySlug`) to return `result.Result[T]`.
  - **Status:** Completed

- [x] **Subtask-02:** [subtask-02.md](../subtasks/295-cg-execute-wrapped-error-remediation/subtask-02.md)
  - **Targets:** `cli/store/repo_query_filter.go`, `cli/store/repo.go`
  - **Objective:** Refactor delegating methods (`ListReposByTag`, `ListReposUnpushed`, `ListReposByBranch`, `ListReposDirtyStatus`) to return `result.Result[[]model.ScanRecord]`.
  - **Status:** Completed

- [x] **Subtask-03:** [subtask-03.md](../subtasks/295-cg-execute-wrapped-error-remediation/subtask-03.md)
  - **Targets:** `cli/cmdpipeline/pipeline_history_render.go`
  - **Objective:** Update `queryRepoPathFromConn` to consume `result.Result[string]`.
  - **Status:** Completed

- [x] **Subtask-04:** [subtask-04.md](../subtasks/295-cg-execute-wrapped-error-remediation/subtask-04.md)
  - **Targets:** FastGate validation & gitmap cpf commit
  - **Objective:** Validate with `python 03-ai-scripts/50-fastgate.py`, verify 0 violations, and commit.
  - **Status:** Completed

---

## 2. Invariants & Quality Gates

1. **Monadic Wrapped Return:** Every fallible repository query returns `result.Result[T]`, eliminating multi-value error tuples.
2. **Guards on Error:** Callers must guard with `res.HasError()` before accessing `res.Value`.
3. **Zero Builds & Zero Tests:** Static AST checking and FastGate verification only.
4. **Sizing:** Functions $\le 15$ lines, files $\le 300$ lines for new files, $\le 500$ lines for legacy files.
