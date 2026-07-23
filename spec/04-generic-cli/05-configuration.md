# Configuration Pattern

> **Related specs:**
> - [04-flag-parsing.md](04-flag-parsing.md) — CLI flags that override config values
> - [10-database.md](10-database.md) — SQLite persistence complementing file-based config
> - [15-constants-reference.md](15-constants-reference.md) — constant naming for config keys

## Three-Layer Config

```
Defaults (hardcoded) → Config file (JSON) → CLI flags (highest priority)
```

| Layer | Source | Priority |
|-------|--------|----------|
| 1. Defaults | Constants in code | Lowest |
| 2. Config file | `./data/config.json` or `--config <path>` | Medium |
| 3. CLI flags | `--mode ssh`, `--output json` | Highest (always wins) |

## Config File

### Location

- Default: `./data/config.json` (relative to binary)
- Override: `--config <path>` flag
- Missing file: use defaults silently (no error)

### Schema

Define a flat JSON structure:

```json
{
  "defaultMode": "https",
  "defaultOutput": "terminal",
  "outputDir": "./toolname-output",
  "excludeDirs": [".cache", "node_modules", "vendor"],
  "notes": ""
}
```

### Rules

| Rule | Detail |
|------|--------|
| Field names | camelCase |
| Array fields | Default to `[]`, never `null` |
| String fields | Default to `""`, never `null` |
| No nested objects | Unless absolutely necessary |
| Struct mirrors JSON | No transformation between file and struct |

## Merge Logic

```go
func LoadAndMerge(configPath, flagMode, flagOutput string) Config {
    cfg := loadDefaults()

    if fileExists(configPath) {
        fileCfg := loadFromFile(configPath)
        cfg = merge(cfg, fileCfg)
    }

    if flagMode != "" {
        cfg.Mode = flagMode
    }
    if flagOutput != "" {
        cfg.Output = flagOutput
    }

    return cfg
}
```

### Merge Rules

1. Load hardcoded defaults into a config struct.
2. If config file exists, overlay its values onto the struct.
3. Apply CLI flags on top — flags always win.
4. If config file is missing, proceed silently with defaults.

## Key Principles

| Principle | Detail |
|-----------|--------|
| Never crash on missing config | Use defaults, warn if needed |
| Flags always win | Explicit user intent overrides everything |
| Config paths relative to binary | Unless absolute |
| Default paths in `constants` | `DefaultConfigPath`, `DefaultOutputDir` |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
