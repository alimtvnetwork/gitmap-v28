# gitmap os update / upgrade

Universal multi-distro system and package manager update across Windows, Linux, and macOS.

## Usage

```bash
gitmap os update [flags]
gitmap os upgrade [flags]
```

## Subcommands

| Command | Description |
|---------|-------------|
| update | Check and refresh package repository indexes without installing packages |
| upgrade | Perform non-interactive package upgrades across all detected package managers |

## Supported Package Managers

| OS | Detected Toolchains |
|----|---------------------|
| Windows | winget (Windows Package Manager) |
| Debian / Ubuntu | apt-get, snap, flatpak |
| Fedora / RHEL | dnf, flatpak |
| Arch Linux | pacman, flatpak |
| macOS | brew (Homebrew) |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| -n, --dry-run | false | Simulate updates without modifying packages |
| -h, --help | false | Display help and usage information |

## Examples

### Check and Refresh Package Indexes

```bash
gitmap os update
```

### Upgrade All Packages Across Host

```bash
gitmap os upgrade
```

### Simulate Upgrades (Dry Run)

```bash
gitmap os upgrade --dry-run
```
