# 167: Visual Layout, CLI Help Text & UX

## 1. CLI Commands & Aliases

Prefix aliases: `gitmap agy <cmd>` or top-level `gitmap <cmd>`:

### 1.1 Running Prompts Backup & Restore
- `gitmap agy backup-running-prompts [-file/-f <path>]`
- `gitmap agy backup-running-prompts ls [--json]`
- `gitmap agy running-prompts backup [-file/-f <path>] [--json]`
- `gitmap agy running-prompts backup ls [--json]`
- `gitmap agy running-prompts restore [--keep/-k] [--json]`
- `gitmap agy restore-running-prompts [--keep/-k] [--json]`
- `gitmap agy running-prompts clean [--force]`

### 1.2 Running Prompts Inspection & Export/Import
- `gitmap agy running-prompts ls [--limit/-l 8] [--wordcount/--wc 100] [--full] [--json]`
- `gitmap agy running-prompts export [-file <path>] [--wc 200]`
- `gitmap agy running-prompts import [-file <path>] [--wc 200]`

### 1.3 Running Projects
- `gitmap agy running-projects [ls] [--json] [-file/-f <path>] [--ssh]`
- `gitmap agy running-projects help`

### 1.4 Automation Watchers
- `gitmap agy finish-prompts-until-green` (alias `fpug`) `<target...> [-t <duration>]`
- `gitmap agy finish-prompts-until-green running-projects [-t <duration>]`
- `gitmap agy shutdown-until-green` (alias `sug`) `ls/help/run/add-projects/rm/agy-running-projects [-t <duration>]`

## 2. Help Text Specifications

Each command must output formatted help text showing synopsis, description, flags, and actionable examples with color-coded syntax.

### 2.1 `running-prompts backup ls` Output Example
```text
  ┌── Antigravity Running Prompts Backup Registry ────────────────────────┐
  │ Database: data/backup-prompts/sql.db (64.0 KB)                        │
  │ Total Batches: 3 · Total Backed Prompts: 14                          │
  └───────────────────────────────────────────────────────────────────────┘
  Batch ID      Created At           Total  Running  Enqueued  Status
  b-9a1b2c3d    2026-09-27 07:15:00  5      2        3         active
  b-4e5f6a7b    2026-09-26 18:30:00  9      1        8         restored
```

### 2.2 Auto-Prune Single Line Notice
When old restored entries (> 1 day old without `--keep/-k`) are pruned during execution:
```text
old data has been removed
```
