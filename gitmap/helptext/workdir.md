# gitmap workdir

Manage and inspect registered work directories, active default workdir,
and quick workspace targets for cd navigation.

## Synopsis

```
gitmap workdir [ls | add <path> | rm <path|id> | set <path|id> | default [path] | path]
gitmap wd                                                # short alias
```

## Behavior

1. **List Work Directories (`ls`):** Displays all registered workspaces with IDs,
   labels, absolute paths, and default status indicators.
2. **Add Work Directory (`add`):** Registers a workspace directory with an optional
   custom label via `--label <name>`.
3. **Set Default (`set` / `default <path>`):** Sets the active default work directory.
   If the path is not yet registered, it is automatically added and set as default.
4. **Inspect Default (`default` / `path`):** Prints the active default workdir info or
   its raw absolute path for shell integration.
5. **CD Integration (`gitmap cd work`):** `gitmap cd work`, `gitmap cd default`, and bare
   `gitmap cd` automatically resolve and navigate to the default work directory.

## Examples

```
$ gitmap workdir default
✓ Default work directory: D:\wp-work\riseup-asia (ID: 1, Label: riseup-asia)

$ gitmap workdir set D:\wp-work\riseup-asia
✓ Default work directory set to: D:\wp-work\riseup-asia

$ gitmap cd work
D:\wp-work\riseup-asia
```

## Commands

- `ls`, `list`: List all registered work directories.
- `add <path> [--label <l>]`: Register a work directory.
- `rm <path|id>`: Remove a registered work directory.
- `set <path|id>`: Set the active default work directory.
- `default [path]`: Show or set the active default work directory.
- `path`, `get`: Print only the absolute path of the default work directory.
