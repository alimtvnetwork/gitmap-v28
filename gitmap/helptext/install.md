# Tool Installer

Install a developer tool by name using the platform package manager.

## Alias

in

## Usage

    gitmap install <tool> [flags]

## Flags

| Flag      | Default | Description                                        |
|-----------|---------|----------------------------------------------------|
| --manager | (auto)  | Force package manager (choco, winget, apt, brew)   |
| --version | latest  | Install a specific version                         |
| --verbose | false   | Show full installer output                         |
| --dry-run | false   | Show install command without executing             |
| --check   | false   | Only check if tool is installed                    |
| --list    | false   | List all supported tools, grouped by category, with installed status (● = installed, ○ = not installed) |

## Supported Tools

| Tool            | Binary         | Description                      |
|-----------------|----------------|----------------------------------|
| vscode          | code           | Visual Studio Code editor        |
| node            | node           | Node.js JavaScript runtime       |
| yarn            | yarn           | Yarn package manager             |
| bun             | bun            | Bun JavaScript runtime           |
| pnpm            | pnpm           | pnpm package manager             |
| python          | python3        | Python programming language      |
| go              | go             | Go programming language          |
| git             | git            | Git version control              |
| git-lfs         | git-lfs        | Git Large File Storage           |
| gh              | gh             | GitHub CLI                       |
| github-desktop  | —              | GitHub Desktop application       |
| cpp             | g++            | C++ compiler (MinGW/g++)         |
| php             | php            | PHP programming language         |
| powershell      | pwsh           | PowerShell shell                 |

## Notepad++ Variants

| Command         | Shortcut        | Description                              |
|-----------------|-----------------|------------------------------------------|
| npp             | NPP + Settings  | Install Notepad++ and sync settings      |
| npp-settings    | NPP Settings    | Sync Notepad++ settings only             |
| install-npp     | Install NPP     | Install Notepad++ only (no settings)     |

Settings are extracted from a bundled zip to `%APPDATA%\Notepad++`.

## Scripts

| Command         | Description                                          |
|-----------------|------------------------------------------------------|
| scripts         | Clone gitmap scripts to a local folder               |

- **Windows**: Reads deploy drive from `powershell.json`, defaults to `D:\gitmap-scripts`
- **Linux/macOS**: Installs to `~/Desktop/gitmap-scripts`

Copies: `install.ps1`, `install.sh`, `run.ps1`, `run.sh`, `uninstall.ps1`, `Get-LastRelease.ps1`.

## Coding Guidelines

| Command       | Aliases              | Description                                              |
|---------------|----------------------|----------------------------------------------------------|
| clean-code    | code-guide, cg, cc   | Install alimtvnetwork coding-guidelines (v15) one-liner  |

All four invocations dispatch to the same flow:

    $ gitmap install clean-code
    $ gitmap install code-guide
    $ gitmap i cg
    $ gitmap i cc

Each runs:

    irm https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v15/main/install.ps1 | iex

Requires PowerShell on PATH (Windows: `powershell` ships by default; Linux/macOS: install `pwsh` 7+).

## Prerequisites

- Windows: Chocolatey or Winget in PATH
- Linux: apt, dnf, or pacman available
- macOS: Homebrew installed

## Examples

### NPP + Settings — Install Notepad++ with settings

    $ gitmap install npp
      Checking if npp is installed...
      Installing npp...
      Verifying npp installation...
      npp installed successfully.
      Verifying npp binary at: C:\Program Files\Notepad++\notepad++.exe
      Binary confirmed: C:\Program Files\Notepad++\notepad++.exe
      Syncing Notepad++ settings...
      Extracting Notepad++ settings to C:\Users\User\AppData\Roaming\Notepad++...
      Settings synced to C:\Users\User\AppData\Roaming\Notepad++

### NPP Settings — Sync settings only

    $ gitmap install npp-settings
      Skipping Notepad++ installation (settings-only mode)
      Syncing Notepad++ settings...
      Extracting Notepad++ settings to C:\Users\User\AppData\Roaming\Notepad++...
      Settings synced to C:\Users\User\AppData\Roaming\Notepad++

### Install NPP — Install Notepad++ only (no settings)

    $ gitmap install install-npp
      Checking if npp is installed...
      Installing npp...
      Verifying npp installation...
      npp installed successfully.
      Skipping Notepad++ settings (install-only mode)

### Install a tool end-to-end

    $ gitmap install vscode
      Checking if vscode is installed...
      Installing vscode...
      Verifying vscode installation...
      vscode installed successfully.

### Check if a tool is already installed

    $ gitmap in go --check
      Checking if go is installed...
      go is already installed (version: go version go1.22.4 linux/amd64)

### Preview install command with dry-run

    $ gitmap install python --dry-run
      Checking if python is installed...
      [dry-run] Would run: choco install python -y

### Clone gitmap scripts

    $ gitmap install scripts
      → Scripts target: /home/alim/Desktop/gitmap-scripts
      Cloning gitmap repo for scripts...
      ✓ Copied: install.ps1
      ✓ Copied: install.sh
      ✓ Copied: run.ps1
      ✓ Copied: run.sh
      ✅ 6 scripts installed to /home/alim/Desktop/gitmap-scripts

## See Also

- [env](env.md) — Manage environment variables and PATH
- [doctor](doctor.md) — Diagnose PATH and version issues
- [setup](setup.md) — Configure Git global settings

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter install
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
