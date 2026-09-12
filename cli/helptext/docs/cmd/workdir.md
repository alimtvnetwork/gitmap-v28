# gitmap workdir (wd)

Manage multiple codebase workspaces in a single Gitmap installation. This allows Gitmap commands like `gitmap clone` and `gitmap ct` to target specific working directories dynamically.

## Usage

```bash
gitmap workdir ls                     # List all registered work directories
gitmap workdir add [path] [--label <label>] # Register a work directory (defaults to current directory)
gitmap workdir set <path|id>          # Set the default work directory (alias for set-default)
gitmap workdir set-default <path|id>  # Set the default work directory
gitmap workdir default                # Set the current directory as the default work directory
gitmap workdir rm <path|id>           # Remove a registered work directory
```

## Examples

**Register a new work directory:**
```bash
# Windows
gitmap workdir add D:\work --label "work"

# Linux / macOS
gitmap workdir add /work --label "work"
```

**Set current directory as default:**
```bash
cd D:\work        # or: cd /work
gitmap workdir default
```

**Set a workspace as default by ID or alias:**
```bash
gitmap workdir set 2
gitmap workdir set company
```

**List all directories:**
```bash
gitmap workdir ls
```
