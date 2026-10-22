# Migration & Verification Specification: Plan 293 — GitMap Plans Command and Fleet Skills Synchronization

## Status: APPROVED
- **Spec ID:** 293-03
- **Task Slug:** `293-gitmap-plans-command-and-fleet-sync`
- **Architecture Spec:** [01-architecture-spec.md](01-architecture-spec.md)
- **Component Spec:** [02-component-spec.md](02-component-spec.md)
- **Master Plan:** [.ai-memory/plans/pending/293-gitmap-plans-command-and-fleet-sync.md](../../../.ai-memory/plans/pending/293-gitmap-plans-command-and-fleet-sync.md)

---

## 1. Migration Protocol & Non-Destructive Invariants

### 1.1 Non-Destructive Addition
- All new functionality resides under `cli/cmdplans/`. No existing command dispatchers or data schemas are mutated destructively.
- The command reads markdown files read-only; no files under `.ai-memory/plans/` are moved, modified, or deleted by `gitmap plans`.

### 1.2 Backwards Compatibility
- `gitmap plans` and `gitmap plan` are new top-level subcommands.
- Existing flags across the CLI are not collided or overridden.

### 1.3 User Non-Negotiable Invariants
- **ZERO Compiler Builds (`go build` BANNED)**.
- **ZERO Automated Unit Tests (`go test` BANNED)**.
- **STRICTLY Relative Git Paths** (no absolute drive paths or file URI schemes).
- **Git Commit Command Invariant**: `gitmap cpf "<module> - <summary>"` (hyphen-separated only).

---

## 2. Verification Protocol (Zero-Build / Zero-Test Mode)

Since `go build` and `go test` are strictly prohibited per user instructions, verification is conducted through static analysis and structural audits:

1. **AST & Syntax Linting**:
   - Run `python 03-ai-scripts/50-fastgate.py` to inspect staged files for:
     - Forbidden absolute filesystem path patterns (`c:`, `d:`, `/home/`, `file:///`).
     - Go function line-count violations ($> 15$ lines).
     - File size limit ratchet ($> 500$ lines).
     - Forbidden boolean naming patterns (`can*`, `should*`, `was*`, `will*`).
2. **Coding Guideline Conformance**:
   - Zero nested `if` statements verified via guard clause audits.
   - Return styling verified: blank line after `}`, blank line before `return`.
   - Leaf error sovereignty verified: all helpers return `error` or `*apperror.AppError`; zero direct `cliexit.HandleError` or `os.Exit` calls in leaf modules.

---

## 3. Fleet Synchronization Ceremony

1. **Skill Update in Current Repository (`gitmap`)**:
   - Update `.agents/skills/gitmap/SKILL.md` to include:
     - Non-negotiable command replacement table entry for manual plan directory listing.
     - Usage guide for `gitmap plans` with flags (`--json`, `--yaml`, `--all`, `N`).
2. **Upstream Skill Update (`coding-guidelines`)**:
   - Mirror `.agents/skills/gitmap/SKILL.md` updates to `d:/work/coding-guidelines/.agents/skills/gitmap/SKILL.md`.
3. **Fleet Broadcast Execution**:
   - Run `python 03-ai-scripts/38-sync-prompts-skills-scripts.py` from repository root.
   - Confirm successful synchronization across the 43 connected repositories.
