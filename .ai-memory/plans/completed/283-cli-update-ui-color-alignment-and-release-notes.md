# Plan: 283-cli-update-ui-color-alignment-and-release-notes

## Metadata
- **Task ID**: `283-cli-update-ui-color-alignment-and-release-notes`
- **Specification**: `02-spec/21-app/283-cli-update-ui-color-alignment-and-release-notes/spec.md`
- **Total Steps**: 2 Subtasks across 1 Wave (A=2, H=2)
- **Status**: IN_PROGRESS

---

## User Request (Verbatim)

```text
# High Priority Instruction

Can you please fix the coloring for the installation script? What do I mean by that is the blue color with black, it just faded in. So you could have a light color. For example, the right-hand side, when we change it to the different color, it could have been the yellow highlighter color, and the current one can be white. And the arrow color can be yellow as well. Yellow to yellow. So that we understand where to where we have converted. And also at the end, I wanted to see some of the release note, one or two lines. And the success, every time the success shows, there should be a little bit of padding from top. And left-hand side, the success needs to match with the number. For example, where the number actually stands, the four or five. So in that way, the success needs to be in that same row. And the success when the first line starts, be consecutive together, no need to have a gap. But when it starts the first time, it needs to have some gap so that it looks professional. The same thing should happen when the success thing ends. There needs to be a little bit of gap. Make sure that you have this, then we have a perfect and high-quality installer script.
```

---

## Subtask Breakdown

### Subtask-01: Update UI Colors, Spacing, Indentation & Duplicate Flush Prevention
- **Owner**: Worker 01
- **Files Owned**:
  - `cli/cmdupdate/update_ui.go`
  - `cli/cmdupdate/update_ui_notes.go` (new helper)
- **Actions**:
  1. In `update_ui.go`, update version transition: `prevVer` to bold white (`#f8f8f2`), `arrow` to yellow (`#f1fa8c`), `currVer` to bright yellow highlighter (`#ffff55`).
  2. In `update_ui.go`, format `OnStageSuccess` with 4 spaces left indentation (matching column of `[4/5]`), top padding on first success badge, consecutive badges grouped without gaps, and bottom padding before summary card.
  3. In `update_ui.go`, guard `RenderSummaryCard()` and `RenderAlreadyUpdated()` with `r.isClosed = true` to prevent deferred `Close()` from flushing duplicate buffers.
  4. In `update_ui_notes.go`, implement helper `ExtractReleaseNotesSummary` and format release notes bullets in summary card.
  5. Enforce line count <= 300 on `cli/cmdupdate/update_ui.go`.

### Subtask-02: Release Notes Propagation in Update Data Pipeline
- **Owner**: Worker 02
- **Files Owned**:
  - `cli/cmdupdate/update_fallback_loop.go`
  - `cli/cmdupdate/types.go`
  - `cli/cmdupdate/update_binary_swap.go`
- **Actions**:
  1. In `update_fallback_loop.go`, add `Body string `json:"body"`` to `GitHubRelease`, pass `Body` into `ReleaseCandidate.ReleaseNotes`.
  2. In `types.go`, add `ReleaseNotes []string `json:"release_notes,omitempty"`` to `UpdateResult`.
  3. In `update_binary_swap.go`, extract 1–2 bullet lines from `ReleaseCandidate.ReleaseNotes` via `ExtractReleaseNotesSummary` and populate `res.ReleaseNotes` in `finalizeUpdate()`.
  4. Verify line counts <= 300 on all edited files.
