# gitmap installed-dir

Show the full path and directory of the active gitmap binary.

## Alias

id

## Usage

    gitmap installed-dir

## Description

Resolves the currently active gitmap binary location, following
symlinks to the real path. Useful for verifying which binary is
being invoked and for update scripts that need to know the target
install path.

## Examples

### Example 1: Show installed path

    gitmap installed-dir

**Output:**

      📂 Installed directory

      Binary:    /home/alim/.local/bin/gitmap
      Directory: /home/alim/.local/bin

### Example 2: Use alias

    gitmap id

**Output:**

      📂 Installed directory

      Binary:    /usr/local/bin/gitmap
      Directory: /usr/local/bin

### Example 3: Compare with which

    which gitmap
    gitmap id

**Output:**

    /home/alim/.local/bin/gitmap

      📂 Installed directory

      Binary:    /home/alim/.local/bin/gitmap
      Directory: /home/alim/.local/bin

## See Also

- [update](update.md) — Self-update from source repo
- [version](version.md) — Show version number
- [doctor](doctor.md) — Diagnose PATH and version issues
- [install](install.md) — Install developer tools (use `gitmap install scripts` to clone utility scripts)

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter installed-dir
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
