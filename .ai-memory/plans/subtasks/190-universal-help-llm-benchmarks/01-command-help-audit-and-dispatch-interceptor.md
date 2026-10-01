# Subtask 190.1: Universal Command Help Interception & Command Audit

> **Parent Plan:** [Plan 190](../../completed/190-universal-command-help-modernization-llm-train-loop-and-polyglot-benchmarks.md)
> **Status:** Complete
> **Lead Architect:** MD ALIM UL KARIM

---

## 1. Objectives

1. Intercept `help`, `--help`, and `-h` universally in `cli/cmd/root.go` and `cli/cmd/helpcheck.go` for all registered commands before command-specific handlers parse positional arguments.
2. Specifically fix commands that currently fail or execute live logic when passed `help`:
   - `commit`: Currently gives fuzzy suggestions or errors on `help`.
   - `stash`: Currently attempts to execute stash operations.
   - `branch`: Currently errors with `unknown branch subcommand: help`.
   - `aum`: Currently launches the background AI server.
   - `templates`: Currently errors with `unknown 'templates' subcommand: help`.
   - `ui`: Currently launches the web server on port 8080.
   - `revert`: Currently attempts to parse `help` as a commit or release target.
   - `user`: Currently errors on `help`.
3. Verify that every command in the CLI routing table handles `<cmd> help` cleanly without errors.

---

## 2. Implementation Steps

1. In `cli/cmd/helpcheck.go`, export a helper `IsHelpTrigger(token string) bool` matching `help`, `--help`, and `-h`.
2. In `cli/cmd/root.go`, add an early pre-dispatch check for commands:
   If `len(os.Args) >= 3 && IsHelpTrigger(os.Args[2])`, verify if the command has a dedicated help handler or route to `printHelpAndExit(command, os.Args[2:])`.
3. In `cli/cmd/rootutility.go`, ensure `dispatchHelpTopic` has mappings for all aliases and commands.
4. In `cli/cmdaum/`, `cli/cmdtemplates/`, `cli/cmdui/`, `cli/cmdbranch/`, `cli/cmdstash/`, add explicit `checkHelp(cmd, args)` calls at the start of their entry points.
