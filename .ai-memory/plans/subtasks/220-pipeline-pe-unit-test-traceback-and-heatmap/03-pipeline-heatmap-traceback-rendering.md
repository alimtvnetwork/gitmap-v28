# Subtask 03: Pipeline Heatmap & Log Cache Enhancement

> **Parent Plan:** `220-pipeline-pe-unit-test-traceback-and-heatmap`  
> **Status:** DONE  
> **Target:** `cli/cmdpipeline/pipeline_query.go`, `cli/cmdpipeline/pipeline_logs.go`, `cli/cmdpipeline/pipeline_details.go`  

---

## Objectives
1. In `pipeline_query.go`:
   - Guard against caching synthetic step fallbacks as authoritative log files when a run is still active or pending.
   - Re-query actual GitHub Actions logs if a cached file only contains synthetic fallback lines.
2. In `pipeline_logs.go` and `pipeline_details.go`:
   - Ensure the rendered heatmap card includes the failed unit test name and the complete traceback.
   - Ensure copied clipboard logs and combined logs preserve the traceback so AI agents and engineers can immediately diagnose the error.
