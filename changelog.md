# Changelog

## [v6.521.0] - 2026-10-09

### Added
- Routine release v6.521.0

---

## [v6.520.0] - 2026-10-09

### Added
- DB-backed spec number issuance (gitmap spec next)

---

## [v6.519.0] - 2026-10-09

### Added
- Routine release v6.519.0

---

## [v6.518.0] - 2026-10-09

### Added
- Routine release v6.518.0

---

## [v6.517.0] - 2026-10-09

### Fixed
- `gitmap agm update` terminal UI: coordinated pterm progress UI (spinner + live installer log tail + stage checkmarks) fed by piped child output — no more minutes of silence, no raw installer escape sequences corrupting the screen.
- `gitmap agm update` no longer misrouted to `gitignore`: the `agm` alias collision in the core dispatcher is removed, so the documented command reaches the real updater.
- Failure path: structured error panel (failed stage, last installer lines, next-step hint) instead of raw output dumps + Go stack traces; verification failures point at the diagnostic log file.
- `agm update --help` now lists the `--yes` flag.

### Added
- `GITMAP_AGM_INSTALL_SCRIPT` env override for the AGM installer source (test seam; default behavior unchanged).

---

## [v6.515.0] - 2026-10-08

### Added
- multi-user project switching, github cli auth status, and git config manager

---

## [v6.514.0] - 2026-10-08

### Added
- user - multi-user project switching, github cli auth status, and git config manager

---

## [v6.513.0] - 2026-10-08

### Added
- Concise bare `gitmap` CLI display (< 15 lines) with structured suggestions:
  - Sub-point 1: Catalog exploration via `gitmap help` and `gitmap -h`.
  - Sub-point 2: Mandatory AI model self-training via `gitmap llm train` (and `gitmap ai llm-train`) before starting workspace tasks.
- Elimination of redundant blank line between binary location triplet and version identity block.
- Authoritative `gitmap llm train` subsystem:
  - Emits full copy-pasteable Antigravity skill with YAML frontmatter directly to stdout.
  - Automatically writes and updates `.agents/skills/gitmap/SKILL.md`.
  - Added `--urls` flag emitting public raw GitHub documentation URLs for LLM memory ingestion.
  - Directs autonomous AI models to generate and persist their own skill from training output.
  - Embeds recursive Git network learning directives (multi-repo topology, prior mistake prevention via `gitmap pe history-ai`, dynamic ETA wait, memory ledgers).
- Modernized skills documentation in `.agents/skills/gitmap/SKILL.md` and `.cursor/skills/gitmap/skill.md`.

---

## [v6.512.0] - 2026-10-08

### Added
- suggestion engine, fastgate pre-commit runner, and muse help clusters

---

## [v6.511.0] - 2026-10-08

### Added
- fix pipeline pe -1 nested ifs boolean guidelines relative paths and misspell

---

## [v6.510.0] - 2026-10-08

### Added
- chrome profile import-export e2e verification and test logging in repo-secrets

---

## [v6.509.0] - 2026-10-08

### Added
- cli build verification, cpar documentation, and chrome profile automation

---

## [v6.508.0] - 2026-10-08

### Added
- universal file runner, supabase vault, cpc chore commits, pe-ai and muse skills v6.508.0

---

## [v6.507.2] - 2026-10-08

### Added
- sync - add native fleet sync and sqlite agent task engine

---

## [v6.507.1] - 2026-10-07

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.507.0] - 2026-10-07

### Added
- pipeline pe cache invalidation relative path sanitization and error extraction fix v6.507.0

---

## [v6.506.2] - 2026-10-07

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.506.1] - 2026-10-07

### Added
- Synchronize prompts, skills, AI scripts, and coding guidelines

---

## [v6.506.0] - 2026-10-07

### Added
- **Cursor Settings Guard Validation**: Added `requireCursorSettingsFile` in `cli/cmdcursor/cursor_settings.go` with descriptive `E1037` error envelopes prompting users to run `gitmap cursor settings apply`.
- **Views AI Settings UI Modernization**: Refactored `src/pages/Settings.tsx` to invert `/api/instances` fetch validation into clean guard clauses, eliminating nested `if` statements.
- **Process Discovery Spacing Hygiene**: Enforced vertical blank line spacing rules in `cli/cmdagy/agy_instance_discovery.go` and `cli/cmdcursor/cursor_settings_test.go`.
- **Quality Gates Verification**: Verified zero violations across `check-nested-ifs.py`, `check-enum-and-boolean.py`, and `check-relative-paths.py`.

---

## [v6.505.0] - 2026-10-06

### Added
- history purge dispatch fix, recycle bin export, and relative path sanitation

---

## [v6.504.0] - 2026-10-06

### Added
- `gitmap login`: unified GitHub authentication — `--token <PAT>` (validated via api.github.com), `--web` browser login (gh device flow or guided token page), interactive chooser, `--status`
- `gitmap logout`: removes the stored GitHub credential
- Logged-in token is reused automatically by clone/pull/push via token resolution
- Secure no-echo token prompt; fail-fast guidance when stdin is not a TTY

---

## [v6.503.0] - 2026-10-06

### Added
- chrome profile import export e2e test verification and test suite fixes

---

## [v6.502.0] - 2026-10-06

### Added
- resolve CI/CD sends push assertion and ERD parity cluster path

---

## [v6.501.0] - 2026-10-06

### Added
- Deep spec consolidation, canonical 8-cluster reduction, and memory unification

---

## [v6.500.0] - 2026-10-06

### Added
- Pre-deep-consolidation baseline release

---

## [v6.499.0] - 2026-10-06

### Added
- App spec and completed plans consolidation and reduction

---

## [v6.498.0] - 2026-10-06

### Added
- Pre-consolidation baseline release

---

## [v6.497.0] - 2026-10-06

### Added
- **Ubuntu IDE and GitHub Desktop Scan Sync Suite**: Resolved omission where `gitmap scan` on Ubuntu failed to add repositories into VS Code and GitHub Desktop.
- **Root Cause Analysis (Issue 71)**: Grounded 4-part RCA in `02-spec/22-app-issues/71-ubuntu-ide-github-desktop-scan-omission.md` identifying missing directory creation in VS Code Project Manager extension path and missing Linux CLI resolution candidates for GitHub Desktop.
