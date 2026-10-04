# Subtask 78.4: `gitmap ssh create` Overwrite Guard, Backup & Undo/Redo

## 1. Context & Objective
When the user executes `gitmap ssh create`:
- If an SSH key already exists, prompt `[y/N]` before overwriting.
- If the user provides `-y`, `--yes`, `--confirm`, or `-f`, bypass the prompt.
- Create a timestamped backup of the existing key pair (`id_rsa.bak.<unix-timestamp>`).
- Enqueue a task record in `TaskHistory` in SQLite database to support undo and redo.
- Suppress the 50-line command reference help menu during `gitmap ssh create`.

---

## 2. Technical Scope
- `cli/cmdssh/sshgen.go`
- `cli/cmdssh/sshexisting.go`
- `cli/cmdssh/ssh_history_tasks.go`

---

## 3. Remediation Checklist
- [ ] Implement `backupKeyWithTimestamp` in `sshexisting.go`.
- [ ] Implement interactive confirmation prompt in `sshgen.go`.
- [ ] Connect `RecordSSHKeyBackupTask` to write into SQLite `TaskHistory`.
- [ ] Suppress help menu from `dispatchCreateSSH`.
