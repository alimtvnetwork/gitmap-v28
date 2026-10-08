# gitmap user

Manage Git author profiles, per-project repository bindings, GitHub CLI auth inspection, and cross-platform OS-level user accounts.

## Usage

```bash
gitmap user <command> [arguments] [flags]
```

## Git Profile & Identity Commands

| Command | Description |
|---------|-------------|
| `info`, `status` `[--json]` | Display consolidated Git identity and GitHub CLI authorization status card |
| `list`, `ls` | List configured Git user profiles with active and bound indicators |
| `switch`, `use <alias> [--global] [--project]` | Switch active Git profile or bind to current project |
| `add <alias> --name <n> --email <e>` | Register a new Git author profile |
| `project [bind <alias>\|unbind\|status]` | Manage per-repository profile bindings |
| `config [global\|local] [--name <n>] [--email <e>]` | Inspect or configure Git `user.name` and `user.email` |
| `sync` | Apply bound project user profile to current repository git config |

## Operating System User Commands

| Command | Description |
|---------|-------------|
| `add <username> [--password <pwd>]` | Create a new standard OS user |
| `create-root <username> [flags]` | Create root/sudo user with ZSH, sudoers NOPASSWD, and SSH keys |
| `rm <username> [--remove-home] [--kill]` | Remove an OS user, clean sudoers, and delete home directory |
| `kill <username> [--force]` | Terminate all active processes owned by user |
| `add-ssh-key <username> <key-or-file>` | Install SSH public key into authorized_keys with 0600 permissions |

## Examples

```bash
# Inspect GitHub CLI auth and active Git identity
gitmap user info
gitmap user info --json

# List configured Git profiles
gitmap user list

# Switch active profile for current project
gitmap user switch work --project

# Register a new Git author profile
gitmap user add work --name "Work Dev" --email "dev@work.com"

# Bind current repository to profile
gitmap user project bind work

# Add a standard local OS user
gitmap user add johndoe

# Create a root/sudo user with ZSH and SSH public key
gitmap user create-root deployer --password secret123 --theme fletcherm --ssh-key "ssh-ed25519 AAA..."
```
