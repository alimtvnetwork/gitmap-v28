# gitmap vscode group

Manage named groups of VS Code projects and workspaces.

## Subcommands

- `gitmap vscode group ls`: List all configured VS Code groups and their target paths.
- `gitmap vscode group add <name> <targets...>`: Add workspace path(s) to a VS Code group.
- `gitmap vscode group rm <name> [target]`: Remove a target from a group or delete the entire group.

## Examples

```bash
# List all VS Code groups
gitmap vscode group ls

# Create/add projects to a group
gitmap vscode group add frontend D:\projects\web D:\projects\ui
gitmap vscode group add backend D:\projects\api

# Remove a path from a group
gitmap vscode group rm frontend D:\projects\ui

# Delete an entire group
gitmap vscode group rm frontend
```
