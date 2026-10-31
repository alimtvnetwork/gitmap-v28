# Specification 367.04: GitMap Skill Anti-Pattern Ban Matrix & Restrict Condition Policy

- **Sequence ID:** `367`
- **Subtask ID:** `TASK-04`
- **Feature Slug:** `sync-script-versioning-and-gitmap-skill-enhancement - Task 04`

---

## 1. Objective

Update `.agents/skills/gitmap/skill.md` across repositories with an authoritative, permanent Ban and Replacement Reference Matrix targeting the exact ad-hoc shell, Python, and Node anti-patterns observed during agent workflows.

---

## 2. Anti-Pattern Replacement Matrix

| Banned Anti-Pattern | Reason for Ban | GitMap Native Replacement |
| :--- | :--- | :--- |
| `powershell -Command "Get-Content '...' -Tail N \| Select-String -Pattern '...'"` | Slow object pipeline, terminal buffer bloat, Windows-only syntax | `gitmap cat "<file>" -t N -q "<pattern>"` |
| `gitmap py -c "import time; time.sleep(N)"` or shell sleep | Blocks processes, wastes CPU, violates no-sleep guidelines | Use Antigravity `schedule`, `gitmap watch`, or reactive task wakeup |
| `Get-ChildItem -Path "..." -File \| Select-Object ...` or `dir /s /b` | Slow COM/PowerShell startup, Windows-only, high latency | `gitmap lf "<dir>"` or `gitmap find "*.*" "<dir>"` |
| `node -e "const fs = require('fs'); ... console.log(...)"` | Spawns heavy Node V8 VM to read disk bytes; breaks on missing Node | `gitmap cat -L N "<file>"` or `gitmap cat -L N:M --num "<file>"` |
| `gitmap py -c "import os; [print(...) for ... in os.walk(...) if ... ]"` | Uncached filesystem walk, high memory usage, redundant script overhead | `gitmap find "*<pattern>*" [dir]` or `gitmap aum search "<pattern>" [dir]` |
| `gitmap py -c "print(open('...').read())"` | Unbuffered Python process spawn, risks OOM on large files | `gitmap cat "<file>"` or `gitmap cat -L 1:100 "<file>"` |
| Bypassing V7 restrict conditions during fleet sync | Violates weekday morning window policy (< 08:00 AM) | Use `38-sync-prompts-skills-scripts-v7.py` with `--auto-window` or `--schedule` |

---

## 3. Reference Examples for Agent Learning
Provide concrete executable command snippets directly inside `gitmap/skill.md` so that agents invoking `/learn` from `gitmap` ingest the correct patterns immediately.
