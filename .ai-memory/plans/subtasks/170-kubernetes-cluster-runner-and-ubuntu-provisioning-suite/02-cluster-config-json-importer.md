# Subtask 02: Cluster Config JSON Importer

## Objective
Implement `gitmap cluster import <config.json>` (and `gitmap sj import-cluster <config.json>`) to ingest cluster topologies (`01-config.json`), encrypt node passwords via SSH RSA-OAEP with SHA-256 (`~/.ssh/id_rsa`), and persist nodes into SQLite with designated roles (`control`, `worker`).

## Disjoint Files Assigned
- `cli/cmdssh/cluster_config.go`
- `cli/cmdssh/cluster_import_cmd.go`
- `cli/cmdssh/cluster_import_cmd_test.go`
- `cli/cmdssh/sshjoin_cmd.go`

## Implementation Details
1. `cli/cmdssh/cluster_config.go`:
   - Define JSON struct mapping to `01-config.json`:
     ```go
     type ClusterConfigJSON struct {
         Control map[string]string `json:"control"`
         Nodes   map[string]string `json:"nodes"`
         User    ClusterUserJSON   `json:"user"`
     }
     type ClusterUserJSON struct {
         Name     string `json:"name"`
         Password string `json:"password"`
     }
     ```
   - Implement `LoadClusterConfigFile(path string) (*ClusterConfigJSON, error)` with file existence and JSON validation.

2. `cli/cmdssh/cluster_import_cmd.go`:
   - Implement `RunClusterImportCLI(args []string) error`:
     - Resolves file path (defaults to `./01-config.json` if omitted).
     - Loads JSON topology.
     - Validates `user.name` and `user.password`.
     - Encrypts password using `EncryptSSHPassword(cfg.User.Password)`.
     - Iterates `cfg.Control`:
       - Enrolls with alias = key (e.g. `master`), IP = val, user = `cfg.User.Name`, role = `"control"`.
     - Iterates `cfg.Nodes`:
       - Enrolls with alias = key (e.g. `worker-1`), IP = val, user = `cfg.User.Name`, role = `"worker"`.
     - Persists all hosts using `store.EnrollSSHHost(ctx, host, history, db)`.
     - Prints formatted terminal summary table of all imported nodes, showing role, alias, IP, username, and encrypted RSA status.

3. `cli/cmdssh/sshjoin_cmd.go`:
   - Register subcommand `import-cluster` (aliases: `import`, `import-config`) on `SSHJoinCmd`.
   - Wire routing in `routeSJSubcommands`.

4. `cli/cmdssh/cluster_import_cmd_test.go`:
   - Add unit tests for JSON parsing, RSA encryption validation, and host enrollment.

## Constraints
- Max function lines <= 15 (target <= 8 lines).
- Affirmative booleans only (is*, has*). No negative booleans.
- Universal AppError wrapping.
- TOTAL BAN on running `go test`, `go build`, or local runner during execution.
- Strict Unix LF line endings.
