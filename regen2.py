import os
import json

task_dir = ".lovable/plans/subtasks/ssh-login-and-join"
plan_file = ".lovable/plans/pending/01-ssh-login-and-join.md"

tasks = []
# Group 1: Database and Store (Tasks 1-10)
tasks.append({"title": "Create SQLite table for SSH Hosts", "sym": "SQLCreateSSHHosts", "file": "gitmap/store/migrations_ssh.go", "phase": "Scaffold", "deps": [], "sig": "func RegisterSSHHostMigration(db *sql.DB, version int, force bool) error", "body": "1. Add a new raw SQL string `CREATE TABLE ssh_hosts (id TEXT PRIMARY KEY, alias TEXT, ip TEXT, username TEXT, created_at DATETIME);`\n2. Register this migration in the `store.Up()` sequence.\n3. Ensure `IF NOT EXISTS` is present to avoid panic on re-run."})
tasks.append({"title": "Define Go struct for SSHHost", "sym": "SSHHost", "file": "gitmap/store/models_ssh.go", "phase": "Scaffold", "deps": [], "sig": "type SSHHost struct { ID string; Alias string; IP string }", "body": "1. Create `SSHHost` struct with fields `ID`, `Alias`, `IP`, `Username`, `CreatedAt`.\n2. Tag fields with `json:\"...\"` and `db:\"...\"`.\n3. Keep fields strictly typed (e.g., time.Time for CreatedAt)."})
tasks.append({"title": "Implement SSHHost Insert", "sym": "InsertSSHHost", "file": "gitmap/store/ssh_repo.go", "phase": "Implement", "deps": [1,2], "sig": "func InsertSSHHost(ctx context.Context, host SSHHost, tx *sql.Tx) error", "body": "1. Write `InsertSSHHost(ctx context.Context, host SSHHost) error`.\n2. Use named parameters `INSERT INTO ssh_hosts ... VALUES (:id, ...)`.\n3. Return wrapped `apperror` on constraint violation."})
tasks.append({"title": "Implement SSHHost Get by Alias", "sym": "GetHostByAlias", "file": "gitmap/store/ssh_repo.go", "phase": "Implement", "deps": [1,2], "sig": "func GetHostByAlias(ctx context.Context, alias string, db *sql.DB) (SSHHost, error)", "body": "1. Write `GetHostByAlias(ctx context.Context, alias string) (SSHHost, error)`.\n2. Query `SELECT * FROM ssh_hosts WHERE alias = ?`.\n3. Map `sql.ErrNoRows` to `apperror.ErrNotFound`."})
tasks.append({"title": "Implement SSHHost Delete", "sym": "DeleteHostByIP", "file": "gitmap/store/ssh_repo.go", "phase": "Implement", "deps": [1,2], "sig": "func DeleteHostByIP(ctx context.Context, ip string, db *sql.DB) error", "body": "1. Write `DeleteHostByIP(ctx context.Context, ip string) error`.\n2. Execute `DELETE FROM ssh_hosts WHERE ip = ?`.\n3. Return nil if 0 rows affected (idempotency)."})
tasks.append({"title": "Create SQLite table for SSH History", "sym": "SQLCreateSSHHistory", "file": "gitmap/store/migrations_ssh_hist.go", "phase": "Scaffold", "deps": [], "sig": "func RegisterSSHHistoryMigration(db *sql.DB, v int, force bool) error", "body": "1. Define `CREATE TABLE ssh_history (id TEXT PRIMARY KEY, host_ip TEXT, joined_at DATETIME, user TEXT);`.\n2. Apply migration using the standard store package pattern."})
tasks.append({"title": "Define Go struct for SSHHistory", "sym": "SSHHistory", "file": "gitmap/store/models_ssh_hist.go", "phase": "Scaffold", "deps": [], "sig": "type SSHHistory struct { ID string; HostIP string; JoinedAt time.Time }", "body": "1. Declare `SSHHistory` with `ID`, `HostIP`, `JoinedAt`, `User`.\n2. Ensure proper mapping for time formats."})
tasks.append({"title": "Implement SSHHistory Log", "sym": "LogSSHJoin", "file": "gitmap/store/ssh_hist_repo.go", "phase": "Implement", "deps": [6,7], "sig": "func LogSSHJoin(ctx context.Context, h SSHHistory, db *sql.DB) error", "body": "1. Create `LogSSHJoin(ctx context.Context, h SSHHistory) error`.\n2. Use standard SQL driver to insert the record."})
tasks.append({"title": "Implement SSHHistory List", "sym": "ListSSHHistory", "file": "gitmap/store/ssh_hist_repo.go", "phase": "Implement", "deps": [6,7], "sig": "func ListSSHHistory(ctx context.Context, limit int, offset int) ([]SSHHistory, error)", "body": "1. Create `ListSSHHistory(ctx context.Context) ([]SSHHistory, error)`.\n2. Order by `joined_at DESC`.\n3. Return empty slice (not nil) if no records."})
tasks.append({"title": "Integration Test for SSH Store", "sym": "TestSSHRepoIntegration", "file": "gitmap/store/ssh_repo_test.go", "phase": "Wire+Test", "deps": [3,4,5,8,9], "sig": "func TestSSHRepoIntegration(t *testing.T, db *sql.DB, ctx context.Context)", "body": "1. Setup in-memory SQLite (`:memory:`).\n2. Run migrations, insert a host, log history, assert counts.\n3. Use `require.NoError` for setup steps."})

