# Subtask 04: Ubuntu Google Chrome Installer & Special Handler

## Scope
- Create `gitmap/cmd/install_chrome_linux.go`:
  - Download official `.deb` (`https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb`) into staging `/tmp`.
  - Execute `sudo apt install -y /tmp/google-chrome-stable_current_amd64.deb` to handle dependencies and configure repository automatically.
  - Fallback: configure Google signing key (`/etc/apt/keyrings/google-chrome.gpg`) and sources list (`/etc/apt/sources.list.d/google-chrome.list`), then `sudo apt update && sudo apt install -y google-chrome-stable`.
- Update `gitmap/cmd/install.go`:
  - Wire `chrome` and `google-chrome` to `specialInstallHandler` invoking `runInstallChromeLinux` on Linux.
- Update `gitmap/cmd/chrome.go`:
  - In `runChromeInstall`, call `runInstallChromeLinux` when on Linux instead of generic package manager.
- Update `gitmap/constants/constants_install.go`:
  - Register `ToolChrome` and `ToolGoogleChrome` in `InstallToolDescriptions` and `InstallToolCategories`.

## Files Touched
- `gitmap/cmd/install_chrome_linux.go`
- `gitmap/cmd/install.go`
- `gitmap/cmd/chrome.go`
- `gitmap/constants/constants_install.go`
