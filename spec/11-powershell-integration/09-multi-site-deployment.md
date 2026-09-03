# Multi-Site Deployment

## Overview

Upload plugins to multiple WordPress sites from a single command using Base64-encoded credentials stored in `powershell.json`.

## Configuration

### `powershell.json` — `wpPlugins.sites` array

```json
{
  "wpPlugins": {
    "sites": [
      {
        "name": "Test V1",
        "url": "https://testv1.developers-organism.com",
        "enabled": true,
        "credentials": [
          {
            "appName": "test-plg-v1",
            "usernameBase64": "<base64-encoded-username>",
            "passwordBase64": "<base64-encoded-app-password>",
            "isDefault": true
          }
        ]
      }
    ]
  }
}
```

### Credential encoding

All credentials are Base64-encoded for basic obfuscation (not encryption):

```powershell

# Encode

[Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes("my-username"))

# Decode (done automatically by run.ps1)

[System.Text.Encoding]::UTF8.GetString([Convert]::FromBase64String("bXktdXNlcm5hbWU="))
```

### Site properties

| Property | Type | Required | Description |
|---|---|---|---|
| `name` | string | ✓ | Human-readable site name (used with `-site` filter) |
| `url` | string | ✓ | WordPress site URL (no trailing slash) |
| `enabled` | boolean | | Whether included in `-uas` (default: true) |
| `credentials` | array | ✓ | Array of credential objects |

### Credential properties

| Property | Type | Required | Description |
|---|---|---|---|
| `appName` | string | ✓ | Application password name for identification |
| `usernameBase64` | string | ✓ | Base64-encoded WordPress username |
| `passwordBase64` | string | ✓ | Base64-encoded WordPress application password |
| `isDefault` | boolean | | Whether this is the default credential (first match used) |

## Usage

```powershell

# Upload all plugins to ALL enabled sites

.\run.ps1 -uas

# Upload all plugins to a specific site by name

.\run.ps1 -uas -site "Test V1"

# Upload all plugins to site #1 (1-based index from -ls)

.\run.ps1 -uas -index 1

# Upload all plugins to site #2

.\run.ps1 -uas -index 2
```

## Flow

1. Validates `wpPlugins.sites` configuration exists
2. Lists all configured sites (shows enabled/disabled status)
3. Filters by `-site` name, `-index` number, or uses all enabled sites
4. Discovers uploadable plugins (respects skipPlugins only, no hardcoded exclusions)
5. ZIPs all plugins once (reuses existing ZIP logic)
6. For each target site:
   - Resolves default credential (first with `isDefault: true`, or first entry)
   - Decodes Base64 username and password
   - Uploads each plugin via QUpload API (`upload-plugin-U-Q.ps1`)
7. Prints multi-site summary with per-site per-plugin results

## Schema

Defined in `02-spec/11-powershell-integration/schemas/powershell.schema.json` under `wpPlugins.sites`.
