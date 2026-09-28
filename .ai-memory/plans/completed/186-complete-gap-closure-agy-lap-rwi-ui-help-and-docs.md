# Consolidated Plan 186: Complete Gap Closure — AGY LAP/RP/RWI/RWC Shorthands, Fast-Forward Alias, UI Help & Command Docs

Spec Reference: [02-spec/21-app/90-agy-add-read-lap-rwi-machine-alias-telegram-and-os-help.md](../../../02-spec/21-app/90-agy-add-read-lap-rwi-machine-alias-telegram-and-os-help.md)
Status: COMPLETED (`v6.371.0`)
Loops Executed: 1

## Summary of Completed Subtasks

### Subtask 01: LAP/RP/RWI/RWC Flag Shorthands, Config Default Hours & Fast-Forward Alias
- Enabled `DisableFlagParsing: true` on `AgyRunningProjectsCmd` (`cli/cmdagy/agy_running_projects.go`) and implemented `parseRunningProjectsFlags` so `--wc`, `--wordcount`, `--ww`, `-ww`, `--json`, and `-f`/`--file` are preserved and forwarded intact to `RunRunningProjectsPromptsTreeCLI`.
- Updated `parseLapArguments` (`cli/cmdagy/agy_lap_and_rp_tree.go`) to read `lap.default_hours` from `config.GetVariable("global", "lap.default_hours")` (fallback `24`), and added support for `p2` / `P2` / `-p2` and `page <N>` page shorthands plus `--ww` / `-ww` / `-wordcount` wordcount flags.
- Enabled `DisableFlagParsing: true` on `AgyRerunWithIDCmd` and `AgyRerunWithConvIDCmd` (`cli/cmdagy/agy_rerun_with_id.go`) with `parseRwiFlexibleArgs` supporting `-p`, `--prompt`, and single-dash `-prompt` anywhere in the argument list and allowing `rwi`/`rwc`/`rwp` with `-p`/`-prompt` without requiring a redundant positional text argument.
- Added `"fast-forward"` and `"ff"` aliases to `AgyAccountSwitchCmd` (`cli/cmdagy/agy_account_switch_cmd.go`) and `normalizeRunningAndSwitchSubcommands` (`cli/cmdagy/agy_cmd.go`), and routed `gitmap agy machine` and `gitmap agy alias` in `tryDispatchAgyShortcut`.

### Subtask 02: Machine/Alias `y` Confirmation, UI Help & Speed Settings, and Command Docs
- Updated `filterPositionalMachineArgs` (`cli/cmdos/os_machine_alias.go`) to strip trailing confirmation tokens (`"y"`, `"yes"`, case-insensitive) when `len(out) >= 2` so `gitmap machine set my-node y` and `gitmap alias set my-alias y` apply automatically without treating `y` as the target name.
- Added the **Speed Settings** cards (`lap.default_hours`, `account_switch.threshold`, `telegram.bot_token` / `telegram.chat_id`, `email.smtp_host` / `email.from` / `email.to`, `machine.name` / `machine.alias`) and a dedicated **📖 CLI & AGY Help** tab (`#tab-help`) in `cli/cmdui/ui_assets.go`.
- Updated `cli/cmdagy/agy_help.go`, `cli/helptext/agy.md`, `cli/helptext/os.md`, `cli/helptext/alias.md`, and `cli/helptext/catalog.go` to document all new commands and aliases across CLI help, UI help, and command docs.
