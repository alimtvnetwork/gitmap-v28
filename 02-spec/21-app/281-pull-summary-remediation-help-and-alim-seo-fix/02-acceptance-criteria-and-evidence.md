# Specification 281: Acceptance Criteria & Verifiable Evidence

> **Specification Slug:** `281-pull-summary-remediation-help-and-alim-seo-fix`  
> **Status:** `Active / Authoritative`  
> **Date:** `2026-10-17`  

---

## Acceptance Criteria

### AC-01: Pull Summary Up-to-Date Repositories Categorization
- **Criterion:** Post-pull concise output (`renderConciseActiveResults`) MUST categorize and display `Up-to-Date Repositories (%d)` alongside `Updated Repositories (%d)`, `Dirty Repositories (%d)`, and `Failed Repositories (%d)`.
- **Formatting:** Use comma-wrapped indentation to prevent terminal flooding when 50+ repositories are up-to-date.
- **Evidence:** Verified via unit test and local dry-run output formatting.

### AC-02: Real-Time Worker Progress Name Visibility
- **Criterion:** In multi-worker pull progress bar, `formatActiveWorkersList()` MUST display up to 4 active repo badges when terminal width >= 100 characters.
- **Evidence:** Output contains active repo names during parallel execution.

### AC-03: Remediation Menu Help Option `[h/help]`
- **Criterion:** Interactive prompt `[?] Remediate failed repositories?` in `pull_remediation.go` MUST accept `h`, `help`, or `?`.
- **Behavior:** Entering `h` prints `renderRemediationHelpGuide()` explaining auto-merge, step-through, and rebase strategies with safety guarantees, and re-prompts the user without exiting or failing.
- **Evidence:** Interactive test verifies help guide renders cleanly and returns to prompt.

### AC-04: Safe-Pull Conflict Detection Hardening
- **Criterion:** `isMergeConflict` in `cli/cloner/safe_pull.go` MUST check both error string and command output for merge conflict keywords (`"unmerged files"`, `"unresolved conflict"`, `"automatic merge failed"`, `"conflict ("`).
- **Behavior:** Stops repeated 4-attempt retries when repository is already in an unmerged conflict state.
- **Evidence:** Unit tests in `safe_pull_test.go` or `safe_pull_conflicts_test.go`.

### AC-05: `Generate Drift Check` CI Fix
- **Criterion:** Running `go generate ./...` in `cli` produces zero git diff drift.
- **Evidence:** `git diff --exit-code cli/` exits 0.

### AC-06: `alim-seo-writing` LFS Pointer Normalization
- **Criterion:** Resolve the persistent dirty status in `alim-seo-writing` caused by raw binary commits matching LFS filters.
- **Evidence:** `git -C D:/work/seo-writing/alim-seo-writing status --short` exits clean with 0 modified files.
