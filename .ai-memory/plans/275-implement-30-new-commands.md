# Program 275 — implement-30-new-commands

**Owner order:** implement all 30 proposals (273 P-01..P-20 + 274 C-01..C-10),
E2E test, release.
**Spec:** `02-spec/21-app/275-implement-30-new-commands/`
**Backup:** `backup/275-implement-30-new-commands` (pushed before code)

## Status
- [x] Step 0: pull, spec 275, backup pushed, mindset + proposals read
- [ ] WS1: spec + plan + subtasks (this file)
- [ ] WS2: implement 30 commands (3 workers)
- [ ] WS3: E2E matrix (30 commands)
- [ ] WS4: fix E2E findings
- [ ] WS5: minor bump + release + green CI

## Worker assignment (disjoint file boxes)
- **Worker 1 — Groups A+B:** `cli/cmdmerge/` (mr list/view/merge/checkout);
  `cli/cmdrepo/`, `cli/cmdscan/` (repo export/import, scan diff/verify)
- **Worker 2 — Groups C+D:** `cli/cmdhelp/`, `cli/cmdtheme/` (help topic,
  howto, theme list/preview); `cli/cmdinstall/`, `cli/cmdrelease/`
  (install verify/urls, release verify/winscan/train)
- **Worker 3 — Group E:** `cli/cmdpe/`, `cli/cmdfleet/`, `cli/cmdhealth/`,
  `cli/cmdprompts/`, `cli/cmdagent/`, `cli/cmdspec/`, `cli/cmdagm/`
  (pe triage/rerun, nodes jump/version, errors last, db doctor,
  hooks install, health, prompts drift, agent report, ci open,
  spec progress, agm queue audit)

Shared files owned by LEAD ONLY: `cli/cmd/rootdata.go`,
`cli/constants/*`, `cli/helpdoc/examples*.go`, `.ai-memory/plans/readme.md`,
`02-spec/21-app/readme.md`. Workers report needed registrations to lead;
lead applies them.

## Waves
| Wave | Work |
|------|------|
| 1 | WS1 spec/plan (lead) → commit + push |
| 2 | WS2 implementation (3 workers parallel) |
| 3 | Lead: registration + helpdoc wiring, `go build ./...` |
| 4 | WS3 E2E matrix |
| 5 | WS4 fixes |
| 6 | WS5 release (lead) |

## Rules
- `gitmap aum search` ONLY. Relative paths. No binaries. No deletions.
- Workers NEVER run git commands. Commits via `gitmap cpf` by lead only.
- Files ≤300 lines. Aliases kept. No `go test` locally.
- No live side effects (no real PRs/merges/releases).
