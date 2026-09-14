# Plan 161: Archive URL Downloader, Download Caching, --download-must & Antigravity Logo Visibility Suite

## Overview
autonomously implement remote URL archive downloading, intelligent download caching and reuse, the --download-must flag, pure Go multi-resolution XDG icon registration, and fix the Antigravity desktop icon in the Ubuntu/GNOME App Center "Installed" section across GitMap and scripts-fixer.

## Key Goals
1. Remote URL Archive Downloader: Enable gitmap install tar <url> (and install zip, install gz, install archive) to accept remote HTTP/HTTPS/FTP URLs.
2. Download Caching & Reuse: Cache downloads in ~/.gitmap/downloads/ (fallback to tempdir.RepoTempDir('downloads')). If a valid archive already exists (>1KB, valid magic bytes), announce `Reusing valid download from cache: <path>` and reuse it without downloading again.
3. Download Cache Bypass Flag (--download-must): Support --download-must (and aliases --force-download, --redownload) to purge existing cached file and force a fresh download.
4. Multi-Tier Fast Downloader: Download using aria2c with 16 parallel connections (-s 16 -x 16), falling back to curl (-fL --retry 3), falling back to atomic streaming Go net/http client.
5. Cross-Command Antigravity Integration: Integrate download caching and --download-must into gitmap install agy (and gitmap agy install), eliminating premature defer os.Remove deletions of downloaded archives.
6. Antigravity Logo & Installed Section Visibility:
- Fix Icon=antigravity (unextended name, not absolute path) in .desktop file to satisfy AppStream specifications.
- Pure Go multi-resolution icon resizer (16, 24, 32, 48, 64, 128, 256, 512) deploying to ~/.local/share/icons/hicolor/<size>/apps/antigravity.png and pixmaps.
- Ensure index.theme exists in ~/.local/share/icons/hicolor and trigger gtk-update-icon-cache -f -t.
- Provide embedded high-res fallback icon in cli/assets/assets.go via go:embed.
7. scripts-fixer Alignment: Update d:/work/scripts-fixer/scripts/os/ubuntu/install-antigravity.sh and install-archive.sh to deploy all resolutions, ensure index.theme, update GTK cache, and use Icon=antigravity.
## Custom Rules
1. Functions strictly <= 15 lines (target <= 8 liness).
2. Affirmative booleans only (isDownloadMust, isValid, hasExisting). No negatives.
3. All errors wrapped in *apperror.AppError with proper codes.
4. Absolute ban on go test and go build during routine turns.

> **Task Origin**: Integration of intelligent archive installer & Antigravity logo fix.
> **Total Loops**: 1 continuous loop.

## Consolidated Subtasks Detail

# Subtask 161-01: Archive URL Downloader, Download Caching & --download-must

## Target Files
- cli/cmdinstall/install_archive_types.go
- cli/cmdinstall/install_archive_cmd.go
- cli/cmdinstall/install_archive_detect.go
- cli/cmdinstall/install_archive_download.go (NEW)
- cli/cmdinstall/install.go
- cli/cmdinstall/agy_install.go
- cli/cmdinstall/installantigravity_deploy_linux.go
- cli/cmdinstall/installantigravity_deploy_windows.go
- cli/cmdinstall/installantigravity_deploy_darwin.go

## Requirements
1. Define ArchiveDownloadParams and ArchiveCacheValidationResult in install_archive_types.go.
2. Implement install_archive_download.go:
   - isRemoteURL(raw string) bool
   - validateRemoteArchiveURL(rawURL string) error
   - filenameFromArchiveURL(rawURL string) string
   - resolveArchiveDownloadCacheDir() string (~/.gitmap/downloads, fallback tempdir.RepoTempDir('downloads'))
   - CheckCachedArchive(filePath string) ArchiveCacheValidationResult (>1KB, valid magic bytes)
   - VerifyArchiveHeaderMagic(filePath string) bool
   - FetchOrReuseArchive(params ArchiveDownloadParams) (string, error):
     - If !params.IsDownloadMust && cache is valid -> print 'Reusing valid download from cache: <path>' and return path.
     - If params.IsDownloadMust && file exists -> os.Remove(cachedPath).
     - Download via multi-tier: aria2c 16 connections, curl -fL, Go net/http.
3. Update parseInstallTarArgs in install_archive_cmd.go:
   - Register --download-must, --force-download, --redownload into opts.IsDownloadMust.
   - If isRemoteURL(opts.ArchivePath), call FetchOrReuseArchive before extract.
4. Update agy_install.go and install.go to bind --download-must flags into installOptions.
5. In installantigravity_deploy_linux.go, _windows.go, _darwin.go:
   - Replace defer os.Remove with FetchOrReuseArchive so downloads are retained and reused across runs!

---

# Subtask 161-02: Antigravity Desktop Icon, Multi-Resolution XDG & Installed Section

## Target Files
- cli/assets/assets.go (NEW)
- cli/cmdinstall/install_icon_types.go (NEW)
- cli/cmdinstall/install_icon_resize.go (NEW)
- cli/cmdinstall/install_icon_xdg.go (NEW)
- cli/cmdinstall/installantigravity_deploy_linux.go
- cli/cmdinstall/installantigravity_cleanup.go

## Requirements
1. Create cli/assets/assets.go with go:embed of icon-256.png as AntigravityDefaultIcon.
2. Create install_icon_types.go:
   - StandardIconResolutions = []int{16, 24, 32, 48, 64, 128, 256, 512}
   - AntigravityCanonicalIconName = 'antigravity'
   - HicolorIndexThemeContent string
3. Create install_icon_resize.go (Pure Go with zero CGO/shell dependencies):
   - Bilinear resizing using stdlib image/png.
   - GenerateMultiResolutionIcons(sourceBytes []byte) ([]ResizedIcon, error).
4. Create install_icon_xdg.go:
   - DeployMultiResolutionIcons(opts IconDeployOptions) error
   - Ensure ~/.local/share/icons/hicolor/index.theme exists.
   - Write all 8 resolutions to <hicolor>/<size>x<size>/apps/antigravity.png.
   - Write pixmaps fallback.
   - Run gtk-update-icon-cache -f -t if available.
   - Run update-desktop-database.
5. In installantigravity_deploy_linux.go:
   - Set Icon=antigravity (strictly unextended name, no absolute paths to satisfy AppStream / GNOME Installed section!).
   - Call DeployMultiResolutionIcons.
6. In installantigravity_cleanup.go:
   - Sweep all 8 icon resolutions and pixmaps on uninstall.

---

# Subtask 161-03: scripts-fixer Antigravity Icon & Caching Alignment

## Target Files
- d:/work/scripts-fixer/scripts/os/ubuntu/install-antigravity.sh
- d:/work/scripts-fixer/scripts/os/ubuntu/install-archive.sh

## Requirements
1. In d:/work/scripts-fixer/scripts/os/ubuntu/install-antigravity.sh:
   - Fix Icon=antigravity in desktop entry (do not use absolute path like $ICON_PATH).
   - Generate multi-resolution icons (16, 24, 32, 48, 64, 128, 256, 512) in every [Home/System] hicolor dir.
   - Ensure index.theme exists so gtk-update-icon-cache does not fail.
   - Run gtk-update-icon-cache -f -t on both user and system hicolor directories.
   - If no icon found in ide_dir, use embedded base64 fallback Antigravity PNG icon.
2. In install-archive.sh:
   - Ensure --download-must and --force pass through properly and cache reuse prints 'Reusing valid download from temp folder: ...'.