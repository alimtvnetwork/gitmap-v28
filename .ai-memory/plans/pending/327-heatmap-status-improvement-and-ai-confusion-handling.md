# Plan: Heatmap Status Improvement and AI Confusion Handling

- **Task Sequence:** 327
- **Task Slug:** `heatmap-status-improvement-and-ai-confusion-handling`
- **Spec Path:** `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/`
- **Status:** `COMPLETED`
- **Created At:** 2026-10-25

## Executive Summary

This plan completes the architectural refinements for GitMap:
1. Root SQLite database caching of repo-cache sequence folders (`RepoCacheSequenceIndex`), mapping `RepoSlug` to `XX-<repo>` folders in `repo-cache`.
2. Normalized 1-to-N database schema for AI confusion file references (`AiConfusion` and `AiConfusionFile`).
3. Autonomous synchronization workflow in `gitmap rc ai-confusion status` (export pending deltas, revert db, git pull --rebase repo-cache, replay via upsert, commit, display status).
4. Relative file context querying (`gitmap rc ai-confusion view <file.go>`) that discovers the owning repo and displays full agent reasoning, spec links, and line snippets.
5. Ghost repository sanitization: rejecting commit messages (`COMMIT - ...`) from failed repository remediation lists.
6. Execution of the minor version release ceremony (`v6.611.0`).

## Subtasks

- [x] `01-subtask-spec-and-planning.md` — Author comprehensive 327 specifications.
- [x] `02-subtask-sequence-folder-caching-and-1-to-n-schema.md` — Implement root SQLite sequence folder cache (`RepoCacheSequenceIndex`) and 1-to-N `AiConfusionFile` schema & store.
- [x] `03-subtask-rc-confusion-status-sync-and-file-view.md` — Implement `rc ai-confusion status` automatic sync and relative file context resolution.
- [x] `04-subtask-pull-output-collapsing-and-repo-sanitization.md` — Sanitize repository name candidates in remediation queries and confirm pull output collapsing.
- [x] `05-subtask-verification-and-release-ceremony.md` — Verification, documentation updates, and minor version release ceremony (`v6.611.0`).
