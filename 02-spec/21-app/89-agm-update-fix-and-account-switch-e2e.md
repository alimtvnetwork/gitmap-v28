# SPEC-APP-89: AGM Update Version Pinning Fix, Parallel Running-Prompts Backup with Media, Account Switch E2E & Threshold Governance

## User Request (Verbatim)

```text
I think for testing purpose, you need to test it, end-to-end test, and confirm that it is working, and also make sure that all the commands that we have crafted for the CLI tool, that has nice help text, okay, with examples. Okay. So let's discuss about the account switch. Okay, you have done the account switch before. But then again, we are trying to do the re-verification one more time, okay? So what I want you to do is... Okay, the way that this is going to work is that first you're going to put the threshold for the account switch at 98%. Okay? And then you use some token or use some credits, add some test prompts to somewhere, okay, so that it starts running. Okay. So once the credit goes less than 98%, it switches. Okay? But before it switches, you know the algorithm. The algorithm is just, once we found the right fit, the top one, we check the refresh. We do a refresh on the account, and if the refresh says the same credit remains, same number follows, then we select that. Okay? But again, we do not do anything yet. The first thing we should do is take a backup of that instance to SQLite database using the AGM. Take the prompts backup, the running prompts, actually. Running prompts only. So you took a backup, let's say. And then, so you need to probably add a running prompts command, like running prompts export or backup and restore, something like this. So it does the restoring, I mean, backuping first. So it should be very small, one to one prompts that will be saved to the database, parallely. So you should try to do it parallely if multiple projects are running. Multiple projects, that single prompt will be saved to the database along with the picture information, picture path, and most of the things that these are attached, it needs to be there, synced. And then when we found the account we want to switch, we do the fast-forward button click. That means dedicate that to fast-forward button. Fast-forward button usually just switches the account, right? So it will go to that account, click on the appropriate button to delegate the call to that appropriate button to refresh the IDE with the new account. Up to this part should be very clear. We have discussed this before. Now, it should restore the running prompt and reinject and see like it's running even using another AGM command, it will see commands are running or not. If not, it will notify in the Telegram. Telegram also should let us know the steps that we mentioned. It should also reveal that. That means it's going to take a backup. The backup means we mentioned like taking a backup of running prompts. So it would show like five projects or six projects, whatever the number is. Projects are getting into backup. Projects name can be there listed. No IDs, but the names of the projects. So once it's saved into the backup using the command, it will again restore the running prompts, and it would reinject, and this status also needs to be in the Telegram, also in the email altogether. Do you understand? Okay, write into the spec so that you never forget And also at the end, once all this testing is done. So first you code, because if you switch it, then the code running would be stopped. Okay? Make sure that you have a process that makes sure that the prompts are running. Okay? That is very, very important. Now, at the end, you would have some count that would give you all this information and finalize the process. So at the end, what you do after the testing and everything is done, you reduce the threshold to, let's say, 15%. Okay? Keep it under 15%, and then you make a minor bump and release, and check the CI/CD using GitLab EE licenT. That is all green. So that's it. That's what I'm expecting for. So you need to have this testing. Now, the same thing, you want to test it, end-to-end test in the instance mode. Or you could try to create a new instance and try to test this, that does this work on instance mode or not. And then finally you remove that instance. Okay? So you do the end-to-end testing here, no problem in the machine. You can try to do anything that you like. If anything goes wrong, I can revert this back. I have a snapshot. So you don't have to worry about this. Do you understand? Can you please follow through all this and then finally make a release? That means before you release, you also change the default threshold to 15%. Okay? And make sure that you have the proper email sending and also email checking. You need to check the email. If the email is connected, you need to check in the email if this user profile or the account is already selected by some other DM. If it is, then you need to skip and move to the next one. Remember that, that's a very crucial step I forgot to remind you. Also, if we have the super base, then it is very easy. Every one of the system will check who is not using the account and has the highest value by the algorithm you have to specify. Then it will just switch to this account. That is the standard process. Is it clear? Do you have any question and confusion?

---

You need to fix it properly

PS D:\work> gitmap agm update
Updating Antigravity Manager via PowerShell...
    [*] Detected pinned version from download URL / invocation: v2.336.0
...
    [!] Attempt 1 failed for release v2.336.0: Failed to download release package for v2.336.0 from https://github.com/alimtvnetwork/Antigravity-Manager/releases/download/v2.336.0/agm-alim_2.336.0_x64-setup.exe
...
  ✖ Post-install verification failed for ag-manager.
gitmap: [E9000:EXECUTION] cmd.verifyInstallation: post-install verification failed for "ag-manager" (binary: "ag-manager")
```

---

