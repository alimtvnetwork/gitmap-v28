# Subtask 02: Kubernetes Cluster Command & Join Pipeline

## Objective
Implement `gitmap cluster k8s` CLI subcommand engine with target resolution, automated token extraction on `join`, live stream output, and wire routing into `cli/cmd/cluster.go`.

## Disjoint Files Assigned
- `cli/cmdssh/cluster_k8s_cmd.go`
- `cli/cmdssh/cluster_k8s_cmd_test.go`
- `cli/cmd/cluster.go`

## Implementation Details
1. `cli/cmdssh/cluster_k8s_cmd.go`:
   - Implement `RunClusterK8sCLI(args []string) error`:
     - Subcommands:
       * `prereq <target>`
       * `install <target> [--version <ver>] [--hostname <name>]`
       * `init <target> [--pod-cidr <cidr>]`
       * `cni <target> [--plugin weave|calico]`
       * `join-command <target>`
       * `join <target> [--command "<join-command>"]`
       * `nfs <target> [--export-dir <dir>]`
       * `helm-install <target> [--version <ver>]`
       * `helm-nfs <target> [--nfs-server <ip>] [--export-dir <dir>]`
       * `status <target>`
       * `reset <target>`
     - Target resolution: calls `store.ListHostsByTarget(ctx, target, db.Conn())`.
     - Automatic Join Pipeline on `gitmap cluster k8s join <workers>`:
       * If `--command` not provided, queries control plane (`store.ListHostsByRole(ctx, "control", db.Conn())`), runs `sudo kubeadm token create --print-join-command`, extracts the command string, and executes across target worker hosts.
     - Executes recipe scripts using `ExecuteNodeCommand(ctx, host, script, true)` with sudo elevation.
     - Live stream output and summary table.

2. `cli/cmd/cluster.go`:
   - Wire routing in `dispatchClusterSubcommand`:
     * `case "k8s", "kube", "kubernetes": return cmdssh.RunClusterK8sCLI(args[1:])`

3. `cli/cmdssh/cluster_k8s_cmd_test.go`:
   - Unit tests for command parsing, flag resolution, and routing.

## Constraints
- Max function lines <= 15 (target <= 8 lines).
- Affirmative booleans only (is*, has*). No negative booleans.
- Universal AppError wrapping.
- TOTAL BAN on running `go test`, `go build`, or local runner during execution.
- Strict Unix LF line endings.
