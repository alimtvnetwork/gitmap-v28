# Heatmap Auto-Read & Auto-Touch Architecture Specification

> **Task:** 369-heatmap-command-enhancements-and-high-speed-sync  
> **Subsystem:** Heatmap Subsystem & File Inspection (`cli/heatmap/`, `cli/cmdmacro/`)  
> **Date:** 31-Oct-26  
> **Status:** APPROVED  

---

## 1. Executive Summary & Objective

In modern autonomous coding workflows, AI agents repeatedly read, inspect, and analyze source files. However, agents have frequently resorted to raw Python scripts (`open(...).readlines()`), bypassing local cognitive recording in GitMap's SQLite Heatmap (`repodb/heatmap.db`). When using `gitmap heatmap touch`, agents were previously forced to specify complex user-defined flags (`--reason`, `--action`).

This specification mandates that:
1. When passed an existing file path (relative or absolute), `gitmap heatmap <filepath> [flags]` automatically:
   - Normalizes the path to a clean repo-relative path.
   - Reads and displays the file content via the high-speed `cat` reader (supporting `-L start:end`, `-n head`, `-t tail`, `-q filter`, `--num`, `--raw`).
   - Automatically records a touch/action log in `repodb/heatmap.db` with `action="read"`, `tool="cat"`, `reasoningState="EXPLORATION"` without requiring any User-Defined Fields/Flags (UDFs).
2. When an agent runs `gitmap cat <filepath> [flags]`, GitMap automatically records the inspection trace into `repodb/heatmap.db` in the background (non-blocking, zero-crash error suppression).

---

## 2. CLI Command Contracts

```bash
# Read file and automatically touch heatmap (replaces ad-hoc Python reading)
gitmap heatmap src-tauri/src/modules/repo_db/agy.rs -L 1130:1155 --num

# Read whole file and log touch
gitmap heatmap cli/cmd/root.go

# Standard cat command automatically logs heatmap touch
gitmap cat cli/cmdui/ui_server_settings.go
```

### Argument Resolution & Auto-Detection
- If `pos[0]` is a valid filesystem path (file exists):
  - Automatically route to `cat` rendering + background heatmap action logging.
  - No need to supply `--reason` or UDF flags.
- If `pos[0]` is a recognized heatmap administrative subcommand (`status`, `tree`, `live`, `search`, `log`, `overview`, `tables`, `query`, `rc`), dispatch normally to the admin handler.

---

## 3. Database Schema Mapping

Every file inspection produces a record in `HeatmapActionLog`:
- `TaskSlug`: Active task slug (ambient env or active task in memory).
- `RepoName`: Active repo name.
- `ActionType`: `"read"` / `"analyze"`.
- `FilePath`: Normalized relative path (e.g. `src-tauri/src/modules/repo_db/agy.rs`).
- `StartLine`: Start line from `-L` flag (or 1).
- `EndLine`: End line from `-L` flag (or file line count).
- `ToolName`: `"cat"`.
- `ReasoningState`: `"EXPLORATION"`.
- `ReasoningNotes`: Automatically populated as `"inspected file via gitmap heatmap"`.
