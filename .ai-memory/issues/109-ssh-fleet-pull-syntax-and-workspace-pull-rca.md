# Root Cause Analysis: SSH Fleet Pull Shell Discrepancy & Workspace-Wide Pull Isolation

- **Issue ID:** RCA-109
- **Date:** 2026-10-23
- **Severity:** High
- **Status:** Resolved

## 1. Problem Statement
When executing `gitmap pa --ssh`, remote node `w1` failed with `(invalid JSON: The filename, directory name, or volume label syntax is incorrect.)`. Furthermore, the SSH fleet enqueue was sluggish due to sequential probing timeouts on offline nodes, and `gitmap pull` without parameters was locked to the running directory instead of operating across all registered workspaces.

## 2. Root Cause Analysis
1. `runRemotePullJSON` executed `secrets.RunCommand(client, remoteCmd, "")` with an empty `shellType`. Because OpenSSH on Windows defaults to `cmd.exe`, PowerShell environment variables (`$env:GCM_CREDENTIAL_STORE='none'; ...`) caused `cmd.exe` to fail with "The filename, directory name, or volume label syntax is incorrect."
2. `probeHostHealth` enforced an unconditional 1.5s retry dial on failed TCP probes, meaning offline nodes took 4.5s each to identify. In addition, `queryFleetNodeVersion` probed REST port 49152 (1.5s per endpoint) before reading `c.BuildVersion` from the database.
3. `isPullCWDEnabled` checked `isGitRepoCWD()`, restricting default `gitmap pull` to the current repository instead of querying all registered workspaces.

## 3. Remediation
1. Updated `runRemotePullJSON` to pass `shellType := resolveRemoteShell(osType)`, wrapping Windows commands in `powershell -NoProfile -Command "..."` and Linux commands in `bash -c "..."`.
2. Parameterized `buildRemoteFleetPullCmd(osType, concurrencyFlags)` to use platform-native environment setting syntax and increased default concurrency to `--parallel 4`.
3. Bounded fleet liveness TCP timeout to `500ms`, skipping retries on fast health checks, and prioritizing stored `c.BuildVersion` in `queryFleetNodeVersion`.
4. Refactored `isPullCWDEnabled` to only match explicit CWD targets (`opts.slug == "." || opts.slug == "cwd"`), defaulting unargumented `gitmap pull` to `opts.all = true` across all workspaces.
