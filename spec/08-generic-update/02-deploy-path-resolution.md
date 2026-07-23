# 02 — Deploy Path Resolution

## Purpose

Define how the build/deploy script determines **where to install** the
new binary. The resolution must handle first-time installs, existing
installs, and explicit overrides.

---

## Flow Diagram

See [`images/deploy-path-resolution-flow.mmd`](images/deploy-path-resolution-flow.mmd)

---

## Core Principle: Deploy to Running Location

The default deploy target is **the directory from which the currently
running executable was launched**. This ensures updates always replace
the active binary in-place, preserving the user's PATH configuration
and co-located data folder.

```
Priority 1: CLI flag override       → always wins
Priority 2: Running executable path → where the binary actually lives
Priority 3: Global PATH lookup      → binary found elsewhere on PATH
Priority 4: Interactive prompt      → ask the user (first-time install only)
Priority 5: Config file default     → fallback from config/JSON
```

---

## 4-Tier Resolution Strategy

| Priority | Source | When Used |
|----------|--------|-----------|
| 1 | **CLI flag** (`--deploy-path`) | User explicitly specifies a path |
| 2 | **Running executable location** | Binary is running; deploy to its own directory |
| 3 | **Global PATH lookup** | Binary is on PATH but invoked differently |
| 4 | **Interactive prompt / Config default** | First-time install or binary not found |

### Tier 1 — CLI Flag Override

If the user passes `--deploy-path <dir>`, use it unconditionally:

```powershell
# PowerShell
if ($DeployPath.Length -gt 0) {
    return $DeployPath
}
```

```bash
# Bash
if [[ -n "$DEPLOY_PATH" ]]; then
    echo "$DEPLOY_PATH"
    return
fi
```

### Tier 2 — Running Executable Location

Resolve the **physical location** of the currently running binary. This
is the primary default — the build script deploys back to where the
tool is already installed:

```powershell
# PowerShell — resolve the active binary's physical directory
$activeCmd = Get-Command <binary> -ErrorAction SilentlyContinue
if ($activeCmd) {
    $activePath = $activeCmd.Source
    $resolvedPath = (Resolve-Path $activePath).Path
    $activeDir = Split-Path $resolvedPath -Parent
    $dirName = Split-Path $activeDir -Leaf

    # Binary lives in <deploy-target>/<binary>/<binary>.exe
    if ($dirName -eq "<binary>") {
        return Split-Path $activeDir -Parent
    }

    return $activeDir
}
```

```bash
# Bash — resolve the active binary's physical directory
active_cmd=$(command -v <binary> 2>/dev/null || true)
if [[ -n "$active_cmd" ]] && [[ -f "$active_cmd" ]]; then
    resolved=$(readlink -f "$active_cmd" 2>/dev/null || echo "$active_cmd")
    active_dir=$(dirname "$resolved")
    dir_name=$(basename "$active_dir")

    if [[ "$dir_name" == "<binary>" ]]; then
        echo "$(dirname "$active_dir")"
        return
    fi

    echo "$active_dir"
    return
fi
```

#### Symlink Resolution

On Linux/macOS, the PATH binary may be a symlink. Always resolve
symlinks before extracting the directory:

```bash
resolved=$(readlink -f "$active_cmd")
```

On macOS, `readlink -f` may not be available. Fall back:

```bash
resolved=$(python3 -c "import os; print(os.path.realpath('$active_cmd'))" 2>/dev/null || echo "$active_cmd")
```

#### Nested Directory Detection

Most CLI deploy structures use a nested folder:

```
<deploy-target>/
└── <binary>/
    ├── <binary> (or <binary>.exe)
    └── data/
```

When detecting the running location, check if the binary's parent
directory matches the binary name. If so, the deploy target is the
**grandparent** directory.

### Tier 3 — Interactive Prompt (First-Time Only)

If the binary is not currently installed anywhere and no PATH entry
exists, prompt the user:

```
No existing installation found.
Where should <binary> be installed?

Default: C:\Users\<user>\bin-run (press Enter to accept)
> _
```

After the user provides a path:
1. Create the directory if it does not exist.
2. Deploy the binary into it.
3. Register the directory in the system PATH (see PATH Registration below).
4. Reload the terminal environment so the binary is immediately accessible.

### Tier 4 — Config File Default

If running non-interactively (CI, scripts), fall back to the config file:

**Windows default** (from `powershell.json`):
```json
{
  "deployPath": "E:\\bin-run"
}
```

**Linux/macOS default**:
```bash
DEPLOY_TARGET="$HOME/bin-run"
```

---

## PATH Registration

When deploying to a new directory (first-time install or changed path),
the directory **must** be added to the system PATH:

### Windows (User Environment Variable)

```powershell
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$deployDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$deployDir", "User")
    Write-Success "Added $deployDir to User PATH"
}
```

