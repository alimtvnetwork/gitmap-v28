# Application Specification: Pure Binary Updater & Modern Terminal UI Installer

**Feature Slug:** `bug-fix-and-terminal-ui-improvement`  
**Status:** In Progress  
**Category:** CLI Tooling & Distribution  
**Scope:** `cli/cmdupdate/`, `install.ps1`, `install.sh`, `local-install.ps1`, `local-install.sh`

---

## 1. Executive Summary & Problem Statement

### 1.1 The Source-Build Updater Bug
When users execute `gitmap update`, the CLI previously executed a source-discovery cascade (`resolveRepoPath()` in `cli/cmdupdate/updaterepo.go`). If any repository path was discovered (embedded constants, SQLite settings, current working directory, or tracked repos), the updater launched `run.ps1 -Update`, which in turn invoked `git pull` and `go build`.
This is a critical architectural violation:
- An end-user updating a pre-compiled binary CLI must **never** be required to have Go, Git, or source repositories installed.
- Updates must download the pre-compiled, tested release binary directly from GitHub Releases (`gitmap-<tag>-windows-amd64.zip` / `gitmap-<tag>-linux-amd64.tar.gz`), verify SHA256 integrity, and swap the executable via a console-safe handoff process.

### 1.2 Terminal UI Deficiencies
Previous installer scripts (`install.ps1`, `install.sh`) displayed flat, unstyled ASCII text with muted 16-color escapes. Furthermore, `install.ps1` explicitly suppressed ANSI color rendering via `$PSStyle.OutputRendering = 'PlainText'`. The installer experience must be upgraded to a vibrant, Lipgloss-inspired 24-bit Truecolor terminal UI with rounded box cards, staged execution steps, interactive spinners, and completion status cards across Windows, macOS, and Linux. Additionally, `local-install.ps1` and `local-install.sh` must be implemented for offline/source-local deployments.

---

## 2. Architecture & Detailed Technical Requirements

### 2.1 Pure Binary Self-Update Engine (`cli/cmdupdate/`)
1. **Default Strategy Inversion:**
   - The default and primary update path must be **Binary Download & In-Place Handoff**.
   - Source building (`run.ps1 -Update`) is deprecated and must only execute if the user explicitly passes `--source-rebuild` or `--from-source`.
2. **Release Resolution & Integrity Gate:**
   - Resolve target release version via GitHub Releases API or sibling probe (`updateprobe.go`).
   - If current version matches target and `--force` is false, report clean up-to-date status and exit 0.
   - Download the platform asset (`.zip` for Windows, `.tar.gz` for POSIX) and `checksums.txt` into a clean temporary staging directory (`os.TempDir()`).
   - Verify SHA256 hash before touching any deployed files.
3. **In-Process Unpacking:**
   - Unpack `gitmap.exe` (Windows) or `gitmap` (POSIX) directly using standard Go packages (`archive/zip`, `archive/tar`, `compress/gzip`).
4. **Console-Safe Windows Handoff:**
   - Active `gitmap.exe` copies itself to `gitmap-update-<pid>.exe`.
   - Active `gitmap.exe` launches `gitmap-update-<pid>.exe update-apply --new-binary <tempPath> --target <targetPath>` and exits immediately with code 0.
   - `gitmap-update-<pid>.exe` renames running target `gitmap.exe` -> `gitmap.exe.old`, moves new binary to `gitmap.exe`, and launches `gitmap.exe update-cleanup` in a detached hidden process with 1500ms delay.
   - `update-cleanup` purges `gitmap-update-*.exe`, `gitmap.exe.old`, and temporary staging files.
5. **Modern Terminal UI in CLI:**
   - Utilize `pterm` spinners, multi-printers, and Lipgloss-style cards to display progress (Checking -> Downloading -> Verifying -> Applying -> Complete).
   - Suppress animations gracefully when `--quiet` or `--json` is supplied or when TTY is absent.

### 2.2 Modern Terminal UI Installers (`install.ps1` & `install.sh`)
1. **ANSI Truecolor & Box Drawing Design System:**
   - Enable Windows VT console mode (`SetConsoleMode(hOut, mode | 0x0004)`).
   - Remove `$PSStyle.OutputRendering = 'PlainText'`.
   - Dracula / Modern Cyber Palette:
     - Accent Cyan (`#8be9fd` / `\033[38;2;139;233;253m`)
     - Brand Purple (`#bd93f9` / `\033[38;2;189;147;249m`)
     - Success Green (`#50fa7b` / `\033[38;2;80;250;123m`)
     - Warning Orange (`#ffb86c` / `\033[38;2;255;184;108m`)
     - Danger Red (`#ff5555` / `\033[38;2;255;85;85m`)
     - Muted Gray (`#6272a4` / `\033[38;2;98;114;164m`)
   - Unicode rounded borders (`╭`, `╮`, `╰`, `╯`, `│`, `─`) with ASCII fallback if `$env:NO_COLOR` or `TERM=dumb`.
2. **Visual Structure:**
   - **Header Banner:** Universal GitMap CLI Installer title and repository link.
   - **Stage Cards (1 to 5):**
     - `[1/5] RESOLVE TARGET RELEASE` (version, channel, arch)
     - `[2/5] VERIFY & DOWNLOAD ASSET` (asset name, speed, sha256 check)
     - `[3/5] EXTRACT & INSTALL BINARY` (target path, alias creation)
     - `[4/5] SEED CONFIGURATION` (config files, docs site)
     - `[5/5] SHELL ENVIRONMENT ACTIVATION` (PATH, profile wrappers, setup check)
   - **Completion Summary Card:** Highlight installed version, binary location, quick start commands.
3. **Environment Parity:**
   - Ensure both `install.ps1` and `install.sh` respect `$GITMAP_DOWNLOAD_URL` / `$env:GITMAP_DOWNLOAD_URL` for isolated local or staging testing.

### 2.3 Standalone Local Installers (`local-install.ps1` & `local-install.sh`)
1. Provide offline local installation scripts in the repository root.
2. If Go compiler is installed, compile binary from local `./cli` source into `./bin/gitmap.exe`.
3. If pre-built binary already exists in `./bin/`, install directly.
4. Deploy to local user binary directory (`$env:LOCALAPPDATA\gitmap-cli` or `~/.local/bin/gitmap-cli`), copy configuration seeds, and configure profile.
5. Render using the exact same colorful Lipgloss ANSI terminal UI design system.

---

## 3. Non-Negotiable Quality & Coding Invariants

1. **Zero Builds or Test Suites (R1):** Routine agent steps do NOT run full `go test ./...` or `go build ./...`. Verification uses targeted file-scoped linters and pre-flight parameter checks.
2. **Positive Booleans Only:** In Go, TypeScript, PowerShell, and Bash, all boolean variables and functions use positive prefixes (`isBinaryUpdate`, `hasExtracted`, `isVTSupported`, `isValidChecksum`).
3. **Strict Relative Paths (R11):** No absolute paths or `file:///` URIs.
4. **Hyphenated GitMap Commits (R8/R9):** Single atomic commit via `gitmap cpb "<module> - <summary>"` or `gitmap cpf "<module> - <summary>"` without colons in the argument.
