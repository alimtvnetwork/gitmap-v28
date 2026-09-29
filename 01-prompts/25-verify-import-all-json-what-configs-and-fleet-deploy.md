# 25-verify-import-all-json-what-configs-and-fleet-deploy.md

Use this prompt in any CI pipeline, local test runner, or autonomous AI agent to rigorously verify that `gitmap import-all-json`, `gitmap what-configs` (`wc`), and fleet deployment hardening are operational and regression-free.

## Pre-flight Quality Gate Checklist

1. **Verify `gitmap import-all-json` Command**:
   - Run `gitmap import-all-json *` or `gitmap import-all-json *.json` -> Scans directory, detects envelope and flat JSON types, and imports configurations into respective subsystems (`ssh-nodes`, `macro`, `commit-pull-config`, `ui-settings`, `templates`).
   - Run `gitmap import-all-json <specific.json>` -> Directly inspects and imports single configuration file.
   - Run with `--dry-run` -> Previews actions without applying changes.
   - Run with `-y` -> Bypasses interactive confirmation.
   - Verify aliases: `import-all-json`, `importalljson`, `import-json-all`, `importall`.
   - Verify unit test passes: `go test -v ./cmd -run TestImportAllJson`.

2. **Verify `gitmap what-configs` (`wc`) Command**:
   - Run `gitmap what-configs <specific.json>` -> Classifies JSON envelope or structural schema.
   - Run `gitmap what-configs *` -> Polyglot scan of JSON, YAML, TOML, SQLite databases, INI, and ENV files.
   - Run `gitmap what-configs *.json` -> Scans all JSON files in the working directory.
   - Verify ANSI table output includes: FILE, CATEGORY, RECOMMENDED COMMAND, and SYSTEM IMPACT.
   - Verify aliases: `what-configs`, `wc`, `what-config`, `whatconfigs`.
   - Verify unit test passes: `go test -v ./cmd -run TestWhatConfigs`.

3. **Verify Fleet Deploy Windows Command-Line Length Hardening**:
   - Verify fleet deployments (`gitmap deploy config ssh`) do not fail with `The command line is too long.` on Windows nodes.
   - Verify payloads >= 1,500 characters stream directly to `.gitmap/fleet_config_deploy.json` via `streamDirectToRemote`.
   - Verify remote execution triggers `gitmap ssh nodes import-json .gitmap/fleet_config_deploy.json` followed by deferred cleanup.

4. **Verify Symmetric SSH Join & Special Character Passwords**:
   - Verify `gitmap ssh join <alias> <user@host>` and `gitmap ssh join <user@host> <alias>` both resolve properly.
   - Verify passwords containing trailing `@` (e.g. `rtyrty123@`) are not misidentified as host addresses.
   - Verify bidirectional connectivity check executes with PowerShell syntax on Windows targets.

5. **Verify `gitmap deploy import` Discovery**:
   - Verify `gitmap deploy import [file]` routes cleanly to `RunSSHNodesImportJSON`.
   - Verify documentation present in `gitmap deploy --help`.

6. **Verify Workstation Hygiene & Zero Secrets Leakage**:
   - Run `python 03-ai-scripts/42-clean-test-and-build-caches.py`.
   - Ensure zero production credentials or secrets committed outside `D:\work\repo-secrets`.
