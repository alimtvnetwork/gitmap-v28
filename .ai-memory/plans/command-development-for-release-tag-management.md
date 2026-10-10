# Plan: Command Development for Release Tag Management (`gitmap fix release tags`)

## User Request (Verbatim)
A, add a new command and also the help text for the command. The command should actually remove releases and release tags if that release tag does not have the proper binary or the release published, or at that time, if that has CI/CD, then the release tag will be removed. Can you please work on this as a command and also add that command to the new command list, and also add the example and help text for this, and also add this command information in the UI as well? Please bump the version, minor version, bump it and release it. And also at the same time, check and test it on top of the `gitmap`. The command format should be `gitmap` space fix space release tags. And that would also give a dry run if you want it to, so that it shows what will be the end result if you run it, and that would prompt it. But user can also do a hyphen Y or confirm, confirm to make sure that it works automatically.

## Slug
`command-development-for-release-tag-management`

## Architecture & Subtask Decomposition
- **Subtask 01 (`01-release-tags-audit-and-deletion-engine.md`)**:
  - Audit engine inspecting git tags, GitHub releases, uploaded binary assets, and CI/CD workflow status.
  - Safe 3-tier deletion executor (GitHub release, remote origin tag, local tag).
  - Safety rules: active binary version protection, latest healthy release immunity, CI in-progress grace window.
- **Subtask 02 (`02-cli-routing-help-ui-and-documentation.md`)**:
  - Command routing in `cli/cmdfix/fix_cmd.go` for `gitmap fix release tags`, `fix-release-tags`, and `frt`.
  - Flags support: `--dry-run`, `-n`, `-y`, `--yes`, `--confirm`, `--json`.
  - Terminal UI: Diagnostic table preview via `termout.PrintTable`, two-column help menu, interactive confirmation prompt.
  - Command list and documentation integration in `constants_cli.go`, `roothelp_clusters.go`, and `cli/helpdoc/fix-release-tags.md`.

## Quality Gates & Verification
1. Strictly relative Git paths across all specs, plans, and source files.
2. Unit tests in `cli/cmdfixreleasetags/` verifying detection, deletion orchestration, mock GH responses, and flag handling.
3. Live test of `gitmap fix release tags --dry-run` and `-h` on the `gitmap` repository.
4. Pass relative path and forbidden strings linters.
5. Minor version bump and release ceremony.
