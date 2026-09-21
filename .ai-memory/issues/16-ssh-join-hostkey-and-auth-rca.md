# Root Cause Analysis (RCA-16): SSH Join Routing, Known Hosts Collision, Key Masking, and Keyboard-Interactive Auth

## 1. Symptom
1. **`gitmap ssh add <user@ip> <alias>` Failed:**
   Running `gitmap ssh add a@192.168.1.9 w2` threw:
   ```text
   gitmap ssh: execute failed: [E_INTERNAL_ERROR:] ParseSSHTarget: 'add' is a reserved command, not a host target (ctx=map[raw:add])
   ```
2. **`gitmap ssh` Public Key Masking & Missing Clipboard Copy:**
   Invoking `gitmap ssh` displayed a redacted public key (`...[redacted, pass --raw to view]...`) and failed to automatically copy existing keys to the system clipboard.
3. **`WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!` on `gitmap ssh <alias>`:**
   After joining, running `gitmap ssh w2` failed with OpenSSH host key verification failure:
   ```text
   @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
   @    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
   @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
   IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
   Host key verification failed.
   ```
4. **`Permission denied (publickey,password,keyboard-interactive)` on Join/Login:**
   Joining with password failed silently on modern Linux systems enforcing PAM keyboard-interactive auth, yet `gitmap ssh join` reported success without deploying authorized keys.
5. **Stack Traces Stripped on SSH Command Errors:**
   Errors in `InteractiveSSHClient.Run` and `SpawnSSHWithPassword` lacked captured stack traces and caller metadata.

---

## 2. Root Cause
1. **Missing Subcommand Route in `dispatchNodeSSH`:**
   In `cli/cmdssh/ssh.go`, `dispatchNodeSSH` only handled `"join"` and `"sj"`. The `"add"` command was not caught, falling through to `dispatchSSH` -> `dispatchFallbackOrLogin` -> `runSSHLogin` -> `executeSSHLoginWithPassword` -> `ParseSSHTarget("add", ...)`. Because `"add"` is in `reservedSSHCommands`, `ParseSSHTarget` returned `E_INTERNAL_ERROR`.
2. **Key Redaction & Clipboard Omission on Disk Cache:**
   - In `cli/cmdssh/ssh_key_display.go`, `formatDisplayPublicKey` called `maskKeyBlob` by default unless `--raw` was provided.
   - In `cli/cmdssh/sshexisting.go`, `printExistingKeyOnDisk` called `printExistingKeyPublic` but omitted `copyPubKeyAndAnnounce`.
3. **Stale Host Key Retention in `known_hosts`:**
   - In `cli/cmdssh/ssh_known_hosts_scan.go`, `persistTrustedHost` appended newly discovered host keys to `known_hosts` without removing old/stale entries for that host. When a remote host was reinstalled or changed its key type (e.g. ECDSA -> ED25519), the old key remained earlier in `known_hosts`, triggering OpenSSH MITM warnings.
   - In `cli/cmdssh/ssh_login_cmd.go`, `executeSSHLoginWithPassword` never invoked `autoTrustTargetHost` prior to spawning `ssh`.
4. **Omission of PAM `keyboard-interactive` Auth in Go SSH Client:**
   In `cli/cmdssh/sshjoin_enroll.go`, `dialNodeWithPassword` only provided `ssh.Password(pass)` in `ssh.ClientConfig.Auth`. Modern Linux distributions (Ubuntu 22.04+, Debian 12+) frequently disable RFC 4252 password auth in favor of PAM `keyboard-interactive`. Because `dialNodeWithPassword` failed, `session.hasClient` was `false`, `deployHostPublicKey` never ran, and `ExecuteSSHJoinEnrollment` still printed success.
5. **Missing Stack and Caller in `executeClientCmd`:**
   In `cli/cmdssh/ssh_client.go`, `executeClientCmd` instantiated `apperror.AppError` without `Stack` or `Caller`, and `cliexit/handle.go` did not consider `E_INTERNAL_ERROR` when formatting stack traces.

---

## 3. Resolution
1. **Route `add` to `RunSSHJoinCLI`:**
   Updated `dispatchNodeSSH` in `cli/cmdssh/ssh.go` to handle `case "join", "sj", "add": return result.MatchWrapper(RunSSHJoinCLI(args))`.
2. **Unconditional Full Public Key Display & Clipboard Copy:**
   - Updated `cli/cmdssh/ssh_key_display.go` so `formatDisplayPublicKey` returns `strings.TrimSpace(pubKey)` unconditionally; removed redaction logic.
   - In `cli/cmdssh/sshexisting.go`, added `copyPubKeyAndAnnounce(strings.TrimSpace(string(pub)))` in `printExistingKeyOnDisk`.
3. **Auto-Purge Stale Entries in Known Hosts:**
   - In `cli/cmdssh/ssh_known_hosts_scan.go`, `persistTrustedHost` now calls `RemoveFromKnownHostsFile(path, kh.Host)` before `AppendKnownHostFile`.
   - In `cli/cmdssh/ssh_known_hosts_file.go`, enhanced `isMatchingHostLine` with `stripHostPort` to properly match host entries with or without port brackets.
   - In `cli/cmdssh/ssh_login_cmd.go`, called `autoTrustTargetHost(ctx, sshTarget)` in `executeSSHLoginWithPassword` prior to spawning `ssh`.
4. **Support `keyboard-interactive` and Enforce Auth Gate on Join:**
   - In `cli/cmdssh/sshjoin_enroll.go`, updated `dialNodeWithPassword` to supply both `ssh.Password(pass)` and `ssh.KeyboardInteractive(...)`.
   - In `ExecuteSSHJoinEnrollment`, if a password was provided and connection failed, the error is surfaced immediately instead of falsely printing success.
   - In `autoTrustTargetHost`, updated target resolution to include custom ports via `resolveTargetAddr`.
5. **Preserve Error Stack Traces:**
   - In `cli/cmdssh/ssh_client.go`, `executeClientCmd` now populates `Stack: apperror.CaptureStackTrace(apperror.DefaultStackTraceSkip)` and `Caller: apperror.CaptureCaller(apperror.DefaultCallerSkip)`.
   - In `cli/cliexit/handle.go`, updated `isStackTraceEnabled` to include `e.Code == "E_INTERNAL_ERROR"`.
   - In `cli/cmdssh/ssh_askpass.go`, ensured temp scripts use `<temp_dir>/gitmap/ssh/`.
6. **Strictly Avoid Rule:**
   Added a TOTAL BAN rule in `.ai-memory/strictly-avoid.md` prohibiting SSH public key masking and clipboard omission.

---

## 4. Prevention & Learnings
1. **Never Mask Public Keys:** Public keys are inherently public and designed for distribution; redacting them breaks usability and developer automation.
2. **Purge Before Appending Known Hosts:** Whenever trusting or re-trusting a host key, always remove stale records for that host first to prevent OpenSSH key mismatch warnings.
3. **Always Pair Password with Keyboard-Interactive in Go SSH:** Any SSH client dialing Linux servers with password credentials must support both RFC 4252 password auth and PAM keyboard-interactive.
4. **Gate Enrollment on Confirmed Client Dial:** Never report machine enrollment success if the initial SSH connection attempt failed when credentials were provided.
