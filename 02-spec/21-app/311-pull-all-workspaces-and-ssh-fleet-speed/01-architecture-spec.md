# Plan 311 Architecture Specification: Workspace-Wide Default Pull & Accelerated SSH Fleet Remote Execution

## User Request (Verbatim)
```text
PS C:\Users\Administrator> gitmap pa --ssh

  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [w3] (192.168.1.12): Offline (skipped, no task enqueued)
    • Remote Node [w1] (192.168.1.3) [GitMap v6.589.0]: Online → Enqueued (async)
    • Remote Node [w2] (192.168.1.7): Local machine (covered by local execution, skipped)
    • Remote Node [w4] (192.168.1.13): Offline (skipped, no task enqueued)
    • Remote Node [u1] (192.168.1.22): Offline (skipped, no task enqueued)
    • Remote Node [main] (192.168.1.20): Offline (skipped, no task enqueued)
    • Current Machine [AI-Main (127.0.0.1)] [GitMap v6.590.0]: Running locally (direct execution, not enqueued)

  • Local VM (127.0.0.1 - localhost): 68 pulled (42 active, 26 up-to-date)
    ...
  ▶ Node [w1] (IP: 192.168.1.3): failed
      (invalid JSON: The filename, directory name, or volume label syntax is incorrect.)

gitmap pa: execute failed: [E9000:EXECUTION] SSH fleet pull-all failed: [E_SSH_FLEET_FAIL:EXECUTION] SSH fleet pull-all had 1 node failure(s): (at=cmdssh/ssh_pull_fleet.go:63) (at=cmdpull/pull_dispatch.go:102)

PS D:\> gitmap pa
  → gitmap pull-all (cwd: D:\)
PS D:\> gitmap cd $work
D:\work
PS D:\work> gitmap pa
  → gitmap pull-all (cwd: D:\work)

FIX ALL these bugs, 
gitmap pull will run on all workspaces not on the running dir, fix the code and fix the bugs 
and remote `gitmap pa --ssh ` is super slow now figure out the root cause and fix it and release please
```

## 1. 4-Part Root Cause Analysis (RCA)

### 1.1 Failure & Symptoms
1. **Remote Node [w1] Failure:**
   `gitmap pa --ssh` failed with:
   `▶ Node [w1] (IP: 192.168.1.3): failed (invalid JSON: The filename, directory name, or volume label syntax is incorrect.)`
2. **Slow Fleet Enqueue & Probing:**
   `gitmap pa --ssh` took multiple seconds before enqueuing because each offline node spent 4.5 seconds in retry timeouts, and online nodes spent 3.0 seconds querying non-existent REST endpoints on port 49152.
3. **Running-Dir Bound Pull:**
   `gitmap pull` without arguments only pulled CWD when inside a git repository, or searched for sub-repositories under CWD when outside a repository, instead of running across all registered workspaces.

### 1.2 Root Cause Analysis
1. **Windows OpenSSH Shell Discrepancy:**
   In `cli/cmdssh/ssh_pull_fleet_dispatch.go`, `runRemotePullJSON` invoked `secrets.RunCommand(client, remoteCmd, "")` with an empty string for `shellType`. Because `shellType == ""`, `wrapCommandForShell` returned the command without wrapping. On Windows OpenSSH, the default execution shell is `cmd.exe`. The command string contained PowerShell environment syntax (`$env:GCM_CREDENTIAL_STORE='none'; ...`). `cmd.exe` failed with `The filename, directory name, or volume label syntax is incorrect.`, causing JSON parsing to fail.
2. **Probing Bottlenecks:**
   - In `ssh_health.go`, `probeWithRetry` executed a 3.0s dial followed by an unconditional 1.5s retry for every offline host.
   - In `ssh_pull_fleet_liveness.go`, `queryFleetNodeVersion` prioritized REST endpoints on port 49152 (1.5s timeout per endpoint) before falling back to `c.BuildVersion`, which was already stored in the database.
   - In `ssh_pull_fleet_dispatch.go`, `executeRemoteNodePull` re-probed OS type and version via SSH even though `c.OS` and `c.BuildVersion` were already recorded.
3. **Pull CWD Coupling:**
   In `cli/cmdpull/pull_cwd.go`, `isPullCWDEnabled` returned `isGitRepoCWD()`, binding unargumented `gitmap pull` to the current working directory instead of pulling all workspaces.

### 1.3 Corrective Actions & Architectural Invariants
1. **Shell-Aware Remote Execution:**
   - Resolve `osType` (`c.OS` fast-path) and `shellType := resolveRemoteShell(osType)`.
   - On Windows: format command for PowerShell (`powershell -NoProfile -Command "..."`) via `shellType: "ps"`.
   - On Linux/macOS: format command for POSIX shell (`export GCM_... && gitmap pa --json`) via `shellType: "bash"`.
2. **Accelerated Fleet Probing:**
   - Bound TCP health check timeout to `500ms` for fleet liveness, skipping retries on fast probes (<1000ms).
   - In `queryFleetNodeVersion`, check `c.BuildVersion` first to resolve in 0ms.
   - Increase remote concurrency to `--parallel 4`.
3. **Workspace-Wide Default Pull:**
   - Restrict `isPullCWDEnabled` to explicit current-directory targets (`opts.slug == "." || opts.slug == "cwd"`).
   - Default unargumented `gitmap pull` to `opts.all = true`, enqueuing all tracked repositories across all registered workspaces.
   - Update invocation banner to display `(all workspaces)`.

### 1.4 Verification & Release Invariant
- Verify via `python 03-ai-scripts/50-fastgate.py` and `python linter-scripts/check-relative-paths.py`.
- Bump version from `6.590.0` to `6.591.0` via `python 03-ai-scripts/37-bump-version.py -t minor`.
