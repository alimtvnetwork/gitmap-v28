# Subtask 17: Modernize Pipeline Failure Cards with TokyoBoxStyle and TokyoTag

## Status: Completed
## Parent Plan: .ai-memory/plans/291-improvements-planning-and-execution.md
## Component Spec: 02-spec/21-app/291-improvements-planning-and-execution/02-component-spec.md
## Target Files: `cli/cmdpipeline/pipeline_history_render.go`

---

## 1. Objective

Modernize pipeline failure cards in `cli/cmdpipeline/pipeline_history_render.go` by replacing raw ASCII box borders with `termout.TokyoBoxStyle` and structuring run metadata with `termout.TokyoTag` pills for unified visual rhythm.

---

## 2. Scope & Implementation Details

1. **Refactor `renderSingleCachedFailureCard`**:
   - Current implementation prints manual border characters (`┌─`, `│`, `└─`).
   - Refactor to assemble the card content in a `strings.Builder` and wrap the result in `termout.TokyoBoxStyle(content)`.

2. **Card Structure**:
   - **Header Line**:
     - `termout.TokyoTag("FAIL", fmt.Sprintf("[%d/%d] #%d", idx, total, r.RunId))`
     - Followed by workflow name: `SafeTokyoStyle(TokyoPrimary).Bold(true).Render(r.WorkflowName)`.
   - **Metadata Line**:
     - `termout.TokyoTag("Branch", r.Branch)`
     - `termout.TokyoTag("SHA", r.Sha)`
     - `termout.TokyoTag("Duration", formatDurationSeconds(r.DurationSeconds))`
     - `termout.TokyoTag("Date", formatRunTimestamp(r.CreatedAt))`
   - **Error Diagnostics Block**:
     - Call existing diagnostic formatters (`renderRunErrorsByDetailMode(db, r.RunId, isDetailed)`) or buffer into card.
     - Steps marked with `termout.TokyoErrorPill("STEP")` and error details indented.

3. **Function Decomposition**:
   - `buildFailureCardHeader(r pipelinedb.PipelineRunRecord, idx, total int) string`
   - `buildFailureCardMeta(r pipelinedb.PipelineRunRecord) string`
   - `renderSingleCachedFailureCard(db *pipelinedb.PipelineSplitDb, r pipelinedb.PipelineRunRecord, idx, total int, isDetailed bool)`
   - All helper functions strictly $\le 15$ lines.

---

## 3. Strict Coding Guidelines & Invariants

- **Function Length Limit**: Every function must be $\le 15$ lines of executable code.
- **Zero Nested Conditionals**: Invert conditions and use early returns.
- **Positive Booleans**: Strict use of `isDetailed`.
- **Return Styling**: Blank line after `}` and blank line before `return`.
- **Zero Build/Test**: Do not run `go test` or `go build`.

---

## 4. Verification Checklist

- [x] Raw ASCII border strings (`┌─`, `│`, `└─`) removed from `renderSingleCachedFailureCard`.
- [x] Card content rendered through `termout.TokyoBoxStyle`.
- [x] Run, branch, SHA, duration, and timestamp formatted with `termout.TokyoTag`.
- [x] Step names and error indicators use Tokyo pill styling.
- [x] All functions $\le 15$ lines without nested `if` statements.
