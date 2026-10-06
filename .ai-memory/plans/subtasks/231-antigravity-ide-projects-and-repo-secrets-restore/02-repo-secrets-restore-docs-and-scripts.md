# Subtask 02: Repo-Secrets Migration Documentation & Automation Toolchain

> **Parent Plan:** `.ai-memory/plans/231-antigravity-ide-projects-and-repo-secrets-restore.md`  
> **Spec Reference:** `02-spec/21-app/231-antigravity-ide-projects-and-repo-secrets-restore/01-architecture-spec.md`  
> **Component Spec Reference:** `02-spec/21-app/231-antigravity-ide-projects-and-repo-secrets-restore/02-component-spec.md`  
> **Status:** `PENDING`  
> **Worker Assignment:** Worker 02  
> **Target Subsystems:** Migration Vault (`repo-secrets/09-antigravity-backup/`), Manifest Stores, Restoration Scripts, Verification Engine, SOP Runbook  
> **Relative Affected Paths:**  
> - `repo-secrets/09-antigravity-backup/readme.md`  
> - `repo-secrets/09-antigravity-backup/vault/projects-manifest.json`  
> - `repo-secrets/09-antigravity-backup/vault/pinned-projects.json`  
> - `repo-secrets/09-antigravity-backup/vault/settings-manifest.json`  
> - `repo-secrets/09-antigravity-backup/vault/plugins-and-skills.json`  
> - `repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.sh`  
> - `repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.ps1`  
> - `repo-secrets/09-antigravity-backup/scripts/verify-antigravity-backup.sh`  
> - `repo-secrets/09-antigravity-backup/temp/.gitkeep`  
> - `repo-secrets/readme.md`  

---

## 1. Technical Objective

Establish the comprehensive Antigravity backup vault in `repo-secrets/09-antigravity-backup/` containing portable JSON manifests for all 78 repositories, 23 pinned access items, user preferences, and installed skills/plugins using the portable `${WORKSPACE_ROOT}` substitution token; author the master Standard Operating Procedure (SOP) runbook; implement cross-platform zero-touch restoration scripts (`restore-antigravity-all.sh`, `restore-antigravity-all.ps1`, `verify-antigravity-backup.sh`); and register Section 09 in `repo-secrets/readme.md`.

---

## 2. Implementation Execution Plan

### Step 1: Scaffold Directory Topography
Create the standardized folder structure:
```bash
mkdir -p repo-secrets/09-antigravity-backup/vault
mkdir -p repo-secrets/09-antigravity-backup/scripts
mkdir -p repo-secrets/09-antigravity-backup/temp
touch repo-secrets/09-antigravity-backup/temp/.gitkeep
```

### Step 2: Author Master SOP Runbook (`readme.md`)
Create `repo-secrets/09-antigravity-backup/readme.md` documenting:
1. Executive overview of the Antigravity restoration system.
2. Complete directory breakdown of `vault/`, `scripts/`, and `temp/`.
3. Step-by-step SOP for migrating configurations to a new developer machine or container.
4. Detailed usage manual for both POSIX Bash and PowerShell runners.
5. Verification gates and security compliance rules.
6. Absolute enforcement of relative paths and positive booleans.

### Step 3: Author Portable JSON Manifests in `vault/`
Generate the four core manifest files adhering strictly to the JSON schemas defined in `02-component-spec.md`:

1. **`vault/projects-manifest.json`:**
   - Contains all 78 repositories classified into the 10 logical project sections.
   - Symmetrically maps UUIDs, preserving core UUIDs (`gitmap`, `scripts-fixer`, `wp-html-automate`, `antigravity-manager`, `letsmarknow`, `letsmarknow-ui`).
   - Uses `uriTemplate: "file://${WORKSPACE_ROOT}/<relativePath>"`.
   - Uses positive booleans: `isPortable: true`, `isSanitized: true`, `isEnabled: true`, `isWorkspaceOnly: false`, `v2Migrated: true`.
   - Records total count: `totalProjects: 78`.

