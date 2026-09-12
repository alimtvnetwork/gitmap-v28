# Tool Installer

Install developer tools, runtimes, database servers, AI models, custom scripts, and workstation profiles across Windows, Linux, and macOS using native platform package managers.

## Alias

in

## Usage

```bash
gitmap install <tool|profile> [flags]
gitmap install profile [name] [--tree]
gitmap in <tool|profile> [flags]
```

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `profile [name]` | Run an installation profile bundle, or list all profiles |
| `logs [tool]` | View installer execution and error logs (`--tail N`, `--clear`) |
| `add <name> [ver]` | Interactively create and register a custom tool installer |
| `export <tool>` | Export custom installer definitions to YAML/JSON (`--all`) |
| `import <file>` | Import custom installer configurations |

## Flags

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--manager <name>` | — | (auto) | Force package manager (`choco`, `winget`, `apt`, `brew`, `snap`, `dnf`, `pacman`) |
| `--version <ver>` | — | latest | Install a specific version |
| `--verbose` | — | false | Show full installer output |
| `--dry-run` | — | false | Show install command without executing |
| `--check` | — | false | Only check if tool is installed |
| `--tree` | `-t` | false | Preview full tool hierarchy of a profile before installing |
| `--list` | `ls` | false | List all supported tools grouped by category with installed status (● = installed, ○ = not installed) |
| `--yes` | `-y` | false | Auto-confirm installation prompts without interactive questions |
| `--explain` | — | false | Print exact resolved shell command before executing |
| `--status` | — | false | Show installed tools recorded in the local SQLite database |
| `--upgrade` | — | false | Upgrade an already-installed tool to latest version |

## Supported Tools Across 6 Categories

### 1. Core Tools

| Tool | Binary | Description |
|------|--------|-------------|
| `gitmap` | `gitmap` | Gitmap core orchestration binary |
| `vscode` | `code` | Visual Studio Code editor |
| `node` | `node` | Node.js JavaScript runtime |
| `yarn` | `yarn` | Yarn package manager |
| `bun` | `bun` | Bun JavaScript runtime |
| `pnpm` | `pnpm` | pnpm package manager |
| `python` | `python3` | Python programming language |
| `go` | `go` | Go programming language |
| `git` | `git` | Git version control |
| `git-lfs` | `git-lfs` | Git Large File Storage |
| `gh` | `gh` | GitHub CLI |
| `github-desktop` | — | GitHub Desktop application |
| `cpp` | `g++` | C++ compiler (MinGW/g++) |
| `php` | `php` | PHP programming language |
| `powershell` | `pwsh` | PowerShell shell |
| `chocolatey` | `choco` | Chocolatey package manager (Windows) |
| `winget` | `winget` | Winget package manager (Windows) |
| `chrome` | `google-chrome` | Google Chrome web browser |
| `build-essential` | `gcc, g++` | Ubuntu compiler toolchain & dev libraries |
| `wordpress` | `wp` | WP-CLI WordPress command line interface |
| `antigravity` | `agy` | Antigravity CLI autonomous coding assistant |
| `ag-manager` | — | Antigravity Manager GUI desktop application |
| `ag-ctx` | — | Add Antigravity to right-click context menu |
| `npp` | `notepad++` | Notepad++ text editor with synced settings |
| `npp-settings` | — | Sync Notepad++ configuration settings only |
| `install-npp` | `notepad++` | Install Notepad++ only (without settings) |
| `vscode-settings` | — | Sync VS Code settings and keybindings |
| `obs` | `obs` | OBS Studio screen recording and streaming |
| `obs-settings` | — | Sync OBS Studio profiles and scenes |
| `wt-settings` | — | Sync Windows Terminal settings |
| `dbeaver` | `dbeaver` | Universal database GUI client |
| `sticky-notes` | — | Sticky Notes desktop utility |
| `scripts` | — | Clone gitmap scripts to local folder |

### 2. Databases (All 12 Supported DBs)

| Database | Canonical Slug | Default Port | Description |
|----------|----------------|--------------|-------------|
| MySQL | `mysql` | 3306 | High-performance relational database |
| MariaDB | `mariadb` | 3306 | Open-source MySQL-compatible database fork |
| PostgreSQL | `postgresql` | 5432 | Advanced open-source object-relational database |
| SQLite | `sqlite` | file-based | Self-contained, serverless embedded database |
| MongoDB | `mongodb` | 27017 | Scalable document-oriented NoSQL database |
| CouchDB | `couchdb` | 5984 | Document database with native HTTP/JSON REST API |
| Redis | `redis` | 6379 | In-memory key-value data structure store & cache |
| Apache Cassandra | `cassandra` | 9042 | Distributed wide-column NoSQL database |
| Neo4j | `neo4j` | 7474 / 7687 | Native graph database management platform |
| Elasticsearch | `elasticsearch` | 9200 | Distributed search and analytics engine |
| DuckDB | `duckdb` | in-process | High-performance columnar analytical SQL database |
| LiteDB | `litedb` | file-based | Embedded NoSQL document store for .NET |

### 3. Languages & Runtimes

| Tool | Aliases | Description |
|------|---------|-------------|
| `rust` | `rustup`, `cargo` | Rust programming language and Cargo package manager |
| `dotnet` | `dotnet-sdk` | .NET SDK and developer runtime |
| `java` | `jdk`, `openjdk` | OpenJDK Java Development Kit |
| `flutter` | — | Flutter SDK for cross-platform app development |
| `laravel` | `artisan`, `laravel-installer` | Laravel PHP framework installer |
| `composer` | — | Dependency manager for PHP |

### 4. Local AI Suite

| Tool | Aliases | Description |
|------|---------|-------------|
| `ollama` | — | Local LLM runner and model serving daemon |
| `llama-cpp` | `llamacpp` | llama.cpp high-performance local inference engine |
| `python-libs` | — | Core AI/ML Python stack (`numpy`, `pandas`, `torch`, `transformers`) |
| `antigravity` | `agy`, `ag` | Autonomous AI coding assistant CLI |
| `ag-manager` | `manager`, `agy-manager` | Antigravity Manager GUI desktop application |

### 5. DevOps & Containers

| Tool | Aliases | Description |
|------|---------|-------------|
| `docker` | — | Docker container runtime and orchestration engine |
| `kubernetes` | `k8s`, `kubectl` | Kubernetes command line container orchestration |
| `jenkins` | — | Jenkins CI/CD automation server |
| `nginx` | `ngx`, `engine-x` | High-performance HTTP server and reverse proxy |
| `vmware` | `open-vm-tools`, `vmtools` | VMware workstation virtualization and guest tools |
| `open-vm-tools` | `vmtools`, `vmware-tools` | VMware guest utilities and shared folders |

### 6. Terminal & Utilities

| Tool | Aliases | Description |
|------|---------|-------------|
| `zsh` | — | Z shell command environment |
| `flameshot` | — | Feature-rich screen capture and annotation tool |
| `conemu` | — | Windows console emulator with tabs and split panes |
| `vlc` | — | VLC media player multimedia framework |
| `qbittorrent` | `qtorrent`, `qbit` | qBittorrent open-source BitTorrent client |
| `utorrent` | `u-torrent` | uTorrent lightweight BitTorrent client |

## Custom Script Tools

Gitmap provides one-liner automated bootstrapping for specialized developer tools:

| Command | Platform Script | Description |
|---------|-----------------|-------------|
| `scripts-fixer` | `install.ps1` / `install.sh` | Gitmap repository scripts and path fixer / auto-repair suite |
| `coding-guidelines` (`cg`, `cc`, `code-guide`) | `install.ps1` / `install.sh` (v24) | AlimTV Network Coding Guidelines (v24) automated compliance installer |
| `macro-ahk` | `download-extension.ps1` / `install.sh` (v55) | AutoHotkey v2 automation and shortcut extension suite |

### Custom Script Examples

```bash
# Install Gitmap scripts fixer
$ gitmap install scripts-fixer

