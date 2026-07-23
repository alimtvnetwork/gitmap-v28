# Post-Install Shell Activation — Contract

> **Parent spec:** [../21-post-install-shell-activation.md](../21-post-install-shell-activation.md)
> **Sibling files:**
> - [02-snippets.md](02-snippets.md) — Per-shell profile snippet bodies
> - [03-doctor.md](03-doctor.md) — `doctor` wrapper status detection
> - [04-idempotency.md](04-idempotency.md) — Rewrite, removal, and version-bump rules

## Purpose

After `setup` (or the bootstrap installer) runs, the user must be able
to invoke the CLI **and any of its shell-integrated subcommands** in
the **current terminal session** without restarting it. This contract
defines the generic, cross-platform behaviour every CLI must follow so
that:

1. The installer/setup step writes a profile snippet that exports
   `PATH` and shell wrappers (e.g. `cd`, `go`) for both Windows and
   Unix shells.
2. The setup step **auto-sources** the user's profile in the current
   session whenever it can.
3. When auto-source is impossible (Windows install in a fresh PS host,
   non-interactive shell, etc.), the CLI prints a deterministic
   one-liner the user can paste to activate the session.
4. The CLI exposes a runtime check (e.g. `<tool> doctor`) that detects
   "binary is on PATH but the wrapper is not loaded" and prints the
   exact reload command.

This contract eliminates the "PATH not active after install" and
"wrapper subcommand silently a no-op" classes of bugs.

---

## Required Behaviours

| ID | Behaviour | Required For |
|----|-----------|--------------|
| PIA-1 | `setup` writes shell snippet to user's profile, idempotent via marker comment. | All shells |
| PIA-2 | `setup` exports a shell-detection env var (e.g. `<TOOL>_WRAPPER=1`) so the binary can tell if the wrapper is active. | All shells |
| PIA-3 | `setup` attempts in-process activation: dot-source `$PROFILE` (PowerShell) or `source ~/.bashrc` / `~/.zshrc` (Bash/Zsh). | Interactive shells |
| PIA-4 | When PIA-3 cannot run (different parent shell, non-interactive, Windows installer host), `setup` prints the **exact** reload one-liner for the detected shell. | All shells |
| PIA-5 | `doctor` reports wrapper status with one of three outcomes: `LOADED`, `INSTALLED_BUT_NOT_LOADED`, `NOT_INSTALLED`. | All shells |
| PIA-6 | Shell-dependent subcommands (anything that must change the parent shell state) print a stderr warning when invoked without `<TOOL>_WRAPPER=1`. | All shells |
| PIA-7 | Profile snippet first-line marker uses the format `# <tool> shell wrapper v<N>` so future versions can rewrite it deterministically. | All shells |

---

## Activation Flow

```
install / setup
      │
      ▼
detect shell (PowerShell / Bash / Zsh / Fish)
      │
      ▼
resolve profile path
      │
      ▼
inject snippet (idempotent via marker)
      │
      ▼
try in-session activation ───── success ──► print "Active in this session"
      │ failure
      ▼
print exact reload one-liner   (e.g. `. $PROFILE` / `source ~/.zshrc`)
      │
      ▼
print fallback: "Or open a new terminal window."
```

---

## In-Session Activation

`setup` MUST try to activate the wrapper in the **current** shell
before falling back to a printed instruction.

### PowerShell

If `setup` is invoked from a PowerShell session, it dot-sources
`$PROFILE` in-process:

```powershell
. $PROFILE
if ($env:TOOLNAME_WRAPPER -eq "1") {
    Write-Host "  ✓ Wrapper active in this session" -ForegroundColor Green
}
```

If the parent host is not PowerShell (e.g. user ran the `.exe`
installer from `cmd.exe` or File Explorer), in-session activation is
skipped and the printed one-liner path is taken.

### Bash / Zsh / Fish

The CLI **cannot** source the profile of its parent shell from a
child process. Instead, it detects the parent shell via `$SHELL` and
prints the exact one-liner:

```
  To start using toolname right now, run:

      source ~/.zshrc

  Or open a new terminal window.
```

The `source ~/.<rc>` line MUST match the active profile (`~/.bashrc`,
`~/.zshrc`, or `~/.config/fish/config.fish`).

---

## Shell Detection Rules

| Detection Source | Used For |
|------------------|----------|
| `$env:PSVersionTable` exists | PowerShell |
| `$ZSH_VERSION` set | Zsh |
| `$BASH_VERSION` set | Bash |
| `$FISH_VERSION` set | Fish |
| Fallback: `basename $SHELL` | Bash/Zsh on Linux/macOS |
| Fallback: `$ComSpec` | cmd.exe (no wrapper supported — print install instruction only) |

If the shell cannot be detected, `setup` prints the snippet for
**both** Bash and PowerShell and asks the user to paste the matching
block into their profile.

---

## Stderr Warning From Shell-Dependent Subcommands

Any subcommand that requires the wrapper (typically anything that
would change the parent shell's CWD or env) MUST detect missing
wrapper state and print a stderr warning, then continue with reduced
behaviour where possible:

```
  ⚠ Shell wrapper not active. The current command will print the path
    instead of changing directory. Run `toolname setup` (and reload
    your shell) to enable shell-integrated behaviour.
```

The warning text MUST include both:
1. The action the user should run (`toolname setup`).
2. The reload step required after setup (`. $PROFILE`, `source ~/.<rc>`).

---

## Constraints

- Snippets MUST be ASCII only — no em-dashes, no Unicode arrows. The
  PowerShell parser fails on UTF-8 in some hosts.
- The CLI MUST NOT modify any line outside its marker block.
- The CLI MUST NOT depend on the user editing their profile manually.
- Snippet body MUST be small (under 10 lines) so users can audit it.
- Reload instructions MUST be a single copy-pasteable line.
