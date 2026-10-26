# Subtask 01: Commit Author Information
**Sequence ID:** 335  
**Parent Plan:** [335-gitmap-author-information-and-search-feature.md](file:///d:/work/gitmap/.ai-memory/plans/pending/335-gitmap-author-information-and-search-feature.md)  
**Status:** Pending  

---

## Actions
1. Extend `CommitMetadata` in `cli/cmdcommit/commit_types.go` with `Author string`.
2. Add `--author` / `-a` flag parsing in `cli/cmdcommit/commit_flags.go` and implement `resolveAuthorFallback`.
3. In `cli/cmdcommit/commit_push.go`, pass `--author` to `git commit` when present.
4. In `cli/cmdcommit/commit_ui.go`, render `Author: <author>` row in `renderCommitPushSummaryCard`.
5. Update tests in `cli/cmdcommit/commit_flags_test.go` and `cli/cmdcommit/commit_push_test.go`.
