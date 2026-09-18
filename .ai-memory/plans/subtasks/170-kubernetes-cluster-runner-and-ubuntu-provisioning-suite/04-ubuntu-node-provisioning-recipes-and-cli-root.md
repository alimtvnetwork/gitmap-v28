# Subtask 04: Ubuntu Node Provisioning Recipes & CLI Root Integration

## Objective
Implement native CLI subcommands for Ubuntu node provisioning (Netplan static IP, base packages, root/sudo user with Oh-My-Zsh, auto-purge) ported from `kubernetes-training/02-ubuntu-install`, wire them into `gitmap cluster node <recipe>`, register `cluster` routing in CLI root, and write comprehensive markdown help documentation.

## Disjoint Files Assigned
- `cli/cmdssh/cluster_node_recipes.go`
- `cli/cmdssh/cluster_node_cmd.go`
- `cli/cmdssh/cluster_node_cmd_test.go`
- `cli/cmd/cluster.go`
- `cli/helptext/cluster.md`
- `cli/helptext/catalog.go`

## Implementation Details
1. `cli/cmdssh/cluster_node_recipes.go`:
   - Port bash logic from `02-ubuntu-install`:
     * `GenerateNetplanScript(ip string, routeIP string) string`:
       Creates `/etc/netplan/00-installer-config.yaml` with ens33 static IP, route, nameservers [8.8.8.8, 1.1.1.1], and runs `netplan apply`.
     * `GenerateBasePackagesScript() string`:
       Installs `curl wget git zsh net-tools htop build-essential`.
     * `GenerateCreateUserScript(username, password, theme string) string`:
       Creates user with `/usr/bin/zsh`, configures `/etc/sudoers.d/`, installs oh-my-zsh, configures theme.
     * `GeneratePurgeScript() string`:
       Runs `apt-get autoremove --purge -y && apt-get clean`.

2. `cli/cmdssh/cluster_node_cmd.go`:
   - Implement `RunClusterNodeCLI(args []string) error`:
     - Subcommands:
       * `set-ip <target> <new-ip> [--route-ip <gw>]`
       * `install-base <target>`
       * `create-user <target> <username> [password] [--theme <theme>]`
       * `set-theme <target> <theme>`
       * `purge <target>`
     - Dispatches generated recipe script using `ExecuteNodeCommand` with sudo.

3. `cli/cmd/cluster.go`:
   - Wire `runCluster`:
     * `"exec"`, `"run"` -> `cmdssh.RunClusterExecCLI(args[1:])`
     * `"run-script"`, `"script"` -> `cmdssh.RunClusterScriptCLI(args[1:])`
     * `"node"` -> `cmdssh.RunClusterNodeCLI(args[1:])`
     * `"import"` -> `cmdssh.RunClusterImportCLI(args[1:])`

4. `cli/helptext/cluster.md`:
   - Write comprehensive user guide and help text with:
     * Architecture overview.
     * Cluster topology JSON import syntax.
     * Role targeting (`all`, `control`, `workers`, `<alias>`).
     * Sudo elevation rules.
     * Ubuntu node recipes.

5. `cli/helptext/catalog.go`:
   - Register topics: `cluster`, `cluster-import`, `cluster-exec`, `cluster-run-script`, `cluster-node`.

## Constraints
- Max function lines <= 15 (target <= 8 lines).
- Affirmative booleans only (is*, has*). No negative booleans.
- Universal AppError wrapping.
- TOTAL BAN on running `go test`, `go build`, or local runner during execution.
- Strict Unix LF line endings.
