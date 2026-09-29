# Plan 194: Import-All-JSON, What-Configs (`wc`), and Fleet Deploy Hardening

## User Request (Verbatim)

```text
You can have this

gitmap commitin --config jsonenvelope/fixtures/commit-pull-config.json && gitmap macro import jsonenvelope/fixtures/macro.json -y && gitmap sj import jsonenvelope/fixtures/ssh-nodes.json -y && gitmap cmdui import-settings jsonenvelope/fixtures/ui-settings.json


buit also you can have like

gitmap import-all-json *
gitmap import-all-json *.json # only json
gitmap import-all-json specific.json
gitmap what-configs(wc) specific.json
gitmap what-configs(wc) * # yaml, json, db etc
gitmap what-configs(wc) *.json

Fix test and release new version please
```

## Spec Reference
- `02-spec/21-app/184-import-all-json-what-configs-and-fleet-deploy-hardening.md`

## Actionable Deliverables & Task Breakdown

- **Task-01: Machine `main` Connectivity & SSH Join Disambiguation**
  - Diagnose connection failure to `main` (`192.168.1.20`) with password ending in `@` (`rtyrty123@`).
  - Update `sshjoin_cmd.go`, `ssh_parser.go`, and `sshjoin_add_pass_cmd.go` to support symmetric argument order (`[alias] [user@host]` and `[user@host] [alias]`).
  - Guard trailing `@` in passwords to prevent misidentification as host address targets.
  - Verify live connectivity and command execution to `main`.

- **Task-02: Windows Fleet Deploy Command-Line Length Hardening**
  - Fix Windows `cmd.exe` 8,191-character buffer limit (`The command line is too long.`) during fleet deployment.
  - Stream payloads >= 1,500 characters directly to `.gitmap/fleet_config_deploy.json` via `streamDirectToRemote`.
  - Execute decoupled `gitmap ssh nodes import-json .gitmap/fleet_config_deploy.json` with deferred cleanup.
  - Fix Windows bidirectional verification syntax in `sshjoin_enroll.go` to use PowerShell conditional separator.

- **Task-03: `gitmap deploy import` Command Discovery & Rich Help**
  - Add `deploy import` and `deploy import-json` to `cli/cmdssh/ssh_deploy_router.go`.
  - Add to known keywords in `cli/cmdssh/ssh_deploy_cmd.go`.
  - Update `cli/cmdssh/ssh_deploy_rich_help.go` documentation.

- **Task-04: `gitmap import-all-json` Command Implementation**
  - Create `cli/cmd/import_all_json_cmd.go` supporting `*`, `*.json`, `specific.json`, and file arrays.
  - Inspect envelope and structural signatures, routing to respective subsystem handlers (`ssh-nodes`, `macro`, `commit-pull-config`, `ui-settings`, `templates`).
  - Support `--dry-run` and `-y` flags.
  - Register aliases `import-all-json`, `importalljson`, `import-json-all`, `importall` in `cli/cmd/roottooling.go`.
  - Add unit tests in `cli/cmd/import_all_json_cmd_test.go`.

- **Task-05: `gitmap what-configs` (`wc`) Inspection Command**
  - Create `cli/cmd/what_configs_cmd.go` supporting `*`, `*.json`, `specific.json`, and polyglot configs (`.json`, `.yaml`, `.yml`, `.toml`, `.db`, `.sqlite`, `.env*`, `.ini`).
  - Render ANSI table with file path, detected format category, recommended CLI command, and system state impact.
  - Register aliases `what-configs`, `wc`, `what-config`, `whatconfigs` in `cli/cmd/roottooling.go`.
  - Add unit tests in `cli/cmd/what_configs_cmd_test.go`.

- **Task-06: Verification, Testing & Release Ceremony**
  - Run all unit tests in `cli/cmd/...` and `cli/cmdssh/...`.
  - Ensure zero nested ifs across all new and modified Go files.
  - Execute release orchestrator to bump minor version, generate changelog, create release branch/tag, and push.
