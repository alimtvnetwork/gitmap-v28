# Subtask 06: OS User Command Runner Decoupling & Hermetic Mock Testing

## Objective
Decouple `cli/osuser` from unmocked host user and process modification commands (`net user`, `net localgroup`, `useradd`, `deluser`, `userdel`, `chpasswd`, `pkill`, `taskkill`) by introducing an injectable `OSCommandRunner`.

## Assigned Files
- `cli/osuser/runner.go`
- `cli/osuser/create_root.go`
- `cli/osuser/remove.go`
- `cli/osuser/kill.go`
- `cli/osuser/sudoers.go`
- `cli/osuser/osuser_test.go`

## Implementation Steps
1. In `cli/osuser/runner.go`:
   - Introduce `type OSCommandRunner func(cmd *exec.Cmd) ([]byte, error)`.
   - Define `var defaultOSCommandRunner OSCommandRunner = func(cmd *exec.Cmd) ([]byte, error) { return cmd.CombinedOutput() }`.
2. In `cli/osuser/create_root.go`, `remove.go`, `kill.go`, `sudoers.go`:
   - Route all `cmd.CombinedOutput()` and `cmd.Run()` calls through `defaultOSCommandRunner(cmd)`.
3. In `cli/osuser/osuser_test.go`:
   - Provide `setupMockRunner()` with `defer` restoration.
   - Add unit tests for `CreateRootUser`, `RemoveEnhancedUser`, and `KillUserProcesses`.
