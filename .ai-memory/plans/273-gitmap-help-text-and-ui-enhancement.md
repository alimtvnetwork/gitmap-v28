# Program 273 — gitmap-help-text-and-ui-enhancement

## Status
- [x] Step 0: Sync, spec 273, backup branch `backup/273-gitmap-help-text-and-ui-enhancement` pushed
- [ ] WS1: Spec + plan + subtasks (FIRST — before any implementation)
- [ ] WS2: Help text + theme overhaul (coding-guidelines-v24 design system)
- [ ] WS3: Missing commands — implement `merge write` / `mr`
- [ ] WS4: 20 new command proposals (planned task, proposals ONLY)
- [ ] WS5: Windows virus RCA (commit 385aba38e ground truth)
- [ ] WS6: Lead final verify + commit + push (NO version bump, NO release)

## Scope authority
User order 2026-10-16: improve terminal help text + UI help, new themes/colors from
coding-guideline design system, `help example` per command, missing commands
(esp. `merge write` / `mr` for Lovable migration), 20 new command proposals in
.ai-memory, Windows virus root-cause analysis in .ai-memory.

## Hard rules (non-negotiable)
- Spec + plan + subtasks BEFORE implementation (WS1 first).
- `gitmap aum search` / `gitmap find` / `gitmap cat` / `gitmap ps` / `gitmap py` / `gitmap llm train` ONLY. TOTAL BAN on grep/rg/git grep/Select-String.
- RELATIVE paths only in all artifacts. NEVER absolute paths (incl. release page/notes).
- NEVER commit binaries. NEVER delete files. ALIASES KEPT.
- Files ≤300 lines. No `go test` locally.
- Atomic `gitmap cpf "<module> - <summary>"` (hyphen, NO colons). Commits MUST push.
- Dumb-terminal safety: piped output must not corrupt (no ANSI in pipes).

## Waves
### WS1 — Spec + plan + subtasks
- [ ] Spec under 02-spec/21-app/gitmap-help-text-and-ui-enhancement/ (overview, theme design, help-example format, missing-command audit, acceptance criteria)
- [ ] Plan in .ai-memory/plans/gitmap-help-text-and-ui-enhancement.md
- [ ] Subtasks in .ai-memory/plans/subtasks/gitmap-help-text-and-ui-enhancement/
- [ ] Commit + push spec BEFORE implementation

### WS2 — Help text + theme overhaul
- [ ] Read coding-guidelines-v24 design system FIRST — themes/colors from it, not invented
- [ ] `help example` for EVERY command: what it does + concrete invocation + expected output
- [ ] Data-driven, single source of truth; UI help screens same treatment
- [ ] No piped-output corruption on dumb terminals

### WS3 — Missing commands
- [ ] Audit dispatcher vs help; implement `merge write` / `mr` (title, body, source→target, draft flag)

### WS4 — 20 new command proposals
- [ ] Planned task in .ai-memory/plans/ — proposals ONLY, do NOT implement

### WS5 — Windows virus RCA
- [ ] Ground truth: commit 385aba38e "eliminated VirTool:MSIL/Subti.C heuristic triggers"
- [ ] Document AV, heuristic, fix, professional prevention (signing, SmartScreen, vendor submission, reproducible builds) with costs — for another AI to continue
