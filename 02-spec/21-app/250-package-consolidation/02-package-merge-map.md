# 250-package-consolidation — Package Merge Map

Concrete merge groups. File counts measured 2026-10-09 (`find <pkg> -name '*.go'` under `cli/`). "Importer impact" = packages importing any member (leaf-or-few recommended first). All paths relative to repo root. Merged package names lowercase.

## Group 1 — error/diagnostic core (FOUNDATION, ordered after orphan triage)

Members: `apperror` (5 files), `appfault` (1, orphan), `suggestion` (5), `errreport` (3). Total: 14 files.
Rationale: all four own error construction, fault reporting, suggestion text, and error reporting. Cycle-safety: zero import cycles in the module (Tarjan SCC confirmed), so the merge is cycle-safe by construction.
RECOMMENDATION: do NOT fold `apperror` (191 importers — rewrite cost is the program's largest single churn). Instead merge `appfault` + `suggestion` + `errreport` (combined ~6 external importers expected) into a new hub-adjacent package `apperror` is NOT touched:
- Target: `cli/diag` (new). Members: `appfault`, `suggestion`, `errreport`. Total 9 files.
- Importer impact: low (all three are leaf/near-leaf; `apperror` importers untouched).
- Ordering: first foundation group; rewrites via `gitmap aum` bottom-up.

## Group 2 — JSON helpers (FOUNDATION)

Members: `stablejson` (3), `jsonenvelope` (5), `jsonenv` (2, orphan). Total: 10 files.
Target: `cli/jsonx`. Rationale: deterministic JSON marshaling, envelope types, and env-backed JSON helpers share importers (`cliexit`, output paths) and never import each other — leaf-safe. Triage note: `jsonenv` was already merged-or-keep candidate; include it here.
Ordering: second; after group 1 (error text flows into JSON error envelopes).

## Group 3 — terminal output (FOUNDATION)

Members: `termtable` (4), `termpad` (3), `termhelp` (6), `theme` (4). Total: 17 files.
Target: `cli/termout`. Rationale: table rendering, padding, help-text layout, and theme tokens are one concern (styled terminal output); importers are all presentation leaves. Cycle check: none of the four import each other today.
Ordering: third; parallel-safe with groups 1–2 (disjoint importer sets).

## Group 4 — secrets (FOUNDATION)

Members: `ghtoken` (4), `secretsresolver` (2), `crypto` (6). Total: 12 files.
Target: `cli/secrets`. Rationale: token storage, secret resolution, and encryption are one trust boundary; `secretsresolver` imports `crypto` semantics already adjacent.
Ordering: fourth; after group 1 (error types referenced by secret failures).

## Group 5 — path/temp/lock (FOUNDATION)

Members: `tempdir` (3), `localdirs` (1), `scripts` (1), `lockfile` (2), `lockcheck` (4). Total: 11 files.
Target: `cli/fspath`. Rationale: temp dirs, well-known local dirs, bundled scripts, file locks, and lock preflight checks = filesystem plumbing. All leaf-level.
Ordering: fifth; no cross-group dependencies (disjoint from 1–4).

## Group 6 — scan pipeline (MIDDLE)

Members: `scanner` (8), `mapper` (10), `indexer` (2), `worker` (1), `probe` (4). Total: 25 files.
Target: `cli/scanpipe`. Rationale: the scan→map→index→work→probe pipeline is one dataflow; each package is already single-stage. Merging keeps the stage files separate (still ~300-line files) inside one package.
Ordering: after foundation groups (uses term output + jsonx + fspath); cycle-safe by construction.

## Group 7 — committransfer subs (MIDDLE)

Members: `committransfer/graph` (2), `committransfer/prdesc` (2), `committransfer/prclean` (2, orphan). Total: 6 files.
Target: `cli/committransfer` (fold subpackages into parent). Rationale: graph, PR description, and PR cleanup are all committransfer internals; `prclean` orphan status makes it a merge-not-delete candidate.
Ordering: after orphan triage (resolves prclean verdict); independent of groups 1–6.

## Group 8 — cmd/commitin leaves (MIDDLE)

Members: `cmd/commitin/checkpoint` (3), `cmd/commitin/dedupe` (2), `cmd/commitin/finalize` (3), `cmd/commitin/replay` (5), `cmd/commitin/walk` (4), `cmd/commitin/prompt` (2, orphan). Total: 19 files.
Target: `cli/cmdcommitin` (fold leaves into the parent; keep orchestrator/message/profile/runlog/workspace/funcintel/e2e separate for now).
Rationale: checkpoint/dedupe/finalize/replay/walk/prompt are leaf stages of the commit-in flow, imported only by their parent. `prompt` orphan → merge verdict. `e2e` (9 files, orphan) stays separate (test-only).
Ordering: after orphan triage; independent of groups 1–6.

## Group 9 — dispatcher micros (LATE)

Members: the 24 one-file non-test `cmd*` packages: `cmdappend`, `cmdaum`, `cmdbash`, `cmdcfrppriorversion`, `cmdclean`, `cmdct`, `cmdexplorer`, `cmdheadtail`, `cmdjsonextract`, `cmdlocate`, `cmdmkdir`, `cmdmuse`, `cmdpin`, `cmdpowershell`, `cmdrest`, `cmdsafe`, `cmdsafety`, `cmdsee`, `cmdsends`, `cmdsetsourcerepo`, `cmdsf`, `cmdstash`, `cmdversion`, `cmdwhoami`. Total: 24 files.
Target: `cli/cmd` (fold into the existing dispatcher package as `cmd_*.go` files).
Rationale: verified `cli/cmd` is imported ONLY by the root package and NO cmd* package imports `cli/cmd` — every leaf is fold-safe. `cmdct` and `cmdjsonextract` are orphans → merge verdict already applies.
Ordering: LAST — every other merge must be green first (dispatcher is the widest importer of everything).

## Group 10 — help/docs (FOUNDATION)

Members: `helptext` (10 files incl. `docs/cmd`), `fixtureversion` (7), `tool/helptextemitter` (1, orphan: `tool` has no other files). Total: 18 files.
Target: `cli/helpdoc`. Rationale: help text, fixture versioning, and the help-text emitter are one documentation surface. Orphan triage verdict for `tool/helptextemitter`: merge into `helpdoc` (the `tool` wrapper package disappears; it held nothing else).
Ordering: alongside foundation groups 1–5 (disjoint importer sets).

## Program totals

Foundation groups 1–5 + 10: 6 groups, 63 files → 5 target packages (group 1a keeps `apperror` untouched).
Middle groups 6–8: 3 groups, 50 files. Late group 9: 24 files.
Package count reduction: ~28 packages absorbed (20 orphan triage + merges), target packages all stay within the ~300-line-per-file convention.
