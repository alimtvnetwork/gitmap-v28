# Subtask 212-02: VMware Shared Folder Deep Verification & Health Audit

> **Parent Plan:** [82-ubuntu-fleet-full-customization-and-embedded-runner.md](../../82-ubuntu-fleet-full-customization-and-embedded-runner.md)  
> **Spec Reference:** [01-architecture-spec.md](../../../../02-spec/21-app/212-ubuntu-fleet-full-customization-and-embedded-runner/01-architecture-spec.md)  
> **Status:** PENDING  
> **Target Node:** Ubuntu U1 (`192.168.1.22`)  
> **Target Files:**  
> - `d:/work/repo-secrets/04-ubuntu-migration/verify-vmware-mount.sh`  

---

## 1. Technical Context & Scope

VMware Workstation Shared Folders (`vmhgfs-fuse`) allow host folders (e.g. `d:\work` or `SharedDirectories`) to be dynamically mapped into `/mnt/hgfs` on Ubuntu guest `u1`. Persistent automounting is configured via `/etc/systemd/system/mnt-hgfs.automount` and `/etc/systemd/system/mnt-hgfs.mount`.

However, guest sleep states, host hypervisor reboots, or FUSE mount drops can cause stale file handles or silent unmounting.

This subtask delivers `verify-vmware-mount.sh`: an idempotent diagnostic and self-healing verification tool that executes an 8-point system audit, checks permissions and write integrity for user `a` (UID 1000), triggers auto-remount remediation when degradation is detected, and exports structured reports (ANSI terminal table and machine-readable JSON).

---

## 2. Technical Specification & Component Contracts

### 2.1 CLI Interface
```bash
./verify-vmware-mount.sh [OPTIONS]
```
- **Flags:**
  - `-j, --json`: Output audit results strictly as a JSON scorecard.
  - `-r, --remount`: Force restart of `mnt-hgfs.automount` and remount `/mnt/hgfs` before probing.
  - `-v, --verbose`: Show full systemd unit logs and mount options.
  - `-h, --help`: Display usage and flag documentation.

---

### 2.2 8-Point Diagnostic Architecture

| Check ID | Inspection Target | Verification Command / Check | Healthy Criteria |
| :--- | :--- | :--- | :--- |
| **CHK-01** | Hypervisor Environment | `systemd-detect-virt` | Returns `vmware` |
| **CHK-02** | FUSE Binary & Tools | `which vmhgfs-fuse` && `dpkg -s open-vm-tools` | Binaries and packages installed |
| **CHK-03** | FUSE Security Policy | `/etc/fuse.conf` | Contains active `user_allow_other` |
| **CHK-04** | Automount Unit State | `systemctl is-active mnt-hgfs.automount` | Returns `active` |
| **CHK-05** | Mount Target Directory | `test -d /mnt/hgfs` | Directory exists with permissions >= 755 |
| **CHK-06** | Host Share Enumeration | `vmware-hgfsclient` | Returns >= 1 configured share |
| **CHK-07** | Read & Traversal Access | `ls -la /mnt/hgfs/<share>` | Readable by non-root user `a` (UID 1000) |
| **CHK-08** | Write & Integrity Test | Create, read back, and delete `.verify-probe-$$` | Data matches and file removes cleanly |

---

### 2.3 Self-Healing Auto-Remount Logic

When `CHK-04`, `CHK-05`, or `CHK-07` detect an unmounted or stale state:
1. Reload systemd daemon: `systemctl daemon-reload`.
2. Restart automount unit: `systemctl restart mnt-hgfs.automount`.
3. Trigger automount tripwire by traversing directory: `ls -la /mnt/hgfs > /dev/null 2>&1`.
4. Re-evaluate mount state (`findmnt -n /mnt/hgfs` or `mountpoint /mnt/hgfs`).
5. If recovered, record status as `RECOVERED`; if still failing, record as `CRITICAL_ERROR`.

---

### 2.4 Diagnostic JSON Schema

