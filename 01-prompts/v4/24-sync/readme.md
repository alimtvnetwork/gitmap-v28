# Multi-Repository Synchronization Prompts (`24-sync`) — Index & Catalog

> [!IMPORTANT]
> Category: `24-sync`
> Architecture: V6 Parameter-Driven Multi-Agent Protocol
> Scope: Canonical Prompts, Skills, Guidelines, and AI Scripts Distribution

---

## Overview

The `24-sync` prompt category governs autonomous, safe, and lossless synchronization of canonical architecture from the central Prompt Architect repository to connected downstream repositories.

Downstream codebases require standardized prompts (`01-prompts/`), agent skills (`.agents/skills/`, `.cursor/skills/`), coding guidelines (`02-spec/02-coding-guidelines/`), and automation scripts (`03-ai-scripts/`, `.agents/scripts/`). However, synchronization must never destroy repository identity, erase custom application logic, overwrite project-specific versioning tools, or leak sensitive credentials.

This category codifies the non-negotiable boundaries, preflight safety protocols, branch protection ceremonies, and multi-agent execution pipelines required to mirror updates seamlessly across multiple repositories.

---

## Directory Index & Catalog

| Prompt File | Version | Scope | Key Capabilities |
| :--- | :--- | :--- | :--- |
| [`01-sync.md`](01-sync.md) | V6 (6.0.0) | Full Fleet (43 Repos) | Multi-agent autonomous synchronization of canonical prompts (`01-prompts/`), skills, shared specs (`02-spec/01-*` to `02-spec/20-*`), and additive AI scripts across all 43 connected repositories with automated pre-flight, backup branches, 5 non-negotiable boundaries, and atomic GitMap commits. |
| [`02-sync-other-codebase.md`](02-sync-other-codebase.md) | V6 (6.0.0) | Parameter-Driven (Cross-Repo) | Multi-agent autonomous synchronization of canonical prompts, skills, guidelines, and additive scripts with backup branching, 5 non-negotiable boundaries, zero hardcoded paths, dynamic `SOURCE_REPO` / `TARGET_REPOS`, and full release ceremony. |

---

## Core Synchronization Architecture

### 1. The 5 Non-Negotiable Boundaries

Every synchronization operation across connected codebases must strictly enforce five core boundaries:

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                   5 NON-NEGOTIABLE SYNCHRONIZATION BOUNDARIES               │
├─────────────────────────────────────────────────────────────────────────────┤
│ 1. Spec 21 Exclusion      │ NEVER touch 02-spec/21-* (app/domain specs).    │
│ 2. Bump Script Protection │ NEVER overwrite repo-specific bump scripts.     │
│ 3. Additive-Only Scripts  │ Copy new AI scripts; diff/preserve modified.   │
│ 4. Memory & Plans Safe    │ NEVER touch .ai-memory/memory/ or plans/.       │
│ 5. Zero Secrets Leakage   │ NEVER sync .env or private credentials.         │
└─────────────────────────────────────────────────────────────────────────────┘
```

1. **Spec 21 Exclusion (TOTAL BAN):**
   - Directories matching `02-spec/21-*` (such as `21-app/`, `21-app-issues/`, `21-app-db/`, `21-app-ui-design-system/`) belong strictly to the target repository's private application domain.
   - Synchronizing or touching these directories is strictly forbidden as it overwrites domain logic and wipes application identity.

2. **Bump Script Protection (IMMUTABLE):**
   - Files containing version bumping logic (e.g. `bump-version.mjs`, `bump_versions.py`, `37-bump-version.py`, `scripts/bump*.py`, `version.json`) have custom hooks, tags, and file targets per repository.
   - These files must never be overwritten by central generic bump scripts.

3. **Additive-Only AI Scripts (NO BLIND OVERWRITES):**
   - New utility scripts in `03-ai-scripts/` and `.agents/scripts/` that do not exist in the target repository are copied cleanly.
   - Existing scripts in the target repository that share names but contain repo-specific modifications must not be overwritten blindly; agents must inspect diffs and preserve local customizations.

4. **Memory & Plans Protection (TOTAL ISOLATION):**
   - Target repositories maintain their own operational logs, active execution plans, and architectural decisions inside `.ai-memory/memory/` and `.ai-memory/plans/`.
   - These areas are completely excluded from central mirroring.

5. **Zero Secrets Leakage (STRICT PURGE):**
   - Environment variables, tokens, API keys, and credentials (`.env`, `.env.*`, secret stores) must never be transferred across repositories.
   - Any credentials needed across development must reside in the `repo-secrets` folder in the default work directory via `gitmap rs`.

---

## Automated Tooling & Reference Scripts

- **Primary Synchronizer Script:**
  - `03-ai-scripts/38-sync-prompts-skills-scripts.py`: High-speed multi-threaded Python automation script that handles branch pulling, backup branch creation (`backup/pre-v3-nsteps-sync-<timestamp>`), directory mirroring, diff checking, and automated git tagging/release ceremony.
- **GitMap High-Speed Command Primacy:**
  - `gitmap` CLI: Used for live symbol searching (`gitmap aum search`), file discovery (`gitmap find`), atomic commits (`gitmap cpf`, `gitmap cpb`), and non-polling CI status checks (`gitmap pipeline-ai status`).
- **SQLite Concurrency & Crash Forensics:**
  - `03-ai-scripts/46-agent-sqlite-task-manager.py`: Atomic worker action logging and crash diagnostic engine.

---

## Parameter-Driven Execution (Zero Hardcoded Paths)

Prompts in this category enforce a strict **Zero Hardcoded Paths Mandate**. Paths are never embedded into prompt markdown or hardcoded into agent instructions.

Callers must supply parameters dynamically in the prompt preamble or invocation payload:

```text
SOURCE_REPO  = <path-to-source-repo>
TARGET_REPOS = [
    "<path-to-target-repo-1>",
    "<path-to-target-repo-2>",
]
```

The orchestrator reads these values at runtime and configures all subagent worker waves accordingly.

---

## Core Invariants

1. **Strict Relative Git Paths Only:** Only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well (TOTAL BAN on `file:///` URIs and absolute filesystem paths).
