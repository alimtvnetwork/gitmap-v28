# Issue 46 RCA: SSH Fleet Update Authentication Dial Failure & Missing Command Suggestions

## 1. Reproduction

1. **Authentication Dial Failure on Online Nodes:**
   - Execute `gitmap ssh nodes`: 3 nodes (`w1`, `w2`, `w3`) enrolled and ready.
   - Execute `gitmap ssh exec ip`: successfully detects offline machines (`alpha-win`, `beta-linux`, `gamma-mac`) and executes on online nodes (`w1`, `w2`, `w3`).
   - Execute `gitmap update all`:
     - Reports `FAILED: ssh dial failed for administrator@... (11075ms)` on all online nodes (`w1`, `w2`, `w3`).
     - Reports `FAILED: ssh dial failed for administrator@... (7300ms)` on all offline nodes (`alpha-win`, `beta-linux`, `gamma-mac`).

2. **Missing Command Suggestions for Typo & Subcommands:**
   - User mistypes `gitmap ssh hosts` (aiming to inspect registered machines): output displays `SSH host alias 'hosts' not found in registry.` without suggesting `gitmap ssh nodes`.
   - User runs `gitmap update al` or `gitmap update apps`: output silently falls through to local source repo rebuild instead of validating and suggesting `gitmap update all` or `gitmap update ls`.

## 2. Root Cause Analysis

1. **Raw Ciphertext Passed as SSH Password (`cli/cmdupdate` and `cli/cmdmacro`):**
   - In `loadDefaultFleetTargets` and `mergeFleetTargets`, node targets were populated with `Password: h.EncryptedPassword`.
   - `tryDialFleetPassword` in `cli/cmdupdate/update_fleet.go` directly invoked `crypto.ConnectWithPassword(target.IP, target.Username, target.Password)`, passing the raw AES-GCM ciphertext as the plaintext SSH password.
   - The remote SSH daemon rejected authentication on every node, producing `ssh dial failed for ...`.

2. **Absence of Pre-Flight Liveness Check:**
   - Unlike `gitmap ssh exec` which performs `isNodeAvailable` / `CheckConnLiveness`, `cmdupdate` immediately attempted heavy SSH dials on all registered nodes in parallel.
   - Offline nodes hung until TCP timeouts elapsed (7–11s) and were misreported as failed SSH authentication attempts instead of `OFFLINE`.

3. **Target OS Hardcoded to Linux:**
   - Targets loaded from `store.SSHHost` were populated with `OS: "linux"` regardless of remote host OS, attempting to run bash/sh curl scripts on Windows machines.

4. **Missing Command Suggestions:**
   - `primaryTopCommands` in `cli/cmd/rootsuggest.go` omitted top-level commands: `update`, `ua`, `ssh`, `ssh-join`, `ssh-exec`, `ssh-nodes`, `install-exec`, `deploy`.
   - `rankCandidateCommands` lacked length-delta tiebreaking, causing short typo inputs (`updat`, `ss`, `shs`) to be eclipsed by unrelated long candidates.
   - `gitmap update <unknown>` lacked argument validation and suggestion routing.
   - `gitmap ssh <alias>` lacked suggestion routing for known subcommands (e.g. `hosts` -> `nodes`).

## 3. Corrective Implementation

1. **Unified Password Decryption (`cli/crypto/ssh_password.go`):**
   - Implemented `crypto.DecryptStoredPassword(cipherText string) (string, error)` supporting AES-GCM (with primary and fallback keys), RSA-OAEP private key decryption, and plaintext fallback.
   - Updated `dialFleetSSH`, `tryDialFleetPassword` (`cmdupdate`), and `tryDialPassword` (`cmdmacro`) to decrypt stored passwords prior to SSH dialing.
   - Updated `resolveTCPAddress` in `cli/crypto/ssh_client.go` to handle IP addresses with existing port specifications (`ip:port`).

2. **Pre-Flight Connectivity Liveness Probe (`cli/cmdupdate/update_fleet.go`):**
   - Integrated `CheckConnLivenessFn` with a 1000ms TCP connection check before attempting SSH handshake.
   - Offline machines are classified immediately as `OFFLINE` with duration ~15ms and details `machine is off or unreachable (<reason>)`.
   - Metrics updated to display `Total: N | Succeeded: X | Failed: Y | Offline: Z | Excluded: E`.

3. **Dynamic Remote OS Probing (`cli/cmdupdate/update_fleet.go`, `update_fleet_ls.go`):**
   - Exported `cmdssh.ProbeRemoteOSType(client *ssh.Client) string`.
   - Dynamically resolves remote OS after connection to determine whether PowerShell (`powershell -NoProfile -ExecutionPolicy Bypass ...`) or bash/sh curl scripts should execute.

4. **Command & Subcommand Suggestions (`cli/cmd`, `cli/cmdssh`):**
   - Added `update`, `ua`, `ssh`, `ssh-join`, `ssh-exec`, `ssh-nodes`, `install-exec`, `deploy` to `primaryTopCommands`.
   - Enhanced `rankCandidateCommands` with `candidateLenDiff` tiebreaker so typo inputs rank closest length matches first.
   - Implemented `suggestUpdateTarget` and `handleUnknownUpdateTarget` in `cli/cmd/rootutility.go`.
   - Implemented `suggestSSHSubcommand` in `cli/cmdssh/ssh_login_cmd.go` to suggest `nodes` when `hosts` is requested.

## 4. Verification & Prevention

1. **Unit Test Verification:**
   - `go test -v -count=1 ./cmdupdate/...`: 100% PASS (0.059s).
   - `go test -v -count=1 ./crypto/...`: 100% PASS (0.021s).
   - `go test -v -count=1 ./cmd -run "TestSuggest"`: 100% PASS (0.030s).
   - `go test -v -count=1 ./cmdssh/...`: 100% PASS (33.458s).

2. **Linter Suite Verification:**
   - `python linter-scripts/check-nested-ifs.py`: 0 violations.
   - `python linter-scripts/check-enum-and-boolean.py`: 0 violations.
   - `python linter-scripts/check-error-management.py`: 0 violations.
   - `python linter-scripts/check-relative-paths.py`: 0 violations.

3. **Live Execution Output:**
   - `gitmap ssh hosts`:
     `gitmap ssh: SSH host alias 'hosts' not found in registry.`
     `  Did you mean: gitmap ssh nodes?`
   - `gitmap update al`:
     `gitmap update: Unknown update target 'al'.`
     `  Did you mean: gitmap update all?`
   - `gitmap update all`:
     - Offline nodes (`alpha-win`, `beta-linux`, `gamma-mac`): `OFFLINE (15ms)`.
     - Online nodes (`w1`, `w2`, `w3`): Successfully authenticated, probed Windows OS, and upgraded to `v6.337.0` (`SUCCESS`).
