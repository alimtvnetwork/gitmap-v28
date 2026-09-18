# 04 — Build Scripts

## Purpose

Define the cross-platform build scripts (`run.ps1` and `run.sh`) that
automate the full pipeline: pull → build → deploy. These scripts are
the primary way developers build and deploy the CLI tool.


## Flow Diagram

See [`images/build-scripts-flow.mmd`](images/build-scripts-flow.mmd)

---

## Script Responsibilities

Both scripts implement the same 4-step pipeline (with a Windows-only
sub-step before build):

```
[1/4]   Pull latest changes      (git pull, branch check)
[2/4]   Resolve dependencies     (go mod tidy)
[2.5/4] Generate Windows         (go-winres make — Windows only;
        resources                 see 11-windows-icon-embedding.md)
[3/4]   Build binary             (go build with ldflags)
[4/4]   Deploy                   (rename-first to resolved target)
```

Step 2.5 is mandatory on Windows because `.syso` files are not
committed (per [11-windows-icon-embedding.md](11-windows-icon-embedding.md)).
On Linux/macOS this step is skipped entirely.

---

## PowerShell Script (`run.ps1`)

### Parameters

```powershell
param(
    [switch]$NoPull,           # Skip git pull
    [switch]$NoDeploy,         # Skip deploy step
    [switch]$ForcePull,        # Discard local changes before pulling
    [string]$DeployPath = "",  # Override deploy target
    [switch]$Update,           # Enable update mode (rename-first PATH sync)
    [switch]$R,                # Run the binary after build
    [switch]$Test,             # Run unit tests instead of build
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RunArgs         # Arguments forwarded to the binary
)
```

### Configuration Loading

Read build settings from a JSON config file:

```powershell
function Load-Config {
    $configPath = Join-Path $ToolDir "powershell.json"
    if (Test-Path $configPath) {
        return Get-Content $configPath | ConvertFrom-Json
    }

    # Sensible defaults if config is missing
    return @{
        deployPath  = "<default-deploy-path>"
        buildOutput = "./bin"
        binaryName  = "<binary>.exe"
        copyData    = $true
    }
}
```

### Build Step

```powershell
function Build-Binary {
    param($Config)

    $binDir  = Join-Path $RepoRoot $Config.buildOutput
    $outPath = Join-Path $binDir $Config.binaryName
    mkdir -Force $binDir | Out-Null

    Push-Location $ToolDir
    try {
        # [2.5/4] Windows-only: regenerate .syso resources
        # See 11-windows-icon-embedding.md — .syso files are NOT
        # committed and MUST be regenerated before every Windows build
        if ($IsWindows) {
            Write-Step "2.5/4" "Generating Windows resources (.syso)"
            if (-not (Get-Command go-winres -ErrorAction SilentlyContinue)) {
                throw "go-winres not installed. Run: go install github.com/tc-hib/go-winres@latest"
            }
            go-winres make
            if ($LASTEXITCODE -ne 0) {
                throw "go-winres make failed (exit $LASTEXITCODE)"
            }
        }

        $absRepoRoot = (Resolve-Path $RepoRoot).Path
        $ldflags = "-X '<module>/constants.RepoPath=$absRepoRoot'"
        go build -ldflags $ldflags -o $outPath .
    } finally {
        Pop-Location
    }

    return $outPath
}
```

### Deploy Step

Uses `Resolve-DeployTarget` (see [02-deploy-path-resolution.md](02-deploy-path-resolution.md))
and rename-first deploy (see [03-rename-first-deploy.md](03-rename-first-deploy.md)).

### Post-Deploy PATH Sync

After deploying, check if the PATH binary differs from the deployed
binary. If so, sync them:

```powershell
$activeCmd = Get-Command <binary> -ErrorAction SilentlyContinue
if ($activeCmd) {
    $activeResolved = (Resolve-Path $activeCmd.Source).Path
    $deployedResolved = (Resolve-Path $deployedBinaryPath).Path

    if ($activeResolved -ne $deployedResolved) {
        # Use rename-first to sync the PATH binary
        $activeBackup = "$activeBinaryPath.old"
        Rename-Item $activeBinaryPath $activeBackup -Force
        Copy-Item $deployedBinaryPath $activeBinaryPath -Force
    }
}
```

---

## Bash Script (`run.sh`)

### Parameters

```bash
--no-pull       # Skip git pull
--no-deploy     # Skip deploy step
--force-pull    # Discard local changes before pulling
--deploy-path   # Override deploy target
--update        # Enable update mode
-r | --run      # Run the binary after build (remaining args forwarded)
-t | --test     # Run unit tests
```

### Configuration Loading

Use `python3` or `jq` to read JSON config:

```bash
get_config_value() {
    local key="$1"
    local default="$2"
    local config_path="$TOOL_DIR/powershell.json"

    if [[ -f "$config_path" ]] && command -v python3 &>/dev/null; then
        python3 -c "import json; d=json.load(open('$config_path')); print(d.get('$key','$default'))" 2>/dev/null || echo "$default"
    elif [[ -f "$config_path" ]] && command -v jq &>/dev/null; then
        jq -r ".$key // \"$default\"" "$config_path" 2>/dev/null || echo "$default"
    else
        echo "$default"
    fi
}
```

### Build Step

