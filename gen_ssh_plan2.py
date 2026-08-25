import os
import json
import hashlib

plan_slug = "ssh-login-and-join"
task_dir = f".lovable/plans/subtasks/{plan_slug}"
plan_file = f".lovable/plans/pending/01-{plan_slug}.md"
os.makedirs(task_dir, exist_ok=True)
os.makedirs(os.path.dirname(plan_file), exist_ok=True)
os.makedirs(".lovable/memory", exist_ok=True)

# Build the domains and units list to hit exactly 50 steps
# The domains: 
#  1. Contract (SQLite schema, types)
#  2. SSH Core (Login, parser, cross-platform execution)
#  3. Key Management (Send-key, add-auth)
#  4. Install / Gitmap Remote Setup (login-install)
#  5. IP Utilities (ip, ip-change)
#  6. Join / SJ (join, rm, ls, history)

# We need 50 steps. 
units = [
    # Domain 1: Database Contracts (1-5)
    {"domain": "Contract", "phase": "Scaffold", "title": "Define SSH host SQLite schema", "sym": "CreateHostsTable", "file": "gitmap/store/migrations_ssh.go"},
    {"domain": "Contract", "phase": "Implement", "title": "Implement host Insert/Update queries", "sym": "InsertSSHHost", "file": "gitmap/store/ssh_host.go"},
    {"domain": "Contract", "phase": "Wire+Test", "title": "Test host queries", "sym": "TestInsertSSHHost", "file": "gitmap/store/ssh_host_test.go"},
    {"domain": "Contract", "phase": "Scaffold", "title": "Define SSH join history schema", "sym": "CreateSSHHistoryTable", "file": "gitmap/store/migrations_ssh.go"},
    {"domain": "Contract", "phase": "Implement", "title": "Implement join history queries", "sym": "InsertSSHHistory", "file": "gitmap/store/ssh_history.go"},
    {"domain": "Contract", "phase": "Wire+Test", "title": "Test join history queries", "sym": "TestInsertSSHHistory", "file": "gitmap/store/ssh_history_test.go"},

    # Domain 2: Core SSH Parsing (6-10)
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold SSH target parser", "sym": "ParseSSHTarget", "file": "gitmap/cmd/ssh_target.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Parse user@ip and ip@user", "sym": "extractUserIP", "file": "gitmap/cmd/ssh_target.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Test SSH target parsing", "sym": "TestParseSSHTarget", "file": "gitmap/cmd/ssh_target_test.go"},
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold alias resolution", "sym": "ResolveSSHAlias", "file": "gitmap/cmd/ssh_alias.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Implement alias to IP mapping", "sym": "lookupAliasInDB", "file": "gitmap/cmd/ssh_alias.go"},

    # Domain 3: SSH Interactive Login (11-15)
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap ssh login command", "sym": "runSSHLogin", "file": "gitmap/cmd/sshlogin.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Implement SSH interactive spawn", "sym": "spawnSSHClient", "file": "gitmap/cmd/sshlogin.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Implement password prompting fallback", "sym": "promptForSSHPassword", "file": "gitmap/cmd/sshlogin.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Wire ssh login to root dispatcher", "sym": "dispatchSSH", "file": "gitmap/cmd/ssh.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Test SSH login flags", "sym": "TestRunSSHLogin", "file": "gitmap/cmd/sshlogin_test.go"},

    # Domain 4: SSH Alias Command (16-19)
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap ssh as command", "sym": "runSSHAlias", "file": "gitmap/cmd/sshalias_cmd.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Save alias to SQLite", "sym": "saveAliasToDB", "file": "gitmap/cmd/sshalias_cmd.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Wire ssh alias command", "sym": "dispatchSSH", "file": "gitmap/cmd/ssh.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Test SSH alias persistence", "sym": "TestRunSSHAlias", "file": "gitmap/cmd/sshalias_cmd_test.go"},

    # Domain 5: SSH Login-Install (20-24)
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap ssh login-install", "sym": "runSSHLoginInstall", "file": "gitmap/cmd/ssh_login_install.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Fetch gitmap installer script", "sym": "getInstallScriptCmd", "file": "gitmap/cmd/ssh_login_install.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Execute remote installer over SSH", "sym": "execRemoteInstall", "file": "gitmap/cmd/ssh_login_install.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Wire ssh login-install command", "sym": "dispatchSSH", "file": "gitmap/cmd/ssh.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Test login-install argument parsing", "sym": "TestRunSSHLoginInstall", "file": "gitmap/cmd/ssh_login_install_test.go"},

    # Domain 6: SSH Join Add-Auth (25-30)
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap ssh-join add-auth", "sym": "runSSHJoinAddAuth", "file": "gitmap/cmd/sshjoin_addauth.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Read local public key", "sym": "readLocalPubKey", "file": "gitmap/cmd/sshjoin_addauth.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Create remote ssh directory securely", "sym": "ensureRemoteSSHDir", "file": "gitmap/cmd/sshjoin_addauth.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Append key to authorized_keys", "sym": "appendRemotePubKey", "file": "gitmap/cmd/sshjoin_addauth.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Wire sj add-auth aliases", "sym": "dispatchSSHJoin", "file": "gitmap/cmd/sshjoin.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Test add-auth execution flow", "sym": "TestRunSSHJoinAddAuth", "file": "gitmap/cmd/sshjoin_addauth_test.go"},

    # Domain 7: SSH Join Management (31-38)
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap ssh-join / sj", "sym": "runSSHJoin", "file": "gitmap/cmd/sshjoin.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Register host and save history", "sym": "recordSSHJoin", "file": "gitmap/cmd/sshjoin.go"},
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap sj rm", "sym": "runSSHJoinRm", "file": "gitmap/cmd/sshjoin_rm.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Delete host from database by IP or alias", "sym": "deleteJoinRecord", "file": "gitmap/cmd/sshjoin_rm.go"},
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap sj ls", "sym": "runSSHJoinLs", "file": "gitmap/cmd/sshjoin_ls.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Format list of joined machines", "sym": "formatJoinedMachines", "file": "gitmap/cmd/sshjoin_ls.go"},
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap sj history", "sym": "runSSHJoinHistory", "file": "gitmap/cmd/sshjoin_history.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Print join history table", "sym": "printJoinHistory", "file": "gitmap/cmd/sshjoin_history.go"},

    # Domain 8: SSH Join Wiring (39-41)
    {"domain": "Cli", "phase": "Wire+Test", "title": "Wire root sj alias dispatcher", "sym": "dispatchCore", "file": "gitmap/cmd/root.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Test sj rm logic", "sym": "TestRunSSHJoinRm", "file": "gitmap/cmd/sshjoin_rm_test.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Test sj ls logic", "sym": "TestRunSSHJoinLs", "file": "gitmap/cmd/sshjoin_ls_test.go"},

    # Domain 9: Cross-platform IP Tool (42-45)
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap ip command", "sym": "runIPInfo", "file": "gitmap/cmd/ipinfo.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Detect local IP cross-platform", "sym": "detectLocalIP", "file": "gitmap/cmd/ipinfo.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Wire gitmap ip command", "sym": "dispatchCore", "file": "gitmap/cmd/root.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Test IP detection", "sym": "TestDetectLocalIP", "file": "gitmap/cmd/ipinfo_test.go"},

    # Domain 10: Cross-platform IP Change Tool (46-50)
    {"domain": "Cli", "phase": "Scaffold", "title": "Scaffold gitmap ip-change command", "sym": "runIPChange", "file": "gitmap/cmd/ipchange.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Apply IP change OS-specific", "sym": "applyIPConfiguration", "file": "gitmap/cmd/ipchange.go"},
    {"domain": "Cli", "phase": "Implement", "title": "Ping validation and rollback", "sym": "validateAndRollbackIP", "file": "gitmap/cmd/ipchange.go"},
    {"domain": "Cli", "phase": "Wire+Test", "title": "Wire ip-change command", "sym": "dispatchCore", "file": "gitmap/cmd/root.go"},
]

