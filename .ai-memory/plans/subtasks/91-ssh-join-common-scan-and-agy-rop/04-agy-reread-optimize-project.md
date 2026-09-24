# Subtask 91-04: Antigravity (AGY) Project Re-Read & Optimization (`gitmap agy rop`)

Spec Reference: [02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md](../../../../02-spec/21-app/141-ssh-join-common-scan-and-agy-rop.md)
Parent Plan: [.ai-memory/plans/completed/91-ssh-join-common-scan-and-agy-rop.md](../../completed/91-ssh-join-common-scan-and-agy-rop.md)

## Objective
Implement `gitmap agy reread-optimize-project` (short alias: `rop [N]`, default `N = 5`) to optimize recently active Antigravity projects, clear cache with retention, safely back up conversations to a Split-DB (`data/agy/<repo-slug>/agy.db`), and allow storage reset.

## Functional Requirements
1. **Activity Filtering:**
   - Discover projects active/communicated in the last 24 hours (`time.Since(updatedAt) < 24*time.Hour`).
   - If fewer than N projects were active, optimize only those without error.
2. **Cache Clean & Retention:**
   - Execute cache clear keeping `--keep 10` conversations.
3. **Split-DB Archival Backup:**
   - Before removing or resetting conversations, serialize the conversation transcripts into a dedicated SQLite Split-DB at `data/agy/<repo-slug>/agy.db`.
   - Store table `AGYConversation` with fields: `ConversationID`, `ProjectSlug`, `ProjectPath`, `StepCount`, `TranscriptJSON`, `CreatedAt`.
4. **Storage Reset Integration:**
   - Wire `gitmap storage reset` and `gitmap storage reset-errors` to optionally purge or reset `data/agy/` backup databases when requested with `--all` or specific target flags.

## Files to Create/Modify
- `cli/cmdagy/agy_rop_types.go` [NEW]
- `cli/cmdagy/agy_rop_backup.go` [NEW]
- `cli/cmdagy/agy_rop.go` [NEW]
- `cli/cmdagy/agy_rop_cmd.go` [NEW]
- `cli/cmdagy/agy_rop_test.go` [NEW]
- `cli/cmdagy/agy_cmd.go` [MODIFY]
- `cli/cmd/storage_cmd.go` & `cli/cmd/storage_ls.go` [MODIFY]
