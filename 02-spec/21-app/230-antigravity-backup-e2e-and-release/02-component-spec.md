# Component Specification: Minor Version Bump & Release Ceremony Engine

## Metadata
- **Specification ID**: `SPEC-230-COMPONENT`
- **Parent Task**: `230-antigravity-backup-e2e-and-release`
- **Target Subsystem**: Version Manifests, Release Orchestration, `.ai-memory/release/`, and Git Release Lifecycle
- **Author**: Spec Author 02
- **Status**: APPROVED
- **Review Date**: 2026-10-06
- **Compliance Rules**:
  - `isRelativePathsOnly`: true
  - `isPositiveBooleanConventions`: true
  - `isZeroAbsolutePaths`: true
  - `isRuleZeroMinorBump`: true

---

## 1. Executive Summary & Architectural Scope

This component specification formalizes the autonomous minor version bump architecture, version manifest synchronization across polyglot configurations, canonical release notes generation, and the mandatory 5-step release branching lifecycle and tagging ceremony for GitMap.

Following the implementation and bidirectional end-to-end verification of the Google Antigravity IDE backup, restoration, registration, and pin management toolchain in `repo-secrets/09-antigravity-backup/`, GitMap undergoes an atomic minor version increment from `v6.493.0` to `v6.494.0`.

All artifacts, manifests, documentation references, and release notes specified herein MUST adhere strictly to relative Git paths (`02-spec/...`, `.ai-memory/...`, `repo-secrets/...`, `03-ai-scripts/...`, `cli/...`). The inclusion of absolute filesystem paths, operating system drive roots (such as Windows drive prefixes), Unix user root prefixes, or local file URI schemes is strictly forbidden across all version-controlled files and generated release documentation.

---

## 2. Non-Negotiable Engineering Directives

### 2.1 Verbatim Ingested Directives
```text
1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first.
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String.
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well.
4. Perform a minor version bump.
5. Make a release with the release notes.
```

### 2.2 Core Mandates
1. **Rule 0 Minor Arithmetic**: Transition version strictly from `6.493.0` to `6.494.0`. Under Rule 0, minor is the canonical default tier; the minor component increments by 1 (`493` -> `494`), and the patch component resets strictly to `0`.
2. **Atomic Manifest Synchronization**: Synchronize all 8 repository manifest locations in a single unified operation without manual edits or omitted files.
3. **100% Relative Git Paths**: Enforce relative repository paths across all release notes, changelog entries, and documentation badges.
4. **Positive Boolean Hygiene**: All state indicators, validation predicates, and operational flags must employ affirmative naming (`isCleanState`, `isManifestSynced`, `isValidSemVer`, `hasReleaseNotes`, `canProceed`).
5. **Deterministic 5-Step Branching Lifecycle**: Execute release preparation on an isolated release branch (`release/v6.494.0`), perform manifest updates, commit, tag, and merge back to `main`.

---

## 3. Minor Version Bump Architecture (`v6.493.0` -> `v6.494.0`)

### 3.1 Single Canonical Source of Truth
The canonical source of truth for version state across the entire repository is `version.json` located at the repository root. All auxiliary manifests, language-specific project descriptors, CLI build constants, and markdown badges derive their state from this file.

```text
                           ┌──────────────────────────────┐
                           │   Canonical Source of Truth  │
                           │         version.json         │
                           │          (v6.494.0)          │
                           └──────────────┬───────────────┘
                                          │
       ┌──────────────────┬───────────────┼───────────────┬──────────────────┐
       │                  │               │               │                  │
       ▼                  ▼               ▼               ▼                  ▼
┌──────────────┐   ┌──────────────┐ ┌───────────┐ ┌──────────────┐ ┌───────────────────┐
│ readme.md    │   │changelog.md  │ │package.json│ │cli/constants/│ │.gitmap/release/   │
│ what-to-read │   │spec98-change │ │template.json││constants.go │ │latest.json        │
│ (pins/badges)│   │(history log) │ │(manifests)│ │(var Version) │ │(branch/tag/ver)   │
└──────────────┘   └──────────────┘ └───────────┘ └──────────────┘ └───────────────────┘
```

### 3.2 SemVer 2.0.0 Transition Arithmetic
The version transition satisfies the SemVer 2.0.0 specification augmented by repository Rule 0:

$$\text{SemVer}_{\text{next}} = \text{MAJOR} \,.\, (\text{MINOR} + 1) \,.\, 0$$

- **Current Canonical Version**: `6.493.0`
  - `MAJOR`: 6
  - `MINOR`: 493
  - `PATCH`: 0
- **Computed Next Version**: `6.494.0`
  - `MAJOR`: 6 (constant)
  - `MINOR`: 494 (incremented)
  - `PATCH`: 0 (reset)
