# Subtask 77.6: Interactive SSH Web UI

- **Parent Plan:** [77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md](../../pending/77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md)
- **Spec Reference:** [02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md](../../../../02-spec/21-app/207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca/01-architecture-spec.md)
- **Status:** Pending
- **Target Area:** `cli/cmdssh/`, `cli/cmdui/`, web assets

## Objective
Implement the interactive browser-based SSH management web interface triggered by `gitmap ssh ui` (aliases: `web`, `dashboard`), which starts an embedded HTTP server, automatically launches the system browser to `http://localhost:8080/ssh`, and provides intuitive key viewing with 1-click clipboard copy, port and firewall management, and fleet node connectivity testing.

## Scope & Implementation Details

### 1. Command Invocation & Browser Auto-Launch (`cli/cmdssh/ssh_ui.go`)
- Command synonyms: `gitmap ssh ui`, `gitmap ssh web`, `gitmap ssh dashboard`.
- Pluggable invocation through `RunSSHUIFn("ssh", 8080)` hooked to `cmdui.RunUI`.
- Fallback browser launch using OS-native handlers:
  - Windows: `cmd /c start http://localhost:8080/ssh`
  - macOS: `open http://localhost:8080/ssh`
  - Linux: `xdg-open http://localhost:8080/ssh`
- Graceful degradation in headless environments (prints clickable URL without exiting with error).

### 2. Embedded Web Server & API Endpoints (`cli/cmdui/ui_server.go`)
- Port allocation: binds to preferred port 8080 (scans up to 50 subsequent ports if 8080 is in use).
- Serves SPA routing under `/ssh` rendering the embedded GitMap Web UI (`IndexHTML` in `cli/cmdui/ui_assets.go`).
- API Endpoints:
  - `GET /api/ssh/key`: Returns local public key information:
    - `publicKey`: formatted public key string.
    - `path`: absolute path to public key file.
    - `algorithm`: key type (RSA, ED25519, etc.).
    - `fingerprint`: SHA256 key fingerprint.
    - `backups`: list of available timestamped backup keys.
  - `POST /api/ssh/key/create`: Generates new SSH key pair with optional overwrite and backup.
  - `GET /api/ssh/ports`: Returns configured listening ports in `sshd_config` and active firewall status.
  - `POST /api/ssh/ports/add`: Adds new listening port and updates host firewall.
  - `POST /api/ssh/ports/rm`: Removes listening port and cleans firewall rule.
  - `POST /api/ssh/ports/public`: Enables WAN / public accessibility on target port.
  - `GET /api/ssh/firewall`: Returns current firewall rules relevant to SSH.
  - `GET /api/ssh/nodes`: Returns enrolled fleet nodes with status, IP, port, and OS.
  - `POST /api/ssh/nodes/test`: Triggers real-time ping/TCP probe to specified node.

### 3. Web UI Dashboard Interface (`cli/cmdui/ui_assets.go`)
- **Key Display & Clipboard Action:**
  - Card displaying the active public key in styled monospace block.
  - Prominent "Copy Key to Clipboard" button triggering native browser clipboard API with visual confirmation toast.
  - Key regeneration modal with overwrite warning and timestamped backup confirmation.
- **Port & Firewall Management Panel:**
  - Visual status badges for configured ports (`22`, `2222`, etc.) with green/red firewall indicators.
  - Modal form to add listening ports or toggle public access.
  - Real-time firewall rule list table showing rule name, port, direction, and action.
- **Fleet Node Health & Testing:**
  - Grid cards for enrolled nodes (`u1`, `w3`, etc.) displaying IP, SSH port, latency, and status.
  - 1-click test probe button triggering `/api/ssh/nodes/test` with latency badge update.

## Acceptance Criteria
- [ ] Running `gitmap ssh ui` starts embedded HTTP server and opens browser to `http://localhost:8080/ssh`.
- [ ] Running aliases `gitmap ssh web` and `gitmap ssh dashboard` perform identical UI launch.
- [ ] Web UI displays local SSH public key with working 1-click clipboard copy button.
- [ ] Web UI displays accurate listening port and firewall status queried from backend APIs.
- [ ] Port add/remove operations from Web UI correctly invoke backend firewall and `sshd_config` updates.
- [ ] Fleet node status cards reflect active node inventory with test probe execution.
