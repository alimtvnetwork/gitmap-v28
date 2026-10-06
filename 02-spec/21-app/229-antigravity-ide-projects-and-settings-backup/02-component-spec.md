# Component Specification: Antigravity IDE Projects & Settings Migration and Backup Engine

## Metadata
- **Specification ID**: `SPEC-229-COMPONENT`
- **Parent Task**: `229-antigravity-ide-projects-and-settings-backup`
- **Target Subsystem**: `repo-secrets/09-antigravity-backup/` & Antigravity IDE Configuration Store
- **Author**: Spec Author 02
- **Status**: APPROVED
- **Review Date**: 2026-10-06

---

## 1. Executive Summary & Scope

This component specification establishes the migration blueprint, data schemas, automation script requirements, and verification protocols for backing up and restoring Google Antigravity IDE projects, pinned project lists, workspace settings, and agent skills across machines and runtime instances.

The output artifacts specified herein reside within the version-controlled `repo-secrets/09-antigravity-backup/` directory, providing an immutable, sanitized vault and self-healing restoration toolchain for polyglot environments spanning Windows PowerShell and POSIX Shell (Linux / macOS / WSL).

---

## 2. User Request & Non-Negotiable Requirements

### 2.1 Verbatim Request
```text
please all repos following the antigravity ide here in the project sections along with pin projects, and also list out your work and steps into the repo-secrets folder for restoring projects and settings from one antigravity to another create folder and also can use gitmap required to
```

### 2.2 Core Mandates
1. **Repository Inventory Parity**: Catalog all 78 detected Git repositories into the Antigravity project ecosystem, mapping default branches, relative workspace roots, and project identifiers.
2. **Pinned Projects Tiering**: Define and persist Tier 1 (core infrastructure and orchestrators) and Tier 2 (active development services) project pins.
3. **Restoration SOP & Automation in `repo-secrets`**: Architect `repo-secrets/09-antigravity-backup/` containing full runbooks, JSON vault manifests, and zero-dependency restoration scripts.
4. **Strict Path Relativity**: All documentation, scripts, and configuration templates within git tracking MUST strictly reference relative Git paths (`repo-secrets/...`, `02-spec/...`, `.ai-memory/...`). No absolute paths or local filesystem URIs are permitted in tracked files.
5. **Positive Boolean Hygiene**: All state variables, parameter switches, and JSON properties must adhere to affirmative positive boolean naming (`isEnabled`, `isClean`, `hasBackup`, `shouldVerify`).

---

## 3. Directory Topography: `repo-secrets/09-antigravity-backup/`

The backup bundle is organized into three segregated tiers: documentation, serialized vault manifests, and cross-platform automation scripts.

```text
repo-secrets/09-antigravity-backup/
├── readme.md                           # Master Restoration Standard Operating Procedure (SOP)
├── vault/                              # Serialized, sanitized JSON backup manifests
│   ├── projects-manifest.json          # 78 repositories with relative roots & branch metadata
│   ├── pinned-projects.json            # Tier 1 & Tier 2 pinned project registries
│   ├── settings-manifest.json          # Sanitized IDE preferences & model configurations
│   └── plugins-and-skills.json         # Workspace skills & extension capability registry
└── scripts/                            # Cross-platform zero-touch restoration scripts
    ├── restore-antigravity-all.ps1     # Master PowerShell orchestrator (Windows)
    ├── restore-antigravity-all.sh      # Master POSIX Shell orchestrator (Linux/macOS/WSL)
    ├── restore-antigravity-projects.ps1 # Individual project JSON generator
    ├── restore-antigravity-settings.ps1 # Settings injector & schema normalizer
    ├── restore-antigravity-pins.ps1    # Pinned projects store synchronizer
    └── verify-antigravity-backup.ps1   # Post-restoration verification & health check
```

---

## 4. Vault Manifest Specifications

All vault manifests are encoded in UTF-8 with standard 2-space JSON formatting and validated against affirmative boolean schemas.

### 4.1 Projects Manifest: `vault/projects-manifest.json`
Catalogs all 78 active repositories discovered across the developer workspace. Replaces machine-specific root prefixes with the portable `${WORKSPACE_ROOT}` token.

