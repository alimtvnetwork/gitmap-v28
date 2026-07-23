# Windows Icon Embedding (go-winres)

## Overview

Embed a custom icon and version metadata into compiled Windows
binaries so they appear professional in File Explorer, taskbar,
and Task Manager. Uses [`go-winres`](https://github.com/tc-hib/go-winres),
a Go tool that generates Windows resource `.syso` files.

## How It Works

1. A `winres/` directory sits next to `main.go` containing:
   - `winres.json` — metadata manifest (version, description, icon refs)
   - `icon.png` (or `icon.ico`) — the application icon
   - Optionally `icon16.png` for small sizes

2. `go-winres make` reads the manifest and produces
   `rsrc_windows_*.syso` files (one per architecture).

3. `go build` automatically links any `.syso` file it finds in the
   package directory — no extra flags needed.

4. On non-Windows builds, `.syso` files are ignored.

## Directory Layout

```
gitmap-v28/
├── main.go
├── winres/
│   ├── winres.json      # Metadata manifest
│   ├── icon.png         # 256x256+ app icon (PNG)
│   └── icon16.png       # 16x16 small icon (optional)
├── rsrc_windows_amd64.syso  # Generated — committed to repo
└── rsrc_windows_arm64.syso  # Generated — committed to repo
```

## Manifest — `winres.json`

```json
{
  "RT_GROUP_ICON": {
    "APP": {
      "0000": [
        "icon.png",
        "icon16.png"
      ]
    }
  },
  "RT_MANIFEST": {
    "APP": {
      "0000": {
        "identity": {
          "name": "gitmap-v28",
          "version": "0.0.0.0"
        },
        "description": "gitmap-v28 - Git repository scanner, mapper, and manager",
        "minimum-os": "win7",
        "execution-level": "asInvoker",
        "dpi-awareness": "per-monitor-v2",
        "use-common-controls": true
      }
    }
  },
  "RT_VERSION": {
    "#1": {
      "0000": {
        "fixed": {
          "file_version": "0.0.0.0",
          "product_version": "0.0.0.0"
        },
        "info": {
          "0409": {
            "CompanyName": "Riseup Asia LLC",
            "FileDescription": "gitmap-v28 CLI",
            "FileVersion": "",
            "InternalName": "gitmap-v28",
            "LegalCopyright": "© 2026 Riseup Asia LLC",
            "OriginalFilename": "gitmap.exe",
            "ProductName": "gitmap-v28",
            "ProductVersion": ""
          }
        }
      }
    }
  }
}
```

### Manifest Fields

| Field | Purpose |
|-------|---------|
| `RT_GROUP_ICON` | References icon files (PNG or ICO) |
| `RT_MANIFEST` | Windows app manifest (DPI, elevation) |
| `RT_VERSION` | Version info shown in Properties → Details |
| `CompanyName` | Displayed in file properties |
| `FileDescription` | Tooltip in File Explorer |
| `ProductName` | Shown in Task Manager |

## Build Integration

### Prerequisites

Install `go-winres` as a Go tool:

```bash
go install github.com/tc-hib/go-winres@latest
```

### Generate `.syso` Files

```bash
cd gitmap-v28
go-winres make
```

This creates `rsrc_windows_amd64.syso` (and arm64 if configured).
The `.syso` files should be **committed to the repository** so that
builds work without requiring `go-winres` to be installed.

### `run.ps1` Integration

Add before the `go build` step:

```powershell
# Step 2b/4: Generate Windows resources (icon + version info)
$winresDir = Join-Path $ProjectDir "winres"
$winresJson = Join-Path $winresDir "winres.json"
if (Test-Path $winresJson) {
    $goWinres = Get-Command go-winres -ErrorAction SilentlyContinue
    if ($goWinres) {
        Write-Step "2b" "4" "Generating Windows resources"
        Push-Location $ProjectDir
        go-winres make
        if ($LASTEXITCODE -ne 0) {
            Write-Warn "go-winres failed; continuing with existing .syso"
        }
        Pop-Location
    } else {
        Write-Info "go-winres not installed; using committed .syso files"
    }
}
```

### `run.sh` Integration (Unix)

On Unix, skip resource generation (`.syso` is Windows-only):

```bash
# go-winres is Windows-only; .syso files are ignored on Unix builds
```

### CI Pipeline Integration

In `.github/workflows/release.yml`, add before cross-compilation:

```yaml
- name: Generate Windows resources
  run: |
    go install github.com/tc-hib/go-winres@latest
    cd gitmap-v28 && go-winres make
```

## Version Stamping at Build Time

For dynamic version embedding, update `winres.json` before building:

```powershell
# Update version in winres.json before go-winres make
$version = & .\bin\gitmap.exe version 2>&1
$winresContent = Get-Content $winresJson | ConvertFrom-Json
$winresContent.RT_VERSION.'#1'.'0000'.info.'0409'.ProductVersion = $version
$winresContent | ConvertTo-Json -Depth 10 | Set-Content $winresJson
```

Alternatively, use `go-winres make --product-version $version` to
override at generation time without modifying the JSON file.

## Icon Design Guidelines

| Property | Value |
|----------|-------|
| Format | PNG (preferred) or ICO |
| Sizes | 256×256 primary, 16×16 small (optional) |
| Style | Terminal/CLI themed with git branch nodes |
| Colors | Emerald green (#2EA043, #1A7F37) on dark (#161B22) |
| Background | Transparent |

The icon should be recognizable at 16×16 (taskbar) and look
detailed at 256×256 (File Explorer large icons).

## Committed Files

The `.syso` files are committed to the repo so that:
- Builds work without `go-winres` installed
- CI doesn't need the tool for every build (only when icon/manifest changes)
- Contributors can build without extra tooling

Add to `.gitignore` only if you want to force regeneration every build.

## What the User Sees

| Context | Before | After |
|---------|--------|-------|
| File Explorer | Generic `.exe` icon | Custom gitmap-v28 icon |
| Properties → Details | Empty metadata | Version, company, description |
| Taskbar | Generic icon | Branded icon |
| Task Manager | `gitmap.exe` | `gitmap-v28` with description |

## Cross-References

| Topic | Spec |
|-------|------|
| Build pipeline | [09-build-deploy.md](../01-app/09-build-deploy.md) |
| PowerShell build patterns | [02-powershell-build-deploy.md](02-powershell-build-deploy.md) |
| Build scripts (generic) | [04-build-scripts.md](../08-generic-update/04-build-scripts.md) |
| Release pipeline | [02-release-pipeline.md](../07-generic-release/02-release-pipeline.md) |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
