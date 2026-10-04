# Subtask 68-02: SSH Stale Password Invalidation, Live Dial Verification & RSA Vault Update

**Subtask ID:** 68-02  
**Status:** Pending  
**Spec Reference:** `cli/cmdssh/ssh_login_pass_prompt.go`, `cli/cmdssh/ssh_login_cmd.go`  
**Assigned Worker:** Worker 02  
**Target Files:**  
- `cli/cmdssh/ssh_login_pass_prompt.go` [MODIFY]  
- `cli/cmdssh/ssh_login_cmd.go` [MODIFY]  
- `cli/cmdssh/ssh_login_cmd_test.go` [MODIFY]  
- `cli/cmdssh/ssh_auth_key_deploy.go` [MODIFY]  

---

## 1. Technical Context & Scope
If `ssh_hosts` contains an outdated, invalid, or expired password:
1. `interceptSSHPasswordIfNeeded` previously saw `currentPass != ""` and bypassed prompting.
2. OpenSSH ran with the stale password, failed, and fell back to prompting directly on the terminal: `administrator@node-t1's password:`.
3. If the remote host terminated the connection, OpenSSH crashed with exit status 255 (`Connection reset by node-t1 port 22`).

## 2. Implementation Directives
1. In `cli/cmdssh/ssh_login_pass_prompt.go`:
   - In `interceptSSHPasswordIfNeeded(ctx context.Context, target string, sshTarget *SSHTarget, currentPass string)`:
     - If `currentPass != ""`:
       Actively test with remote sshd via `verifyTargetPassword(sshTarget, currentPass)`.
       If valid: return `currentPass, nil`.
       If invalid:
         Print `⚠ Stored password for %s@%s is invalid or expired. Prompting for updated password.\n`
         Proceed directly to interactive prompt!
     - If public key works (`hasKeyAuth(sshTarget)`): return `"", nil`.
     - If not interactive terminal: return `"", nil`.
     - In interactive prompt:
       Prompt user via masked `term.ReadPassword`: `Enter password for %s@%s: `
       Verify with remote sshd (up to 3 attempts).
       Prompt for RSA consent:
         `Do you want to save this password for future use? (y/n) [Encrypted locally with RSA algorithm]: `
       If affirmative: encrypt with RSA-OAEP and save to `ssh_hosts` via `saveExplicitPassword`.
       Return verified password.
2. In `cli/cmdssh/ssh_auth_key_deploy.go`:
   - In `promptAndConnectSSH`:
     When user enters a password to deploy keys, after connecting successfully, call `saveExplicitPassword(context.Background(), c.Alias, &target, pass)` so the password is automatically saved in the local RSA vault for future logins!
3. In `cli/cmdssh/ssh_login_cmd_test.go`:
   - Add unit tests validating:
     - Stored password invalidation when `verifyTargetPassword` returns false.
     - Valid stored password bypass when `verifyTargetPassword` returns true.
