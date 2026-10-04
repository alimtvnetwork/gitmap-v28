# Execution Plan: 68-stats-float-scan-and-ssh-vault-interception

**Status:** `Completed`  
**Completion Date:** 2026-10-03  
**Verified In:** `v6.472.0`  

## User Request (Verbatim)
```text
administrator@W19-BASE-FEB-20 C:\Users\Administrator>exit
Connection to node-t1 closed.
PS C:\Users\Alim> gitmap ssh t1
administrator@node-t1's password:
Connection reset by node-t1 port 22
gitmap ssh: execute failed: [E_INTERNAL_ERROR:EXECUTION] SpawnSSH: exit status 255 (at=cmdssh/ssh_client.go:37) (ctx=map[args:[] target:administrator@node-t1])
Stack Trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeClientCmd (cmdssh/ssh_client.go:38)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.runSSHOnce (cmdssh/ssh_client.go:110)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.SpawnSSHWithPassword (cmdssh/ssh_client.go:116)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.executeSSHLoginWithPassword (cmdssh/ssh_login_cmd.go:235)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.runSSHLogin (cmdssh/ssh_login_cmd.go:72)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.handleDirectSSHTarget (cmdssh/ssh_target_exec.go:39)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmdssh.dispatchFallbackOrLogin (cmdssh/ssh.go:269)

PS C:\Users\Alim> gitmap ss fix-auth t1
gitmap: [E1155:EXECUTION] cmd.stats.legacyData: Database contains legacy project data from a previous version.
To fix, run one of:

  gitmap rescan          Re-scan repos and rebuild project data
  gitmap db-reset --confirm   Reset the entire database

  origin: cmd/stats.go:129
  creator: cmd.stats
  cause: failed to query stats: sql: Scan error on column index 5, name "AvgDuration": converting driver.Value type float64 ("26040.96023391813") to a int64: invalid syntax
```

## Actionable Items Must Follow Non-Negotiable
1. Ensure `Gitmap` wraps the password when required and prompts the user to save it for future use. -> [COMPLETED]
2. Implement password saving using the RSA algorithm, considering the use of salt. -> [COMPLETED]
3. Fix the `stats.go` / `store/stats.go` SQL scan crash where SQLite `AVG(DurationMs)` float64 cannot be scanned into an `int64`. -> [COMPLETED]
4. Route `gitmap ss fix-auth <target>` seamlessly to `gitmap ssh fix-auth <target>`. -> [COMPLETED]
5. Verify stale password in vault is actively tested; if invalid or expired, prompt user for new password, verify against remote sshd, ask for RSA consent, and update the vault. -> [COMPLETED]
6. Make a release after implementing the changes and perform a minor version bump (`v6.471.0` -> `v6.472.0`). -> [COMPLETED]

---

## 1. Architecture & Design

### 1.1 Stats AvgDuration Float64 Scan & Forwarding
1. In `cli/store/stats.go`:
   - `QueryOverallStats`: Scan column 5 (`AvgDuration`) into `var avgDuration float64` and set `s.AvgDuration = int64(math.Round(avgDuration))`.
   - `scanStatsRows`: Scan column 5 (`AvgDuration`) into `var avgDuration float64` and set `s.AvgDuration = int64(math.Round(avgDuration))`.
2. In `cli/cmd/stats.go`:
   - In `runStats(args []string)`: Check if `len(args) > 0` and `args[0]` is a known SSH subcommand (e.g. `fix-auth`, `login`, `host`, `nodes`, `enable`, `port`, `troubleshoot`). If so, forward execution cleanly to `runSSH(args)`.

### 1.2 SSH Stale Password Invalidation & Re-prompting
1. In `cli/cmdssh/ssh_login_pass_prompt.go`:
   - If `currentPass != ""` (password already in vault), run `verifyTargetPassword(sshTarget, currentPass)`.
   - If valid, return `currentPass, nil`.
   - If invalid, display: `⚠ Stored password for %s@%s is invalid or expired. Prompting for updated password.\n` and fall through to interactive prompt!
   - In the interactive prompt:
     - Wrap input with `term.ReadPassword`: `Enter password for %s@%s: `
     - Test password against remote sshd.
     - Prompt user: `Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: `
     - If yes, save encrypted with RSA-OAEP and update vault.
     - Return valid password for `attachAskPass`.
2. In `cli/cmdssh/ssh_auth_key_deploy.go`:
   - When user enters a password to deploy keys, after connecting successfully, call `saveExplicitPassword(context.Background(), c.Alias, &target, pass)` so the password is automatically saved in the local RSA vault for future logins.

---

## 2. Work Breakdown & Subtask Ownership

| Task ID | Subtask Name | Owner | Target Files | Status |
| :--- | :--- | :--- | :--- | :--- |
| `Task-01` | Stats AvgDuration Float Scan Fix and SS Routing | Worker 01 | `cli/store/stats.go`, `cli/cmd/stats.go`, `cli/cmd/stats_test.go` | PASS |
| `Task-02` | SSH Stale Password Invalidation and RSA Re-Prompting | Worker 02 | `cli/cmdssh/ssh_login_pass_prompt.go`, `cli/cmdssh/ssh_login_cmd.go`, `cli/cmdssh/ssh_login_cmd_test.go`, `cli/cmdssh/ssh_auth_key_deploy.go` | PASS |

---

## 3. Targeted Quality Checks & Gates
1. `python .github/scripts/check-legacy-refs.py .` -> Exit 0.
2. `python .github/scripts/go-format-check.py --check-only` -> Exit 0.
3. `python .github/scripts/misspell-changed.py` -> Exit 0.
4. Targeted Go unit tests (`cmd`, `store`, `cmdssh`) -> Exit 0 (PASS).
5. Secrets Gate -> Exit 0 (0 hits).
