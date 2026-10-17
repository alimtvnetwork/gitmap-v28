# Plan: 280-large-files-split-flatten-and-muse-sync

## Status: COMPLETED
- **Task ID:** 280-large-files-split-flatten-and-muse-sync
- **Milestone:** Milestone 36 (Large File Splitting, Git LFS Conversion, Repository Flattening & Muse Sync)
- **Specification:** `02-spec/21-app/280-large-files-split-flatten-and-muse-sync/`

## User Request (Verbatim)
```text
PS D:\work\story-telling> gitmap cfr https://github.com/alimtvnetwork/story-telling-v1
Cloning story-telling-v1 into story-telling...
...
fatal: cannot create directory at 'story/audio/chhotto-lal-murgi': No space left on device
warning: Clone succeeded, but checkout failed.
You can inspect what was checked out with 'git status'
and retry with 'git restore --source=HEAD :/'
```

## Subtask Breakdown

- [x] **Subtask 01:** CI/CD error fix & Muse training source prioritization in `cli/cmd/llm/llm_train_muse.go` (`decisionLog` restore, `repo-secrets` primary path).
- [x] **Subtask 02:** Large files splitting engine in `cli/cmdsplit/split_cmd.go` & `cli/cmdsplit/multipart_zip.go` ($> 10\text{ MB}$, multipart zip, STORE for media, DEFLATE for data, history purge, dry-run).
- [x] **Subtask 03:** Git LFS conversion command in `cli/cmdlfs/lfs_convert.go` ($> 5\text{ MB}$, `.gitattributes` tracking, LFS migration, dry-run).
- [x] **Subtask 04:** Repository flattening engine & undo manager in `cli/cmdflatten/flatten_cmd.go` (numbered successor, single flat commit, marked-for-deletion archive, swap, undo).
- [x] **Subtask 05:** CLI command registration and hyphen-to-space normalization in `cli/cmd/rootutility.go` and `cli/cmd/rootdata275.go`.
- [x] **Subtask 06:** Help documentation (`cli/helpdoc/split-large-files.md`, `lfs-convert.md`, `flatten.md`), `cli/cmdhelp/howto.go`, and UI data in `src/data/commands.ts` and `src/data/troubleshootingIssues.ts`.
- [x] **Subtask 07:** Quality verification (linters, Spec 280 evidence, plan completion, version bump to `v6.562.0`, and atomic push via `gitmap cpf`).
