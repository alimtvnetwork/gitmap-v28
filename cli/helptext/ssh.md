# SSH Key Management

Manages SSH key pairs for Git authentication.

## Usage

    gitmap ssh [subcommand] [flags]

## Subcommands

| Subcommand | Alias | Description                               |
|------------|-------|-------------------------------------------|
| *(none)*   |       | Generate a new SSH key pair               |
| cat        |       | Display the public key for a named key    |
| list       | ls    | List all stored SSH keys                  |
| delete     | rm    | Delete a key record (optionally files)    |
| config     |       | Regenerate ~/.ssh/config managed entries  |
| join       | sj    | Enroll machine by user@ip or IP into SSH registry |
| login      |       | Connect to a host by alias or user@host   |
| as         |       | Create an SSH alias mapping for a host IP |
| exec       | se    | Execute a remote command on a target host |
| scan       |       | Probe and test liveness across SSH fleet nodes |
| check      | health, ping | Check connectivity, open port 22, and health across SSH machines |
| install    | i     | Install or update GitMap on remote node(s) |
| update     | u     | Update GitMap binary across remote node(s) |
| agy        |       | Run Antigravity CLI or open remote folder |
| code       |       | Open remote folder in VS Code via SSH Remote |
| compare    | matrix| Display comparison table: SSH vs Cluster vs SC |

## Flags (generate)

| Flag      | Short | Description                            | Default        |
|-----------|-------|----------------------------------------|----------------|
| --name    | -n    | Label for the key in the database      | default        |
| --path    | -p    | File path for the private key          | ~/.ssh/id_rsa  |
| --email   | -e    | Email comment for the key              | git global     |
| --force   | -f    | Skip prompt if key already exists      | false          |
| --host    | -H    | Git provider hostname                  | github.com     |
| --confirm |       | Require explicit yes before generating | false          |

## Flags (list)

| Flag      | Short | Description                            |
|-----------|-------|----------------------------------------|
| --json    |       | Output keys as JSON for scripting      |

## Flags (delete)

| Flag      | Short | Description                            |
|-----------|-------|----------------------------------------|
| --name    | -n    | Name of the key to delete              |
| --files   |       | Also delete key files from disk        |

## Flags (clone integration)

| Flag      | Short | Description                            |
|-----------|-------|----------------------------------------|
| --ssh-key | -K    | SSH key name to use for cloning        |

## Prerequisites

`ssh-keygen` must be available on PATH (included with OpenSSH).

## Examples

### Generate a default SSH key

    $ gitmap ssh
      Generating public/private rsa key pair.
      ✓ SSH key "default" generated
        Path:        ~/.ssh/id_rsa
        Fingerprint: SHA256:abc123...
        Public key:

      ssh-rsa AAAA... user@example.com

      ℹ  Copy the public key above and add it to your Git provider.

### Generate a named key for work

    $ gitmap ssh --name work --path ~/.ssh/id_rsa_work
      ✓ SSH key "work" generated
        Path:        ~/.ssh/id_rsa_work
        Fingerprint: SHA256:def456...

### Display the public key (`view` / `v` / `cat`)

    $ gitmap ssh view --name work
    ssh-rsa AAAA... user@example.com

### Copy the public key to the clipboard (`copy` / `cp`)

    $ gitmap ssh copy --name work
    ssh-rsa AAAA... user@example.com
      ✓ Public key copied to clipboard (412 bytes)

Uses `clip` on Windows, `pbcopy` on macOS, and `wl-copy` / `xclip` / `xsel`
on Linux (first one found). If none are available, the key is still
printed and a one-line warning is emitted — never fails.

### Create a key (`create` is an explicit alias for the default)

    $ gitmap ssh create --name work --path ~/.ssh/id_rsa_work

### Clone using a specific SSH key

    $ gitmap clone repos.json --ssh-key work
      → Cloning with SSH key "work" (~/.ssh/id_rsa_work)
      ✓ Cloned 5 repos, 0 failed.

### List all stored keys

    $ gitmap ssh list

      SSH Keys (2):

      Name            Path                            Fingerprint                Created
      default         ~/.ssh/id_rsa                   SHA256:abc123...           2026-03-22
      work            ~/.ssh/id_rsa_work              SHA256:def456...           2026-03-22

### Enroll remote SSH machine (ssh-join / sj)

    $ gitmap ssh join alim@192.168.1.14 devbox
      ✓ Machine 'devbox' (alim@192.168.1.14) joined successfully.
        Recall anytime: gitmap ssh devbox
        Or connect directly: gitmap ssh 192.168.1.14

    # Enroll multiple machines simultaneously (comma-separated):
    $ gitmap ssh join 192.168.1.10,192.168.1.11,192.168.1.12

### Enroll with encrypted password (add-with-pass)

    $ gitmap ssh join add-with-pass alim@192.168.1.14 secret123 devbox
      ✓ Machine 'devbox' (alim@192.168.1.14) joined successfully with password.
        Password encrypted and stored securely using SSH RSA key.
        Recall anytime: gitmap ssh devbox (auto-login via AskPass)

### Execute remote command with automatic liveness checks (exec / se)

    $ gitmap ssh exec "gitmap --version"
      [devbox|192.168.1.14] gitmap version v6.260.0 linux/amd64
      [node2|192.168.1.20] OFFLINE (skipped: connection timeout)

    $ gitmap ssh exec --target devbox "uname -a"
      [devbox|192.168.1.14] Linux devbox 5.15.0-107-generic x86_64

    # Execute across multiple targeted machines (comma-separated):
    $ gitmap ssh exec devbox,worker-1 "free -m"

