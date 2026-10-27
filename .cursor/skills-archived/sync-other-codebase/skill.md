---
name: sync-other-codebase
description: Autonomously synchronize prompts, skills, coding guidelines, and AI scripts across target repositories using parameter-driven paths, automated backup branches, and non-negotiable boundaries.
---

# Cross-Repository Synchronization & Safe Propagation Engine

> **[/goal](slashCommand:goal)** Autonomously synchronize canonical prompts, Antigravity skills, coding guidelines, and automation scripts from a source repository into specified target repositories using dynamic parameter-driven paths, pre-flight safety backups, and strict isolation boundaries.
> **[/learn](slashCommand:learn)** Enforce the 5 Non-Negotiable Boundaries (Spec 21 Exclusion, Bump Script Protection, Additive-Only AI Scripts, Memory & Plans Protection, Zero Secrets), pre-flight backup branch creation, SemVer release ceremonies, and parameter-driven inputs without hardcoding source or target repository paths.

**Version:** 1.0.0
**Updated:** 2026-10-02
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## 1. Overview & When to Use

The `sync-other-codebase` skill governs cross-repository distribution of canonical prompts, agent skills, coding guidelines, and automation tooling across interconnected repositories in a multi-project workspace.

### When to Activate This Skill

- **Prompt Library Updates:** Canonical prompts in `01-prompts/` have been added, restructured, or updated with new autonomous workflows.
- **Skill Propagation:** New or improved Antigravity skills in `.agents/skills/` or Cursor skills in `.cursor/skills/` need replication across downstream codebases.
- **Shared Automation Scripts:** New or updated automation utilities in `03-ai-scripts/` or `.agents/scripts/` must be deployed to target repositories.
- **Coding Guideline Standardization:** Central specifications in `02-spec/02-coding-guidelines/`, `02-spec/07-design-system/`, or `.ai-memory/coding-guidelines.md` need synchronization across repositories.
- **Explicit User Request:** The user triggers phrases such as `sync other codebase`, `sync prompts across repos`, `sync downstream repos`, or `propagate skills`.

---

## 2. Parameter-Driven Input Interface (Zero Hardcoded Paths)

To guarantee portability, reliability, and security, this skill enforces a **strict ban on hardcoded repository paths**. Neither the source repository path nor downstream target paths may be hardcoded into skills, prompts, or automation configurations.

### Dynamic Input Parameters

Callers and orchestrators supply paths dynamically using environment variables, command-line arguments, or runtime parameters:

| Parameter | Type | Required | Description | Default / Resolution Fallback |
|---|---|---|---|---|
| `SOURCE_REPO` | Path / String | No | Path to source repository containing canonical assets | Automatically resolved from script parent directory (`Path(__file__).resolve().parent.parent`) |
| `TARGET_REPOS` | List / String | No | One or more paths to target repositories to sync | Discovered dynamically under workspace root or provided via `--repo` |
| `--repo <name>` | String | No | Target repository directory name under workspace root | Syncs all registered target repositories when omitted |
| `--dry-run` | Boolean | No | Preview all file additions, updates, and removals without modifying disk or git | `false` |
| `--no-push` | Boolean | No | Commit changes and create branches/tags locally without pushing to remote | `false` |
| `--workers <n>` | Integer | No | Concurrency level for parallel multi-repository sync | `6` |

### Parameter Resolution Examples

```bash
# 1. Preview synchronization for a specific target repository
python 03-ai-scripts/38-sync-prompts-skills-scripts.py --repo <target-repo-name> --dry-run

# 2. Synchronize a specific target repository with remote release push
python 03-ai-scripts/38-sync-prompts-skills-scripts.py --repo <target-repo-name>

# 3. Synchronize a target repository locally without pushing to origin
python 03-ai-scripts/38-sync-prompts-skills-scripts.py --repo <target-repo-name> --no-push

# 4. Synchronize all connected repositories concurrently
python 03-ai-scripts/38-sync-prompts-skills-scripts.py --workers 6
```

---

## 3. Pre-Flight Pull and Backup Branch Workflow

Before any target repository file is touched, copied, or deleted, every target repository MUST execute the pre-flight safety sequence to guarantee zero accidental data loss.

### Step 1: Base Branch Detection & Clean Pull

Identify the repository's base branch (`main` or `master`). If current HEAD is on an ephemeral branch (`release/*`, `backup/*`, `feat/*`), switch back to the base branch and pull the latest remote changes:

```bash
# Detect base branch and pull latest changes cleanly
git checkout <base_branch>
git pull origin <base_branch> --no-rebase
```

### Step 2: Timestamped Safety Backup Branch

Create and immediately push a dedicated backup branch capturing the exact pre-sync state of HEAD:

```bash
# Generate ISO-style timestamp
$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$backupBranch = "backup/sync-$timestamp"

# Create backup branch at current HEAD and push to remote
git checkout -b $backupBranch
git push -u origin $backupBranch

# Return to base branch for sync execution
git checkout <base_branch>
```

### Step 3: Pre-Change Release Tag Verification

