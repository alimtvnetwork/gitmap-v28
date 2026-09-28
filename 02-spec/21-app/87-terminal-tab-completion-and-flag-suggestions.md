# Canonical Spec: Terminal Tab Completion, AGY Flag Suggestions & Shell Integration

**Spec ID:** SPEC-APP-87  
**Status:** Active  
**Author:** Antigravity Master Orchestrator  
**Date:** 2026-09-28  
**Version:** v6.365.1  

---

## 1. User Request (Verbatim)

```text
Hi. Look, all the flags that we have in many places. Flag, especially the AGY commands. Okay. I want all these to be in suggestion in terminal, regardless of the Windows or Ubuntu. Can you please work on it and make sure that all the things actually come as a help text in the terminal and also a suggestion, like if I do a Tab, it just completes. Some places we can type file path or things like that, then it will start providing these suggestions. It will also go for the machine aliasing. Okay, if we go into SSH, all this command needs to have proper suggestions. Also, for in the future, let's say we typed a command, used it, so in future, when we try to reuse it, it will try to provide a suggestion. Can you please work on it? Is it possible for you to do it?
```

---

## 2. Technical Architecture & Invariants

1. **Cross-Platform Shell Tab Completion (Windows PowerShell + Ubuntu Bash/Zsh):**
   - Native Cobra `__complete` dynamic completion engine must be fully wired for root commands, subcommands, and flags.
   - Command `gitmap completion install` (alias `gitmap completion in`, `gitmap setup completion`) automatically detects the current shell (PowerShell on Windows, Bash/Zsh on Linux/macOS) and registers completion script hooks in:
     - PowerShell: `$PROFILE` (incorporating `Register-ArgumentCompleter -Native -CommandName gitmap` and PSReadLine predictive IntelliSense).
     - Bash: `~/.bashrc` or `/etc/bash_completion.d/gitmap`.
     - Zsh: `~/.zshrc` or `~/.zfunc/_gitmap`.
2. **Comprehensive AGY Subcommands & Flags Completion:**
   - All flags across `gitmap agy` and root AGY aliases (`sug`, `shutdown-until`, `rerun`, `rr`, `running-prompts`, `rp`, `fpug`, `pins`) must feature informative Cobra help annotations and complete on Tab.
   - Value completions:
     - `--model`: `flash_lite`, `flash`, `pro`.
     - `--time` / `-t`: durations (`1m`, `2m`, `5m`, `10m`).
     - `rerun [target]`: sequence numbers `1..N`, project names.
     - `sug [subcommand]`: `ls`, `run`, `add-projects`, `rm`, `agy-running-projects`, `help`.
3. **Context-Aware Dynamic File Path & SSH Node Completion:**
   - File path flags (e.g. `--file`, `-f`, `--path`, `--dest`, `--config`) must use Cobra's file completion directives (`cobra.ShellCompDirectiveDefault` / `cobra.MarkFlagFilename`).
   - SSH and Cluster commands (`gitmap ssh`, `gitmap deploy`, `gitmap deploy-right`, `gitmap deploy-left`, `gitmap update --remote`) must dynamically suggest registered cluster node aliases (`w1`, `w2`, `alpha-win`), IP addresses, and database host IDs from SQLite.
4. **Command Execution History & Future Reuse Suggestions:**
   - Persist executed commands with timestamp, frequency, and exit status into a dedicated SQLite Split-DB table (`CommandHistoryRecord`).
   - Expose recent command suggestions when typing `gitmap` without arguments or requesting completions.
   - Integrate with PowerShell PSReadLine History Prediction to automatically suggest previously executed commands inline as ghost text.
