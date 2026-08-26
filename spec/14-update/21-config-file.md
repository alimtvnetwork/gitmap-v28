# Configuration File

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Purpose

Define the location, schema, and lifecycle of the CLI configuration file used for deploy path defaults, update preferences, and persistent settings.

---

## File Location

| Platform | Path |
|----------|------|
| Windows | `%APPDATA%\<binary>\config.json` |
| Linux | `~/.config/<binary>/config.json` |
| macOS | `~/Library/Application Support/<binary>/config.json` |

### XDG Compliance (Linux)

If `$XDG_CONFIG_HOME` is set, use it instead of `~/.config`:

```bash
config_dir="${XDG_CONFIG_HOME:-$HOME/.config}/<binary>"
```

### Go Implementation

```go
import (
    "os"
    "path/filepath"
    "runtime"
)

func ConfigDir() string {
    switch runtime.GOOS {
    case "windows":
        return filepath.Join(os.Getenv("APPDATA"), "<binary>")
    case "darwin":
        home, _ := os.UserHomeDir()
        return filepath.Join(home, "Library", "Application Support", "<binary>")
    default: // linux and others
        xdg := os.Getenv("XDG_CONFIG_HOME")
        if xdg == "" {
            home, _ := os.UserHomeDir()
            xdg = filepath.Join(home, ".config")
        }
        return filepath.Join(xdg, "<binary>")
    }
}

func ConfigPath() string {
    return filepath.Join(ConfigDir(), "config.json")
}
```

---

## JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "properties": {
    "deployPath": {
      "type": "string",
      "description": "Default directory for binary installation (Tier 3 of deploy path resolution)"
    },
    "repoPath": {
      "type": "string",
      "description": "Path to local source repository for source-based updates"
    },
    "updateChannel": {
      "type": "string",
      "enum": ["stable", "beta", "alpha"],
      "default": "stable",
      "description": "Which release channel to follow for updates"
    },
    "checkUpdatesOnStart": {
      "type": "boolean",
      "default": false,
      "description": "Whether to check for updates on every CLI invocation"
    },
    "lastUpdateCheck": {
      "type": "string",
      "format": "date-time",
      "description": "ISO 8601 timestamp of the last update check (set automatically)"
    }
  },
  "additionalProperties": false
}
```

### Example

```json
{
  "deployPath": "E:\\bin-run",
  "repoPath": "<repo-root>",
  "updateChannel": "stable",
  "checkUpdatesOnStart": false,
  "lastUpdateCheck": "2026-04-13T10:30:00Z"
}
```

---

## First-Time Creation

The config file is created when:

1. **First install** — the install script creates it with default `deployPath`
2. **First run** — if the binary detects no config file, it creates one with defaults
3. **Manual creation** — user creates it manually

### Default Values by Platform

| Field | Windows Default | Linux/macOS Default |
|-------|----------------|---------------------|
| `deployPath` | `%LOCALAPPDATA%\<binary>` | `$HOME/.local/bin` |
| `repoPath` | (empty — detected from embedded constant) | (empty) |
| `updateChannel` | `stable` | `stable` |
| `checkUpdatesOnStart` | `false` | `false` |

### Creation Logic

```go
func EnsureConfig() (*Config, error) {
    path := ConfigPath()

    data, err := os.ReadFile(path)
    if err == nil {
        var cfg Config
        if err := json.Unmarshal(data, &cfg); err != nil {
            return nil, fmt.Errorf("parse config %s: %w", path, err)
        }
        return &cfg, nil
    }

    // Create default config
    cfg := &Config{
        DeployPath:    defaultDeployPath(),
        UpdateChannel: "stable",
    }

    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        return nil, fmt.Errorf("create config dir: %w", err)
    }

    data, _ = json.MarshalIndent(cfg, "", "  ")
    if err := os.WriteFile(path, data, 0644); err != nil {
        return nil, fmt.Errorf("write config: %w", err)
    }

    return cfg, nil
}
```

---

## Integration with Deploy Path Resolution

The config file provides **Tier 3** of the [Deploy Path Resolution](./02-deploy-path-resolution.md) strategy:

```
Priority 1: CLI flag (--deploy-path) → highest
Priority 2: PATH lookup (existing binary location)
Priority 3: Config file default (config.json → deployPath) → lowest
```

---

## PowerShell Integration

For PowerShell build scripts, the config file is `powershell.json` in the project root (not the user config). See [PowerShell Integration](../11-powershell-integration/examples/server-client-project.json) for the project-level config format.

The **user-level** `config.json` described here is different from the **project-level** config — it stores per-user preferences, not project build settings.

---

## Constraints

- Config file uses JSON format — no YAML, TOML, or INI
- All fields are optional — missing fields use defaults
- The binary must never crash if the config file is missing or malformed
- Config file is created lazily (on first need), not eagerly at install time
- `lastUpdateCheck` is written automatically — users should not edit it
- Unknown fields are rejected (`additionalProperties: false`) to catch typos

---

## Cross-References

- [Deploy Path Resolution](./02-deploy-path-resolution.md) — Tier 3 uses config file defaults
- [Self-Update Overview](./01-self-update-overview.md) — Source repo path from config
- [Updater Binary](./19-updater-binary.md) — Reads repo and deploy path from config
- [Install Scripts](./18-install-scripts.md) — May create initial config on install

---

*Configuration file — v3.2.0 — 2026-04-13*