- **Annotated Git Tag**: `v6.494.0`
- **Dedicated Release Branch**: `release/v6.494.0`

### 3.3 Automated Bump Engine: `03-ai-scripts/37-bump-version.py`
The version update is executed autonomously by `03-ai-scripts/37-bump-version.py`. The engine executes the following core workflow:
1. `read_canonical_version()`: Ingests `version.json` and parses `Version: "6.493.0"`.
2. `calculate_next_version(current_ver, tier="minor")`: Emits `6.494.0`.
3. `update_version_json(next_ver, today_str)`: Updates `version.json` with new version and ISO date.
4. `update_package_json(next_ver)`: Updates root `package.json` if present.
5. `update_template_version(next_ver)`: Synchronizes `prompt-version.template.json`.
6. `update_latest_json(next_ver)`: Updates `.gitmap/release/latest.json`.
7. `update_readme_pins(current_ver, next_ver)`: Updates badges and install snippets in `readme.md`.
8. `update_what-to-read_pins(current_ver, next_ver)`: Updates badges and pins in `what-to-read.md`.
9. `update_constants_go(next_ver)`: Replaces `var Version` in `cli/constants/constants.go`.
10. `create_release_notes(next_ver, scope)`: Authors initial release note scaffold in `.ai-memory/release/`.
11. `update_changelogs(next_ver, scope, today_str)`: Prepends changelog entries to `changelog.md` and `02-spec/19-main-worker-service/98-changelog.md`.
12. `run_repo_sync_if_available()`: Triggers `go generate ./...` in `cli/` and `npm run sync` if defined.

---

## 4. Version Manifests Synchronization Topology

Every release requires complete synchronization across 8 target manifest locations:

| Target File (Relative Git Path) | Modification Scope | Validation Criteria |
| :--- | :--- | :--- |
| `version.json` | `"Version": "6.494.0"`, `"releaseDate": "YYYY-MM-DD"` | Valid JSON, matches regex `^6\.494\.0$` |
| `package.json` | `"version": "6.494.0"` | Valid JSON, top-level version field updated |
| `prompt-version.template.json` | `"version": "6.494.0"` | Valid JSON template, matches version |
| `.gitmap/release/latest.json` | `"version": "6.494.0"`, `"tag": "v6.494.0"`, `"branch": "release/v6.494.0"` | Valid JSON, branch/tag synchronization |
| `readme.md` | Pinned badges, quick-install URLs, header versions | Zero stale `v6.493.0` strings in pins |
| `what-to-read.md` | Header badges, install commands, guide version references | Zero stale `v6.493.0` strings in pins |
| `changelog.md` | Prepend `## [v6.494.0] - YYYY-MM-DD` entry block | Clean Markdown header, categorized bullets |
| `cli/constants/constants.go` | `var Version = "6.494.0"` | Valid Go syntax, compiles without errors |

### 4.1 `version.json` Manifest Schema
```json
{
  "_purpose": "version.json is the single canonical source of truth for repository versioning...",
  "Version": "6.494.0",
  "Title": "Coding Guidelines",
  "RepoSlug": "coding-guidelines-v24",
  "RepoUrl": "https://github.com/alimtvnetwork/coding-guidelines-v24",
  "releaseDate": "2026-10-06"
}
```

### 4.2 `readme.md` & `what-to-read.md` Version Pinning Rules
All references matching the following patterns MUST be replaced:
1. `**Pinned version: v6.493.0**` $\rightarrow$ `**Pinned version: v6.494.0**`
2. `#### Pinned Version Install (v6.493.0)` $\rightarrow$ `#### Pinned Version Install (v6.494.0)`
3. Raw GitHub release URL endpoints:
   - PowerShell: `https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.494.0/install.ps1`
   - POSIX Shell: `https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.494.0/install.sh`

### 4.3 `changelog.md` Entry Format
```markdown
# Changelog

## [v6.494.0] - 2026-10-06

### Added
- Antigravity IDE backup, restoration, registration, and pin management bidirectional E2E test verification (`repo-secrets/09-antigravity-backup/`).
- Automated roundtrip audit logging verification across PowerShell and POSIX restoration scripts.
- Minor version bump synchronization across version manifests, CLI constants, and release artifacts.

### Changed
- Promoted GitMap release version from `v6.493.0` to `v6.494.0`.

---
```

---

## 5. Canonical Release Notes Specification

### 5.1 Location & Naming Standard
- **Relative Git Path**: `.ai-memory/release/release-notes-v6.494.0.md`
- **Encoding**: UTF-8 without BOM, LF line endings
- **Strict Prohibition**: Under NO circumstance shall this file contain absolute filesystem paths, Windows drive identifiers, local machine user names, or local file URI hyperlinks.

