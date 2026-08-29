# gitmap user

Manage cross-platform OS-level user accounts across Windows, Ubuntu, Debian, and Fedora environments.

This command natively binds to `net user` on Windows and `useradd`/`userdel` on Linux distributions to seamlessly provision or decommission system users.

## Usage

```bash
gitmap user <command> [arguments]
```

## Commands

| Command | Description |
|---------|-------------|
| `add <username> [--password <pwd>]` | Create a new OS user |
| `rm <username>` | Remove an OS user and their profile data |

## Examples

```bash

# Add a local user on Windows or Linux

gitmap user add johndoe

# Add a user with an initial password

gitmap user add janedoe --password secret123

# Completely remove the user and their home directory / profile

gitmap user rm johndoe
```
