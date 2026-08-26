import os
import hashlib
import re
import json

PLAN_SLUG = "gitmap-installer"
TOTAL_STEPS = 40
PLAN_FILE = f".lovable/plans/pending/02-{PLAN_SLUG}.md"
SUBTASK_DIR = f".lovable/plans/subtasks/{PLAN_SLUG}"
os.makedirs(SUBTASK_DIR, exist_ok=True)

tasks = [
    {"domain": "Contract", "phase": "Scaffold", "title": "Define Installer Models", "files": ["gitmap/model/installer.go", "gitmap/model/installer_test.go"], "symbol": "InstallerScript, InstallerVersion", "cmd": "go test ./model -run TestInstallerModel"},
    {"domain": "Plugin", "phase": "Implement", "title": "Add SQLite Migrations for Installers", "files": ["gitmap/store/migrations.go"], "symbol": "installer_scripts table, installer_versions table", "cmd": "go test ./store -run TestMigrations"},
    {"domain": "Plugin", "phase": "Implement", "title": "Store CreateInstaller", "files": ["gitmap/store/installer_create.go", "gitmap/store/installer_create_test.go"], "symbol": "CreateInstaller(script *model.InstallerScript)", "cmd": "go test ./store -run TestCreateInstaller"},
    {"domain": "Plugin", "phase": "Implement", "title": "Store GetInstallerBySlug", "files": ["gitmap/store/installer_get.go", "gitmap/store/installer_get_test.go"], "symbol": "GetInstallerBySlug(slug string)", "cmd": "go test ./store -run TestGetInstaller"},
    {"domain": "Plugin", "phase": "Implement", "title": "Store SaveVersion", "files": ["gitmap/store/installer_version.go", "gitmap/store/installer_version_test.go"], "symbol": "SaveVersion(version *model.InstallerVersion)", "cmd": "go test ./store -run TestSaveVersion"},
    {"domain": "Plugin", "phase": "Implement", "title": "Store ListInstallers", "files": ["gitmap/store/installer_list.go", "gitmap/store/installer_list_test.go"], "symbol": "ListInstallers() ([]model.InstallerScript, error)", "cmd": "go test ./store -run TestListInstallers"},
    {"domain": "Plugin", "phase": "Implement", "title": "Store ResetInstallers", "files": ["gitmap/store/installer_reset.go", "gitmap/store/installer_reset_test.go"], "symbol": "ResetInstallers(slug string, all bool)", "cmd": "go test ./store -run TestResetInstallers"},
    {"domain": "Plugin", "phase": "Implement", "title": "Store DeleteInstaller", "files": ["gitmap/store/installer_delete.go", "gitmap/store/installer_delete_test.go"], "symbol": "DeleteInstaller(slug string)", "cmd": "go test ./store -run TestDeleteInstaller"},
    {"domain": "Cli", "phase": "Scaffold", "title": "Root Installer Command", "files": ["gitmap/cmd/installer.go", "gitmap/cmd/installer_test.go"], "symbol": "installerCmd", "cmd": "go test ./cmd -run TestInstallerCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer Create Command", "files": ["gitmap/cmd/installer_create.go", "gitmap/cmd/installer_create_test.go"], "symbol": "installerCreateCmd, parseCreateFlags", "cmd": "go test ./cmd -run TestInstallerCreateCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer Update Command", "files": ["gitmap/cmd/installer_update.go", "gitmap/cmd/installer_update_test.go"], "symbol": "installerUpdateCmd", "cmd": "go test ./cmd -run TestInstallerUpdateCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer Update Win Command", "files": ["gitmap/cmd/installer_update_win.go", "gitmap/cmd/installer_update_win_test.go"], "symbol": "installerUpdateWinCmd", "cmd": "go test ./cmd -run TestInstallerUpdateWinCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer Install Win Command", "files": ["gitmap/cmd/installer_install_win.go", "gitmap/cmd/installer_install_win_test.go"], "symbol": "installerInstallWinCmd", "cmd": "go test ./cmd -run TestInstallerInstallWinCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer Export Commands", "files": ["gitmap/cmd/installer_export.go", "gitmap/cmd/installer_export_test.go"], "symbol": "installerExportCmd, installerExportAllCmd", "cmd": "go test ./cmd -run TestInstallerExportCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer Import Command", "files": ["gitmap/cmd/installer_import.go", "gitmap/cmd/installer_import_test.go"], "symbol": "installerImportCmd", "cmd": "go test ./cmd -run TestInstallerImportCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer Reset Command", "files": ["gitmap/cmd/installer_reset.go", "gitmap/cmd/installer_reset_test.go"], "symbol": "installerResetCmd", "cmd": "go test ./cmd -run TestInstallerResetCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer Revert Commands", "files": ["gitmap/cmd/installer_revert.go", "gitmap/cmd/installer_revert_test.go"], "symbol": "installerUndoCmd, installerRedoCmd, installerRevertCmd", "cmd": "go test ./cmd -run TestInstallerRevertCmd"},
    {"domain": "Cli", "phase": "Implement", "title": "Installer List Command", "files": ["gitmap/cmd/installer_ls.go", "gitmap/cmd/installer_ls_test.go"], "symbol": "installerLsCmd", "cmd": "go test ./cmd -run TestInstallerLsCmd"},
    {"domain": "Contract", "phase": "Implement", "title": "Path Normalization Utilities", "files": ["gitmap/fsutil/path_normalize.go", "gitmap/fsutil/path_normalize_test.go"], "symbol": "NormalizeToForwardSlashes(path), MakeRelativeToRoot(base, path)", "cmd": "go test ./fsutil -run TestPathNormalize"},
    {"domain": "Contract", "phase": "Scaffold", "title": "OS Targets Constants", "files": ["gitmap/constants/os_targets.go"], "symbol": "OSWindows, OSUbuntu, OSCentOS, OSDebian, OSArch", "cmd": "go build ./constants"},
    {"domain": "Plugin", "phase": "Scaffold", "title": "Manager Struct", "files": ["gitmap/installer/manager.go", "gitmap/installer/manager_test.go"], "symbol": "Manager, NewManager(db)", "cmd": "go test ./installer -run TestManagerStruct"},
    {"domain": "Plugin", "phase": "Implement", "title": "Manager Create Logic", "files": ["gitmap/installer/create.go", "gitmap/installer/create_test.go"], "symbol": "(m *Manager) Create(name, desc)", "cmd": "go test ./installer -run TestManagerCreate"},
    {"domain": "Plugin", "phase": "Implement", "title": "Manager Update Logic", "files": ["gitmap/installer/update.go", "gitmap/installer/update_test.go"], "symbol": "(m *Manager) Update(slug, osTarget)", "cmd": "go test ./installer -run TestManagerUpdate"},
    {"domain": "Contract", "phase": "Implement", "title": "Semantic Versioning Helper", "files": ["gitmap/installer/versioning.go", "gitmap/installer/versioning_test.go"], "symbol": "NextSemanticVersion(current)", "cmd": "go test ./installer -run TestSemanticVersioning"},
    {"domain": "Plugin", "phase": "Implement", "title": "Manager Undo Logic", "files": ["gitmap/installer/revert.go", "gitmap/installer/revert_test.go"], "symbol": "(m *Manager) Undo(slug)", "cmd": "go test ./installer -run TestManagerUndo"},
    {"domain": "Plugin", "phase": "Implement", "title": "Manager Redo Logic", "files": ["gitmap/installer/redo.go", "gitmap/installer/redo_test.go"], "symbol": "(m *Manager) Redo(slug)", "cmd": "go test ./installer -run TestManagerRedo"},
    {"domain": "Plugin", "phase": "Implement", "title": "Manager Exact Revert Logic", "files": ["gitmap/installer/revert_exact.go", "gitmap/installer/revert_exact_test.go"], "symbol": "(m *Manager) RevertTo(slug, version)", "cmd": "go test ./installer -run TestManagerRevertExact"},
    {"domain": "Plugin", "phase": "Implement", "title": "Export Single ZIP", "files": ["gitmap/installer/export.go", "gitmap/installer/export_test.go"], "symbol": "(m *Manager) ExportToZip(slug, path)", "cmd": "go test ./installer -run TestExportToZip"},
    {"domain": "Plugin", "phase": "Implement", "title": "Export All Installers", "files": ["gitmap/installer/export_all.go", "gitmap/installer/export_all_test.go"], "symbol": "(m *Manager) ExportAllToZip(path)", "cmd": "go test ./installer -run TestExportAllToZip"},
    {"domain": "Plugin", "phase": "Implement", "title": "Export Global State", "files": ["gitmap/installer/export_global.go", "gitmap/installer/export_global_test.go"], "symbol": "(m *Manager) ExportGlobalState(path)", "cmd": "go test ./installer -run TestExportGlobalState"},
    {"domain": "Plugin", "phase": "Implement", "title": "Import ZIP Dispatcher", "files": ["gitmap/installer/import.go", "gitmap/installer/import_test.go"], "symbol": "(m *Manager) ImportFromZip(path)", "cmd": "go test ./installer -run TestImportFromZip"},
    {"domain": "Plugin", "phase": "Implement", "title": "Import JSON Content", "files": ["gitmap/installer/import_json.go", "gitmap/installer/import_json_test.go"], "symbol": "(m *Manager) ImportFromJson(jsonStr)", "cmd": "go test ./installer -run TestImportFromJson"},
    {"domain": "Plugin", "phase": "Implement", "title": "Import Global State", "files": ["gitmap/installer/import_global.go", "gitmap/installer/import_global_test.go"], "symbol": "(m *Manager) ImportGlobalState(path)", "cmd": "go test ./installer -run TestImportGlobalState"},
    {"domain": "Plugin", "phase": "Implement", "title": "Import Conflict Resolver", "files": ["gitmap/installer/conflict.go", "gitmap/installer/conflict_test.go"], "symbol": "(m *Manager) resolveImportConflict(slug)", "cmd": "go test ./installer -run TestConflictResolver"},
    {"domain": "Plugin", "phase": "Implement", "title": "Installer Execution Dispatch", "files": ["gitmap/installer/execute.go", "gitmap/installer/execute_test.go"], "symbol": "(m *Manager) Execute(slug, osTarget)", "cmd": "go test ./installer -run TestExecuteDispatch"},
    {"domain": "Plugin", "phase": "Implement", "title": "Instruction Parser", "files": ["gitmap/installer/parse_instructions.go", "gitmap/installer/parse_instructions_test.go"], "symbol": "ParseInstructions(jsonBlob)", "cmd": "go test ./installer -run TestParseInstructions"},
    {"domain": "Plugin", "phase": "Implement", "title": "OS Command Runner", "files": ["gitmap/installer/runner.go", "gitmap/installer/runner_test.go"], "symbol": "RunInstallerCommand(cmd)", "cmd": "go test ./installer -run TestRunner"},
    {"domain": "Plugin", "phase": "Implement", "title": "Format Installer Table", "files": ["gitmap/installer/ui_ls.go", "gitmap/installer/ui_ls_test.go"], "symbol": "FormatInstallerTable(list)", "cmd": "go test ./installer -run TestFormatTable"},
    {"domain": "Plugin", "phase": "Implement", "title": "Detailed Help Printer", "files": ["gitmap/installer/ui_help.go", "gitmap/installer/ui_help_test.go"], "symbol": "PrintDetailedHelp()", "cmd": "go test ./installer -run TestDetailedHelp"},
    {"domain": "Plugin", "phase": "Implement", "title": "Composer Example Printer", "files": ["gitmap/installer/ui_examples.go", "gitmap/installer/ui_examples_test.go"], "symbol": "PrintComposerExample()", "cmd": "go test ./installer -run TestComposerExample"}
]

