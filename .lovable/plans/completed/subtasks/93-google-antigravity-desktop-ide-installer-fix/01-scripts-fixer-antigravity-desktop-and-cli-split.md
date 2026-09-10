# Subtask 93.01: Scripts-Fixer Antigravity Desktop & CLI Split

## Goal
In `D:\work\scripts-fixer`, refactor script `69-install-antigravity` to install the **Google Antigravity Desktop IDE Application**, and create `78-install-antigravity-cli` for the `agy` command-line tool.

## Files Impacted
- `D:\work\scripts-fixer\scripts\69-install-antigravity\run.ps1`
- `D:\work\scripts-fixer\scripts\os\ubuntu\install-antigravity.sh`
- `D:\work\scripts-fixer\scripts\78-install-antigravity-cli\run.ps1` (NEW)
- `D:\work\scripts-fixer\scripts\os\ubuntu\install-antigravity-cli.sh` (NEW)
- `D:\work\scripts-fixer\scripts\shared\install-keywords.json`
- `D:\work\scripts-fixer\scripts\os\ubuntu\profile-ubuntu-dev-ai.sh`
- `D:\work\scripts-fixer\run.ps1`
- `D:\work\scripts-fixer\scripts\run.sh`

## Acceptance Criteria
1. `69-install-antigravity\run.ps1` downloads `Antigravity-x64.exe` from Google Storage, checks `%LOCALAPPDATA%\Programs\Antigravity\Antigravity.exe`, runs silent install `/S`, and supports `uninstall`.
2. `install-antigravity.sh` downloads `Antigravity.tar.gz` from Google Storage, extracts to `/opt/antigravity` (or `~/.local/share/antigravity`), creates `/usr/local/bin/antigravity` symlink, and installs `antigravity.desktop`.
3. Dedicated `78-install-antigravity-cli\run.ps1` and `install-antigravity-cli.sh` maintain the `agy` CLI tool installation.
4. Keyword mappings: `"antigravity"` -> 69 (Desktop), `"agy"` / `"antigravity-cli"` -> 78 (CLI).
