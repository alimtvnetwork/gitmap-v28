# Subtask 02: Commit `--ai-agent` Flag & Heatmap Provenance Attribution

- **Parent Plan:** `359-gitmap-feature-implementation-and-update`
- **Slug:** `gitmap-feature-implementation-and-update - Task 02`
- **Status:** `pending`
- **Assigned Subagent:** Worker 02

---

## Deliverables
1. Define `FlagAIAgent = "--ai-agent"` in `cli/cmdcommit/constants/constants.go` and `cli/cmdcommit/commit_types.go`.
2. Update `cli/cmdcommit/commit_flags.go`:
   - Handle `FlagAIAgent` in `hasAgentOrHeatmapPrefix` and `isAgentFlag`.
   - Set `meta.AgentID = name`, `meta.AIName = name`, `meta.IsAIMode = true`.
   - In `resolveAuthorFallback`: if `--author` was not explicitly passed and `meta.AgentID != ""` or `meta.AIName != ""`, default author to `fmt.Sprintf("gitmap automation bot - %s - runs autonomously <info@rgitmap.com>", meta.AIName)`.
3. In `cli/cmdcommit/commit_history_heatmap.go`:
   - Ensure `buildHeatmapActionRecord` populates `ModelName: rec.AiModelName` and `Notes: rec.AgentId`.
   - In `insertHeatmapCommitLogRow`, record agent attribution in `Notes`.
4. Add unit test in `cli/cmdcommit/commit_flags_test.go` or `commit_push_test.go` asserting `--ai-agent` flag parsing.
5. Coding Guidelines Compliance:
   - Functions <= 15 lines.
   - Zero nested ifs.
   - Positive boolean naming (`isAIMode`, `hasPush`).
   - Leaf error sovereignty.
