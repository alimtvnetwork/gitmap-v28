# Subtask 06: V6 Prompt Synchronization, Fleet Sync & SemVer Minor Release Ceremony

## Status: READY
- **Task ID:** `315-subtask-06`
- **Parent Spec:** [02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/00-overview.md](../../../../02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/00-overview.md)
- **Module:** `01-prompts/14-execute`, `.agents/skills/execute-parent-task-with-n-steps-v6`, `cli/constants`, Release
- **Owned Files:**
  - `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`
  - `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`
  - `cli/constants/constants.go`
  - `.ai-memory/release/release-notes-v6.593.0.md`
  - `.ai-memory/plans/completed/315-gitmap-commit-tracking-and-database-integration.md`

---

## 1. Objectives

1. **Update Canonical V6 Execution Prompt in `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`:**
   - Extend Phase 3 execution and commit instructions: instruct AI agents to pass `--agent-id`, `--heatmap-id`, and `--task` flags whenever invoking `gitmap cpf` or intermediate `gitmap commit` commands.
   - Document ambient environment fallbacks (`ANTIGRAVITY_AGENT_ID`, `GITMAP_TASK`, `GITMAP_HEATMAP_ID`, `AI_MODEL`) so agents understand automated metadata ingestion.
   - Reiterate the non-negotiable hyphen-separated commit message invariant (`<module> - <summary>`) without colons.

2. **Update V6 Execution Skill in `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`:**
   - Synchronize prompt updates into the skill markdown to ensure all autonomous subagents inherit the updated commit attribution behavior.
   - Add explicit guidance on zero unattributed AI commits in swarm workflows.

3. **Fleet Synchronization across Connected Repositories:**
   - Execute `03-ai-scripts/38-sync-prompts-scripts-and-skills.py` (or repository sync toolchain) to distribute updated prompts, skills, and guidelines across the workspace fleet.

4. **Targeted Quality Gate Validation:**
   - Execute `03-ai-scripts/01-check-relative-paths.py` (verify 100% relative paths; zero drive letters).
   - Execute `03-ai-scripts/18-check-enum-and-boolean.py` (verify positive booleans and zero nested ifs).
   - Execute `03-ai-scripts/27-check-forbidden-strings.py` (verify zero forbidden strings or stale identifiers).

5. **Execute SemVer Minor Release Ceremony (`v6.593.0`):**
   - Increment minor version in `cli/constants/constants.go` to `6.593.0`.
   - Author comprehensive release notes in `.ai-memory/release/release-notes-v6.593.0.md`.
   - Transition master execution plan from `.ai-memory/plans/pending/` to `.ai-memory/plans/completed/`.
   - Lead Orchestrator performs single atomic commit via `gitmap cpf` with provenance flags:
     `gitmap cpf "commit-tracking - universal commit info, split-db history, ai provenance flags, and v6 prompt sync" --agent-id <id> --heatmap-id <id> --task gitmap-commit-tracking-and-database-integration`.
   - Push release branch `release/v6.593.0`, tag `v6.593.0`, and `main` upstream.

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Update V6 Execution Prompt in `01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`
1. Locate Phase 3 Commit & Push instructions (around line 240+).
2. Insert mandatory AI provenance flag instructions:
   ```markdown
   > **AI Commit Provenance Mandate (Task 315):**
   > When creating commits via `gitmap cpf` (or `gitmap cm`), AI agents MUST pass provenance flags:
   > `gitmap cpf "<module> - <summary>" --agent-id <id> --heatmap-id <id> --task <slug>`
   > Alternatively, ensure ambient environment variables are set (`ANTIGRAVITY_AGENT_ID`, `GITMAP_TASK`, `GITMAP_HEATMAP_ID`).
   > Flags are stripped before creating the commit message; never put flags inside the message string.
   ```
3. Verify formatting and save.

### Step 2: Synchronize V6 Skill in `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`
1. Replicate prompt changes into `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`.
2. Ensure frontmatter description accurately reflects V6 provenance tracking.

### Step 3: Run Multi-Repository Synchronization Script
1. Execute `python 03-ai-scripts/38-sync-prompts-skills-scripts.py` across connected repositories.
2. Confirm prompt and skill mirrors are updated without conflicting changes.

### Step 4: Execute Quality Linters
1. Run `python 03-ai-scripts/01-check-relative-paths.py`:
   - Verify zero absolute paths (`C:`, `D:`, `file:///`) in all modified files.
2. Run `python 03-ai-scripts/18-check-enum-and-boolean.py`:
   - Verify all booleans follow positive naming (`is*`, `has*`) and functions have zero nested ifs.
3. Run `python 03-ai-scripts/27-check-forbidden-strings.py`:
   - Confirm zero forbidden strings.

### Step 5: Version Bump & Release Notes Authoring
1. Update `cli/constants/constants.go`:
   - Update `Version = "6.593.0"`.
2. Author `.ai-memory/release/release-notes-v6.593.0.md` detailing:
   - Executive summary of Task 315.
   - Universal `gitmap commit info <sha>` feature with ASCII tree visuals.
   - External Split-DB `repodb/commit_history.db` architecture.
   - Provenance flags and ambient fallbacks.
   - Guideline 12 and V6 prompt synchronization.
3. Move `.ai-memory/plans/pending/315-gitmap-commit-tracking-and-database-integration.md` to `.ai-memory/plans/completed/`.

### Step 6: Atomic Commit & Upstream Push
1. Confirm working tree status via `gitmap status`.
2. Execute atomic commit via Lead Orchestrator:
   ```bash
   gitmap cpf "commit-tracking - universal commit info, split-db history, ai provenance flags, and v6 prompt sync" --agent-id lead-orchestrator --task gitmap-commit-tracking-and-database-integration
   ```
3. Create annotated tag `v6.593.0` and push upstream.

---

## 3. Verification & Quality Gates

| Quality Gate | Tool / Script | Expected Outcome |
| :--- | :--- | :--- |
| **Relative Path Hygiene** | `03-ai-scripts/01-check-relative-paths.py` | 100% PASS: Zero absolute paths or drive letters. |
| **Boolean & Enum Conventions** | `03-ai-scripts/18-check-enum-and-boolean.py` | 100% PASS: Positive boolean conventions maintained. |
| **Prompt Synchronization** | `03-ai-scripts/38-sync-prompts-skills-scripts.py` | 100% PASS: All prompts and skills synchronized. |
| **Working Tree Cleanliness** | `git status` | 100% PASS: Zero untracked or dangling files. |
| **Upstream Synchronization** | `git push` | 100% PASS: Release branch and tag `v6.593.0` pushed upstream. |

---

## 4. Definition of Done
- [ ] V6 execution prompt and skill updated with commit provenance instructions.
- [ ] Connected repositories synchronized via `38-sync-prompts-skills-scripts.py`.
- [ ] All targeted quality linters pass with 100% green status.
- [ ] Minor version bumped to `v6.593.0` in `cli/constants/constants.go`.
- [ ] Release notes published in `.ai-memory/release/release-notes-v6.593.0.md`.
- [ ] Master plan transitioned to completed.
- [ ] Atomic commit and tag `v6.593.0` pushed upstream.
