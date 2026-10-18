# Subtask 02: Fix Progress Bar Column Width & Waterfall Wrapping Bug

**Parent Task:** `gitmap-and-antigravity-manager-agm-package-enhancements`  
**Subtask ID:** `subtask-02`  
**Status:** `IN_PROGRESS`  
**Owned Files:**  
- `cli/cmdpull/pull_progress_bar_render.go`
- `cli/cmdpull/pull_progress_bar_visual.go`

## Objective
Fix the defect shown in `assets/screenshots/gitmap-and-antigravity-manager-agm-package-enhancements-01.png` where `gitmap pa` / `pull-all` prints multiple progress bars stacking into a staircase/waterfall on each update, and screenshot 02 where non-TTY background tasks spam newlines on every git progress tick.

## Implementation Details
1. In `cli/cmdpull/pull_progress_bar_render.go`:
   - Replace `utf8.RuneCountInString(plain)` with `runewidth.StringWidth(plain)` using `github.com/mattn/go-runewidth`.
   - In `truncateANSI`, correctly increment visible character width by `runewidth.RuneWidth(r)` rather than `1`, accounting for wide Unicode runes and emojis (`⚡`, `🔍`, `🔀`).
   - Introduce a defensive right-margin: `maxLen := termWidth - 4` (clamped to a minimum of 20). This ensures the terminal cursor never touches the last column $W$, preventing ConHost / Terminal emulator auto-wrap to row Y+1.
2. In `cli/cmdpull/pull_progress_bar_visual.go`:
   - In `renderNonTTY()`, suppress intermediate milestone ticks (counting/compressing) so background runners do not spam newlines.
