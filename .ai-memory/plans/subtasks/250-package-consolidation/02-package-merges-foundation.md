# Subtask: 250 foundation package merges (groups 1–5, 10)

Owner: Worker. Runs after orphan triage (01). Foundation-first order: each group lands on a compiling tree. Hubs (`constants`, `apperror`, `store`, `cliexit`, `model`) are untouched (D4). After every group: `go build ./...` clean AND `go list` confirms no new import cycles.

## Group 1 — error/diagnostic core → `cli/diag` (new)

- [ ] Move `appfault` (1f), `suggestion` (5f), `errreport` (3f) into `cli/diag` (9 files).
- [ ] Rewrite importers via `gitmap aum` (low count; `apperror` with 191 importers is NOT touched).
- [ ] Acceptance: build clean, no new cycles, importers rewritten.

## Group 2 — JSON helpers → `cli/jsonx` (new)

- [ ] Move `stablejson` (3f), `jsonenvelope` (5f), `jsonenv` (2f) into `cli/jsonx` (10 files).
- [ ] Rewrite importers via `gitmap aum`.
- [ ] Acceptance: build clean, no new cycles, importers rewritten.

## Group 3 — terminal output → `cli/termout` (new)

- [ ] Move `termtable` (4f), `termpad` (3f), `termhelp` (6f), `theme` (4f) into `cli/termout` (17 files).
- [ ] Rewrite importers via `gitmap aum`.
- [ ] Acceptance: build clean, no new cycles, importers rewritten.

## Group 4 — secrets → `cli/secrets` (new)

- [ ] Move `ghtoken` (4f), `secretsresolver` (2f), `crypto` (6f) into `cli/secrets` (12 files).
- [ ] Rewrite importers via `gitmap aum`.
- [ ] Acceptance: build clean, no new cycles, importers rewritten.

## Group 5 — path/temp/lock → `cli/fspath` (new)

- [ ] Move `tempdir` (3f), `localdirs` (1f), `scripts` (1f), `lockfile` (2f), `lockcheck` (4f) into `cli/fspath` (11 files).
- [ ] Rewrite importers via `gitmap aum`.
- [ ] Acceptance: build clean, no new cycles, importers rewritten.

## Group 10 — help/docs → `cli/helpdoc` (new)

- [ ] Move `helptext` (10f incl. `docs/cmd`), `fixtureversion` (7f), `tool/helptextemitter` (1f) into `cli/helpdoc` (18 files); dissolve empty `tool` wrapper.
- [ ] Rewrite importers via `gitmap aum`.
- [ ] Acceptance: build clean, no new cycles, importers rewritten.

## Ordering constraints

1. Orphan triage (01) completes first.
2. Groups 1–5 + 10 may run in sequence in any order (disjoint importer sets; group 1 first recommended since JSON error envelopes reference error types).
3. Groups 6–8 (middle) and group 9 (dispatcher micros) are separate follow-on subtasks — not in this file.

Global acceptance: all six groups landed, 63 files re-homed into 5 target packages, behavior identical, `go build ./...` clean, zero import cycles.