#### Schema Definition
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "AntigravityProjectsManifest",
  "type": "object",
  "required": ["version", "generatedAt", "totalCount", "isPortable", "projects"],
  "properties": {
    "version": { "type": "string" },
    "generatedAt": { "type": "string", "format": "date-time" },
    "totalCount": { "type": "integer", "minimum": 1 },
    "isPortable": { "type": "boolean" },
    "projects": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "name", "relativeWorkspacePath", "defaultBranch", "isActive", "hasGitFolder"],
        "properties": {
          "id": { "type": "string" },
          "name": { "type": "string" },
          "relativeWorkspacePath": { "type": "string" },
          "defaultBranch": { "type": "string" },
          "isActive": { "type": "boolean" },
          "hasGitFolder": { "type": "boolean" },
          "category": { "type": "string" }
        }
      }
    }
  }
}
```

#### Field Specifications
- `id`: Deterministic alphanumeric identifier (kebab-case slug matching repository name).
- `name`: Human-readable display label matching the repository folder.
- `relativeWorkspacePath`: Portable path relative to workspace parent (e.g., `gitmap`, `movie-cli-v8`, `notes-app-v3`).
- `defaultBranch`: Primary tracking branch (`main` or `master`).
- `isActive`: Boolean indicating active inclusion in IDE indexing.
- `hasGitFolder`: Boolean confirming presence of valid `.git` root.

---

### 4.2 Pinned Projects Manifest: `vault/pinned-projects.json`
Maintains the priority rankings and quick-access definitions for critical repositories.

#### Schema Definition
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "AntigravityPinnedProjects",
  "type": "object",
  "required": ["version", "updatedAt", "isSyncEnabled", "pinnedProjects"],
  "properties": {
    "version": { "type": "string" },
    "updatedAt": { "type": "string", "format": "date-time" },
    "isSyncEnabled": { "type": "boolean" },
    "pinnedProjects": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "name", "tier", "priorityOrder", "isPinned", "relativeWorkspacePath"],
        "properties": {
          "id": { "type": "string" },
          "name": { "type": "string" },
          "tier": { "type": "integer", "enum": [1, 2] },
          "priorityOrder": { "type": "integer" },
          "isPinned": { "type": "boolean" },
          "relativeWorkspacePath": { "type": "string" }
        }
      }
    }
  }
}
```

#### Tier Allocation Standard
- **Tier 1 (Core Infrastructure & Orchestration)**:
  - `gitmap`: Central autonomous companion, multi-agent orchestrator, Split-DB engine.
  - `coding-guidelines`: Universal architectural specifications and coding standards source repo.
- **Tier 2 (Active Development Repositories)**:
  - `movie-cli-v8`: High-velocity CLI application and testbed.
  - `notes-app-v3`: Core application workload.
  - `portfolio-v2`: Web delivery and client portal.

---

### 4.3 Settings Manifest: `vault/settings-manifest.json`
Contains sanitized Antigravity IDE global and instance configurations, eliminating host-dependent file paths and credentials while preserving developer workflow ergonomics.

#### Structure & Sanitization Rules
- **Editor Configurations**: Tab sizes, format-on-save, whitespace trimming, font ligatures.
- **AI Agent Settings**: Telemetry mode, context token budget, auto-compaction triggers.
- **Strictly Excluded (Sanitized)**:
  - API keys (`geminiApiKey`, `anthropicApiKey`, `openaiApiKey`).
  - Session tokens and OAuth refresh cookies.
  - Machine-specific hardware identifiers and local IP addresses.

```json
{
  "version": "1.0.0",
  "sanitizedAt": "2026-10-06T12:00:00Z",
  "isSanitized": true,
  "hasSecretsExcluded": true,
  "antigravitySettings": {
    "telemetry": {
      "isEnabled": false,
      "isCrashReportingActive": false
    },
    "agent": {
      "maxConcurrentSubagents": 4,
      "isAutoCompactionEnabled": true,
      "contextWindowTokens": 200000,
      "isStrictRelativePathEnforced": true
    },
    "editor": {
      "formatOnSave": true,
      "insertSpaces": true,
      "tabSize": 2,
      "trimTrailingWhitespace": true,
      "renderWhitespace": "selection"
    }
  }
}
```

