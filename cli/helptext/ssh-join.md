# gitmap ssh-join

Enroll remote machines into GitMap SSH host registry, register memorable aliases, authorize SSH public keys, and recall hosts instantly without hostname collision or OpenSSH errors.

## Aliases

```bash
gitmap sj
gitmap ssh join
```

## Usage

```bash
gitmap ssh-join <user@ip|ip> [alias] [flags]
gitmap ssh-join add <user@ip|ip> [alias] [flags]
gitmap sj <user@ip|ip> [alias] [flags]
gitmap ssh join <user@ip|ip> [alias] [flags]
gitmap sj <subcommand> [flags]
```

## Description

`gitmap ssh-join` (or shortcut `gitmap sj`) registers SSH target machines into the local SQLite database (`ssh_hosts` and `ssh_history`). Once joined, you can connect directly using the alias via `gitmap ssh <alias>`, connect by IP via `gitmap ssh <ip>` (which automatically uses your enrolled user), or execute remote commands via `gitmap ssh exec <alias> <cmd>`.

Target format accepts:
- User and IP: `alim@192.168.1.14` or `root@10.0.0.12` (stores enrolled user, default alias `host-<ip>`)
- Plain IP: `192.168.1.50` (uses current user, default alias `host-192.168.1.50`)
- IP with Port: `192.168.1.50:2222`
- User, IP, and Port: `admin@10.0.0.15:2222`

`gitmap ssh join` is an internal GitMap command and is never forwarded to OpenSSH. If an unknown alias is queried, GitMap provides mistake recovery guidance with registered machines and exact copy-pasteable join instructions.

## Subcommands

| Subcommand | Alias | Description |
|------------|-------|-------------|
| add | enroll, join | Enroll remote SSH machine with alias and optional key authorization |
| scan | find, discover, probe | Scan local subnet or CIDR network for active SSH machines on port 22 |
| status | ping, health, check | Inspect connectivity, latency, and reachability of SSH machines |
| ls | list | List all enrolled SSH machines and aliases |
| rm | remove, delete | Remove an enrolled machine from registry by alias or IP |
| add-auth | auth | Authorize local public key on the target machine |
| history | hist | Display enrollment history audit logs |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| --user | -u | current user | Remote SSH username (fallback: `$USER`, `$USERNAME`, `root`) |
| --name, --alias | -n | `host-<ip>` | Memorable alias name for host recall |
| --port | -p | 22 | Target SSH port |
| --auth | | false | Push local public key to remote `~/.ssh/authorized_keys` |
| --force | -f | false | Overwrite existing alias or host mapping |
| --json | | false | Output result in JSON format for scripting |
| --help | -h | false | Show command help |

## Examples

### 1. Join Machine by User and IP (user@ip)

Enroll remote machines directly using `user@ip`:

```bash
# Join user@ip directly (alias defaults to host-<ip>)
gitmap ssh-join alim@192.168.1.14

# Join user@ip with memorable alias
gitmap ssh-join alim@192.168.1.14 devbox

# Join user@ip with sj shortcut
gitmap sj root@192.168.1.14 prod-server

# Join user@ip and push SSH public key in one step
gitmap ssh-join alim@192.168.1.14 devbox --auth

# Join user@ip using explicit add subcommand
gitmap ssh-join add dev@192.168.1.50 devbox
```

### 2. Join Machine by Plain IP

Enroll a server using default current OS username and generated alias:

```bash
gitmap ssh-join 192.168.1.50
gitmap sj 192.168.1.50
```

### 3. Join Machine with Memorable Alias

Assign a custom alias directly as the second argument:

```bash
gitmap ssh-join 192.168.1.50 prod-db
gitmap sj 192.168.1.50 prod-db
```

### 4. Join with Custom Port

Target non-standard SSH ports:

```bash
gitmap sj 192.168.1.100:2222 staging
```

### 5. Join Using Flags

Configure user, port, and alias using explicit flags:

```bash
gitmap sj 10.0.0.5 --user deploy --port 2222 --name worker-1
```

### 6. Join and Push Public Key Authorization

Automatically append local SSH public key (`~/.ssh/id_rsa.pub` or `~/.ssh/id_ed25519.pub`) to remote `~/.ssh/authorized_keys`:

```bash
gitmap sj 192.168.1.50 prod-db --auth
```

### 7. Instant Host Recall

Connect directly to any enrolled machine using its alias:

```bash
gitmap ssh prod-db
gitmap ssh dev-box
```

### 8. List Enrolled Machines

View all tracked SSH machines, aliases, users, and ports:

```bash
gitmap sj ls
gitmap ssh join ls
```

### 9. Authorize SSH Key on an Enrolled Machine

Push public key credentials to a machine that is already enrolled:

```bash
gitmap sj add-auth prod-db
```

### 10. Remote Command Execution

Run commands on enrolled hosts using aliases:

```bash
gitmap ssh exec prod-db "uptime"
gitmap se dev-box "df -h"
```

### 11. Remove Machine from Registry

Unregister an alias or IP:

```bash
gitmap sj rm prod-db
```

### 12. View Enrollment History

Inspect timestamps and audit trail of joined machines:

```bash
gitmap sj history
```

### 13. Subnet Network Discovery & SSH Scanning

Scan the local `/24` subnet or a custom CIDR network for machines listening on SSH port 22:

```bash
# Auto-detect local interface subnet and probe port 22
gitmap sj scan

# Scan specific CIDR block
gitmap sj scan 192.168.1.0/24

# Probe custom port with increased timeout and worker pool
gitmap sj find 10.0.0.0/24 --port 2222 --timeout 1s --workers 50
```

Discovered machines are cross-referenced against enrolled hosts, indicating `[ENROLLED: alias]` or `[NEW]` alongside copy-pasteable enrollment suggestions.

### 14. Machine Connectivity Health & Latency Ping

Inspect connectivity, round-trip latency, and online/offline status of machines:

```bash
# Probe all registered SSH machines in the registry
gitmap sj status

# Probe an individual machine by alias or IP
gitmap sj status devbox
gitmap sj ping 192.168.1.50

# Inspect non-standard port with custom timeout
gitmap sj health 10.0.0.12 --port 2222 --timeout 2s
```

Offline machines display specific failure reasons (e.g. `connection refused` or `connection timed out`).

## Mistake Recovery

If you attempt to connect to an unknown alias:

```bash
gitmap ssh db-staging
```

GitMap intercepts the command before calling OpenSSH, displays an error envelope, lists all currently enrolled machines, and provides the exact command needed to enroll the host:

```text
Error: alias 'db-staging' not found in SSH registry.

Available hosts:
  • prod-db (192.168.1.50:22)
  • dev-box (10.0.0.12:22)

To enroll this machine:
  gitmap ssh-join user@<ip> db-staging
  gitmap ssh-join <ip> db-staging
  gitmap sj user@<ip> db-staging
```

## See Also

- `gitmap ssh` - SSH key generation and managed configuration
- `gitmap ssh exec` - Execute remote commands on enrolled hosts
- `gitmap cluster` - Cluster-wide SSH node delegation

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter ssh-join
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
