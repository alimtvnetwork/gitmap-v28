# Master Plan: update-antigravity-ide

## User Request (Verbatim)

```text
can you please update antigravity IDE and is there


# Update Antigravity IDE: high priority instruction, non-negotiable task

Hi there. Can you please update the Antigravity IDE? And also please confirm if we have a command using `gitmap` to update the Antigravity IDE using the Ubuntu or other machine, and also at the same time, try to check in the scripts fixture, do we have the update option for the Antigravity? Okay. Please confirm this and then let me know. So usually we should have a command in the `gitmap`. We should update the Antigravity IDE, and also inside the Antigravity tools, we should have a command to update the IDE if it is not updated, for example, for Linux. So work on it, and then finally update and release for `gitmap`, scripts fixture, and Antigravity tools, and also update the IDE here. Do you understand?

## slug: update-antigravity-ide

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/update-antigravity-ide/ and enqueue plan task in .ai-memory/plans/update-antigravity-ide.md (subtasks in .ai-memory/plans/subtasks/update-antigravity-ide/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Use `gitmap` AI agents to enter data
5. Task completion includes committing and pushing to Git
6. Confirm if there is a command using `gitmap` to update the Antigravity IDE on Ubuntu or other machines
7. Check in the scripts fixture for an update option for the Antigravity IDE
8. Ensure there is a command in the Antigravity tools to update the IDE for Linux if not already present
9. Update and release for `gitmap`, scripts fixture, and Antigravity tools, and update the IDE

## Must follow and spawn agent using

@[\1](file;.agents/skills/\1)

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[\1](file;.agents/skills/\1) skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

---

## 1. Executive Summary & Objective

The primary objective is to implement end-to-end Antigravity IDE update capabilities across the developer ecosystem:
1. **GitMap CLI**: Confirm whether a command exists (verified: currently missing) and add first-class support for `gitmap update antigravity` / `gitmap update ide` (and `gitmap agy update ide`).
2. **Scripts Fixture (`scripts-fixer`)**: Inspect the fixture repository (verified: missing `update` action in `scripts-linux/69-install-antigravity/run.sh`) and add an explicit `update` action/option.
3. **Antigravity Tools (`antigravity-manager`)**: Ensure a robust Linux update command (`agm agy update-ide` or `scripts/update-antigravity-ide.sh`) is provided.
4. **Live System Update**: Execute update on the host Ubuntu machine, verifying launcher symlinks and binaries in `~/.local/share/antigravity-ide/` and `~/.local/bin/antigravity`.
5. **Release & Sync**: Update versions and changelogs across all affected repositories (`gitmap`, `scripts-fixer`, `antigravity-manager`), committing and pushing atomically via GitMap.

---

## 2. Traceable Deliverables & Task Breakdown

| Task ID | Component / Area | Status | Target Files |
| :--- | :--- | :--- | :--- |
| **Task-01** | Specifications & Plan | IN_PROGRESS | `02-spec/21-app/update-antigravity-ide/`, `.ai-memory/plans/` |
| **Task-02** | GitMap IDE Update Command | QUEUED | `cli/cmd/rootutility_update.go`, `cli/cmd/rootutility_update_args.go`, `cli/cmdinstall/agy_install.go`, `cli/cmdinstall/installantigravity.go` |
| **Task-03** | Scripts Fixture Update Option | QUEUED | `git-work/scripts-fixer/scripts-linux/69-install-antigravity/run.sh`, `manifest.json`, `config.json` |
| **Task-04** | Antigravity Tools Update Command | QUEUED | `git-work/antigravity-manager/scripts/update-antigravity-ide.sh`, `package.json` |
| **Task-05** | Live Ubuntu IDE Update | QUEUED | `~/.local/share/antigravity-ide/`, `~/.local/bin/antigravity` |
| **Task-06** | Version Bumps, Release & GitMap Commits | QUEUED | `version.json`, `changelog.md` across repositories |

---

## 3. Disjoint Subtask Allocation

- **Subtask 01 (`01-gitmap-ide-update.md`)**: Implementation of `gitmap update antigravity` and `gitmap update ide` routing to `cmdinstall.RunInstallAntigravityWithOpts(cmdinstall.InstallOptions{Force: true})`.
- **Subtask 02 (`02-antigravity-manager-tools.md`)**: Implementation of Antigravity tools update script `scripts/update-antigravity-ide.sh` in `antigravity-manager`.
- **Subtask 03 (`03-scripts-fixer-update.md`)**: Implementation of `update` action in `scripts-fixer` (`scripts-linux/69-install-antigravity/run.sh`).
- **Subtask 04 (`04-acceptance-and-release.md`)**: Verification of live IDE update, version manifest bumps, changelog entries, and atomic push via GitMap.

---

## 4. Operational Invariants & Rules

- **R1 Zero Builds / Full Test Suites**: Run only targeted syntax and unit checks on touched files.
- **R3 Concrete Evidence**: Record diffstats, exit codes, and file existence.
- **R5 Multi-Agent Spawning**: A = 2 agents executed via `invoke_subagent` at each phase.
- **R6 Disjoint File Bounding Boxes**: Each agent operates strictly in non-overlapping files.
- **R8/R9 Atomic GitMap Commit**: Single atomic commit via `gitmap cpf "<module> - <summary>"` per modified repository.
- **R11 Strictly Relative Git Paths**: Zero absolute paths (`/home/a/...` or `C:\...`) in code, specs, commits, or release notes.
