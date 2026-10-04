# Spec 207: GitMap Database Reset, SSH Key Lifecycle, Cross-OS Firewall, and DetectedProject Foreign Key RCA

## Status: Active
- **Spec ID:** 207
- **Slug:** `207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca`
- **Scope:** Database Engine, SSH Key Lifecycle, Cross-OS Firewall, Web UI, Scanner Project Detection
- **Created At:** 2026-10-04
- **Parent Goal:** Implement full database purge and reseed (`gitmap reset`), reorder SSH terminal output, add `gitmap ssh view`, safeguard `gitmap ssh create` with overwrite prompt and undo/redo, fix DetectedProject foreign key bug, implement cross-OS SSH port and firewall automation, introduce `gitmap ssh ui`, and build intelligent SSH connection troubleshooting.

---

## 1. Executive Summary

This specification governs seven architectural enhancements to GitMap:
1. **Foreign Key Constraint Resolution (RCA):** Fix `FOREIGN KEY constraint failed (787)` on `INSERT INTO DetectedProject` by ensuring `RepoId` and `ProjectTypeId` exist in their respective parent tables before project upserts, with fallback foreign key synchronization.
2. **Comprehensive Database Reset & Reseed Engine (`gitmap reset` / `gitmap db reset`):** Systematically purge all active SQLite database files across `.gitmap/`, report all removed database files with human-readable confirmations, and reseed clean schemas and default dictionaries.
3. **SSH Output Layout Alignment & Dedicated Key Viewer (`gitmap ssh view`):** Reorder `gitmap ssh` terminal presentation so SSH key status and public key info appear cleanly at the bottom, and introduce `gitmap ssh view` (`v`, `show`, `key`) that displays exclusively the SSH public key/details and copies it to the system clipboard without printing command help.
4. **SSH Key Creation Overwrite Guard, Backup & Undo/Redo (`gitmap ssh create`):** Detect existing keys and prompt for interactive user confirmation (`[y/N]`) unless `-y`/`--yes`/`--confirm` is provided; automatically create timestamped backups with task undo/redo registration; sync regeneration in `scripts-fixer`.
5. **Cross-OS SSH Port Management & Firewall Automation:** Enable and disable SSH on custom ports (single or multiple), configure host firewall rules across Windows (`netsh`/PowerShell), Ubuntu/Debian (`ufw`/`iptables`), CentOS/RHEL (`firewalld`), and macOS (`pfctl`), and provide `gitmap ssh enable-public <port>` to expose SSH externally with security routing.
6. **Interactive Browser SSH Management UI (`gitmap ssh ui` / `gitmap ssh web`):** Launch a local web server serving an intuitive interface for SSH keys, active ports, firewall rules, and fleet node health.
7. **Guided SSH Connection Diagnostics & 7-Step Troubleshooting Engine:** Intercept SSH connection errors and render structured, step-by-step remediation advice to guide users from ping testing to public key deployment.

---

## 2. Root Cause Analysis (RCA): DetectedProject Foreign Key Constraint (787)

### Symptom
During scan execution with project detection:
```text
[QueryWrapper Error]: exec failed: constraint failed: FOREIGN KEY constraint failed (787)
query: INSERT INTO DetectedProject
        (RepoId, ProjectTypeId, ProjectName, AbsolutePath, RepoPath, RelativePath, PrimaryIndicator)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(RepoId, ProjectTypeId, RelativePath) DO UPDATE SET ...
```

### Root Cause
1. In SQLite, `PRAGMA foreign_keys = ON;` enforces that `DetectedProject.RepoId` must match an existing `Repository.Id` and `DetectedProject.ProjectTypeId` must match an existing `ProjectType.Id`.
2. When scanning multi-repo workspaces or workspaces where some repositories fail clone validation, are excluded, or have not yet committed their parent `Repository` row within the same transaction/connection, the child `DetectedProject` insertion executes against an uncommitted or missing `RepoId`.
3. In addition, when dynamic project types (e.g. newly registered framework indicators) are detected, if the lookup table `ProjectType` has not been pre-seeded or upserted prior to project insertion, SQLite throws error 787.

### Architecture Fix
1. Pre-seed and verify all detected `ProjectType` slugs before scanning child files.
2. Ensure repository upsert commits the `Repository` record and retrieves the verified primary key before dispatching project detection.
3. Add defensive fallback: if `RepoId` is missing or uncommitted, verify existence via `SELECT Id FROM Repository WHERE Id = ?` or query by `Path`. If absent, upsert the repository stub or skip project attachment with structured error logging.

---

## 3. Database Reset & Reseed Engine (`gitmap reset`)

### Command Interface
- Primary: `gitmap reset`
- Aliases: `gitmap db reset`, `gitmap database reset`
- Flags:
  - `-y`, `--yes`: Confirm non-interactive database reset.
  - `--keep-cache`: Purge databases but retain search/AUM caches.
  - `--dry-run`: Preview database files that would be removed.

### Execution Workflow
1. Disconnect and close active database connection pools (`gitmap.db`, `installation.db`, `repodb/`, `templates.db`, `pipeline.db`).
2. Discover all `.db`, `.db-shm`, `.db-wal` files under `.gitmap/` and user cache directories.
3. Remove discovered database files safely.
4. Report deleted databases to the console:
   ```text
   ✔ Removed database: .gitmap/gitmap.db (1.2 MB)
   ✔ Removed database: .gitmap/installation.db (256 KB)
   ✔ Removed database: .gitmap/templates.db (64 KB)
   ✔ Total 3 database(s) purged.
   ```
