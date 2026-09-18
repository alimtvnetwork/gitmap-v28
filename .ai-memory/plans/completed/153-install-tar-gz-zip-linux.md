# Plan 153: Linux Archive Installer (`gitmap install tar <archive>`)

- **Slug**: `153-install-tar-gz-zip-linux`
- **Date**: 2026-09-13
- **Status**: completed
- **Execution Budget**: N = 150 loops
- **Loops Taken**: 2 self-loop iterations (Phase 1: Research, Skill Bootstrap & Decomposition; Phase 2: Implementation, Coverage & Consolidation)

---

## 1. Executive Summary & Task Genesis

### Task Genesis
The user requested dedicated command support to install `.tar`, `.gz`, or `.zip` archive packages in Ubuntu/Linux:
> `<cli> install tar <zip, tar, gz> # file installation, intelleigent enough to figure out tar, gz files with steps to instal for linux only, please`
> `clear???`

### Platform Scope
Archive package installation is strictly scoped to Linux/Ubuntu (`runtime.GOOS == "linux"`). When invoked on non-Linux platforms (Windows, macOS), the command halts cleanly with a structured `*apperror.AppError` (`E_LINUX_ONLY`) stating that archive package installation is supported on Linux / Ubuntu only.

---

## 2. Architecture & Intelligent Strategy Detection

### A. Supported Formats
Supports automatic identification and decompression of:
- Compressed Tarballs: `.tar.gz`, `.tgz`, `.tar.xz`, `.txz`, `.tar.bz2`, `.tbz2`, `.tar.zst`, `.tzst`
- Uncompressed Tarballs: `.tar`
- Zip Archives: `.zip`
- Standalone Compressed Binaries: `.gz`

### B. Extraction Pipeline
- Extracts packages into repository-scoped temporary storage (`tempdir.RepoTempDir("archive-install")`).
- Uses `archive.CompactExtract` to automatically flatten unnecessary single root wrapper folders.
- Single `.gz` files are decompressed directly using `compress/gzip`.

### C. Four Intelligent Installation Strategies
The extractor walks the package directory and automatically selects the optimal installation path:
1. **Pre-Compiled Binary Application (`StrategyBinaryApp`)**:
   - Searches for executable ELF binaries (`\x7fELF` magic bytes or `0111` execution bits) in `bin/` or root directory matching the application name.
   - Deploys application directory to `~/.local/share/<appname>/` (user mode) or `/opt/<appname>/` (root mode).
   - Creates execution symlink in `~/.local/bin/<appname>` (or `/usr/local/bin/<appname>`) with `0755` permissions.
   - Installs `.desktop` file (or creates one if an icon is detected) in `~/.local/share/applications/` and refreshes GNOME cache via `update-desktop-database`.
2. **Installer Script (`StrategyScript`)**:
   - Detects `install.sh` or `setup.sh`.
   - Executes script via `bash` in the extracted root.
3. **Source Build (`StrategySource`)**:
   - Detects `configure` script and `Makefile`.
   - Executes `./configure` and `make install`.
4. **Standalone Compressed Binary (`StrategySingleGz`)**:
   - Directly extracts single compressed binary file to `~/.local/bin/<appname>`, setting executable permissions `0755`.

### D. Step-by-Step Progress Display
Every installation outputs five distinct, numbered steps:
```
Installing archive package: myapp.tar.gz
  [1/5] Verifying Linux environment and resolving archive: myapp.tar.gz
  [2/5] Inspecting archive format and extracting to staging...
  [3/5] Analyzing package structure and detecting installation strategy...
      -> Detected strategy: binary_app
  [4/5] Deploying package files to destination for myapp...
  [5/5] Finalizing PATH symlinks, desktop launcher, and database records...

  ✓ myapp successfully installed!
```

---

## 3. Database Tracking & Universal Uninstaller Parity

1. **Dual Database Recording**:
   - Successfully installed archive packages are recorded in both `installation.db` and master `gitmap.db` with package manager = `"archive"`.
2. **Universal Uninstallation (`gitmap uninstall <appname>`)**:
   - Recognizes archive applications automatically via `IsCustomStandaloneTool` and `hasInstalledArchiveFiles`.
   - Purges `~/.local/share/<appname>`, `/opt/<appname>`, symlink `~/.local/bin/<appname>`, desktop file `~/.local/share/applications/<appname>.desktop`, and removes records from both databases.

---

## 4. Consolidated Subtasks

### Subtask 01: Install Tar CLI Dispatch and Validation
- **Target**: CLI dispatch in `cmdinstall`, argument parsing for `install tar`, direct archive file detection, and Linux-only platform guards.
- **Status**: Completed.
- **Key Files**:
  - `cli/cmdinstall/install.go`
  - `cli/cmdinstall/exports.go`
  - `cli/cmdinstall/install_archive_cmd.go`

### Subtask 02: Archive Extractor and Strategy Detector
- **Target**: Intelligent format detection, extraction to staging, package analysis (ELF binary, scripts, makefile, gz), and installation step runner.
- **Status**: Completed.
- **Key Files**:
  - `cli/cmdinstall/install_archive_types.go`
  - `cli/cmdinstall/install_archive_detect.go`
  - `cli/cmdinstall/install_archive_extract.go`

### Subtask 03: Deployment, Desktop and DB Tracking
- **Target**: Deployment to `~/.local/share/` or `/opt/`, PATH symlink creation (`0755`), `.desktop` registration, dual-database tracking, and uninstaller integration.
- **Status**: Completed.
- **Key Files**:
  - `cli/cmdinstall/install_archive_deploy.go`
  - `cli/cmdinstall/install_uninstall_paths.go`
  - `cli/cmdinstall/install_uninstall_custom.go`
  - `cli/cmdinstall/install_archive_test.go`
  - `cli/helptext/install.md`
  - `readme.md`

---

## 5. Verification & Quality Gates

1. **Unit Tests**:
   - `cli/cmdinstall/install_archive_test.go`:
     - `TestIsInstallTarCommand_IdentifiesArchiveCommands`
     - `TestIsInstallTarCommand_RejectsNonArchiveCommands`
     - `TestExtractInstallTarArgs_ExtractsPositionalArgs`
     - `TestDeriveArchiveAppName_StripsKnownExtensions`
     - `TestParseInstallTarArgs_ParsesFlagsAndAppName`
     - `TestParseInstallTarArgs_RejectsMissingPath`
     - `TestVerifyLinuxArchivePlatform_ValidatesCurrentOS`
     - `TestResolvePackageStrategy_ResolvesBinaryApp`
     - `TestResolvePackageStrategy_ResolvesScriptAndSource`
2. **Coding Guidelines Adherence**:
   - All newly authored/refactored functions adhere strictly to <= 15 lines.
   - Affirmative booleans (`is*`, `has*`) strictly enforced.
   - Zero explicit `== true` evaluations.
   - All errors use structured `*apperror.AppError` with stack traces.
3. **Execution Constraints**:
   - Total ban on routine test running (`go test`) and build checking (`go build`) strictly maintained.
