# Installer PATH Not Active in Current Shell After Install

## Ticket

After installing gitmap-v27 via `curl | bash` on macOS (or `irm | iex` on Windows),
running `gitmap-v27` immediately returns "command not found" / "not recognized".
The user must manually source their profile or open a new terminal.

## Symptoms

### Unix (macOS / Linux)

1. User runs `curl -fsSL .../install.sh | bash` on macOS (zsh default).
2. Installer adds PATH entry to `~/.zprofile` (or `~/.zshrc`).
3. Installer prints `export PATH=...` and reload instructions.
4. User types `gitmap-v27` → `zsh: command not found: gitmap-v27`.

### Windows

1. User runs `irm .../install.ps1 | iex` in PowerShell.
2. Installer adds install directory to User PATH in the registry.
3. User types `gitmap-v27` → `gitmap-v27: The term 'gitmap-v27' is not recognized`.
4. CMD and Git Bash sessions also cannot find `gitmap-v27`.

## Root Cause

Three separate issues compound across platforms:

### 1. Subshell / process isolation (unfixable)

- **Unix**: `curl | bash` runs the installer in a child process. The
  `export PATH=...` only affects that subshell — the parent interactive
  shell never receives the updated PATH (fundamental POSIX behavior).
- **Windows**: `irm | iex` runs in the current session, but registry-based
  PATH changes only take effect in **new** processes. CMD and Git Bash
  sessions that are already open never see the update.

### 2. Single-profile / single-environment write (fixable)

- **Unix**: The installer wrote the PATH entry to only **one** profile file,
  chosen by `$SHELL` detection. Users in a different shell or terminal
  emulator that reads a different profile wouldn't find gitmap.
- **Windows**: The installer only updated the Windows Registry User PATH.
  PowerShell `$PROFILE` and Git Bash profiles (`~/.bashrc`,
  `~/.bash_profile`) were not touched, so those environments required
  the user to open a brand-new terminal.

### 3. No immediate activation instruction

The post-install message was either absent or buried, making users expect
`gitmap-v27` to just work immediately after install.

## Fix

### Phase 1: Multi-profile PATH registration

Write the PATH entry to **all** detected profile files that exist or are
standard for the platform. Each write is idempotent — skipped if the
entry already exists in the file.

#### Unix (`install.sh`)

| Shell | Profiles written |
|-------|-----------------|
| zsh   | `~/.zshrc` AND `~/.zprofile` (create `.zshrc` if neither exists) |
| bash  | `~/.bashrc` AND `~/.bash_profile` (or `~/.profile`) |
| fish  | `~/.config/fish/config.fish` |

Additionally, always write to `~/.profile` as a catch-all for POSIX `sh`.

#### Windows (`install.ps1`)

| Environment | Target |
|-------------|--------|
| System-wide | Windows Registry — User `PATH` environment variable |
| PowerShell  | `$PROFILE` (`Microsoft.PowerShell_profile.ps1`) |
| Git Bash    | `~/.bashrc` AND `~/.bash_profile` |

All profile writes use a `# gitmap-path` marker for idempotent
insertion and future cleanup during uninstall.

### Phase 2: Immediate activation guidance

After installation, print a **prominent**, environment-specific activation
command that the user can copy-paste immediately.

#### Unix example

```
  ✓ Installed! To start using gitmap-v27 right now, run:

      source ~/.zshrc

  Or open a new terminal window.
```

#### Windows example

```
  ✓ Installed! To start using gitmap-v27 right now:

    PowerShell:  $env:PATH += ";C:\Users\you\.local\bin"
    CMD:         set PATH=%PATH%;C:\Users\you\.local\bin
    Git Bash:    export PATH="$PATH:/c/Users/you/.local/bin"

  Or open a new terminal window.
```

### Phase 3: Session PATH export

- **Unix**: The script does `export PATH="${PATH}:${dir}"` which works when
  sourced directly (`source install.sh`) but not via pipe. This limitation
  is documented in the post-install output.
- **Windows**: The script does `$env:PATH += ";${dir}"` which takes effect
  in the current PowerShell session immediately. Registry changes require
  a new process; `SendMessageTimeout` broadcasts `WM_SETTINGCHANGE` to
  notify other applications.

## Prevention

1. All installer scripts must write PATH to **multiple** profile files /
   environments to cover login shells, interactive shells, POSIX sh,
   PowerShell, CMD, and Git Bash.
2. Post-install output must show a single, copy-pasteable activation command
   **prominently** (not buried in a summary block) for each detected
   environment.
3. The installer must explicitly state that pipe-based invocation
   (`curl | bash`, `irm | iex`) cannot always modify the parent shell's
   environment.
4. Profile writes must use an idempotent marker (e.g. `# gitmap-path`) to
   prevent duplicate entries and enable clean uninstall.

## Related

- `spec/02-app-issues/20-path-not-available-in-other-shells.md` — cross-shell visibility (superseded)
- `spec/01-app/94-install-script.md` — installer specification
- `gitmap-v27/scripts/install.sh` — Unix implementation
- `gitmap-v27/scripts/install.ps1` — Windows implementation
