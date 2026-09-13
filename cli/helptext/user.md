# gitmap user

Manage cross-platform OS-level user accounts, root/administrative access, and SSH keys across Windows, Ubuntu, Debian, and Fedora environments.

This command natively binds to `net user` on Windows and `useradd`/`userdel` on Linux distributions to seamlessly provision or decommission system users.

## Usage

```bash
gitmap user <command> [arguments] [flags]
```

## Commands

| Command | Description |
|---------|-------------|
| `add <username> [--password <pwd>]` | Create a new standard OS user |
| `create-root <username> [flags]` | Create root/sudo user with ZSH, sudoers NOPASSWD, and SSH keys |
| `rm <username> [--remove-home] [--kill]` | Remove an OS user, clean sudoers, and delete home directory |
| `kill <username> [--force]` | Terminate all active processes owned by user |
| `add-ssh-key <username> <key-or-file>` | Install SSH public key into authorized_keys with 0600 permissions |

## Examples

```bash
# Add a standard local user
gitmap user add johndoe

# Create a root/sudo user with ZSH and SSH public key
gitmap user create-root deployer --password secret123 --theme fletcherm --ssh-key "ssh-ed25519 AAA..."

# Install SSH public key for an existing user
gitmap user add-ssh-key deployer ~/.ssh/id_ed25519.pub

# Terminate all processes owned by user
gitmap user kill johndoe

# Remove user and their home directory / profile
gitmap user rm johndoe
```
