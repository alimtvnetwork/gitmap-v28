# Issue Domain 03: SSH and Remote Connection Errors

- **Domain:** Multi-Node Network Probing and SSH Daemon Lifecycle
- **Status:** Consolidated Problem & Resolution Matrix

## 1. Firewall ICMP Drops with Active SSH Port
- **Symptoms:** `gitmap nodes ping` falsely reported online machines as OFFLINE on Windows targets.
- **Root Cause:** Windows Defender Firewall dropped ICMP echo requests while allowing TCP port 22.
- **Resolution:** Implemented dual-stack probing: when ICMP fails, fallback to concurrent TCP handshake on port 22, classifying as `REACHABLE (TCP)`.

## 2. Password Capture Echo Leaks
- **Symptoms:** Passwords typed into terminal were echoed in plaintext on certain Linux terminals.
- **Root Cause:** Standard `fmt.Scanln` used instead of POSIX termios / Win32 raw terminal mode.
- **Resolution:** Integrated `golang.org/x/term.ReadPassword` with character bullet masking (`*`).
