# Subtask 03: Asynchronous Concurrent Multi-Node Dispatch & Execution Timeout

> **Parent Plan:** `67-nodes-cfr-remote-fleet-clone-enhancement.md`  
> **Status:** `PENDING`  
> **Target Files:**  
> - `cli/cmdnodes/nodes_clone.go`  
> - `cli/cmdnodes/nodes_clone_remote.go`  
> - `cli/cmdnodes/nodes_clone_types.go`  
> - `cli/cmdnodes/nodes_clone_test.go`  

---

## 1. Technical Context & Scope

In distributed fleet deployments, running remote commands sequentially causes unacceptable performance bottlenecks. If a single machine is bogged down by slow network disk I/O, large git packfiles, or a frozen SSH channel, sequential processing blocks all subsequent machines. Furthermore, without a hard execution timeout ceiling, the CLI can hang indefinitely.

Currently:
1. `executeFleetNodesParallel` launches unbounded goroutines but lacks an execution context with an explicit timeout deadline (e.g. 3 minutes).
2. Remote commands do not uniformly enforce structured `--json` execution and response deserialization.
3. Errors from hung network sockets leak into OS-level timeouts (often several minutes) instead of being bounded deterministically.

This subtask updates the dispatch pipeline in `cli/cmdnodes/nodes_clone.go` and `cli/cmdnodes/nodes_clone_remote.go` to provide:
- Concurrent asynchronous fan-out across all online nodes.
- A 3-minute hard execution timeout ceiling managed via `context.WithTimeout`.
- Structured remote command invocation passing `--json` to `gitmap cfr` / `gitmap clone`.
- Structured JSON response extraction into `RemoteCloneNodeResult` with seamless legacy fallback.

---

## 2. Technical Specification & Implementation Details

### 2.1 Asynchronous Worker Dispatch Engine

- **Function:** `executeFleetNodesParallel(conns []db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string, timeout time.Duration) []RemoteCloneNodeResult`
- **Concurrency Architecture:**
  - Create root execution context: `ctx, cancel := context.WithTimeout(context.Background(), timeout)`. Default timeout is `3 * time.Minute`.
  - Filter target connections to only those confirmed reachable/online during pre-flight probing (or all candidate connections if pre-flight is skipped).
  - Launch each node worker in an isolated goroutine.
  - Synchronize completion with `sync.WaitGroup`.
  - Collect individual node outcomes into a slice protected by `sync.Mutex`.

### 2.2 Context-Aware Remote Execution & Timeout Bounds

- **Function:** `runRemoteNodeWorkerWithContext(ctx context.Context, conn db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string) RemoteCloneNodeResult`
- **Lifecycle & Cancellation:**
  - Establish SSH connection using `cmdssh.ConnectSSHClientWithContext(ctx, conn)` or connect with cancel monitoring.
  - If `ctx.Done()` fires prior to or during execution:
    - Mark node status as `"failed"`.
    - Record error: `"execution timed out (3m limit exceeded)"`.
    - Close client connection immediately to terminate hung remote child processes.
  - Record execution duration in milliseconds (`DurationMs`).

### 2.3 Remote JSON Protocol Execution & Parsing

- **Command Generation:**
  - Format command: `gitmap cfr <args> --json` (or `gitmap clone <args> --json`).
  - Target directory is dynamically set via `Set-Location` (PowerShell) or `cd` (Unix).
- **Structured Response Extraction:**
  - Capture remote stdout.
  - Invoke `cmdclone.ParseCloneJSONResponse(stdout)`.
  - On valid JSON payload (`DirectCloneJSONResponse`):
    - `Success == true`: set `res.Status = "success"`, `res.Details = payload.Message`.
    - `Success == false`: set `res.Status = "failed"`, `res.Error = payload.Message`.
- **Fallback for Legacy Remote Binaries:**
  - If remote stderr or stdout contains `"flag provided but not defined: --json"`, automatically invoke `retryWithoutJSON(client, opts, fileName, isWin, shell)`.
  - Parse legacy string output (`"already exists on disk"`, `"cloned successfully"`) to populate `res.Details`.

### 2.4 Local Host Coordination

- **Function:** `executeLocalCloneWithContext(ctx context.Context, opts NodesCloneOptions) (bool, string, time.Duration)`
- **Behavior:**
  - If `opts.IsSkipLocal` is true, immediately return `(true, "skipped local execution (except-self)", 0)`.
  - Pass `--json` flag to `dispatchLocalKind` and extract outcome via `cmdclone.ParseCloneJSONResponse`.
  - Respect `ctx` cancellation if local execution exceeds the deadline.

---

## 3. Function Decomposition Plan (<= 15 Lines Per Function)

To adhere to the <= 15 lines per function coding guideline, the implementation is decomposed into:

1. `executeFleetNodesParallel(...) []RemoteCloneNodeResult` (<= 14 lines)
2. `createFleetExecutionContext(timeout time.Duration) (context.Context, context.CancelFunc)` (<= 6 lines)
3. `spawnRemoteWorkerGoroutine(ctx context.Context, wg *sync.WaitGroup, mu *sync.Mutex, conn db.SSHConnection, opts NodesCloneOptions, b []byte, fn string, out *[]RemoteCloneNodeResult)` (<= 14 lines)
4. `runRemoteNodeWorkerWithContext(...) RemoteCloneNodeResult` (<= 14 lines)
5. `connectRemoteSSHWithTimeout(ctx context.Context, conn db.SSHConnection) (*ssh.Client, error)` (<= 12 lines)
6. `executeRemoteCommandWithContext(ctx context.Context, client *ssh.Client, cmd, shell string) (string, error)` (<= 14 lines)
7. `parseRemoteCloneJSONOutcome(res *RemoteCloneNodeResult, stdout string) bool` (<= 12 lines)
8. `handleRemoteExecutionError(res *RemoteCloneNodeResult, err error, stdout string)` (<= 12 lines)
9. `isExecutionTimeoutError(err error, ctx context.Context) bool` (<= 8 lines)
10. `executeLocalCloneWithContext(...) (bool, string, time.Duration)` (<= 14 lines)

---

## 4. Coding Guidelines & Invariant Rules

- **Positive Booleans Only:** Use `isLocalOk`, `isTimedOut`, `hasJSONResponse`, `isAsync`. Avoid negative boolean variables.
- **Error Wrapping:** Wrap system errors with `apperror.WrapSimple(err, "<context>")`.
- **Safe Resource Cleanup:** Always ensure `defer client.Close()` and `defer cancel()` are called.
- **Zero-Git Command Rule:** Do not invoke any git commands during development or test execution.

---

## 5. Verification & Testing Protocol

1. **Unit Tests for Concurrency & Timeout (`cli/cmdnodes/nodes_clone_test.go`):**
   - Test `executeFleetNodesParallel` launches concurrent workers and aggregates results.
   - Test timeout cancellation fires when context deadline expires before worker completion.
   - Test `parseRemoteCloneJSONOutcome` correctly parses structured success and failure payloads.
   - Test legacy fallback triggers when `--json` flag is reported unrecognized.
2. **Quality Gate:**
   - Confirm all functions are 15 lines or fewer.
   - Verify code compiles cleanly with zero lint regressions.
