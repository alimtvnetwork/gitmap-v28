# Subtask 193.2: Format Inspection CLI Command (`gitmap which-format`)
Traceability ID: Task-02
Spec Reference: [02-spec/21-app/183-typed-json-envelope-and-format-inspection.md](../../../02-spec/21-app/183-typed-json-envelope-and-format-inspection.md)
Target Files: cli/cmd/which_format_cmd.go, cli/cmd/which_format_cmd_test.go, cli/cmd/root.go, cli/cmd/roottooling.go
Action: Implement `gitmap which-format [files...]` and directory auto-scanner. Inspect JSON files, identify matched format, display suggested import command and system impact summary. Flag unmatched JSONs. Provide single-line batch import command and `-y` flag advice.
Acceptance Criteria:
1. `gitmap which-format` scans folder when no files passed.
2. Identifies known GitMap formats (`ssh-nodes`, `macro`, `commit-pull-config`, etc.).
3. Emits suggested import commands and system impact for matched files, and explains non-match for unmatched files.
4. Outputs copy-pasteable single-line batch command with `-y` advice.
Targeted Verification: `go test -v ./cmd -run TestWhichFormat`
