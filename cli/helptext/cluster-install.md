# gitmap cluster install

Install GitMap remotely across cluster nodes via SSH using official installation one-liners and automated bootstrapping.

## Usage

```bash
gitmap cluster install [component] [target] [flags]
gitmap sj install [component] [target] [flags]
```

## Description

`gitmap cluster install` connects to target cluster nodes over SSH and automates the installation or upgrade of GitMap.

Supported target topologies:
- `all`: Broadcast installation across all registered nodes in parallel.
- `control` / `workers`: Target nodes matching specific cluster roles.
- `<alias>`: Target a specific host by its memorable alias (e.g., `devbox`, `k8s-w1`).
- `<user@ip|ip>`: Target an ad-hoc IP address or host directly.

### Platform-Specific Installers
- **Linux & macOS**: Fetches and executes official script via `curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | bash`.
- **Windows**: Executes official PowerShell one-liner via `irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex`.

### Automatic Preflight Bootstrapping Engine
When running commands remotely with `gitmap cluster exec` or `gitmap sc bash`, GitMap inspects whether the command requires `gitmap`. If `gitmap` is not found on the remote machine (`command -v gitmap` fails), the cluster runner **automatically bootstraps and installs GitMap** on the fly before executing the command.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--version` | `-v` | `latest` | Specific GitMap version to install (e.g. `v6.256.0` or `latest`) |
| `--sudo` | `-s` | `true` | Execute installer with elevated privileges using stored or prompted sudo credentials |
| `--parallel` | `-p` | `4` | Number of concurrent SSH worker threads for fan-out installation |
| `--os` | | `linux` | Target operating system runtime (`linux`, `darwin`, `windows`) |
| `--help` | `-h` | | Show command help |

## Examples

```bash
# 1. Install latest GitMap across all registered cluster nodes in parallel
gitmap cluster install gitmap all

# 2. Install GitMap on a single node by alias using cluster or sj syntax
gitmap cluster install gitmap devbox
gitmap sj install gitmap devbox

# 3. Install a specific GitMap version across all worker nodes
gitmap cluster install gitmap workers --version v6.256.0

# 4. Install GitMap without sudo elevation
gitmap cluster install gitmap devbox --sudo=false

# 5. Concurrently install GitMap across 8 nodes in parallel
gitmap cluster install gitmap all --parallel 8

# 6. Install on a Windows remote host via PowerShell
gitmap cluster install gitmap win-worker --os windows

# 7. Auto-bootstrapping in action: remotely scheduled shutdown installs GitMap if missing
gitmap cluster exec devbox "gitmap schedule shutdown 1:45hr"
```

See also: `gitmap cluster exec`, `gitmap cluster bootstrap`, `gitmap cluster nodes`, `gitmap ssh-join`
