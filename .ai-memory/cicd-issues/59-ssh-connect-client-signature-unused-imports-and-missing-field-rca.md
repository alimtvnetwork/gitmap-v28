# 59 — SSH Connect Client Signature, Unused Imports & db.SSHConnection Field RCA

- **Date:** 2026-09-19
- **Status:** Resolved
- **Impact:** CI/CD Build Failures (History Rewrite Smoke, Race Detector, Cross-Platform Build)

---

### 1. Symptom

Compilation failed across 5 GitHub Actions workflow sections during CI build on commit `f5fcf7a`:
```text
cmdssh/ssh_agy_cmd.go:71:45: too many arguments in call to connectSSHClient
    have (db.SSHConnection, string)
    want (db.SSHConnection)
cmdssh/ssh_code_remote.go:15:45: too many arguments in call to connectSSHClient
    have (db.SSHConnection, string)
    want (db.SSHConnection)
cmdssh/ssh_install_remote.go:67:45: too many arguments in call to connectSSHClient
    have (db.SSHConnection, string)
    want (db.SSHConnection)
cmdssh/ssh_update_remote.go:41:45: too many arguments in call to connectSSHClient
    have (db.SSHConnection, string)
    want (db.SSHConnection)
cmdssh/sshexec.go:4:2: "context" imported and not used
cmdssh/sshexec.go:10:2: "github.com/charmbracelet/lipgloss" imported and not used
cmdssh/sshexec.go:160:31: c.ID undefined (type db.SSHConnection has no field or method ID)
cmdssh/ssh_transfer.go:17:2: "github.com/alimtvnetwork/gitmap-v28/cli/termpad" imported and not used
```

---

### 2. Root Cause

1. `connectSSHClient` signature in `cli/cmdssh/sshexec.go` was tightened to `connectSSHClient(c db.SSHConnection)` without maintaining backward compatibility for existing callers in `ssh_agy_cmd.go`, `ssh_code_remote.go`, `ssh_install_remote.go`, and `ssh_update_remote.go` that pass a diagnostic `header string`.
2. Unused imports `"context"` and `"github.com/charmbracelet/lipgloss"` in `cli/cmdssh/sshexec.go` and `"github.com/alimtvnetwork/gitmap-v28/cli/termpad"` in `cli/cmdssh/ssh_transfer.go` remained after UI extraction.
3. `isConnExcluded` in `cli/cmdssh/sshexec.go` attempted to read `c.ID`, but `db.SSHConnection` uses `Alias` as its primary key and does not declare an `ID` field.

---

### 3. Resolution

1. **Variadic Header Signature:** Updated `connectSSHClient` signature in `cli/cmdssh/sshexec.go` to `connectSSHClient(c db.SSHConnection, headers ...string) (*ssh.Client, bool)`, supporting both 1-argument and 2-argument call sites seamlessly.
2. **Removed Unused Imports:**
   - Stripped `"context"` and `"github.com/charmbracelet/lipgloss"` from `cli/cmdssh/sshexec.go`.
   - Stripped `"github.com/alimtvnetwork/gitmap-v28/cli/termpad"` from `cli/cmdssh/ssh_transfer.go`.
3. **Removed Non-existent Field Reference:** Removed `idStr := fmt.Sprintf("%d", c.ID)` from `isConnExcluded` in `cli/cmdssh/sshexec.go`, preserving exclusion matching on `c.Alias`, `c.IPAddress`, and `userHost` (`Username@IPAddress`).

---

### 4. Prevention & Learnings

- When modifying helper signatures shared across sibling files in the same package (`cmdssh`), use variadic optional parameters (`headers ...string`) to preserve API compatibility across all callers.
- Cross-reference database model definitions in `cli/db/` (`SSHConnection` uses `Alias` as primary key) before accessing struct fields in command filters.
- Run `gofmt -l` and targeted file-level linter checks to detect unused imports prior to committing changes.
