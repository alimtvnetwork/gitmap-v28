# Spec 142: Chrome Profile Auth Session, Refresh Token & Cookies Full-Fidelity Export/Import

**Version:** 1.0.0
**Updated:** 2026-09-24
**Status:** Canonical
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## User Request (Verbatim)

```text
I want you to test out chrome export, import profiles here with e2e local testing not for cicd or run these tests anytime later, flag those for temp  running only

I want you to test out the Chrome profile import, export properly. Okay? So you can do anything. Any removal or anything, it would not harm the system. Okay, so the system is VM, so you could do anything in the Chrome profile. It's not going to destroy anything. I give you the full permission to do whatever that you like, but make sure that the Chrome profile end-to-end testing works properly. So make sure that you have the refresh token and everything exported and tested properly here. Yeah. Make sure that the code is effective enough so that it does not break, because when you start to import, first try to do it with one account. Okay? One account at a time, so that you don't lose all the account, because you need the accounts to test this out, right? So let's say if you export and try to import and the import does not work, then you will be doomed. So what you could do is take one profile, for instance, try to export that properly, and then remove it, and then try to re-import and open the Chrome, go to Gmail, see that if it is logged in. If it does not, let's say you did not have the refresh token and other stuff, that means it does not work. Now, in that case, the recommendation is that you fix your code basically, to export the refresh token and try it with another account. Okay? And you do it until you prove yourself that it works. It's bulletproof, nothing is wrong there, we can trust the code, and then and again, then you proceed with the confidence. Is it clear? Can you please test these all the import and export profile of the Chrome? Only in this VM, do not use SSH to deal with the other machines Chrome, please. Okay? Is it understood? Can you please help me with this?
```

---

## 1. Problem Statement & Root Cause Analysis

When a user exports and re-imports a Google Chrome profile using `gitmap chrome export` and `gitmap chrome import`:
Chrome opens but all Google accounts and Gmail sessions are lost (user is redirected to `https://accounts.google.com/ServiceLogin` or `signin/identifier`).

Four distinct failure mechanisms caused this session loss:
1. **Omitted Token Service Restoration:** `importSingleSnapshotWithStepLogging` in `chromeprofile_smart_import.go` parsed `exp.TokenVault` from the JSON snapshot, but NEVER called `restoreChromeTokenService(dest.Path, exp.TokenVault)`. The SQLite table `token_service` in `Web Data` was never restored.
2. **Aggressive Preferences Authentication Scrubbing:** `patchImportedChromeProfilePreferences` hardcoded `keepSignin = false`. When `keepSignin` was false, `scrubImportedPreferencesAuth` explicitly deleted `account_info`, `google`, `gaia_cookie`, `sync`, and injected `signin.allowed = false` into `Preferences`. Chrome's startup sequence read `signin.allowed = false` and forcibly dropped all remaining Google authentication.
3. **Omission of Session Cookies:** `writeChromeExport` only captured Bookmarks, Preferences, and TokenVault. The session cookies database (`<ProfileDir>/Network/Cookies`) containing `SID`, `HSID`, `SSID`, `SAPISID` was omitted from exports.
4. **Stripping of Local State GAIA Metadata:** `registerChromeProfileWithFullSchema` deleted `gaia_id`, `gaia_name`, and `gaia_given_name` from `Local State` (`profile.info_cache`), so Chrome failed to associate the profile tile with the primary Google account.

---

## 2. Architecture & Remediation Contracts

### 2.1 Full-Fidelity Export Model (`chromeExport`)
`chromeExport` is extended with:
- `GaiaID string`: Primary account GAIA identifier (e.g. `113336085041923585255`).
- `GaiaName string`: User's full name.
- `GaiaGivenName string`: User's first name.
- `CookiesRawBase64 string`: Base64-encoded binary payload of `<ProfileDir>/Network/Cookies` (or `<ProfileDir>/Cookies`).
- `WebDataRawBase64 string`: Base64-encoded binary payload of `<ProfileDir>/Web Data`.
- `TokenVault *ChromeTokenVault`: Structured extract of all `token_service` rows for multi-cipher reversibility.

### 2.2 Preservation of Preferences Authentication
- `patchImportedChromeProfilePreferences` defaults to `keepSignin = true`.
- When `keepSignin = true`:
  - `account_info` array is strictly preserved.
  - `google.services.signin` is strictly preserved.
  - `gaia_cookie` is strictly preserved.
  - `signin.allowed` is set to `true`.

### 2.3 Comprehensive Import Pipeline
During `importSingleSnapshotWithStepLogging` and `applyChromeExport`:
1. Restore `Bookmarks` and `Preferences`.
2. Apply `patchImportedChromeProfilePreferencesWithOptions(dest.Path, dest.DisplayName, true)`.
3. If `WebDataRawBase64` is present, decode and write `<dest>/Web Data`.
4. Restore `TokenVault`: call `restoreChromeTokenService(dest.Path, exp.TokenVault)`.
5. If `CookiesRawBase64` is present, decode and write `<dest>/Network/Cookies`.
6. Register in `Local State` with `GaiaID`, `GaiaName`, `GaiaGivenName`, and `Email` preserved.

### 2.4 ZIP Archive Parity
- In `exportChromeProfileSQLite`: explicitly add `Network/Cookies` and `Cookies` to the zip payload.
- In `extractZipFile`: ensure parent directories (`os.MkdirAll(filepath.Dir(dest))`) are created so nested paths unpack reliably.

---

## 3. Local E2E Verification & Isolation Gate

- All local E2E tests are implemented in `cli/cmdchromeprofile/chrome_profile_auth_e2e_test.go`.
- Tagged with `//go:build e2e` so CI/CD quality runners (`06-cicd-local-runner.py run-smart`, `go test ./...`) skip them automatically.
- Verification tests round-trip export, removal, and re-import with verification of:
  - `Web Data` table `token_service` row count and token byte size matching source.
  - `Network/Cookies` SQLite table `cookies` Google auth cookie count matching source.
  - `Preferences` `account_info` array containing valid GAIA ID and email.
  - `Local State` `info_cache` entry containing `gaia_id` and `user_name`.
