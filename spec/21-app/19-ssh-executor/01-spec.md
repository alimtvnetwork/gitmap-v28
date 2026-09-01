# Spec: SSH Executor & Delegation (gitmap se)

## Goal

Implement a highly concurrent, context-aware SSH executor (`gitmap se`) capable of natively delegating core `gitmap` commands over SSH or falling back to platform-appropriate shells (`bash`/`ps`). It must seamlessly auto-install `gitmap` and `powershell` (on Unix) remotely if required.

## Actionable Non-Negotiable Items

- **Shell Defaulting:** `gitmap se "cmd"` MUST execute via `ps` if the target is Windows, and `bash` if Unix.
- **Native Delegation:** Commands prefixed with `mkdir`, `cat`, or `ssh` MUST delegate to the remote `gitmap` binary (i.e. `gitmap mkdir -p ...`).
- **Auto-Installation (Gitmap):** Before executing a command on the remote, it MUST verify if `gitmap` is available. If missing, it MUST install it using the cross-platform one-liners.
- **Auto-Installation (PowerShell):** If a Unix target is invoked with the `ps` shell explicitly, the executor MUST ensure `pwsh` is installed.
- **`gitmap reset-and-rescan`:** MUST trigger a full database reset via `--confirm` and immediately run a `--rescan`.
- **Function Constraints:** All new Go functions MUST adhere to the 15-line hard cap (8 preferred) and avoid nested `if` statements. No magic strings or raw booleans (follow `Result` envelope).

## Validation & Testing

- MUST include positive and negative tests for OS-based shell selection (`determineSSHCommand`).
- MUST include tests verifying early exits when arguments are empty.
- MUST test `parseSEFlags` successfully excluding machines by comma-separated values.
