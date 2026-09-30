# Subtask 02: SSH Join Interactive Vault Samples & Clean Dev Lint Fix
Traceability ID: Task-03, Task-05
Spec Reference: [02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md](../../../02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md)
Target Files: cli/cmdssh/sshjoin_cmd.go, cli/cmd/clean_dev_entry.go
Action:
1. In `cli/cmdssh/sshjoin_cmd.go`, ensure the help text and examples explicitly showcase `gitmap ssh join {user}@ip <alias>` prompting interactively for a password on first-time join, encrypting and storing it via salted RSA (`EncryptSSHPassword`), reusing it on `gitmap ssh <alias>`, and skipping storage if left blank.
2. In `cli/cmd/clean_dev_entry.go`, wire `runDevToolTopLevel` to call `runDevTopLevel(args)` so `runDevTopLevel` and `isDevCleanActionVerb` are actively used and `golangci-lint` (`unused`) passes with zero warnings.
Acceptance Criteria:
- `gitmap ssh join --help` displays the interactive salted-RSA password prompt sample and explanation.
- `runDevTopLevel` and `isDevCleanActionVerb` in `cli/cmd/clean_dev_entry.go` have active callers.
Targeted Verification: gofmt -l cli/cmdssh cli/cmd
