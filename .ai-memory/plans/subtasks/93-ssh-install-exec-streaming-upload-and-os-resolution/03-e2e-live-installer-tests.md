# Subtask 03: Live Fleet E2E Verification & Error Surfacing
Traceability ID: Task-04 / Task-05
Spec Reference: [02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md](../../../02-spec/21-app/155-ssh-install-exec-streaming-upload-and-os-resolution.md)
Target Files: [cli/cmdssh/ssh_install_exec.go, cli/tests/e2e/ssh_exec_os_filter_and_install_exec_tempe2e_test.go]
Action: Surface r.Error in renderInstallExecResultsTable when status indicates failure, author automated tests in cli/tests/e2e/, and execute live end-to-end installation test with real binary on fleet nodes.
Acceptance Criteria:
- Table outputs descriptive error message on upload/install failures.
- Real 17.3MB installer uploads successfully to online node and executes properly.
Targeted Verification: go test -v -tags=tempe2e ./cli/tests/e2e -run TestTempE2E_SSHInstallExec
