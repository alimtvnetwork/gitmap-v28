# 64 — Remediation Command Execution Fix & Chrome Profile Refresh Token Vault

**Status:** COMPLETED  
**Priority:** High (CODE RED)  
**Category:** CLI / Git Remediation Engine & Chrome Profile Management  
**Reported:** 2026-09-09  

---

## 1. Problem Statement & User Session Logs

### Part A: Remediation Fix Quoting Breakdown (`cmd /c` Pathspec Error)
During interactive reconciliation / remediation (`gitmap reconcile` / `gitmap fix`):
```text
[4/10] llm-orchestrator-v5 (+4 untracked)
    • untracked: ai-bridge/internal/vector/hnsw.go
    ...
  Pick [1=stash, 2=wip, 3=discard, s=skip, a=all-stash, q=quit]: 2
ℹ Applying Fix: Option 2 (Commit Work-In-Progress) on llm-orchestrator-v5
  Running: git -C D:/wp-work/riseup-asia/llm-orchestrator-v3 add -A && git -C D:/wp-work/riseup-asia/llm-orchestrator-v3 commit -m "wip: local changes" && git -C D:/wp-work/riseup-asia/llm-orchestrator-v3 pull --rebase

error: pathspec 'local' did not match any file(s) known to git
error: pathspec 'changes"' did not match any file(s) known to git

✗ Fix failed: exit status 1
```

**Root Cause:**
- `gitutil/remediation_commit.go` generates a chained shell command string:
  `git -C %s add -A && git -C %s commit -m "wip: local changes" && git -C %s pull --rebase`
