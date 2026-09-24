# Plan 93: SSH Install-Exec Streaming Upload Protocol, DB Resolution Fallback & Dynamic Remote OS Probing

Spec Reference: [02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md](../../../02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md)  
Issue Reference: [02-spec/22-app-issues/42-ssh-install-exec-upload-failed-and-os-misclassification-rca.md](../../../02-spec/22-app-issues/42-ssh-install-exec-upload-failed-and-os-misclassification-rca.md)  

## 1. Architectural Context & Blast Radius
When running `gitmap ssh install-exec <setup.exe>`, previous code converted binaries into inlined Base64 command strings. Real-world installers (17.3 MB) produced 23 million character commands that exceeded OS buffer limits (32,767 chars on Windows) and failed instantly in 0ms with `✖ UPLOAD FAILED`.
Additionally, `ssh_hosts` table lacked an explicit `OS` column, causing merged connections to default to `"linux"`, preventing `w1` (Windows Server) from being targeted with Windows arguments or recognized under `--os win`.
Finally, `store.BinaryDataDir()` had no fallback to user AppData data directory when run from arbitrary folders (e.g. `Downloads`), reporting `No matching SSH machines`.

## 2. Actionable Deliverables & Subtask Mapping
- **Subtask 93.1**: Implement Pure SSH Stdin Streaming Protocol (`StreamFileToRemote`) via POSIX Tar stream (`tar.exe -xf - -C <dir>` on Windows, `tar -xf - -C <dir>` on Unix) with StdinPipe fallback, replacing Base64 command-line string inlining.
- **Subtask 93.2**: Implement Dynamic Remote OS Probing upon connection in `ssh_install_exec.go`, persist updated OS to DB, fix `cli/db/sshconnection.go` merge logic, and add AppData global DB fallback in `store/location.go`.
- **Subtask 93.3**: Surface full error messages in `renderInstallExecResultsTable`, author E2E integration tests in `cli/tests/e2e/`, and execute real installation test on live fleet nodes.