### 5.2 Release Notes Schema Definition
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "GitMapReleaseNotesSpecification",
  "type": "object",
  "required": [
    "version",
    "releaseDate",
    "hasQuickInstallOneLiners",
    "hasRelativePathsOnly",
    "sections"
  ],
  "properties": {
    "version": { "type": "string", "pattern": "^v\\d+\\.\\d+\\.\\d+$" },
    "releaseDate": { "type": "string", "format": "date" },
    "hasQuickInstallOneLiners": { "type": "boolean" },
    "hasRelativePathsOnly": { "type": "boolean" },
    "sections": {
      "type": "object",
      "required": ["quickInstall", "whatsChanged", "relativeLedger", "verificationMatrix"],
      "properties": {
        "quickInstall": { "type": "object" },
        "whatsChanged": { "type": "array" },
        "relativeLedger": { "type": "array" },
        "verificationMatrix": { "type": "object" }
      }
    }
  }
}
```

### 5.3 Canonical Content Template: `.ai-memory/release/release-notes-v6.494.0.md`
```markdown
# GitMap v6.494.0

## Quick Install v6.494.0

### Windows (PowerShell)
```powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.494.0/install.ps1 -OutFile install.ps1; .\install.ps1 -TargetDir ".ai-memory/prompts" -Version "v6.494.0"
```

### Unix / Linux / macOS (Bash)
```bash
curl -sL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.494.0/install.sh | bash -s -- ".ai-memory/prompts" "v6.494.0"
```

---

## What's Changed in v6.494.0

### Added
- **Antigravity IDE Backup E2E Suite**: Full bidirectional validation of backup, export, ingestion, restoration, and verification scripts in `repo-secrets/09-antigravity-backup/scripts/`.
- **Bidirectional Roundtrip Verification**: Validated that exported projects manifest (`repo-secrets/09-antigravity-backup/vault/projects-manifest.json`) and pinned projects manifest (`repo-secrets/09-antigravity-backup/vault/pinned-projects.json`) reconstruct active workspace configurations identically.
- **Audit Logging at Each Step**: Added structured telemetry output across `backup-antigravity-state.ps1`, `restore-antigravity-all.ps1`, and `verify-antigravity-backup.ps1`.
- **Minor Version Bump v6.494.0**: Synchronized `version.json`, root `readme.md`, `changelog.md`, `cli/constants/constants.go`, and `.gitmap/release/latest.json`.

### Verified
- Bidirectional roundtrip export and restoration verified.
- Pre-flight quality gates and manifest parity verified.
- 100% relative Git paths compliance verified across all release artifacts.

---

## Relative Artifacts & Specifications Ledger

| Artifact Type | Relative Git Path | Purpose |
| :--- | :--- | :--- |
| Component Spec | `02-spec/21-app/230-antigravity-backup-e2e-and-release/02-component-spec.md` | Minor bump & release ceremony specification |
| Architecture Spec | `02-spec/21-app/230-antigravity-backup-e2e-and-release/01-architecture-spec.md` | E2E architecture & test topology |
| Subtask Plan | `.ai-memory/plans/subtasks/230-antigravity-backup-e2e-and-release/02-minor-bump-and-release-ceremony.md` | Operational release instructions |
| Canonical Version | `version.json` | Root single source of truth |
| Release Notes | `.ai-memory/release/release-notes-v6.494.0.md` | Published release changelog |
| Backup SOP | `repo-secrets/09-antigravity-backup/readme.md` | Restoration Standard Operating Procedure |
| Projects Vault | `repo-secrets/09-antigravity-backup/vault/projects-manifest.json` | 78 cataloged repositories manifest |
| Pinned Vault | `repo-secrets/09-antigravity-backup/vault/pinned-projects.json` | Tier 1 & Tier 2 pinned IDE projects |
| Restoration Script | `repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.ps1` | Zero-dependency PowerShell restoration |
| Verification Script | `repo-secrets/09-antigravity-backup/scripts/verify-antigravity-backup.ps1` | Post-restoration integrity test |

---