When invoked with `--json`, the script emits structured telemetry:
```json
{
  "timestamp": "2026-10-04T21:30:00Z",
  "node": "u1",
  "user": "a",
  "overallStatus": "HEALTHY",
  "checks": [
    { "id": "CHK-01", "name": "Virtualization Check", "passed": true, "details": "vmware" },
    { "id": "CHK-02", "name": "VMware Tools", "passed": true, "details": "open-vm-tools 13.0.10" },
    { "id": "CHK-03", "name": "FUSE Configuration", "passed": true, "details": "user_allow_other enabled" },
    { "id": "CHK-04", "name": "Automount Unit", "passed": true, "details": "active (running)" },
    { "id": "CHK-05", "name": "Mount Directory", "passed": true, "details": "/mnt/hgfs accessible" },
    { "id": "CHK-06", "name": "Host Shared Folders", "passed": true, "details": ["SharedDirectories"] },
    { "id": "CHK-07", "name": "User Read Access", "passed": true, "details": "UID 1000 can read" },
    { "id": "CHK-08", "name": "Write Integrity Probe", "passed": true, "details": "Write/read/delete successful" }
  ]
}
```

---

## 3. Function Decomposition Plan (<= 15 Lines Per Function)

1. `check_virtualization()` (<= 10 lines): Probes `systemd-detect-virt` and sets `is_vmware`.
2. `check_vmware_tools()` (<= 12 lines): Verifies `vmhgfs-fuse` and `open-vm-tools` presence.
3. `check_fuse_config()` (<= 10 lines): Inspects `/etc/fuse.conf` for `user_allow_other`.
4. `check_systemd_automount()` (<= 10 lines): Queries `systemctl is-active mnt-hgfs.automount`.
5. `enumerate_host_shares()` (<= 12 lines): Runs `vmware-hgfsclient` and stores discovered shares.
6. `check_read_access()` (<= 12 lines): Tests directory traversal under `/mnt/hgfs`.
7. `execute_write_probe(share_dir)` (<= 14 lines): Writes temporary probe file, verifies hash/content, and unlinks.
8. `perform_auto_remount()` (<= 14 lines): Restarts automount unit and touches mount path.
9. `render_terminal_table()` (<= 15 lines): Prints formatted ANSI table with Pass/Fail badges.
10. `render_json_scorecard()` (<= 15 lines): Generates compliant JSON payload.
11. `show_verify_help()` (<= 12 lines): Prints CLI arguments and usage instructions.
12. `main("$@")` (<= 15 lines): Parses options, runs 8 checks, performs auto-remount if requested, and formats output.

---

## 4. Coding Guidelines & Invariant Rules

- **Positive Booleans Only:** Use `is_vmware`, `has_fuse_config`, `is_unit_active`, `has_shares`, `is_readable`, `is_writable`. Avoid negative boolean logic.
- **Safe Fallbacks:** Ensure read/write probe creates files with unique PIDs (`probe-$$`) and cleans up in a `trap` handler.
- **Strict Error Handling:** Use `set -euo pipefail`.
- **Zero Git Commands:** Verification tool must strictly execute zero git commands.
- **Relative Path Hygiene:** Use relative links across specifications and plans.

---

## 5. Verification & Acceptance Protocol

### Test Case 1: Healthy System Audit
```bash
ssh u1 "bash -s" < verify-vmware-mount.sh
```
- **Expected Output:**
  - 8/8 checks pass with `[PASS]` ANSI green badges.
  - Overall status: `HEALTHY`.
- **Exit Code:** `0`

### Test Case 2: JSON Output Verification
```bash
ssh u1 "bash -s" < verify-vmware-mount.sh --json | jq .overallStatus
```
- **Expected Output:** `"HEALTHY"`
- **Exit Code:** `0`

### Test Case 3: Self-Healing Remount Verification
```bash
ssh u1 "sudo umount -l /mnt/hgfs; bash -s" < verify-vmware-mount.sh --remount
```
- **Expected Output:**
  - Auto-remount activates `mnt-hgfs.automount`.
  - Share directory `/mnt/hgfs/SharedDirectories` immediately re-emerges and write probe succeeds.
- **Exit Code:** `0`