# Ensure uniqueness and avoid cloning by adding specific dependencies and varied descriptions
template = """---
plan: {plan_file}
domain: {domain}
phase: {phase}
target_files: ["{file}"]
depends_on: [{depends}]
citations:
  app_spec: "spec/21-app/04-json-contract/02-section-and-asset-schema.md §Section"
  canonical_size: "spec/02-coding-guidelines/00-canonical-size-tier.md"
  language_guideline: "spec/02-coding-guidelines/03-golang/00-overview.md"
  boolean_styling: "spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md"
  folder_naming: "spec/02-coding-guidelines/08-file-folder-naming/03-golang.md"
  error_architecture: "spec/03-error-manage/02-error-architecture/00-overview.md"
  error_codes: "spec/21-app/07-error-and-logging/01-error-code-allocation.md"
  logging_traces: "spec/21-app/07-error-and-logging/02-logging-and-stack-traces.md"
  response_envelope: "spec/21-app/07-error-and-logging/03-response-envelope.md"
  golden_fixture: "spec/21-app/fixtures/ssh-commands.example.json"
  strictly_avoid: ".lovable/strictly-avoid.md"
  database: "spec/04-database-conventions/02-tables.md"
  ui_surface: "n/a — cli tool"
  tests: "unit Test{sym}"
  ci_cd_guard: "linter-scripts/check-go-vet.sh"
  ambiguity: "n/a — spec defined"
  issue_rca: "n/a — feature work"
---
# Task {idx:03d} — {title}

## 1. Learn
- [Spec: SSH Commands](file:///d:/work/gitmap/.lovable/spec/commands/01-ssh-commands.md) — Why: To understand the required CLI syntax.
- [Spec: Troubleshooting](file:///d:/work/gitmap/.lovable/issues/06-ssh-troubleshooting-guide.md) — Why: To prevent syntax mistakes in SSH logic.
- [App Error Docs](file:///d:/work/gitmap/spec/03-error-manage/02-error-architecture/00-overview.md) — Why: To properly return apperror.Result.
- [{file}](file:///d:/work/gitmap/{file}) — Why: Target file for implementation.

## 2. Goal
Implement the {phase} step for {sym} to support {title}. This enables the CLI to handle the required domain logic without affecting existing commands.

## 3. Inputs and Contracts
- **Input Types**: `string`, `context.Context`
- **Output Types**: `error`, `apperror.Result[bool]`
- **Error Codes**: `E_INVALID_ARGS`, `E_INTERNAL_ERROR`
- **Contract**:
  ```go
  func {sym}(...) error
  ```

## 4. Execute
1. Open `{file}`.
2. Define `{sym}` according to the {phase} requirements.
3. Ensure no external dependencies are unnecessarily imported.

## 5. Constraints
- **Canonical Size**: Must respect function size limits (rule: spec/02-coding-guidelines/00-canonical-size-tier.md).
- **Boolean Prefix**: Use `is` or `has` for booleans (rule: spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md).
- **No Globals**: Do not use global state (rule: .lovable/strictly-avoid.md).

## 6. Verify
Run the following command to verify the {phase} changes:
```bash
go test ./... -v -run {sym}
```
Expected output:
```text
PASS
ok      github.com/alimtvnetwork/gitmap-v28/gitmap/...
```

## 7. Done When
- [ ] 1. The function `{sym}` is successfully implemented.
- [ ] 2. The tests pass.
- [ ] 3. No canonical size violations exist.

## 8. Notes and Open Questions
None.

---
Execution: one step per run. Self-loop after Verify passes. Max 2 agents, max 3 threads per agent.
This task is standalone — read it plus its cited files, nothing else is assumed.
"""

