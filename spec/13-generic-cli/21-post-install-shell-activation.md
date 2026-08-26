# Post-Install Shell Activation — Generic CLI Spec

> **Related specs:**
> - [11-build-deploy.md](11-build-deploy.md) — install/deploy step that places the binary on PATH
> - [19-shell-completion.md](19-shell-completion.md) — completion install uses the same profile-injection pattern
> - [13-checklist.md](13-checklist.md) — implementation phases that include setup
> - Implementation reference: [06-version-and-help.md](../12-cicd-pipeline-workflows/06-version-and-help.md) — shell-integrated commands and help output
> - Historical sibling-app issue references live outside this repo; this spec keeps the activation contract local

## Purpose

After `setup` (or the bootstrap installer) runs, the user must be able
to invoke the CLI **and any of its shell-integrated subcommands** in
the **current terminal session** without restarting it. This spec
defines the generic, cross-platform contract every CLI must follow so
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

## Profile Snippet Contract

Every snippet MUST start with the marker line and MUST end with a
matching closing marker so the CLI can rewrite or remove it without
disturbing surrounding content:

```

# <tool> shell wrapper v2 — managed by `<tool> setup`. Do not edit manually.

...snippet body...

# <tool> shell wrapper v2 end

```

### PowerShell (`$PROFILE`)

```powershell

# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.

$env:TOOLNAME_WRAPPER = "1"
function gcd { Set-Location (toolname cd @args) }

# toolname shell wrapper v2 end

```

### Bash / Zsh (`~/.bashrc`, `~/.zshrc`)

```bash

# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.

export TOOLNAME_WRAPPER=1
gcd() { cd "$(toolname cd "$@")" ; }

# toolname shell wrapper v2 end

```

### Fish (`~/.config/fish/config.fish`)

```fish

# toolname shell wrapper v2 — managed by `toolname setup`. Do not edit manually.

set -gx TOOLNAME_WRAPPER 1
function gcd; cd (toolname cd $argv); end

# toolname shell wrapper v2 end

```

The detection variable name MUST follow `<TOOL>_WRAPPER` (uppercased,
underscores) so multiple CLIs can coexist in one profile.

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

### Bash / Zsh

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

## `doctor` Wrapper Check

`doctor` MUST emit one of these three outcomes:

| Status | Stdout | Exit |
|--------|--------|------|
| LOADED | `[OK] Shell wrapper active (TOOLNAME_WRAPPER=1)` | 0 |
| INSTALLED_BUT_NOT_LOADED | `[!!] Shell wrapper installed but not loaded — run: source ~/.zshrc` | 1 |
| NOT_INSTALLED | `[!!] Shell wrapper missing — run: toolname setup` | 1 |

Detection algorithm:

1. Read `<TOOL>_WRAPPER` from environment → if `"1"` → **LOADED**.
2. Read profile file (`$PROFILE` / `~/.bashrc` / `~/.zshrc` /
   `~/.config/fish/config.fish`) and search for the marker line
   `# <tool> shell wrapper v<N>`. If found → **INSTALLED_BUT_NOT_LOADED**.
3. Otherwise → **NOT_INSTALLED**.

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

## Idempotency Rules

- Profile rewrites MUST be safe to run repeatedly (`setup`, `update`,
  CI provisioning).
- The marker comment is the **only** legal anchor for rewrites. Tools
  MUST NOT use line counts, absolute offsets, or content matching of
  the wrapper body.
- Bumping the version (`v2` → `v3`) MUST first remove every previous
  marker block before injecting the new one.
- Removing the wrapper (`<tool> uninstall --shell-wrapper`) MUST
  delete the marker block in full and leave surrounding content
  byte-identical.

---

## Cross-Platform Parity Table

| Capability | PowerShell | Bash | Zsh | Fish |
|------------|-----------|------|-----|------|
| Profile detection | ✅ `$PROFILE` | ✅ `~/.bashrc` | ✅ `~/.zshrc` | ✅ `~/.config/fish/config.fish` |
| Marker-based snippet | ✅ | ✅ | ✅ | ✅ |
| Wrapper detection env var | ✅ `$env:` | ✅ `export` | ✅ `export` | ✅ `set -gx` |
| In-session activation | ✅ dot-source | ❌ (print one-liner) | ❌ (print one-liner) | ❌ (print one-liner) |
| `doctor` LOADED check | ✅ | ✅ | ✅ | ✅ |
| `doctor` INSTALLED_BUT_NOT_LOADED | ✅ | ✅ | ✅ | ✅ |
| Reload one-liner printed | ✅ `. $PROFILE` | ✅ `source ~/.bashrc` | ✅ `source ~/.zshrc` | ✅ `source ~/.config/fish/config.fish` |

---

## Implementation Checklist For New CLIs

1. Add `<TOOL>_WRAPPER` constant to the constants package.
2. Define `ShellWrapperMarkerPrefix` and `ShellWrapperMarkerSuffix`
   constants (e.g. `# toolname shell wrapper v2`).
3. Implement `setup/wrapper.go` with:
   - `DetectShell()`
   - `ResolveProfilePath(shell)`
   - `InjectSnippet(profilePath, shell)` (idempotent via marker)
   - `RemoveSnippet(profilePath)` for upgrades
   - `TryInSessionActivate(shell)` (PowerShell only succeeds today)
   - `PrintReloadInstruction(shell)`
4. Add `doctor` check `checkShellWrapper()` returning one of the three
   statuses above.
5. Add stderr warnings to every shell-dependent subcommand using
   `os.Getenv("<TOOL>_WRAPPER")`.
6. Cover with tests:
   - Snippet injection on a fresh profile.
   - Snippet re-injection (no duplicates).
   - Marker-based removal.
   - `doctor` returns the correct status for each of the three states.

---

## Constraints

- Snippets MUST be ASCII only — no em-dashes, no Unicode arrows. The
  PowerShell parser fails on UTF-8 in some hosts.
- The CLI MUST NOT modify any line outside its marker block.
- The CLI MUST NOT depend on the user editing their profile manually.
- Snippet body MUST be small (under 10 lines) so users can audit it.
- Reload instructions MUST be a single copy-pasteable line.
- Tests MUST cover injection, idempotency, removal, and `doctor`
  status detection (per [12-testing.md](12-testing.md)).

---

## Why This Spec Exists

These bugs in the gitmap project triggered this generic spec:

| Issue | Root Cause |
|-------|------------|
| `22-installer-path-not-active-after-install` | Installer wrote to PATH but never told the user to reload, and never auto-activated. |
| `24-cd-command-does-not-change-shell-directory` | The `cd` subcommand silently no-op'd because the wrapper function was not loaded. |
| `25-powershell-cd-wrapper-not-loaded` | Same as 24 but on Windows — wrapper installed in `$PROFILE` but the running session never sourced it. |

By following the contract above, every new CLI in this framework
inherits a deterministic, AI-implementable post-install activation
flow on day one.

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
