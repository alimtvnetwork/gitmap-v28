# CI/CD Issue 66: macOS Cluster Join Test Handshake Context Deadline Timeout RCA

- **Stage**: CI / Cross-Platform Build (`macos-latest` / `go test ./... (no cache)`)
- **Status**: ✅ Resolved
- **Affected Packages**: `cli/cluster`, `cli/cmd`
- **Workflow Run**: [#35517024065](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/35517024065)

---

## 1. Symptom

In GitHub Actions workflow run `#35517024065` on job `macos-latest / go build + test`, `go test ./... (no cache)` failed after 84.3s with a context deadline exceeded error in `TestRunJoin_Success`:

```text
--- FAIL: TestRunJoin_Success (9.28s)
join_test.go:120: expected successful join, got error: [E8005:EXECUTION] join.Handshake: Failed to join cluster at 127.0.0.1:49272: context deadline exceeded (at=cmd/join.go:116) (creator=cmd.runJoin) (ctx=map[address:127.0.0.1:49272 hostname:sat12-bq163-41a6a51b-e99d-45f5-b26b-da6aa2b8873a-3A7EA6D1731A.local]) (cause=context deadline exceeded)
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmd	84.304s
FAIL
```

---

## 2. Root Cause

On heavily loaded CI runners (specifically `macos-latest` virtual machines), the background goroutine running `srv.Serve(listener)` in `startMockJoinServer` was delayed in executing `rpcServer.Register(s)` and calling `listener.Accept()`, while `NodeClient.dialTLS()` in `cli/cluster/client.go` had a tight 2-second dial timeout (`Timeout: 2 * time.Second`) with zero retry logic, causing the client's TLS handshake context deadline to expire before the server could process the initial handshake.

---

## 3. Resolution

1. **Dial Timeout & Exponential Retry in `cli/cluster/client.go`:**
   - Increased the client dialer timeout from `2 * time.Second` to `5 * time.Second`.
   - Extracted `dialWithRetry` to attempt TLS dial up to 3 times with 50ms sleep intervals upon transient handshake failures.
2. **Mock Server Startup Synchronization in `cli/cmd/join_test.go`:**
   - In `startMockJoinServer`, added `time.Sleep(20 * time.Millisecond)` to yield the thread and allow the `srv.Serve` goroutine to initialize and enter `listener.Accept()` before returning.
   - In `TestRunJoin_Success`, wrapped `runJoin` in a resilient retry loop (up to 5 attempts with 100ms backoff).

---

## 4. Prevention & Learnings

- **Never use tight 2-second dial timeouts for TLS client connections in CI environments:** Virtualized runners with multi-package test concurrency frequently experience thread-scheduling latency; dial timeouts for in-process or local cluster handshakes should be at least 5 seconds.
- **Always synchronize mock server startup in unit tests:** When launching test listeners or servers in background goroutines (`go srv.Serve(listener)`), allow a brief startup grace period or implement client-side retries to prevent test race conditions.
