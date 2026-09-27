# Issue 48 RCA: SSH Fleet Authentication Handshake Failure and Misleading Node Status

## 1. Reproduction

1. **SSH Authentication Handshake Abort:**
   - Execute `gitmap agm update --ssh` across the cluster fleet.
   - Nodes `w1` (192.168.1.3), `w2` (192.168.1.7), and `w3` (192.168.1.12) are active and running.
   - Output immediately fails with connection errors:
     ```text
     [w2|192.168.1.7] Connect error: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
     [w1|192.168.1.3] Connect error: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
     [w3|192.168.1.12] Connect error: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain
     [FLEET FAIL]  [w2|192.168.1.7]: FAILED - unable to connect to [w2|192.168.1.7] (took 104ms)
     [FLEET FAIL]  [w3|192.168.1.12]: FAILED - unable to connect to [w3|192.168.1.12] (took 122ms)
     [FLEET FAIL]  [w1|192.168.1.3]: FAILED - unable to connect to [w1|192.168.1.3] (took 120ms)
     ```
   - Total execution reports 0 succeeded and 6 failed.

2. **Misleading `● ready` Status on Offline Nodes:**
   - Run `gitmap ssh nodes`.
   - The registered fleet contains unreachable machines (`alpha-win`, `beta-linux`, `gamma-mac` on `10.20.0.11-13`).
   - The output table reported:
     ```text
       ALIAS            ROLE           HOST (IP:PORT)         USER           STATUS     ENROLLED
       ----------------------------------------------------------------------------------------------------
       alpha-win        worker         10.20.0.11:22          admin          ● ready    2026-09-26 02:24:54
       beta-linux       worker         10.20.0.12:22          ubuntu         ● ready    2026-09-26 02:24:54
       gamma-mac        worker         10.20.0.13:22          devops         ● ready    2026-09-26 02:24:54
     ```
   - Offline nodes were misleadingly labeled as `● ready` (green), giving the false impression that they were reachable.

## 2. Root Cause Analysis

1. **Absence of Password & Vault Fallback in SSH Dialing:**
   - In `cli/cmdssh/sshexec.go`, `connectSSHClient` strictly invoked `connectWithKeyPath`. If the configured private key failed or was not accepted by the remote host, it immediately printed `Connect error: ssh: handshake failed: ...` to standard output and returned failure without attempting fallback authentication methods.
   - The system did not check default key files (`~/.ssh/id_ed25519`, `~/.ssh/id_rsa`), did not inspect SQLite vaulted passwords (`encrypted_password` in `ssh_hosts`), and did not fall back to `vmpass.json` credential stores.

2. **Windows Server 2022 OpenSSH Key Authentication Isolation:**
   - On Windows Server nodes (`administrator@192.168.1.7` / `w2`), OpenSSH server runs under the Windows security subsystem. For accounts in the `Administrators` local group, OpenSSH strictly ignores `%USERPROFILE%\.ssh\authorized_keys`.
   - Instead, Windows OpenSSH requires keys to be placed in `C:\ProgramData\ssh\administrators_authorized_keys` with strict ACL permissions (`icacls` restricted to `SYSTEM` and `Administrators`). When public key authentication failed against `w2`, the absence of password fallback completely blocked SSH commands.

3. **Static Default Status in `RenderSSHHostsTable`:**
   - In `cli/cmdssh/sshjoin_table.go`, `renderHostTableRow` hardcoded `status := defaultHostStatus` (where `defaultHostStatus = "ready"`).
   - `printSJList` in `cli/cmdssh/sshjoin_ls_cmd.go` fetched rows from SQLite and rendered the table immediately without conducting a lightweight live TCP connectivity probe. Consequently, dead IP addresses were always displayed with the green bullet `● ready`.

4. **Premature Error Logging Before Fallback Exhaustion:**
   - `connectWithKeyPath` called `printHeaderError(header, "Connect error", err)` immediately upon receiving an authentication failure from the SSH library, polluting standard output even when a secondary password fallback would succeed.

## 3. Corrective Implementation

