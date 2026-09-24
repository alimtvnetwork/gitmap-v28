# Spec 153: Dynamic Column Width Alignment for Efficient Pull (`gitmap pae`)

> **Spec ID:** `SPEC-153`  
> **Version:** `v6.331.0`  
> **Status:** Completed / Released  
> **Date:** 2026-09-24  

---

## 0. User Request (Verbatim) & Actionable Deliverables

### Request 1:
```text
fix the output column width for gitmap pae https://prnt.sc/ZviMUeV9slNP please
```

### Request 2:
```text
imrpove the column display for gitmap pae, can you please fix it properly
```

### Problem Statement:
Screenshot telemetry at `https://prnt.sc/ZviMUeV9slNP` and `media_1790261308371.png` exhibited two display flaws in `gitmap pae` (`gitmap pull all-efficient`):
1. **Active Repositories Column Misalignment**:
   In `cli/cmdpull/pull_efficient.go`, `renderConciseActiveResults` hardcoded `%-26s`. Repositories whose name exceeded 26 characters (e.g. `bsrm-presentation-hiltrax-v4` [29 chars], `riseup-asia-website-project-v6` [30 chars], `kita-social-media-content-calender-v2` [37 chars]) overflowed the 26-char boundary, displacing the status column horizontally.
2. **Inactive Repositories Word-Wrap Splitting**:
   In `printInactiveSkipSummary`, 55 inactive repository names were printed as a single comma-separated line with `strings.Join(names, ", ")`, causing the terminal to hard-wrap mid-word (`bsrm-presentation-hiltr` ... `ax-v4`, `gitlo` ... `gger-new-v2`, `movie-cli-` ... `v8`).
3. **Local Deployment Out-of-Date**:
   The user's local binary at `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` was running `v6.325.0`.

### Extracted Actionable Deliverables:
1. **Dynamic Column Width Calculation**:
   - `ResolveConciseRepoColWidth(states []*PullRepoState) int`:
     - Default baseline width: `26`.
     - Dynamically expands to fit the widest active repository name (e.g. 37+ chars).
     - Bounded by `detectTerminalWidth() - 30` to prevent unintended line wrapping.
2. **Vivid Status Syntax Highlighting**:
   - `StyleRepoStatusLabel(statusLabel string) string`:
     - `up-to-date`: bright bold green (`ColorGreen`)
     - `dirty`: bright bold yellow (`ColorYellow`)
     - `failed`: bright bold red (`ColorRed`)
     - Diff stats (e.g. `+576/-119 (32)`): `+576` in green, `/-119` in red, `(32)` in dim white.
3. **Clean Comma-Bound Word Wrapping for Inactive Repositories**:
   - `FormatWrappedInactiveList(names []string, indent string, maxLineLen int) string`:
     - Wraps lines at comma boundaries before reaching `maxLineLen`.
     - Preserves clean 6-space indentation (`      `) across all wrapped lines.
     - Never splits a repository name across lines.
4. **Local Deployment Synchronization**:
   - Compile and deploy the updated binary to `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`.
5. **Quality Gates & Tests**:
   - 6 unit tests in `cli/cmdpull/pull_efficient_render_test.go`.
   - Isolated temporary E2E test in `cli/tests/e2e/pae_column_width_tempe2e_test.go` (`//go:build tempe2e`, `RUN_TEMP_E2E=1`).
6. **Release & Push**:
   - Tag `v6.331.0` and push to `origin main`.
