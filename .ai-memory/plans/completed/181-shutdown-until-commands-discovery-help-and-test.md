# Completed Plan 181: Shutdown-Until Command Discovery, Root Dispatch & Rich Help Integration

**Spec Reference:** [02-spec/21-app/86-shutdown-until-commands-discovery-help-and-test.md](../../../02-spec/21-app/86-shutdown-until-commands-discovery-help-and-test.md)  
**Execution Date:** 2026-09-28  
**Release Target:** v6.365.0  

---

## 1. Verified Outcomes & Achievements

- [x] **Subtask 01: Root Dispatch and Top-Level Aliases**
  - Registered `gitmap shutdown-until`, `gitmap shutdown-until-green`, and `gitmap sug` in `cli/cmd/rootutility.go` (`utilitySystemEntries`) and `cli/cmd/root.go` (`dispatchAgySubsystem`).
  - Added `CmdShutdownUntil = "shutdown-until"`, `CmdShutdownUntilAlias = "sug"`, `CmdShutdownUntilGreen = "shutdown-until-green"` to `cli/constants/constants_cli.go`.
  - Added `shutdown-until` alias to `AgySUGCmd` and `normalizeWorkflowSubcommands` in `cli/cmdagy/agy_cmd.go`.

- [x] **Subtask 02: Rich Terminal Help Menu & Documentation**
  - Created `cli/cmdagy/agy_sug_help.go` implementing `RenderAgySugHelp()` with a cyan/yellow two-column boxed terminal layout (`termhelp.HelpMenu`).
  - Connected `gitmap sug help`, `gitmap shutdown-until --help`, `gitmap help shutdown-until`, and Cobra help handlers directly to `RenderAgySugHelp()`.
  - Promoted `shutdown-until (sug)`, `finish-prompts-until-green (fpug)`, and `running-prompts (rp)` to the very top of `Protocols & Automation` in `gitmap agy help` (`cli/cmdagy/agy_help_automation.go`).
  - Integrated `HelpShutdownUntil` into `printGroupIntegrations()` in `cli/cmd/rootusage_groups.go` and `cli/constants/constants_helpgroups.go`.
  - Authored comprehensive documentation in `cli/helptext/shutdown-until.md` with synopsis, options table, subcommands, and real-world examples.
  - Added unit test suite `cli/cmd/root_sug_help_test.go` verifying topic normalization and rich help rendering.

- [x] **Subtask 03: Safe Non-Destructive Live Testing of Watch Loop**
  - Added `--dry-run` (`-n`) and `--once` (`-1`) flag support in `cli/cmdagy/agy_sug.go` and `cli/cmdagy/agy_sug_executor.go`.
  - Implemented safe simulation in `triggerSystemShutdown`: when `--dry-run` is active, displays the exact OS command that would execute (`shutdown /s /t 60` or platform equivalent) without calling OS shutdown.
  - Added informative guidance when `DiscoverRunningProjects` finds no active queues instead of silently writing an empty config.
  - Created safety unit test suite `cli/cmdagy/agy_sug_executor_test.go` proving `OSShutdownExecutorFn` is never invoked under `--dry-run`.
  - Performed live testing on `gitmap shutdown-until --help`, `gitmap sug help`, `gitmap agy help`, `gitmap sug ls`, `gitmap sug agy-running-projects`, and `gitmap sug run --dry-run --once`.

---

## 2. Modified Files

- `02-spec/21-app/readme.md`
- `cli/constants/constants_cli.go`
- `cli/constants/constants_helpgroups.go`
- `cli/cmd/rootusage_groups.go`
- `cli/cmd/rootutility.go`
- `cli/cmd/root.go`
- `cli/cmdagy/agy_cmd.go`
- `cli/cmdagy/agy_sug.go`
- `cli/cmdagy/agy_sug_executor.go`
- `cli/cmdagy/agy_sug_help.go`
- `cli/cmdagy/agy_help_automation.go`
- `cli/cmd/rich_help_dispatcher.go`
- `cli/helptext/print.go`
- `cli/helptext/shutdown-until.md`
- `cli/cmd/root_sug_help_test.go`
- `cli/cmdagy/agy_sug_executor_test.go`