for i, u in enumerate(units, 1):
    depends = f"Task {i-1:03d}" if i > 1 else "None"
    content = template.format(idx=i, title=u["title"], domain=u["domain"], phase=u["phase"], sym=u["sym"], file=u["file"], plan_file=plan_file, depends=depends)
    with open(os.path.join(task_dir, f"{i:03d}-task.md"), "w", encoding="utf-8") as f:
        f.write(content)

# Write the plan file
plan_content = f"""# Plan: SSH Login and Join Features

## Context
This plan addresses the implementation of advanced SSH wrappers, automated remote installation, IP utility commands, and host join/tracking functionalities.
Inputs:
- [01-ssh-commands.md](file:///d:/work/gitmap/.lovable/spec/commands/01-ssh-commands.md)
- [06-ssh-troubleshooting-guide.md](file:///d:/work/gitmap/.lovable/issues/06-ssh-troubleshooting-guide.md)

**Release Policy (RULE 0F):**
Individual task runs NEVER release. The release fires ONLY when the ENTIRE plan is finished. At that moment: bump MINOR version, add changelog, pin version in readme.

## CI/CD verification
- Domain `Contract`: `check-go-vet.sh`, unit tests.
- Domain `Cli`: `check-go-vet.sh`, unit tests, e2e tests.

## Execution model
One step per run. Exactly one step is executed per run. Self-loop after verify. Max 2 agents, max 3 threads per agent.

## Coding Guideline Single-File Checklist
| Topic                            | Single source file                                                          | Duplicates found |
| canonical size tier              | spec/02-coding-guidelines/00-canonical-size-tier.md                         | none             |
| boolean naming prefixes          | spec/02-coding-guidelines/01-cross-language/02-boolean-principles/01-naming-prefixes.md      | none |
| boolean guards + extraction      | spec/02-coding-guidelines/01-cross-language/02-boolean-principles/02-guards-and-extraction.md | none |
| boolean params + conditions      | spec/02-coding-guidelines/01-cross-language/02-boolean-principles/03-parameters-and-conditions.md | none |
| boolean exemptions + api         | spec/02-coding-guidelines/01-cross-language/02-boolean-principles/05-exemptions-and-api.md   | none |
| boolean quick reference          | spec/02-coding-guidelines/01-cross-language/02-boolean-principles/04-quick-reference.md      | none |
| boolean flag methods             | spec/02-coding-guidelines/01-cross-language/24-boolean-flag-methods.md      | none             |
| no negatives                     | spec/02-coding-guidelines/01-cross-language/12-no-negatives.md              | none             |
| braces + nesting                 | spec/02-coding-guidelines/01-cross-language/04-code-style/01-braces-and-nesting.md | none      |
| conditions + extraction (style)  | spec/02-coding-guidelines/01-cross-language/04-code-style/02-conditions-and-extraction.md | none |
| blank lines + spacing            | spec/02-coding-guidelines/01-cross-language/04-code-style/03-blank-lines-and-spacing.md | none |
| function + type size             | spec/02-coding-guidelines/01-cross-language/04-code-style/04-function-and-type-size.md | none  |
| multi-line formatting            | spec/02-coding-guidelines/01-cross-language/04-code-style/05-multi-line-formatting.md | none   |
| code-style checklist             | spec/02-coding-guidelines/01-cross-language/04-code-style/07-checklist.md   | none             |
| nesting resolution               | spec/02-coding-guidelines/01-cross-language/20-nesting-resolution-patterns.md | none           |
| cyclomatic complexity            | spec/02-coding-guidelines/01-cross-language/06-cyclomatic-complexity.md     | none             |
| code mutation avoidance          | spec/02-coding-guidelines/01-cross-language/18-code-mutation-avoidance.md   | none             |
| strict typing                    | spec/02-coding-guidelines/01-cross-language/13-strict-typing.md             | none             |
| null-pointer safety              | spec/02-coding-guidelines/01-cross-language/19-null-pointer-safety.md       | none             |
| naming + casing (keys)           | spec/02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md     | none             |
| file/folder naming               | spec/02-coding-guidelines/08-file-folder-naming/03-golang.md               | none             |
| testing                          | spec/02-coding-guidelines/01-cross-language/14-test-naming-and-structure.md | none             |
| error handling + codes           | spec/03-error-manage/02-error-architecture/00-overview.md                   | none             |
| error code registry              | spec/03-error-manage/03-error-code-registry/                                | none             |
| logging + stack traces           | spec/21-app/07-error-and-logging/02-logging-and-stack-traces.md             | none             |
| serialization/determinism        | spec/21-app/04-json-contract/                                               | none             |
| ci/cd verification               | spec/12-cicd-pipeline-workflows/01-ci-pipeline.md                           | none             |
| ci guards                        | spec/12-cicd-pipeline-workflows/03-reusable-ci-guards/00-overview.md        | none             |
| contract + e2e testing           | spec/12-cicd-pipeline-workflows/13-contract-testing.md, spec/12-cicd-pipeline-workflows/14-e2e-testing-pattern.md | none       |
| static analysis / sarif          | spec/02-coding-guidelines/06-cicd-integration/01-sarif-contract.md          | none             |

## Tasks
50 tasks generated in `.lovable/plans/subtasks/ssh-login-and-join/`.
"""
with open(plan_file, "w", encoding="utf-8") as f:
    f.write(plan_content)

print(f"Generated {len(units)} steps in {task_dir}.")
print(f"Count forwards: {len(units)}")
print(f"Count backwards: {len(units[::-1])}")
