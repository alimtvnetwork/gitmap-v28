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
| copy       | cp    | Transfer file to remote machines with path macro expansion |
| mv         | move  | Move file from host to remote machines and remove source |
| known-hosts| kh    | Manage, list, trust, and sync ~/.ssh/known_hosts & SQLite DB |
| trust      |       | Auto-scan & trust remote host key in known_hosts and DB |
| untrust    |       | Untrust and remove a machine from known_hosts and DB |
| fix-auth   |       | Deploy SSH public key to remote authorized_keys (Unix/Windows) |
| macro      | m     | Synchronize, export, or import macros across SSH machines |
| export-oneliner | eo | Generate single-line import command and copy to clipboard |
| deploy     |       | Deploy node-config or mesh public keys across all fleet nodes |
| export-all |       | Export all settings, macros, config & SSH data from local to nodes |
| import-all |       | Import all settings, macros, config & SSH data from a remote node |
| error-logs | err, errors, logs | Query and display error logs and execution trace from last SSH operations |

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

    # Enroll multiple machines simultaneously (comma- or space-separated):
    $ gitmap ssh join 192.168.1.10,192.168.1.11,192.168.1.12
    $ gitmap ssh join 192.168.1.10 192.168.1.11 192.168.1.12

### Enroll with encrypted password (add-with-pass)

    $ gitmap ssh join add-with-pass alim@192.168.1.14 secret123 devbox
      ✓ Machine 'devbox' (alim@192.168.1.14) joined successfully with password.
        Password encrypted and stored securely using SSH RSA key.
        Recall anytime: gitmap ssh devbox (auto-login via AskPass)

### Execute remote command with automatic liveness checks (exec / se)

Execute remote commands across nodes. Any GitMap function or verb can be executed directly without repeating `gitmap`:

    # Run any GitMap verb directly without mentioning 'gitmap':
    $ gitmap ssh exec status
    $ gitmap ssh exec open google.com
    $ gitmap ssh exec macro sync --all
    $ gitmap ssh exec pull-all
    $ gitmap ssh exec doctor

    # Standard shell commands or explicit gitmap invocations also work:
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

    # Check multiple targeted machines (comma- or space-separated):
    $ gitmap ssh check devbox,192.168.1.20
    $ gitmap ssh check 127.0.0.1 192.168.1.50

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

### Known Hosts Management (`known-hosts` / `kh` / `trust` / `untrust`)

List all tracked known hosts in a terminal UI table:

    $ gitmap ssh known-hosts ls
    Known SSH Hosts (3 tracked):
    HOST / IP              KEY TYPE       FINGERPRINT                                        UPDATED
    ----------------------------------------------------------------------------------------------------------
    192.168.1.5            ssh-ed25519    SHA256:ROdbMLgtCZ47lGGUryHnDCpik9/G5V8H/ie6tuqai2k 2026-09-20 01:21
    192.168.1.9            ssh-ed25519    SHA256:ROdbMLgtCZ47lGGUryHnDCpik9/G5V8H/ie6tuqai2k 2026-09-20 01:21
    github.com             ssh-ed25519    SHA256:+DiY3wvvV6TuJJhbpZisF/zPTOZ736+Fnkas3BjZ7DA 2026-09-18 10:00

Auto-scan and trust a new host without interactive prompts:

    $ gitmap ssh trust 192.168.1.5
      ✓ Host '192.168.1.5' trusted successfully!
        Key Type:    ssh-ed25519
        Fingerprint: SHA256:ROdbMLgtCZ47lGGUryHnDCpik9/G5V8H/ie6tuqai2k

Untrust and remove a host from known_hosts and SQLite DB:

    $ gitmap ssh untrust 192.168.1.5
      ✓ Removed '192.168.1.5' from known_hosts and database.

Synchronize `~/.ssh/known_hosts` into the SQLite database:

    $ gitmap ssh known-hosts sync
      ✓ Synchronized 3 known host(s) to database.

### Deploy Authorized Keys to Remote Hosts (`fix-auth`)

Deploy your local public SSH key to remote `authorized_keys` for single or multiple machines:

    $ gitmap ssh fix-auth machineid, ip, id
    $ gitmap ssh fix-auth devbox --unix

### Export All Settings to Nodes (`export-all`)

Export all settings, macros, configuration, SSH connections, and known hosts from the current machine to all online nodes (or targeted nodes):

    # Export everything to all online nodes:
    $ gitmap ssh export-all nodes
    $ gitmap ssh export-all all

    # Export to specific nodes with overwrite:
    $ gitmap ssh export-all devbox,worker-1 --force

    # Simulate export without making changes:
    $ gitmap ssh export-all --dry-run

### Import All Settings from a Node (`import-all`)

Import all settings, macros, configuration, SSH connections, and known hosts from a specified remote node to the current machine:

    # Import by node name, IP, or alias:
    $ gitmap ssh import-all node "devbox"
    $ gitmap ssh import-all node "192.168.1.5"
    $ gitmap ssh import-all "main"

    # Import with overwrite:
    $ gitmap ssh import-all node "devbox" --force

    # Simulate import without modifying local files:
    $ gitmap ssh import-all node "devbox" --dry-run

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
