# Spec 153: Dynamic Column Width Alignment for Efficient Pull (`gitmap pae`)

> **Spec ID:** `SPEC-153`  
> **Version:** `v6.330.0`  
> **Status:** Completed / Released  
> **Date:** 2026-09-24  

---

## 0. User Request (Verbatim) & Actionable Deliverables

```text
fix the output column width for gitmap pae https://prnt.sc/ZviMUeV9slNP please
```

### Problem Statement:
Screenshot telemetry at `https://prnt.sc/ZviMUeV9slNP` exhibits column misalignment in `gitmap pae` (`gitmap pull all-efficient`).
In `cli/cmdpull/pull_efficient.go`:
```go
func renderConciseActiveResults(states []*PullRepoState) {
    fmt.Println()
    for _, s := range states {
        statusLabel := s.Changes
        if statusLabel == "" || statusLabel == "synced" {
            statusLabel = "up-to-date"
        }
        fmt.Printf("    • %-26s %s\n", s.RepoName, statusLabel)
    }
}
```
Because the column width was hardcoded to `%-26s`, any repository whose name exceeds 26 characters (e.g. `bsrm-presentation-hiltrax-v4` [29 chars], `kita-social-media-content-calender-v2` [37 chars]) overflowed the 26-char boundary, printing 0 padding and displacing the status column horizontally. Repositories with names of 26 characters (e.g. `ai-empathy-prompt-tuner-v1`) had 0 trailing spaces, leaving only 1 space before `up-to-date`.

### Extracted Actionable Deliverables:
1. **Dynamic Column Width Calculation**:
   - Introduce `resolveConciseRepoColWidth(states []*PullRepoState) int`:
     - Default baseline width: `26`.
     - Scan all `s.RepoName` in `states`.
     - If the maximum name length exceeds 26, expand `colWidth` to `max(len(s.RepoName))`.
     - Bound against detected terminal width (`detectTerminalWidth() - 30`) to avoid unwanted line wrapping in narrow terminals.
2. **Column Spacing & Alignment**:
   - Provide a clean 2-space column separator between repo names and statuses: `%-*s  %s`.
   - Ensure every status label (`up-to-date`, `dirty`, `+4941/-484 (68)`, etc.) aligns perfectly vertically across all rows.
3. **Modular Function Decomposition**:
   - Decompose into small, testable single-responsibility functions adhering to `02-spec/02-coding-guidelines/`:
     - `resolveConciseRepoColWidth(states []*PullRepoState) int`
     - `resolveRepoStatusLabel(changes string) string`
     - `formatConciseActiveResultLine(colWidth int, repoName, statusLabel string) string`
     - `renderConciseActiveResultsTo(w io.Writer, states []*PullRepoState)`
4. **Unit & Isolated Temporary E2E Tests**:
   - Add unit tests in `cli/cmdpull/pull_efficient_test.go`.
   - Add isolated temporary E2E test in `cli/tests/e2e/pae_column_width_tempe2e_test.go` (`//go:build tempe2e`, `RUN_TEMP_E2E=1`).
5. **Version Bump & Release**:
   - Bump version to `v6.330.0` (or next patch version), commit atomically, tag, and push.
