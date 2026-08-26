# Specification: Multi-OS Installer, Cluster SSH, OS Update & Deep Help

## 1. Multi-OS Installer Subsystem
### Dedicated OS Commands
- `gitmap in install-ubuntu <slug>`: Runs Ubuntu installer sequence (`apt`, `snap`, bash).
- `gitmap in install-debian <slug>`: Runs Debian installer sequence (`apt`, `dpkg`, bash).
- `gitmap in install-arch <slug>`: Runs Arch Linux installer sequence (`pacman`, `yay`, `paru`, bash).
- `gitmap in install-centos <slug>`: Runs CentOS/RHEL/Rocky/Alma installer sequence (`dnf`, `yum`, bash).
- `gitmap in install-fedora <slug>`: Runs Fedora installer sequence (`dnf`, bash).
- `gitmap in install-mac <slug>`: Runs macOS installer sequence (`brew`, zsh, bash).
- `gitmap in install-win <slug>`: Runs Windows installer sequence (PowerShell, `winget`, `choco`, CMD).
- `gitmap in install-unix <slug>`: Runs universal POSIX Unix installer sequence.
- `gitmap in install <slug>`: Auto-detects current host OS and selects matching block.

### OS Updates & Script Modification
- `gitmap in update-ubuntu <slug>` / `update-arch` / `update-centos` / `update-debian` / `update-fedora` / `update-mac` / `update-unix`
- Bumps semantic version and updates OS-specific payload in `installer_scripts` and archives historical state in `installer_versions`.

### Universal Unix Execution Ordering & Language Runtimes
- `--order unix-first`: Executes universal Unix check/setup first, then OS-specific script.
- `--order os-first`: Executes OS-specific script first, then universal Unix post-verification.
- `--order os-only`: Executes OS-specific script only.
- `--order fallback`: Falls back to universal Unix script if OS-specific script is absent.
- Runtime support for interpreters: `sh`, `bash`, `powershell`, `cmd`, `python`, `node`, `bun`, `php`, `go`.

### Installer & Target Removal
- `gitmap in rm <slug>` / `delete <slug>`: Deletes installer and version history.
- `gitmap in rm <slug> --os <os_name>`: Strips specific OS target payload.
- `gitmap in rm-version <slug> <vX.Y.Z>`: Purges specific historical version tag.

### Git Direct Export & Auto-Commit
- `gitmap in export-all-git <repo-url | folder> [--file <name>] [--branch <bname>] [--message <msg>] [--push]`
- `gitmap in export-git <slug> <repo-url | folder>`
- Local Folder: copies JSON/ZIP, performs `git add`, `git commit -m "<msg>"`, and optional `git push`.
- Remote URL: clones repo in temporary sandbox, writes portable forward-slash JSON, commits and pushes.

---

## 2. Distributed SSH & Multi-Node Cluster
- Multi-IP Public Key Distribution: `gitmap sj distribute-keys <ip1,ip2...> [--user <u=root>] [--key <path>]`
- Multi-Machine Login & Execution: `gitmap sj broadcast <ip1,ip2...> --cmd "<command>"`
- ZSH & Profile Configuration Sync: `gitmap sj sync-profile <ip1,ip2...> [--zsh]`

---

## 3. Universal OS Update & Regional Mirror Auto-Fix
- `gitmap os update`:
  - Windows: Powershell Windows Update API (`usoclient` / `wuauclt`).
  - Ubuntu/Debian: `apt-get update && apt-get upgrade -y`.
  - CentOS/RHEL/Fedora: `dnf upgrade -y` / `yum update -y`.
  - Arch: `pacman -Syu --noconfirm`.
  - macOS: `softwareupdate -ia --verbose`.
- `gitmap os full-upgrade`: Executes OS version upgrade (`do-release-upgrade` or `dist-upgrade`).
- `gitmap os fix-mirrors` (alias `gitmap os update-fix`):
  - Automatically fixes regional Ubuntu/Debian mirror glitches (e.g., Malaysia `my.archive.ubuntu.com` or corrupted regional mirrors).
  - Backs up `/etc/apt/sources.list` to `/etc/apt/sources.list.bak-<timestamp>`.
  - Rewrites regional URLs to canonical official US/Global repositories (`archive.ubuntu.com`, `security.ubuntu.com`).

---

## 4. Comprehensive Help System & Deep `commit-in` Help
- Detailed CLI help on root `gitmap --help` and subcommands `gitmap <cmd> help`.
- `gitmap commit-in help` / `gitmap commit-write help`:
  - Full multi-section documentation explaining JSON configurations, author rotation, SEO scheduling templates, deduplication heuristics, AST function intelligence (funcintel), and profile synchronization with copy-pasteable JSON schemas.
