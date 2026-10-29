# Executable spec ticket

**Version:** 4.3.0
**Updated:** 2026-10-02
**Status:** Active

A ticket is one change a stranger can implement and a stranger can fail. Fill every heading below, in this order.

Fill `Context` with five labeled lines: Who, Current, Desired, Why now, Done when.

`Current state` names a verification date and at least one `path:line` citation.

`Acceptance criteria` are numbered. Each one is pass or fail.

`Testing plan` is a table with columns Layer, What, Count.

`Files` is a table with columns File, Change. Paths start at the repo root. No drive letters. No `file:///` URIs.

`Out of scope` and `Do not touch` each contain at least one bullet.

`Checklist` repeats the acceptance criteria as `- [ ]` items.

A ticket with more than 8 rows in `Files` must be split before implementation starts.

## Context

Who: the next agent asked to add a heading checker for this ticket shape.
Current: module templates describe a whole folder, not one pass-or-fail change.
Desired: one checker that fails when a ticket heading is missing or out of order.
Why now: a weak agent treats a long template as permission to invent the missing heading.
Done when: the checker prints `PASS` on this file and fails when `## Rollback` is absent.

## Current state

Verified 2026-10-02. `02-spec/01-spec-authoring-guide/08-non-cli-module-template.md:1` is a module shape. It is not a one-change ticket.

## Proposed change

Add a standard-library Python checker. It reads one or more paths. With no path, it reads this file only. A file is in scope when it contains the exact line `# Executable spec ticket`, or both `## Context` and `## Do not touch`. In-scope files must contain the ten level-2 headings in the order listed in this template. Other files are ignored.

## Acceptance criteria

1. This file contains the ten level-2 headings in the template order. Pass or fail.
2. `python linter-scripts/check-spec-ticket-headings.py` prints `PASS` and exits 0. Pass or fail.
3. The same command exits non-zero on a copy that lacks the `## Rollback` line. Pass or fail.
4. The file table in this example lists at most the checker and this template. Pass or fail.

## Testing plan

| Layer | What | Count |
| :--- | :--- | :--- |
| Command | Checker on this template | 1 pass |
| Command | Checker on a copy missing `## Rollback` | 1 fail |

## Rollback

Revert the commit that added the checker. No data migration.

## Files

| File | Change |
| :--- | :--- |
| `linter-scripts/check-spec-ticket-headings.py` | Create the heading checker |
| `02-spec/01-spec-authoring-guide/15-executable-spec-ticket.md` | Keep this template as the heading source |

## Out of scope

- Rewriting existing modules into tickets.
- Registering the checker in CI.

## Do not touch

- `version.json`
- Token values under `02-spec/07-design-system/`

## Checklist

- [ ] This file contains the ten level-2 headings in the template order
- [ ] `python linter-scripts/check-spec-ticket-headings.py` prints `PASS` and exits 0
- [ ] The same command exits non-zero on a copy that lacks the `## Rollback` line
- [ ] The file table in this example lists at most the checker and this template
