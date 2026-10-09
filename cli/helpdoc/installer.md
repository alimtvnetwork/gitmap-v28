# gitmap installer

Manage, author, export, import, update, and version multi-platform installation scripts and packages.

## Usage

```bash
gitmap installer [subcommand] [flags]
```

## Description

The `gitmap installer` command family provides a script installation creation, export, import, universal Unix execution ordering, Git-direct auto-committing, and versioning system. It manages both built-in and user-authored installers stored in SQLite.

## Subcommands

| Subcommand | Description |
|------------|-------------|
| `ls` | List all registered installers, versions, and target platforms |
| `add <name> [version]` | Create or update a custom installer interactively |
| `create <name>` | Generate a new installer script scaffolding with template headers |
| `export <name> [path]` | Export an installer script to a file or directory |
| `export-all [dir]` | Export all stored installers to a target directory |
| `export-git <name>` | Export installer directly committed to the local Git repository |
| `export-all-git` | Batch export and auto-commit all installers to Git |
| `import <path>` | Import an installer YAML/JSON or shell script into the database |
| `update <name>` | Update an existing installer definition, URL, or script content |
| `rm <name>` | Delete an installer from the database |
| `install <name>` | Execute smart installer with multi-manager fallback |
| `win <name>` | Execute Windows-specific installer package or PowerShell hook |
| `update-win <name>` | Update Windows-specific installer definitions |
| `os-update <os>` | Run OS package repository updates (`ubuntu`, `debian`, `arch`, `fedora`, `mac`, `unix`) |
| `os-install <os>` | Run OS-specific package installation for target distribution |
| `undo` / `redo` / `revert` | Navigate installer version history and revert changes |
| `reset` | Reset installer database to factory defaults |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--manager <name>` | (auto) | Target package manager (`choco`, `winget`, `apt`, `brew`, `snap`) |
| `--platform <os>` | (current) | Target operating system (`windows`, `linux`, `darwin`) |
| `--dry-run (-n)` | false | Show commands without executing or writing to disk |
| `--verbose (-v)` | false | Display full debug output and execution logs |
| `--force (-f)` | false | Overwrite existing scripts or bypass safety checks |
| `--json` | false | Format command output as structured JSON |

## Examples

### List Registered Installers

```bash
gitmap installer ls
```

### Create a New Custom Installer

```bash
gitmap installer create my-cli-tool
```

### Interactively Add or Update an Installer

```bash
gitmap installer add ripgrep 14.1.0
```

### Export an Installer to File

```bash
gitmap installer export ripgrep ./scripts/install-ripgrep.sh
```

### Export Directly with Git Auto-Commit

```bash
gitmap installer export-git ripgrep
```

Output:
```text
✓ Exported installer: ripgrep
✓ Committed to repository: "chore(installer): export ripgrep v14.1.0"
```

### Import an Installer Definition

```bash
gitmap installer import ./custom-installer.yaml
```

### Execute Distribution OS Update

```bash
gitmap installer os-update ubuntu
```

### Delete an Installer

```bash
gitmap installer rm obsolete-tool --force
```

## See Also

- [install](install.md) — Install developer tools, runtimes, and workstation profiles
- [os](os.md) — Manage operating system configurations and diagnostics
- [env](env.md) — Manage environment variables and PATH
