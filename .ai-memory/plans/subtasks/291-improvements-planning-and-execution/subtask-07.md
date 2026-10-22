# Subtask 07: Pipeline Log Stream Closing & Pipe Error Handling

> **Subtask ID:** Subtask-07  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-07.md`  
> **Status:** PENDING  

---

## 1. Objectives & Scope

1. Audit and refactor `cli/cmdpipeline/pipeline_logs_cmd.go` and its companion modules (`pipeline_logs_output.go`, `pipeline_logs_fetch.go`) to guarantee proper error handling on stream closures, pipe writes, and file exports.
2. Eliminate unhandled error discards (`_ = ...`) on pipe writer closers and file synchronizations.
3. Propagate stream errors back to the caller instead of silently truncating log output.
4. Enforce strict coding guidelines:
   - Functions strictly $\le 15$ lines ($\le 8$ lines preferred).
   - Zero nested conditionals.
   - Positive boolean naming (`isDetailed`, `isJSON`, `hasFailure`, `isFromCache`).
   - Blank line before every `return` and after `}`.

---

## 2. Detailed Technical Plan

### 2.1 Refactor Stream & Pipe Operations in `pipeline_logs_cmd.go`
- Locate all file writers, tempfile handlers, and pipe transfers in `pipeline_logs_cmd.go`:
  - When writing formatted logs to disk or piping into clipboard utilities, capture write and close errors:
    ```go
    if err := writer.Close(); err != nil {
        return apperror.WrapSimple(err, "closeLogStream")
    }
    ```
  - Replace unchecked `defer file.Close()` with safe deferred closing helper that logs or records closure errors if the main return was successful.

### 2.2 Refactor Output Generation in `pipeline_logs_output.go`
- In `renderAndCopyTerminalReport(p PipelineErrorLogsPayload)`:
  - Check clipboard write status; if clipboard fails (e.g., in headless CI or SSH session without display), fallback gracefully to terminal display without crashing.
  - Decompose rendering routines exceeding 15 lines into concise child helpers.

### 2.3 Style & Guideline Compliance
- Verify all boolean options in `PipelineLogsOptions` conform to positive naming:
  - `FailedOnly` -> `isFailedOnly` (or refactor to `hasFailedOnly`).
  - `IsDetailed`, `IsJSON`, `DoCopy` -> strictly positive `is*` or `has*`.
- Ensure clean spacing before all `return` statements and after closing brace blocks.

---

## 3. Verification Criteria

- [ ] All stream closers and pipe writers in `cli/cmdpipeline/pipeline_logs_cmd.go` handle errors explicitly.
- [ ] No ignored `_ = ...` assignments exist on critical file/pipe IO calls.
- [ ] Every function in `pipeline_logs_cmd.go` and `pipeline_logs_output.go` is $\le 15$ lines.
- [ ] All booleans adhere to positive naming (`is*`, `has*`).
- [ ] Blank line precedes every `return` and follows every `}` block.
