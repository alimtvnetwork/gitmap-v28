# Antigravity (AGY) Integration

Manage Google Antigravity workspaces, prompts, deduplication, and conversation health.

## Command Files in this Folder

| File | Subcommand | Description |
|---|---|---|
| [`ls.md`](./ls.md) | `gitmap agy ls` | Status table of all projects & empty conversation detection |
| [`remove-projects-with-empty-conversations.md`](./remove-projects-with-empty-conversations.md) | `gitmap agy remove-projects-with-empty-conversations` | Prune projects having 0 or aborted sessions (`--except`, `--dry-run`) |
| [`optimize-projects.md`](./optimize-projects.md) | `gitmap agy optimize-projects` | Deduplicate repeated project paths (`--repeat-fix`) |
| [`clear.md`](./clear.md) | `gitmap agy clear` | Remove missing or stale workspace records (`--except`) |
| [`prompt.md`](./prompt.md) | `gitmap agy prompt` | Send prompts to single or all project sessions |
| [`sync-and-backup.md`](./sync-and-backup.md) | `gitmap agy sync`, `export`, `import` | Sync repos to workspaces and backup configurations |

## Overview
Antigravity workspaces are stored as JSON files under `~/.gemini/config/projects/*.json` and conversation histories are recorded in SQLite databases under `~/.gemini/antigravity/conversations/*.db`.