5. Trigger automatic reseeding:
   - Run schema migrations (`InitSchema`).
   - Reseed core lookup tables (`ProjectType`, `RepoTag`, `DefaultSettings`).
   - Print confirmation: `✔ All databases reseeded with clean baseline schemas.`

---

## 4. SSH Command Enhancements & Dedicated Key Viewer

### 4.1 Output Layout Alignment
When `gitmap ssh` is executed without arguments:
- Display help categories and commands first.
- Move local SSH key status, fingerprint, and public key preview to the very bottom as the focal summary card.

### 4.2 `gitmap ssh view` (Key Viewer & Auto-Clipboard)
- Synopsis: `gitmap ssh view` (Aliases: `v`, `show`, `key`, `pubkey`)
- Actions:
  1. Locates active SSH public key (`~/.ssh/id_rsa.pub`, `~/.ssh/id_ed25519.pub`, or instance key).
  2. Copies public key content directly to OS clipboard (Windows `clip.exe` / PowerShell, Linux `xclip` / `wl-copy`, macOS `pbcopy`).
  3. Renders formatted terminal card:
     ```text
     ╔════════════════════════════════════════════════════════════════════════════════╗
     ║ SSH Public Key Information                                                     ║
     ╚════════════════════════════════════════════════════════════════════════════════╝
       Path:        %USERPROFILE%\.ssh\id_rsa.pub
       Algorithm:   RSA (3072-bit) / ED25519
       Fingerprint: SHA256:...
       Status:      Copied to clipboard! ✔

       [tip] To regenerate this key:
         gitmap ssh create [email] [-y]
     ```

### 4.3 `gitmap ssh create` Overwrite Guard & Undo/Redo
- Checks if target key already exists.
- If present and `-y` is NOT specified, prompts interactively:
  `SSH key already exists at ~/.ssh/id_rsa. Overwrite? [y/N]: `
- If user confirms:
  - Creates backup file: `~/.ssh/id_rsa.bak.<timestamp>`.
  - Registers task undo action in GitMap task journal.
  - Generates new key.
- If rejected: exits cleanly without modifying keys.

---

## 5. Cross-OS SSH Port Management & Firewall Automation

### Command Suite
- `gitmap ssh port ls`: List active listening SSH ports and firewall rules.
- `gitmap ssh port add <port>`: Open port in firewall and configure sshd (e.g. `2222`).
- `gitmap ssh port rm <port>`: Close port in firewall and remove from sshd.
- `gitmap ssh enable [--port <port>]`: Ensure SSH service is running and firewall allowed.
- `gitmap ssh disable`: Stop SSH service and close firewall ports.
- `gitmap ssh enable-public <port>`: Allow external inbound SSH traffic through firewall and bind to `0.0.0.0`.
- `gitmap firewall add --port <port> --proto <tcp|udp> --action <allow|deny>`
- `gitmap firewall app <allow|deny> --path <binary>`

### OS Implementation Adapters
- **Windows:** PowerShell `New-NetFirewallRule`, `Set-NetFirewallRule`, `Get-NetFirewallPortFilter`, `sshd_config` update.
- **Ubuntu/Debian:** `ufw allow <port>/tcp`, `ufw delete allow <port>/tcp`, `/etc/ssh/sshd_config.d/gitmap.conf`.
- **CentOS/RHEL:** `firewall-cmd --permanent --add-port=<port>/tcp && firewall-cmd --reload`.
- **macOS:** `pfctl` anchor rules and `System Preferences` remote login sharing.

---

## 6. Interactive SSH Management Web UI (`gitmap ssh ui`)

- Command: `gitmap ssh ui` (Aliases: `gitmap ssh web`, `gitmap ssh gui`)
- Starts embedded lightweight HTTP server (random or default port 7890).
- Automatically opens system browser to `http://127.0.0.1:7890`.
- UI Features:
  - **Key Management:** View public key, 1-click copy, generate new key with confirmation, view backups.
  - **Port & Firewall:** Visual toggles for active ports, public access switch, firewall rules table.
  - **Fleet Nodes:** Live status cards for enrolled nodes (`u1`, `w3`, etc.) with 1-click ping and test.

---

## 7. Guided SSH Troubleshooting Engine

When any SSH connection or dial failure occurs:
Render a structured, formatted diagnostic box:
```text
  ▲ SSH Connection Failed to a@node-u1:22
  ────────────────────────────────────────────────────────────────────────
  Recommended Troubleshooting Steps:
    1. Ensure GitMap & SSH are enabled on target:
       Run on target machine: gitmap ssh enable --port 22
    2. Check target network reachability:
       Run locally: gitmap ping node-u1
    3. Verify target username and credentials:
       Run locally: gitmap ssh pass show u1
    4. Attempt manual connection probe:
       Run locally: ssh a@node-u1
    5. Test authentication via GitMap interactive prompt:
       Run locally: gitmap ssh join a@node-u1 u1
    6. Deploy current machine's public SSH key:
       Run locally: gitmap ssh copy-id u1
    7. Inspect target firewall status:
       Run on target machine: gitmap ssh port ls
```
