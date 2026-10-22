# Migration & Verification Specification 295: Monadic Wrapped Result Architecture

## Status: APPROVED
- **Spec ID:** 295-03
- **Slug:** `295-cg-execute-wrapped-error-remediation`
- **Domain:** Verification & Quality Gates

---

## 1. Migration Strategy

1. **Zero External Breakages:**
   Existing legacy methods like `db.ListRepos()` (with 60+ callers) are preserved to respect the micro-batch rule.
2. **Targeted Refactoring:**
   The repository methods (`RepoRepository` and its delegators `ListReposUnpushed`, `ListReposByTag`, `ListReposByBranch`, `ListReposDirtyStatus`, `FindPathBySlug`) are transformed into canonical monadic wrapped containers (`ScanRecordSliceWrap`, `ScanRecordWrap`, `RepoPathWrap`) declared in `cli/store/repository_types.go`.
3. **Consumer Updates:**
   `cli/cmdpipeline/pipeline_history_render.go` updated to check `pathRes.HasError()` and read `pathRes.Value`.

---

## 2. Verification Protocol

- **Zero-Build & Zero-Test Invariant:** No compiler commands (`go build`) or test executions (`go test`) are permitted.
- **FastGate Validation:** `python 03-ai-scripts/50-fastgate.py --files ...` must pass 100% green.
- **AST Inspections:**
  - Function line limit: $\le 15$ lines ($\le 8$ lines preferred).
  - File line limit: $\le 300$ lines for new files, $\le 500$ lines for legacy files.
  - Zero nested ifs.
  - Positive booleans (`is*`, `has*` only).
  - Vertical spacing: blank line before `if`, blank line after `}`, blank line before `return`.