# Group 2: Core SSH Utilities and Parsing (Tasks 11-20)
tasks.append({"title": "Define SSH Target Struct", "sym": "SSHTarget", "file": "gitmap/cmd/ssh_parser.go", "phase": "Scaffold", "deps": [], "sig": "type SSHTarget struct { Username string; IP string; Port int }", "body": "1. Declare `SSHTarget` with `Username` and `IP` strings.\n2. Add a `String()` method returning `username@ip` format."})
tasks.append({"title": "Implement Target Parsing", "sym": "ParseSSHTarget", "file": "gitmap/cmd/ssh_parser.go", "phase": "Implement", "deps": [11], "sig": "func ParseSSHTarget(raw string, defaultUser string, defaultPort int) (*SSHTarget, error)", "body": "1. Parse input string `user@ip` or `ip@user`.\n2. Detect if it's just an IP or alias.\n3. Return pointer to `SSHTarget` and potential error."})
tasks.append({"title": "Test SSH Target Parser", "sym": "TestParseSSHTarget", "file": "gitmap/cmd/ssh_parser_test.go", "phase": "Wire+Test", "deps": [12], "sig": "func TestParseSSHTarget(t *testing.T, raw string, expected *SSHTarget)", "body": "1. Table-driven tests for inputs: `192.168.1.9@a`, `a@192.168.1.9`, `m1`.\n2. Assert correct assignment to `Username` and `IP` fields."})
tasks.append({"title": "Define Interactive SSH Runner", "sym": "InteractiveSSHClient", "file": "gitmap/cmd/ssh_client.go", "phase": "Scaffold", "deps": [], "sig": "type InteractiveSSHClient struct { Stdin io.Reader; Stdout io.Writer; Stderr io.Writer }", "body": "1. Define interface or struct wrapping `os/exec` for SSH.\n2. Must accept `os.Stdin`, `os.Stdout`, `os.Stderr` mapping."})
tasks.append({"title": "Implement Spawn SSH", "sym": "SpawnSSH", "file": "gitmap/cmd/ssh_client.go", "phase": "Implement", "deps": [14], "sig": "func SpawnSSH(ctx context.Context, target SSHTarget, args []string) error", "body": "1. Create command `ssh target.String()`.\n2. Wire I/O streams directly to terminal.\n3. Call `cmd.Run()` and wrap `*exec.ExitError`."})
tasks.append({"title": "Implement SSH Password Fallback logic", "sym": "PromptSSHPassword", "file": "gitmap/cmd/ssh_client.go", "phase": "Implement", "deps": [14], "sig": "func PromptSSHPassword(ctx context.Context, prompt string, fd int) (string, error)", "body": "1. Use `golang.org/x/term` to prompt for password securely.\n2. (Future-proofing: return string without echoing to stdout)."})
tasks.append({"title": "Test SSH Client Config", "sym": "TestSSHClient", "file": "gitmap/cmd/ssh_client_test.go", "phase": "Wire+Test", "deps": [15,16], "sig": "func TestSSHClient(t *testing.T, target SSHTarget, expected []string)", "body": "1. Use mock or echo binary to simulate SSH spawn.\n2. Verify the arguments passed to `exec.Command` are correct."})
tasks.append({"title": "Scaffold `gitmap ip` resolver struct", "sym": "IPResolver", "file": "gitmap/cmd/ip_resolver.go", "phase": "Scaffold", "deps": [], "sig": "type IPResolver struct { Cache map[string]string; Timeout time.Duration }", "body": "1. Empty struct with method signatures for fetching local IP.\n2. Prepare cross-platform implementation structure."})
tasks.append({"title": "Implement Local IP Detection", "sym": "GetLocalIP", "file": "gitmap/cmd/ip_resolver.go", "phase": "Implement", "deps": [18], "sig": "func GetLocalIP(ctx context.Context, skipLoopback bool, ifaceName string) (string, error)", "body": "1. Iterate `net.Interfaces()` and find first non-loopback IPv4.\n2. Handle Windows, Unix, Centos gracefully."})
tasks.append({"title": "Test Local IP Detection", "sym": "TestGetLocalIP", "file": "gitmap/cmd/ip_resolver_test.go", "phase": "Wire+Test", "deps": [19], "sig": "func TestGetLocalIP(t *testing.T, skipLoopback bool, expected string)", "body": "1. Assert returned IP matches standard IPv4 regex.\n2. Ensure it does not return `127.0.0.1`."})