## Verification & Integrity Checksums
- `isCleanWorkingTree`: true
- `hasRelativePathsOnly`: true
- `isManifestSynced`: true
- `isValidSemVer`: true
```

---

## 6. Mandatory 5-Step Release Branching Lifecycle & Tagging Ceremony

The Git release ceremony executes through five distinct, sequential phases to guarantee zero regression on the default branch and total traceability:

```text
  [Active Branch (main)]
            │
            ▼
 1. Check Working Tree Cleanliness (git status --porcelain)
            │
            ▼
 2. [STEP 1] Dedicated Release Branch Creation:
             git checkout -b release/v6.494.0
            │
            ▼
 3. [STEP 2] Execute Version Bump & Manifest Synchronization:
             python 03-ai-scripts/37-bump-version.py -t minor --scope "Antigravity IDE Backup E2E Validation and Release"
            │
            ▼
 4. [STEP 3] Atomic Staging & Commit on Release Branch:
             git add version.json package.json readme.md changelog.md what-to-read.md \
                     cli/constants/constants.go .gitmap/release/latest.json \
                     .ai-memory/release/release-notes-v6.494.0.md
             git commit -m "release: v6.494.0 Antigravity IDE backup E2E validation and release"
            │
            ▼
 5. [STEP 4] Cryptographic Annotated Tag Creation:
             git tag -a v6.494.0 -m "Release v6.494.0: Antigravity IDE Backup E2E & Release Ceremony"
            │
            ▼
 6. [STEP 5] Merge to Main, Push Upstream, & Revert:
             git checkout main
             git merge release/v6.494.0
             git push origin main
             git push origin release/v6.494.0
             git push origin v6.494.0
             git checkout <original_branch>
```

### 6.1 Phase Breakdown & Operational Directives

#### Step 1: Dedicated Release Branch Creation
- **Action**: Spawn and switch to an isolated branch dedicated exclusively to version `v6.494.0`.
- **Command**: `git checkout -b release/v6.494.0`
- **Rule**: Direct bumps on `main` without an intermediate release branch are strictly prohibited.

#### Step 2: Automated Manifest Bump via Python Script
- **Action**: Execute `37-bump-version.py` targeting the minor tier.
- **Command**:
  ```bash
  python 03-ai-scripts/37-bump-version.py -t minor --scope "Antigravity IDE Backup E2E Validation and Release"
  ```
- **Rule**: All 8 manifest locations must be updated synchronously in this step.

#### Step 3: Atomic Staging & Commit on Release Branch
- **Action**: Stage all modified version files, changelog updates, and the newly generated release notes document.
- **Commit Format**: Hyphen-separated or conventional commit string:
  ```bash
  git commit -m "release: v6.494.0 Antigravity IDE backup E2E validation and release"
  ```
- **Rule**: Commit must be isolated strictly to release-related manifest modifications.

#### Step 4: Cryptographic Git Tagging Ceremony
- **Action**: Affix an annotated, immutable Git tag pointing directly to the release commit.
- **Command**:
  ```bash
  git tag -a v6.494.0 -m "Release v6.494.0: Antigravity IDE Backup E2E & Release Ceremony"
  ```
- **Verification**: `git describe --tags --exact-match` must return `v6.494.0`.

#### Step 5: Merge Back to Main, Upstream Push, & Restoration
- **Action**: Fast-forward merge the release branch into `main`, push commits and tags to the remote repository, and restore the initial branch context.
- **Commands**:
  ```bash
  git checkout main
  git merge release/v6.494.0
  git push origin main
  git push origin release/v6.494.0
  git push origin v6.494.0
  ```
- **Rule**: Remote references (`origin`) receive the release branch, the updated `main`, and the new tag `v6.494.0`.

---

## 7. Positive Boolean Compliance & Quality Gates

All configuration objects, validation results, and CLI flags used across the release subsystem adhere to positive boolean naming conventions:

| Negative Identifier (PROHIBITED) | Positive Alternative (MANDATORY) | Semantic Role |
| :--- | :--- | :--- |
| `skipTests` | `shouldRunTests` / `isTestExecutionBypassed` | Quality gate test controller |
| `dirtyTree` | `isCleanState` | Git working tree status indicator |
| `noPush` | `shouldPushToRemote` | Remote upload toggle |
| `unverified` | `isVerified` | Cryptographic signature / manifest check |
| `notInherited` | `isInherited` | Version inheritance resolver |
| `disableSync` | `isSyncEnabled` | Manifest synchronization flag |
| `isNotDryRun` | `isLiveExecution` | Execution simulation guard |

---

## 8. Failure Modes & Self-Healing Protocols

1. **Dirty Working Tree Detected**:
   - *Condition*: `git status --porcelain` returns untracked or unstaged modifications prior to Step 1.
   - *Protocol*: Refuse execution. Require Worker 01 or Lead Orchestrator to stage or stash changes before initiating release branch.
2. **Version Collision / Half-Executed Bump**:
   - *Condition*: `version.json` already contains `6.494.0` while tag `v6.494.0` does not exist.
   - *Protocol*: Idempotency guard triggers. Resume from Step 3 without repeating arithmetic increment to prevent skipping to `v6.495.0`.
3. **Absolute Path Contamination**:
   - *Condition*: Release notes or changelog contains absolute filesystem roots, drive prefixes, or local file URI schemes.
   - *Protocol*: Validation pre-commit hook rejects the release. Path scrubber normalizes paths against the Git repository root.