### Linux/macOS (Shell Profiles)

```bash
add_to_path() {
    local dir="$1"
    local profiles=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.profile")

    for profile in "${profiles[@]}"; do
        if [[ -f "$profile" ]] && ! grep -q "$dir" "$profile"; then
            echo "export PATH=\"$dir:\$PATH\"  # <binary>-path" >> "$profile"
        fi
    done
}
```

### Terminal Reload

After PATH registration, refresh the current session:

```powershell
# PowerShell — refresh PATH in the current session
$env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
            [Environment]::GetEnvironmentVariable("Path", "User")
```

```bash
# Bash — source the profile or export directly
export PATH="$deploy_dir:$PATH"
```

---

## Data Folder Co-Location

The data folder (`data/`) is **always co-located with the binary** at
its physical location. This means:

```
<binary-dir>/
├── <binary>.exe        ← the executable
└── data/
    ├── config.db       ← SQLite database
    └── ...             ← other data files
```

### Resolution in Code

```go
func BinaryDataDir() string {
    exe, err := os.Executable()
    if err != nil {
        return filepath.Join(".", "data")
    }

    resolved, err := filepath.EvalSymlinks(exe)
    if err != nil {
        resolved = exe
    }

    return filepath.Join(filepath.Dir(resolved), "data")
}
```

### Why Co-Location Matters

- The data folder is **always discoverable** — no config file needed.
- Moving the binary moves the data too.
- `run.ps1` / `run.sh` copies the data folder alongside the binary
  during deploy:

```powershell
# Copy data folder to deploy target alongside the binary
if ($Config.copyData) {
    $srcData = Join-Path $BuildOutput "data"
    $destData = Join-Path $appDir "data"
    if (Test-Path $srcData) {
        Copy-Item $srcData $destData -Recurse -Force
    }
}
```

---

## Complete Resolution Function

### PowerShell

```powershell
function Resolve-DeployTarget {
    param($Config, $OverridePath)

    # 1) CLI override
    if ($OverridePath.Length -gt 0) {
        Write-Info "Deploy target: CLI override -> $OverridePath"
        return $OverridePath
    }

    # 2) Running executable / PATH detection
    $activeCmd = Get-Command <binary> -ErrorAction SilentlyContinue
    if ($activeCmd) {
        $activePath = (Resolve-Path $activeCmd.Source).Path
        $activeDir = Split-Path $activePath -Parent
        $dirName = Split-Path $activeDir -Leaf

        if ($dirName -eq "<binary>") {
            $target = Split-Path $activeDir -Parent
        } else {
            $target = $activeDir
        }

        Write-Info "Deploy target: detected from running location -> $target"
        return $target
    }

    # 3) Config default
    Write-Info "Deploy target: config default -> $($Config.deployPath)"
    return $Config.deployPath
}
```

### Bash

```bash
resolve_deploy_target() {
    # 1) CLI override
    if [[ -n "$DEPLOY_PATH" ]]; then
        echo "$DEPLOY_PATH"
        return
    fi

    # 2) Running executable / PATH detection
    local active_cmd
    active_cmd=$(command -v <binary> 2>/dev/null || true)
    if [[ -n "$active_cmd" ]] && [[ -f "$active_cmd" ]]; then
        local resolved
        resolved=$(readlink -f "$active_cmd" 2>/dev/null || echo "$active_cmd")
        local active_dir
        active_dir=$(dirname "$resolved")
        local dir_name
        dir_name=$(basename "$active_dir")

        if [[ "$dir_name" == "<binary>" ]]; then
            echo "$(dirname "$active_dir")"
        else
            echo "$active_dir"
        fi
        return
    fi

    # 3) Config default
    echo "$DEPLOY_TARGET"
}
```

---

## Installed Directory Command

Provide a utility command for users to check where the binary is installed:

```
<binary> installed-dir
```

Output:

```
  Installed directory

  Binary:    /home/user/.local/bin/<binary>
  Directory: /home/user/.local/bin/
```

---

## Constraints

- CLI flag always takes highest priority — no exceptions.
- Default deploy target is the running executable's own directory.
- PATH detection must resolve symlinks before extracting directories.
- The config file default is only used when the binary is not found
  on PATH (first-time installs or CI environments).
- On first-time install, register the deploy directory in the system
  PATH and reload the terminal so the binary is immediately usable.
- Data folder is always co-located with the binary at its physical
  location — never stored in a separate global directory.
- Log which tier was used so the user can see where the binary will
  be deployed.

## Application-Specific References

| App Spec | Covers |
|----------|--------|
| [02-powershell-build-deploy.md](../03-general/02-powershell-build-deploy.md) | Deploy target resolution and nested deploy structure |
| [09-build-deploy.md](../01-app/09-build-deploy.md) | gitmap-v28 deploy target on PATH, nested `gitmap-v28/` subfolder |

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
