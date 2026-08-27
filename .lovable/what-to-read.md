# What to Read

> Canonical map of what the AI must read before working on this project.
> Last updated: 2026-08-10T04:40:24Z

## Changelog

- 2026-08-22T04:46:47Z, Enqueued massively granular codebase audit plan for v1.4.5 coding guidelines (453 violations mapped).

- 2026-08-10T04:40:24Z, Implemented LFS smudge error detection and auto-fix fallback in `clonefrom`.
- 2026-08-09T12:08:25Z, Verified canonical what-to-read map for the project.
- 2026-08-09T10:23:25Z, Initialized canonical what-to-read map for the project.

## Before any task (always)

- `.lovable/memory/index.md`, why: Contains the central index of all project features, constraints, bugs, and design decisions.
- `.lovable/coding-guidelines.md`, why: Enforces strict TS Enum suffix rules, TS/Python query wrappers, Golang styles, and structural bans.
- `.lovable/plans/index.md`, why: Lists active and completed tasks to ensure efforts don't conflict with ongoing work.
- `.lovable/plan.md`, why: The current single-file tracker for the active AI tasks.
- `.lovable/strictly-avoid.md`, why: Contains the absolute list of project prohibitions.
- `.lovable/ambiguous-questions/01-new-ambiguity/`, why: Open questions that need resolution from the user.

## Before writing code

- `.lovable/memory/release-architecture-map.md`, why: Details the repository Single Source of Truth versioning, component inheritance, and release propagation rules.

- `.lovable/memory/style/code-quality-improvement.md`, why: Defines architectural resilience, dynamic paths, generic wrapping, and iteration reduction rules.
- `.lovable/memory/style/01-ts-enums-and-query-wrappers.md`, why: Defines the strict requirement for `*Type` enums instead of TS string unions, and query logging wrappers.
- `.lovable/memory/tech/01-version-json-architecture.md`, why: Explains that `version.json` is the sole source of truth and `constants.go` is injected via `-ldflags`. Do not use powershell arrays.

## Before adding a feature

- `spec/readme.md`, why: The spec folder structure dictates how feature blueprints are organized.

## Before writing a spec

- `.lovable/memory/project/what-to-read.md`, why: The internal project memory version of what to read, which outlines JSON output contracts and step-by-step feature recipes.

## Before adding a unit test

- `spec/01-app/57-skipmeta-integration-test.md`, why: Details the project's integration test strategy and assertions against metadata behavior.

## See also

- Root `README.md` (must stay in sync with this file)

.lovable/plans/pending/02-git-rm-and-folder.md
gitmap/cmd/folder/folder.go
gitmap/cmd/gitrm/gitrm.go

.lovable/plans/pending/03-ignore-and-add.md
gitmap/cmd/ignore/ignore.go
gitmap/cmd/add/add.go
