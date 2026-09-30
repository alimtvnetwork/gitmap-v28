# Agent Memory Lifecycle

## 1. Purpose of `.ai-memory/memories/`

The `.ai-memory/memories/` folder acts as the institutional knowledge hub for AI sessions. However, to prevent this directory from becoming a bloated, contradictory mess, a strict lifecycle must be enforced.

## 2. Time-To-Live (TTL) and Cleanup

- **Pending/Planned Tasks:** Files in `pending/` and `planned/` must be moved to `done/` upon task completion.
- **Deduplication:** AI agents must periodically scan `.ai-memory/memories/` for duplicated patterns or rules that have since been formalized into the `02-spec/` folder.
- **Stale Memory Purge:** Any memory file that contradicts the canonical `02-spec/` folder is considered stale and MUST be deleted. The `02-spec/` folder ALWAYS wins.

## 3. Formatting

Every memory file must include an `Updated:` date stamp. This allows agents to determine which memory file is more recent if a conflict arises within the memory folder itself.
