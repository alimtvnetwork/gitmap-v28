# Master Plan 363: GitMap Synchronization and Version Control

- **Sequence ID:** 363  
- **Task Slug:** `gitmap-synchronization-and-version-control`  
- **Root Slug:** `gitmap-synchronization-and-version-control`  
- **Spec Reference:** `02-spec/21-app/363-gitmap-synchronization-and-version-control/01-architecture-spec.md`  
- **Status:** `pending`  
- **Created At:** 2026-10-31  

---

## 1. Executive Summary & Architectural Overview

As GitMap scales across multi-machine clusters, developer workstations, and autonomous AI swarms, maintaining strict synchronization of canonical specifications, prompts, skills, and version manifests becomes critical. Currently, synchronization scripts (`03-ai-scripts/38-sync-prompts-skills-scripts.py`) operate ad-hoc, lack native Go CLI commands, risk overwriting human modifications, and do not track asset integrity through cryptographically verified manifests.

Plan 363 delivers a comprehensive, production-grade synchronization and version control subsystem natively integrated into the GitMap CLI:
1. **Canonical Asset Ingestion & Manifest Generation:** Computes deterministic SHA-256 hash trees for prompts (`01-prompts/`), skills (`.agents/skills/`), specifications (`02-spec/`), and scripts (`03-ai-scripts/`), recording them inside `version.json` at repository root.
2. **Non-Destructive 3-Way Hash Verification & Human Conflict Protection:** Distinguishes between pristine canonical states, upstream changes, and local human edits. Automatically preserves modified local files, creates non-destructive backups (`.bak`), and halts overwriting unless explicitly forced.
3. **Dedicated CLI Subcommands:** Native first-class CLI routing for `gitmap cg sync`, `gitmap prompt sync`, and `gitmap skill sync` with `--dry-run`, `--force`, `--manifest-only`, and `--json` telemetry.
4. **On-The-Fly Terminal Execution Prompts:** CLI emission of standardized V6 (Tree Swarm Hierarchy) and V7 (Time-Restricted & Autobot Governed) prompt envelopes via `gitmap cg execute v6` and `gitmap cg execute v7`.
5. **Muse AI Onboarding & Help Subsystem:** Automated bootstrapping (`gitmap muse onboard`) and themed terminal guidance (`gitmap muse help`) integrating Muse memory IO and training vaults.
6. **Tokyo Night Themed Lipgloss Cards:** Terminal UI output featuring high-contrast Tokyo Night palettes, visual progress bars, diff indicators, and audit cards.
7. **Zero-Downtime Binary Swap & Release Ceremony:** In-place binary replacement via `gitmap swap-binary`, minor version bump to `v6.636.0`, and fleet-wide sync verification.

---

## 2. Architectural Invariants & Non-Negotiable Standards

1. **Leaf Error Sovereignty:** Leaf functions returning `error` or `*apperror.AppError` must bubble errors upward. Never call `os.Exit` or CLI exit handlers from leaf functions.
2. **Positive Boolean Naming:** All booleans must strictly begin with `is*` or `has*` (`isDryRun`, `isForce`, `hasConflict`, `isHumanModified`, `hasManifest`). Negative or banned prefixes (`can*`, `should*`, `was*`, `will*`) are strictly prohibited.
3. **Function Sizing & Zero Nesting:** Functions must remain within 8–15 lines. Flatten all conditional branching with early returns and guard clauses.
4. **Return Styling:** Mandatory blank line after closing brace `}` and blank line before `return`.
5. **Relative Git Paths:** 100% relative repository paths only (`02-spec/...`, `cli/...`). Zero hardcoded absolute paths (`C:\...`, `/home/...`).
6. **Zero-Build & Zero-Test Invariant:** No unrequested builds or E2E test suites during maintenance and planning. Fast-forwarding file-level operations only.
7. **Semantic Commit Disciplines:** Commits must use dedicated verbs (`gitmap cpr` for releases, `gitmap cpf` for features, `gitmap cbf` for bug fixes, `gitmap cfr` for refactors, `gitmap cpc` for chores).

---

## 3. Subtasks Breakdown Matrix

