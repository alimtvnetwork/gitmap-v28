# RCA-104: Root Cause Analysis - Unit Test Assertion Failure on Delegated Python Command in `cmdssh`

## 1. Executive Summary
- **Task Identification**: `ci-cd-fix-gitmap-release`
- **Current Version**: `v6.490.0` (Commit: `0ac46658`)
- **Failing CI Workflow**: `Cross-Platform Build` (Run `#37323808237`)
- **Impacted Subsystems**:
  * SSH execution command validation in `cli/cmdssh/ssh_exec_command_test.go`

---

## 2. Part 1: Symptom & Environmental Manifestation
- **Workflow Run**: `Cross-Platform Build #37323808237`
- **Jobs Failing**:
  1. `ubuntu-latest / go build + test` (Step: `go test ./... (cached)`)
  2. `macos-latest / go build + test` (Step: `go test ./... (cached)`)
  3. `windows-latest / go build + test` (Step: `go test ./... (cached)`)
- **Observed Error Messages**:
  ```text
  --- FAIL: TestIsGitmapCommand (0.00s)
      ssh_exec_command_test.go:95: expected isGitmapCommand("python") to be false
  FAIL
  FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdssh	0.142s
  ```

---

## 3. Part 2: Root Cause Analysis (RCA)
- **Root Cause**:
  In commit `0ac46658` (`Feature: ssh - add py python and cursor delegation to sshexec and update gitmap skills`), `py`, `python`, `cursor`, and `cur` were added as recognized utility commands in `isGitmapUtilityCommand()` (`cli/cmdssh/ssh_exec_command.go:151`).
  Consequently, `isGitmapCommand("python")` now returns `true`.
  However, in `cli/cmdssh/ssh_exec_command_test.go:92`, the test fixture `invalidCommands` still contained `"python"`, causing `TestIsGitmapCommand` to fail when asserting that `isGitmapCommand("python")` must be `false`. Additionally, the newly added valid commands (`py`, `python`, `cursor`, `cur`) had not yet been appended to the test's `validCommands` slice.

---

## 4. Part 3: Surgical Technical Solution
- **File Modified**: `cli/cmdssh/ssh_exec_command_test.go`
- **Refactoring Applied**:
  1. Added `"py"`, `"python"`, `"cursor"`, and `"cur"` to `validCommands` in `TestIsGitmapCommand`.
  2. Replaced `"python"` in `invalidCommands` with `"ruby"`, ensuring truly invalid commands are properly validated without conflicting with the new delegated command set.

---

## 5. Part 4: Verification & Proof Matrix

| Check / Requirement | Tool / Command Executed | Result | Status |
|:---|:---|:---|:---:|
| SSH Exec Command Test | `go test -short ./cmdssh/...` | PASS (`ok github.com/alimtvnetwork/gitmap-v28/cli/cmdssh 50.040s`) | **VERIFIED** |
| Specific Test Unit | `go test ./cmdssh -run TestIsGitmapCommand` | PASS (`0.146s`) | **VERIFIED** |
| Nested If Linter | `python linter-scripts/check-nested-ifs.py` | PASS (0 violations across 2 files) | **VERIFIED** |
| Boolean & Enum Linter | `python linter-scripts/check-enum-and-boolean.py` | PASS (0 violations across 3208 files) | **VERIFIED** |
| Boolean Guidelines Linter | `python linter-scripts/check-boolean-guidelines.py` | PASS (0 violations across 2 files) | **VERIFIED** |
| Error Management Linter | `python linter-scripts/check-error-management.py` | PASS (0 violations across 4269 files) | **VERIFIED** |
| Relative Paths Linter | `python linter-scripts/check-relative-paths.py` | PASS (0 absolute paths across 8777 files) | **VERIFIED** |
| Newline Styling Linter | `python linter-scripts/check-newline-styling.py` | PASS (All files Unix LF) | **VERIFIED** |
