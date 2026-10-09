# Plan 246 — Native Parallel Fix Command + LLM Train Integration + UI/UX Overhaul

## User request (verbatim)
Owner 2026-10-09: "there was a AI script for fixing the file and coding right you can check want you to have that feature inside GitMap and also add this feature in the LLM train and also skills the reason I'm saying this because GitMap should be able to fix this Unicode encoding or any other issues parallelly yeah applying the grouping something like this so it it should have this process dealing with and fixing this parallelly and we can change the worker threads so this is how it should be changeable from the CLI" + "check the settings part settings UI these were not good and most of the UI that we have not very up so please define or have your knowledge to improve the UI UX and UI themes make sure these are very lucrative" — "Go ahead if there is any confusion let me know. Let me know the reports."

## Scope
- Repo: `~/workspace/repos/gitmap-v28`, main at d018f68 (pulled 2026-10-09, tree clean).
- Backup branch `backup/246-parallel-fix-command` created from d018f68 and pushed to origin (via git directly — `gitmap space backup-branch` correctly refused on task-245's untracked spec files).
- Spec 246 is the next number (245 exists from the parallel stream).

## Key findings from research
- 8 fixer scripts in `03-ai-scripts/` (~1096 lines): encoding, newlines, guideline-composite, relative-paths, naming (audit-only), gofmt (worker pool), misspell (30-word dict, has case-mangling bug), md-gap. Inconsistent CLIs; shared engine `02-shared-engine.py`.
- **Name collision (lead decision D2):** `gitmap fix` / `cli/cmdfix/` are TAKEN (git-state remediation: stash/wip/discard). New command = **`gitmap autofix`** (alias `afx`), new package **`cli/cmdautofix/`**.
- **LLM train is print-only fiction:** no phase registry, no execution. Program 246 builds a real `TrainPhase` interface + registry; phase 6 "Heal & Fix" = (a) git-state via existing `cmdfix`, (b) content audit via new `cmdautofix` (report-only in train).
- **UI defects (verified):** dead light-theme switcher, fake page routing, unparsed `-p/--port/--host/--no-browser` flags, SSE claim with no endpoint, help menu advertises non-existent tabs, emoji iconography, `alert()` dialogs, undefined `--accent`, hardcoded demo data. Assets = Go raw-string constants (no build step) — keep that model.

## Conflicts
- Skill R1 bans `go build`/`go vet`; the OWNER explicitly authorized builds for this task. Top-Instruction Priority Mandate: owner instruction wins. `go test` suites stay off-limits; `go vet` + CLI verification is the bar. Logged, not hidden.

## Wave plan
- Wave 0 (lead): spec 246 (4 files) + 6 subtasks + this plan + index updates — one atomic commit.
- Wave 1 (A=2): Worker 1 → subtasks 01+02 (`cli/cmdautofix/` engine + CLI, dispatch registration). Worker 2 → subtasks 03+04 (llm TrainPhase + heal phase; skills docs).
- Wave 2 (A=2): Worker 1 → subtask 05 (UI themes: real dark+light, fix dead switcher, `--accent`). Worker 2 → subtask 06 (UI polish: SVG icons, toasts, scales, responsive, demo-data removal, flag parsing, help-menu alignment).
- Wave 3 (lead): `go build ./...` + `go vet` on touched packages (personally verified, never trusted from workers); smoke-test `gitmap autofix --help`, dry-run on scratch tree, `gitmap llm train --text-only`, load settings UI page; atomic `gitmap cpf` + push per wave.

## Checkboxes
- [x] Step 0: git pull (d018f68), backup branch `backup/246-parallel-fix-command` pushed.
- [x] Research wave: 8-script deep-dive + llm-train/UI deep-dive (both reports in).
- [x] Spec wave: 4 spec files + 6 subtasks written; lead reconciled heal-phase to drive both `cmdfix` (git-state) and `cmdautofix` (content).
- [x] Wave 0 commit: spec + plan + indexes (commit pushed via `gitmap cpf`).
- [x] Wave 1a: llm train phase + skills (03, 04) — DONE (Worker 2): `llm_phases.go`, `llm_heal.go` (two sub-steps), text renames, skills §10 in all 3 homes byte-identical; `go build` + `go vet` exit 0.
- [x] Wave 1b: `cli/cmdautofix/` engine + CLI (01, 02) — DONE (Worker 1): 17 files, `fix` parent + 9 subcommands, check→summary→prompt flow, scan cache, old `fix` displaced from rootcore (stash/wip/discard kept). Lead fixes: rich-help `fix` routing, non-TTY prompt hang (300ms piped-input window), declined-prompt exit code 1→correct. Lead-verified: `go build ./...` exit 0, all exit codes correct, `fix --help` renders new parent help.
- [ ] Wave 1 commit: NOTE — parallel spec-248 stream's `git add -A` swept 246's new files into commit d5d3e55; lead's follow-up fixes committed as 26df08b. All content pushed. Lesson logged in AGENTS.md (targeted `git add` when streams share a branch).
- [x] Wave 2a: UI polish (06) — DONE (Worker B): flag parsing (`-p/--port/--host/--no-browser`), help-menu honesty (10 real tabs, SSE claim removed), demo-data removal, sticky save bar, icon+label pills, alert→toast; `go build` exit 0 (lead-verified). Note: touched `ui_server.go` minimally (2 spots) to wire host/no-browser through — required, flagged.
- [ ] Wave 2b: UI themes (05) — worker running.
- [ ] Wave 2: UI themes (05); UI polish (06).
- [ ] Wave 3: lead verification (build, vet, smoke tests, settings UI load); final report to owner.

## Decisions (lead)
- D1: Command named `autofix`/`afx` (fix/cmdfix taken). D2: dry-run default, `--apply` writes. D3: `--workers/-w`, default = NumCPU(). D4: exit codes 0/1/2 (CI compat). D5: misspell case-preserving (fix script bug, don't replicate). D6: REPO_FILE_URI parameterized via `--uri-pattern`. D7: train heal phase = cmdfix (git-state, `--heal-apply` applies) + cmdautofix (content, report-only always). D8: UI keeps Go-string asset model (no build pipeline). D9: help menu aligned to 10 real tabs (no invented pages). D10: SSE claim removed, not implemented.
- D11 (OWNER REFINEMENT 2026-10-09, overrides D1–D4): `fix` IS the parent command with per-category subcommands (`encoding|newlines|naming|paths|gofmt|misspell|markdown|guidelines|all`). Check-driven: scan → mandatory summary → prompt `[y/N]`; `-y`/`--yes` applies. `fix all` = one summary, one prompt. `--workers N` (default CPU) + scan-result cache. The old git-state `fix` is DISPLACED (entry removed from rootcore.go); `stash`/`wip`/`discard` entries untouched. Package stays `cli/cmdautofix/`; entry `RunFixCmd`. Train sub-step B calls `cmdautofix.Scan` directly (never prompts).

## Assumptions
- `cli/cmdautofix` will not import `cli/cmd/llm` (worker verifies at build; fallback: heal phase skips content sub-step with a warning).
- Task-245's untracked spec files remain untouched (another workstream).

## Follow-ups (out of scope)
- 500–1000-line band splits (program 243 ledger).
- HelpDisplay SampleOutput gap (program 243 ledger).
- `go test` suite for the new packages (owner has not authorized test suites).
