# Developer Tools Cache Cleanup (`clean-dev` / `os dev-clean` / `dev-cleanup`)

> **Status:** Authoritative Feature Specification & Implementation Record
> **Command:** `gitmap clean-dev`, `gitmap os dev-clean`, `gitmap os dev-cleanup`
> **Package:** `cli/osclean`, `cli/cmdos/os_dev_clean.go`, `cli/cmd/clean_dev_entry.go`

---

## 1. Problem Statement & Motivation
Developers working on Windows and cross-platform polyglot stacks accumulate massive build caches, module dependencies, and package archives:
- **Go**: `go build` caches (`%LOCALAPPDATA%\go-build`) and module downloads (`pkg/mod`) routinely exceed 10+ GB with >100,000 files. Crucially, Go marks all downloaded module files as Read-Only (`[System.IO.FileAttributes]::ReadOnly` / mode `0444`), causing standard file removal on Windows to fail with `access denied` (`UnauthorizedAccessException`).
- **pnpm**: Global content-addressable store (`pnpm/store`) and custom `dev-tool\pnpm\store` directories grow indefinitely without pruning.
- **npm**: Cache archives in `%LOCALAPPDATA%\npm-cache` and `~/.npm`.
- **Chocolatey**: Package installer download archives in `%LOCALAPPDATA%\Chocolatey\Cache`, `%ProgramData%\chocolatey\cache`, and `%TEMP%\chocolatey`.
- **Other Package Managers**: Yarn, Bun, Python/pip, Cargo/Rust, Gradle, Maven, and .NET/NuGet.

---

## 2. Solution Architecture

### 2.1 Routing & Aliases
- Accessible via top-level command: `gitmap clean-dev` (aliases: `cleandev`, `dev-cleanup`, `devcleanup`, `dev-clean`).
- Accessible via OS dispatcher: `gitmap os dev-clean`, `gitmap os dev-cleanup`, `gitmap os cleanup dev`, `gitmap os clean dev`.
- Flags:
  - `--dry-run`, `-n`, `-d`: preview reclaimable space without modifying files.
  - `--yes`, `-y`: bypass confirmation prompt.
  - `--json`: machine-readable JSON output.
  - `--verbose`, `-v`: show notes and non-fatal warnings per category.
  - `--only <categories>`: comma-separated filter list.

### 2.2 10-Step Sweeper Engine (`cli/osclean/`)
1. **Go**: `go clean -cache -modcache -testcache -fuzzcache` + disk sweeps across `GOCACHE`, `GOMODCACHE`, `%USERPROFILE%\go\pkg\mod`, and `$DEV_DIR\go\pkg\mod`. Strips read-only attributes before removal.
2. **pnpm**: `pnpm store prune` + disk sweeps across `%LOCALAPPDATA%\pnpm\store`, `~/.pnpm-store`, `$DEV_DIR\pnpm\store`, and `C:\dev-tool\pnpm`.
3. **npm**: `npm cache clean --force` + directory sweep.
4. **Chocolatey**: `choco cache clean -y` + directory sweep across `%LOCALAPPDATA%\Chocolatey\Cache`, `%ProgramData%\chocolatey\cache`, and `%TEMP%\chocolatey`.
5. **Yarn**: `yarn cache clean` + directory sweep.
6. **Bun**: `bun pm cache rm` + directory sweep.
7. **Python / pip**: `pip cache purge` + directory sweep.
8. **Cargo / Rust**: `~/.cargo/registry/cache`, `~/.cargo/registry/src`, `~/.cargo/git/db`.
9. **.NET / NuGet**: `dotnet nuget locals all --clear` + `%LOCALAPPDATA%\NuGet\v3-cache`.
10. **Gradle / Maven**: `~/.gradle/caches`, `~/.m2/repository`.

### 2.3 Safe Invariant
All executable binaries in `~/go/bin`, global npm packages, rustup toolchains, installed software packages in Chocolatey (`lib`), and project source code remain strictly preserved.
