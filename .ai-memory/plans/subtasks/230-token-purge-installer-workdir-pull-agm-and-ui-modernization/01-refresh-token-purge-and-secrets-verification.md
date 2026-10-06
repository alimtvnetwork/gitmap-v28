# Subtask 01: Refresh Token Purge and Secrets Verification

> **Parent Plan:** [230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md](../../pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)  
> **Spec Reference:** [02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/01-architecture-spec.md](../../../../02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/01-architecture-spec.md)  
> **Status:** `COMPLETED`  
> **Target Areas:**  
> - `/home/a/.cache/vmware/drag_and_drop/`  
> - `~/.antigravity_tools/accounts.json`  
> - `~/.antigravity_tools/accounts/*.json`  
> - `cli/gitignoreagm/sanitizer.go`  

---

## 1. Technical Objective

Conduct a comprehensive system-wide search across temporary directories, VMware shared caches, and user profile folders to locate and purge any unencrypted JSON files containing OAuth 2.0 refresh tokens (specifically `agm_accounts_backup_2026-10-06.json`). Verify that canonical credential stores maintain strict POSIX permissions and that no token files leak into git-tracked repositories.

---

## 2. Implementation & Remediation Details

1. **System Cache Audit & Purge:**
   - Identified exposed account backup in VMware drag-and-drop temporary cache:
     `/home/a/.cache/vmware/drag_and_drop/hfPh77/agm_accounts_backup_2026-10-06.json`.
   - Executed safe deletion of the unencrypted backup file and verified zero remnants in parent cache paths.
   - Verified that `/tmp` and other shared volatile directories are clear of any `accounts_backup*.json` archives.

2. **Canonical Credential Boundary Verification:**
   - Ensured active AGM credentials reside exclusively in the official configuration directory:
     - Directory: `~/.antigravity_tools/` (Permissions: `0700` / `rwx------`).
     - Main Accounts Store: `~/.antigravity_tools/accounts.json` (Permissions: `0600` / `rw-------`).
     - Sub-Account Profiles: `~/.antigravity_tools/accounts/*.json` (Permissions: `0600` / `rw-------`).

3. **Git Hygiene & Pre-Commit Protection:**
   - Checked repository status using `cli/gitignoreagm/sanitizer.go` rules.
   - Verified that `.gitignore` explicitly ignores `*.json`, `*tokens*.db`, and `*backup*.json` patterns across all working repositories.

---

## 3. Verification Commands & Evidence

```bash
# 1. Verify VMware drag and drop cache is clean
ls -la /home/a/.cache/vmware/drag_and_drop/hfPh77/ 2>&1

# 2. Verify file permissions on canonical credentials
stat -c "%a %n" ~/.antigravity_tools/accounts.json ~/.antigravity_tools/accounts/

# 3. Confirm zero git-tracked accounts or tokens
git ls-files | grep -iE 'agm_accounts|user_tokens|refresh_token' || echo "Git tracking clean"
```

---

## 4. Acceptance Criteria

- [x] Target unencrypted backup file deleted without residue.
- [x] Zero refresh tokens detected in untracked or volatile directories.
- [x] Canonical AGM credential directories maintain `0700` and files maintain `0600` permissions.
- [x] Repository tracking verified clean against secrets leakage.
