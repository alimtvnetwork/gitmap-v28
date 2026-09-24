# Subtask 01: Pure SSH Stdin Streaming Protocol for Remote File Transfers
Traceability ID: Task-01 / Task-03
Spec Reference: [02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md](../../../02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md)
Target Files: [cli/cmdssh/ssh_install_exec.go, cli/cmdssh/ssh_stream.go, cli/cmdssh/ssh_transfer.go]
Action: Replace inlined Base64 command-line string execution with native pure SSH session StdinPipe streaming. Implement StreamFileToRemote with POSIX Tar extraction (tar.exe -xf - -C <dir> on Windows, tar -xf - -C <dir> on Unix) and direct stdin pipe fallbacks.
Acceptance Criteria:
- Payloads > 15MB upload smoothly without command-line buffer overflow or 0ms failure.
- Exact file permissions (0755) and file lengths preserved.
Targeted Verification: go test ./cli/cmdssh -run TestStream
