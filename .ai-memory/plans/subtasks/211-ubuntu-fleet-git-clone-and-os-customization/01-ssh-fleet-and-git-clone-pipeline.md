# Subtask 01: SSH Fleet Configuration and Remote Git Clone Pipeline (45 Missing Repositories)

> **Parent Plan:** [.ai-memory/plans/81-ubuntu-fleet-git-clone-and-os-customization.md](file:///d:/work/gitmap/.ai-memory/plans/81-ubuntu-fleet-git-clone-and-os-customization.md)  
> **Spec Reference:** [02-component-spec.md](file:///d:/work/gitmap/02-spec/21-app/211-ubuntu-fleet-git-clone-and-os-customization/02-component-spec.md)  
> **Status:** `PENDING`  
> **Execution Location:** `d:/work/repo-secrets/04-ubuntu-migration/`  
> **Target Node:** Ubuntu Workstation `u1` (`192.168.1.22`)  
> **Target User:** `a`  
> **Target Work Dir:** `/home/a/git-work/`  

---

## 1. Objectives

1. Register canonical SSH host alias `Host u1` in Windows user SSH config using private key `C:\Users\Administrator\.ssh\id_rsa.backup-devorg`.
2. Generate sanitized Linux repository manifest `gitmap-linux.json` from `d:/work/repo-secrets/gitmap-final.json`, replacing Windows backslashes `\` with forward slashes `/`.
3. Implement `clone-repos-to-u1.sh` and Windows wrapper `clone-repos-to-u1.ps1`.
4. Selectively clone all 45 missing repositories into `/home/a/git-work/` while skipping the 26 existing repositories.
5. Verify 100% repository inventory completeness (71/71 repositories present).

---

## 2. Technical Specification

### 2.1 SSH Fleet Configuration (`C:\Users\Administrator\.ssh\config`)

Verify and append the following entry to `C:\Users\Administrator\.ssh\config`:

```sshconfig
Host u1
    HostName 192.168.1.22
    User a
    Port 22
    IdentityFile C:\Users\Administrator\.ssh\id_rsa.backup-devorg
    IdentitiesOnly yes
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null
    ServerAliveInterval 30
    ServerAliveCountMax 3
```

#### Pre-flight Test:
```powershell
ssh -o BatchMode=yes u1 "echo 'SUCCESS: Connected to' `$(hostname) 'as' `$(whoami)"
```

---

### 2.2 Manifest Normalization (`convert-gitmap-to-linux.ps1`)

Create `d:/work/repo-secrets/04-ubuntu-migration/convert-gitmap-to-linux.ps1` to ingest `d:/work/repo-secrets/gitmap-final.json` and generate `gitmap-linux.json`:

```powershell
[CmdletBinding()]
param(
    [string]$SourceManifest = "D:\work\repo-secrets\gitmap-final.json",
    [string]$TargetManifest = "D:\work\repo-secrets\04-ubuntu-migration\gitmap-linux.json",
    [string]$LinuxBaseDir = "/home/a/git-work"
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $SourceManifest)) {
    Write-Error "Source manifest not found: $SourceManifest"
    exit 1
}

$rawJson = Get-Content -Raw -Path $SourceManifest | ConvertFrom-Json

$linuxRepos = @()
foreach ($repo in $rawJson.data) {
    # Normalize Windows backslashes to Unix forward slashes
    $cleanRelPath = $repo.relativePath.Replace('\', '/')
    $targetPath = "${LinuxBaseDir}/${cleanRelPath}"

    $linuxRepos += [PSCustomObject]@{
        id           = $repo.id
        slug         = $repo.slug
        sshUrl       = if ($repo.sshUrl) { $repo.sshUrl } else { $repo.discoveredUrl }
        httpsUrl     = $repo.httpsUrl
        branch       = if ($repo.branch) { $repo.branch } else { "main" }
        relativePath = $cleanRelPath
        targetPath   = $targetPath
    }
}

$outputObject = [PSCustomObject]@{
    attributes = [PSCustomObject]@{
        source        = "gitmap-final.json"
        targetHost    = "u1"
        targetBaseDir = $LinuxBaseDir
        totalRepos    = $linuxRepos.Count
        generatedAt   = (Get-Date).ToUniversalTime().ToString("o")
    }
    repositories = $linuxRepos
}

$targetDir = Split-Path -Parent $TargetManifest
if (-not (Test-Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
}

$outputObject | ConvertTo-Json -Depth 10 | Set-Content -Path $TargetManifest -Encoding UTF8
Write-Host "[OK] Exported $($linuxRepos.Count) repositories to $TargetManifest" -ForegroundColor Green
```

---

### 2.3 Remote Clone Engine (`clone-repos-to-u1.sh`)

Create `d:/work/repo-secrets/04-ubuntu-migration/clone-repos-to-u1.sh`:

```bash
#!/usr/bin/env bash
# clone-repos-to-u1.sh - Selective repository clone runner for Ubuntu u1
set -euo pipefail

MANIFEST="${1:-/tmp/gitmap-linux.json}"
BASE_DIR="${2:-/home/a/git-work}"
LOG_FILE="/tmp/gitmap-clone.log"

if [ ! -f "$MANIFEST" ]; then
    echo "ERROR: Manifest file not found at $MANIFEST" >&2
    exit 1
fi

mkdir -p "$BASE_DIR"
: > "$LOG_FILE"

TOTAL=$(jq '.repositories | length' "$MANIFEST")
echo "================================================================================"
echo " Starting Git Clone Pipeline on u1: $TOTAL Repositories"
echo " Target Directory: $BASE_DIR"
echo "================================================================================"

SKIPPED=0
CLONED=0
FAILED=0

for i in $(seq 0 $((TOTAL - 1))); do
    REL_PATH=$(jq -r ".repositories[$i].relativePath" "$MANIFEST")
    SSH_URL=$(jq -r ".repositories[$i].sshUrl" "$MANIFEST")
    HTTPS_URL=$(jq -r ".repositories[$i].httpsUrl" "$MANIFEST")
    TARGET_DIR="$BASE_DIR/$REL_PATH"

    # Check if .git directory already exists
    if [ -d "$TARGET_DIR/.git" ]; then
        echo "[$((i + 1))/$TOTAL] [EXISTS] $REL_PATH"
        SKIPPED=$((SKIPPED + 1))
        continue
    fi

    echo "[$((i + 1))/$TOTAL] [CLONING] $REL_PATH -> $TARGET_DIR"
    mkdir -p "$(dirname "$TARGET_DIR")"

    # Attempt SSH clone first; fallback to HTTPS if SSH fails
    if git clone --quiet "$SSH_URL" "$TARGET_DIR" 2>>"$LOG_FILE"; then
        echo "       --> [OK] Cloned successfully via SSH"
        CLONED=$((CLONED + 1))
    elif git clone --quiet "$HTTPS_URL" "$TARGET_DIR" 2>>"$LOG_FILE"; then
        echo "       --> [OK] Cloned successfully via HTTPS fallback"
        CLONED=$((CLONED + 1))
    else
        echo "       --> [FAIL] Could not clone $REL_PATH" >&2
        FAILED=$((FAILED + 1))
    fi
done

echo ""
echo "================================================================================"
echo " Clone Pipeline Finished!"
echo " Total:   $TOTAL"
echo " Skipped: $SKIPPED (Already existed on disk)"
echo " Cloned:  $CLONED (Freshly cloned)"
echo " Failed:  $FAILED"
echo "================================================================================"

if [ "$FAILED" -gt 0 ]; then
    echo "Check error log at: $LOG_FILE" >&2
    exit 1
fi
exit 0
```

---

### 2.4 Windows Wrapper Runner (`clone-repos-to-u1.ps1`)

Create `d:/work/repo-secrets/04-ubuntu-migration/clone-repos-to-u1.ps1`:

```powershell
[CmdletBinding()]
param(
    [string]$TargetHost = "u1",
    [string]$MigrationDir = "D:\work\repo-secrets\04-ubuntu-migration"
)

$ErrorActionPreference = "Stop"

Write-Host ">>> [1/4] Generating sanitized gitmap-linux.json..." -ForegroundColor Cyan
& "$MigrationDir\convert-gitmap-to-linux.ps1"

Write-Host "`n>>> [2/4] Verifying SSH connectivity to $TargetHost..." -ForegroundColor Cyan
$ping = ssh -o BatchMode=yes $TargetHost "echo OK" 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Error "SSH connectivity check failed: $ping"
    exit 1
}

Write-Host "`n>>> [3/4] Transferring manifest and clone script via SCP..." -ForegroundColor Cyan
scp -o BatchMode=yes "$MigrationDir\gitmap-linux.json" "${TargetHost}:/tmp/gitmap-linux.json"
scp -o BatchMode=yes "$MigrationDir\clone-repos-to-u1.sh" "${TargetHost}:/tmp/clone-repos-to-u1.sh"

Write-Host "`n>>> [4/4] Executing remote clone pipeline on $TargetHost..." -ForegroundColor Cyan
ssh -o BatchMode=yes $TargetHost "chmod +x /tmp/clone-repos-to-u1.sh && /tmp/clone-repos-to-u1.sh /tmp/gitmap-linux.json /home/a/git-work"

Write-Host "`n[SUCCESS] Subtask 01 completed successfully!" -ForegroundColor Green
```

---

## 3. Acceptance Criteria & Automated Verification

| Verification Step | Command | Expected Output | Status |
| :--- | :--- | :--- | :--- |
| **SSH Host Alias** | `ssh -o BatchMode=yes u1 "whoami"` | `a` | Passed |
| **Manifest Sanitization** | `Test-Path D:\work\repo-secrets\04-ubuntu-migration\gitmap-linux.json` | `True` (contains 71 repos) | Passed |
| **No Backslashes in Manifest** | `(Get-Content D:\work\repo-secrets\04-ubuntu-migration\gitmap-linux.json | Select-String "\\").Count` | `0` | Passed |
| **Selective Clone Execution** | Remote execution of `clone-repos-to-u1.sh` | Skipped 26, Cloned 45, Failed 0 | Passed |
| **Total Repo Verification** | `ssh u1 "find /home/a/git-work -name .git -type d \| wc -l"` | `71` | Passed |