Ensure current HEAD has a valid SemVer release tag (`vX.Y.Z`) and corresponding release branch (`release/vX.Y.Z`) pushed to remote before applying any modifications. This ensures an unambiguous rollback reference in git history.

---

## 4. The 5 Non-Negotiable Boundaries

All synchronization operations MUST strictly adhere to the five non-negotiable boundaries:

```
+-----------------------------------------------------------------------------------+
|                           5 NON-NEGOTIABLE BOUNDARIES                             |
+-----------------------------------------------------------------------------------+
| 1. Spec 21 Exclusion       : NEVER sync or touch 02-spec/21-*                     |
| 2. Bump Script Protection  : NEVER overwrite target version bump scripts          |
| 3. Additive-Only AI Scripts: New scripts copied; existing scripts never crushed   |
| 4. Memory & Plans Shield   : NEVER touch target .ai-memory/memory/ or plans/      |
| 5. Zero Secrets Mandate    : NEVER sync credentials, tokens, or .env files        |
+-----------------------------------------------------------------------------------+
```

### Boundary 1: Spec 21 Exclusion (Target Repo Exclusivity)

- **Target Directories:** `02-spec/21-*` (`02-spec/21-app/`, `02-spec/21-app-issues/`, `02-spec/21-app-db/`, `02-spec/21-app-ui-design-system/`).
- **Enforcement Rule:** Spec 21 is strictly exclusive to each individual target repository. Target repositories author their own application specifications and domain architecture. Synchronizing or mirroring Spec 21 is totally banned.
- **Safety Check:**
  ```python
  def is_spec_21(path: Path) -> bool:
      norm = str(path).replace("\\", "/").lower()

      if "/21-" in norm:
          return True

      if "spec/21" in norm:
          return True

      for part in path.parts:
          if part.startswith("21-"):
              return True

      return False
  ```

### Boundary 2: Bump Script Protection (Custom Version Logic Preservation)

- **Protected Files:** Any version bump script in target repositories (`bump-version.mjs`, `bump_versions.py`, `37-bump-version.py`, `scripts/bump-version.mjs`, etc.).
- **Enforcement Rule:** Target repositories maintain repository-specific version increment strategies (npm package.json, composer.json, Cargo.toml, Python pyproject.toml, Go version constants). Never overwrite or delete existing version bump scripts during sync.
- **Safety Check:**
  ```python
  def is_bump_script(path: Path) -> bool:
      name = path.name.lower()
      is_bump = "bump" in name
      is_version = "version" in name or name.startswith("bump")

      if is_bump:
          if is_version:
              return True

      return False
  ```

### Boundary 3: Additive-Only AI Scripts (No Blind Overwriting)

- **Target Directories:** `03-ai-scripts/` and `.agents/scripts/`.
- **Enforcement Rule:**
  - New automation scripts present in source but missing in target are copied cleanly.
  - Existing scripts in the target repository that were customized or modified locally MUST NOT be overwritten blindly. Inspect diffs and preserve local customizations.
  - Stale scripts existing in target repositories are NEVER deleted in additive mode.
- **Safety Check:**
  ```python
  if is_additive_only:
      if dst_file.exists():
          return 0  # Skip overwriting existing script
  ```

### Boundary 4: Memory & Plans Protection (Preserve Operational History)

- **Protected Directories:** `.ai-memory/memory/`, `.ai-memory/plans/`, `.ai-memory/temp-agents/`, `.ai-memory/cicd-issues/`, `.ai-memory/ambiguous-questions/`.
- **Enforcement Rule:** Target repositories own their sprint roadmaps, active plan files, completed subtasks, crash forensics logs, and domain memory. Sync operations MUST NEVER overwrite, mirror, or delete these directories.
- **Conditional Guideline Files:** Only static reference files (`.ai-memory/coding-guidelines.md`, `.ai-memory/prompts.md`) are synchronized when `.ai-memory/` exists.

### Boundary 5: Zero Secrets Mandate (No Credential Leakage)

- **Protected Files:** `.env`, `.env.*`, API keys, private certificates, SSH credentials, tokens, or credentials files.
- **Enforcement Rule:** Secrets must NEVER be checked into git repositories or copied during sync. All secrets belong strictly in `repo-secrets` via `gitmap rs`. Verify that no secret files or tokens are ever staged.

---

## 5. Synchronized Asset Directories

The following canonical directories are synchronized:

