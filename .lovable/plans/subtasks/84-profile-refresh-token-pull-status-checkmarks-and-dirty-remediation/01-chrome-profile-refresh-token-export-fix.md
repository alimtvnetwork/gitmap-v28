# Subtask 01: Chrome Profile Refresh Token Export Fix

## Objective
Ensure all Chrome profile exports (single-profile, multi-profile JSON, YAML, and SQLite) properly extract, encode, and serialize Google OAuth refresh tokens from `Web Data` (`token_service` table) into `TokenVault`.

## Files to Touch
- `gitmap/cmd/chromeprofile_export_all.go`
- `gitmap/cmd/chromeprofile_tokens.go`
- `gitmap/cmd/chromeprofile_tokens_test.go`

## Detailed Implementation Steps
1. In `gitmap/cmd/chromeprofile_export_all.go`:
   - Refactor `loadSingleChromeProfileExport(name string)` to call `readChromeTokenService(srcPath)`:
     ```go
     exp.TokenVault, _ = readChromeTokenService(srcPath)
     ```
   - Resolve `displayName` and `email` using `resolveProfileNameAndEmail(name, exp.Preferences)`.
   - In `initChromeSQLiteTables`, add `chrome_tokens` table schema:
     ```sql
     CREATE TABLE IF NOT EXISTS chrome_tokens (
         profile_name TEXT,
         service TEXT,
         account_id TEXT,
         raw_base64 TEXT,
         double_base64 TEXT,
         PRIMARY KEY (profile_name, service)
     );
     ```
   - In `populateProfilesInSQLite`, insert token entries into `chrome_tokens` when `exp.TokenVault` is non-nil.
2. Ensure all functions are $\le 15$ lines, blank line before returns, zero nested `if` statements.
3. Author/update unit tests in `gitmap/cmd/chromeprofile_tokens_test.go` to verify `loadSingleChromeProfileExport` includes `TokenVault`.
