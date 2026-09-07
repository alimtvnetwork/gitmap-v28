# gitmap workdir

Manage and inspect registered work directories, active default workdir,
and quick workspace targets for cd navigation.

## Synopsis

```
gitmap workdir [ls | add [path] | rm <path|id> | set <path|id> | default [path] | path]
gitmap wd                                                # short alias
```

## Behavior

1. **List Work Directories (`ls`):** Displays all registered workspaces with IDs,
   labels, absolute paths, and default status indicators.
2. **Add Work Directory (`add`):** Registers a workspace directory (defaults to current
   directory if omitted) with an optional custom label via `--label <name>`. If the directory
   is already registered, it confirms it is already added.
3. **Set Default (`set` / `default <path>`):** Sets the active default work directory.
   If the path is not yet registered, it is automatically added and set as default.
4. **Inspect Default (`default` / `path`):** Prints the active default workdir info or
   its raw absolute path for shell integration.
5. **CD Integration (`gitmap cd work`):** `gitmap cd work`, `gitmap cd default`, and bare
   `gitmap cd` automatically resolve and navigate to the default work directory.

## Examples

```
# Windows:
$ gitmap workdir default
✓ Default work directory: D:\work (ID: 1, Label: work)

$ gitmap workdir set D:\work
✓ Default work directory set to: D:\work

$ gitmap cd work
D:\work

# Linux / macOS:
$ gitmap workdir set /work
✓ Default work directory set to: /work

$ gitmap cd work
/work
```

## Commands

- `ls`, `list`: List all registered work directories.
- `add [path] [--label <l>]`: Register a work directory (defaults to current directory).
- `rm <path|id>`: Remove a registered work directory.
- `set <path|id>`: Set the active default work directory.
- `default [path]`: Show or set the active default work directory.
- `path`, `get`: Print only the absolute path of the default work directory.
