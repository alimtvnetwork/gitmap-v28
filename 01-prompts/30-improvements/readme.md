# Prompt Category 30: Codebase Improvements Architecture

This directory houses the canonical, production-grade prompts for discovering, cataloging, and refactoring 100-item improvement campaigns across the repository:

| Prompt File | Domain & Primary Focus | Key Invariants & Tenets |
|---|---|---|
| [`01-find-100-improvements-in-ui-level.md`](01-find-100-improvements-in-ui-level.md) | **Terminal UI & Web UX Layer** | Lipgloss adaptive color tokens, TrueColor 24-bit ANSI, responsive column folding, multi-track progress bars, and high-contrast accessibility. |
| [`02-find-100-improvements-in-codebase-coding-guidelines.md`](02-find-100-improvements-in-codebase-coding-guidelines.md) | **Source Code & Error Management** | Leaf error return sovereignty, zero swallowed errors, 15-line function caps, zero nested ifs, positive booleans (`is*`, `has*`), and Split-DB repository encapsulation. |
| [`03-find-100-improvements-in-whole-system-ux.md`](03-find-100-improvements-in-whole-system-ux.md) | **Whole-System UX & Operational Resilience** | Cause-and-effect analysis, non-destructive execution, autostash protections, safe defaults, zero data loss, and backward compatibility. |

---

## Associated Skills & Toolchain Mappings

Each canonical prompt is mirrored directly into both Google Antigravity Agent Skills and Cursor Skills:

| Prompt Reference | Antigravity Agent Skill | Cursor Skill |
|---|---|---|
| `01-find-100-improvements-in-ui-level.md` | [find-100-ui-improvements](../../.agents/skills/find-100-ui-improvements/skill.md) | [find-100-ui-improvements](../../.cursor/skills/find-100-ui-improvements/skill.md) |
| `02-find-100-improvements-in-codebase-coding-guidelines.md` | [find-100-codebase-improvements](../../.agents/skills/find-100-codebase-improvements/skill.md) | [find-100-codebase-improvements](../../.cursor/skills/find-100-codebase-improvements/skill.md) |
| `03-find-100-improvements-in-whole-system-ux.md` | [find-100-system-ux-improvements](../../.agents/skills/find-100-system-ux-improvements/skill.md) | [find-100-system-ux-improvements](../../.cursor/skills/find-100-system-ux-improvements/skill.md) |

---

## Multi-Repository Synchronization

This entire directory is part of the canonical prompt suite synchronized across all 43 connected repositories via:
`03-ai-scripts/38-sync-prompts-skills-scripts.py`
