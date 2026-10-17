# 50 — Prompts + Skills Sync to Targets (`50-sync-prompts-skills-to-targets.py`)

Re-runnable, idempotent sync that pushes the canonical **prompts and skills**
from coding-guidelines-v24 into target repos that **do not need spec**
(e.g. `story-telling-v1`).

## What syncs

| # | Source (this repo) | Target destination |
|---|---|---|
| 1 | `01-prompts/` | `<target>/01-prompts/` |
| 2 | `.agents/skills/` | `<target>/.agents/skills/` |
| 3 | `.cursor/skills/` | `<target>/.cursor/skills/` |

## What NEVER syncs (fail-safe default-deny)

- **`02-spec/` — NEVER.** Not shared specs (01–20), not app specs (21+), nothing.
  The script allows ONLY the three roots above (`SYNC_ROOTS`); the `is_denied()`
  gate rejects `02-spec/` on every copy and walk path. This is the program's
  hard boundary: repos that do not need spec must never receive spec.
- **`.ai-memory/`** — target repos own their memory and plans.
- **`03-ai-scripts/`** — sync scripts ship via the fleet sync (script 38) or by hand.
- Build/cleanup debris: `.pyc`, `.pyo`, `.tmp`, `.DS_Store`, `__pycache__`.

## How to add a target

**Option A — registry file** (recommended): add an entry to
`scripts/prompt-sync-targets.json`:

```json
{
  "name": "my-repo",
  "note": "what this repo is",
  "path": "/absolute/path/to/my-repo"
}
```

**Option B — CLI** (repeatable, bypasses the registry):

```bash
gitmap py 03-ai-scripts/50-sync-prompts-skills-to-targets.py --target /path/to/my-repo --target /path/to/other-repo
```

The registry ships with `story-telling-v1` as the example target.

## How to run

```bash
# List registry targets
gitmap py 03-ai-scripts/50-sync-prompts-skills-to-targets.py --list

# Preview changes (no writes)
gitmap py 03-ai-scripts/50-sync-prompts-skills-to-targets.py --dry-run

# Sync every registry target
gitmap py 03-ai-scripts/50-sync-prompts-skills-to-targets.py

# Sync one ad-hoc target (safe playground: ~/workspace/custom-skills has no git)
gitmap py 03-ai-scripts/50-sync-prompts-skills-to-targets.py --target ~/workspace/custom-skills
```

## Re-runnable + idempotent

- Files are compared by SHA-256 content hash: identical files are skipped.
- Canonical skills/prompts update in place; **target-local skill folders
  (present only in the target) are kept, never deleted** — same convention as
  the fleet sync since `0c5fad76`.
- Second run reports 0 changes: `changed` → `clean (idempotent no-op)`.

## Per-target report

Each run prints per-target counts (copied new / updated / skipped-identical /
kept target-local / denied-skipped), the fail-safe proof line
(`02-spec/ denied by is_denied() gate`), then a summary table.

## Relation to script 38

Script `38-sync-prompts-skills-scripts.py` is the fleet ceremony (43 repos,
backup branches, release tags, and it *does* carry shared `02-spec/01-*–20-*`).
Script 50 is the light, spec-free variant for small target repos: no git
ceremony, no releases — just prompts + skills, with `02-spec/` hard-denied.
