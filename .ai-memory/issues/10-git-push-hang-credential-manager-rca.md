# RCA: Git Push Hanging Indefinitely on Windows Headless Tasks

## 1. Symptom
When executing `git push origin main`, the command hangs indefinitely in background subshells with zero terminal output and 0 bytes transmitted. In the agent task runner, the task status stays in `RUNNING` for over 10 minutes without completing or failing.

**Process inspection snapshot (`Get-Process -Name "*git*"`)**:
```text
ProcessName                Id  CPU(s)
-----------                --  ------
git                      9472    0.02
git                     16932    0.00
git                     28264    0.00
git-credential-manager  22336    1.42
git-remote-https        13772    0.16
```

---

## 2. Root Cause
Git Credential Manager (`git-credential-manager.exe`) blocked indefinitely attempting an interactive OAuth/GUI prompt within a headless, non-interactive subshell because the cached GitHub credential for `alimtvnetwork` was invalid/expired.

---

## 3. Resolution

1. **Killed Orphaned Processes & Cancelled Stalled Tasks:**
   - Terminated hanging `git-credential-manager.exe` and `git-remote-https.exe` processes.
   - Cancelled background tasks `task-39315` and `task-39322`.

2. **Configured Anti-Hang & Timeout Guards in Repository Git Config:**
   ```bash
   git config credential.interactive never
   git config http.lowSpeedLimit 1000
   git config http.lowSpeedTime 30
   ```
   Now, any missing or expired credential fails immediately in 1.5 seconds (`fatal: Cannot prompt because user interactivity has been disabled`) rather than silently freezing for 10+ minutes.

3. **Re-Authentication Options Provided:**
   - **Option 1 (GitHub Device Login):** Navigate to `https://github.com/login/device` and enter code `5E22-3EAB`.
   - **Option 2 (GitHub Desktop):** GitHub Desktop is running on the workstation; push directly with one click.
   - **Option 3 (Personal Access Token):** Store PAT in Windows Credential Manager via `cmdkey /generic:git:https://github.com /user:alimtvnetwork /pass:<TOKEN>`.

---

## 4. Prevention & Learnings
- **Rule:** Never execute `git push` or `git fetch` in automated agent subshells without non-interactive guards (`credential.interactive=never`, `GCM_INTERACTIVE=never`, `GIT_TERMINAL_PROMPT=0`).
- **Timeout Safety:** Always maintain `http.lowSpeedTime 30` in git configuration to automatically kill stalled HTTP connections.
