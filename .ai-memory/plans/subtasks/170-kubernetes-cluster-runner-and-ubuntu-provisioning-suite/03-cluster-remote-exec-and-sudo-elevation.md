# Subtask 03: Cluster Remote Exec & Sudo Elevation

## Objective
Implement `gitmap cluster exec <target> "<command>" [--sudo]` and `gitmap cluster run-script <target> <script.sh> [--sudo]` with remote execution, in-memory RSA password decryption, `sudo -S` password elevation, parallel concurrency, and live line tagging.

## Disjoint Files Assigned
- `cli/cmdssh/cluster_runner.go`
- `cli/cmdssh/cluster_exec_cmd.go`
- `cli/cmdssh/cluster_exec_cmd_test.go`
- `cli/cmdssh/cluster_script_cmd.go`

## Implementation Details
1. `cli/cmdssh/cluster_runner.go`:
   - Define `ClusterRunResult` struct: `Host store.SSHHost`, `ExitCode int`, `Stdout string`, `Stderr string`, `Duration time.Duration`, `Err error`.
   - Implement `ExecuteNodeCommand(ctx context.Context, host store.SSHHost, shellCmd string, isSudo bool) ClusterRunResult`:
     - Decrypts password using `DecryptSSHPassword(host.EncryptedPassword)`.
     - Formats command:
       - If `isSudo`: `echo '<decrypted>' | sudo -S bash -c '<escaped_cmd>'`
       - Else: `bash -c '<escaped_cmd>'`
     - Connects via OpenSSH using AskPass (`injectAskPassEnv`) or stdin piping.
     - Live line streaming: prefixes each line with `[<alias>]` (e.g. `[master]`, `[worker-1]`).
   - Implement `DispatchClusterRun(ctx context.Context, hosts []store.SSHHost, shellCmd string, isSudo bool, concurrency int) []ClusterRunResult`:
     - Dispatches across worker pool with bounded concurrency.
     - Collects results and prints clean summary table.

2. `cli/cmdssh/cluster_exec_cmd.go`:
   - Implement `RunClusterExecCLI(args []string) error`:
     - Syntax: `gitmap cluster exec <all|control|workers|<alias>> "<command>" [flags]`
     - Parses `--sudo` / `-s` flag.
     - Parses `--parallel` / `-p` flag (default 4).
     - Resolves target nodes via `store.ListHostsByTarget(ctx, target, db)`.
     - Calls `DispatchClusterRun`.

3. `cli/cmdssh/cluster_script_cmd.go`:
   - Implement `RunClusterScriptCLI(args []string) error`:
     - Syntax: `gitmap cluster run-script <all|control|workers|<alias>> <script-path> [flags]`
     - Reads local script file content.
     - For each target node:
       - Uploads script to `/tmp/on-the-fly-cmd/script-<timestamp>.sh` using base64 echo or scp.
       - Runs `chmod +x` on the remote script.
       - Executes remote script with sudo if `--sudo` requested.
       - Cleans up the remote script file.

4. `cli/cmdssh/cluster_exec_cmd_test.go`:
   - Unit tests for argument parsing, target group expansion, and command formatting.

## Constraints
- Max function lines <= 15 (target <= 8 lines).
- Affirmative booleans only (is*, has*). No negative booleans.
- Universal AppError wrapping.
- TOTAL BAN on running `go test`, `go build`, or local runner during execution.
- Strict Unix LF line endings.
