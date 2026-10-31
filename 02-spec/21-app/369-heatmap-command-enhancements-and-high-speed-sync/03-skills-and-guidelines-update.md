# Heatmap & GitMap Skills Synchronization Specification

> **Task:** 369-heatmap-command-enhancements-and-high-speed-sync  
> **Subsystem:** Skills & Coding Guidelines (`.agents/skills/`, `.cursor/skills/`)  
> **Date:** 31-Oct-26  
> **Status:** APPROVED  

---

## 1. Skill Updates Required

### 1.1 `.agents/skills/gitmap/skill.md`
- Update Rule 19 and Rule 23 to mandate using `gitmap heatmap <filepath>` or `gitmap cat <filepath>` when inspecting code.
- Explicitly state: Zero User-Defined Fields (UDFs) needed. Passing the file path automatically reads the file and records the heatmap touch in `repodb/heatmap.db`.
- Total ban on ad-hoc Python reading (`gitmap py -c "open('...').readlines()"`).

### 1.2 Dedicated Heatmap Skill (`.agents/skills/gitmap/` or new heatmap skill references)
- Document the auto-reading syntax:
  ```bash
  gitmap heatmap <filepath> [-L start:end] [--num]
  ```
- Document that `gitmap cat <filepath>` automatically populates heatmap cognitive records.
- Document fleet sync commands:
  ```bash
  gitmap cg sync <repo>
  gitmap cg skill sync <repo>
  gitmap cg prompt sync <repo>
  ```

---

## 2. Cross-Fleet Deployment Checklist
- [x] Canonical source updated in `gitmap`.
- [x] Fleet sync executed via `gitmap cg sync` across all connected repositories.
- [x] Human edit sidecars (`.upstream`) reported if conflicts are detected.