2. **`vault/pinned-projects.json`:**
   - Contains all 23 priority repositories.
   - Symmetrically maps Tier 1 (7 Core Orchestrators) and Tier 2 (16 Ecosystem Services).
   - Preserves sequential priority order (1 to 23).
   - Uses positive booleans: `isPortable: true`, `isSyncEnabled: true`, `isPinned: true`, `isEnabled: true`.
   - Records total count: `totalPinned: 23`.

3. **`vault/settings-manifest.json`:**
   - Captures portable IDE user preferences, agent default model routing (`inherit`), and telemetry configurations.
   - Strictly sanitizes credentials (zero raw API tokens, session cookies, or secrets).
   - Uses positive booleans: `isPortable: true`, `isSanitized: true`, `isAutoExecutionAllowed: true`, `isContextCompactionEnabled: true`, `isLocalTelemetryActive: true`, `isCredentialMaskingStrict: true`.

4. **`vault/plugins-and-skills.json`:**
   - Documents active plugins (`gitmap-companion`, `project-manager-sync`) and skills (`gitmap`, `execute-parent-task-with-n-steps-v6`, `coding-guidelines`).
   - Uses positive booleans: `isPortable: true`, `isEnabled: true`, `isBuiltin: true`, `isGlobal: true`.

### Step 4: Implement POSIX Bash Restoration Runner (`scripts/restore-antigravity-all.sh`)
Author `repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.sh`:
- Include `set -euo pipefail`.
- Implement positive boolean flags: `isDryRun`, `shouldBackup`, `canOverwrite`, `isVerbose`, `isPortable`, `isValid`.
- Support CLI options: `--workspace-root <path>`, `--dry-run`, `--backup`, `--force`, `--verbose`, `--help`.
- Default workspace root to parent directory (`$PWD/..`) if not specified.
- Pre-flight backup: snapshot existing target config to `backup/pre-restore-<timestamp>/`.
- Ingestion routine: parse `projects-manifest.json`, substitute `${WORKSPACE_ROOT}`, write individual `<uuid>.json` files into `$HOME/.gemini/config/projects/`.
- Pinned projects routine: write `$HOME/.gemini/config/pinned_projects.json`.
- Settings merge routine: merge portable settings into `$HOME/.gemini/config/config.json`.
- Output summary metrics and return exit code 0.
- Make executable via `chmod +x repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.sh`.

### Step 5: Implement PowerShell Restoration Runner (`scripts/restore-antigravity-all.ps1`)
Author `repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.ps1`:
- Parameter block: `[CmdletBinding()]` supporting `-WorkspaceRoot`, `-DryRun`, `-Backup`, `-Force`, `-Verbose`.
- Symmetrical execution parity with Bash runner.
- Format Windows file URIs using forward slashes: `file://${WORKSPACE_ROOT}/<relativePath>`.
- Target destination: `$env:USERPROFILE\.gemini\config\projects\`.
- Atomic pre-restore backup to `$env:USERPROFILE\.gemini\config\backup\pre-restore-<timestamp>\`.
- Write individual project descriptors, `pinned_projects.json`, and merge settings.
- Robust error handling using positive boolean conditions.

### Step 6: Implement Verification Audit Script (`scripts/verify-antigravity-backup.sh`)
Author `repo-secrets/09-antigravity-backup/scripts/verify-antigravity-backup.sh`:
- Performs automated audit of verification gates VG-01 through VG-07.
- Checks:
  1. Manifest existence in `vault/`.
  2. JSON syntax validity via `python3 -m json.tool`.
  3. Absolute path scanning: zero machine drive letters or host path prefixes in manifests.
  4. Repository count parity: exactly 78 in `projects-manifest.json`.
  5. Pinned count parity: exactly 23 in `pinned-projects.json`.
  6. Positive boolean rule: all boolean properties use affirmative naming (`isEnabled`, `isPinned`, `isSanitized`, `isPortable`, `isValid`).
  7. Secret sanitization: zero raw API tokens or private keys.
- Returns exit code 0 on pass, 1 on failure.
- Make executable via `chmod +x repo-secrets/09-antigravity-backup/scripts/verify-antigravity-backup.sh`.

### Step 7: Update `repo-secrets/readme.md` Table of Contents
Safely edit `repo-secrets/readme.md` to index section `09-antigravity-backup/`:
- Add entry to the top-level directory table of contents.
- Document description: "Antigravity IDE project manifests, pinned configurations, cross-platform restoration automation, and migration runbook."
- Ensure strictly relative link: `09-antigravity-backup/readme.md`.

---

## 3. Verification & Validation Commands

Worker 02 must execute the following validation commands upon completing authoring:

```bash
# 1. Run the automated backup verification audit
bash repo-secrets/09-antigravity-backup/scripts/verify-antigravity-backup.sh

