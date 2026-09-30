# Subtask 01: Multi-Pass JSON Envelope Variables & Single-Arg `change-password` Mode
Traceability ID: Task-01, Task-02
Spec Reference: [02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md](../../../02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md)
Target Files: cli/jsonenvelope/variables.go, cli/jsonenvelope/envelope.go, cli/jsonenvelope/variables_test.go, cli/cmdos/change_password_cmd.go, cli/cmdos/change_password_cmd_test.go, cli/cmdos/os_help_modern.go
Action:
1. In `cli/jsonenvelope/variables.go`, update `MergeVariables` to resolve chained variable references (`${secretsDir}` inside `summaryPath`) across 2 passes, and update `ExpandVariables` to run 2 passes so nested variables always resolve completely.
2. In `cli/jsonenvelope/envelope.go`, ensure `ExtractPayload` also expands top-level `variables` / `workDirectory.variables` when `attributes` is present even if `data` is omitted (flat envelopes like `03-git-setup.json` and `06-vmpass.json`).
3. In `cli/cmdos/change_password_cmd.go`, update `parseChangePasswordArgs` and `assignPositionalArgs` so that:
   - `-u` / `--user <user>` and `-p` / `--password <pass>` flags are supported.
   - When 1 positional argument is given (`gitmap os change-password <new-password>`), if `opts.Username == ""`, treat `pos[0]` as `opts.Password` (so `resolvePasswordTarget` defaults `opts.Username` to the current OS user and prompts for confirmation).
   - When 2 positional arguments are given (`gitmap os change-password <user> <new-password>`), set `opts.Username = pos[0]` and `opts.Password = pos[1]`.
4. Update `cli/cmdos/change_password_cmd_test.go` and `cli/cmdos/os_help_modern.go` to test and showcase both single-argument current-user mode and two-argument admin mode.
Acceptance Criteria:
- Chained variables in `MergeVariables` and `ExpandVariables` resolve cleanly.
- `gitmap os change-password MyNewPass` sets `opts.Password = "MyNewPass"` and resolves `opts.Username` to the current user.
- `gitmap os change-password admin pass123` sets `opts.Username = "admin"` and `opts.Password = "pass123"`.
Targeted Verification: gofmt -l cli/jsonenvelope cli/cmdos
