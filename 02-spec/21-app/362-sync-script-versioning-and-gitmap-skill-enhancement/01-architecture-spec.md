# Spec 362: Sync Script Versioning, GitMap Cat Range Inspection, and Ad-Hoc Python Reading Ban

- **Spec Number:** 362
- **Title:** Sync Script Versioning and GitMap Skill Enhancement
- **Slug:** `sync-script-versioning-and-gitmap-skill-enhancement`
- **Status:** Approved
- **Created Date:** 31-Oct-26
- **Architecture Standards:** GitMap V7 Engine, Split SQLite DBs, Leaf Error Sovereignty, Positive Booleans, <= 15 Line Functions

---

## 1. Executive Summary

AI coding agents across multi-repository workspaces frequently degrade to ad-hoc, brittle Python one-liners (`gitmap py -c "open_f = open(...).read(); [print(l) for l in open_f.splitlines() if ...]"` or `gitmap py -c "lines = open(...).readlines(); print(''.join(lines[140:220]))"`) or raw PowerShell `Get-Content` commands when inspecting files. This occurs because `gitmap cat` previously lacked native line slicing, line numbering, and query filtering flags.

Furthermore, multi-repository asset synchronization from `coding-guidelines` lacked clear dual versioning between **V6 (Unrestricted)** and **V7 (Time-Restricted)** operating engines. Under V7 execution, tasks must enforce calendar-aware morning commit cutoffs (< 08:00 AM on weekdays or auto-window scheduling) and autonomous bot identity governance (`gitmap login set ai-auto <agent>`).

This specification architects:
1. **Advanced `gitmap cat` Engine**: Native line range (`-L start:end`), line numbering (`--num`), head/tail limits (`--head N`, `--tail N`), and query filters (`-q <query>`).
2. **Dual-Versioned Sync Scripts**: `38-sync-prompts-skills-scripts.py` (V6 execution, unrestricted) and `38-sync-prompts-skills-scripts-v7.py` (V7 execution, time-restricted & autobot governed).
3. **Canonical Skill Hardening & Bans**: Explicit total bans on ad-hoc Python file reading and restrict condition bypasses, paired with a complete replacement matrix and on-the-fly teaching.
4. **Fleet Synchronization**: Propagating hardened skills and prompts across all 43 connected repositories.

---

## 2. Command Architecture & Replacements

### 2.1 Banned Ad-Hoc Snippets vs. Native GitMap Commands

| Banned Prohibited Snippet | GitMap Mandatory Native Command | Why GitMap is Required |
| :--- | :--- | :--- |
| ❌ `gitmap py -c "print(open('file.go').read())"` | ✅ `gitmap cat file.go` | Zero-disk streaming via `output.Raw()` with safe image protection cards. |
| ❌ `gitmap py -c "lines=open('file.go').readlines(); print(''.join(lines[140:220]))"` | ✅ `gitmap cat file.go -L 140:220` | Sub-millisecond line range slicing with Tokyo Night header banner. |
| ❌ `gitmap py -c "[print(l) for l in open('file.md').splitlines() if 'ban' in l.lower()][:15]"` | ✅ `gitmap cat file.md -q "ban" --limit 15` | High-speed case-insensitive line filtering with line numbers. |
| ❌ `pwsh -Command "Get-Content file.go -TotalCount 30"` | ✅ `gitmap cat file.go --head 30` | Cross-platform fast line extraction without PowerShell overhead. |
| ❌ `gitmap py -c "lines=open('file.go').readlines(); print(''.join(lines[-30:]))"` | ✅ `gitmap cat file.go --tail 30` | Direct tail extraction with line number preservation. |
| ❌ Ad-hoc Python line numbering loops | ✅ `gitmap cat file.go -L 140:180 --num` | Formatted gutter line numbering (`  140 | ...`). |

---

## 3. Sync Script Dual Versioning Architecture

### 3.1 V6 Unrestricted Engine (`38-sync-prompts-skills-scripts.py`)
- Executes fleet synchronization as-is across connected repositories.
- Disables time-window constraints; operates immediately.
- Preserves pre-change and post-change tags and multi-branch backups.

### 3.2 V7 Time-Restricted & Autobot Engine (`38-sync-prompts-skills-scripts-v7.py`)
- Invoked during V7 parent task execution.
- Evaluates operating mode:
  - Weekday (< 08:00 AM in target timezone or `--auto-window` slotting).
  - Weekend (24-hour unconstrained).
- Validates active bot identity (`gitmap login set ai-auto <agent>`).
- Enforces strict boundary checks and passes `--tz`, `--schedule`, and `--non-interactive` to downstream commit ceremonies.

---

## 4. Acceptance Criteria & Invariants

1. `gitmap cat <filepath> -L <start>:<end>` correctly slices lines and renders headers.
2. `gitmap cat <filepath> -q <query>` accurately filters matching lines.
3. `gitmap cat <filepath> --head <N>` and `--tail <N>` display exact counts.
4. `38-sync-prompts-skills-scripts-v7.py` enforces V7 restrict conditions cleanly.
5. `.agents/skills/gitmap/skill.md` in both `coding-guidelines` and `gitmap` contains Rule 19 & Rule 20.
6. All modified Go functions adhere to <= 15 lines and positive booleans.
7. Fleet sync succeeds without errors.
