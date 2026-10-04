# Issue 85 RCA: SSH Authentication Handshake Fallback and Accurate Node Status Reporting

## 1. Reproduction & Symptoms
1. **Premature Authentication Failure on Fleet Update (`gitmap agm update --ssh` / `gitmap update --ssh`):**
   - When executing `gitmap agm update --ssh` across live cluster nodes (`w1` at `node-w1`, `w2` at `node-w2`, `w3` at `node-w3`), nodes failed with:
     ```text
     [w2|node-w2] Connect error: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
     [w1|node-w1] Connect error: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
     [w3|node-w3] Connect error: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
     ```
   - If a node's configured key file (`KeyPath`) was rejected by the remote OpenSSH server (e.g. key missing from server's `administrators_authorized_keys`), the client immediately logged `Connect error` and failed before falling back to decrypted password authentication.
   - Offline nodes (`alpha-win`, `beta-linux`, `gamma-mac` at `10.20.0.11-13`) were misreported as `FAILED: unable to connect` in the fleet summary table instead of `○ OFFLINE`.

2. **Inaccurate Node Status Reporting (`gitmap ssh nodes`):**
   - Running `gitmap ssh nodes` displayed `STATUS: ● ready` for all registered nodes, even offline/unreachable machines (`alpha-win`, `beta-linux`, `gamma-mac`), because `renderHostTableRow` hardcoded `defaultHostStatus = "ready"` and never performed liveness or reachability checks.
   - Running `gitmap ssh nodes -v` displayed `STATUS: ● ready` even when connection failed due to authentication rejection (`GITMAP VERSION: not installed`).

## 2. Root Cause Analysis
1. **Short-Circuiting and Premature Error Logging in SSH Client Connection:**
   - In `cli/cmdssh/sshexec.go`, `connectSSHClient` called `tryConnectWithKey` before fallback password resolution.
   - When `crypto.ConnectWithKey` failed, `connectWithKeyPath` immediately printed `printHeaderError(header, "Connect error", err)` to standard output.
   - Furthermore, `c.EncryptedPassword` was unpopulated for nodes enrolled via key only, and `queryHostPasswordFromDB` only queried the SQLite database without consulting default cluster credentials (`vmpass.json`).
2. **Missing Live Status Probing in Table Listing:**
   - `printSJList` in `cli/cmdssh/sshjoin_ls_cmd.go` fetched hosts from SQLite and directly rendered the table via `RenderSSHHostsTable`.
   - `renderHostTableRow` in `cli/cmdssh/sshjoin_table.go` ignored actual node availability, unconditionally rendering `formatStatusColored(defaultHostStatus, 10)` with green `● ready`.
3. **Improper Offline Classification in Fleet Parallel Worker:**
   - In `cli/cmdssh/ssh_update_remote.go` and `cli/cmdssh/fleet_parallel_render.go`, when `isNodeAvailable` classified a host as offline, the worker returned an unqualified error, causing `PrintFleetDone` to print `[FLEET FAIL]` and `PrintFleetSummary` to report the host under `Failed` rather than `Offline` / `[FLEET SKIP]`.

## 3. Code Fix
1. **Fallback SSH Authentication (`cli/crypto/ssh_client.go`, `cli/cmdssh/ssh_dial_fallback.go`, `cli/cmdssh/ssh_vmpass.go`):**
   - Added `crypto.ConnectWithFallback(ip, user, keyPath, password string) (*ssh.Client, error)` attempting key authentication first and falling back to password authentication.
   - Created `ssh_vmpass.go` with `ResolveFallbackCredentials(user, osType string) string` resolving fallback credentials from `vmpass.json` across workspace, binary, and home directories.
   - Created `ssh_dial_fallback.go` with `dialNodeWithFallback(c db.SSHConnection, header string) (*ssh.Client, error)` that tries configured key, default user keys (`id_ed25519`, `id_rsa`), and vaulted/fallback passwords without premature error logging. When password authentication succeeds, it automatically encrypts and vaults the password into SQLite for subsequent connections.
   - Updated `connectSSHClient` in `cli/cmdssh/sshexec.go` and `dialSSHNodeClient` in `cli/cmdagy/agy_running_projects_ssh.go` to use fallback dialing.
2. **Live Node Status Probing (`cli/store/models_ssh.go`, `cli/cmdssh/sshjoin_probe.go`, `cli/cmdssh/sshjoin_table.go`, `cli/cmdssh/sshjoin_ls_cmd.go`):**
   - Added `Status string` field to `store.SSHHost`.
   - Created `sshjoin_probe.go` with `probeHostsLiveStatus` concurrently probing TCP connectivity (1200ms timeout) and SSH authentication.
   - Updated `formatStatusColored` and `renderHostsTableHeader` in `sshjoin_table.go` to support column width 21, correctly displaying `● ready` (green), `○ offline (timeout)` (dim), and `▲ auth failed` (yellow).
   - Updated `probeOnlineNodeVersion` and `resolveNodeStatusColor` in `cli/cmdssh/ssh_node_version.go` to set `▲ auth failed` when connection fails on online nodes.
3. **Offline Node Handling in Fleet Parallel Execution (`cli/cmdssh/fleet_parallel_render.go`, `cli/cmdssh/ssh_update_remote.go`):**
   - Updated `executeSingleSSHNodeUpdate` to return explicit offline error when `isNodeAvailable` fails.
   - Updated `PrintFleetDone` and `PrintFleetSummary` in `fleet_parallel_render.go` to classify unreachable nodes as `[FLEET SKIP] ...: OFFLINE` and `○ OFFLINE`, reporting separate `Offline: N` metrics.

## 4. Verification & Prevention
1. **E2E & Unit Test Verification:**
   - `go test -v -tags=tempe2e -run TestSSHFleetAuthFallback_TempE2E ./tests/e2e/...`: PASS (0.41s across `w1`, `w2`, `w3`).
   - `go test -v -run TestRenderSSHHostsTable ./cmdssh/...`: PASS (100% across all table formats including offline and auth failed).
   - `go test -v ./crypto/...`: PASS (100%).
2. **Live Execution Verification:**
   - `gitmap ssh nodes`: Accurately displays `○ offline (timeout)` for `alpha-win`, `beta-linux`, `gamma-mac`, and `● ready` for live nodes `w1`, `w2`, `w3`.
   - `gitmap ssh nodes -v`: Accurately queries and displays installed GitMap versions for `w1`, `w2`, `w3` with zero connection errors, and `○ offline` for offline machines.
   - `gitmap agm update --ssh`: Succeeded across all 3 live nodes (`w1`, `w2`, `w3`) with 0 failures and 3 offline nodes cleanly skipped.
   - `gitmap ssh pass ls`: Accurately shows all 3 live nodes vaulted with `Encrypted (RSA/AES)` passwords.
3. **Prevention:**
   - Always run pre-flight connectivity checks and comprehensive multi-credential fallbacks for fleet SSH commands.
   - Never hardcode static status indicators in terminal listing tables.
