# PATH Not Available in Other Shells After Install

> **Superseded** — The multi-profile PATH registration introduced in
> [issue 22](./22-installer-path-not-active-after-install.md) resolves this
> problem. The installer now writes to all standard shell profiles, making
> this single-profile limitation obsolete.

## Ticket

After installing gitmap-v27 via `curl | sh` on macOS (zsh default), the binary
is not found when the user switches to `sh` or opens a non-zsh shell.

## Symptoms

1. User runs `curl -fsSL .../install.sh | sh` on macOS (zsh is default shell).
2. Installer adds `export PATH="$PATH:~/.local/bin"` to `~/.zshrc`.
3. gitmap-v27 works in zsh terminals.
4. User opens `sh` or switches to bash → `gitmap-v27: command not found`.

## Root Cause

The installer detects the user's **login shell** (`$SHELL`) and writes the
PATH entry to that shell's profile only (e.g. `~/.zshrc` for zsh). Other
shells (`sh`, `bash`, `fish`) read different profile files (`~/.profile`,
`~/.bashrc`, `~/.config/fish/config.fish`) and never see the PATH addition.

This is expected POSIX behavior — each shell has its own startup files —
but the installer did not clearly communicate this to the user.

## Fix

Updated the post-install summary in `install.sh` to:
1. Show which profile file was modified and for which shell.
2. Explicitly warn that other shells will NOT have gitmap.
3. Print the exact `export PATH=...` line the user needs to add to other
   shell profiles manually.

## Prevention

1. **Installer Output Contract** (in `.lovable/plan.md`): Unix installers
   must explicitly warn that other shells need manual PATH configuration.
2. All installer scripts (`.sh`, `.ps1`) must show version, path, and
   shell-specific PATH status in the post-install summary.
3. Never assume that adding PATH to one shell profile makes the binary
   globally available across all shells on the system.

## Related

- `spec/02-app-issues/19-missing-macos-binaries-and-lint-regression.md`
- `.lovable/plan.md` — Installer Output Contract guardrail
