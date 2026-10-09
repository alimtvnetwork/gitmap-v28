# Subtask: 250 orphan triage

Owner: Worker 1. Prerequisite for ALL merges (D2). For each of the 20 orphans: confirm zero importers via `gitmap aum search`, record the verdict, then move to the merge groups that consume the verdict. Owner hard rule: NEVER delete without explicit ask — default verdict is merge-or-keep.

Verdicts confirmed by Worker 1 on 2026-10-09. "Confirmed" = importer count verified via `gitmap aum search "gitmap-v28/cli/<path>" -e .go` from `cli/`; files read for a one-line rationale.

| # | Orphan package | Files | Verdict | Evidence (importers) | Notes |
|---|---|---|---|---|---|
| 1 | `appfault` | 1 | merge into `cli/diag` (group 1) | 0 importers (confirmed 2026-10-09) | — |
| 2 | `cmdct` | 1 | merge into `cli/cmd` (group 9) | 0 importers (confirmed 2026-10-09) | package-main dispatcher shim; 1 file |
| 3 | `cmdcodingguidelines` | 3 | keep | 0 importers (confirmed 2026-10-09) | keep — coding-guideline flow dispatcher (handler + commit + compat, 3 files); wire-up may come later |
| 4 | `cmddaemon` | 4 | keep | 1 importer: `tests/fleet_deploy_clone_ui_e2e_test.go` (test-only) | keep — REST daemon client/server/types entry point (4 files incl. own test); wire-up may come later |
| 5 | `cmdjsonextract` | 1 | merge into `cli/cmd` (group 9) | 0 importers (confirmed 2026-10-09) | package-main extractor shim; 1 file |
| 6 | `cmd/gitmap-node-join` | 1 | merge into `cli/cmdcluster` (sibling) | 0 importers (confirmed 2026-10-09) | merge — package-main cluster node-join helper (1 file); node-join is a cluster subcommand in spirit |
| 7 | `cmd/commitin/e2e` | 9 | keep (test-only) | 0 importers (confirmed 2026-10-09) | keep — commin e2e test surface (9 files); test-only, stays separate |
| 8 | `cmd/commitin/prompt` | 2 | merge into `cli/cmdcommitin` (group 8) | 0 importers (confirmed 2026-10-09) | prompt package (bufio-based, 2 files) |
| 9 | `committransfer/prclean` | 2 | merge into `cli/committransfer` (group 7) | 0 importers (confirmed 2026-10-09) | 2 files |
| 10 | `completion/internal/gencommands` | 1 | keep (main package) | 0 importers (confirmed 2026-10-09) | keep — `package main` go:generate generator (1 file) emitting completion constants; cannot merge into a library |
| 11 | `db/zombiezen` | 1 | merge into `cli/db` (sibling) | 0 importers (confirmed 2026-10-09) | merge — single pure-Go SQLite driver adapter behind a build tag (1 file) |
| 12 | `enums` | 2 | keep | 0 importers for `cli/enums` root; `cli/enums/ctxmodetype` has 9 importers (`cli/cmdinstall`, incl. tests) | keep (CORRECTION to the prior "merge into cli/cmdcommitin"): ctxmodetype is a Windows shell context-menu enum actively imported by cmdinstall — unrelated to commitin; the `enums` tree is the designated closed-set home per spec 243 Wave D; 2 files (doc.go + variant.go) |
| 13 | `helptext/docs/cmd` | 4 incl. md | merge into `cli/helpdoc` (group 10) | none | docs generation leaf |
| 14 | `jsonenv` | 2 | merge into `cli/jsonx` (group 2) | 0 importers (confirmed 2026-10-09; the 15 substring hits were `jsonenvelope`, a different package) | uniform `--json` envelope package; 2 files |
| 15 | `logging` | 2 | keep | 0 importers (confirmed 2026-10-09) | keep — structured `--log-json` NDJSON sink reusing the jsonenv envelope (2 files incl. test); candidate for coding-guidelines reuse scan first |
| 16 | `tests` | 81 | keep | 0 importers (confirmed 2026-10-09) | keep — integration test root (81 go files / 82 total); never merge into library packages |
| 17 | `tests/e2e` | 15 | keep | 0 importers (confirmed 2026-10-09) | keep — e2e surface (15 files); never merge |
| 18 | `tool/helptextemitter` | 1 | merge into `cli/helpdoc` (group 10) | 0 importers (confirmed 2026-10-09) | `tool` wrapper dissolves (holds nothing else); 1 file |
| 19 | `utils` | 1 | keep | 0 importers (confirmed 2026-10-09) | keep — single-file goroutine pool (`async_pool.go`, 1 file); candidate for coding-guidelines reuse scan first |
| 20 | TBD (brief lists 20, enumerated 19) | — | confirm during triage | none | owner to confirm the 20th orphan; do not guess |

- [x] Rows 1–20: importer count confirmed via `gitmap aum search` for each package (Worker 1, 2026-10-09; row 20 TBD still open — not an importer question).
- [x] Rows 1–19: verdict recorded (merge target / keep) — no deletions. One correction: row 12 `enums` changed from "merge into `cli/cmdcommitin`" to keep (evidence above).
- [ ] Merge-consuming verdicts handed to the group owners (groups 1, 2, 7, 8, 9, 10).
- [ ] Keep verdicts listed in the final report with their rationale (coding-guidelines reuse scan may still absorb `logging`/`utils` later).

Acceptance: all 20 rows checked, zero deletions, `go build ./...` clean after the verdict-driven moves.
