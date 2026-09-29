# Subtask 193.3: Subsystem Importers Dual-Format Ingestion
Traceability ID: Task-03
Spec Reference: [02-spec/21-app/183-typed-json-envelope-and-format-inspection.md](../../../02-spec/21-app/183-typed-json-envelope-and-format-inspection.md)
Target Files: cli/cmdssh/ssh_import.go, cli/cmdmacro/macro_export_import.go, cli/cmd/commitin/config_json.go, cli/cmdui/ui_server.go
Action: Update SSH nodes importer, macro importer, commit-in config parser, and UI settings loader to transparently accept both the new typed envelope (`attributes` + `data`) and legacy flat formats.
Acceptance Criteria:
1. SSH import parses both envelope and raw nodes JSON.
2. Macro import parses both envelope and raw macro array.
3. Commit-in config parses both envelope and raw config.
4. UI settings loader parses both envelope and raw settings.
Targeted Verification: `go test -v ./cmdssh/... ./cmdmacro/... ./cmd/commitin/... ./cmdui/...`
