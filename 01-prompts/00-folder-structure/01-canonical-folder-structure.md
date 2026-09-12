# Canonical Folder Structure & Sizing Rules — Architecture Spec (must follow)

> **Prompt Version:** 2.1.0
> **Synchronization:** Single Source of Truth Mirror for `.lovable/folder-structure.md`
> **Target:** `01-prompts/00-folder-structure/01-canonical-folder-structure.md`

Every AI prompt and agent MUST strictly adhere to the central repository layout defined below and in `.lovable/folder-structure.md`. No invented folders, no alternate paths.

All dates are UTC. All filenames are lowercase kebab-case with a two-digit zero-padded `XX-` prefix where sequencing applies. `XX-` is monotonic within its folder scope.

```text
/                                      # Root of the repository
  spec/                             # Canonical layered specifications (source of truth)
    01-spec-authoring-guide/           # Meta-guide for writing and maintaining specs
    02-coding-guidelines/              # Detailed cross-language coding architectures
    03-error-manage/                   # Error management, envelopes, error codes
    04-database-conventions/          # Database schema design, ORM mappings, views
    05-split-db-architecture/          # Operational vs. configuration DB patterns
    06-seedable-config-architecture/   # Seedable configuration with versioned changelogs
    07-design-system/                  # Design tokens, UI component hierarchy, dark mode
    08-docs-viewer-ui/                 # Documentation viewer UI & presentation system
    09-code-block-system/              # Syntax highlighting pipeline & interaction system
    10-research/                       # Technology benchmarks and comparative studies
    11-powershell-integration/         # PowerShell scripting and automation conventions
    12-cicd-pipeline-workflows/        # CI/CD workflows, GitHub Actions, quality gates
    13-generic-cli/                    # Generic CLI UX, terminal colors, flag semantics
    14-update/                         # CLI self-update architecture & release delivery
    15-distribution-and-runner/        # Cross-platform installers, runners, forwarding
    16-generic-release/                # Universal release blueprint & asset matrices
    17-consolidated-guidelines/        # Single-file AI summaries of all major spec modules
    18-wp-plugin-how-to/               # WordPress plugin architecture & REST APIs
    19-main-worker-service/            # Split-tier worker service & proxy nodes
    21-app/ ... 24-app-ui-design-system/ # Application-specific business features & schemas

  .lovable/                            # AI metadata, cognitive memory, and automation layer
    01-index.md                        # Master repository index and directory router
    02-user-preferences                # Explicit user communication preferences
    03-strictly-avoid.md               # Universal hard prohibitions & CODE RED constraints
    04-suggestions.md                  # Architectural suggestions & pending improvements
    05-plan.md                         # Active high-level roadmap and milestone tracker
    06-what-to-read.md                 # Reading sequence and entrypoint priority list
    folder-structure.md                # 📄 AUTHORITATIVE SINGLE SOURCE OF TRUTH FOR FOLDER STRUCTURE

    ai-fix-scripts/                    # Reusable high-speed Python automation toolchain
      01-index.md                      # Master catalog & search tag registry
      02-shared-engine.py              # Central engine: constants, lazy regex, locks, cache
      03..20-*.py                      # Specialized linters, fixers, and local CI runners

    plans/                             # Micro-task execution center
      01-index.md                      # Master registry of active and completed plans
      pending/                         # High-level active parent plans (xx-<slug>.md)
      subtasks/                        # Bounded micro-tasks (xx-<slug>/01-step.md)
      completed/                       # Archived completed parent plans & subtasks

    memory/                            # Long-term institutional cognitive memory (SINGULAR)
      01-index.md                      # Master table of contents & CODE RED rules
      architecture/                    # System architecture decisions & split-DB patterns
      constraints/                     # Non-negotiable technical constraints
      done/                            # Historical completed milestones
      features/                        # Feature-specific knowledge and requirements
      issues/                          # Root cause analyses and bug postmortems
      patterns/                        # Reusable design and coding patterns
      processes/                       # Operational and maintenance procedures
      project/                         # Project metadata and author attributions
      standards/                       # Coding standards and enum specifications
      style/                           # Code style and naming conventions
      suggestions/                     # Granular suggestion trackers
      workflow/                        # Sprint and migration state trackers

    prompts/                           # Prompt repository & category hierarchy
      01-prompts-category/             # Canonical source for categorized prompt modules (01-22)
      *.md                             # Synced flat prompts generated from categories

    release/                           # Release automation & version bumping hub
      release-method.md                # Authoritative version bump registry
      bump_versions.py                 # Automated regex-based version mutator
      issues/                          # Release failure logs and diagnostics

    question-and-ambiguity/            # Ambiguity resolution & iteration counters
      readme.md                        # Open and resolved ambiguity log
      task-counter.md                  # Task iteration metrics

    suggestions/                       # Detailed suggestion proposals
      01-index.md                      # Index of active suggestions
      completed/                       # Archived implemented proposals

    cicd-issues/                       # CI pipeline failure diagnostics and RCAs
    assets/                            # Mockups, diagrams, and reference media
```

---

## Numbering Rules (`XX-` Monotonic Sequencing)

1. Two-digit zero-padded, monotonic per folder scope (`01-`, `02-`, ...).
2. Next `XX-` = max existing `XX-` in that folder scope + 1. Never reuse or decrement.
3. Plans: `XX-` is monotonic across `pending/`, `subtasks/`, and `completed/` combined.

---

## Hard Rules & Prohibitions

- **Single Authority:** [`.lovable/folder-structure.md`](.lovable/folder-structure.md) is the absolute single source of truth.
- **Strict Plurality:** Always `.lovable/memory/` (singular), `.lovable/plans/` (plural), `01-prompts/` (plural), `.lovable/suggestions/` (plural).
- **No Orphaned Folders:** Never invent unapproved folders (e.g., `memories/`, `tasks/`, `todos/`).
- **Relative Paths Only:** All file paths must be relative from the git root. Absolute paths (`C:\...`, `file:///...`) are banned.
