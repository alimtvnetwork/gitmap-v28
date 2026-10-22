# Master Plan 293: GitMap Plans Command and Fleet Skills Synchronization

## Status: COMPLETED
- **Plan ID:** 293
- **Slug:** `293-gitmap-plans-command-and-fleet-sync`
- **Master Spec Reference:** [01-architecture-spec.md](../../../02-spec/21-app/293-gitmap-plans-command-and-fleet-sync/01-architecture-spec.md)
- **Component Spec Reference:** [02-component-spec.md](../../../02-spec/21-app/293-gitmap-plans-command-and-fleet-sync/02-component-spec.md)
- **Migration & Verification Spec Reference:** [03-migration-and-verification-spec.md](../../../02-spec/21-app/293-gitmap-plans-command-and-fleet-sync/03-migration-and-verification-spec.md)
- **Scope:** Implementation of `gitmap plans` command, Tokyo Night Lipgloss rendering, JSON/YAML serializers, skill updates in gitmap and coding-guidelines, and fleet synchronization.
- **Execution Mandate:** Zero-Build, Zero-Test, Strictly Relative Lowercase Git Paths

---

## 1. Executive Summary & Wave Organization

This master plan coordinates the implementation of the new `gitmap plans` CLI command and the subsequent multi-repository skill propagation across the 43 connected fleet repositories.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ Wave 1: Types, Parser & Formatter Engine (Subtasks 01–03)                   │
│ Spec: 02-spec/21-app/293-gitmap-plans-command-and-fleet-sync/01-architecture│
├─────────────────────────────────────────────────────────────────────────────┤
│ Wave 2: Lipgloss UI Rendering & Command Wiring (Subtasks 04–05)             │
│ Spec: 02-spec/21-app/293-gitmap-plans-command-and-fleet-sync/02-component   │
├─────────────────────────────────────────────────────────────────────────────┤
│ Wave 3: Skill Documentation & Fleet Synchronization (Subtasks 06–07)        │
│ Spec: 02-spec/21-app/293-gitmap-plans-command-and-fleet-sync/03-migration   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Complete 7-Subtask Registry & Traceability Matrix

### Wave 1: Types, Parser & Formatter Engine (Subtasks 01–03)

- [x] **Subtask-01:** [subtask-01.md](../subtasks/293-gitmap-plans-command-and-fleet-sync/subtask-01.md)
  - **Target:** `cli/cmdplans/plans_types.go`, `cli/constants/constants_commands_extra.go`
  - **Objective:** Define data structures, options, enums ending in `Type`, and command name constants.
  - **Status:** Completed

- [x] **Subtask-02:** [subtask-02.md](../subtasks/293-gitmap-plans-command-and-fleet-sync/subtask-02.md)
  - **Target:** `cli/cmdplans/plans_parser.go`
  - **Objective:** Implement filesystem scanner and metadata extractor for `.ai-memory/plans/`.
  - **Status:** Completed

- [x] **Subtask-03:** [subtask-03.md](../subtasks/293-gitmap-plans-command-and-fleet-sync/subtask-03.md)
  - **Target:** `cli/cmdplans/plans_format.go`
  - **Objective:** Implement JSON and YAML serializers for machine-readable output.
  - **Status:** Completed

### Wave 2: Lipgloss UI Rendering & Command Wiring (Subtasks 04–05)

- [x] **Subtask-04:** [subtask-04.md](../subtasks/293-gitmap-plans-command-and-fleet-sync/subtask-04.md)
  - **Target:** `cli/cmdplans/plans_render.go`
  - **Objective:** Render Tokyo Night Lipgloss tables, status pills, and pending backlog overflow guidance.
  - **Status:** Completed

- [x] **Subtask-05:** [subtask-05.md](../subtasks/293-gitmap-plans-command-and-fleet-sync/subtask-05.md)
  - **Target:** `cli/cmdplans/plans.go`, `cli/cmd/rootutility.go`
  - **Objective:** Implement `RunPlans(args []string) error` and wire into command dispatch.
  - **Status:** Completed

### Wave 3: Skill Documentation & Fleet Synchronization (Subtasks 06–07)

- [x] **Subtask-06:** [subtask-06.md](../subtasks/293-gitmap-plans-command-and-fleet-sync/subtask-06.md)
  - **Target:** `.agents/skills/gitmap/SKILL.md`, `../coding-guidelines/.agents/skills/gitmap/SKILL.md`
  - **Objective:** Document `gitmap plans` command, flag taxonomy, and replacement matrix.
  - **Status:** Completed

- [x] **Subtask-07:** [subtask-07.md](../subtasks/293-gitmap-plans-command-and-fleet-sync/subtask-07.md)
  - **Target:** `03-ai-scripts/38-sync-prompts-skills-scripts.py`, `03-ai-scripts/50-fastgate.py`
  - **Objective:** Execute fleet synchronization and fastgate pre-commit audit.
  - **Status:** Completed
