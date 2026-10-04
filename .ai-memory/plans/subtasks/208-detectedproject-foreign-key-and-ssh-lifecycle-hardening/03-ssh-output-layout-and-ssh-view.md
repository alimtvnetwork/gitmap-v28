# Subtask 78.3: SSH Output Layout Reorder & Dedicated Key Viewer (`gitmap ssh view`)

## 1. Context & Objective
When the user executes `gitmap ssh`:
- The command reference help should appear FIRST.
- The SSH public key and fingerprint card should appear at the VERY END.
When the user executes `gitmap ssh view` (or `v`, `show`, `key`):
- Only the SSH key card, fingerprint, and public key are displayed.
- The 50-line command reference help is suppressed.
- The public key is automatically copied to the OS clipboard.
- A suggestion is printed: `To regenerate this key: gitmap ssh create -y (or pass --force)`.

---

## 2. Technical Scope
- `cli/cmdssh/ssh.go`
- `cli/cmdssh/sshcat.go`
- `cli/cmdssh/ssh_help_menu.go`
- `cli/cmdssh/ssh_fallback.go`

---

## 3. Remediation Checklist
- [ ] In `cli/cmdssh/ssh.go`, adjust execution flow so command help renders before public key summary.
- [ ] In `cli/cmdssh/sshcat.go`, support `gitmap ssh view` with formatted card, clipboard copy, and regeneration hint.
- [ ] Support `--raw` / `-r` flag to emit raw key text without headers for scripting.
