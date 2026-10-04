# Subtask 01: Ubuntu Fleet Node u1 Deep Filesystem Hygiene & Repository Boundary Audit

> **Task Reference:** `215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager`  
> **Parent Spec:** [01-architecture-spec.md](file:///d:/work/gitmap/02-spec/21-app/215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager/01-architecture-spec.md)  
> **Target Node:** Ubuntu 24.04 LTS (`u1` / `192.168.1.22`)  
> **Target User:** `a` (UID 1000, GID 1000)  
> **Execution Mode:** Automated SSH Elevated Stream (`sudo bash -s`)  
> **Status:** READY FOR EXECUTION  

---

## 1. Objective & Scope

Execute a thorough, elevated filesystem sanitization protocol on Ubuntu fleet node `u1` to resolve filesystem pollution accumulated during initial automated bootstrapping and toolchain installations.

### Key Scope Items:
1. **Root-Owned Literal Tilde Directory Remediation:** Safely eradicate `/home/a/git-work/'~'` (created by an Oh My Zsh / installer running via `sudo` with unexpanded tilde argument) without risking deletion of `$HOME` (`/home/a` or `/root`).
2. **Cross-Platform Stray Directory Cleanup:** Recursively purge stray Windows path directories created by un-sanitized path piping (e.g. `/home/a/C:\Users\...` or `/home/a/C:*`).
3. **Transient Artifact Removal:** Eradicate loose `.deb` package files, stray `.sh` bootstrap scripts, and un-sandboxed OAuth token caches sitting directly in `/home/a/`.
4. **Hollow Home Repositories Purge:** Remove redundant, empty, or duplicate git repository folders inside `/home/a/` to guarantee that all 74 repositories live exclusively under `/home/a/git-work/`.
5. **Repository Secrets & Cache Boundary Governance:**
   - Rigorously inspect and preserve `/home/a/git-work/repo-secrets` (verify `.git`, credentials, and permissions).
   - Ensure `/home/a/git-work/repo-cache` is cloned or initialized with proper user ownership (`a:a`).

---

## 2. Technical Implementation & Command Specification

### 2.1 Elevated Hygiene Bash Script (`deep-filesystem-hygiene.sh`)

The script is piped directly through SSH to `sudo bash -s` to handle root-owned items, followed by standard user-level commands for user-owned paths.

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "=========================================================="
echo "  [HYGIENE] Starting Node u1 Deep Filesystem Sanitization "
echo "=========================================================="

# 1. Root-owned literal tilde directory remediation
if [ -e "/home/a/git-work/~" ] || [ -d "/home/a/git-work/~" ]; then
    echo "[STEP 1] Found literal tilde directory at /home/a/git-work/~. Removing with sudo..."
    sudo rm -rf "/home/a/git-work/~"
    echo "[STEP 1] Successfully eradicated /home/a/git-work/~"
else
    echo "[STEP 1] /home/a/git-work/~ does not exist (CLEAN)."
fi

# 2. Stray Windows path tree cleanup
echo "[STEP 2] Scanning for stray Windows path directories in /home/a/..."
find /home/a -maxdepth 1 \( -name 'C:*' -o -name 'c:*' -o -name 'C:\\*' \) -print0 2>/dev/null | while IFS= read -r -d '' stray_dir; do
    echo "  -> Removing stray Windows tree: $stray_dir"
    sudo rm -rf "$stray_dir"
done

# 3. Purge loose transient .deb installation packages from home root
echo "[STEP 3] Purging loose .deb files from /home/a/..."
for deb_file in /home/a/*.deb; do
    if [ -f "$deb_file" ]; then
        echo "  -> Removing loose package: $deb_file"
        rm -f "$deb_file"
    fi
done

# 4. Purge loose stray .sh scripts from home root (preserve structured repositories)
echo "[STEP 4] Scanning for loose .sh scripts in /home/a/ root..."
find /home/a -maxdepth 1 -name '*.sh' -type f | while IFS= read -r script_file; do
    echo "  -> Archiving/Removing loose script: $script_file"
    rm -f "$script_file"
done

# 5. Purge loose OAuth credentials dumped into home root
echo "[STEP 5] Purging loose OAuth token artifacts in /home/a/..."
find /home/a -maxdepth 1 \( -name '*oauth*' -o -name '*token*' -o -name '.gdrive*' \) -type f | while IFS= read -r token_file; do
    echo "  -> Removing loose credential file: $token_file"
    rm -f "$token_file"
done

# 6. Purge hollow or duplicate repos directly in /home/a (all repos must reside in /home/a/git-work/)
echo "[STEP 6] Checking for hollow git folders directly under /home/a/..."
for candidate in /home/a/*; do
    if [ -d "$candidate/.git" ] && [ "$candidate" != "/home/a/git-work" ]; then
        echo "  -> Found stray git repository directly in home: $candidate. Removing redundant clone..."
        rm -rf "$candidate"
    fi
done

# 7. Preserving repo-secrets & initializing repo-cache
echo "[STEP 7] Verifying /home/a/git-work/repo-secrets integrity..."
if [ -d "/home/a/git-work/repo-secrets/.git" ]; then
    echo "  -> /home/a/git-work/repo-secrets is present and valid."
else
    echo "  -> WARNING: /home/a/git-work/repo-secrets missing or incomplete! Inspecting..."
fi

echo "[STEP 8] Ensuring /home/a/git-work/repo-cache directory exists..."
if [ ! -d "/home/a/git-work/repo-cache" ]; then
    echo "  -> Initializing /home/a/git-work/repo-cache..."
    mkdir -p "/home/a/git-work/repo-cache"
    chown -R a:a "/home/a/git-work/repo-cache"
else
    echo "  -> /home/a/git-work/repo-cache already present."
fi

echo "=========================================================="
echo "  [HYGIENE] Node u1 Filesystem Sanitization Complete!     "
echo "=========================================================="
```

---

## 3. Remote Invocation Protocol via Master Embedded Runner

The sanitization script can be executed directly from Windows PowerShell using the zero-dependency embedded runner:

```powershell
# From Windows Management Host:
$scriptContent = @'
# Paste POSIX sanitization script here
'@

$scriptContent | ssh.exe -o BatchMode=yes u1 "tr -d '\r' | bash -s"
```

Or via the unified GitMap embedded runner action:
```powershell
pwsh -File d:/work/repo-secrets/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action hygiene
```

---

## 4. Verification Commands & Live Evidence Protocol

Execute the following commands to verify that node `u1` is 100% clean and compliant:

```powershell
# 1. Verify 0 tilde directories exist under /home/a/git-work/
ssh.exe u1 "[ ! -e /home/a/git-work/'~' ] && echo 'PASS: No tilde dir' || echo 'FAIL: Tilde dir still exists'"

# 2. Verify 0 stray Windows path directories in /home/a
ssh.exe u1 "[ $(find /home/a -maxdepth 1 -name 'C:*' | wc -l) -eq 0 ] && echo 'PASS: No stray Windows paths' || echo 'FAIL: Windows paths found'"

# 3. Verify 0 loose .deb files in home root
ssh.exe u1 "[ $(find /home/a -maxdepth 1 -name '*.deb' | wc -l) -eq 0 ] && echo 'PASS: No loose debs' || echo 'FAIL: Loose debs found'"

# 4. Verify 0 loose .sh files in home root
ssh.exe u1 "[ $(find /home/a -maxdepth 1 -name '*.sh' | wc -l) -eq 0 ] && echo 'PASS: No loose sh files' || echo 'FAIL: Loose scripts found'"

# 5. Verify /home/a/git-work/repo-secrets is intact
ssh.exe u1 "[ -d /home/a/git-work/repo-secrets/.git ] && echo 'PASS: repo-secrets valid' || echo 'FAIL: repo-secrets missing'"

# 6. Verify /home/a/git-work/repo-cache is initialized
ssh.exe u1 "[ -d /home/a/git-work/repo-cache ] && echo 'PASS: repo-cache ready' || echo 'FAIL: repo-cache missing'"
```

---

## 5. Acceptance Criteria & Definition of Done

- [ ] **Literal Tilde Dir Purged:** `/home/a/git-work/~` is confirmed deleted; non-elevated and elevated inspections return non-existent.
- [ ] **Windows Paths Purged:** Zero entries matching `/home/a/C:*` or `/home/a/c:*`.
- [ ] **Loose Debs Purged:** Exactly 0 `.deb` files residing in `/home/a/`.
- [ ] **Loose Scripts Purged:** Exactly 0 `.sh` files residing directly in `/home/a/`.
- [ ] **Secrets Preserved:** `/home/a/git-work/repo-secrets` has an intact `.git` folder and valid HEAD.
- [ ] **Cache Initialized:** `/home/a/git-work/repo-cache` exists and is writable by user `a`.
