---
name: sync-codebase
description: Autonomously synchronize prompts, skills, coding guidelines, and AI scripts across target repositories using parameter-driven paths, automated backup branches, and non-negotiable boundaries.
---

# Multi-Repository Codebase Synchronization (Alias)

Source prompt: `01-prompts/24-sync/02-sync-other-codebase.md`
Companion skill: `.agents/skills/sync-other-codebase/skill.md`

## Instructions

1. Read `01-prompts/24-sync/02-sync-other-codebase.md` in full before doing the task.
2. Follow `sync-other-codebase` specifications.
3. Pre-flight pull on base branches (`git pull origin <base> --no-rebase`).
4. Create and push safety backup branches (`backup/sync-<timestamp>`).
5. Enforce 5 Non-Negotiable Boundaries:
   - Spec 21 Exclusion: Never sync or touch `02-spec/21-*`.
   - Bump Script Protection: Never overwrite version bump scripts.
   - Additive-Only AI Scripts: Copy new scripts, diff and inspect existing scripts before modifying.
   - Memory & Plans Protection: Never overwrite or delete `.ai-memory/memory/` or `.ai-memory/plans/` in target repositories.
   - Zero Secrets: Never copy `.env` or credentials.
6. Commit changes and push safely per repository.
