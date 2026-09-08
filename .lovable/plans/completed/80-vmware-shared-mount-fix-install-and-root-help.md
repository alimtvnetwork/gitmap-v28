# 80-vmware-shared-mount-fix-install-and-root-help.md

**Title:** VMware Shared Folders Mount Resilience, Install Subcommand & Root Help Integration  
**Status:** Completed  
**Parent Task Workflow:** Parent Task N-Step Continuous Loop & Multi-Agent Orchestration (Prompt v2.1.0)  
**Budget (N):** 100 steps (Phase 1: Steps 1..50, Phase 2: Steps 51..100)  
**Target Codebase:** `gitmap/constants/`, `gitmap/cmd/`, `gitmap/helptext/`  

---

## 1. Task Overview & Root Cause Analysis

### A. The Mount Failure Bug
- **Reported Error:**
  ```text
  ▶ gitmap vmware shared enable
    ✓ Verified mount point /mnt/hgfs
  gitmap vmware: execute failed: [E4003:EXECUTION] cmd.vmware.mountHostShare: mount failed: Error -107 cannot open connection!
   (exit status 149) (at=cmd/vmware_shared.go:115) (creator=cmd.vmware)
  ```
- **Root Cause Analysis (RCA):**
  1. Linux kernel error code -107 (`ENOTCONN` - Transport endpoint is not connected) is emitted by `vmhgfs-fuse` when the guest FUSE process cannot establish a backchannel communication with the host VMware hypervisor.
  2. The primary cause is that the host virtual machine settings have "Shared Folders" disabled (or no folders added under VM -> Settings -> Options -> Shared Folders).
  3. Secondary causes include:
     - `open-vm-tools` or `vmtoolsd` service not currently running in the guest.
     - Stale or corrupted previous FUSE mount at `/mnt/hgfs` holding a dead descriptor.
     - Lack of `open-vm-tools-desktop` package or missing `user_allow_other` in `/etc/fuse.conf`.
  4. The previous implementation in `cmd/vmware_shared.go` blindly ran `sudo vmhgfs-fuse ...` without:
     - Unmounting stale mounts first (`fusermount -u` or `umount -l`).
     - Starting/restarting `open-vm-tools` service.
     - Verifying exported host shares with `vmware-hgfsclient`.
     - Attempting alternative mount mechanisms (`mount -t fuse.vmhgfs-fuse`).
     - Emitting user-actionable remediation instructions when Error -107 occurs.

### B. Missing Install Command
- `gitmap vmware` lacked an `install` subcommand.
- `gitmap install` lacked `vmware` / `open-vm-tools` tooling package definitions in `constants_install.go`, `install_packages.go`, `install_packages_extra.go`, and `installverify.go`.

### C. Missing Root Help Entry
- `vmware` (alias `vm`) was completely omitted from root help screens (`constants_helpgroups.go`, `rootusage_groups.go`, and `rootusagefilter_rows.go`).

---

## 2. Task-Specific Rule Set & Architectural Invariants

1. **Rule 1 (Resilient Mount Pipeline):** The VMware mount engine must implement a 4-tier resilience pipeline:
   - Tier 1: Pre-mount cleanup (unmount stale mounts if present).
   - Tier 2: Service assurance (`systemctl start open-vm-tools`).
   - Tier 3: Primary `vmhgfs-fuse` attempt with fallback to `mount -t fuse.vmhgfs-fuse`.
   - Tier 4: Diagnostic detection using `vmware-hgfsclient` to provide clear step-by-step host UI instructions if Error -107 occurs.
2. **Rule 2 (Discrete Executions):** Never invoke raw shell strings via `sh -c` or `cmd /c`. All commands (`vmhgfs-fuse`, `fusermount`, `systemctl`, `apt-get`, `crontab`) must pass structured argument slices to `exec.Command(name, args...)`.
3. **Rule 3 (Universal AppError Wrapping):** All errors must return typed `*apperror.AppError` values with error code `E4003:EXECUTION` or `E4001:VALIDATION`.
4. **Rule 4 (Coding Guidelines):** Functions must strictly follow the repository limits (<= 15 lines), blank line before returns, and positive boolean names (`is*`, `has*`).

---

## 3. Subtask Decomposition Ledger

| Subtask File | Scope & Deliverables | Target Files |
| :--- | :--- | :--- |
| `01-task-vmware-mount-resilience.md` | Fix mountHostShare with stale unmount, service start, fallback mount, and Error -107 diagnostics | `gitmap/cmd/vmware_shared.go`, `gitmap/cmd/vmware_shared_test.go` |
| `02-task-vmware-install-command.md` | Add `gitmap vmware install` subcommand and integrate `vmware` into `gitmap install` tooling engine | `gitmap/cmd/vmware.go`, `gitmap/cmd/vmware_install.go`, `gitmap/constants/constants_install.go`, `gitmap/cmd/install_packages.go`, `gitmap/cmd/install_packages_extra.go`, `gitmap/cmd/installverify.go` |
| `03-task-root-help-integration.md` | Register `vmware` in root helptext groups, compact listing, filter rows, and update markdown help | `gitmap/constants/constants_helpgroups.go`, `gitmap/cmd/rootusage_groups.go`, `gitmap/cmd/rootusagefilter_rows.go`, `gitmap/helptext/vmware.md` |
| `04-task-verification-and-ci.md` | Unit test execution, linters validation (nested ifs, booleans, error management), and binary sync | Unit tests, `linter-scripts/`, `bin/gitmap.exe` |

---

## 4. Verification & Quality Gates

- Unit tests for vmware shared, mount error diagnostics, install command, and help integration.
- `python linter-scripts/check-nested-ifs.py` -> 0 violations.
- `python linter-scripts/check-boolean-guidelines.py` -> 0 violations.
- `python linter-scripts/check-error-management.py` -> 0 violations.
- Full binary compilation and synchronization across all 4 executable targets.