```bash
build_binary() {
    local bin_dir="$REPO_ROOT/$BUILD_OUTPUT"
    local out_path="$bin_dir/$BINARY_NAME"
    mkdir -p "$bin_dir"

    cd "$TOOL_DIR"

    # [2.5/4] Windows-only: regenerate .syso resources
    # See 11-windows-icon-embedding.md — .syso files are NOT
    # committed and MUST be regenerated before every Windows build.
    # Detect Git Bash / MSYS / Cygwin via uname.
    case "$(uname -s)" in
        MINGW*|MSYS*|CYGWIN*)
            write_step "2.5/4" "Generating Windows resources (.syso)"
            if ! command -v go-winres &>/dev/null; then
                write_fail "go-winres not installed. Run: go install github.com/tc-hib/go-winres@latest"
                exit 1
            fi
            go-winres make || { write_fail "go-winres make failed"; exit 1; }
            ;;
    esac

    local abs_repo_root
    abs_repo_root=$(cd "$REPO_ROOT" && pwd)
    local ldflags="-X '<module>/constants.RepoPath=$abs_repo_root'"

    go build -ldflags "$ldflags" -o "$out_path" .

    BINARY_PATH="$out_path"
}
```

### Platform Detection

```bash
BINARY_NAME="<binary>"
if [[ "$(uname -s)" == *MINGW* ]] || [[ "$(uname -s)" == *MSYS* ]]; then
    BINARY_NAME="<binary>.exe"
fi
```

### File Size Display

macOS and Linux use different `stat` flags:

```bash
if [[ "$(uname -s)" == "Darwin" ]]; then
    size=$(stat -f%z "$out_path")
else
    size=$(stat -c%s "$out_path")
fi
size_mb=$(echo "scale=2; $size / 1048576" | bc)
```

---

## Logging

Both scripts use the same color scheme:

| Color | Prefix | Purpose |
|-------|--------|---------|
| Magenta | `[1/4]` | Step headers |
| Green | `OK` | Success messages |
| Cyan | `->` | Informational messages |
| Yellow | `!!` | Warnings |
| Red | `XX` | Errors |

### PowerShell

```powershell
function Write-Step    { param($Step, $Message); Write-Host "  [$Step] " -Fore Magenta -NoNewline; Write-Host $Message }
function Write-Success { param($Message); Write-Host "  OK " -Fore Green -NoNewline; Write-Host $Message -Fore Green }
function Write-Info    { param($Message); Write-Host "  -> " -Fore Cyan -NoNewline; Write-Host $Message -Fore Gray }
function Write-Warn    { param($Message); Write-Host "  !! " -Fore Yellow -NoNewline; Write-Host $Message -Fore Yellow }
function Write-Fail    { param($Message); Write-Host "  XX " -Fore Red -NoNewline; Write-Host $Message -Fore Red }
```

### Bash

```bash
write_step()    { echo -e "  \033[0;35m[$1]\033[0m $2"; }
write_success() { echo -e "  \033[0;32mOK\033[0m \033[0;32m$1\033[0m"; }
write_info()    { echo -e "  \033[0;36m->\033[0m \033[0;37m$1\033[0m"; }
write_warn()    { echo -e "  \033[0;33m!!\033[0m \033[0;33m$1\033[0m"; }
write_fail()    { echo -e "  \033[0;31mXX\033[0m \033[0;31m$1\033[0m"; }
```

---

## Source File Validation

Before building, validate that critical source files exist:

```bash
required_files=("main.go" "go.mod" "cmd/root.go" "constants/constants.go")

missing=()
for file in "${required_files[@]}"; do
    [[ ! -f "$TOOL_DIR/$file" ]] && missing+=("$file")
done

if [[ ${#missing[@]} -gt 0 ]]; then
    echo "  XX Missing source files (${#missing[@]}):"
    for f in "${missing[@]}"; do echo "  - $f"; done
    exit 1
fi
```

---

## Data Folder Copy

If the config specifies `copyData: true`, copy the `data/` directory
alongside the binary (both to build output and deploy target):

```bash
if [[ -d "$TOOL_DIR/data" ]]; then
    rm -rf "$bin_dir/data"
    cp -r "$TOOL_DIR/data" "$bin_dir/data"
fi
```

---

## Git Pull with Conflict Resolution

The pull step should handle local changes gracefully:

```
Pull fails due to local changes?
├── --force-pull: auto-discard + clean + retry
├── Interactive: offer choices
│   ├── [S] Stash changes (git stash push, then retry)
│   ├── [D] Discard changes (git checkout -- ., then retry)
│   ├── [C] Clean all (discard + git clean -fd, then retry)
│   └── [Q] Quit
```

---

## Constraints

- Both scripts must read from the same JSON config file.
- Bash scripts require `bash 4+` (macOS ships `bash 3` — use
  `#!/usr/bin/env bash`).
- Error handling: `set -euo pipefail` in Bash,
  `$ErrorActionPreference = "Stop"` in PowerShell.
- All paths resolved to absolute before use.
- `go mod tidy` runs before every build for dependency consistency.
- Version is displayed immediately after build: `<binary> version`.
- **On Windows**, `go-winres make` MUST run before `go build` to
  regenerate `rsrc_windows_*.syso` resource files. The build script
  MUST fail loudly if `go-winres` is not installed — never silently
  skip resource generation, or the resulting `.exe` will lack its
  icon and version metadata. Install with:
  `go install github.com/tc-hib/go-winres@latest`.
  See [11-windows-icon-embedding.md](11-windows-icon-embedding.md).

## Application-Specific References

| App Spec | Covers |
|----------|--------|
| [11-build-deploy.md](../13-generic-cli/11-build-deploy.md) | Build/deploy scripting, config, logging, ldflags, and last-release detection |
| [17-self-update-app-update.md](../17-consolidated-guidelines/17-self-update-app-update.md) | Consolidated build-script and update workflow guidance |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
