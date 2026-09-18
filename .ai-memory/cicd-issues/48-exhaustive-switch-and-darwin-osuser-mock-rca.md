# CI/CD Issue 48: Exhaustive Switch and macOS osuser Mock Failures RCA

## 1. Why It Happened
CI run #35244186467 and Cross-Platform Build run #35244186285 failed due to:
1. `exhaustive` linter finding missing enum constant `ReplaceModeTypeUnknown` in `cmd/replace.go:31`.
2. `macos-latest` failing 3 unit tests in `cli/osuser` because `CreateRootUser`, `RemoveEnhancedUser`, and `KillUserProcesses` return `unsupported os: darwin`.

---

## 2. How It Happened
- `ReplaceModeType` enum in `replace_classify.go` includes `ReplaceModeTypeUnknown`, which was not explicitly covered in `cmd/replace.go`'s `switch mode`.
- `osuser` tests `TestCreateRootUser_Mocked`, `TestRemoveEnhancedUser_Mocked`, and `TestKillUserProcesses_Mocked` mock the command runner, but the OS dispatch functions return an `AppError` on macOS before reaching the runner.

---

## 3. Root Cause
- `cli/cmd/replace.go:31`: Missing `case ReplaceModeTypeUnknown:` in `switch mode`.
- `cli/osuser/osuser_test.go:101, 114, 127`: Missing platform guard `if !isSupportedUserPlatform() { t.Skipf(...) }`.

---

## 4. Code Fix
- Added explicit `case ReplaceModeTypeUnknown:` in `cli/cmd/replace.go`.
- Added `isSupportedUserPlatform()` and `t.Skipf(...)` to mocked tests in `cli/osuser/osuser_test.go`.
