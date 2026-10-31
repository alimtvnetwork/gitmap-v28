---
name: gitmap-heatmap
description: Autonomously inspect files, track in-flight touches without UDF flags, and query cognitive decision trees via GitMap Heatmap.
---

# GitMap Heatmap & File Inspection Skill (`gitmap-heatmap`)

> **Invocation:** `@[gitmap-heatmap](file;.cursor/skills/gitmap-heatmap)`

## Mission & Purpose
This skill provides authoritative architectural guidance, CLI commands, and execution rules for inspecting files, tracking cognitive traces, logging agent actions, and preventing file write collisions using GitMap's SQLite Heatmap tier (`repodb/heatmap.db` and `<EffectiveDataDir>/heatmap.db`).

---

## 1. Key Architectural Invariants & Automatic Touch

1. **Automatic Heatmap Touch on File Read (Zero UDFs):**
   - Whenever an AI agent inspects, reads, or slices a file, it MUST use `gitmap heatmap <filepath>` or `gitmap cat <filepath>`.
   - **No User-Defined Fields/Flags (UDFs) Required:** Do NOT pass `--reason`, `--action`, or manual JSON payloads for simple reads. GitMap automatically resolves relative/absolute paths, sets `action="read"`, `tool="cat"`, and records line ranges in `repodb/heatmap.db`.
   - The file content is printed directly in the terminal, and the heatmap record is stored atomically in SQLite.

2. **Total Ban on Ad-Hoc Python File Reading:**
   - ❌ **BANNED:** `gitmap py -c "lines = open('path/to/file.ext').readlines(); print(''.join(lines[1130:1155]))"`
   - ❌ **BANNED:** `python -c "open('...').read()"`
   - 👉 **MANDATORY REPLACEMENT:** `gitmap heatmap path/to/file.ext -L 1130:1155 --num` (or `gitmap cat path/to/file.ext -L 1130:1155 --num`)

---

## 2. Essential Commands

```bash
# Read file and automatically touch heatmap (replaces ad-hoc Python reading)
gitmap heatmap <filepath>
gitmap heatmap src-tauri/src/modules/repo_db/agy.rs -L 1130:1155 --num

# View file length and line metadata
gitmap heatmap <filepath> --info

# Native cat command automatically logs heatmap touch in background
gitmap cat <filepath> [-L start:end] [-n head] [-t tail] [-q query]

# View active in-flight file touches and leases across all swarm agents
gitmap heatmap live
gitmap heatmap status

# View hierarchical execution tree of agent reasoning
gitmap heatmap tree
gitmap heatmap tree --task <task-slug>

# Search past agent decisions and touched files
gitmap heatmap search "<query>"
gitmap heatmap search --file <path>

# Cross-repository heatmap inspection (RC Heatmap)
gitmap rc heatmap status
gitmap rc heatmap tree
gitmap rc heatmap live
gitmap rc heatmap sync
```
