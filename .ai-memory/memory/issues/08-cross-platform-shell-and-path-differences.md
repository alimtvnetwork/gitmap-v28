# Issue Domain 08: Cross-Platform Shell and Path Differences

- **Domain:** PowerShell, Bash, Zsh, and Line Ending Normalization
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Oh-My-Zsh Theme Escape Pollution
- **Symptoms:** Remote SSH JSON outputs contained ANSI color escape codes and update prompts.
- **Root Cause:** Non-interactive SSH sessions triggered interactive `.zshrc` profile scripts.
- **Resolution:** Wrapped remote execution commands in `sh -c` or set `TERM=dumb` to suppress interactive shell hooks.

## 2. CRLF Line Ending Breakage in Shell Scripts
- **Symptoms:** Bash scripts failed with `\r: command not found` on Linux VMs.
- **Root Cause:** Windows git checkouts converted line endings to CRLF.
- **Resolution:** Enforced `.gitattributes` setting `* text eol=lf` across all repository shell and source files.
