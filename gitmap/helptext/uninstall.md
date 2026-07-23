# gitmap uninstall

Remove a third-party tool **or** the gitmap binary itself, depending on
whether you pass a tool name.

## Two modes

| Form                              | What it does                                                |
|-----------------------------------|-------------------------------------------------------------|
| `gitmap uninstall <tool> [flags]` | Remove a third-party tool (vscode, npp, …) via its package manager. |
| `gitmap uninstall [flags]`        | Shortcut → `gitmap self-uninstall [flags]` (removes the gitmap binary, data dir, PATH snippet). |

The "no tool name" shortcut was added in v3.75.0 so users do not have to
remember the `self-` prefix. Flags pass through verbatim, so:

```
gitmap uninstall --confirm --keep-data
```

is exactly equivalent to:

```
gitmap self-uninstall --confirm --keep-data
```

See `gitmap help self-uninstall` for the full self-uninstall flag table
(including `--shell-mode`, `--keep-snippet`, etc.).

## Tool-uninstaller flags

| Flag        | Description                                       |
|-------------|---------------------------------------------------|
| `--dry-run` | Print the package-manager command; do not run it. |
| `--force`   | Skip confirmation; ignore "not installed" errors. |
| `--purge`   | Also remove configuration files (apt purge, choco -x). |

## Examples

### Remove gitmap itself

```
gitmap uninstall                       # interactive
gitmap uninstall --confirm             # skip prompt
gitmap uninstall --confirm --keep-data # keep ~/.config/gitmap or %APPDATA%\gitmap
```

### Remove a third-party tool

```
gitmap uninstall vscode
gitmap uninstall npp --dry-run
gitmap uninstall git --force --purge
```

## See Also

- `gitmap self-uninstall` — the underlying binary remover
- `gitmap install <tool>` — install third-party tools
- `gitmap reinstall` — uninstall + reinstall the binary in one step

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter uninstall
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