## 1. Architectural Overview & Root Causes

### 1.1 `gitmap agm update` Version Hijack & Verification Fix
- **Root Cause 1 (Pinned Version Hijack):** When `gitmap agm update` executed `install.ps1` without an explicit `-Version` flag, `Resolve-PinnedVersion` inside `install.ps1` scanned `ConsoleHost_history.txt` (`Get-PSReadLineOption.HistorySavePath`), matched an unrelated command version (`v2.336.0`), and attempted to download a non-existent `agm-alim_2.336.0_x64-setup.exe` asset.
- **Fix 1:**
  - In `cli/cmdinstall/installagmanager.go`, dynamically resolve the latest release version from `https://api.github.com/repos/alimtvnetwork/Antigravity-Manager/releases/latest` (or redirect header from `/releases/latest`) when `version == ""`, pass `-Version '<resolved>' -Update -NoLaunch` with a cache-busting query string (`?cb=<timestamp>`), and inject `AGM_VERSION=<resolved>` into the subprocess environment.
- **Root Cause 2 (Post-Install Binary Verification Failure):** `isAgManagerInstalled()` and `findAgManagerWindowsPath()` searched only for `ag-manager.exe` or `Antigravity.Tools.exe`, missing `C:\Users\Administrator\AppData\Local\Programs\agm-alim\agm-alim.exe`.
- **Fix 2:**
  - Register `agm-alim` and `%LOCALAPPDATA%\Programs\agm-alim\agm-alim.exe` as the primary binary candidate in `cli/cmdinstall/installagmanager_windows.go` and `cli/cmdinstall/installagmanager_exec.go`.

### 1.2 Parallel 1-to-1 Running Prompts Backup & Restore with Media Attachments
- **Parallel Discovery & Persistence:**
  - Across all active Antigravity projects, capture running/active prompts in parallel (`sync.WaitGroup` worker pool, 1-to-1 per running project) into SQLite (`gitmap-running-prompts.db` / `backup-prompts/sql.db`).
  - Extract and persist any attached picture/media paths (`MediaPathsJSON`, image attachments, screenshot references) alongside the prompt text, project name, and sequence metadata.
- **Privacy & Naming Rule:**
  - Notifications (Telegram, Email, CLI summary) MUST list human-readable **project names** (e.g. `gitmap`, `Antigravity-Manager`, `scripts-fixer`) and count (e.g. `4 projects backed up`), and MUST NEVER expose raw project UUIDs.

### 1.3 Account Switch Algorithm, Email/Supabase Lock Check, Fast-Forward Delegation & Threshold Lifecycle
1. **Threshold Evaluation:**
   - Support configurable credit threshold (`--threshold`, default `15%`, testable at `98%` via `--threshold 98` or `account-switch test --threshold 98`).
2. **Candidate Ranking & Double-Check Refresh:**
   - Rank available accounts by highest remaining credit.
   - For the top candidate, perform an explicit **Account Refresh** probe to verify the remaining credit matches the pre-refresh value (`refreshedCredit == candidateCredit`).
3. **Distributed Multi-VM Lock Check (Supabase & Email):**
   - Query **Supabase** (if configured) and **Email inbox/lock state** (if email is connected) to check if the candidate account is already claimed/active on another VM or DM instance.
   - If the candidate account is in use by another VM/DM, skip it and evaluate the next highest candidate.
4. **Pre-Switch Parallel Running Prompts Backup:**
   - Before switching accounts, execute parallel 1-to-1 running prompts SQLite backup across all active workspaces (preserving prompt text + attached image paths).
   - Dispatch step notification to **Telegram** and **Email** listing the count and names of projects backed up (zero raw UUIDs).
5. **Fast-Forward Account Switch Delegation:**
   - Trigger the Fast-Forward (`fast-forward` / `switch-account`) action to switch the active profile and refresh the Antigravity IDE session with the selected account.
6. **Post-Switch Restore, Re-Injection & Liveness Verification:**
   - Restore the backed-up running prompts from SQLite and re-inject them into their respective workspaces.
   - Verify via `gitmap agy running-projects` / `running-prompts ls` that all prompts are actively running.
   - Dispatch completion status report to **Telegram** and **Email** detailing backed-up project names, target account, and re-injection verification status.
7. **Instance Mode E2E Verification & Final Threshold Reset (`15%`):**
   - Support `--instance <name>` (creating an isolated test instance directory, running the 98% threshold switch + backup/restore lifecycle end-to-end, verifying all steps, and cleaning up the test instance).
   - Reset default production threshold to `15%` (`DefaultAccountSwitchThreshold = 15`), execute minor version bump & release, and verify CI/CD pipeline is green.
