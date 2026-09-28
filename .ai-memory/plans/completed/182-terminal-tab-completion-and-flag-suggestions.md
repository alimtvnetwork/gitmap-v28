# Consolidated Plan 182: Terminal Tab Completion, AGY Flag Suggestions & Shell Integration

Spec Reference: [02-spec/21-app/87-terminal-tab-completion-and-flag-suggestions.md](../../../02-spec/21-app/87-terminal-tab-completion-and-flag-suggestions.md)
Execution Summary: Completed in 4 subtasks under N-step continuous loop.

## User Request (Verbatim)

```text
Hi. Look, all the flags that we have in many places. Flag, especially the AGY commands. Okay. I want all these to be in suggestion in terminal, regardless of the Windows or Ubuntu. Can you please work on it and make sure that all the things actually come as a help text in the terminal and also a suggestion, like if I do a Tab, it just completes. Some places we can type file path or things like that, then it will start providing these suggestions. It will also go for the machine aliasing. Okay, if we go into SSH, all this command needs to have proper suggestions. Also, for in the future, let's say we typed a command, used it, so in future, when we try to reuse it, it will try to provide a suggestion. Can you please work on it? Is it possible for you to do it?
```

---

## Consolidated Outcomes & Verifications

### 1. Dynamic Shell Tab-Completion Engine & Profile Auto-Installer
- Implemented Cobra-driven dynamic completion architecture with unified `GetRootCompletionCmd()` in `cli/cmd/root_cobra_completion.go`.
- Intercepted `gitmap __complete` and `gitmap __completeNoDesc` in `cli/cmd/root.go` for zero-overhead, clean stdout completion.
- Updated shell completion generators for PowerShell, Bash, and Zsh in `cli/completion/` (`powershell.go`, `bash.go`, `zsh.go`, `completion.go`).
- Implemented `gitmap completion install [shell]` in `cli/cmd/completion.go` and `cli/completion/install.go`, which registers dynamic tab completion and PSReadLine prediction into `$PROFILE` / `~/.bashrc` / `~/.zshrc`.
- Verified `gitmap completion powershell`, `gitmap completion bash`, `gitmap completion zsh`, and `gitmap completion install`.

### 2. AGY Subcommands and Flag Dynamic Completion Provider
- Configured Cobra flags, flag completion functions, and `ValidArgsFunction` for all Antigravity commands:
  - `agy rerun` / `gitmap rerun`: `--prompt`, `--model` (completes `flash_lite`, `flash`, `pro`), `--restart`, `--dry-run`, `--all`, `--queue`, project target numbers (`1`, `2`).
  - `agy shutdown-until-green` / `gitmap sug`: subcommands (`ls`, `run`, `add-projects`, `rm`, `agy-running-projects`, `help`) and flags (`--time`, `--dry-run`, `--once`).
  - `agy running-prompts`: subcommands (`ls`, `backup`, `restore`, `clean`, `export`, `import`, `help`) and flags (`--limit`, `--wc`, `--full`, `--json`, `--file`).
  - `agy running-projects`: subcommands (`ls`) and flags (`--json`, `--ssh`, `--file`).
  - `agy finish-prompts-until-green` / `gitmap fpug`: targets completion (`running-projects`).
- Tagged all file flags with `cobra.MarkFlagFilename`.
- Verified `gitmap __complete agy rerun ""` and `gitmap __complete agy rerun --model ""`.

### 3. File Path and SSH Fleet Node Dynamic Auto-Completion
- Implemented `GetFleetNodeCompletions()` in `cli/cmdssh/ssh_completion.go` querying registered cluster nodes from database.
- Implemented `DeployValidArgsFunction` for `gitmap deploy`, `deploy-right`, `deploy-left`:
  - Argument 0: cluster node aliases (`alpha-win`, `beta-linux`, `gamma-mac`, `w1`, `w2`, `w3`, `all`).
  - Argument 1: local source file/folder path with `cobra.ShellCompDirectiveDefault`.
  - Argument 2: remote destination directory suggestions.
- Registered fleet node completion on `gitmap update --remote <node>`.
- Verified `gitmap __complete deploy ""` and `gitmap __complete deploy w1 ""`.

### 4. Command History Logging and Shell Predictor Integration
- Created dedicated SQLite Split-DB in `cli/store/command_history_split_db.go` (`data/history/commands.db`).
- Integrated automatic execution recording in `finishCommandAudit` (`cli/cmd/audit_finish.go`) with command line, execution time, duration, and status.
- Implemented `gitmap history [ls] [--limit 20] [--json]`, `gitmap history suggest <prefix>`, and `gitmap history clear` in `cli/cmd/command_history_cli.go`.
- Added PSReadLine predictor integration in PowerShell profiles with `Set-PSReadLineOption -PredictionSource History` and `Set-PSReadLineOption -PredictionViewStyle ListView`.
- Verified `gitmap history ls` and `gitmap history suggest comp`.
