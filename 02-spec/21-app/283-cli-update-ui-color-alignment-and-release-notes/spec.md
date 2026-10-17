# Specification: CLI Update UI Color Alignment, Spacing & Release Notes

## Status
- **Task ID**: `283-cli-update-ui-color-alignment-and-release-notes`
- **Specification Version**: 1.0.0
- **Status**: APPROVED
- **Date**: 2026-10-17

---

## 1. User Request (Verbatim)

```text
# High Priority Instruction

Can you please fix the coloring for the installation script? What do I mean by that is the blue color with black, it just faded in. So you could have a light color. For example, the right-hand side, when we change it to the different color, it could have been the yellow highlighter color, and the current one can be white. And the arrow color can be yellow as well. Yellow to yellow. So that we understand where to where we have converted. And also at the end, I wanted to see some of the release note, one or two lines. And the success, every time the success shows, there should be a little bit of padding from top. And left-hand side, the success needs to match with the number. For example, where the number actually stands, the four or five. So in that way, the success needs to be in that same row. And the success when the first line starts, be consecutive together, no need to have a gap. But when it starts the first time, it needs to have some gap so that it looks professional. The same thing should happen when the success thing ends. There needs to be a little bit of gap. Make sure that you have this, then we have a perfect and high-quality installer script.
```

---

## 2. Problem Statement & Root Cause Analysis

### 2.1 Low Contrast / Faded Blue on Black
- In `cli/cmdupdate/update_ui.go`, the previous version string was rendered using `DraculaGray = "#6272a4"` (`RGB 98, 114, 164`). On standard dark terminal themes (PowerShell, Command Prompt, macOS/Linux dark xterms), `#6272a4` has heavy blue undertones that render dark and muddy against black backgrounds, making the previous version barely legible.
- The arrow `→` was styled with `DraculaPurple = "#bd93f9"`.
- The current version was styled with `DraculaGreen = "#50fa7b"`.
- The user requested:
  - Previous version: **White** (clean, crisp, readable foreground `#f8f8f2` / white).
  - Arrow: **Yellow** (`#f1fa8c` / yellow).
  - Target/Current version: **Yellow Highlighter** (`#ffff55` / bright yellow highlighter).
  - This establishes a clear visual semantic transition: white source version converted through a yellow arrow to a highlighted yellow target version.

### 2.2 Release Notes Absence
- When users run `gitmap update`, the summary card confirms the binary update but does not display what changed in the new version.
- GitHub Release API responses (`GET /repos/{owner}/{repo}/releases/tags/{tag}`) already return the `"body"` markdown field containing release bullets.
- The Go struct `GitHubRelease` in `cli/cmdupdate/update_fallback_loop.go` dropped this field.
- By parsing the first 1–2 bullet lines from `GitHubRelease.Body` and propagating them to `UpdateResult`, the CLI can render a clean 1–2 line release note summary in the card.

### 2.3 `SUCCESS` Badge Misalignment & Missing Padding
- Step progress lines (`■ [4/5] Applying binary handoff...`) indent the step tag `[4/5]` by 4 columns (2-space margin + 2-column spinner symbol).
- `pterm.Success` badges were printed flush-left at column 0 (`SUCCESS Target release resolved...`), creating an awkward visual misalignment with the step number above.
- There was no top margin before the first `SUCCESS` badge, and no bottom margin after the last `SUCCESS` badge before the summary card.
- Indenting `SUCCESS` badges by 4 spaces ensures the badge aligns directly under `[4/5]`.
- Adding an initial blank line before the first badge and a trailing blank line after the last badge provides professional visual breathing room, while keeping intermediate badges consecutive without extra gaps.

### 2.4 Duplicate Post-Card Execution Flushes
- In `cli/cmdupdate/update_binary_swap.go`, `runBinaryReleaseUpdate` defers `ui.Close()`.
- `RenderSummaryCard()` in `cli/cmdupdate/update_ui.go` called `r.stopPrinters()` without setting `r.isClosed = true`.
- When the function returned and `ui.Close()` fired, `r.stopPrinters()` ran a second time.
- Calling `r.multi.Stop()` again caused `pterm.MultiPrinter` to re-flush its entire buffer to `stdout` below the summary card.
- Guarding `RenderSummaryCard()`, `RenderAlreadyUpdated()`, and `Close()` with `r.isClosed = true` and an atomic stop flag eliminates this duplicate rendering.

