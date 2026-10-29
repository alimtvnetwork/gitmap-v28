# Plan: 347-gitmap-user-change-command

## Status: COMPLETED
## Slug: gitmap-user-change-command
## Category: Feature & CLI Ergonomics

## User Request (Verbatim)
# GitMap User Change Command: high priority instruction, non-negotiable task

Hi there. Can you please add a command for `gitmap` to change the user, let's say change user or change login that it want every one of them. Also, they could provide in one liner what would be the email, name, and also we should have a change to autobot. The autobot would work is that we can give an email address or no email address. It would say as AI automation bot as a name and it would put as automation agent. It's good. And also in the current machine, I want you to change all the Git name to AI automation bot and the email to local at local. Keep the email empty for now. Can you please do that? And make a bump in the `gitmap`, make sure that everything is released with it. Does this make sense?

## task: GitMap User Change Command
## slug: gitmap-user-change-command

# Actionable Items Must Follow Non-Negotiable
1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-gitmap-user-change-command/ and enqueue plan task in .ai-memory/plans/pending/XX-gitmap-user-change-command.md (subtasks in .ai-memory/plans/subtasks/XX-gitmap-user-change-command/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Add a command in `gitmap` to change the user, allowing input for email and name in a one-liner.
6. Implement a feature to switch to an autobot user with the name "AI automation bot" and email "local at local", leaving the email empty if not provided.
7. Change all Git user names on the current machine to "AI automation bot" and the email to "local at local".
8. Ensure a version bump in `gitmap` and release all changes accordingly.

## High-Level Execution Phases
- **Phase 1: Planning & Specifications**
  - Task-01: Codebase discovery across CLI routing, usercontext, and migration. (DONE)
  - Task-02: Author comprehensive architecture, component, and verification specs in `02-spec/21-app/347-gitmap-user-change-command/` and subtask plans in `.ai-memory/plans/subtasks/347-gitmap-user-change-command/`. (DONE)
- **Phase 2: Core Implementation**
  - Task-03: Implement `change-user`, `change-login`, `autobot` command handlers and root dispatcher routes in `cli/cmduser/` and `cli/cmd/`. (DONE)
  - Task-04: Integrate with `cli/usercontext/git_identity.go` for zero-no-op empty email handling and unsetting support. (DONE)
- **Phase 3: Machine Migration & Release**
  - Task-05: Update help menu documentation, migrate machine's Git identity globally to "AI automation bot" and "local@local". (DONE)
  - Task-06: Bump version to v6.628.0, synchronize changelogs, and execute `gitmap cpr` release. (DONE)

## Subtasks Manifest
1. `01-cli-commands.md`: CLI command routing and one-liner input parsing. (DONE)
2. `02-autobot-feature.md`: Autobot user identity presets and profile seed. (DONE)
3. `03-machine-migration-and-release.md`: Global machine git config switch, docs, and version bump release ceremony. (DONE)

## Verification Evidence
- `go test -C cli -count=1 ./cmduser/... ./usercontext/... ./constants/...`: PASS (All unit tests passing).
- `python 03-ai-scripts/50-fastgate.py`: PASS (0.026s).
- `python 03-ai-scripts/05-guideline-autofixer.py --check-only cli/cmduser`: PASS (Clean newlines & valid boolean conventions).
- `python 03-ai-scripts/05-guideline-autofixer.py --check-only cli/usercontext`: PASS (Clean newlines & valid boolean conventions).
- Global Git user.name migrated to `AI automation bot`.
- Global Git user.email migrated to empty / `local@local`.
