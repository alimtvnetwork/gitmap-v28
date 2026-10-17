# 01 — Overview: implement-30-new-commands

**Version:** 1.0.0
**Updated:** 2026-10-17
**Spec:** 275

---

## 1. Owner order (verbatim)

> "Go ahead implement and do e2e testing and realease please"
> (approving implementation of all 30 proposals from programs 273 + 274,
> with E2E testing and a release)

## 2. Build specs (source of truth per command)

- `.ai-memory/plans/273-twenty-command-proposals.md` — P-01..P-20
- `.ai-memory/plans/274-ten-new-command-proposals.md` — C-01..C-10

Each proposal defines: problem, invocation, expected output, aliases,
effort/value, gain, depends-on. Workers implement exactly what the
proposal describes. Proposals are NOT re-designed here.

## 3. The 30 commands, grouped for workers

### Group A — mr family (cli/cmdmerge/)
| ID | Command | Aliases | Notes |
|----|---------|---------|-------|
| P-01 | `mr list` | `mrl` | `--state open`; table of PRs via `gh pr list` |
| P-02 | `mr view <n>` | `mrv` | checks/reviews/mergeable via `gh pr view` |
| P-03 | `mr merge <n>` | `mrm` | `--strategy squash`; requires `--yes` or confirm |
| P-04 | `mr checkout <n>` | `mrco` | `gh pr checkout` |

Build on existing `cli/cmdmerge/` (merge_write.go pattern). Reuse the
`gh` host layer (merge_write_gh.go).

### Group B — repo + scan (cli/cmdrepo/, cli/cmdscan/)
| ID | Command | Aliases | Notes |
|----|---------|---------|-------|
| P-05 | `repo export` | `repox` | reuse `scan export` manifest format |
| P-06 | `repo import` | `repoi` | clone + register + restore remotes |
| P-17 | `scan diff` | `scand` | vs last-scan snapshot |
| C-08 | `scan verify` | — | index↔filesystem both directions; `--repair` |

### Group C — help + theme (cli/cmdhelp/, cli/cmdtheme/)
| ID | Command | Aliases | Notes |
|----|---------|---------|-------|
| P-07 | `help topic` | `ht` | topic index: releases, pipeline, fleet, migration |
| P-08 | `howto <task>` | — | task→recipe; keyword matching, no AI call |
| P-09 | `theme list` | `themes` | from 273 theme slugs (navy-purple, tokyo-night, vscode-dark-plus, monokai-pro + legacy) |
| P-10 | `theme preview <name>` | `themep` | render sample without switching |

### Group D — install + release (cli/cmdinstall/, cli/cmdrelease/)
| ID | Command | Aliases | Notes |
|----|---------|---------|-------|
| C-01 | `install verify` | — | binary/PATH/symlink/source checklist |
| C-02 | `install urls <tag>` | — | HEAD-check installer script URLs |
| C-03 | `release verify` | — | asset count, draft, tag→HEAD parity |
| C-04 | `release winscan <tag>` | — | 7 Defender patterns from `.ai-memory/plans/273-windows-virus-rca.md` |
| P-15 | `release train` | `rtrain` | ordered multi-repo release; `--dry-run` default-safe |

### Group E — pipeline + fleet + hygiene (cli/cmdpe/, cli/cmdfleet/, cli/cmdhealth/)
| ID | Command | Aliases | Notes |
|----|---------|---------|-------|
| P-13 | `pe triage` | `pet` | worst-failure bundle |
| P-14 | `pe rerun <run>` | `perr` | `gh run rerun`; failed-jobs-only |
| P-11 | `nodes jump <node>` | `nj` | ssh + cd (prints command if no TTY) |
| P-12 | `nodes version` | `nv` | version skew table |
| P-16 | `errors last` | `errl` | last tracked error from error module |
| P-18 | `db doctor` | `dbdoc` | sizes + staleness; `--repair` |
| P-19 | `hooks install --all` | `hkin` | commit-msg guard across repos; `--dry-run` |
| P-20 | `health` | `hlth` | aggregation dashboard |
| C-05 | `prompts drift` | — | canonical vs copies by content hash |
| C-06 | `agent report --since` | — | from agent task SQLite DB |
| C-07 | `ci open [repo]` | — | deep link to failing job; Chrome |
| C-09 | `spec progress --spec` | — | parse spec parts + ledger checkboxes |
| C-10 | `agm queue audit` | — | queue vs dispatch receipts |

## 4. Non-goals

- No live side effects in implementation or E2E (no real PRs opened/merged,
  no real releases triggered, no hooks installed without `--dry-run`).
- No re-design of proposals. No new commands beyond the 30.
- The 272-rejected items stay parked (not applicable here — none of the 30
  touch them).
