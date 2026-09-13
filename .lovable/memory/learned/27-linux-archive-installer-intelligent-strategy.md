# Linux Archive Package Installation & Intelligent Strategy Engine

- **Date**: 2026-09-13
- **Plan Reference**: `.lovable/plans/completed/153-install-tar-gz-zip-linux.md`

## Overview
This memory record details the design, dispatching, format detection, and deployment pipeline for installing `.tar`, `.gz`, `.tgz`, and `.zip` archive packages on Linux via `gitmap install tar <archive-file>`.

## 1. Platform Guard & Command Dispatch
- **Linux-Only Enforcement**: The command enforces `runtime.GOOS == "linux"`. On Windows or macOS, execution halts with structured `*apperror.AppError` code `E_LINUX_ONLY`.
- **Command Routing**: `gitmap install tar <file>`, `gitmap in tar <file>`, or direct file invocation `gitmap install <file.tar.gz>` automatically dispatches to `runInstallTar`.
- **Arguments**: Supports optional `--name` / `-n` custom application names, `--dry-run`, and `-v`.

## 2. Extraction & Intelligent Package Strategies
- **Staging**: Always extracts into repository-scoped temporary storage (`tempdir.RepoTempDir("archive-install")`), keeping OS temp clean.
- **Auto-Flattening**: Leverages `archive.CompactExtract` to unwrap nested single-directory tarballs.
- **Decompression**: Standalone `.gz` files use `compress/gzip` directly to unpack single binary files.
- **Strategy Detection**:
  1. `StrategyBinaryApp`: Discovers ELF regular executables in `bin/` or root directory matching the app name. Copies app to `~/.local/share/<appname>/` (or `/opt/<appname>/`), creates execution symlink in `~/.local/bin/<appname>`, sets permissions `0755`, installs desktop launcher, and runs `update-desktop-database`.
  2. `StrategyScript`: Detects `install.sh` or `setup.sh` and runs it via `bash`.
  3. `StrategySource`: Detects `Makefile` / `configure` and executes source build.
  4. `StrategySingleGz`: Installs standalone binary into `~/.local/bin/<appname>`.

## 3. Step-by-Step Progress Feedback
Outputs 5 clear numbered steps:
- `[1/5]` Linux environment verification & archive path resolution.
- `[2/5]` Archive format inspection & extraction to staging.
- `[3/5]` Package analysis & strategy detection (`binary_app`, `install_script`, etc.).
- `[4/5]` Destination deployment & file installation.
- `[5/5]` PATH symlinking, desktop database refresh, and database recording.

## 4. Dual-Database Tracking & Uninstallation
- Records tool in `installation.db` and master `gitmap.db` with manager = `"archive"`.
- Recognizes custom archive applications in `IsCustomStandaloneTool` via disk path inspection.
- `gitmap uninstall <appname>` cleanly deletes `~/.local/share/<appname>`, `/opt/<appname>`, symlink `~/.local/bin/<appname>`, desktop file, and database tracking rows.