# Group 3: Command Dispatches - Login and Alias (Tasks 21-30)
tasks.append({"title": "Scaffold `gitmap ssh login` cmd", "sym": "runSSHLogin", "file": "gitmap/cmd/ssh_login_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runSSHLogin(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Define cobra command or custom CLI runner function.\n2. Bind args `[target]`."})
tasks.append({"title": "Wire Login with DB & Parser", "sym": "executeSSHLogin", "file": "gitmap/cmd/ssh_login_cmd.go", "phase": "Implement", "deps": [4,12,15,21], "sig": "func executeSSHLogin(ctx context.Context, target string, force bool) error", "body": "1. Call `ParseSSHTarget`.\n2. If alias, lookup via `GetHostByAlias`.\n3. Call `SpawnSSH`."})
tasks.append({"title": "Scaffold `gitmap ssh as` alias cmd", "sym": "runSSHAlias", "file": "gitmap/cmd/ssh_alias_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runSSHAlias(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Define command taking `ip` and `alias name`.\n2. Prepare help text indicating 'as' keyword support."})
tasks.append({"title": "Implement Alias Saving", "sym": "saveAliasCommand", "file": "gitmap/cmd/ssh_alias_cmd.go", "phase": "Implement", "deps": [3,23], "sig": "func saveAliasCommand(ctx context.Context, ip string, alias string) error", "body": "1. Extract IP and Alias strings.\n2. Call `InsertSSHHost`.\n3. Print success message to console."})
tasks.append({"title": "Scaffold `gitmap ssh login-install` cmd", "sym": "runSSHLoginInstall", "file": "gitmap/cmd/ssh_login_install_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runSSHLoginInstall(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Command definition to login and immediately install gitmap.\n2. Prepare `target` argument."})
tasks.append({"title": "Implement Remote Install Script Fetch", "sym": "getInstallPayload", "file": "gitmap/cmd/ssh_login_install_cmd.go", "phase": "Implement", "deps": [25], "sig": "func getInstallPayload(ctx context.Context, target string, version string) (string, error)", "body": "1. Define bash payload: `curl -fsSL https://... | bash`.\n2. Ensure it handles Unix targets gracefully."})
tasks.append({"title": "Implement Remote Install Execution", "sym": "executeRemoteInstall", "file": "gitmap/cmd/ssh_login_install_cmd.go", "phase": "Implement", "deps": [15,26], "sig": "func executeRemoteInstall(ctx context.Context, payload string, target SSHTarget) error", "body": "1. Instead of interactive terminal, pass the bash payload via stdin to `SpawnSSH`.\n2. Wait for completion and log output."})
tasks.append({"title": "Wire Login Commands to Dispatcher", "sym": "dispatchSSH", "file": "gitmap/cmd/ssh.go", "phase": "Wire+Test", "deps": [22,24,27], "sig": "func dispatchSSH(ctx context.Context, args []string, parent *cobra.Command) error", "body": "1. Add switch cases for `login`, `login-install`, and fallback for `$username@ip`.\n2. Parse raw arguments properly to handle implicit aliases."})
tasks.append({"title": "E2E Test SSH Login Dispatcher", "sym": "TestDispatchSSHLogin", "file": "gitmap/cmd/ssh_e2e_test.go", "phase": "Wire+Test", "deps": [28], "sig": "func TestDispatchSSHLogin(t *testing.T, args []string, mockCtx context.Context)", "body": "1. Call `dispatchSSH` with mock args `[]string{\"login\", \"a@b\"}`.\n2. Assert correct handlers are invoked without actually spawning SSH."})
tasks.append({"title": "Document Login Commands in Help", "sym": "appendSSHHelp", "file": "gitmap/helptext/docs/cmd/ssh.go", "phase": "Wire+Test", "deps": [28], "sig": "func appendSSHHelp(cmd *cobra.Command, args []string, buf io.Writer) error", "body": "1. Add examples for `gitmap ssh m1`.\n2. Explain alias resolution mechanism in markdown block."})

