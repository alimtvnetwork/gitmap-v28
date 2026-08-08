# spec/03-commit-in

Authoritative, AI-blind-ready specification for the `gitmap-v28 commit-in` (`cin`)
command. Read in numerical order. Implementation is **forbidden** until the
user explicitly says `next` (per `.lovable/plan.md`).

## Iterations

| #  | File                                | Scope                                                               |
|----|-------------------------------------|---------------------------------------------------------------------|
| 01 | `01-overview-and-glossary.md`       | Verbatim, vocabulary, invariants, scope boundaries                  |
| 02 | `02-cli-surface.md`                 | Argv grammar, separators, flags, prompts, exit codes                |
| 03 | `03-pipeline.md`                    | Stage-by-stage pipeline (`CommitInStage` enum), Mermaid flow        |
| 04 | `04-database-schema.md`             | Every table, column, FK, index, ERD, enum-mirror tables             |
| 05 | `05-profiles-and-json-shape.md`     | Profile JSON shape (PascalCase), file layout, default-binding rule  |
| 06 | `06-message-and-function-intel.md`  | Message rules order, weak-word policy, per-language function detect |
| 07 | `07-acceptance-and-edge-cases.md`   | Acceptance matrix, ambiguity resolutions, conformance tests         |
| 08 | `08-tag-mirroring-and-release-branches.md` | Annotated-tag mirroring, version-tag detection, auto release branches, RewrittenCommit columns, T1–T10 |
| 09 | `09-commit-in-replay-map.md`        | `CommitInReplayMap` table — annotated-tag old↔new SHA map, `TagReplayOutcome` enum, cross-run idempotency lookup, R1–R10 |

## Source-of-truth rules (apply to ALL files in this folder)

1. PascalCase for every table name, column name, JSON key, JSON value, enum
   member. Never magic strings — every classifier value is an `Enum` mirrored
   to a join table.
2. Every primary key is `INTEGER PRIMARY KEY AUTOINCREMENT` named
   `<PascalCaseTableName>Id` (database convention §4 of project memory).
3. No file content, no file hash. Only `RelativePath` per source SHA.
4. Replicate BOTH original `AuthorDate` AND `CommitterDate` on re-commit
   unless author override is set (then keep both dates, only swap identity).
5. Source SHA → new SHA mapping is persisted in `ShaMap`. Same source SHA
   ever seen again (any run, any input repo) → SKIP + log `DuplicateSourceSha`.
6. File-system surface limited to: SQLite path, clone temp path, profile
   JSON path, log path. No project structure changes.
7. `<source>` auto-init rule: URL → clone into CWD; existing repo → reuse;
   existing non-repo folder → `git init` in place; missing path → `mkdir -p`
   then `git init`. Never refuse, never prompt for this.