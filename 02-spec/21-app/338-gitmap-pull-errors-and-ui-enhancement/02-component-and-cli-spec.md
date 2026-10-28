# Component & CLI Specification: GitMap Pull Errors and UI Enhancement
# Spec Sequence: 338
# Slug: gitmap-pull-errors-and-ui-enhancement

## 1. Casing Normalization Persistence (`cli/cmdfix`)

In `cli/cmdfix/fix_casing.go`:
- Extend `ReconcileDiscrepancyList` to automatically commit the normalized casing when executed:
```go
func CommitCasingNormalization(repoPath string, d CasingDiscrepancy) error {
	msg := fmt.Sprintf("chore: normalize casing %s to %s", d.CurrentPath, d.CanonicalPath)
	return exec.Command("git", "-C", repoPath, "commit", "-m", msg).Run()
}
```
- Only commit if `git status --porcelain` shows staged changes.
- This guarantees casing normalizations are recorded permanently in Git history and do not re-surface during future pull or fix runs.

## 2. Pull Error Eviction Engine (`cli/cmdpull`)

In `cli/cmdpull/pull_missing_remediate.go`:
- Update `purgeInvalidRepoEntry(repoName, repoPath string)`:
```go
func purgeInvalidRepoEntry(repoName, repoPath string) {
	if repoPath == "" {
		return
	}
	_ = DeleteRepoRecordByPath(repoPath)
	clearPullSplitDBError(repoName, repoPath)
}

func clearPullSplitDBError(repoName, repoPath string) {
	db, err := store.OpenPullSplitDB()
	if err != nil {
		return
	}
	defer db.Close()
	_ = db.ClearPullErrorsForRepo(repoName, repoPath)
}
```
- In `RemediateMissingRepo`:
```go
if remoteURL == "" || strings.Contains(remoteURL, " ") {
    printMissingRemoteURLError()
    purgeInvalidRepoEntry(f.RepoName, f.RepoPath)
    return false
}
```
- This permanently evicts ghost records from both `gitmap.db` and the Pull Split-DB error queue.

## 3. High-Definition Icon & Favicon Suite

In `cli/cmdlogviewer/favicon.go`:
- Replace `GitMapLogoSVG` with modern, crisp, 100% transparent vector logo:
  - Clean branch fork geometry floating freely with 8% quiet zone.
  - Distinct colors: `#89b4fa` (primary blue branch), `#a6e3a1` (merged green node), `#cdd6f4` (white accent).
  - Multi-resolution crisp SVG rendering with zero background card artifacting.
- Synchronize across `public/favicon.ico`, `cli/cmdlogviewer/favicon.ico`, and UI asset headers.

## 4. Mail Web UI & Mailbox Inspector

In `cli/cmdmail/mail_ui_html.go` & `mail_ui_server.go`:
- Add Mailbox Overview Panel:
  - Displays configured senders, host reachability, and unread counts.
  - Interactive sender selection pills with active indicators.
- Composer Enhancements:
  - Live markdown preview with syntax highlighting and crisp modern buttons.
  - Crisp SVG iconography for Send, Clear, Mailbox, and Template actions.
  - Masked / Unmasked address toggle.
