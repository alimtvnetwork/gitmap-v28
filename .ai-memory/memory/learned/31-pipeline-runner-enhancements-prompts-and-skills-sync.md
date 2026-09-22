# 31 — Pipeline Runner Enhancements, Prompts & Skills Synchronization

- **Slug:** pipeline-runner-enhancements-prompts-and-skills-sync
- **Date:** 2026-09-22
- **Version:** v6.297.0
- **Category:** learned
- **Status:** permanent

---

## 1. Executive Summary

In response to the recent updates to the CI/CD pipeline runner (`03-ai-scripts/06-cicd-local-runner.py`) and GitMap's diagnostic CLI tools (`gitmap pipeline-ai`, `gitmap pe`, `gitmap pd`), the following prompts and skills were enhanced across both `coding-guidelines` and `gitmap`:
1. `release-orchestrator` prompt (`01-prompts/17-release-management/06-release-orchestrator.md`)
2. `release-orchestrator` skill (`.agents/skills/release-orchestrator/skill.md`)
3. `ci-cd-fix-with-release` prompt (`01-prompts/16-ci-cd/06-ci-cd-fix-with-release.md`)
4. `ci-cd-fix-with-release-tweak` prompt (`01-prompts/16-ci-cd/02-ci-cd-fix-with-release-tweak.md`)
5. `ci-cd-fix-with-release` skill (`.agents/skills/ci-cd-fix-with-release/skill.md`)

All prompt files, skills, and AI scripts were mirrored, committed, and synchronized to GitHub across both `coding-guidelines` (`e2946b33`, `ea3379a7`) and `gitmap` (`4813ef28`).

---

## 2. CI/CD Pipeline Runner & GitMap Diagnostics Integration

The updated prompts and skills now explicitly document and mandate the following capabilities:

### A. Priority Incremental Runner Shortcuts
- `python 03-ai-scripts/06-cicd-local-runner.py run-smart` (or alias `smart`, `--smart`, `-s`): The fast path that builds only changed Go packages into OS temp and runs the Quad Runner.
- `python 03-ai-scripts/06-cicd-local-runner.py run-incremental` (or `incremental`, `--incremental`).
- Direct package targeting: `python 03-ai-scripts/06-cicd-local-runner.py --pkg <target>` (`--package`, `-p`, `--target-file`, `--file`).
- Heatmap-driven fast path: `python 03-ai-scripts/06-cicd-local-runner.py --fast` (runs only hot and warm tests based on `.ai-memory/test-heatmap.json`, skipping cold tests).
- Change-scoped linting: `python 03-ai-scripts/06-cicd-local-runner.py --changed-only` (default commit window: `--commits 20`).
- Cross-OS build gates: `python 03-ai-scripts/06-cicd-local-runner.py --all-os`.
- Standard development mode (`--no-tests`) vs pre-release verification (`--run-tests`).

### B. GitMap Pipeline AI Suite & Diagnostics
- `gitmap pipeline-ai status --json` / `gitmap pl-ai status --json`: Sub-50ms check of remote CI workflow state, active branch, and ETA.
- `gitmap pipeline-ai status -t <etaSeconds>`: Dynamic timeout wait for running pipelines without burning tokens or user credits in tight loops.
- `gitmap pipeline error-logs` (shortcut `gitmap pe`): Extracts failing step logs directly to file for 4-part RCA.
- `gitmap pe clear -y`: Clears stale past error logs before starting a fresh diagnostic run.
- `gitmap pipeline details` (shortcut `gitmap pd`): Runner targets table, SQLite caching, and branch status.
- `gitmap pipeline purge`: Zero-storage GitHub Actions purge maintaining 0.0 GB footprint.
- `gitmap db status` and `gitmap db prune`: SQLite database health check and 10MB auto-pruning ceiling.

### C. Bounded Stack Trace Extraction (RCA 58)
- Extractions must always be bounded to the failing failure frame (strictly 5 preceding + 20 trailing lines via `extractBoundedStackLines`), halting on exit code or job boundary to prevent clipboard explosion.

---

## 3. Synchronization Ledger

- `coding-guidelines`: Updated prompts & skills, committed and pushed (`e2946b33`, `ea3379a7`).
- `gitmap`: Updated prompts & skills mirrored via `38-sync-prompts-skills-scripts.py`, committed and pushed (`4813ef28`).
- Working trees: Clean and verified on `main` across both repositories.