1. **Multi-Tier SSH Client Fallback (`cli/crypto/ssh_client.go`, `cli/cmdssh/ssh_dial_fallback.go`):**
   - Added `crypto.ConnectWithFallback(ip, user, keyPath, password string) (*ssh.Client, error)` which attempts public key authentication first (if key path exists) and seamlessly falls back to password authentication without emitting stdout error noise.
   - Created `dialNodeWithFallback(c db.SSHConnection, header string) (*ssh.Client, error)` in `cli/cmdssh/ssh_dial_fallback.go`:
     - Tier 1: Try configured `c.KeyPath`.
     - Tier 2: Try standard default keys (`id_ed25519`, `id_rsa`, `id_ecdsa`).
     - Tier 3: Decrypt `c.EncryptedPassword` or query host password from SQLite.
     - Tier 4: Fall back to `ResolveFallbackCredentials(c.Username, c.OS)` from `vmpass.json`.
     - On successful password authentication, automatically encrypts and vaults the password into `ssh_hosts.encrypted_password` for future instant connection.

2. **Concurrent Live Status Probing (`cli/cmdssh/sshjoin_probe.go`):**
   - Implemented `probeHostsLiveStatus(ctx context.Context, hosts []store.SSHHost) []store.SSHHost`:
     - Concurrently dials each host over TCP with a fast 1200ms timeout (`net.DialTimeout`).
     - If TCP connection fails or times out: sets `Status = "○ offline (timeout)"`.
     - If TCP succeeds, verifies authentication: if auth fails, sets `Status = "▲ auth failed"`.
     - If auth succeeds: sets `Status = "● ready"`.
   - Integrated `probeHostsLiveStatus` into `printSJList` (`cli/cmdssh/sshjoin_ls_cmd.go`) so `gitmap ssh nodes` always reflects true live network state.

3. **Status Rendering & Visual Alignment (`cli/cmdssh/sshjoin_table.go`):**
   - Extended table width to 110 characters with a 21-character visual status column:
     - `● ready` displayed in bold green (`constants.ColorGreen`).
     - `○ offline (timeout)` displayed in dim gray (`constants.ColorDim`).
     - `▲ auth failed` displayed in bold yellow (`constants.ColorYellow`).
   - Updated `cli/cmdssh/fleet_parallel_render.go` to cleanly display `[FLEET SKIP] ...: OFFLINE` for unreachable hosts, separating `Offline: N` from true operational failures.

4. **Temporary E2E Integration Test (`cli/tests/e2e/ssh_fleet_auth_fallback_and_node_status_tempe2e_test.go`):**
   - Created comprehensive end-to-end integration test tagged with `tempe2e` verifying:
     - Auth fallback across `w1`, `w2`, and `w3` using password and key permutations.
     - Fast timeout classification (<1500ms) for unreachable subnet addresses (`10.20.0.11`).
     - Status string formatting and ANSI color contract compliance.

## 4. Verification & Prevention

1. **E2E Test Execution:**
   - Ran `$env:RUN_TEMP_E2E="1"; go test -v -tags tempe2e ./tests/e2e/...`.
   - Result: `TestSSHFleetAuthFallbackAndNodeStatus_TempE2E` passed in 0.41s; full suite passed in 4.68s.

2. **Unit Test & Static Analysis:**
   - `go test ./cmdssh/...`: Passed in 36.57s.
   - `python linter-scripts/check-nested-ifs.py`: 0 violations across all files.
   - `python linter-scripts/check-enum-and-boolean.py`: 0 violations across 2,917 source files.
   - `golangci-lint run --issues-exit-code=1 ./...`: 0 issues found (clean exit 0).

3. **Live Execution Verification:**
   - Executed `go run . agm update --ssh`:
     ```text
     ================================================================================
      SSH Fleet Execution Summary: Update agm
      Total: 6 | Succeeded: 3 | Failed: 0 | Offline: 3
     --------------------------------------------------------------------------------
       ALIAS            IP                STATUS                    DURATION  DETAILS
       ───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────
       alpha-win        10.20.0.11        ○ OFFLINE      1501ms  offline: node [alpha-win|10.20.0.11] is unreachable
       beta-linux       10.20.0.12        ○ OFFLINE      1501ms  offline: node [beta-linux|10.20.0.12] is unreachable
       gamma-mac        10.20.0.13        ○ OFFLINE      1501ms  offline: node [gamma-mac|10.20.0.13] is unreachable
       w1               192.168.1.3       SUCCESS         12276ms  OK
       w2               192.168.1.7       SUCCESS         12228ms  OK
       w3               192.168.1.12      SUCCESS         12462ms  OK
     ================================================================================
     ```
   - Executed `go run . ssh nodes`:
     - Displays `alpha-win`, `beta-linux`, `gamma-mac` as `○ offline (timeout)`.
     - Displays `w1`, `w2`, `w3` as `● ready`.