plan_content = f"""---
status: pending
---
# Gitmap Installer Commands & Scripts Implementation

## Context
Goal: Implement the gitmap installer create/export/import commands and sqlite persistence (40 steps).
Reference spec: `.lovable/spec/commands/01-gitmap-installer.md`.

Execution model:
One step per run. Exactly one step is executed per run. Self-loop after Verify.
At most 2 spawned agents, max 3 threads.

## Tasks
"""
for i, t in enumerate(tasks, 1):
    plan_content += f"- [ ] Task {i:03d} — {t['title']}\n"

with open(PLAN_FILE, "w", encoding="utf-8") as f:
    f.write(plan_content)

for i, t in enumerate(tasks, 1):
    task_num = f"{i:03d}"
    dep = f"{i-1:03d}" if i > 1 else "none"
    files_str = json.dumps(t['files'])
    
    task_content = f"""---
plan: .lovable/plans/pending/02-{PLAN_SLUG}.md
domain: {t['domain']}
phase: {t['phase']}
target_files: {files_str}
depends_on: ["Task {dep}"]
citations:
  app_spec: ".lovable/spec/commands/01-gitmap-installer.md §Core Requirements"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  language_guideline: "spec/02-coding-guidelines/03-golang/00-overview.md"
  boolean_styling: "spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md"
  folder_naming: "spec/02-coding-guidelines/08-file-folder-naming/03-golang.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  error_codes: "spec/21-app/07-error-and-logging/01-error-code-allocation.md"
  logging_traces: "spec/21-app/07-error-and-logging/02-logging-and-stack-traces.md"
  response_envelope: "spec/21-app/07-error-and-logging/03-response-envelope.md"
  golden_fixture: "n/a — no wire format in this step"
  strictly_avoid: ".lovable/strictly-avoid.md"
  database: "spec/04-database-conventions/01-sqlite-schema.md"
  ui_surface: "n/a — no direct UI in this step"
  tests: "unit {t['cmd'].split(' ')[-1]}"
  ci_cd_guard: "linter-scripts/check-go-build"
  ambiguity: "n/a — no ambiguity filed"
  issue_rca: "n/a — not a bugfix"
---
# Task {task_num} — {t['title']}

## 1. Learn
- Read `.lovable/spec/commands/01-gitmap-installer.md` to understand the overarching requirement for {t['symbol']}.
- Read `spec/02-coding-guidelines/00-canonical-size-tier.md` to ensure `{t['files'][0]}` remains concisely sized.
- Review `spec/03-error-manage/02-error-architecture/00-overview.md` for proper `apperror` context wrapping.
- Inspect `{t['files'][0]}` dependencies to see how {t['symbol']} interacts with its callers.

## 2. Goal
The objective is to implement `{t['symbol']}` natively in `{t['files'][0]}`. This explicitly unblocks downstream operations dependent on `{t['title']}` in the {t['domain']} domain. No other files should be manipulated.

## 3. Inputs and Contracts
- Exported Symbols: `{t['symbol']}`
- Package: `{t['files'][0].split('/')[-2]}`
- Error wrapping MUST use the `E_INSTALLER_*` code family.

## 4. Execute
1. Open `{t['files'][0]}`.
2. Implement the required structure, type, or function for `{t['symbol']}`.
3. Write unit tests for success and failure boundaries in `{t['files'][-1]}`.
4. Ensure no cross-domain pollution.

## 5. Constraints
- Must adhere strictly to `spec/02-coding-guidelines/00-canonical-size-tier.md` (keep logic segmented under 60 lines).
- Error wrapping must include stack traces.

## 6. Verify
```bash
{t['cmd']}
```
Expected output: The test suite passes cleanly with no panics.

## 7. Done When
- [ ] `{t['symbol']}` is successfully mapped and tested in `{t['files'][0]}`.
- [ ] All CI and `go test` commands exit zero.
- [ ] No hardcoded or dummy assumptions are left in the code.

## 8. Notes and Open Questions
None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.
"""
    with open(f"{SUBTASK_DIR}/{task_num}-task.md", "w", encoding="utf-8") as f:
        f.write(task_content)

print(f"Generated {TOTAL_STEPS} strict, unique tasks for {PLAN_SLUG}.")
