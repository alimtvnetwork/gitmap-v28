# Subtask 03: Harmonize `01-gitmap` JSON Envelopes (`03-`, `06-`, `07-`, `08-`)
Traceability ID: Task-04
Spec Reference: [02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md](../../../02-spec/21-app/189-envelope-variables-os-password-and-secrets-sequence.md)
Target Files: D:\work\repo-secrets\01-gitmap\03-git-setup.json, D:\work\repo-secrets\01-gitmap\06-vmpass.json, D:\work\repo-secrets\01-gitmap\07-seo-templates.json, D:\work\repo-secrets\01-gitmap\08-seo-templates-cli.json, D:\work\repo-secrets\vmpass.json
Action:
1. Update `01-gitmap/07-seo-templates.json` to have a complete `version: "2.0"` `attributes` block with structured `workDirectory` object (`path: "${workDir}"`, `defaultPath: "D:\\work\\gitmap"`, `isApplied: true`, `isEnforced: false`, `variables: {"workDir": "D:\\work\\gitmap"}`), cross-platform `notes`, and top-level `variables`.
2. Update `01-gitmap/08-seo-templates-cli.json` to include a complete `version: "2.0"` `attributes` block (with structured `workDirectory` object and cross-platform `notes`), top-level `variables`, and `data` containing `titles` and `descriptions` (while preserving root-level `titles` and `descriptions` for direct readers).
3. Update `01-gitmap/03-git-setup.json`, `01-gitmap/06-vmpass.json`, and root `vmpass.json` so they include a `"data"` wrapper alongside root keys for full `IsEnvelope` compatibility.
Acceptance Criteria:
- All JSON files in `01-gitmap/` (`01-` through `08-`) have `version: "2.0"`, structured `workDirectory` objects, `variables`, and `data` blocks.
- `gitmap which-format` and `ConvertFrom-Json` parse all JSON files in `01-gitmap/` without error.
Targeted Verification: pwsh -NoProfile -Command "Get-ChildItem D:\work\repo-secrets\01-gitmap\*.json | ForEach-Object { Get-Content $_.FullName -Raw | ConvertFrom-Json | Out-Null }"
