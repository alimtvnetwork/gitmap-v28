# antigravity (agy, ag) — Antigravity Workspace Management

Manage Antigravity projects, project groups, undo/redo state, plugins, and settings.

## Subcommands

- `gitmap agy group [add|rm|ls|show|export|import|prompt]`: Manage and prompt named groups of Antigravity projects.
- `gitmap agy undo`: Revert the last clear or project removal action from an automatic snapshot.
- `gitmap agy redo`: Reapply the undone Antigravity project state.
- `gitmap agy plugin [ls|install <slug>]`: List installed Antigravity plugins and install new plugins.
- `gitmap agy settings [export|import <file.json>]`: Export and import Antigravity configuration in JSON format.
- `gitmap agy ls`: List all discovered Antigravity projects with health status.
- `gitmap agy clear`: Remove stale or missing projects (snapshots state before deletion).
- `gitmap agy optimize-projects`: Deduplicate and optimize registered project entries.
- `gitmap agy export-projects <file.zip>`: Create a zip archive backup of all projects.
- `gitmap agy import-projects <file.zip>`: Restore projects from a zip archive backup.
- `gitmap agy install [ide|manager|cli|all]`: Install Google Antigravity Desktop IDE (default), Manager GUI, or CLI.
- `gitmap agm install`: Install Antigravity Manager GUI / Tools (or `gitmap install agm`).

## Installation Flags

The unified Antigravity installer engine supports the following multi-platform flags:

- `--force`, `-f`: Force reinstallation even if already present on the system.
- `--prefix`, `--dir <path>`: Custom installation directory prefix (overrides system defaults).
- `--dry-run`, `-n`: Simulate installation without downloading or modifying the host filesystem.
- `--version <ver>`: Specific release version to install (defaults to 2.13.0).
- `--yes`, `-y`: Automatic yes to interactive prompts.
- `--verbose`, `-v`: Enable detailed diagnostic output during installation.

## Examples

```bash
# Install Google Antigravity Desktop IDE and CLI (unified engine)
gitmap install antigravity
gitmap install agy

# Install via agy subcommand with multi-platform flags
gitmap agy install ide --force
gitmap agy install cli --dry-run
gitmap agy install all --prefix /opt/custom/antigravity

# Manage project groups
gitmap agy group add frontend repo-a repo-b
gitmap agy group ls
gitmap agy group prompt frontend "Review styling and line gaps"

# Undo a recent clear action
gitmap agy clear --all
gitmap agy undo

# Inspect and install plugins
gitmap agy plugin ls
gitmap agy plugin install firebase-agent-plugin

# Export and import settings
gitmap agy settings export my-settings.json
gitmap agy settings import my-settings.json
```
