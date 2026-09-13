# gitmap zsh

Install, configure, switch, profile, and maintain ZSH and Oh-My-Zsh environments across Ubuntu, Debian, CentOS, and Fedora.

## Usage

```bash
gitmap zsh <command> [arguments] [flags]
```

## Commands

| Command | Description |
|---------|-------------|
| `install [flags]` | Install zsh, unattended Oh-My-Zsh, and plugins |
| `theme [name] [flags]` | Change ZSH theme in ~/.zshrc and append custom configs |
| `switch [flags]` | Change default login shell to ZSH |
| `profile [flags]` | Provision standard workspace directories (scripts, gitlab, github, .ssh) |
| `clean [flags]` | Backup .zshrc, purge oh-my-zsh, and optionally reinstall cleanly |
| `status` | Inspect local ZSH installation, active theme, and plugins |
| `help` | Show command usage and options |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--theme <name>` | `robbyrussell` | Initial or target Oh-My-Zsh theme |
| `--append-rc` | false | Append and deduplicate developer configs in ~/.zshrc |
| `--user <username>` | current | Target user account |
| `--reinstall` | true | Clean and reinstall oh-my-zsh |
| `--backup-rc` | true | Backup existing ~/.zshrc prior to wiping |

## Examples

```bash
# Install ZSH and Oh-My-Zsh with autosuggestions
gitmap zsh install

# Change active theme to agnoster
gitmap zsh theme agnoster

# Provision standard workspace folders with secure permissions
gitmap zsh profile

# Clean and reinstall Oh-My-Zsh
gitmap zsh clean

# Check current ZSH status
gitmap zsh status
```
