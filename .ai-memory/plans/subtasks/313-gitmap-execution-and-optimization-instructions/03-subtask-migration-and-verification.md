# Subtask 03: Migration, 7-Day Purge Verification & Quality Gates (Spec 313)

- **Parent Plan:** [.ai-memory/plans/pending/313-gitmap-execution-and-optimization-instructions.md](../../pending/313-gitmap-execution-and-optimization-instructions.md)
- **Specification:** [02-spec/21-app/313-gitmap-execution-and-optimization-instructions/03-migration-and-verification-spec.md](../../../02-spec/21-app/313-gitmap-execution-and-optimization-instructions/03-migration-and-verification-spec.md)
- **Status:** `ACTIVE`

## Scope & Deliverables
1. Step-by-Step Rollout Matrix:
   - Phase 1: Windows remote SSH quoting fix (`-EncodedCommand`) and separate STDOUT/STDERR isolation (`RunCommandSeparated`).
   - Phase 2: Split-DB pull task gate, 10-second observer ticker loop, rolling ETA calculator, and 7-day retention purge.
   - Phase 3: SQLite SSH Key Vault encryption and auto-relogin hook.
   - Phase 4: Direct JSON payload ingestion for heatmap and agent commands, with ModelName/ToolName logging.
   - Phase 5: Heatmap skill authoring (`.agents/skills/heatmap/skill.md`) and V6 execute prompt updates.
   - Phase 6: Dynamic pull worker concurrency scaling and in-memory update streaming extraction.
   - Phase 7: Repository scan disable/enable CLI, REST API endpoints, and React UI controls.
2. 7-Day Retention Purge Validation:
   - Purge `PullItemHistory` and `PullTask` older than 7 days based on `SettingPullRetentionDays`.
   - Guarantee zero file/diff contents stored in the database.
3. Verification Checklist:
   - Targeted linters on modified files.
   - 100% relative paths.
   - Zero secrets committed.
