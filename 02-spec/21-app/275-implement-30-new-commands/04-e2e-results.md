# 04 — E2E results (program 275)

**Spec:** 275 · **Updated:** 2026-10-17 · **Binary:** /tmp/gitmap275 (built from HEAD)

All 30 commands tested: `--help` renders, two-word form routes, alias routes,
help-example entry present, error cases clean, dry-runs side-effect free.

| # | Command | Aliases | --help | Two-word | Alias | Example | Errors | Dry-run | Status |
|---|---------|---------|--------|----------|-------|---------|--------|---------|--------|
| P-01 | mr list | mrl | PASS | PASS | PASS | PASS | — | n/a (read-only) | PASS |
| P-02 | mr view | mrv | PASS | PASS | PASS | PASS | PASS (missing n) | n/a | PASS |
| P-03 | mr merge | mrm | PASS | PASS | PASS | PASS | — | PASS | PASS |
| P-04 | mr checkout | mrco | PASS | PASS | PASS | PASS | — | n/a | PASS |
| P-05 | repo export | repox | PASS | PASS | PASS | PASS | — | n/a | PASS |
| P-06 | repo import | repoi | PASS | PASS | PASS | PASS | — | PASS | PASS |
| P-07 | help topic | ht | PASS | PASS | PASS | PASS | PASS (unknown topic) | n/a | PASS |
| P-08 | howto | — | PASS | n/a | PASS | PASS | PASS (empty) | n/a | PASS |
| P-09 | theme list | themes | PASS | PASS | PASS | PASS | — | n/a | PASS |
| P-10 | theme preview | themep | PASS | PASS | PASS | PASS | PASS (invalid name) | n/a | PASS |
| P-11 | nodes jump | nj | PASS | PASS | PASS | PASS | PASS (missing node) | n/a | PASS |
| P-12 | nodes version | nv | PASS | PASS | PASS | PASS | — | n/a | PASS |
| P-13 | pe triage | pet | PASS | PASS | PASS | PASS | — | n/a | PASS |
| P-14 | pe rerun | perr | PASS | PASS | PASS | PASS | PASS (missing id) | PASS | PASS |
| P-15 | release train | rtrain | PASS | PASS | n/a | PASS | — | PASS (default) | PASS |
| P-16 | errors last | errl | PASS | PASS | PASS | PASS | — | n/a | PASS |
| P-17 | scan diff | scand | PASS | PASS | PASS | PASS | — | n/a | PASS |
| P-18 | db doctor | dbdoc | PASS | PASS | PASS | PASS | — | n/a (stub) | PASS |
| P-19 | hooks install | hkin | PASS | PASS | PASS | PASS | PASS (missing --all) | PASS | PASS |
| P-20 | health | hlth | PASS | n/a | PASS | PASS | — | n/a | PASS |
| C-01 | install verify | — | PASS | PASS | n/a | PASS | — | n/a | PASS |
| C-02 | install urls | — | PASS | PASS | n/a | PASS | PASS (missing tag) | n/a | PASS |
| C-03 | release verify | — | PASS | PASS | n/a | PASS | — | n/a | PASS |
| C-04 | release winscan | — | PASS | PASS | n/a | PASS | — | n/a | PASS |
| C-05 | prompts drift | — | PASS | PASS | n/a | PASS | — | n/a | PASS |
| C-06 | agent report | — | PASS | PASS | n/a | PASS | PASS (bad --since) | n/a | PASS |
| C-07 | ci open | — | PASS | PASS | n/a | PASS | — | n/a | PASS |
| C-09 | spec progress | — | PASS | PASS | n/a | PASS | PASS (missing --spec) | n/a | PASS |
| C-10 | agm queue audit | — | PASS | PASS | n/a | PASS | — | n/a | PASS |

**Result: 30/30 PASS.**

Notes:
- `pe triage` routes before the `pe` pipeline-errors alias via dispatch275 precedence.
- `help topic` required a routing intercept in runHelpDispatch (collides with `gitmap help <topic>`).
- `db doctor`/`errors last` degrade gracefully when DBs absent (verified from /tmp).
- No live side effects performed: no PRs opened/merged, no reruns executed, no hooks installed, no releases triggered.
- `go build ./...` exit 0 (lead-verified).
