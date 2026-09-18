# CI/CD Issue 50: Pipeline Fallback ResultSlice Methods and Slice Type Drift RCA

## 1. Symptom
In GitHub Actions CI run `#35259849192` on commit `547573e`, 19 pipeline jobs failed during `Go vet`, `typecheck`, and `go build`:

### Error Output
```text
cmdpipeline/pipeline_fallback.go:30:12: runRes.IsFail undefined (type pipelinedb.PipelineRunSliceResult has no field or method IsFail)
cmdpipeline/pipeline_fallback.go:42:13: errRes.IsFail undefined (type pipelinedb.PipelineErrorSliceResult has no field or method IsFail)
cmdpipeline/pipeline_fallback.go:47:17: compactRes.IsFail undefined (type pipelinedb.PipelineCompactErrorSliceResult has no field or method IsFail)
cmdpipeline/pipeline_fallback.go:88:5: p.Notes undefined (type *PipelineErrorLogsPayload has no field or method Notes)
cmdpipeline/pipeline_fallback.go:110:4: p.Notes undefined (type *PipelineErrorLogsPayload has no field or method Notes)
cmdpipeline/pipeline_logs_cmd.go:114:58: cannot use groups (variable of type []CommitPipelineGroup) as []*CommitPipelineGroup value in argument to resolveTargetGroupWithFallback
cmdpipeline/pipeline_logs_cmd.go:126:48: cannot use groups (variable of type []*CommitPipelineGroup) as []CommitPipelineGroup value in argument to ResolveCommitGroupByTarget
cmdpipeline/pipeline_logs_cmd.go:150:30: cannot use items (variable of type []CommitWorkflowItem) as []CommitWorkflowLogItem value in argument to assembleWorkflowLogs
cmdpipeline/pipeline_logs_cmd.go:150:38: cannot use items (variable of type []CommitWorkflowItem) as []CommitWorkflowLogItem value in return statement
cmdpipeline/pipeline_logs_cmd.go:166:10: cannot use []CommitWorkflowLogItem{…} (value of type []CommitWorkflowLogItem) as []CommitWorkflowItem value in return statement
```

---

## 2. Root Cause
1. **`ResultSlice` Method Call Mismatch**: `pipelinedb.PipelineRunSliceResult` is a type alias for `result.ResultSlice[PipelineRunRecord]`. The failure checking predicate is a method call `IsFailed()`, but was invoked as a struct field `.IsFail`.
2. **Missing `Notes` Field**: `PipelineErrorLogsPayload` in `cli/cmdpipeline/pipeline.go` did not define the `Notes string` field assigned during previous run fallback.
3. **Commit Groups Slice Pointer Mismatch**: `GroupRunsByCommit` returns concrete value slice `[]CommitPipelineGroup`. The fallback helper was typed with pointer slice `[]*CommitPipelineGroup`.
4. **Function Return Type Typo**: `fetchTargetWorkflows` was declared returning `[]CommitWorkflowItem` instead of `[]CommitWorkflowLogItem`.

---

## 3. Resolution
1. Replaced all `.IsFail` accesses with `.IsFailed()` on `runRes`, `errRes`, and `compactRes` in `cli/cmdpipeline/pipeline_fallback.go`.
2. Added `Notes string json:"notes,omitempty"` to `PipelineErrorLogsPayload` in `cli/cmdpipeline/pipeline.go`.
3. Standardized parameter types of `ResolveFallbackTargetGroup` and `resolveTargetGroupWithFallback` to `[]CommitPipelineGroup`.
4. Corrected the return type of `fetchTargetWorkflows` to `[]CommitWorkflowLogItem` in `cli/cmdpipeline/pipeline_logs_cmd.go`.
5. Updated `cli/cmdpipeline/pipeline_fallback_test.go` fixture to instantiate `[]CommitPipelineGroup`.

---

## 4. Prevention & Learnings
- **Result Type Envelopes**: Always invoke `.IsFailed()`, `.IsSuccess()`, or `.HasError()` as methods when inspecting `Result[T]` or `ResultSlice[T]` envelopes.
- **Concrete Slice Signatures**: Validate matching slice element types (`[]T` vs `[]*T`) against domain constructors (`GroupRunsByCommit`).