---

### 4.4 Plugins & Skills Manifest: `vault/plugins-and-skills.json`
Indexes both built-in IDE extensions and portable workspace skills under `.agents/skills/`.

#### Schema Structure
- Catalog of all 70+ verified workspace skills (e.g., `gitmap`, `gitmap-agent-orchestrator`, `execute-parent-task-with-n-steps-v6`).
- Portable references using relative path mappings.
- Status flags confirming validation state (`isValid`, `hasSkillMarkdown`, `isEnabled`).

---

## 5. Script Architecture & Component Design

The restoration scripts provide zero-touch deployment across target environments without requiring prior installation of external dependencies (pure PowerShell 5.1+/Core and POSIX Bash).

### 5.1 Master PowerShell Orchestrator: `restore-antigravity-all.ps1`
- **Purpose**: Unified entry point for Windows workstations.
- **Execution Workflow**:
  1. Parse command-line parameters (`-WorkspaceRoot`, `-InstanceTarget`, `-IsDryRun`, `-ShouldForceOverwrite`).
  2. Perform pre-flight environment discovery (detect user home directory, Antigravity instance paths, GitMap CLI).
  3. Invoke `restore-antigravity-projects.ps1` to populate project definitions.
  4. Invoke `restore-antigravity-pins.ps1` to synchronize pinned project lists.
  5. Invoke `restore-antigravity-settings.ps1` to apply sanitized configuration settings.
  6. Execute `verify-antigravity-backup.ps1` to validate restoration integrity against scorecard gates.
- **Safety Flags**:
  - `-IsDryRun`: Preview filesystem actions without modifying disk state.
  - `-ShouldCreateBackup`: Create timestamped backup of existing settings prior to mutation.

### 5.2 Master Bash Orchestrator: `restore-antigravity-all.sh`
- **Purpose**: Unified entry point for Linux, macOS, and WSL environments.
- **Parity Guarantee**: Mirror all parameter flags, path normalization routines, and exit codes of the PowerShell orchestrator.
- **POSIX Safety**: `set -euo pipefail`, trap signals on error, clear descriptive terminal output.

### 5.3 Projects Restorer: `restore-antigravity-projects.ps1`
- **Mechanism**:
  - Reads `vault/projects-manifest.json`.
  - Determines destination project directory:
    - Default: `~/.gemini/config/projects/`
    - Instance-aware: `~/.antigravity_tools/instances/<instance-id>/home/.gemini/config/projects/`
  - Iterates through project entries. For each repository:
    - Generates `<project-id>.json`.
    - Constructs valid `FolderURI` with OS-appropriate format dynamically using target environment root.
    - Writes formatted JSON file with atomic write semantics (`[System.IO.File]::WriteAllText`).

### 5.4 Pins Restorer: `restore-antigravity-pins.ps1`
- **Mechanism**:
  - Ingests `vault/pinned-projects.json`.
  - Targets destination store: `~/.gemini/config/pinned_projects.json`.
  - Re-anchors relative workspace paths against the target machine's active workspace root.
  - Merges or replaces pin registry based on `-ShouldMergeExisting` flag.
  - Optionally calls `gitmap agy pins add` if GitMap binary is discovered on `$env:PATH`.

### 5.5 Settings Restorer: `restore-antigravity-settings.ps1`
- **Mechanism**:
  - Ingests `vault/settings-manifest.json`.
  - Targets `~/.gemini/antigravity/config.json`.
  - Safely deep-merges restored settings into destination without overwriting locally configured authentication tokens.

### 5.6 Backup Verifier: `verify-antigravity-backup.ps1`
- **Mechanism**:
  - Executes validation passes across project files, pinned project stores, and settings JSON files.
  - Calculates checksums and ensures count parity (78 projects restored).
  - Evaluates all Verification Scorecard Gates (VG-01 through VG-07).
  - Emits summary scorecard and returns standard process exit codes (0 = Success, 1 = Gate Failure).

---

## 6. Architectural Solutions for Gaps GAP-01 through GAP-09

