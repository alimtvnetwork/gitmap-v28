# PowerShell Script Changelog

All notable changes to the PowerShell runner script (`run.ps1`) and upload scripts will be documented in this file.

## New Entry Template (copy/paste)

Use this template whenever you change `run.ps1`, upload scripts, or make a functional/config-schema change to `powershell.json`.

```md

## [X.Y.Z] - YYYY-MM-DD

### Added

- ...

### Changed

- ...

### Fixed

- ...

### Notes

- (optional) Migration steps, breaking changes, required config updates
```

---

## [2.1.0] - 2026-03-12

### Added

- **QUpload admin menu system**: Full WordPress admin UI with Dashboard and Error Logs pages, replacing the old Tools submenu
- **QUpload error log viewer**: Tabbed file viewer (Log, Error, Stack Trace) with Copy, Download, Clear, and Live auto-refresh
- **QUpload Admin class with traits**: `AdminMenuTrait`, `AdminErrorAjaxTrait` for menu registration, asset enqueuing, and AJAX file operations
- **QUpload enums**: `AdminPageType`, `AdminTabType`, `NonceType`, `AjaxActionType` for admin infrastructure
- **FileLogger getters**: `getLogFile()`, `getErrorFile()`, `getStacktraceFile()` for AJAX log access
- **QUpload admin assets**: `admin.css`, `admin-errors.css`, `admin-errors.js`

### Changed

- **QUpload bootstrap**: Admin initialization moved from `Plugin` constructor trait to standalone `Admin` singleton in `qupload.php` (matches Riseup Asia pattern)
- **Version bump**: All components synchronized to 2.1.0 (app, PowerShell, both plugins)
- **HookType enum**: Added `AdminInit`, `AdminEnqueue`, and `ajax()` helper method
- **CapabilityType enum**: Added `ManageOptions` case

### Fixed

- **Riseup Asia `admin-settings.php`**: Removed duplicate `<?php` tag causing `unexpected token "<"` parse error on deployment
- **Riseup Asia `Admin.php`**: Replaced `PaginationConfigType::logRetrievalMaxLines()` method call in static property default with literal `500` (PHP does not allow method calls in constant expressions)

---

## [upload-plugin-v2 2.1.0] - 2026-02-10

### Added

- **Self-update OPcache flush**: After uploading `riseup-asia-uploader` to itself, script calls `opcache-reset.php` to force-flush PHP OPcache, then verifies version
- **`opcache-reset.php`**: New standalone PHP file deployed with the plugin for OPcache reset (Basic Auth secured)
- **Imunify360 detection**: JSON responses with "Access denied" / "bot-protection" now throw actionable errors instead of false success
- **ZIP staging progress**: Real-time file count display during staging (e.g., `Staging: 47/47 files (100%)`)
- **Full path display**: ZIP destination path, cache directory, compression ratio shown in output
- **`Accept: application/json` header**: All HTTP requests now include this to prevent HTML challenge pages

### Changed

- **Pipeline expanded to 8 steps**: Step 8 is now self-update-aware (OPcache flush + verify) or standard version check
- **Self-update version priority**: For self-updates, client-sent version takes priority over server response (server returns stale version from cached old code)
- **V1 fallback**: Uses direct named parameters instead of array splatting to prevent JSON mangling

### Fixed

- Version mismatch on self-update caused by PHP OPcache serving old bytecode
- False "PUBLISH COMPLETE" when server returns Imunify360 block message
- V1 fallback failing due to `@fallbackArgs` array splatting mangling JSON strings

---

## [run.ps1 1.2.0] - 2026-02-08

### Added

- **Runtime data cleanup**: Force mode (`-f`, `-r`) now cleans backend sessions, request-sessions, error logs, and standalone log files from `backend/data/`
- **cleanPaths expanded**: `powershell.json` now includes `backend/data/sessions`, `backend/data/request-sessions`, and `backend/data/errors`

### Changed

- `-r` flag description updated to reflect session/log cleanup behavior

### Notes

- Directories cleaned: `data/sessions/`, `data/request-sessions/`, `data/errors/`, `data/log.txt`, `data/error.log.txt`
- Cleanup only runs when `dataDir` is configured in `powershell.json`

---

## [run.ps1 1.1.0] - 2026-02-04

### Added

- **Version tracking**: Script now has version number in header and `powershell.json`
- **PnP artifact cleanup**: Force mode now removes `.pnp.cjs`, `.pnp.loader.mjs`, `.pnp.data.json`
- **Improved install detection**: Respects `EffectiveNodeLinker` (PnP vs isolated) when checking if install is needed

### Changed

- **Rebuild sequence**: `-r` flag now correctly defers frontend install until after force-clean
- **Install always runs**: `-i` and `-r` flags always trigger `pnpm install`, even if `node_modules` exists

### Fixed

- "vite is not recognized" error when using `-r` flag

---

## [run.ps1 1.0.0] - 2026-02-02

### Added

- Initial PowerShell runner with pnpm PnP support
- Git pull, prerequisites check, pnpm install, build, and run steps
- Flags: `-b`, `-s`, `-p`, `-f`, `-i`, `-r`, `-fw`, `-h`, `-v`
- Auto-install of Go, Node.js, and pnpm via winget
- Windows Firewall rule management
- Configurable via `powershell.json`

---

*Keep this file updated when scripts change.*
