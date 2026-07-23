# self-uninstall

Remove gitmap from this machine.

## Synopsis

```
gitmap self-uninstall [--confirm] [--keep-data] [--keep-snippet]
                      [--shell-mode <mode>]
```

## What it removes

| Target              | Path                                                |
|---------------------|-----------------------------------------------------|
| Binary + deploy dir | Folder containing the running `gitmap` / `gitmap.exe` |
| Data dir            | `<deploy>/.gitmap/` (skip with `--keep-data`)       |
| PATH snippet        | `# gitmap shell wrapper v*` block in resolved profiles |
| Completion files    | `gitmap-completion.{bash,zsh,fish}` in deploy dir   |

The set of profiles touched depends on `--shell-mode` (default `auto`,
which strips every known profile across zsh, bash, pwsh, and fish for
the safest full removal). Mirrors the install-side flag — see
`gitmap self-install --help` for the full singleton/combo table.

| `--shell-mode`        | Profiles cleaned (Unix)                                                  |
|-----------------------|--------------------------------------------------------------------------|
| `auto` / `both`       | every known profile across all families                                  |
| `zsh`                 | `~/.zshrc`, `~/.zprofile`                                                |
| `bash`                | `~/.bashrc`, `~/.bash_profile`                                           |
| `pwsh`                | `~/.config/powershell/{profile,Microsoft.PowerShell_profile}.ps1`        |
| `fish`                | `~/.config/fish/config.fish`                                             |
| `zsh+pwsh` (combo)    | strict union of zsh and pwsh files; bash + fish + `~/.profile` skipped   |

On Windows only the pwsh profile family is meaningful; non-pwsh tokens
in a combo resolve to no files (no error).

## Confirmation

Without `--confirm`, the command lists the targets and waits for the
user to type `yes`. Pass `--confirm` to skip the prompt.

## Windows handoff

The running `gitmap.exe` cannot delete itself while loaded. When the
binary lives inside the deploy dir, self-uninstall copies itself to
`%TEMP%\gitmap-handoff-<pid>.exe` and re-execs the hidden
`self-uninstall-runner` from there, which performs the deletion and
then schedules its own removal via `cmd.exe /C ... del`.

## Re-installing afterwards

```
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.sh | bash

# Windows (PowerShell)
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/install.ps1 | iex
```

…or, if you still have a copy of the binary somewhere, just run
`gitmap self-install`.

## Examples

```
gitmap self-uninstall
gitmap self-uninstall --confirm
gitmap self-uninstall --confirm --keep-data
gitmap self-uninstall --confirm --keep-snippet --keep-data
gitmap self-uninstall --confirm --shell-mode zsh+pwsh
gitmap self-uninstall --confirm --shell-mode pwsh
```

## See also

- `gitmap self-install` — install or re-install the gitmap binary
- `spec/04-generic-cli/21-post-install-shell-activation/04-idempotency.md` —
  marker block conventions used to locate and strip the PATH snippet

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter self-uninstall
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