### 6.1 GAP-01: Cross-Platform Path Normalization
- **Defect**: Windows environments use drive letters and backslashes (e.g., repository path delimiters), while Linux/WSL environments use forward slashes. Ingesting Windows backslashes into POSIX Antigravity instances causes broken project references and unresolvable workspace links.
- **Architectural Solution**:
  - Store all paths in vault manifests using portable relative POSIX notation (forward slashes, no drive letters).
  - Restoration scripts detect the runtime OS (`$IsWindows` / `uname -s`) and normalize paths dynamically:
    ```powershell
    function Resolve-NormalizedPath {
        param([string]$BaseDir, [string]$RelativePath)
        $Combined = Join-Path -Path $BaseDir -ChildPath $RelativePath
        return [System.IO.Path]::GetFullPath($Combined)
    }
    ```

### 6.2 GAP-02: Relative URI Remapping
- **Defect**: Antigravity project JSON files store workspace links as folder URIs. Hardcoding absolute, machine-specific URIs prevents cross-machine portability.
- **Architectural Solution**:
  - The manifests store `relativeWorkspacePath`.
  - The restorer script calculates the target machine's base workspace root and constructs the folder URI dynamically using standard RFC 8089 URI conventions:
    ```powershell
    function Convert-ToFolderUri {
        param([string]$AbsolutePath)
        $CleanPath = $AbsolutePath.Replace('\', '/')
        if (-not $CleanPath.StartsWith('/')) {
            $CleanPath = '/' + $CleanPath
        }
        return [System.Uri]::new($CleanPath).AbsoluteUri
    }
    ```

### 6.3 GAP-03: Multi-Instance Handling
- **Defect**: Antigravity supports isolated runtime instances located under `~/.antigravity_tools/instances/<instance-id>/home/` in addition to the global user profile `~/.gemini/`. A restore script targeting only the global profile fails to update the active working instance.
- **Architectural Solution**:
  - The restoration engine implements an instance resolution cascade:
    1. Check for explicit `-InstanceTarget <instance-id>` parameter.
    2. Auto-detect active instances by inspecting `~/.antigravity_tools/instances/`.
    3. Update both the active instance profile AND the global `~/.gemini/` fallback profile when `-IsGlobalSyncEnabled` is true.

### 6.4 GAP-04: Settings Schema Alignment
- **Defect**: Antigravity maintains configuration across multiple tiers: `~/.gemini/antigravity/config.json`, IDE settings `settings.json`, and internal state DBs. Direct file replacement can overwrite newer configuration versions or introduce deprecated keys.
- **Architectural Solution**:
  - Implement non-destructive deep merging.
  - The restorer parses existing JSON (if present), overlays sanitized manifest values, preserves local keys not defined in the manifest, and writes the merged result back atomically.

### 6.5 GAP-05: SQLite WAL & Lock Integrity
- **Defect**: Antigravity conversation and workspace state databases (`conversation_summaries.db`, `state.vscdb`) operate in SQLite WAL (Write-Ahead Logging) mode. Copying `.db` files while active processes hold open locks or while uncommitted WAL files (`-wal`, `-shm`) exist causes corruption.
- **Architectural Solution**:
  - The backup engine enforces pre-flight lock checks.
  - In PowerShell and Shell scripts, require Antigravity to be closed or issue a WAL checkpoint prior to any database copy operations.
  - In GitMap operations, utilize `PRAGMA wal_checkpoint(TRUNCATE)` before exporting database snapshots.

### 6.6 GAP-06: Pinned Projects Dual-Store Synchronization
- **Defect**: Pinned projects are tracked in both `.gemini/config/pinned_projects.json` and internal GitMap database tables (`agy_pin_projects`). Updating one without the other leads to split-brain state.
- **Architectural Solution**:
  - The `restore-antigravity-pins.ps1` script writes the canonical `pinned_projects.json` file first.
  - Upon successful file write, if the `gitmap` executable is available, the script executes `gitmap agy pins reconcile` to synchronize the internal database state with the JSON source of truth.

### 6.7 GAP-07: Plugins & Built-in Skills Catalog Portability
- **Defect**: Skills exist in two locations: built-in IDE skills (located in instance application data directories) and repository-specific workspace skills (under `.agents/skills/`). Moving settings across machines breaks absolute skill symlinks.
- **Architectural Solution**:
  - The `plugins-and-skills.json` manifest records skill names and relative configurations.
  - Restoration scripts only re-link workspace skills relative to the repository workspace root, avoiding tampering with machine-managed built-in IDE binary bundles.

