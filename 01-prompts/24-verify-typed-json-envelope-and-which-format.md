# 24-verify-typed-json-envelope-and-which-format.md

Use this prompt in any CI pipeline, local test runner, or autonomous AI agent to rigorously verify that the Typed JSON Envelope architecture (`attributes` + `data`), format inspection CLI command (`gitmap which-format`), subsystem dual-format ingestion, and `repo-secrets` manifest normalization are operational and regression-free.

## Pre-flight Quality Gate Checklist

1. **Verify Typed JSON Envelope Schema**:
   - Confirm standard envelope format contains two top-level keys:
     - `attributes`: `type` (e.g. `ssh-nodes`, `macro`, `commit-pull-config`, `ui-settings`, `templates`, `repo-manifest`), `source`, `how`, `version` (default `"1.0"`), and `timestamp` (UTC ISO8601).
     - `data`: core subsystem payload (nodes list, macro step tree, config settings, templates, etc.).
   - Verify unit tests: `go test -v ./jsonenvelope/...` passes cleanly.

2. **Verify Format Inspection CLI (`gitmap which-format`)**:
   - Run `gitmap which-format <file1.json> [file2.json...]` -> Displays formatted box summary with matched/unmatched status.
   - Run `gitmap which-format` (zero args) in any directory containing `.json` files -> Automatically discovers and inspects all JSON files in the working directory.
   - Run `gitmap which-format --dir <directory>` -> Auto-scans target folder and displays:
     - Total scanned, matched count, and unmatched count.
     - For matched files: type descriptor, human-readable name, source, exact suggested import command, and system impact description.
     - For unmatched files: clear explanatory notice ("Could not match known GitMap JSON format / unsupported schema - will not import").
     - Single-line copy-pasteable batch command combining all matched import actions with advice to append `-y`.
   - Verify aliases: `gitmap which format`, `gitmap format which`, `gitmap format inspect`, `gitmap which-json`.

3. **Verify Subsystem Dual-Format Ingestion**:
   - `gitmap sj import` / `gitmap ssh import`: Seamlessly ingests both typed envelope (`attributes` + `data`) and legacy flat node JSONs without errors.
   - `gitmap macro import`: Unwraps typed envelope and imports steps into SQLite macro storage.
   - `gitmap commitin --config`: Parses envelope-wrapped or flat configuration objects.
   - `gitmap templates import`: Unwraps typed envelopes and accepts both scalar strings and stringified JSON arrays in variables.
   - `gitmap clone`: Unwraps typed envelope scan manifests and executes repository plans.
   - `gitmap cmdui`: Saves settings wrapped in typed envelope and loads from both legacy and enveloped files.

4. **Verify `repo-secrets` Normalization & Zero Secrets Leakage**:
   - Confirm all manifest JSON files in `D:\work\repo-secrets` (`01-gitmap/`, `04-w1-machine/`, `05-w2-machine/`, `06-w3-machine/`, `07-final-network-machine/`) are migrated to the typed envelope architecture.
   - Strict Security Rule: Never commit passwords, private keys, or sensitive tokens outside `repo-secrets`. Test fixtures in `cli/jsonenvelope/fixtures/` must contain only mock/dummy data.

5. **Verify Clean Cache & Workstation Hygiene**:
   - Run `python 03-ai-scripts/42-clean-test-and-build-caches.py` -> Cleans temp files, Go build/test cache, and temporary test databases.

6. **CI/CD Quality Gate**:
   - Verify all workflows on `main` are 100% green via `gh run list --limit 6` or `gitmap pipeline-ai status`.