# Install coding guidelines compliance checker (v24)
$ gitmap install coding-guidelines
$ gitmap in cg

# Install AutoHotkey macro extension suite
$ gitmap install macro-ahk
```

## Installation Profiles

Installation profiles install a curated bundle of developer tools in a single command. Profiles verify already installed tools and only download missing components.

| Profile | Aliases | Description | Key Tools Included |
|---------|---------|-------------|-------------------|
| `minimal` | `min`, `basic` | Minimal dev workstation | VS Code, Git, Node.js, Python |
| `base` | `ubuntu-basic`, `ub` | Essential OS toolchain & shell | curl, git, build-essential, zsh |
| `dev` | `developer`, `dev-stack` | Standard workstation with AI suite | VS Code, Git, Python, Node, pnpm, Go, Rust, PHP, Antigravity, AG-Manager |
| `small-dev` | `ub+sdev`, `ubuntu-small-dev` | Lightweight dev suite with runtimes | basic + VS Code + Go + Node.js |
| `advance` | `ub+dev`, `ubuntu+dev` | Full developer workstation suite | small-dev + Docker + Python 3 |
| `dev-advance` | — | Advanced full-stack developer stack | Full developer workstation + databases + AI |
| `terminal` | — | Terminal power-user tools | Zsh, ConEmu, Git, curl, wget, nano, vim |
| `ubuntu` | `ubuntu-dev`, `linux-dev` | Ubuntu developer workstation | build-essential, Git, Zsh, VS Code, Chrome, Node, Python, Go, Antigravity, AG-Manager |
| `ai` | `ai-dev`, `ml`, `llm` | AI / ML workstation | Python, Ollama, llama-cpp, Python ML libs, Antigravity, AG-Manager |
| `backend` | `back`, `server` | Backend developer workstation | Minimal stack + Docker, MySQL, PostgreSQL, Redis, Go, .NET, Java |
| `fullstack` | `full`, `web` | Full-stack web developer workstation | Backend stack + pnpm, PHP, Composer, MongoDB, Jenkins CI/CD |
| `web-dev` | — | Modern web development stack | Node.js, pnpm, Yarn, Bun, PHP, Composer, MySQL |
| `devops` | — | Cloud & container infrastructure | Docker, Kubernetes, Jenkins, Nginx, VMware |

### Profile Usage & Tree Preview

```bash
# List all available profiles and installed tool count
$ gitmap install profile
$ gitmap install profile --list