# Group 4: SSH Join and Auth Management (Tasks 31-42)
tasks.append({"title": "Scaffold `gitmap ssh-join` (sj) cmd", "sym": "runSSHJoin", "file": "gitmap/cmd/sshjoin_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runSSHJoin(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Base command handler for `ssh-join` and `sj`.\n2. Identify arguments for add/rm/ls."})
tasks.append({"title": "Implement Join Record Logic", "sym": "executeSSHJoin", "file": "gitmap/cmd/sshjoin_cmd.go", "phase": "Implement", "deps": [3,8,31], "sig": "func executeSSHJoin(ctx context.Context, target string, history SSHHistory) error", "body": "1. When joining, record the host via `InsertSSHHost`.\n2. Log event via `LogSSHJoin`.\n3. Print 'Joined successfully'."})
tasks.append({"title": "Scaffold `gitmap sj rm` cmd", "sym": "runSJRm", "file": "gitmap/cmd/sshjoin_rm_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runSJRm(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Handler for `remove machine-alias` or `rm ip`.\n2. Validate input argument count."})
tasks.append({"title": "Implement SJ Remove Logic", "sym": "executeSJRm", "file": "gitmap/cmd/sshjoin_rm_cmd.go", "phase": "Implement", "deps": [5,33], "sig": "func executeSJRm(ctx context.Context, target string, force bool) error", "body": "1. Call `DeleteHostByIP` or by Alias.\n2. Print 'Machine removed'."})
tasks.append({"title": "Scaffold `gitmap sj ls` cmd", "sym": "runSJLs", "file": "gitmap/cmd/sshjoin_ls_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runSJLs(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Command for listing joined machines.\n2. Needs to access the SQLite store."})
tasks.append({"title": "Implement SJ List Output", "sym": "printSJList", "file": "gitmap/cmd/sshjoin_ls_cmd.go", "phase": "Implement", "deps": [4,35], "sig": "func printSJList(ctx context.Context, out io.Writer, max int) error", "body": "1. Fetch all hosts (requires new store method `ListHosts`).\n2. Format using `text/tabwriter` for columns: Name, IP, ID."})
tasks.append({"title": "Scaffold `gitmap sj history` cmd", "sym": "runSJHistory", "file": "gitmap/cmd/sshjoin_hist_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runSJHistory(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Command for viewing join timeline.\n2. Bind `history` sub-argument."})
tasks.append({"title": "Implement SJ History Print", "sym": "printSJHistory", "file": "gitmap/cmd/sshjoin_hist_cmd.go", "phase": "Implement", "deps": [9,37], "sig": "func printSJHistory(ctx context.Context, out io.Writer, filter string) error", "body": "1. Call `ListSSHHistory`.\n2. Print 'Who joined and when' in a clean tabular format."})
tasks.append({"title": "Scaffold `gitmap sj add-auth`", "sym": "runSJAddAuth", "file": "gitmap/cmd/sshjoin_auth_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runSJAddAuth(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Define command structure for injecting public key to remote host.\n2. Handle `$ip@$user` targeting."})
tasks.append({"title": "Implement Read Local PubKey", "sym": "getLocalPublicKey", "file": "gitmap/cmd/sshjoin_auth_cmd.go", "phase": "Implement", "deps": [39], "sig": "func getLocalPublicKey(ctx context.Context, keyPath string, parse bool) (string, error)", "body": "1. Read `~/.ssh/id_rsa.pub` or similar.\n2. Throw `E_NOT_FOUND` if key missing."})
tasks.append({"title": "Implement Remote Sudo Key Append", "sym": "appendKeyRemote", "file": "gitmap/cmd/sshjoin_auth_cmd.go", "phase": "Implement", "deps": [15,40], "sig": "func appendKeyRemote(ctx context.Context, pubKey string, target SSHTarget) error", "body": "1. Construct `mkdir -p ~/.ssh && echo 'key' >> ~/.ssh/authorized_keys`.\n2. Execute via `SpawnSSH` (potentially requiring sudo prompt)."})
tasks.append({"title": "Wire SJ commands to dispatcher", "sym": "dispatchSJ", "file": "gitmap/cmd/root.go", "phase": "Wire+Test", "deps": [32,34,36,38,41], "sig": "func dispatchSJ(ctx context.Context, args []string, root *cobra.Command) error", "body": "1. Bind `ssh-joined`, `ssh-join`, and `sj` to the same sub-router.\n2. Route `ls`, `rm`, `history`, `add-auth` properly."})

