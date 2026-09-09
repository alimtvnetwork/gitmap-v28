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
| antigravity     | agy            | Antigravity CLI autonomous coding assistant |
| ag-manager      | —              | Antigravity Manager GUI desktop application |
| ag-ctx          | —              | Add Antigravity to right-click context menu |
| build-essential | gcc, g++       | Ubuntu compiler toolchain & dev libraries   |
| chrome          | google-chrome  | Google Chrome web browser        |
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
| rust            | rustc, cargo   | Rust programming language & Cargo |
| docker          | docker         | Docker container platform        |
| qtorrent        | qbittorrent    | qBittorrent BitTorrent client    |
| utorrent        | utorrent       | uTorrent BitTorrent client       |

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
 
## Antigravity AI Suite

Install the Antigravity autonomous coding assistant and Antigravity Manager GUI desktop application:

| Command | Aliases | Description |
|---------|---------|-------------|
| antigravity | agy, antigravity-cli | Autonomous coding assistant CLI (`https://get.antigravity.dev`) |
| ag-manager | manager, agy-manager, ag-m | Antigravity Manager GUI (latest GitHub release via git tags) |
| ag-ctx | — | Add Antigravity right-click context menu integration |

### Antigravity Manager Git Tag Discovery
`ag-manager` automatically queries GitHub releases and remote git tags (`https://github.com/lbjlaq/Antigravity-Manager.git`) using semantic version comparison to select and download the newest release matching your operating system (`.exe`/`.msi` on Windows, `.deb`/`.AppImage` on Linux, `.dmg` on macOS).

```bash
$ gitmap install antigravity
$ gitmap install ag-manager
$ gitmap install ag-manager --version 4.6.9
$ gitmap in agy
```

### AGY Subcommand Invocation
Antigravity tools can also be installed directly via `gitmap agy install`:

```bash
$ gitmap agy install              # Installs Antigravity Manager GUI (default)
$ gitmap agy install manager      # Installs Antigravity Manager GUI
$ gitmap agy install cli          # Installs Antigravity CLI (agy)
$ gitmap agy install all          # Installs both Manager and CLI
$ gitmap agy in manager --dry-run # Preview installation plan
```

## Installation Profiles

Installation profiles install a curated bundle of developer tools in a single command. Profiles verify already installed tools and only download missing components.

| Profile | Aliases | Description | Key Tools Included |
|---------|---------|-------------|-------------------|
| `dev` | developer, dev-stack | Standard developer workstation with AI | VS Code, Git, Python, Node, pnpm, Go, Rust, PHP, Antigravity, AG-Manager |
| `ubuntu` | ubuntu-dev, linux-dev | Ubuntu developer workstation | build-essential, Git, Zsh, VS Code, Chrome, Node, Python, Go, Antigravity, AG-Manager |
| `ubuntu-dev-ai` | ubuntu-ai, linux-ai | Full Ubuntu AI / ML developer workstation | Ubuntu workstation + Ollama, llama-cpp, Python ML libs, Antigravity AI suite |
| `ai` | ai-dev, ml, llm | AI / ML workstation | Python, Ollama, llama-cpp, Python ML libs, Antigravity, AG-Manager |
| `minimal` | min, basic | Essential minimal developer workstation | VS Code, Git, Node.js, Python |
| `backend` | back, server | Backend developer workstation | Minimal stack + Docker, MySQL, PostgreSQL, Redis, Go, .NET, Java |
| `fullstack` | full, web | Full-stack web developer workstation | Backend stack + pnpm, PHP, Composer, MongoDB, Jenkins CI/CD |

### Profile Usage

```bash
$ gitmap install dev
$ gitmap in ubuntu
$ gitmap install ai --dry-run
$ gitmap install profile          # List all profiles and progress
```

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

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
