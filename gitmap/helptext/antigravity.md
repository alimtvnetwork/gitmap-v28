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

## Examples

```bash
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