# Group 5: Network IP Management (Tasks 43-50)
tasks.append({"title": "Scaffold `gitmap ip` cmd implementation", "sym": "runIPCmd", "file": "gitmap/cmd/ip_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runIPCmd(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Create CLI entrypoint for returning IP.\n2. Use standard gitmap command pattern."})
tasks.append({"title": "Wire `gitmap ip` logic", "sym": "executeIPCmd", "file": "gitmap/cmd/ip_cmd.go", "phase": "Implement", "deps": [19,43], "sig": "func executeIPCmd(ctx context.Context, skipLoopback bool, w io.Writer) error", "body": "1. Call `GetLocalIP`.\n2. Print to stdout.\n3. Return error if loopback only."})
tasks.append({"title": "Scaffold `gitmap ip-change` cmd", "sym": "runIPChangeCmd", "file": "gitmap/cmd/ipchange_cmd.go", "phase": "Scaffold", "deps": [], "sig": "func runIPChangeCmd(cmd *cobra.Command, args []string, ctx context.Context) error", "body": "1. Define argument structure `[new-ip]`.\n2. Require root/admin privileges note."})
tasks.append({"title": "Implement OS specific IP swap (Linux)", "sym": "swapIPLinux", "file": "gitmap/cmd/ipchange_linux.go", "phase": "Implement", "deps": [45], "sig": "func swapIPLinux(ctx context.Context, oldIP string, newIP string) error", "body": "1. Use `ip addr add` / `ip addr del` commands.\n2. Needs to capture previous IP state for rollback."})
tasks.append({"title": "Implement OS specific IP swap (Windows)", "sym": "swapIPWindows", "file": "gitmap/cmd/ipchange_windows.go", "phase": "Implement", "deps": [45], "sig": "func swapIPWindows(ctx context.Context, interfaceName string, newIP string) error", "body": "1. Use `netsh interface ip set address` via `os/exec`.\n2. Save previous DHCP/Static state."})
tasks.append({"title": "Implement Ping Validator", "sym": "validatePing", "file": "gitmap/cmd/ipchange_cmd.go", "phase": "Implement", "deps": [45], "sig": "func validatePing(ctx context.Context, targetHost string, count int) bool", "body": "1. Exec `ping -c 3 8.8.8.8` (or `-n 3` on Windows).\n2. Return boolean indicating success."})
tasks.append({"title": "Implement Safe Rollback orchestrator", "sym": "executeIPChange", "file": "gitmap/cmd/ipchange_cmd.go", "phase": "Implement", "deps": [46,47,48], "sig": "func executeIPChange(ctx context.Context, newIP string, doPing bool) error", "body": "1. Call swap function.\n2. Call `validatePing`.\n3. If ping fails, print 'reverting' and swap back. Prompt if not automatic."})
tasks.append({"title": "Wire IP commands to dispatcher", "sym": "dispatchIP", "file": "gitmap/cmd/root.go", "phase": "Wire+Test", "deps": [44,49], "sig": "func dispatchIP(ctx context.Context, args []string, parent *cobra.Command) error", "body": "1. Expose `ip` and `ip-change` in root menu.\n2. Add complete end-to-end test validating argument parsing."})