| Subtask ID | File | Description | Target Package / Files | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Task-01** | [01-subtask-architecture-and-specs.md](.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/01-subtask-architecture-and-specs.md) | Author architecture and component specs defining manifest schemas, 3-way hash algorithms, and CLI contracts. | `02-spec/21-app/363-gitmap-synchronization-and-version-control/` | Pending |
| **Task-02** | [02-subtask-canonical-assets-and-manifest.md](.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/02-subtask-canonical-assets-and-manifest.md) | Implement canonical asset scanner, SHA-256 hash tree generator, and `version.json` sync manifest writer. | `cli/cmdcg/`, `cli/cmdprompt/`, `pkg/manifest/` | Pending |
| **Task-03** | [03-subtask-hash-engine-and-human-protection.md](.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/03-subtask-hash-engine-and-human-protection.md) | Implement non-destructive 3-way hash verification, human modification detection, and conflict backup routines. | `cli/cmdsync/`, `cli/cmdcg/`, `pkg/hashengine/` | Pending |
| **Task-04** | [04-subtask-cli-subcommands-dispatch.md](.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/04-subtask-cli-subcommands-dispatch.md) | Wire `gitmap cg sync`, `gitmap prompt sync`, `gitmap skill sync` into root dispatcher with flag parsers. | `cli/root.go`, `cli/cmdcg/`, `cli/cmdprompt/` | Pending |
| **Task-05** | [05-subtask-cg-execute-terminal-prompts.md](.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/05-subtask-cg-execute-terminal-prompts.md) | Build `gitmap cg execute v6` and `v7` on-the-fly terminal prompt generator with clipboard and batch flags. | `cli/cmdcg/`, `cli/cmdprompt/` | Pending |
| **Task-06** | [06-subtask-muse-onboarding-and-help.md](.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/06-subtask-muse-onboarding-and-help.md) | Implement `gitmap muse onboard` environment setup, training vault link, and rich `gitmap muse help` command. | `cli/cmdmuse/`, `cli/helpdoc/` | Pending |
| **Task-07** | [07-subtask-helpdoc-and-terminal-styling.md](.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/07-subtask-helpdoc-and-terminal-styling.md) | Author Tokyo Night Lipgloss summary cards and markdown help documentation for sync and muse commands. | `cli/helpdoc/`, `cli/richhelp/` | Pending |
| **Task-08** | [08-subtask-verification-release-and-sync.md](.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/08-subtask-verification-release-and-sync.md) | Targeted unit verification, binary swap, minor bump release `v6.636.0`, and fleet-wide sync broadcast. | `version.json`, `pkg/constants/`, changelog | Pending |

---

## 4. Subtasks Detail & Deliverables

### Task-01: Specification & Architecture Contract Authoring
- Author `02-spec/21-app/363-gitmap-synchronization-and-version-control/01-architecture-spec.md` detailing the 3-way sync engine, SHA-256 manifest schema, and human protection invariants.
- Author `02-spec/21-app/363-gitmap-synchronization-and-version-control/02-component-and-cli-spec.md` with concrete Go struct definitions, CLI flag definitions, exit codes, and dispatcher routes.
- Define unambiguous state transition tables for file synchronization: `IN_SYNC`, `OUTDATED`, `HUMAN_MODIFIED`, `NEW_ASSET`, `ORPHANED`.

### Task-02: Canonical Asset Ingestion & version.json Manifest Generation
- Build the asset crawler scanning canonical roots: `01-prompts/`, `.agents/skills/`, `02-spec/`, and `03-ai-scripts/`.
- Compute canonical SHA-256 hashes for each asset, discarding non-semantic CRLF/LF line-ending differences.
- Extend `version.json` with the canonical manifest section:
  ```json
  "syncManifest": {
    "manifestVersion": "1.0.0",
    "generatedAt": "2026-10-31T18:00:00Z",
    "assetCount": 142,
    "promptsHash": "e3b0c442...",
    "skillsHash": "8f43a21b...",
    "assets": {
      "01-prompts/01-unified-ai-prompt-v4.md": "a1b2c3d4...",
      ".agents/skills/gitmap/skill.md": "b2c3d4e5..."
    }
  }
  ```
- Implement deterministic alphabetical key sorting to guarantee zero spurious git diffs.

### Task-03: Non-Destructive Hash Verification & Human Conflict Algorithm
- Implement the 3-way comparison algorithm:
  - `CanonicalHash`: hash of asset in canonical repository.
  - `BaseHash`: hash recorded in target repository `version.json` at last sync.
  - `DiskHash`: current hash of asset file on target repository disk.
