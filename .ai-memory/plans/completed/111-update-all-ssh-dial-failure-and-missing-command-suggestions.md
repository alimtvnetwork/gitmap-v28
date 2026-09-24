# Plan 111: SSH Fleet Update Decrypted Auth, Liveness Probe & Command Suggestions

## Metadata
- Spec Reference: [02-spec/21-app/160-ssh-fleet-update-decrypted-auth-liveness-probe-and-command-suggestions.md](../../../02-spec/21-app/160-ssh-fleet-update-decrypted-auth-liveness-probe-and-command-suggestions.md)
- RCA Reference: [02-spec/22-app-issues/46-update-all-ssh-dial-failure-and-missing-command-suggestions-rca.md](../../../02-spec/22-app-issues/46-update-all-ssh-dial-failure-and-missing-command-suggestions-rca.md)
- Status: Completed
- Duration / Cycles: 1 cycle (Parent Task N-Step Loop v2.5.0)

## Overview & Scope
Resolved `gitmap update all` failure across registered cluster nodes and added user command suggestion discovery:
1. `gitmap update all` reported `FAILED: ssh dial failed for administrator@...` on all 6 nodes because stored encrypted passwords were being passed as raw ciphertext to SSH authentication dialers.
2. `cmdupdate` dialed offline nodes (`alpha-win`, `beta-linux`, `gamma-mac`) without pre-flight connectivity checks, hanging until TCP timeouts and misreporting them as SSH authentication dial failures instead of `OFFLINE`.
3. Remote machines were hardcoded to `OS: "linux"` in update targets, running curl scripts on Windows machines.
4. Typos and unknown subcommands (`gitmap update al`, `gitmap ssh hosts`, `updat`, `ss`, `shs`) lacked helpful "Did you mean?" suggestions.

## Outcomes & Verification
- Created `cli/crypto/ssh_password.go` with `crypto.DecryptStoredPassword()` supporting AES-GCM (primary and fallback keys), RSA-OAEP private key decryption, and plaintext fallback.
- Exported and integrated `cmdssh.ProbeRemoteOSType()` in `cmdupdate` to dynamically resolve PowerShell vs bash/sh install commands.
- Added pre-flight TCP connectivity check `CheckConnLivenessFn` with 1000ms timeout in `cmdupdate`, marking offline nodes as `OFFLINE (15ms)` and isolating unit tests from network dependencies.
- Added `update`, `ua`, `ssh`, `ssh-join`, `ssh-exec`, `ssh-nodes`, `install-exec`, `deploy` to `primaryTopCommands` in `cli/cmd/rootsuggest.go`.
- Added `candidateLenDiff` tiebreaker to `rankCandidateCommands` for high-accuracy typo suggestions on short inputs.
- Implemented `suggestUpdateTarget` and `handleUnknownUpdateTarget` in `cli/cmd/rootutility.go`.
- Implemented `suggestSSHSubcommand` in `cli/cmdssh/ssh_login_cmd.go` to suggest `nodes` when `hosts` is requested.
- Verified unit test suite: `cmdupdate` (PASS, 0.059s), `crypto` (PASS, 0.021s), `rootsuggest` (PASS, 0.030s), `cmdssh` (PASS, 33.458s).
- Verified linters: 0 violations across nested ifs, booleans, relative paths, and error management.
- Live CLI verification:
  - `gitmap update all`: Offline nodes (`alpha-win`, `beta-linux`, `gamma-mac`) flagged `OFFLINE (15ms)`, online nodes (`w1`, `w2`, `w3`) updated to `v6.337.0` (`SUCCESS`).
  - `gitmap ssh hosts`: suggested `gitmap ssh nodes`.
  - `gitmap update al`: suggested `gitmap update all`.
