# Subtask 05: Master Embedded Runner & Health Scorecard

> **Task Reference:** `214-ubuntu-fleet-automation-and-workstation-governance`  
> **Parent Plan:** [.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md](file:///d:/work/gitmap/.ai-memory/plans/214-ubuntu-fleet-automation-and-workstation-governance.md)  
> **Target Node:** Ubuntu 24.04 LTS (`u1` / `192.168.1.22`)  

---

## 1. Objective & Requirements
- Deliver a single standalone master PowerShell script that houses all bash shell automation internally inside `@' ... '@` here-strings, requiring zero external `.sh` file dependencies on the remote node.
- Provide end-to-end verification and emit an automated scorecard covering all workstation metrics.
- Document full command catalog, error forensic analysis, and future OS setup roadmap.

---

## 2. Implementation & Commands
- Implemented in `d:/work/repo-secrets/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1`.
- Fully documented in `d:/work/repo-secrets/04-ubuntu-migration/step-by-step-log-v3.md`.
- Execution command:
  ```powershell
  pwsh -File d:/work/repo-secrets/04-ubuntu-migration/master-embedded-ubuntu-runner.ps1 -Action verify
  ```

---

## 3. Verification & Live Evidence
- Live scorecard output:
  ```text
  =================== VERIFICATION SCORECARD ===================
   [PASS] SCALE : 1.3999999999999999
   [PASS] WALLPAPER : 'file:///usr/share/backgrounds/warty-final-ubuntu.png'
   [PASS] AUTOMOUNT : active
   [PASS] SYMLINK : /home/a/git-work
   [PASS] ANTIGRAVITY : 2.19.1
   [PASS] REPOS : 74
   [PASS] BRAINS : 98
  ==============================================================
  ```
- 100% of checks passed.
- Status: **DONE**
