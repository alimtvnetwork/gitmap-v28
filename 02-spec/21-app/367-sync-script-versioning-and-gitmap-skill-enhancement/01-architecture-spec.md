# Specification 367: Sync Script Versioning, GitMap Skill Anti-Pattern Ban Matrix & Cat Query Window Enhancements

- **Sequence ID:** `367`
- **Feature Slug:** `sync-script-versioning-and-gitmap-skill-enhancement`
- **Status:** `in_progress`
- **Creation Date:** 31-Oct-26
- **Architecture Lead:** V7 Master Autonomous Execution Engine

---

## 1. Executive Summary

Autonomous AI coding agents frequently revert to fragile, platform-dependent shell one-liners, unbuffered Python scripts, or Node.js invocations to perform routine tasks such as reading lines from logs, searching directories, or polling with sleep loops. 

Specifically, telemetry reveals recurring anti-patterns:
1. `powershell -Command "Get-Content '...' -Tail 5 | Select-String -Pattern 'media|USER_INPUT'"`
2. `gitmap py -c "import time; time.sleep(8)"`
3. `Get-ChildItem -Path "cli/cmdinstall" -File | Select-Object -ExpandProperty Name`
4. `node -e "const fs = require('fs'); const lines = fs.readFileSync('.../transcript.jsonl', 'utf-8').trim().split('\n'); console.log(JSON.parse(lines[1519]).content);"`
5. `gitmap py -c "import os; [print(root, files) for root, dirs, files in os.walk('...') if any('subagent' in f.lower() or 'research' in f.lower() for f in files)]"`
6. `gitmap py -c "print(open('cli/cmdui/ui_server_settings.go', encoding='utf-8').read())"`
7. `gitmap py -c "import os; [print(os.path.join(r, f)) for r, d, fs in os.walk('cli/cmdui') for f in fs if 'settings' in f.lower()]"`

This specification establishes:
1. **Total Ban on File Reading & Traversal Anti-Patterns:** Enshrines explicit bans and high-visibility side-by-side alternative reference matrices in `.agents/skills/gitmap/skill.md` across both `gitmap` and `coding-guidelines`.
2. **Dual Sync Script Versioning:** Provides two distinct operating modes for fleet synchronization:
   - **V7 Execution (`38-sync-prompts-skills-scripts-v7.py` or `--v7` / `--mode v7`):** Strictly enforces calendar-aware morning commit cutoffs (< 08:00 AM on weekdays; 24h unconstrained on weekends), target timezone resolution, auto-window slotting, and bot identity governance.
   - **V6 Execution (`38-sync-prompts-skills-scripts.py` or `--mode v6`):** Runs unrestricted for rapid development without calendar restrictions.
3. **`gitmap cat` Enhancements:**
   - Adds `-t` as standard shorthand for `--tail`.
   - Combines query pattern filtering (`-q` / `--query`) with window flags (`--tail`, `--head`, `-L start:end`) so tail-filtering and range-filtering execute seamlessly in a single zero-disk streaming pass.
4. **`gitmap py` Anti-Pattern Guard:** Extends inline evaluation warnings to flag `os.walk` and `time.sleep` anti-patterns and guide agents to native GitMap tools.
5. **Cross-Repository Synchronization:** Propagates all skill updates and sync scripts across connected repositories.