Nodes are probed with an in-memory TTL reachability cache (45 seconds). Offline nodes
are skipped immediately without hanging your terminal.

### Probe fleet liveness & reachability (scan)

    $ gitmap ssh scan

      SSH Fleet Liveness & Reachability Scan:

      ALIAS       IP             USER     PORT   STATUS    LATENCY   DETAILS
      devbox      192.168.1.14   alim       22   ONLINE       12ms   tcp reachable
      worker-1    192.168.1.20   ubuntu     22   OFFLINE         -   connection timeout

      Summary: 1/2 nodes online

### Check connectivity & open port 22 (check / health / ping)

    $ gitmap ssh check
      Checking SSH connectivity on port 22:
      STATUS    ALIAS       IP             USER     PORT   LATENCY   DETAILS
      ONLINE    devbox      192.168.1.14   alim       22      11ms   reachable
      OFFLINE   worker-1    192.168.1.20   ubuntu     22         -   connection timed out

    # Check a single targeted machine:
    $ gitmap ssh check devbox
      ONLINE    devbox      192.168.1.14   alim       22      10ms   reachable

    # Check multiple targeted machines (comma-separated):
    $ gitmap ssh check devbox,192.168.1.20

### Install GitMap or packages across fleet (install / i)

    $ gitmap ssh install gitmap devbox
      Installing / Updating 'gitmap' across SSH fleet (devbox):
      [devbox|192.168.1.14] gitmap missing, installing fresh...
      [devbox|192.168.1.14] Installed successfully!

    $ gitmap ssh install agy devbox
      Installing / Updating 'agy' across SSH fleet (devbox):
      [devbox|192.168.1.14] Installing package 'agy' via gitmap...
      [devbox|192.168.1.14] Installed agy successfully!

    $ gitmap ssh install gitmap all
      [devbox|192.168.1.14] gitmap found, updating to latest...
      [worker-1|192.168.1.20] OFFLINE (skipped: connection timeout)

Installs GitMap if missing or upgrades to the latest release if already present.
Also installs any package supported by GitMap (e.g. `agy`, `devbox`).

### Update GitMap or packages across fleet (update / u)

    $ gitmap ssh update gitmap all
      Updating 'gitmap' across SSH fleet (all):
      [devbox|192.168.1.14] Updated successfully!

    $ gitmap ssh update agy devbox
      Updating 'agy' across SSH fleet (devbox):
      [devbox|192.168.1.14] Updated successfully!

### Antigravity (AGY) CLI & remote folder delegation (agy)

    $ gitmap ssh agy devbox open /var/www/my-project
      ✓ Opening remote workspace in Google Antigravity: devbox:/var/www/my-project

    $ gitmap ssh agy devbox "agy --version"
      [devbox|192.168.1.14] Google Antigravity CLI v1.12.0

    # Diagnostic AppError stack trace on failure or unreachable node:
    $ gitmap ssh agy worker-1 "agy status"
      [worker-1|192.168.1.20] Offline: [EXECUTION] node [worker-1|192.168.1.20] is unreachable: connection timeout (at=cmdssh/ssh_target_nodes.go:27)
      Stack Trace:
        at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.checkRemoteNodeOnline (cmdssh/ssh_target_nodes.go:27)
        at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.establishAgyClient (cmdssh/ssh_agy_cmd.go:67)
        at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.runAgyOnNode (cmdssh/ssh_agy_cmd.go:57)
        at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeAgyOnFleet (cmdssh/ssh_agy_cmd.go:32)

### VS Code remote SSH integration (code)

    $ gitmap ssh code open /opt/app --target devbox
      ✓ Launched VS Code Remote: code --remote ssh-remote+alim@192.168.1.14 /opt/app

### Subsystems Architecture Comparison (compare / matrix)

Display comparison table between `ssh`, `cluster`, and `sc` (`servers-clients`):

    $ gitmap ssh compare

| Subsystem | Primary Focus | Join Command | Exec Command | Monitoring | Best Used When |
|---|---|---|---|---|---|
| `gitmap ssh` | Direct node management | `gitmap ssh join <u@ip>` | `gitmap ssh exec <cmd>` | `gitmap ssh scan` | Ad-hoc terminal commands, install/update, AGY/code open |
| `gitmap cluster` | Multi-node cluster orchestration | `gitmap cluster node add <ip>` | `gitmap cluster exec <target> <cmd>` | `gitmap cluster node ls` | K8s bootstrap, cluster recipes, distributed scripts |
| `gitmap sc` | Servers-clients fleet daemon | `gitmap sc join <server-url>` | `gitmap sc exec <cmd>` | `gitmap sc status` | Master-worker topology, continuous sync, live telemetry |

#### When to use which command:
- **`ssh`**: Fast, lightweight, direct command execution over standard SSH. Ideal for developer workstations, ad-hoc maintenance, and AGY/VS Code remote opening.
- **`cluster`**: Role-based infrastructure orchestration (`control` vs `workers`), provisioning recipes (Netplan IP, users, apt purge), and full Kubernetes lifecycle.
- **`sc` (`servers-clients`)**: High-speed fan-out broadcasts across entire fleet with concurrency pools, multi-shell execution, and continuous daemon synchronization.

## See Also

- `gitmap ssh-join` - Machine enrollment, alias recall, and public key authorization
- `gitmap cluster` - Multi-node orchestration, Kubernetes lifecycle, and node recipes
- `gitmap sc` - Broadcast execution and servers-clients daemon topology
- `gitmap clone` - Clone repositories from structured files
- `gitmap setup` - Configure Git global settings

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter ssh
```

The JSON schema is published at `02-spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
