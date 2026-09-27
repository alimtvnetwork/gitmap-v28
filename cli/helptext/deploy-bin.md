# gitmap ssh deploy-bin

Autonomous one-liner deployment of GitMap binary across remote SSH fleet nodes.

## Synopsis

```bash
gitmap ssh deploy-bin [target] [flags]
gitmap ssh push-bin [target] [flags]
gitmap ssh sync-bin [target] [flags]
gitmap deploy-bin [target] [flags]
```

## Description

`gitmap ssh deploy-bin` completely replaces tedious, error-prone manual `scp` commands such as:

```bash
# Old manual workflow (REQUIRES passwords, manual path entry, and StrictHostKeyChecking bypass):
scp -o StrictHostKeyChecking=no d:\work\gitmap\cli\gitmap.exe Administrator@192.168.1.3:C:/Users/Administrator/AppData/Local/gitmap-cli/gitmap.exe

# New smart GitMap one-liner (uses encrypted credentials, auto-detects OS path, streams binary, verifies version):
gitmap ssh deploy-bin w1
```

The command:
1. Resolves the active GitMap binary on the local machine (or accepts `--file <path>`).
2. Discovers target nodes from the SQLite cluster inventory (`gitmap-ssh-nodes.json` / database).
3. Connects securely via SSH key or vault-encrypted credentials with zero password prompts.
4. Detects remote operating system (Windows, Linux, macOS) and resolves standard installation paths:
   - **Windows:** `C:\Users\<user>\AppData\Local\gitmap-cli\gitmap.exe`
   - **Linux / macOS:** `/usr/local/bin/gitmap` (or `~/.local/bin/gitmap`)
5. Streams the binary payload directly over an encrypted SSH standard input pipe, avoiding command-line argument size limits and base64 memory bloat.
6. Tests execution on the remote host (`gitmap version`) to ensure 100% operational readiness.

## Arguments

- `[target]`: Node alias (e.g. `w1`, `w2`, `w3`), IP address (`192.168.1.3`), or `all` to broadcast to the entire online fleet. Defaults to `all`.

## Options

- `--file, -f <path>`: Local binary to deploy. Defaults to current running executable or newly built binary in workspace.
- `--dry-run`: Previews the target nodes and paths without transmitting data.
- `--timeout, -t <duration>`: Network transfer timeout per node (default: 60s).
- `--help, -h`: Show terminal help guide.

## Examples

```bash
# Deploy latest binary to node w1
gitmap ssh deploy-bin w1

# Deploy to all online fleet machines concurrently
gitmap ssh deploy-bin all

# Deploy a specific build artifact to node w3
gitmap ssh deploy-bin w3 --file ./cli/gitmap.exe

# Preview fleet deployment
gitmap ssh deploy-bin --dry-run
```
