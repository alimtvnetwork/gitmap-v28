# self-install

Install (or re-install) the gitmap binary on this machine.

## Synopsis

```
gitmap self-install [--dir <path>] [--yes] [--version <tag>]
                    [--shell-mode <mode>]
                    [--show-path] [--force-lock]
```

## What it does

1. Resolves the install directory:
   - `--dir <path>` if supplied.
   - Default with prompt otherwise:
     - **Windows**: `D:\gitmap`
     - **Unix**: `~/.local/bin/gitmap`
   - `--yes` accepts the default without prompting.
2. Loads the platform installer from one of two sources:
   - **Embedded**: `install.ps1` / `install.sh` shipped inside the binary
     via `go:embed`. No network needed.
   - **Remote** (fallback): downloaded from
     `raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/gitmap/scripts/`.
3. Writes the script to a temp file (UTF-8 BOM on PowerShell), runs it
   with `-InstallDir` / `--dir`, and forwards `--version` if pinned.

## --shell-mode <mode>

Controls which shell profile files receive the PATH snippet on Unix.
Defaults to `auto`.

### Singleton modes

| Mode   | Writes PATH to                                             |
|--------|------------------------------------------------------------|
| `auto` | Detected shell profiles (current behavior, default)        |
| `both` | zsh + bash + .profile + fish (if installed) + pwsh         |
| `zsh`  | `~/.zshrc` and `~/.zprofile` only                          |
| `bash` | `~/.bashrc` and `~/.bash_profile` only                     |
| `pwsh` | `~/.config/powershell/Microsoft.PowerShell_profile.ps1`    |
| `fish` | `~/.config/fish/config.fish`                               |

### Combo modes (v3.48.0+)

Any `+`-joined combination of the concrete shell families
(`zsh`, `bash`, `pwsh`, `fish`) writes PATH to **only** those families.
Combos are **strict** — `~/.profile` and any unlisted family are skipped.

| Combo            | Writes PATH to                                                    |
|------------------|-------------------------------------------------------------------|
| `zsh+pwsh`       | zsh profiles + pwsh profile only (recommended for macOS pwsh users) |
| `bash+fish`      | bash profiles + fish config only                                  |
| `zsh+bash+pwsh`  | zsh + bash + pwsh, skip fish and `~/.profile`                     |

`auto` and `both` cannot appear inside a combo (they're meta values, not
shell families).

`--profile <mode>` is kept as a hidden alias for `--shell-mode <mode>`.
`--dual-shell` is kept as a hidden alias for `--shell-mode both`.

When the resolved mode includes `pwsh` (`both`, `pwsh` singleton, or any
combo containing `pwsh`), the installer also exports `GITMAP_DUAL_SHELL=1`
into `install.sh`'s environment so `detect_active_pwsh` fires regardless
of the parent shell — guaranteeing the pwsh profile is written even when
the installer is launched from zsh.

## Examples

```
gitmap self-install
gitmap self-install --yes
gitmap self-install --dir D:\dev\gitmap
gitmap self-install --version v3.0.0
gitmap self-install --shell-mode both              # write every detected profile
gitmap self-install --shell-mode zsh+pwsh          # macOS pwsh user, deterministic
gitmap self-install --shell-mode bash+fish         # Linux user with two shells
gitmap self-install --shell-mode zsh+bash+pwsh     # cross-shell dev workstation
gitmap self-install --shell-mode pwsh              # only touch the pwsh profile
gitmap self-install --show-path                    # audit which profiles got written
```

## See also

- `gitmap self-uninstall` — remove gitmap from this machine
- `gitmap update` — pull a newer build from the source repo

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter self-install
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