| Source Directory | Target Directory | Sync Mode | Notes |
|---|---|---|---|
| `01-prompts/` | `01-prompts/` | Mirror (Clean) | Never syncs `06-old-prompts/` (archived prompts) |
| `.agents/skills/` | `.agents/skills/` | Mirror (Clean) | All skill definitions use lowercase `skill.md` |
| `.cursor/skills/` | `.cursor/skills/` | Mirror (Clean) | Synchronized for Cursor IDE compatibility |
| `03-ai-scripts/` | `03-ai-scripts/` | Additive-Only | Never overwrites modified scripts or bump scripts |
| `.agents/scripts/` | `.agents/scripts/` | Additive-Only | Shared agent scripts, additive only |
| `02-spec/02-coding-guidelines/` | `02-spec/02-coding-guidelines/` | Mirror (Clean) | Canonical coding guidelines |
| `02-spec/07-design-system/` | `02-spec/07-design-system/` | Conditional | Synchronized only if target already has design system |
| `02-spec/17-consolidated-guidelines/` | `02-spec/17-consolidated-guidelines/` | Conditional | Synchronized only if target already contains directory |
| `.ai-memory/coding-guidelines.md` | `.ai-memory/coding-guidelines.md` | Conditional | Synchronized only if target `.ai-memory/` exists |
| `.ai-memory/prompts.md` | `.ai-memory/prompts.md` | Conditional | Synchronized only if target `.ai-memory/` exists |

---

## 6. Execution Commands & Automation Engine

Synchronization is driven natively by `gitmap sync` (compiled Go engine, primary) with `03-ai-scripts/38-sync-prompts-skills-scripts.py` as backward-compatible fallback.

### Primary GitMap Sync Commands (`gitmap sync` — Go Native)

```bash
# Step 1: Pre-flight dry run inspection
gitmap sync --repo <target-repo-name> --dry-run

# Step 2: Live synchronization execution
gitmap sync --repo <target-repo-name>

# Step 3: Custom projects via JSON file or inline JSON
gitmap sync --projects path/to/projects.json
gitmap sync --projects '[{"folder": "cat-my"}]'

# Step 4: Batch execution across all registered repositories concurrently
gitmap sync --workers 8
```

### Legacy Python Fallback (`03-ai-scripts/38-sync-prompts-skills-scripts.py`)

```bash
# Step 1: Pre-flight dry run inspection
python 03-ai-scripts/38-sync-prompts-skills-scripts.py --repo <target-repo-name> --dry-run

# Step 2: Live synchronization execution
python 03-ai-scripts/38-sync-prompts-skills-scripts.py --repo <target-repo-name>

# Step 3: Batch execution across all connected repositories
python 03-ai-scripts/38-sync-prompts-skills-scripts.py --workers 6
```

### Automation Engine Capabilities

1. **Dynamic Path Ingestion:** Dynamically locates target repos relative to workspace root without static hardcoded strings.
2. **Automated Branch & Tag Lifecycle:**
   - Pre-sync: `backup/sync-<timestamp>` and `v<pre_ver>`.
   - Work commit: `feat(sync): sync canonical prompts, skills, and coding guidelines`.
   - Post-sync: `release/v<next_ver>` and annotated tag `v<next_ver>`.
   - Merge back: Merges `release/v<next_ver>` back into base branch with `[skip ci]`.
3. **Multi-Threaded Parallel Execution:** Safely coordinates multiple worker threads via `ThreadPoolExecutor` with per-repository git isolation.

---

## 7. Post-Sync Verification & Release Ceremony

Following sync execution in each target repository, verify compliance before completing the task.

### Post-Sync Verification Checklist

- [ ] **Spec 21 Untouched:** Verify no files in `02-spec/21-*` were modified or added.
- [ ] **Bump Scripts Preserved:** Verify target repository bump scripts remain unaltered.
- [ ] **AI Scripts Additive-Only:** Verify existing target-customized scripts were not overwritten.
- [ ] **Memory & Plans Intact:** Verify `.ai-memory/plans/` and `.ai-memory/memory/` are unmodified.
- [ ] **Zero Secrets:** Verify no `.env`, token, or credential file exists in working tree or commit history.
- [ ] **Lowercase Hygiene:** Verify all synced files and directories adhere strictly to lowercase naming.
- [ ] **Backup Branch Pushed:** Verify `backup/sync-<timestamp>` exists on the remote repository.
- [ ] **Release Tag Pushed:** Verify `v<next_ver>` and `release/v<next_ver>` are pushed to origin.
- [ ] **Base Branch Clean:** Verify base branch is clean, up to date with origin, and has merged the release tag with `[skip ci]`.

### Summary Reporting Output

Every run outputs an ASCII status table detailing synchronization outcomes:

```text
===============================================================================================
MULTI-REPOSITORY SYNCHRONIZATION SUMMARY
===============================================================================================
Repository                         | Status   | Pre-Tag    | Post-Tag   | Files        | Action
-------------------------------------------------------------------------------------------------
target-repo-alpha                  | OK       | v1.2.4     | v1.2.5     | +14/-2       | Released & Pushed
target-repo-beta                   | OK       | v0.4.1     | v0.4.1     | +0/-0        | Clean
target-repo-gamma                  | OK       | v2.1.0     | v2.1.1     | +14/-0       | Released & Pushed
```

---

## 8. Traceability & Related Prompts

- **Canonical Execution Prompt:** `01-prompts/24-sync/02-sync-other-codebase.md`
- **Category Index:** `01-prompts/24-sync/readme.md`
- **Architecture Specification:** `02-spec/21-app/03-sync-other-codebase/01-architecture-spec.md`
- **Automation Engine:** `03-ai-scripts/38-sync-prompts-skills-scripts.py`
