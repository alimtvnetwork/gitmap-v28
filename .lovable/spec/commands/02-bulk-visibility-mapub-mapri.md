---
slug: bulk-visibility-mapub-mapri
status: open
scope: cmd/visibility, db/migrations
captured: 2026-06-06
---

# Command: Bulk wildcard visibility flips — `make-all-public` / `make-all-private`

Full verbatim spec preserved in `.lovable/plans/pending/01-bulk-visibility-mapub-mapri.md`
Context section. Summary:

- Long forms: `gitmap make-all-public`, `gitmap make-all-private`
- Shorthands: `gitmap MAPUB`, `gitmap MAPRI`
- Syntax: `<cmd> <owner-or-org-url-or-folder> <pattern1>, <pattern2>, …` `[-Y]`
- Patterns: exact / `prefix*` / `*contains*` / `prefix*suffix`
- Interactive: numbered list → yes/no/exclude `1,3-5` → apply; `-Y` skips
- SQLite: `GitMapRun` (1) → `GitMapRepoResult` (n), enums small-int-backed
- Help MD + changelog + version bump required

**Confirmed NOT pre-existing** as of 2026-06-06:
- `grep` of `gitmap/constants/constants_cli.go` shows only `CmdMakePublic = "make-public"` and `CmdMakePrivate = "make-private"`.
- `gitmap/cmd/visibilitybulk.go` implements `<repo> <count>` (vN-N+1 form), NOT wildcard/owner-wide form.
- No `MAPUB` / `MAPRI` constants anywhere.
