# 92-sync-prompts-and-skills-from-coding-guidelines

## Execution Summary
- **Started By:** User request to sync all prompts and skills from `D:\wp-work\riseup-asia\coding-guidelines` under Parent Task N-Step Continuous Loop & Multi-Agent Orchestration Workflow (v2.1.0, N=150).
- **Duration / Loops:** Completed in 1 self-loop cycle using 2 concurrent planning subagents followed by 2 concurrent execution subagents.
- **Status:** Completed & Consolidated.

---

## Accomplishments

### 1. Skills Synchronization (`.agents/skills/`)
- **13 New Skills Added:**
  - `autonomous-qa-and-testing/skill.md`
  - `ci-cd-create/skill.md`
  - `clean-artifacts-and-git-history/skill.md`
  - `execute-ai-instruction-writer/skill.md`
  - `execute-coding-guideline-fix/skill.md`
  - `execute-parent-task-with-n-steps/skill.md`
  - `fix-spec-from-audit/skill.md`
  - `fix-with-rca/skill.md`
  - `inventory-pending-tasks/skill.md`
  - `plan-coding-guideline-audit/skill.md`
  - `spec-reverse-engineering/skill.md`
  - `write-antigravity/skill.md`
  - `write-memory/skill.md`
- **9 Truncated Stub Skills Upgraded to Full Engines:**
  - `ci-cd-fix/skill.md` (upgraded from 27-line stub to full 352-line engine; adapted `02-spec/` to `spec/`)
  - `coding-guidelines/skill.md` (upgraded from 44-line stub to full 671-line engine; adapted `02-spec/` to `spec/`)
  - `execute-batched-loop/skill.md` (upgraded from 21-line stub to full 226-line engine; renamed to lowercase `skill.md`)
  - `execute-batched-loop-wor/skill.md` (upgraded from 21-line stub to full 217-line engine; renamed to lowercase `skill.md`)
  - `execute-pending-tasks/skill.md` (upgraded from 28-line stub to full 232-line engine)
  - `read-memory-enhanced/skill.md` (updated with git commit inspection and lowercase conventions)
  - `cg-boolean-and-naming/skill.md` (aligned local runner path to `03-ai-scripts/06-cicd-local-runner.py`)
  - `cg-error-management/skill.md` (aligned local runner path to `03-ai-scripts/06-cicd-local-runner.py`)
  - `execute-parent-task/skill.md` (upgraded with full orchestration protocol, keeping name `execute-parent-task`)
- **Case Normalization:**
  - Renamed uppercase `SKILL.md` to lowercase `skill.md` across all skill directories.
- **Gitmap-Specific Skills 100% Preserved:**
  - `cluster-ssh-delegation`, `db-sqlite-conventions`, `gitmap-cli-command`, `parallel-cpu-checkers-orchestration`, `powershell-cross-platform-scripting`, `react-docs-frontend`, `release-and-versioning`, `release-orchestrator`, `spec-authoring`.

### 2. Rules Synchronization (`.agents/rules/`)
- Synced 5 canonical rules and directory anchor:
  - `boolean-standards.md`
  - `code-style-and-formatting.md`
  - `error-management.md` (adapted with `*apperror.AppError`)
  - `go-guidelines.md` (adapted with `*apperror.AppError`)
  - `micro-tasking-and-batching.md`
  - `.gitkeep`
- Preserved `gitmap_guidelines.md` completely intact.

### 3. Prompts Catalog Synchronization (`01-prompts/`, `prompts/`, `.lovable/prompts.md`)
- Synced all 21 categories containing 89 markdown prompt files and 2 `.gitkeep` files from `D:\wp-work\riseup-asia\coding-guidelines\01-prompts/`.
- Merged `01-prompts/01-prompt-library-setup/01-prompt-library-setup.md` with full canonical specification and preserved Gitmap's active macro builder prompts in `## Instructions Index`.
- Copied `prompts/001-conversation-log-and-context-wrapper.md` while preserving Gitmap's `001-interactive-macro-builder-commands.md` and `002-gitmap-interactive-builder-and-clean-sync.md`.
- Installed `.lovable/prompts.md` (Version 3.0.0 master index).
- Created `.lovable/01-index.md` as repository-level context & directory router.
- Purged obsolete `01-prompts/01-release.md` and orphaned `01-general-prompts/`.

### 4. Verification & Quality Gates
- Ran `python 03-ai-scripts/10-encoding-normalizer.py --fix` across all 9,363 files, ensuring 100% UTF-8 without BOM and UNIX LF line endings.
- Verified pure relative paths across all markdown links.