- Conflict Rules:
  - If `DiskHash == BaseHash` and `DiskHash != CanonicalHash`: Safe to update (no local changes).
  - If `DiskHash == CanonicalHash`: Up to date (`IN_SYNC`).
  - If `DiskHash != BaseHash` and `DiskHash != CanonicalHash`: Local file has been modified (`HUMAN_MODIFIED`).
- Human Protection Routine:
  - By default, do NOT overwrite `HUMAN_MODIFIED` files.
  - Emit clear warning with diff pointers.
  - If `--backup` is enabled, write `.bak` sidecar file before touching.
  - If `--force` is enabled, overwrite local file with canonical source.

### Task-04: CLI Subcommand Dispatch (`gitmap cg sync`, `gitmap prompt sync`, `gitmap skill sync`)
- Wire subcommands into `cli/root.go`:
  - `gitmap cg sync [targets...] [--dry-run] [--force] [--json]`
  - `gitmap prompt sync [targets...] [--dry-run] [--force]`
  - `gitmap skill sync [targets...] [--dry-run] [--force]`
  - Universal alias: `gitmap sync [prompts|skills|cg|all]`
- Implement multi-repository target resolver discovering sibling workspaces or respecting explicit directories.
- Return structured `Result[SyncReport]` envelopes with leaf error sovereignty.

### Task-05: On-The-Fly Terminal Execution Prompts (`cg execute v6` / `v7`)
- Implement `gitmap cg execute v6` (Tree Swarm Hierarchy: A=3, H=2, Capacity=6, micro-batches 5-8 files).
- Implement `gitmap cg execute v7` (Time-Restricted & Autobot Governed: `< 08:00 AM` weekday cutoff, single-user bot identity).
- Provide CLI flags:
  - `--clipboard` / `-c`: copy prompt directly to system clipboard.
  - `--raw`: print plain text prompt without Lipgloss decorative cards.
  - `--batch <size>`: parameterize micro-batch file limit.

### Task-06: Muse Onboarding & Muse Help Terminal Command Engine
- Implement `gitmap muse onboard` / `gitmap muse init`:
  - Verify and register AI bot identity: `gitmap login set ai-auto muse`.
  - Validate `Muse/` memory directory structure, chat seeds, and training vault links.
  - Seed initial training digest cache if empty.
- Implement `gitmap muse help` / `gitmap thelp muse`:
  - Detailed workflow guide for Muse multi-day digest consumption, LLM training loops, and on-call skill triggers.

### Task-07: Help Documentation & Terminal UI Tokyo Cards Integration
- Author rich markdown help documents under `cli/helpdoc/`:
  - `cli/helpdoc/cg_sync.md`
  - `cli/helpdoc/prompt_sync.md`
  - `cli/helpdoc/skill_sync.md`
  - `cli/helpdoc/muse.md`
- Integrate with Tokyo Night themed Lipgloss renderers:
  - Display dual-column sync tables: Asset Name, Status, Local State, Action.
  - Highlight human conflicts in bold Amber/Gold (`#ffb86c`).
  - Highlight successful updates in Spring Green (`#50fa7b`).

### Task-08: Release Ceremony, Binary Swap & Fleet Sync Verification
- Perform targeted unit test passes for hash engine, manifest generation, and dispatcher routing.
- Execute zero-downtime binary swap using `gitmap swap-binary` on Windows and POSIX.
- Perform release ceremony:
  - Bump `version.json` and constants to `v6.636.0`.
  - Update `changelog.md` and release notes.
  - Commit using dedicated semantic release verb: `gitmap cpr "v6.636.0 gitmap synchronization and version control"`.
  - Broadcast canonical sync across connected fleet nodes.

---

## 5. Verification & Acceptance Criteria

- [ ] All 8 subtask files exist under `.ai-memory/plans/subtasks/363-gitmap-synchronization-and-version-control/`.
- [ ] `version.json` manifest structure supports deterministic SHA-256 generation.
- [ ] 3-way comparison algorithm correctly flags human-modified files without data loss.
- [ ] `gitmap cg sync`, `gitmap prompt sync`, `gitmap skill sync` route cleanly without exit handler violations.
- [ ] `gitmap cg execute v6` and `v7` emit complete, parameter-accurate execution instructions.
- [ ] `gitmap muse onboard` and `gitmap muse help` render complete environment checks and Tokyo cards.
- [ ] Help documents are fully registered in `cli/helpdoc/` and rich help menus.
- [ ] All code conforms to positive booleans, 8-15 line functions, zero nested ifs, and relative paths.
