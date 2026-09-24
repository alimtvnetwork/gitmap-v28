# Spec 160: SSH Fleet Update Decrypted Auth, Liveness Probe, and Command Suggestions

## User Request (Verbatim)

```text
Find the root cause of the update all issue and fix it properly

PS C:\Users\Administrator> gitmap ssh nodes

  ALIAS            ROLE           HOST (IP:PORT)         USER           STATUS     ENROLLED
  ----------------------------------------------------------------------------------------------------
  w3               worker         192.168.1.12:22        administrator  ● ready    2026-09-24 20:24:25
  w2               worker         192.168.1.7:22         administrator  ● ready    2026-09-24 20:24:25
  w1               worker         192.168.1.3:22         administrator  ● ready    2026-09-24 20:24:24

  Total: 3 registered node(s)

PS C:\Users\Administrator> gitmap ssh exec ip

  Notice: The following machine(s) are currently OFF or unreachable:
    • [alpha-win | 10.20.0.11] (machine is off)
    • [gamma-mac | 10.20.0.13] (machine is off)
    • [beta-linux | 10.20.0.12] (machine is off)

  ● Commands injected for 3 machine(s) [w2:192.168.1.7, w1:192.168.1.3, w3:192.168.1.12] (batch: 50, timeout: 60s)
  ────────────────────────────────────────────────────────────────────────────────
  ...
  Execution completed in 4.39s (3 succeeded, 0 failed, 3 offline)

PS C:\Users\Administrator> gitmap update all

[FLEET UPDATE] Updating 'all' across 6 cluster node(s) in parallel...

  ✖ [alpha-win|10.20.0.11] FAILED: ssh dial failed for administrator@10.20.0.11 (7303ms)
  ✖ [beta-linux|10.20.0.12] FAILED: ssh dial failed for administrator@10.20.0.12 (7289ms)
  ✖ [gamma-mac|10.20.0.13] FAILED: ssh dial failed for administrator@10.20.0.13 (7300ms)
  ✖ [w1|192.168.1.3] FAILED: ssh dial failed for administrator@192.168.1.3 (11105ms)
  ✖ [w2|192.168.1.7] FAILED: ssh dial failed for administrator@192.168.1.7 (11090ms)
  ✖ [w3|192.168.1.12] FAILED: ssh dial failed for administrator@192.168.1.12 (11075ms)

================================================================================
 SSH Fleet Update Summary [all]: Total: 6 | Succeeded: 0 | Failed: 6 | Offline: 0 | Excluded: 0
--------------------------------------------------------------------------------
  ALIAS            IP                STATUS                DURATION  TELEMETRY DETAILS
  ───────────────────────────────────────────────────────────────────────────────────────
  alpha-win        10.20.0.11        FAILED         7303ms  ssh dial failed for administrator@10.20.0.11
  beta-linux       10.20.0.12        FAILED         7289ms  ssh dial failed for administrator@10.20.0.12
  gamma-mac        10.20.0.13        FAILED         7300ms  ssh dial failed for administrator@10.20.0.13
  w1               192.168.1.3       FAILED        11105ms  ssh dial failed for administrator@192.168.1.3
  w2               192.168.1.7       FAILED        11090ms  ssh dial failed for administrator@192.168.1.7
  w3               192.168.1.12      FAILED        11075ms  ssh dial failed for administrator@192.168.1.12
================================================================================

here the 3 machines are on as you can see but it is not updating if a command is mistaken then it shouk,d should show the ssuggestions pelasede
```

## System Architecture & Invariants

1. **Decrypted Password SSH Authentication:**
   - Saved passwords in the node registry are encrypted at rest using AES-GCM or RSA-OAEP.
   - Any remote execution or update command (`cmdupdate`, `cmdmacro`, `cmdssh`) calling SSH dialers must decrypt the stored ciphertext into plaintext memory via `crypto.DecryptStoredPassword()` before initiating SSH password auth. Passing raw ciphertext causes remote SSH auth rejection.

2. **Pre-Flight Connectivity Liveness Probe:**
   - Before attempting heavy SSH connection protocols or commands, `cmdupdate` runs a 1000ms TCP connectivity probe using `CheckConnLivenessFn` against target IP and port.
   - Offline or unreachable nodes are marked immediately with status `OFFLINE`, duration ~15ms, and details `machine is off or unreachable (<reason>)`.
   - Offline nodes are excluded from remote command execution and tallied under `Offline: N` in metrics rather than falsely counted as failed authentication dials.

3. **Dynamic Remote OS Resolution:**
   - Remote machines must not be assumed to run Linux. After SSH connection establishment, `cmdupdate` probes the actual operating system via `cmdssh.ProbeRemoteOSType(client)`.
   - Windows targets receive PowerShell-based update scripts (`powershell -NoProfile -ExecutionPolicy Bypass -Command ...`), while Unix targets receive shell-based bash/sh scripts.

4. **Typo and Unknown Subcommand Suggestions:**
   - Top-level commands (`update`, `ua`, `ssh`, `ssh-join`, `ssh-exec`, `ssh-nodes`, `install-exec`, `deploy`) are registered in `primaryTopCommands`.
   - `suggestTopLevelCommands()` ranks candidates with length-delta tiebreaking (`candidateLenDiff`), ensuring short inputs like `updat` and `ss` accurately suggest `update` and `ssh`.
   - `gitmap update <unknown>` prints available update subcommands and suggests closest matches (e.g., `al` or `hosts` -> `all`, `inventory` -> `ls`).
   - `gitmap ssh <alias>` checks whether the unknown alias matches known SSH subcommands (e.g., `hosts` -> `nodes`) and displays `Did you mean: gitmap ssh nodes?`.

## Verification Commands

```bash
# 1. Verify cmdupdate unit tests and mock liveness
cd cli && go test -v -count=1 ./cmdupdate/...

# 2. Verify crypto password decryption unit tests
go test -v -count=1 ./crypto/...

# 3. Verify root command suggestions unit tests
go test -v -count=1 ./cmd -run "TestSuggest"

# 4. Live CLI verification
gitmap update all --dry-run
gitmap update al
gitmap ssh hosts
```
