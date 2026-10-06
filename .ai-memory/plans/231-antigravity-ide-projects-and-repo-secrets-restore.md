# Master Execution Plan: 231-antigravity-ide-projects-and-repo-secrets-restore

## User Request (Verbatim)
```text
# High Priority Instruction

please all repos following the antigravity ide here in the project sections along with pin projects, and also list out your work and steps into the repo-secrets folder for restoring projects and settings from one antigravity to another create folder and also can use gitmap required to

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. List all repositories following the Antigravity IDE in the project sections along with pin projects.
5. Document your work and steps into the repo-secrets folder for restoring projects and settings from one Antigravity to another.
6. Create the necessary folder structure for the above tasks.
7. Use GitMap as required for the tasks.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[.agents/skills/gitmap] skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

---

## 1. Executive Summary & Architecture Strategy

This plan establishes end-to-end multi-agent orchestration for task `231-antigravity-ide-projects-and-repo-secrets-restore`:
1. **Antigravity IDE Project Sections (78 Repositories)**: Symmetrically register all 78 repositories across direct and nested namespaces into Antigravity IDE project descriptors, mapping valid RFC 3986 file URIs, default git branches, and permissions.
2. **Pinned Projects Quick Access (23 Repositories)**: Establish Tier 1 (7 core infrastructure orchestrators) and Tier 2 (16 key ecosystem services) pinned project hierarchies in `pinned_projects.json`.
3. **Restoration Runbook & Toolchain in `repo-secrets/09-antigravity-backup/`**: Document complete Standard Operating Procedure (SOP) for migrating and restoring Antigravity projects, settings, and pins between environments (Windows, Linux, macOS), complete with portable JSON manifests and automated POSIX Bash / PowerShell restoration scripts.
4. **Strict Path Relativity & Positive Booleans**: Enforce 100% relative paths (`02-spec/...`, `.ai-memory/...`, `repo-secrets/...`) with zero absolute host prefixes or drive letters in tracked files, and affirmative boolean naming conventions.

---

## 2. Multi-Agent Delegation & Ownership Boundaries

| Phase / Role | Subagent Type | Owned Paths & Deliverables | Status |
| :--- | :--- | :--- | :--- |
| **Phase 1: Discovery** | Research 01 & 02 (`research`) | Read-only analysis of `repo-secrets/gitmap-final.json` and spec schemas | **DONE** |
| **Phase 1: Spec Authoring** | Spec Agent 01 (`self`) | `02-spec/21-app/231-antigravity-ide-projects-and-repo-secrets-restore/01-architecture-spec.md` & Subtask 01 | **DISPATCHED** |
| **Phase 1: Spec Authoring** | Spec Agent 02 (`self`) | `02-spec/21-app/231-antigravity-ide-projects-and-repo-secrets-restore/02-component-spec.md` & Subtask 02 | **DISPATCHED** |
| **Phase 2: Execution** | Worker 01 (`self`) | Project section definitions, 78 project descriptors, `pinned_projects.json` | **QUEUED** |
| **Phase 2: Execution** | Worker 02 (`self`) | `repo-secrets/09-antigravity-backup/` runbooks, manifests, restoration scripts | **QUEUED** |
| **Phase 3: Consolidation** | Lead Orchestrator | Register updates, verification gate execution, final consolidation | **QUEUED** |

---

## 3. Subtask Decomposition

### Subtask 01: Antigravity IDE Projects Ingestion & Pinned Projects Suite
- **File**: `.ai-memory/plans/subtasks/231-antigravity-ide-projects-and-repo-secrets-restore/01-antigravity-projects-and-pins.md`
- **Owner**: Worker 01
- **Scope**:
  - Ingest 78 Git repositories (42 direct, 36 nested across 9 namespaces).
  - Structure project sections into 10 logical categories for IDE navigation.
  - Symmetrically configure project descriptors with preserved UUIDs for core repositories (`gitmap`, `scripts-fixer`, `wp-html-automate`, `antigravity-manager`, `letsmarknow`, `letsmarknow-ui`).
  - Configure `pinned_projects.json` with 23 Tier 1 and Tier 2 projects.

### Subtask 02: Repo-Secrets Migration Documentation & Automation Toolchain
- **File**: `.ai-memory/plans/subtasks/231-antigravity-ide-projects-and-repo-secrets-restore/02-repo-secrets-restore-docs-and-scripts.md`
- **Owner**: Worker 02
- **Scope**:
  - Establish `repo-secrets/09-antigravity-backup/` with `vault/`, `scripts/`, `temp/`.
  - Author master SOP runbook `repo-secrets/09-antigravity-backup/readme.md`.
  - Create portable manifests: `projects-manifest.json`, `pinned-projects.json`, `settings-manifest.json`, `plugins-and-skills.json`.
  - Implement cross-platform zero-touch restoration scripts (`restore-antigravity-all.sh`, `restore-antigravity-all.ps1`, `verify-antigravity-backup.sh`).
  - Register `09-antigravity-backup/` in `repo-secrets/readme.md`.

---

## 4. Verification Gates (VG-01 to VG-07)

1. **VG-01 Vault Structure & Manifest Validation**: All manifests exist and validate against affirmative boolean schema.
2. **VG-02 Repository Count Parity**: Exactly 78 repositories indexed in manifests.
3. **VG-03 Path Relativity & Portability**: Zero drive letters (`C:`, `D:`) or `/home/...` paths in manifests or tracked markdown files.
4. **VG-04 Restoration Script Parity**: Both POSIX Bash and PowerShell scripts support dry-run, backup, and restore flags.
5. **VG-05 Atomic Non-Destructive Ingestion**: Restore routines preserve existing local credentials without destruction.
6. **VG-06 Secret Sanitization Gate**: Zero raw API keys, tokens, or plaintext passwords in manifests.
7. **VG-07 Documentation & Relative Link Hygiene**: 100% relative Git paths in all committed documentation.
