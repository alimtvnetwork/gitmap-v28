# CI/CD Issue 05: Cluster TLS Dial Timeout & VSCode PM Sync Test Env Race

- **Stage**: CI `go test -race` / Test Matrix
- **Status**: ✅ Resolved
- **Affected Packages**: `github.com/alimtvnetwork/gitmap-v28/gitmap/cluster`, `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd`

## 1. Problem Description

During CI test runs on Linux runners (`go test -race -count=1 -timeout=15m ./cmd/... ...`), `github.com/alimtvnetwork/gitmap-v28/gitmap/cmd` took over 151 seconds and failed:
```
→ vscode-pm-sync: re-tagging projects.json entries at /tmp/TestVSCodePMSyncDryRunDoesNotMutate1643790492/001/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json
  • dry-run: 1 entries scanned, 1 would change (no write)
  Note: --no-commit set; leaving guideline files uncommitted.
  Note: --no-push set (or no upstream); push step skipped.
FAIL
FAIL	github.com/alimtvnetwork/gitmap-v28/gitmap/cmd	151.660s
```

## 2. Root Cause Analysis

Two interacting factors caused this failure:

1. **Unbounded `tls.Dial` Connection Hanging in `cluster/dispatcher.go`**:
   - `cluster.Dispatch()` invoked `tls.Dial("tcp", node.IP+":8081", conf)` without a connection dialer timeout.
   - When executing `TestClusterCommand`, `node-2` has a dummy unreachable IP (`192.168.1.11`). On Linux CI runners, the TCP SYN retransmission timeout blocked `tls.Dial` for 130-150 seconds per remote node attempt before timing out at the OS socket level.

2. **Unsynchronized Environment Mutation in VSCode PM Sync Test Helpers**:
   - `swapHomeEnv` in `vscodepmsync_testhelper_test.go` used raw `os.Setenv` to mutate `HOME`, `XDG_CONFIG_HOME`, and `APPDATA`.
   - When tests execute concurrently under `go test -race`, concurrent `os.Setenv` / `os.Getenv` access triggers a data race in the Go runtime.

## 3. Solution

1. **Bound TLS Dialing with `net.Dialer` Timeout**:
   - In `gitmap/cluster/dispatcher.go`, replaced bare `tls.Dial` with `tls.DialWithDialer(&net.Dialer{Timeout: 500 * time.Millisecond}, "tcp", node.IP+":8081", conf)` for both PowerShell and Cmd executors.
   - In `gitmap/cluster/client.go`, replaced bare `tls.Dial` in `dialTLS()` with `tls.DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "tcp", c.address, conf)`.
   - `TestClusterCommand` execution time dropped from 21.07s / 151s down to 2.08s.

2. **Adopt Thread-Safe `t.Setenv`**:
   - Refactored `swapHomeEnv` in `vscodepmsync_testhelper_test.go` to accept `t *testing.T` and use `t.Setenv()`.

## 4. What NOT to Repeat

- **Never call `net.Dial` or `tls.Dial` without a timeout**: Always use `tls.DialWithDialer` with an explicit `Dialer.Timeout` or a context (`DialContext`) to prevent OS-level TCP SYN hangs on unreachable IPs.
- **Never mutate environment variables with `os.Setenv` inside tests**: Always use `t.Setenv(key, val)`, which integrates with the Go test harness and guarantees race-free cleanup.
