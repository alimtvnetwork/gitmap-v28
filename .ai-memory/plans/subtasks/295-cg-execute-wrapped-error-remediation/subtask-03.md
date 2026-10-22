# Subtask 03: Update Pipeline History Render Consumer

## Status: Completed
- **Parent Plan:** [295-cg-execute-wrapped-error-remediation.md](../../pending/295-cg-execute-wrapped-error-remediation.md)
- **Target File:** `cli/cmdpipeline/pipeline_history_render.go`
- **Objective:** Update `queryRepoPathFromConn` to consume `result.Result[string]` from `repo.FindPathBySlug(slug)`.

---

## Action Items
1. Check `pathRes.HasError()` via guard clause.
2. If error, log error via `logRepoLookupFailure(slug, pathRes.Err)` and return `""`.
3. If success, return `pathRes.Value`.
4. Keep file strictly $\le 500$ lines.
