# Template vs Project-Specific Differences

> **Location:** `spec/11-powershell-integration/01-template-vs-project-differences.md`  
> **Updated:** 2026-02-09

---

## Overview

The generic template (`spec/powershell-integration/templates/run.ps1`, 745 lines) is a portable subset of the project's actual `run.ps1` (1006 lines). This document catalogs features **intentionally excluded** from the template because they are project-specific.

---

## Project-Specific Features (NOT in Template)

### 1. `-u` / `-upload` Flag (WordPress Plugin Upload)

| Aspect | Detail |
|--------|--------|
| Actual lines | 14, 627–692 |
| Purpose | Upload a plugin to a WordPress site via `upload-plugin-v2.ps1` |
| Dependencies | `wpPlugins` block in `powershell.json`, `wp-plugins/scripts/wp-plugin-config.json`, `upload-plugin-v2.ps1` |
| Why excluded | Specific to WordPress plugin management projects |

### 2. `-z` / `-zip` Flag (Versioned ZIP Archive)

| Aspect | Detail |
|--------|--------|
| Actual lines | 15, 525–622 |
| Purpose | Create a versioned ZIP of a plugin folder (e.g., `riseup-asia-uploader-v1.36.1.zip`) |
| Dependencies | `wpPlugins` config, PHP header version extraction (`Version: x.y.z`) |
| Why excluded | Specific to WordPress plugin packaging |

### 3. `-pp` / `-pluginpath` Parameter

| Aspect | Detail |
|--------|--------|
| Actual lines | 18 |
| Purpose | Override plugin folder path for `-u` or `-z` operations |
| Why excluded | Only relevant when `-u` or `-z` flags exist |

### 4. Build Summary in `-b` Mode

| Aspect | Detail |
|--------|--------|
| Actual lines | 950–955 |
| Purpose | Prints per-step timing breakdown when using `-buildonly` |
| Status | **Could be added** to template — it's generic functionality |

### 5. `wpPlugins` Config Block in `powershell.json`

The actual `powershell.json` includes a `wpPlugins` block for managing multiple WordPress plugins:

```json
{
  "wpPlugins": {
    "defaultUploader": "riseup-asia-uploader",
    "pluginsDir": "wp-plugins",
    "plugins": {
      "riseup-asia-uploader": {
        "name": "Riseup Asia Uploader",
        "path": "wp-plugins/riseup-asia-uploader",
        "mainFile": "riseup-asia-uploader.php",
        "autoUpload": true
      }
    }
  }
}
```

This is not in the template's `powershell.json` because it's WordPress-specific.

### 6. `pnpmVirtualStorePath` Config Key

The actual `powershell.json` has `pnpmVirtualStorePath` (line 19) but neither script reads it — both hardcode `.pnpm` in `Configure-PnpmStore`. This key is vestigial and can be removed from the actual config.

---

## Features Present in BOTH (Confirmed Synchronized)

| Feature | Status |
|---------|--------|
| Core pipeline (Git → Prerequisites → Build → Copy → Run) | ✅ Identical |
| pnpm v10+ `--dangerously-allow-all-builds` | ✅ Identical |
| PnP/ESM `NODE_OPTIONS` injection | ✅ Identical |
| Cross-drive store detection + fallback | ✅ Identical |
| Node v24+ isolated linker fallback | ✅ Identical |
| Auto-install via `winget` (Go, Node, pnpm) | ✅ Identical |
| Force-clean with wildcard support | ✅ Identical |
| Backend runtime data cleaning | ✅ Identical |
| Firewall rule management | ✅ Identical |
| Config file auto-creation from example | ✅ Identical |
| Environment variable injection from config | ✅ Identical |
| Rebuild mode (`-r` = `-f` + `-i`) | ✅ Identical |
| Install deferral in rebuild mode | ✅ Identical |

---

## Recommendation

The template is correctly generic. To adopt WP plugin features in a new project, copy the `-u`, `-z`, and `-pp` sections from the actual script and add the `wpPlugins` block to `powershell.json`.
