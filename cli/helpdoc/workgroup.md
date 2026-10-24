# gitmap workgroup

Manage folder-specific work directories, SSH key bindings, and directory grouping.

## Synopsis

```
gitmap workgroup [add|bind|unbind|ls|show|test-ssh|rm] [arguments...] [flags]
gitmap wg <command>          # Short alias
gitmap group <command>       # Grouping alias
```

## Behavior

1. **Named Workgroups (`add`):** Registers logical workgroup directory roots (e.g. `client-a`, `foss`, `corp`) with dedicated SSH keys and Git author profiles.
2. **Path SSH Key Binding (`bind`):** Directly associates a directory with an SSH private key (`~/.ssh/id_client`) and Git author identity.
3. **Ancestor Lineage Resolution:** Child repositories inside nested directories automatically inherit workgroup settings and SSH credentials from parent folders without per-repo reconfiguration.
4. **Isolated SSH Handshake Testing (`test-ssh`):** Tests authentication against remote Git hosts (`github.com`, `gitlab.com`, `bitbucket.org`) using `ssh -T -i <key>`.
5. **Zero-Prompt Git Subprocess Injection:** Automatically configures isolated `GIT_SSH_COMMAND` parameters and author profiles during git operations.

## Commands

- `add <name> <dir>`: Register a named workgroup for a filesystem directory root.
- `bind <dir> <key>`: Bind a specific directory to an SSH private key.
- `unbind <dir>`: Remove folder-specific SSH key binding (reverts to parent or global).
- `ls`, `list`: Display tabular summary of all registered workgroups and bindings.
- `show [name|dir]`: Inspect detailed configuration, SSH key path, and author profile.
- `test-ssh [name|dir]`: Test remote SSH handshake authentication.
- `rm <name|id>`: Safely remove workgroup registration.

## Flags

- `--ssh-key <path>`: Path to private SSH key file.
- `--desc <text>`: Human-readable description of the workgroup.
- `--user-name <name>`: Custom `user.name` applied to commits within this group.
- `--user-email <email>`: Custom `user.email` applied to commits within this group.
- `--default`: Designates this workgroup as the active default workspace.
- `--host <hostname>`: Remote SSH host for handshake testing (default: `github.com`).
- `-j`, `--json`: Output structured JSON representation.
- `-h`, `--help`: Display terminal help menu.

## Examples

```bash
# Register a client workspace with a dedicated SSH key and commit author
gitmap wg add client-a D:\work\client-a --ssh-key ~/.ssh/id_client_a --user-name "John Doe" --user-email "john@clienta.com"

# Bind an SSH key to a specific repository folder
gitmap wg bind D:\work\client-a\backend ~/.ssh/id_backend

# List all registered workgroups
gitmap wg ls

# Inspect effective configuration and inheritance for a directory
gitmap wg show D:\work\client-a\backend\auth-service

# Verify SSH authentication against GitHub
gitmap wg test-ssh client-a --host github.com

# Safely delete workgroup registration
gitmap wg rm client-a
```
