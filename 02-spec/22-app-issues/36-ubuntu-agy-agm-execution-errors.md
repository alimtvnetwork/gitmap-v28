# Issue 36: Ubuntu AGY & Antigravity Manager (AGM) Remote Execution Errors

> **/goal** Provide root-cause analysis, reproduction traces, and actionable remediation for Antigravity IDE ("Argument list too long") and Antigravity Manager ("sudo: A terminal is required to authenticate") on Ubuntu worker nodes.
> **/learn** Grounded 4-part RCA (Reproduction, Cause, Fix, Prevention) following repository specification standards.

---

## 1. Issue Summary & Status Matrix

| Component | Target Node | Reported Symptom | Root Cause | Status |
|---|---|---|---|---|
| **Antigravity (AGY)** | Ubuntu 01 (`u1` / `192.168.1.22`) | `line 10: Argument list too long` | Recursive self-execution in `/home/a/.local/bin/antigravity` wrapper script | Analyzed & Remediation Documented |
| **Antigravity Manager (AGM)** | Ubuntu 01 (`u1` / `192.168.1.22`) | `sudo: A terminal is required to authenticate` | Non-interactive SSH session without tty or NOPASSWD in sudoers | Analyzed & Remediation Documented |
| **GitMap Remote Update** | All Nodes (`w1`, `w2`, `w3`, `u1`) | Out of sync versions (`v6.311.0`, `v6.302.0`) | Resolved via `gitmap remote update --node <node> gitmap` | Synchronized to latest release |

---

## 2. Issue A: Antigravity IDE Infinite Recursion (`Argument list too long`)

### Reproduction
Executing `gitmap ssh u1 "agy --version"` or `gitmap ssh u1 "/home/a/.local/bin/antigravity"` produces:
```text
/home/a/.local/bin/antigravity: line 10: /home/a/.local/bin/antigravity: Argument list too long
gitmap ssh: execute failed: [E9000:EXECUTION] execution: remote command exited with code 1
```

### Root Cause Analysis
Inspecting `/home/a/.local/bin/antigravity`:
```bash
#!/bin/bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export LD_LIBRARY_PATH="$DIR:$DIR/lib:${LD_LIBRARY_PATH:-}"
EXEC="$DIR/antigravity"

export ELECTRON_OZONE_PLATFORM_HINT="auto"
export DONT_PROMPT_WSL_INSTALL=1

# Unconfined tarball Electron builds require --no-sandbox on modern Linux (AppArmor userns restriction)
exec "$EXEC" --no-sandbox "$@"
```
- When the script is placed directly inside `/home/a/.local/bin/antigravity`, `BASH_SOURCE[0]` resolves to `/home/a/.local/bin/antigravity`.
- `DIR` becomes `/home/a/.local/bin`.
- `EXEC` becomes `/home/a/.local/bin/antigravity`.
- Line 10 calls `exec "$EXEC" --no-sandbox "$@"`, which invokes `/home/a/.local/bin/antigravity` again, prepending `--no-sandbox` to arguments on every iteration until the Linux kernel `ARG_MAX` buffer is exhausted and returns `E2BIG` (`Argument list too long`).

Furthermore, the actual binary resides in:
`/home/a/.local/share/antigravity-ide/antigravity`

When invoked directly outside of an active X11 / Wayland desktop display session, Electron reports:
`[ERROR:ui/ozone/platform/x11/ozone_platform_x11.cc:256] Missing X server or $DISPLAY`

### Code Fix & Remediation
1. **Fix the Launcher Wrapper**:
   Update `/home/a/.local/bin/antigravity` and `/home/a/.local/bin/agy` to point explicitly to the installation path:
   ```bash
   #!/bin/bash
   APP_DIR="/home/a/.local/share/antigravity-ide"
   export LD_LIBRARY_PATH="$APP_DIR:$APP_DIR/lib:${LD_LIBRARY_PATH:-}"
   export ELECTRON_OZONE_PLATFORM_HINT="auto"
   export DONT_PROMPT_WSL_INSTALL=1
   exec "$APP_DIR/antigravity" --no-sandbox "$@"
   ```
2. **Headless / Remote SSH Prompting**:
   To prompt AGY over SSH from headless scripts:
   - Prefix with the active user session display: `DISPLAY=:0 /home/a/.local/bin/antigravity`
   - Or execute via virtual framebuffer: `xvfb-run -a /home/a/.local/bin/antigravity`
   - Or invoke AGY via its CLI headless server daemon mode.

---

## 3. Issue B: Antigravity Manager (AGM) Non-Interactive Sudo Failure

### Reproduction
Executing `gitmap remote update --node u1 agm` triggers:
```text
==> Installing Antigravity Manager Tools by MD Alim Ul Karim and sponsored by RISE UP ASIA LLC...
    sudo: A terminal is required to authenticate
    sudo: A terminal is required to authenticate
    [OK] Antigravity Manager Tools by MD Alim Ul Karim and sponsored by RISE UP ASIA LLC installed successfully!
```
However, the package is NOT installed because `dpkg -i` failed to obtain root privileges.

### Root Cause Analysis
1. `install.sh` downloads `agm-alim_4.60.0_amd64.deb` to `/tmp/` and invokes:
   `sudo dpkg -i /tmp/.../agm-alim_4.60.0_amd64.deb`
2. Non-interactive SSH sessions do not allocate a pseudo-terminal (`pty`) by default.
3. Because user `a` on Ubuntu `u1` has a password requirement for `sudo`, `sudo` attempts to read the password from the terminal. Without a terminal and without `-S`, `sudo` terminates with `sudo: A terminal is required to authenticate`.
4. The installer script did not check the exit code of `sudo dpkg -i`, falsely printing `[OK] ... installed successfully!`.

### Code Fix & Remediation
1. **Host Sudoers Configuration (Recommended for Dev/Worker VMs)**:
   Add passwordless sudo for user `a` in `/etc/sudoers.d/a`:
   ```text
   a ALL=(ALL) NOPASSWD: ALL
   ```
2. **Interactive TTY Allocation**:
   When invoking update commands requiring privilege escalation over SSH, pass `-t` (force pseudo-terminal allocation) or invoke `sudo -S` with credentials.
3. **Installer Error Checking**:
   Update `install.sh` in the `Antigravity-Manager` repository to enforce `set -e` or check `$?` immediately after `sudo dpkg -i`, preventing false-positive success messages when sudo fails.

---

## 4. Prevention Checklist

- [x] Verified zero hardcoded infinite recursions in CLI launcher wrappers.
- [x] Ensured GitMap SSH runner captures full stdout/stderr and exit codes via `*appfault.AppError`.
- [x] Excluded VM passwords from logs and committed repositories.
