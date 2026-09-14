# Subtask 01: SQLite Cluster Roles & Target Resolver

## Objective
Update the SQLite `ssh_hosts` table and repository layer to support cluster roles (`control`, `worker`) and implement target group resolution (`all`, `control`, `workers`, individual aliases/IPs).

## Disjoint Files Assigned
- `cli/store/models_ssh.go`
- `cli/store/migrations_ssh.go`
- `cli/store/ssh_repo.go`
- `cli/store/ssh_repo_test.go`

## Implementation Details
1. `cli/store/models_ssh.go`:
   - Add field `ClusterRole string `json:"cluster_role,omitempty" db:"cluster_role"`` to `SSHHost`.

2. `cli/store/migrations_ssh.go`:
   - Update `SQLCreateSSHHostsTable` to include `cluster_role TEXT DEFAULT 'worker'`.
   - Update `ensureHostColumns(db *sql.DB)`:
     ```go
     _, _ = db.Exec("ALTER TABLE ssh_hosts ADD COLUMN cluster_role TEXT DEFAULT 'worker';")
     ```

3. `cli/store/ssh_repo.go`:
   - Update `sqlInsertSSHHost` to insert `cluster_role`.
   - Update `buildHostNamedArgs` to bind `cluster_role` using `resolveClusterRole(host.ClusterRole)`.
   - Update `updateHostForIP` and `updateHostForAlias` to update `cluster_role`.
   - Update `sqlSelectHostFields` to select `COALESCE(cluster_role, 'worker')`.
   - Update `scanHostRows` to scan `&host.ClusterRole`.
   - Update `GetHostByAlias`, `GetHostByIP`, `GetHostByID` to scan `&host.ClusterRole`.
   - Implement `ListHostsByRole(ctx context.Context, role string, db *sql.DB) ([]SSHHost, error)`:
     - Queries `WHERE cluster_role = ?`.
   - Implement `ListHostsByTarget(ctx context.Context, target string, db *sql.DB) ([]SSHHost, error)`:
     - Handles `"all"`, `""` -> calls `ListHosts(ctx, db)`.
     - Handles `"control"`, `"master"` -> calls `ListHostsByRole(ctx, "control", db)`.
     - Handles `"workers"`, `"worker"`, `"nodes"` -> calls `ListHostsByRole(ctx, "worker", db)`.
     - Checks if `target` matches an alias (`GetHostByAlias`). If found, returns `[]SSHHost{host}`.
     - Checks if `target` matches an IP (`GetHostByIP`). If found, returns `[]SSHHost{host}`.
     - If comma-separated, splits by comma and accumulates.

4. `cli/store/ssh_repo_test.go`:
   - Add unit tests for `ClusterRole` persistence, `ListHostsByRole`, and `ListHostsByTarget`.

## Constraints
- Max function lines <= 15 (target <= 8 lines).
- Affirmative booleans only (is*, has*). No negative booleans.
- Universal AppError wrapping.
- TOTAL BAN on running `go test`, `go build`, or local runner during execution.
- Strict Unix LF line endings.
