# Subtask 01: Audit, Scan & Tooling Infrastructure for Enums & Constants

Parent Plan: [140-constants-and-enums-architecture.md](../../pending/140-constants-and-enums-architecture.md)

## Goals
1. Scan entire codebase for enums missing `*Type` suffix and raw rune casts.
2. Author `.agents/skills/cg-enums/skill.md` documenting coding guideline enforcement for enums and constants.
3. Upgrade `linter-scripts/check-enum-guidelines.py` with comprehensive checks.

## Acceptance Criteria
- [x] `.agents/skills/cg-enums/skill.md` created with YAML frontmatter.
- [x] `check-enum-guidelines.py` verified with Go enum checks (detected 11 violations).
- [x] Subtask completed and logged.