template = """---
plan: {plan_file}
domain: Cli
phase: {phase}
target_files: ["{file}"]
depends_on: [{depends}]
citations:
  app_spec: "spec/19-ssh-executor/01-spec.md §Section"
  canonical_size: "spec/05-coding-guidelines/01-code-quality-improvement.md"
  language_guideline: "spec/05-coding-guidelines/02-go-code-style.md"
  boolean_styling: "spec/05-coding-guidelines/03-naming-conventions.md"
  folder_naming: "spec/05-coding-guidelines/05-file-project-structure.md"
  error_architecture: "spec/05-coding-guidelines/04-error-handling.md"
  error_codes: "spec/05-coding-guidelines/04-error-handling.md"
  logging_traces: "spec/05-coding-guidelines/07-logging-observability.md"
  response_envelope: "spec/05-coding-guidelines/10-api-design.md"
  golden_fixture: "spec/08-json-schemas/ssh-list.schema.json"
  strictly_avoid: ".lovable/strictly-avoid.md"
  database: "spec/05-coding-guidelines/11-database-patterns.md"
  ui_surface: "n/a — cli tool"
  tests: "unit Test{sym}"
  ci_cd_guard: ".github/workflows/ci.yml"
  ambiguity: "n/a"
  issue_rca: "n/a"
---
# Task {idx:03d} — {title}

## 1. Learn
- [SSH Commands](file:///d:/work/gitmap/.lovable/spec/commands/01-ssh-commands.md) — Why: Defines required behaviour.
- [App Error Docs](file:///d:/work/gitmap/spec/05-coding-guidelines/04-error-handling.md) — Why: Standards for returning results.
- [{file}](file:///d:/work/gitmap/{file}) — Why: Target file.

## 2. Goal
Deliver the {phase} step for `{sym}` to support the {title} feature. This is isolated logic for the SSH/IP subdomains.

## 3. Inputs and Contracts
- Types: `string`, `context.Context`
- Outputs: `error`
- Codes: `E_INTERNAL_ERROR`
- Signature:
  ```go
  {sig}
  ```

## 4. Execute
{body}

## 5. Constraints
- **Canonical Size**: spec/05-coding-guidelines/01-code-quality-improvement.md.
- **Error Types**: Must use `apperror`.
- **No Globals**: .lovable/strictly-avoid.md.

## 6. Verify
```bash
go test ./... -v -run {sym}
```
Expected output:
```text
PASS
```

## 7. Done When
- [ ] 1. `{sym}` is fully functional.
- [ ] 2. Tests pass successfully.
- [ ] 3. No canonical size violations exist.

## 8. Notes and Open Questions
None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.
"""

for i, u in enumerate(tasks, 1):
    dep_list = []
    for d in u["deps"]:
        dep_list.append(f"Task {d:03d}")
    depends_str = ", ".join(dep_list) if dep_list else "None"
    
    content = template.format(
        idx=i, title=u["title"], phase=u["phase"], sym=u["sym"], sig=u["sig"],
        file=u["file"], plan_file=plan_file, depends=depends_str,
        body=u["body"]
    )
    with open(os.path.join(task_dir, f"{i:03d}-task.md"), "w", encoding="utf-8") as f:
        f.write(content)
print("Regenerated 50 highly unique tasks with signatures.")
