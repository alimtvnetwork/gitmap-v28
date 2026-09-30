# CI/CD Issue 90: US Locale `misspell`, `apperror.NewExecutionError`, and `gofmt` Fixes in `v6.416.0` / `v6.416.1`

## 1. Reproduction
During the `v6.416.0` CI workflow:
- `Spell Check (misspell, US locale)` failed on `cli/cmdos/change_password_cmd.go`.
- `go build` / `Verify test compilation` failed on `cli/cmdos/change_password_cmd.go:116` (`undefined: apperror.NewUnsupportedError`).
- `test_ci_scripts.py::test_gofmt_check_clean_repo` flagged formatting drift in `cli/cmdos/exports.go` and `cli/jsonenvelope/envelope.go`.

## 2. Cause
- Line 36 of `cli/cmdos/change_password_cmd.go` used British English `"Cancelled by user."` instead of US English `"Canceled by user."`.
- Line 116 of `cli/cmdos/change_password_cmd.go` referenced `apperror.NewUnsupportedError`, whereas `cli/apperror` provides `apperror.NewExecutionError`.
- `cli/cmdos/exports.go` and `cli/jsonenvelope/envelope.go` had extra blank lines detected by `gofmt -l`.

## 3. Fix
- Replaced `"Cancelled by user."` with `"Canceled by user."` in `cli/cmdos/change_password_cmd.go`.
- Replaced `apperror.NewUnsupportedError` with `apperror.NewExecutionError` in `cli/cmdos/change_password_cmd.go`.
- Applied `gofmt -w` on `cli/cmdos/exports.go` and `cli/jsonenvelope/envelope.go`.

## 4. Prevention
Enforce US English spelling (`Canceled`), verify `cli/apperror` constructor signatures (`NewExecutionError`, `NewValidationError`), and run `gofmt -w` on modified Go packages before pushing release commits.