# Preview tool hierarchy before installing (--tree / -t)
$ gitmap install profile dev --tree
$ gitmap in dev -t
$ gitmap install --tree

# Install profile directly
$ gitmap install profile dev
$ gitmap in dev -y
$ gitmap install ubuntu --dry-run
```

## Logs Management

Every install and profile execution is tracked with detailed logs:

```bash
# View recent install logs
$ gitmap install logs

# View logs for a specific tool
$ gitmap install logs node --tail 50

# Shortcut
$ gitmap in logs
```

## Custom Installer Add & Export/Import

```bash
# Interactively register a new tool installer
$ gitmap install add my-tool 1.0.0

# Export installer configuration to JSON/YAML
$ gitmap install export my-tool
$ gitmap install export --all

# Import installer configuration
$ gitmap install import my-tool.yaml
```

## Examples

### Check if a tool is installed

```bash
$ gitmap in go --check
  Checking if go is installed...
  go is already installed (version: go version go1.22.4 linux/amd64)
```

### Preview install command with dry-run

```bash
$ gitmap install python --dry-run
  Checking if python is installed...
  [dry-run] Would run: choco install python -y
```

### Force a specific package manager

```bash
$ gitmap install node --manager winget --yes
```

### Install with verbose output

```bash
$ gitmap install rust --verbose
```

### Explain command before executing

```bash
$ gitmap install ag-manager --explain
```

## See Also

- [installer](installer.md) — Manage custom installer scripts, export, import, and versioning
- [env](env.md) — Manage environment variables and PATH
- [doctor](doctor.md) — Diagnose PATH and version issues
- [setup](setup.md) — Configure Git global settings