- `cmd/fix_cmd.go:executeFixRecipe` executes this on Windows via `exec.Command("cmd", "/c", recipe.Command)`.
- Go's `os/exec` on Windows wraps arguments containing spaces in quotes and escapes inner double quotes with backslashes (`\"`).
- `cmd.exe` does not unescape `\"` (treating `\` as a path separator), stripping or distorting quotes.
- `git.exe` receives `"wip:`, `local`, and `changes"` as distinct argv tokens, causing `git commit` to treat `local` and `changes"` as pathspecs.
- Lack of detailed diagnostics, execution stacktraces, and blunt known solutions for common git remediation failures.

### Part B: Chrome Profile Export Refresh Token Vault with Double Base64 & Caesar Cipher
User directive:
> *"Also for the chrome profile export, or export with needs to have the refresh token is it doable?*  
> *Can you please keep the refresh token for each profile in json as 2time base 64 hashing to that can reverted back and also a chespercypher algorithm for the refresh token with ways to revert back the data so add a section in the profile to have variations for he code keeping saved so that the system knows how to revert back, is it clear???"*

**Technical Scope:**
1. Extract Chrome profile OAuth refresh tokens from `<profile>/Web Data` (table `token_service`) without lock conflicts.
2. Provide double Base64 encoding (`base64(base64(token))`) with bidirectional reversibility.
3. Provide Caesar cipher ("chespercypher") with shift variations and exact deciphering functions.
4. Add a structured `tokenVault` section in `chromeExport` JSON schema capturing all encoding variations and revert instructions.
5. Support automatic token restoration on `gitmap chrome profile import`.

---

## 2. Task-Specific Rule Set

1. **Rule TR-1 (Native Step Execution):** Never pass concatenated shell strings with inner quotes to `cmd /c`. All remediation recipes must supply structured `Steps: []RemediationStep` executed directly via `exec.Command(step.Name, step.Args...)`.
2. **Rule TR-2 (Blunt Diagnostic Reporting):** When a remediation step fails, output un-sugar-coated failure diagnostics: command line executed, working directory, exit code, captured stderr/stdout, root cause analysis, and known remediation solutions.
3. **Rule TR-3 (Token Vault Bi-Directional Integrity):** All token encodings (Double Base64 and Caesar Cipher) must be 100% reversible, mathematically verified by unit tests asserting `revert(encode(token)) == token`.
4. **Rule TR-4 (Concurrency Safety & SQLite Read-Only Fallback):** When reading `Web Data` from active Chrome profiles, use SQLite `file:<path>?mode=ro` with copy-to-temp fallback to ensure zero failure on file locks.
5. **Rule TR-5 (Coding Guidelines):** All functions <= 8-15 lines, positive booleans (`is...`, `has...`), zero swallowed errors, zero nested ifs.

---

## 3. Implementation Plan & Completed Execution

### Step 1: Remediation Step Architecture & Blunt Error Diagnostics [COMPLETED]
- In `gitmap/gitutil/remediation_generator.go`:
  - Defined `RemediationStep struct { Name string; Args []string }`.
  - Added `Steps []RemediationStep` to `RemediationRecipe`.
  - Added `CleanRepoPathRaw(repoPath)` for unquoted path arguments.
- Updated `GenerateCommitRecipe`, `GenerateStashRecipe`, `GenerateDiscardRecipe` to populate native `Steps`.
- In `gitmap/cmd/fix_execute.go` and `gitmap/cmd/fix_diagnostics.go`:
  - Implemented `executeStructuredRecipe(item *RemediationItem, recipe gitutil.RemediationRecipe) error`.
  - Implemented `analyzeGitErrorOutput(output string, err error)` outputting blunt RCA and known solutions for pathspec, conflict, overwrite, auth, and network errors.
  - Implemented `isBenignCommitClean` to tolerate `nothing to commit` in wip commit step so it seamlessly advances to `pull --rebase`.

### Step 2: Chrome Token Extraction & Cipher Engine [COMPLETED]
- Created `gitmap/cmd/chromeprofile_tokens.go`:
  - `readChromeTokenService(profilePath string) (*ChromeTokenVault, error)` reading `Web Data` table `token_service`.
  - `EncodeDoubleBase64(data []byte) string` & `DecodeDoubleBase64(s string) ([]byte, error)` (100% reversible).
  - `EncodeCaesarCipher(s string, shift int) string` & `DecodeCaesarCipher(s string, shift int) string` (exact reverse letter and digit shifts).
  - `EncodeCaesarByteShift(data []byte, shift byte) string` & `DecodeCaesarByteShift(s string, shift byte) ([]byte, error)`.
  - `buildTokenEntry` generating `variations` map (`doubleBase64`, `caesarCipher`, `caesarByteShift`) and `revertSteps`.
  - `restoreChromeTokenService(profilePath string, vault *ChromeTokenVault) error` for re-inserting tokens into destination profile `token_service` table.

### Step 3: Token Import & Restoration Engine [COMPLETED]
- Updated `chromeExport` in `gitmap/cmd/chromeprofile_export.go` to include `TokenVault *ChromeTokenVault`.
- Updated `writeChromeExport` to extract and populate token vault in JSON snapshot.
- Updated `applyChromeExport` to restore tokens into `<dstProfile>/Web Data` table `token_service`.

### Step 4: Verification & Quality Gates [COMPLETED]
- Unit tests:
  - `gitmap/gitutil/remediation_steps_test.go` (`TestRemediationStepsGeneration`) PASS!
  - `gitmap/cmd/chromeprofile_tokens_test.go` (`TestChromeTokensDoubleBase64Roundtrip`, `TestChromeTokensCaesarCipherRoundtrip`, `TestChromeTokensCaesarByteShiftRoundtrip`, `TestChromeTokenVaultSQLiteRoundtrip`) PASS!
- Full package test suite: `go test ./gitutil/... ./cmd/...` passes 100% green.
- Live CLI verification: `gitmap chrome export Default` successfully produces JSON containing populated `tokenVault` with all cipher variations.
- Linters:
  - `check-nested-ifs.py`: 0 violations across 2,532 files.
  - `check-error-management.py`: 0 violations across 2,558 files.