# 2. Execute dry-run restoration on POSIX Bash runner
bash repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.sh --dry-run --verbose

# 3. Verify total project count in manifest
python3 -c "
import json
with open('repo-secrets/09-antigravity-backup/vault/projects-manifest.json') as f:
    d = json.load(f)
assert d.get('totalProjects') == 78, f'Expected 78 projects, found {d.get(\"totalProjects\")}'
assert len(d.get('projects', [])) == 78, f'Expected 78 project entries, found {len(d.get(\"projects\", []))}'
print('Projects manifest validated: 78 / 78 repositories')
"

# 4. Verify pinned count in manifest
python3 -c "
import json
with open('repo-secrets/09-antigravity-backup/vault/pinned-projects.json') as f:
    d = json.load(f)
assert d.get('totalPinned') == 23, f'Expected 23 pinned projects, found {d.get(\"totalPinned\")}'
assert len(d.get('projects', [])) == 23, f'Expected 23 pinned entries, found {len(d.get(\"projects\", []))}'
print('Pinned manifest validated: 23 / 23 repositories')
"

# 5. Audit for accidental absolute paths in repo-secrets/09-antigravity-backup/
test $(grep -rn "\${WORKSPACE_ROOT}" repo-secrets/09-antigravity-backup/vault/ | wc -l) -gt 0
! grep -E -rn "(/home/|[A-Za-z]:\\\\|[A-Za-z]:/)" repo-secrets/09-antigravity-backup/vault/
```

---

## 4. Acceptance Criteria Checklist

- [ ] **AC-01**: Directory topography `repo-secrets/09-antigravity-backup/` created with `readme.md`, `vault/`, `scripts/`, `temp/`.
- [ ] **AC-02**: All 4 manifests (`projects-manifest.json`, `pinned-projects.json`, `settings-manifest.json`, `plugins-and-skills.json`) exist in `vault/` and validate against JSON schemas.
- [ ] **AC-03**: Exactly 78 repositories recorded in `projects-manifest.json` across 10 sections with portable `${WORKSPACE_ROOT}` token.
- [ ] **AC-04**: Exactly 23 pinned repositories recorded in `pinned-projects.json` (7 Tier 1 Core + 16 Tier 2 Ecosystem).
- [ ] **AC-05**: Master SOP runbook authored in `repo-secrets/09-antigravity-backup/readme.md` with complete migration workflows.
- [ ] **AC-06**: `restore-antigravity-all.sh` and `restore-antigravity-all.ps1` implemented with dry-run, backup, and restore flags.
- [ ] **AC-07**: `verify-antigravity-backup.sh` implemented and passes all verification gates (VG-01 through VG-07).
- [ ] **AC-08**: Section `09-antigravity-backup/` registered in `repo-secrets/readme.md`.
- [ ] **AC-09**: Zero absolute host paths or drive letters in tracked manifests and documentation.
- [ ] **AC-10**: Positive boolean conventions applied 100% across all schemas, scripts, and documentation.
