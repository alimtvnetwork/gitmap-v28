# gitmap clean-dev / gitmap os dev-clean

Developer Tools Cache Remover -- reclaims massive disk space occupied by compilers, developer runtimes, package managers, and build caches across polyglot stacks.

## Usage

```bash
gitmap clean-dev [flags]
gitmap os dev-clean [flags]
gitmap os cleanup dev [flags]
```

## Aliases

- `cleandev`, `dev-cleanup`, `devcleanup`, `dev-clean`, `clean-dev`

## Supported Developer Categories

| # | Category | Key | Description & Handled Locations |
|---|----------|-----|---------------------------------|
| 1 | Go | `go`, `go-buildcache` | `go clean -cache -modcache -testcache -fuzzcache` + `%LOCALAPPDATA%\go-build`, `%USERPROFILE%\go\pkg\mod`, `$DEV_DIR\go\pkg\mod`. Strips Windows read-only flags prior to deletion. |
| 2 | pnpm | `pnpm`, `pnpm-store` | `pnpm store prune` + `%LOCALAPPDATA%\pnpm\store`, `~/.pnpm-store`, `$DEV_DIR\pnpm\store`. Runtime binaries preserved. |
| 3 | npm | `npm`, `npm-cache` | `npm cache clean --force` + `%LOCALAPPDATA%\npm-cache`, `~/.npm`. |
| 4 | Chocolatey | `choco`, `choco-cache` | `choco cache clean -y` + `%LOCALAPPDATA%\Chocolatey\Cache`, `%ProgramData%\chocolatey\cache`. Package state preserved. |
| 5 | Yarn | `yarn`, `yarn-cache` | `yarn cache clean` + `%LOCALAPPDATA%\Yarn\Cache`, `~/.yarn/cache`. |
| 6 | Bun | `bun`, `bun-cache` | `bun pm cache rm` + `%LOCALAPPDATA%\bun\install\cache`, `~/.bun/install/cache`. |
| 7 | Python | `pip`, `python`, `pip-cache` | `pip cache purge` + `%LOCALAPPDATA%\pip\cache`, `~/.cache/pip`. |
| 8 | Cargo | `cargo`, `rust`, `cargo-cache` | `~/.cargo/registry/cache`, `~/.cargo/registry/src`, `~/.cargo/git/db`. Toolchains preserved. |
| 9 | .NET | `nuget`, `dotnet`, `nuget-cache` | `dotnet nuget locals all --clear` + `%LOCALAPPDATA%\NuGet\v3-cache`. |
| 10 | Gradle/Maven | `gradle`, `maven`, `gradle-maven-cache` | `~/.gradle/caches`, `~/.m2/repository`. |

## Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--dry-run` | `-n`, `-d` | Preview files, directories, and bytes that would be freed without deleting anything |
| `--yes` | `-y` | Bypass interactive confirmation prompt |
| `--json` | | Output execution results in structured JSON schema |
| `--verbose` | `-v` | Display individual subpaths, command notes, and non-fatal warnings |
| `--only <cats>` | | Filter sweep to specific comma-separated categories (e.g. `--only go,npm,pnpm`) |
| `--help` | `-h` | Show usage instructions |

## Examples

### 1. Preview Disk Space Reclaimable (Dry-Run)

```bash
gitmap clean-dev --dry-run
```

Output:

```text
  OS Dev-Cleanup (Developer Tools Cache Remover)
  ==============================================
  [DRY-RUN] Preview only. No cache files will be deleted.

  • go-buildcache:     Would clean 14205 files, 812 dirs (4120.50 MB freed)
  • pnpm-store:        Would clean 8912 files, 310 dirs (1840.10 MB freed)
  • npm-cache:         Would clean 1205 files, 45 dirs (320.40 MB freed)
  • pip-cache:         Would clean 410 files, 12 dirs (510.00 MB freed)

  ✔ [Dry-Run] Reclaimable Total: 6791.00 MB across 24732 files and 1179 dirs (82ms)
```

### 2. Clean Only Go and pnpm Caches Non-Interactively

```bash
gitmap clean-dev --only go,pnpm -y
```

### 3. Structured Machine Output (JSON)

```bash
gitmap clean-dev --dry-run --json
```

Output JSON:

```json
{
  "totalBytesFreed": 7120891230,
  "totalItemsRemoved": 24732,
  "totalDirsRemoved": 1179,
  "categories": [
    {
      "category": "go-buildcache",
      "label": "Go build cache + module downloads (~/go/bin SAFE)",
      "itemsRemoved": 14205,
      "dirsRemoved": 812,
      "bytesFreed": 4320670000
    }
  ],
  "isDryRun": true,
  "durationMs": 85
}
```
