# Fix commit-right missing commits bug

## Goal

Explain and fix why `gitmap commit-right` only considered 8 commits instead of 179 after a previous failure.

## Root Cause

When `commit-right` executes, it walks the source repository's history by running `git checkout <sha>` for each commit to take a snapshot, and relies on a `defer` block (`defer func() { _ = checkoutRef(plan.SourceDir, plan.SourceHEAD) }()`) to restore the original branch when it finishes.

However, if the source repository has uncommitted changes (e.g., `README.md` was modified), the `git checkout` command inside the loop fails on a conflict. Crucially, the deferred `checkoutRef` ALSO fails because of those same uncommitted changes. 
As a result, the source repository is permanently left in a "Detached HEAD" state at the last successful commit (commit 8).

When the user runs `commit-right` again, it looks at `HEAD` of the source repository, which is now pointing to commit 8, and only sees those 8 commits!

## Plan

1.  **Check for dirty source:** Inject an `isWorkingTreeDirty(sourceDir)` check at the beginning of `runOneDirection`.
2.  If the source tree has uncommitted changes, abort immediately with an actionable error.
3.  Implement `isWorkingTreeDirty` using `git status --porcelain`.
4.  Answer the user's question clearly.
