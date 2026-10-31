# Subtask 01: Specification & Architecture Contract Authoring

- **Parent Plan:** [363-gitmap-synchronization-and-version-control](../../pending/363-gitmap-synchronization-and-version-control.md)
- **Subtask ID:** `TASK-01`
- **Slug:** `gitmap-synchronization-and-version-control - Task 01`
- **Target Path:** `02-spec/21-app/363-gitmap-synchronization-and-version-control/`
- **Status:** `pending`
- **Assigned Worker:** Architecture Subagent

---

## 1. Goal

Author the foundational architecture specification and component CLI contract for GitMap's next-generation synchronization and version control subsystem. Fully specify the 3-way hash verification state machine, the `syncManifest` schema in `version.json`, human conflict protection protocols, and CLI subcommand interfaces prior to implementation.

---

## 2. In Scope

1. **Architecture Spec Authoring (`01-architecture-spec.md`):**
   - High-level system architecture and component interaction flow.
   - Canonical asset topology: `01-prompts/`, `.agents/skills/`, `02-spec/`, `03-ai-scripts/`.
   - Manifest schema specification for `version.json` (SHA-256 trees, metadata headers, ordering rules).
   - 3-Way Hash state machine: Canonical Hash, Target Base Hash, Target Disk Hash.
   - Human Modification Conflict Protection protocol (non-destructive detection, sidecar backups, override switches).
   - Muse AI onboarding lifecycle and memory integration contract.

2. **Component & CLI Spec Authoring (`02-component-and-cli-spec.md`):**
   - Detailed Go structs and types: `SyncManifest`, `AssetHashEntry`, `SyncReport`, `ConflictActionType`, `SyncOptions`.
   - Command-line syntax and flags for:
     - `gitmap cg sync [targets...] [--dry-run] [--force] [--json] [--manifest-only]`
     - `gitmap prompt sync [targets...] [--dry-run] [--force]`
     - `gitmap skill sync [targets...] [--dry-run] [--force]`
     - `gitmap cg execute [v6|v7] [--clipboard] [--raw] [--batch N]`
     - `gitmap muse onboard [--check-only] [--init-vault]`
     - `gitmap muse help`
   - Exit code mappings and strongly-typed AppError integration (`*apperror.AppError`).
   - Tokyo Night Lipgloss palette specifications and card layouts.

---

## 3. Out of Scope

- Go code implementation of the sync engine or CLI dispatchers (deferred to Tasks 02–06).
- Running builds or unit tests during spec authoring.
- Version bumping or release operations.

---

## 4. Architectural Invariants

- **Relative Paths Only:** 100% relative repository paths (`02-spec/...`, `cli/...`).
- **Positive Booleans Only:** All flags and struct fields must use `is*` or `has*` (`isDryRun`, `isForce`, `hasConflict`, `isHumanModified`).
- **Function Limit:** Design helper methods and interfaces adhering to <= 15 lines max.
- **Error Sovereignty:** All APIs must return `error` or `*apperror.AppError` directly to outer callers.

---

## 5. Verification Checklist

- [ ] `02-spec/21-app/363-gitmap-synchronization-and-version-control/01-architecture-spec.md` exists and covers state machines, manifest schemas, and conflict logic.
- [ ] `02-spec/21-app/363-gitmap-synchronization-and-version-control/02-component-and-cli-spec.md` exists and defines Go types, CLI syntax, and Lipgloss tokens.
- [ ] Manifest format matches `version.json` single-source-of-truth guidelines.
- [ ] Human conflict algorithm guarantees zero silent data loss.
