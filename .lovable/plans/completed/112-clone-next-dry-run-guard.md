# 112-clone-next-dry-run-guard.md: Clone-Next Dry-Run Side-Effect Prevention & Preview Verification

**Status: completed**

## 1. Executive Summary

Resolves open Issue 05 from .lovable/pending-issues/01-current-issues.md:
In gitmap/cmd/clonenext.go, executeCloneNextPipeline previously invoked prepareCloneNextTarget (which physically deletes existing target folders on disk) and handleCreateRemote (which creates remote repositories on GitHub) before inspecting cnFlags.DryRun.

This plan moves dry-run detection and preview execution upfront via a dedicated handleCloneNextDryRun helper in gitmap/cmd/clonenextdryrun.go, ensuring that gitmap clone-next <version> --dry-run performs pure planning with zero destructive side-effects.

## 2. Mandatory Rules & Invariants

1. Rule 1 (Strict Relative Git Paths): All paths in markdown files, artifacts, and documentation must be strictly relative to the repository root. Zero file:/// URIs.
2. Rule 2 (Strict File & Function Sizing): Every function <= 15 lines, mandatory blank line before every return statement.
3. Rule 3 (Affirmative Booleans): All booleans start with is or has. No explicit == true.
4. Rule 4 (Zero Swallowed Errors): Wrap all errors with appropriate context.
5. Rule 5 (Cache Tracking): Record touched files in .lovable/temp/recent-file-changes.json.
6. Rule 6 (Consolidated Commit & Push): Commit all touched files together and immediately push to origin main.

## 3. Subtasks Breakdown

1. Add handleCloneNextDryRun(isDryRun bool, url, dest string) bool in gitmap/cmd/clonenextdryrun.go.
2. Refactor executeCloneNextPipeline and performCloneNextOperation in gitmap/cmd/clonenext.go so dry-run is handled upfront before folder removal or remote creation.
3. Add unit test in gitmap/cmd/clonenextdryrun_test.go.
4. Verify with targeted Go test and linters (check-newline-styling.py, check-boolean-guidelines.py).
5. Update .lovable/pending-issues/01-current-issues.md to mark Issue 05 as FIXED.