### 6.8 GAP-08: Cross-Platform Shell Script Parity
- **Defect**: Shell scripts authored for Windows PowerShell frequently utilize non-portable cmdlets (`Get-ItemPropertyValue`, Windows-specific path parsing), while POSIX scripts often fail when executed in minimal environments lacking Bash 5.
- **Architectural Solution**:
  - Enforce complete behavioral parity between `restore-antigravity-all.ps1` and `restore-antigravity-all.sh`.
  - Use standard POSIX-compliant constructs in `.sh` (`/bin/sh` or `/bin/bash` with minimal built-ins) and standard PowerShell Core / Windows PowerShell 5.1 compatible syntax in `.ps1`.

### 6.9 GAP-09: Secret & Credential Sanitization
- **Defect**: IDE settings files may inadvertently capture developer API keys, session tokens, or personal access tokens. Committing these into `repo-secrets/` poses a severe security risk.
- **Architectural Solution**:
  - Implement a mandatory redaction filter during manifest generation.
  - All keys matching patterns `*Token*`, `*Key*`, `*Secret*`, `*Password*`, or `*Auth*` are stripped or replaced with redacted placeholders.
  - Verification gate VG-06 strictly audits the vault manifests to ensure zero sensitive string patterns exist.

---

## 7. Verification Scorecard Gates (VG-01 to VG-07)

Every backup and restoration cycle must pass the seven verification gates before being certified for production use:

| Gate ID | Gate Name | Evaluation Criteria | Enforcement Mechanism |
|:---|:---|:---|:---|
| **VG-01** | Vault Structure & Manifest Validation | All 4 JSON files in `vault/` exist and match their JSON schemas | `verify-antigravity-backup.ps1` schema validation check |
| **VG-02** | Repository Count Parity | `projects-manifest.json` contains exactly 78 verified repositories | Count assertion (`$projects.Count -eq 78`) |
| **VG-03** | Path Relativity & Portability | 100% of paths in manifests and docs are relative (zero drive letters or `/home` roots) | Regex audit for absolute path patterns |
| **VG-04** | Restoration Script Parity | Both `.ps1` and `.sh` scripts execute without syntax errors and support `-IsDryRun` | Script static analysis & dry-run execution test |
| **VG-05** | Atomic Non-Destructive Ingestion | Target files are merged non-destructively; existing files backed up if flag set | File integrity and pre-restore backup check |
| **VG-06** | Secret Sanitization Gate | Zero API keys, tokens, passwords, or personal credentials present in manifests | Keyword & pattern scanning across vault files |
| **VG-07** | Git Hygiene & Relative Path Gate | All links in documentation adhere strictly to repository-relative format | Documentation markdown link inspection |

---

## 8. Master Restoration Runbook Summary (for `readme.md`)

The `repo-secrets/09-antigravity-backup/readme.md` must provide a comprehensive, operator-ready manual structured as follows:
1. **Overview & Architecture**: Purpose of the backup vault and component layout.
2. **Prerequisites**: Supported operating systems (Windows 10/11, Ubuntu 20.04+, macOS 12+), PowerShell 5.1+/Core or Bash, GitMap CLI (recommended).
3. **Quick Start Workflows**:
   - One-touch Windows restore: `powershell -ExecutionPolicy Bypass -File scripts/restore-antigravity-all.ps1 -WorkspaceRoot <path>`
   - One-touch POSIX restore: `bash scripts/restore-antigravity-all.sh --workspace-root <path>`
4. **Surgical Step-by-Step Procedures**:
   - Restoring project registries only (`restore-antigravity-projects.ps1`).
   - Restoring pinned projects only (`restore-antigravity-pins.ps1`).
   - Restoring sanitized IDE settings only (`restore-antigravity-settings.ps1`).
5. **Post-Restoration Verification**: Running `verify-antigravity-backup.ps1` to confirm scorecard gate pass.
6. **Troubleshooting Guide**: Resolving locked database files, path mismatch issues, and multi-instance path redirections.
