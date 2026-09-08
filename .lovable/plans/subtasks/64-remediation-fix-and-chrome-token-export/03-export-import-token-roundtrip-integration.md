# Subtask 03 — Export & Import Token Vault Integration

**Parent Plan:** [.lovable/plans/pending/64-remediation-fix-and-chrome-token-export.md](../../pending/64-remediation-fix-and-chrome-token-export.md)  
**Status:** COMPLETED  
**Files:**
- `gitmap/cmd/chromeprofile_export.go`
- `gitmap/cmd/chromeprofile_smart_import.go`

---

## Objective
Integrate the token vault into `writeChromeExport` and `applyChromeExport` so that exported profile JSON snapshots contain the token vault variations and importing a profile restores the decrypted tokens into the destination profile's `Web Data`.

## Requirements
1. Update `chromeExport` struct in `chromeprofile_export.go` to include `TokenVault *ChromeTokenVault`.
2. In `writeChromeExport`, call `readChromeTokenService(srcProfile)` and attach to `exp.TokenVault`.
3. In `applyChromeExport`, if `exp.TokenVault` is present, restore tokens into `<dstProfile>/Web Data` table `token_service`.
4. Handle table creation if `token_service` does not yet exist in the destination database.
