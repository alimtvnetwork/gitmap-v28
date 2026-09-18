# 4-Part Root Cause Analysis (RCA): Exhaustive Switch ReplaceModeType and macOS osuser Test Failures

## 1. Why It Happened (High-Level Business & Architectural Context)
In pipeline run #35244186467, two quality gates failed:
1. The **Exhaustive Switch Diff** gate in `CI.yml` detected 1 new finding in `cli/cmd/replace.go`: `missing cases in switch of type cmd.ReplaceModeType: cmd.ReplaceModeTypeUnknown`. The repository enforces that all enum constants must have explicit cases in `switch` blocks.
2. The **macOS Cross-Platform Build** gate in `crossplatform-build.yml` failed in `cli/osuser` during `go test ./...` on `macos-latest`: `unexpected error: [E9000:EXECUTION] execution: unsupported os for user creation: darwin`. The `osuser` module intentionally restricts user creation, removal, and termination to Linux and Windows, but unit tests did not check `runtime.GOOS` before executing.

---

## 2. How It Happened (Technical Execution Flow)
1. **Exhaustive Linter Check**:
   - In `cli/cmd/replace_classify.go`, `ReplaceModeType` is defined with 5 constants: `ReplaceModeTypeUnknown`, `ReplaceModeTypeLiteral`, `ReplaceModeTypeVersionN`, `ReplaceModeTypeAll`, `ReplaceModeTypeAudit`.
   - In `cli/cmd/replace.go`, `dispatchReplaceMode` had:
     ```go
     switch mode {
     case ReplaceModeTypeLiteral: ...
     case ReplaceModeTypeAudit: ...
     case ReplaceModeTypeAll, ReplaceModeTypeVersionN: ...
     default: ...
     }
     ```
   - Because `ReplaceModeTypeUnknown` was omitted from the explicit `case` clauses, `golangci-lint` analyzer `exhaustive` flagged it as a new unhandled enum case.
2. **macOS Cross-Platform Runner Execution**:
   - GitHub Actions executes `go test ./...` on `macos-latest` where `runtime.GOOS == "darwin"`.
   - `CreateRootUser`, `RemoveEnhancedUser`, and `KillUserProcesses` delegate to `dispatchPlatformCreate`, `dispatchPlatformRemove`, and `dispatchPlatformKill`.
   - In all three functions, `switch runtime.GOOS` only handles `"windows"` and `"linux"`. Any other OS hits `default: return apperror.NewExecutionError(...)`.
   - In `cli/osuser/osuser_test.go`, the mocked unit tests (`TestCreateRootUser_Mocked`, `TestRemoveEnhancedUser_Mocked`, and `TestKillUserProcesses_Mocked`) expected the operation to succeed with mock runner output, but failed at the OS dispatch switch with `unsupported os: darwin`.

---

## 3. Root Cause (Exact File, Line, and Mechanism)
- `cli/cmd/replace.go:31`: The `switch mode` block lacked an explicit `case ReplaceModeTypeUnknown:` clause.
- `cli/osuser/osuser_test.go:101, 114, 127`: The mocked tests lacked a platform gate (`runtime.GOOS == "windows" || runtime.GOOS == "linux"`) and failed when executed on macOS CI runners.

---

## 4. Code Fix (Exact Code Snippets)

### Fix 1: Add Explicit `ReplaceModeTypeUnknown` Case in `cli/cmd/replace.go`
```go
func dispatchReplaceMode(mode ReplaceModeType, positional []string, opts replaceOpts) {
	switch mode {
	case ReplaceModeTypeLiteral:
		runReplaceLiteral(positional[0], positional[1], opts)
	case ReplaceModeTypeAudit:
		runReplaceAudit(opts)
	case ReplaceModeTypeAll, ReplaceModeTypeVersionN:
		dispatchVersionMode(mode, positional, opts)
	case ReplaceModeTypeUnknown:
		fmt.Fprint(os.Stderr, constants.ErrReplaceNeedsArgs)
		cliexit.HandleError(nil, constants.ExitCodeError)
	default:
		fmt.Fprint(os.Stderr, constants.ErrReplaceNeedsArgs)
		cliexit.HandleError(nil, constants.ExitCodeError)
	}
}
```

### Fix 2: Add Platform Skip Guard in `cli/osuser/osuser_test.go`
```go
func isSupportedUserPlatform() bool {
	return runtime.GOOS == "windows" || runtime.GOOS == "linux"
}

func TestCreateRootUser_Mocked(t *testing.T) {
	if !isSupportedUserPlatform() {
		t.Skipf("skipping on unsupported platform: %s", runtime.GOOS)
	}
...
func TestRemoveEnhancedUser_Mocked(t *testing.T) {
	if !isSupportedUserPlatform() {
		t.Skipf("skipping on unsupported platform: %s", runtime.GOOS)
	}
...
func TestKillUserProcesses_Mocked(t *testing.T) {
	if !isSupportedUserPlatform() {
		t.Skipf("skipping on unsupported platform: %s", runtime.GOOS)
	}
```
