# 01-ssh-commands: SSH & Network Utilities Specification

## 1. Overview
This specification details the behavior, cross-platform implementation, and database architecture for the new `gitmap ssh`, `gitmap ssh-join` (`sj`), and `gitmap ip` command suites.

## 2. Command Definitions and Routing

### 2.1 SSH Login (`gitmap ssh login`)
**Behavior:** Initiates an interactive SSH session.
- Supports both `user@ip` and `ip@user` string formats interchangeably.
- If called via `gitmap ssh <target>`, it routes to `login` automatically if `<target>` contains an `@` or matches a known alias.
- Securely prompts for a password using `golang.org/x/term` if SSH key authentication fails.
- Caches successful authentication metadata (username, IP, alias) into the local SQLite database.

### 2.2 SSH Aliases (`gitmap ssh <ip> as <alias>`)
**Behavior:** Registers a shorthand alias for an IP/User combination.
- Usage: `gitmap ssh 192.168.1.9 as m1`
- Resolves subsequent calls to `gitmap ssh m1` by querying the SQLite `ssh_hosts` table.

### 2.3 Remote Installation (`gitmap ssh login-install <target>`)
**Behavior:** Connects to the remote machine and executes the gitmap installation script.
- Pipes the bootstrap script via stdin over the SSH transport.
- Triggers standard installation without requiring interactive user input, aside from initial authentication.

### 2.4 SSH Join and Tracking (`gitmap ssh-join` / `gitmap sj`)
**Behavior:** Manages a fleet of known SSH machines and tracks connection history.
- `gitmap sj add <target>` (or just `gitmap sj <target>`): Joins the machine, stores its ID/IP in SQLite, and performs a first-time install if gitmap is missing.
- `gitmap sj ls`: Lists all joined machines (Name/Alias, IP, ID) in a tabulated format.
- `gitmap sj rm <ip|alias>`: Removes the machine from the tracked database.
- `gitmap sj history`: Displays a chronological log of when users joined machines.

### 2.5 Key Distribution (`gitmap sj add-auth <target>`)
**Behavior:** Pushes the local host's SSH public key to the guest machine.
- Attempts password-based login to append `~/.ssh/id_rsa.pub` to the guest's `~/.ssh/authorized_keys`.
- Handles `sudo` privilege escalation safely if writing to protected directories.

### 2.6 IP Management (`gitmap ip` / `gitmap ip-change`)
**Behavior:** Cross-platform network utility tools.
- `gitmap ip`: Detects and prints the primary non-loopback IPv4 address. Adapts to Windows (`ipconfig`/WMI), macOS/Unix (`ifconfig`/`ip addr`), and Debian/CentOS quirks.
- `gitmap ip-change <new-ip>`: Modifies the host's IP address.
  - Implements a safety rollback orchestrator: Swaps the IP, pings `8.8.8.8` 3 times.
  - If the ping fails, automatically reverts to the previous IP configuration to prevent network lockouts.

## 3. Database Schema (SQLite)
All data is stored in the local `gitmap` SQLite instance.

### `ssh_hosts` Table
- `id` (TEXT PRIMARY KEY)
- `alias` (TEXT UNIQUE)
- `ip` (TEXT)
- `username` (TEXT)
- `created_at` (DATETIME)

### `ssh_history` Table
- `id` (TEXT PRIMARY KEY)
- `host_ip` (TEXT, foreign key to ssh_hosts)
- `joined_at` (DATETIME)
- `user` (TEXT)

## 4. End-to-End Testing & Validation
Every command must have:
1. **Unit Tests:** Validating argument parsers (e.g., `user@ip` vs `ip@user`).
2. **Integration Tests:** Verifying SQLite persistence (`save alias`, `history` queries).
3. **E2E Tests:** Stubbing the `os/exec` SSH binary to assert the correct arguments and scripts are piped to the remote host.
