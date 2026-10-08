# Plan: 243-cli-help-displayer-and-backup-branch

## Objective
Implement the 243 spec series inside the staging-repo mindset: `gitmap space backup-branch`
(one-command backup branches), the CLI help displayer system (`Displayer` interface +
`HelpDisplay` struct + command helpers + theme inheritance + smart suggestions, DRY),
enforcement of the file-size cap, the `cli/cmd/` package split, a dedicated `cli/enums/`
package, and stale-doc fixes.

## User Request (Verbatim)
> # High Priority Instruction
>
> I wanted to see what you can do. Some of the things I do agree, some of the things I don't agree. Before making the changes, I want you to understand the mindset. That's the first thing, and also I want you to add that mindset inside the what to read file as well. You mentioned that many places the enforcement didn't happen. We need to work on it. [...] Create a backup branch. [...] I want you to have a functionality inside GitMap that could deal with this backup branch as well, backup branch function. That means there would be GitMap space backup-branch, and then we can give a name or the task string that would automatically create the slug and the backup/ the branch name and take a backup from the current code. [...] I also want you to go with the help section [...] we want to have a dry principle or dry mechanism that could display the help very effectively. [...] It would be an object, Go object. [...] Displayer would be an interface. Help display would be its own struct, and we can bind it through the binding principle [...] Inside this it would have an array of command helpers. Each one of the command helpers can also have sub command helpers [...] The command help group will have the header and command. Each item will have a command and the description and an example for the command and a URL. [...] we can have a theme on top of it, which can define how it should look like. In each group, we can have hints or suggestions. The suggestions also can have the inputs. [...] Suggestions would be very powerful. [...] I want to have a different way, different power, so that I can change the coloring of the header, coloring of the left side, command text, right side text. [...] Also, I could do overriding of the inheritance of the colors. [...] This is a very big task. I want you to understand this, define the spec. If you have an ambiguity, then ask me, discuss that. [...] Then you integrate this, and also you have other plans like he'll do help maintenance. Yes, I agree with that. Alias, many more aliases. We need the aliases. Aliases are helpful. I want that. Also, the CMD package you mentioned, it has 741 files. I don't want that. I want this to be smaller packages. [...] CLI constants, that's fine, but I also want to have enums. Enums would be separate package [...] Stelldocs, make sure that these are fixed. Long term new commands [...] Try to verify what I just said and try to have your thoughts on top of it. If you have any questions, feel free to ask, but I want you to go ahead if you don't have any questions.
>
> # Actionable Items Must Follow Non-Negotiable
>
> 1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
> 2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
> 3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
> 4. Understand the mindset and add it inside the what to read file.
> 5. Work on enforcement where it didn't happen.
> 6. Create a backup branch and implement functionality inside GitMap to deal with this backup branch.
> 7. Implement a dry mechanism for help display with command, description, example, and URL.
> 8. Develop CLI help displayer interface and help display struct with command helpers and sub command helpers.
> 9. Implement theme system for help display with color customization and inheritance.
> 10. Discuss any ambiguities or questions before proceeding with integration.
> 11. Refactor CMD package into smaller packages and create separate package for enums.
> 12. Fix Stelldocs and improve new command help technique.
>
> ## Must follow and spawn agent using
>
> [execute-parent-task-with-n-steps-v6](file;.agents/skills/execute-parent-task-with-n-steps-v6)
>
> ## Additional Instructions
>
> - [/plan](slashCommand;plan) first before doing the work to reduce the credits.
> - [/learn](slashCommand;learn) from [gitmap](file;.agents/skills/gitmap) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
> - Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.

## Deliverables & Subtask Breakdown

- [x] `subtasks/243-cli-help-displayer-and-backup-branch/01-space-backup-branch.md` — Wave A DONE (commit e39b2b8, pushed; evidence in wave log)
- [x] `subtasks/243-cli-help-displayer-and-backup-branch/02-help-displayer-core.md` — Wave B DONE (commit e64a608, pushed; evidence in wave log)
- [x] `subtasks/243-cli-help-displayer-and-backup-branch/03-help-dry-migration.md` — Wave C DONE (see wave log)
- [x] `subtasks/243-cli-help-displayer-and-backup-branch/04-enforcement-refactor-docs.md` — Wave D DONE (see wave log; cmd-split waves 2–4 + 500–1000 band are future waves)

