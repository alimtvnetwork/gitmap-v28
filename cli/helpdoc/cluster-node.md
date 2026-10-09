# gitmap cluster node

Manage cluster nodes and execute Ubuntu node provisioning recipes.

## Usage

```bash
gitmap cluster node <subcommand> [args...] [flags]
```

## Lifecycle Subcommands

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| `add` | `join`, `enroll`, `new` | Register and enroll a node into cluster |
| `rm` | `remove`, `delete` | Remove a node from cluster inventory |
| `ls` | `list`, `nodes` | List all registered cluster nodes |

## Provisioning Recipes

| Recipe | Syntax | Description |
|--------|--------|-------------|
| `set-ip` | `gitmap cluster node set-ip <target> <new-ip> [--route-ip <gw>]` | Configure static Netplan IP and default gateway route |
| `install-base` | `gitmap cluster node install-base <target>` | Install essential development tools and shell runtimes |
| `create-user` | `gitmap cluster node create-user <target> <username> [pass] [--theme <t>]` | Create sudo user with zsh and Oh-My-Zsh |
| `set-theme` | `gitmap cluster node set-theme <target> <theme>` | Update `ZSH_THEME` in `~/.zshrc` across target nodes |
| `purge` | `gitmap cluster node purge <target>` | Purge unneeded packages and clean apt cache |

## Target Selectors

All recipes accept semantic cluster target selectors:
- `all`: Targets every node in the cluster
- `control`: Targets control-plane master nodes
- `workers`: Targets worker nodes
- `<alias|ip>`: Targets a single node by alias or IP address

## Flags

- `-h`, `--help`: Show command help.
- `--route-ip <gw>`: Gateway route IP for `set-ip` (default: `192.168.0.1`).
- `--theme <t>`, `-t <t>`: Oh-My-Zsh theme for `create-user` (default: `fletcherm`).

## Examples

```bash
# 1. Enroll a new node into the cluster
gitmap cluster node add ubuntu@192.168.0.101 k8s-w1

# 2. List all registered cluster nodes
gitmap cluster node ls

# 3. Remove a node from the cluster
gitmap cluster node rm k8s-w1

# 4. Configure static Netplan IP with default gateway
gitmap cluster node set-ip k8s-w1 192.168.0.101 --route-ip 192.168.0.1

# 5. Install base developer tools (curl, git, zsh, build-essential) across all nodes
gitmap cluster node install-base all

# 6. Create dedicated sudo user with Oh-My-Zsh on worker nodes
gitmap cluster node create-user workers kube SecretPass123 --theme agnoster

# 7. Update ZSH theme across all cluster nodes
gitmap cluster node set-theme all fletcherm

# 8. Clean up unused packages and apt cache to reduce VM disk footprint
gitmap cluster node purge all
```

See also: `gitmap cluster add`, `gitmap cluster remove`, `gitmap cluster nodes`, `gitmap cluster`
