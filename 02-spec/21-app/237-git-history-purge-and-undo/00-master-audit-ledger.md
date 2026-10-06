# Ledger: 237-git-history-purge-and-undo

Request slug: 237-git-history-purge-and-undo
Request first line: Based on what you have done with this repository, how you can actually remove the Git history, create several commands in `gitmap` to do that.
Status: COMPLETED
Phase: 3    Wave: 2 / 2    Step: 18 / 300
Last completed action: Phase 3 Consolidation & Register Updates
Next action: Atomic GitMap Commit
Workers in flight: none
Commits: none    Pushed: no
Branch: main | Tree at start: clean
Tools: invoke_subagent=yes send_message=yes ask_question=yes gitmap=yes

| Task-ID | Subtask | Owner | Owned files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | 01-architecture-spec | Subagent 01 | `02-spec/21-app/237-git-history-purge-and-undo/01-architecture-spec.md` | DONE | PASS (Exit 0) |
| Task-02 | 02-component-and-cli-spec | Subagent 02 | `02-spec/21-app/237-git-history-purge-and-undo/02-component-and-cli-spec.md` | DONE | PASS (Exit 0) |
| Task-03 | 01-splitdb-schema-and-temp-backup | Subagent 01 | `.ai-memory/plans/subtasks/237-git-history-purge-and-undo/01-splitdb-schema-and-temp-backup.md` | DONE | PASS (Exit 0) |
| Task-04 | 02-preflight-analyzer-and-graph-diff | Subagent 01 | `.ai-memory/plans/subtasks/237-git-history-purge-and-undo/02-preflight-analyzer-and-graph-diff.md` | DONE | PASS (Exit 0) |
| Task-05 | 03-core-purge-engine-and-release-cleanup | Subagent 02 | `.ai-memory/plans/subtasks/237-git-history-purge-and-undo/03-core-purge-engine-and-release-cleanup.md` | DONE | PASS (Exit 0) |
| Task-06 | 04-undo-engine-and-cli-wiring | Subagent 02 | `.ai-memory/plans/subtasks/237-git-history-purge-and-undo/04-undo-engine-and-cli-wiring.md` | DONE | PASS (Exit 0) |

Assumptions:
- SQLite SplitDB Tier 1/3 stores the `HistoryPurgeOperation`, `HistoryPurgeCommitMap`, `HistoryPurgeFile`, and `HistoryUndoOperation` records.
- Temp backup directory isolates all extracted original blobs under `os.TempDir()/gitmap/history-backup/<repo-slug>/<operation-id>/`.
- Pre-flight visual report renders before/after commit graph using terminal box glyphs and commit count statistics.
- Undo command prints non-warranty advisory disclaimer: "You can undo this if you wanted to. We do not confirm this, but you can try: gitmap history undo <id>".

Conflicts: none

Stage list:
- `02-spec/21-app/237-git-history-purge-and-undo/00-master-audit-ledger.md`
- `02-spec/21-app/237-git-history-purge-and-undo/01-architecture-spec.md`
- `02-spec/21-app/237-git-history-purge-and-undo/02-component-and-cli-spec.md`
- `.ai-memory/plans/237-git-history-purge-and-undo.md`
- `.ai-memory/plans/subtasks/237-git-history-purge-and-undo/01-splitdb-schema-and-temp-backup.md`
- `.ai-memory/plans/subtasks/237-git-history-purge-and-undo/02-preflight-analyzer-and-graph-diff.md`
- `.ai-memory/plans/subtasks/237-git-history-purge-and-undo/03-core-purge-engine-and-release-cleanup.md`
- `.ai-memory/plans/subtasks/237-git-history-purge-and-undo/04-undo-engine-and-cli-wiring.md`
- `cli/store/purge_history.go`
- `cli/store/purge_history_models.go`
- `cli/tempdir/history_backup.go`
- `cli/cmdpurge/preflight.go`
- `cli/cmdpurge/graph_preview.go`
- `cli/cmdpurge/purge_engine.go`
- `cli/cmdpurge/release_prune.go`
- `cli/cmdpurge/purge_undo.go`
- `cli/cmd/history_purge_cmd.go`
- `cli/cmd/roottooling.go`
- `cli/cmd/root_cobra_completion.go`