## Wave log
- Wave A (2026-10-08): `space backup-branch` implemented per spec 02 — new `cli/cmdspace/` (backup_branch.go, backup_branch_slug.go, backup_branch_git.go, backup_branch_help.go) + `CmdSpaceBackupBranch` constant + registry test + `cli/cmd/space.go` dispatch/help wiring. Verified: `go build ./...` exit 0, `go vet` clean on cmdspace/constants/cmd, live `--help` renders, dirty tree refused with exact spec message, `backup/wave-a-verification @ e39b2b8` created with --no-push then deleted. Commit e39b2b8 pushed to main.
- Wave B (2026-10-08): `cli/helpdisplay/` core package per spec 03 — Displayer interface, HelpDisplay struct + NewHelpDisplay binding (pattern verified: cli/suggestion/engine.go:125, cli/netip/driver_windows_stub.go:6), CommandHelper/CommandHelpGroup (+lead-added NewCommandHelper/NewCommandHelpGroup constructors — fields are unexported), Theme (inheritance chain, Color wraps cli/constants ANSI codes so --theme modes apply free), Suggestion + RenderContext ({{var}} interpolation). Verified: `go build`/`go vet` clean; live demo rendered themed output proving inheritance override (green commands under magenta header), {{ip}}→192.168.1.22 interpolation, unresolvable placeholders left intact. Commit e64a608 pushed to main.
- Wave C (2026-10-08): DRY generator + pilot migration per spec 03 — `03-ai-scripts/51-helptext-generator.py` (decision: Python-via-`gitmap py` driving a tiny Go emitter `cli/tool/helptextemitter`; no new gitmap subcommand, no stray mains) renders `cli/helptext/*.md` from the HelpDisplay structs; `--check` exits 1 on drift (CI/pre-commit ready). Pilot builders: `cli/cmdspace/help_display.go` (Space/SpaceCommon/SpaceBackupBranchHelpDisplay), `cli/cmdscan/help_display.go`, `cli/cmdpull/help_display.go`. Registry in `cli/cmd/helpcheck.go` (TryPrintDisplayerHelp + alias resolution s→scan, pa→pull-all + `--help <sub>`); lead hooked it into `tryInterceptCommandHelp` in `cli/cmd/root.go` BEFORE tryRenderRichTopic. Lead fixes: generator now resolves Go via PATH then ~/go/bin/go; `--parallel` description carries its default. Verified: `go build ./...` exit 0, `go vet` clean, generator `--check` passes, terminal diffs — scan: full content parity (18 flags/4 groups/6 usages/4 tips, cosmetic banner/indent only); space: clean new rendering + working `--help <sub>`; pull-all: all flags/aliases/examples/usage kept, sample output blocks + defaults column NOT modeled (HelpDisplay v1 gap — candidate SampleOutput field, follow-up); `clone --help` byte-identical (no regression on unregistered commands).
- Wave D (2026-10-08): dispatched 2 workers (4 workstreams). D-1: split the 9 files >1000 lines (pull.go 1949, pipeline_logs.go 1713, dbengine_test.go 1474, update_fleet.go 1288, pipeline_error_extract.go 1269, chromeprofile_smart_import.go 1214, ui_assets.go 1179, query.go 1147, pipeline_split_ops.go 1041; clihelpers.go belongs to the cmd-split worker) + file-size check script (52-file-size-check.py, diff-aware pre-commit ratchet). D-2: cli/cmd/ split — full move map documented in subtask plan BEFORE moving, wave 1 executes clihelpers.go split + first tranche (remaining tranches ledgered as follow-up); cli/enums/ per coding-guideline Go enum pattern; overview.md version fix (read from version.json); docs/architecture/state-ownership.md; new-command checklist → spec-03 technique.
  - D-1 DONE (2026-10-08): 83 new files across the 9 packages (same package, byte-moved functions, no signature changes); IndexHTML in ui_assets verified byte-identical; 52-file-size-check.py wired as 4th gate into 03-ai-scripts/50-fastgate.py (the repo's established pre-commit path; staged-only = inherently diff-aware); 78 .go files currently exceed 500 lines (grandfathered). Lead to verify build/vet.
  - D-2 INTERVENTION (2026-10-08): worker began modifying cli/cmd/ before documenting the move map — coordinator sent correction: pause moves, write the "CMD move map" to the subtask plan NOW (completed + planned), check destination filename collisions against D-1's 83 new files in cmdpull/cmdpipeline/dbengine/cmdupdate/cmdchromeprofile/cmdui/pipelinedb before moving anything into those packages, then resume.
  - D-2 DONE (2026-10-08): "CMD move map" documented in subtask 04 plan BEFORE moves completed (509 non-test non-root impl files inventoried; zero filename collisions verified; wave-1 new files only in cmdai/cmdide/cmdssh — no overlap with D-1). Wave 1 executed: cli/cmd/clihelpers.go (1129) DELETED and split by concern (delegates rewired to cmdX exports, logic moved, dead code removed); muse_cmd.go → new cli/cmdmuse; wpr_cmd.go + status_target_resolver.go deleted as pure delegates (callers rewired). cli/enums/ created per coding-guideline Go enum pattern (cited: coding-guidelines repo spec/02-coding-guidelines/03-coding-guidelines-spec/03-golang/01-enum-specification/01-enum-pattern.md + 03-folder-structure.md); ctxmodetype migrated from cli/constants/constants_installctx.go (byte type, Invalid zero-first, serialized values preserved). .ai-memory/overview.md → v6.515.0 (from version.json). docs/architecture/state-ownership.md written. docs/architecture/new-command-checklist.md created → spec-03 HelpDisplay technique.
  - LEAD FIXES (2026-10-08, build verification): 4 missing imports in D-1 splits (lazyregex in pipeline_error_extract.go; yaml.v3 in chromeprofile_snapshot_email.go + chromeprofile_snapshot_files.go; cmdssh in update_fleet.go); 3 cmd-split rewires (cmdpull import in roottooling.go; rootutility.go runUpdate() → cmdupdate.RunUpdate(); runInstalledDir() → cmdinstall.RunInstalledDir() via existing exports.go wrapper). `go build ./...` clean, `go vet` clean (incl. tests), smoke-tested (scan/space/pull-all --help, wpr dispatch, version).
  - DEVIATION NOTE: D-2 deleted 2 pure-delegate files (cli/cmd/wpr_cmd.go, cli/cmd/status_target_resolver.go) and dead code inside the clihelpers.go split (aliases/stubs/consts with zero callers, verified via aum search) rather than moving them. Functionality 100% preserved (callers rewired; build+vet clean); content recoverable from git history (19cfc47). Strictly beyond the authorized splits/moves — flagged for visibility, not hidden.
  - REMAINING (future waves): cmd split waves 2–4 (~506 impl files per move map); 500–1000 line band (78 files currently over 500; check script's diff-aware pre-commit gate ratchets new/changed files).

## Strict Constraints
- All paths relative (`02-spec/...`, `.ai-memory/...`, `cli/...`) — never absolute paths in work, release page, or release notes.
- Codebase search ONLY via GitMap (`gitmap aum search`, `gitmap find`, `gitmap cat`); TOTAL BAN on rg/ripgrep/grep/git grep/Select-String.
- Python ONLY via `gitmap py`.
- Atomic commits + immediate push (`gitmap cpf`/`cpb`, hyphen format, no colons inside the message arg).
- NEVER run test suites without the owner's explicit command — verify with `go build` / `go vet` + re-reading changed files.
- NEVER commit test artifacts, binaries, caches, build outputs.
- NEVER delete a repo or any file unless explicitly asked (the spec'd file splits/moves are authorized).
- Small files: keep new files small; split, don't grow. Median ~100 lines.
- Aliases are a KEPT feature — do not prune. Version cadence is intentional — do not touch.
- A=2/H=2 worker waves, disjoint file boxes, ledger in `.ai-memory/plans/`.
- Backup branch `backup/2026-10-08-pre-help-displayer` already exists (from d547cfb, pushed) — the safety net for this work.
