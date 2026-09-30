# CI/CD Issue 90: US Locale `misspell` Failure on `Cancelled` in `change_password_cmd.go`

## 1. Reproduction
During the `v6.416.0` CI workflow (`Spell Check (misspell, US locale)`), step #7 (`Run misspell on changed files`) failed on `cli/cmdos/change_password_cmd.go`.

## 2. Cause
Line 36 of `cli/cmdos/change_password_cmd.go` printed `"Cancelled by user."` using British English spelling (`Cancelled` with double-L) instead of US English (`Canceled`), which violates the repository's `misspell -locale US` CI quality gate.

## 3. Fix
Replaced `"Cancelled by user."` with `"Canceled by user."` in `cli/cmdos/change_password_cmd.go`.

## 4. Prevention
Enforce US English spelling (`Canceled`, `Labeled`, `Initializes`) across all CLI user-facing prompts and log strings prior to release commits.