---

## 3. Architecture & Functional Requirements

### 3.1 Color Palette & Version Transition
1. Define / update color styles in `cli/cmdupdate/update_ui.go`:
   - `prevVer`: `lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f8f8f2"))` (Bold White).
   - `arrow`: `lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1fa8c"))` (Yellow).
   - `currVer`: `lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffff55"))` (Yellow Highlighter).
2. Card border: keep crisp `lipgloss.RoundedBorder()` with clean readable border foreground (`DraculaPurple` / `DraculaCyan`).
3. Card title: `✔ Gitmap CLI Updated Successfully` in bold bright green.

### 3.2 Release Notes Extraction & Rendering
1. In `cli/cmdupdate/update_fallback_loop.go`:
   - Add `Body string `json:"body"`` to `GitHubRelease`.
   - Propagate `ReleaseNotes: ghRel.Body` into `ReleaseCandidate`.
2. In `cli/cmdupdate/types.go`:
   - Add `ReleaseNotes []string `json:"release_notes,omitempty"`` to `UpdateResult`.
3. In `cli/cmdupdate/update_ui_notes.go` (new helper file to maintain <= 300 lines limit on `update_ui.go`):
   - Implement `ExtractReleaseNotesSummary(body string, maxLines int) []string`.
   - Strip markdown headers (`#`, `##`, `###`), horizontal rules (`---`), and extract list items starting with `- `, `* `, or `• `. Clamps long strings to 72 chars.
4. In `cli/cmdupdate/update_ui.go`:
   - In `formatSummaryLines()`, append `Release Notes:` header and bullet items if present in `result.ReleaseNotes`.
   - In `renderPlainSummary()`, format bullet lines if plain mode is active.

### 3.3 `SUCCESS` Badge Indentation & Padding Protocol
1. In `cli/cmdupdate/update_ui.go`:
   - Introduce `hasPrintedFirstSuccess bool` in `updateUIRendererImpl`.
   - In `OnStageSuccess()`:
     - If `!r.hasPrintedFirstSuccess`:
       - Output an empty line before the first badge (top padding).
       - Set `r.hasPrintedFirstSuccess = true`.
     - Left-indent the `SUCCESS` badge by 4 spaces so that `SUCCESS` aligns with column 4 (matching `[4/5]`).
     - Do NOT add blank lines between consecutive `SUCCESS` badges.
2. In `RenderSummaryCard()` / `renderLipglossSummary()`:
   - Ensure an empty line is printed after the last `SUCCESS` badge before the card box (bottom padding).
3. In `stopPrinters()` & `Close()`:
   - Mark `r.isClosed = true` in `RenderSummaryCard()` and `RenderAlreadyUpdated()`.
   - Ensure `r.stopPrinters()` is idempotent and stops the spinner with `WithRemoveWhenDone(true)`.

---

## 4. File Boundaries & Ownership

- **Worker 01**:
  - `cli/cmdupdate/update_ui.go` (card colors, version transition formatting, padding, isClosed flag)
  - `cli/cmdupdate/update_ui_notes.go` (new helper for release notes parsing and bullet formatting)
- **Worker 02**:
  - `cli/cmdupdate/update_fallback_loop.go` (add `Body` to `GitHubRelease`, pass to `ReleaseCandidate`)
  - `cli/cmdupdate/types.go` (add `ReleaseNotes` to `UpdateResult`)
  - `cli/cmdupdate/update_binary_swap.go` (populate `res.ReleaseNotes` in `finalizeUpdate`)

---

## 5. Verification & Acceptance Criteria

1. Version transition renders with previous version in bold white (`#f8f8f2`), arrow in yellow (`#f1fa8c`), and new version in bright yellow highlighter (`#ffff55`).
2. Release notes (1–2 lines) are cleanly parsed from GitHub release body and displayed in the summary card.
3. `SUCCESS` badges have an empty line padding at top, are indented by 4 spaces to match the column of `[4/5]`, remain consecutive without intervening gaps, and have an empty line gap before the summary card.
4. No duplicate step output or buffer dump appears after the summary card.
5. All modified Go files stay strictly <= 300 lines with nesting depth <= 1.
