# Plan 82: Ubuntu Fleet Full Customization & Self-Contained Embedded Runner

> **Plan Status:** Active  
> **Traceability IDs:** Subtask 212-01 .. Subtask 212-05  
> **Spec Reference:** [02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md](../../02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md)  
> **Target Node:** Ubuntu U1 (`node-u1`)  
> **Execution Location:** `./repo-secrets\04-ubuntu-migration\`  

---

## 1. Architectural Context & Subtask Mapping

| Subtask ID | Focus Area | Target Deliverable | Status |
| :--- | :--- | :--- | :--- |
| **Subtask 212-01** | Wallpaper & GUI Launching | `set-desktop-wallpaper.sh`, `launch-gui-app.sh` | In-Progress |
| **Subtask 212-02** | VMware Shared Folders | `verify-vmware-mount.sh` | In-Progress |
| **Subtask 212-03** | Antigravity Brain Migration | `sync-antigravity-deep.ps1` | In-Progress |
| **Subtask 212-04** | Embedded Master Runner | `master-embedded-ubuntu-runner.ps1` | In-Progress |
| **Subtask 212-05** | Retrospective Log & Roadmap | `step-by-step-log-v2.md` | In-Progress |

---

## 2. 5-Subtask Detailed Breakdown

### Subtask 212-01: GNOME Desktop Wallpaper & Remote GUI App Launch
- Author `set-desktop-wallpaper.sh` setting both `picture-uri` and `picture-uri-dark` via D-Bus session.
- Author `launch-gui-app.sh` launching applications via `systemd-run --user` (Antigravity, Text Editor, VS Code).
- Test remotely on U1 over SSH.

### Subtask 212-02: VMware Shared Folder Deep Verification
- Author `verify-vmware-mount.sh` checking `/mnt/hgfs` mount status, `vmware-hgfsclient`, permissions, and auto-remount.
- Verify live on U1.

### Subtask 212-03: Antigravity Deep Conversation History & Workspace Portability
- Author `sync-antigravity-deep.ps1` exporting conversation histories, tar streaming over SSH, and running path normalizer on U1.
- Update `conversation_summaries.db` on U1.

### Subtask 212-04: Self-Contained Master PowerShell Runner
- Author `master-embedded-ubuntu-runner.ps1` embedding all bash scripts as here-strings `@' ... '@` and streaming to remote bash with CRLF stripping (`tr -d '\r' | bash -s`).
- Support parameters `-Step all|wallpaper|gui|vmware|antigravity` and `-DryRun`.

### Subtask 212-05: Detailed Engineering Log V2 & Verification
- Author `step-by-step-log-v2.md` capturing all commands, thinking, outputs, errors, and future OS setup roadmap.
