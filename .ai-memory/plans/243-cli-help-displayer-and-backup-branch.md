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
- [ ] `subtasks/243-cli-help-displayer-and-backup-branch/02-help-displayer-core.md`
- [ ] `subtasks/243-cli-help-displayer-and-backup-branch/03-help-dry-migration.md`
- [ ] `subtasks/243-cli-help-displayer-and-backup-branch/04-enforcement-refactor-docs.md`

## Wave log
- Wave A (2026-10-08): `space backup-branch` implemented per spec 02 — new `cli/cmdspace/` (backup_branch.go, backup_branch_slug.go, backup_branch_git.go, backup_branch_help.go) + `CmdSpaceBackupBranch` constant + registry test + `cli/cmd/space.go` dispatch/help wiring. Verified: `go build ./...` exit 0, `go vet` clean on cmdspace/constants/cmd, live `--help` renders, dirty tree refused with exact spec message, `backup/wave-a-verification @ e39b2b8` created with --no-push then deleted. Commit e39b2b8 pushed to main.

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
