# CI/CD Issue Diagnostics, Auto-Repair & Rerun ETA

## Context

When a pipeline fails, inspecting logs is only half the battle. Developers must diagnose which linter, style rule, or compiler gate caused the failure and quickly remediate it before triggering a rerun.

## Architectural Specification

### 1. Historical Baseline Rerun ETA

When displaying error logs (or in `--json` payload), GitMap computes the historical average duration of previous successful runs of the targeted workflow:
```text
● Estimated pipeline rerun duration (ETA): ~Xs
  (Derived from historical successful pipeline runs baseline)
```
In JSON payloads, this is emitted as `rerunEtaSeconds: <int>`.

### 2. Internal CI/CD Diagnostic Probes

GitMap embeds a local execution harness for internal repository checks:

| Probe Name | Target Script / Tool | Description | Auto-Fix Capability |
|---|---|---|---|
| `gofmt formatting` | `gofmt -l .` in `gitmap/` | Checks for unformatted Go code | `gofmt -w .` |
| `Nested If Linter` | `linter-scripts/check-nested-ifs.py` | Enforces zero nested `if` statements | Manual / Flattening |
| `Boolean & Enum Linter` | `linter-scripts/check-enum-and-boolean.py` | Enforces positive prefixes and enum suffixes | Manual / Renaming |
| `Relative Paths Linter` | `linter-scripts/check-relative-paths.py` | Prevents absolute drive letter paths | Manual / Path fix |
| `Error Management Linter` | `linter-scripts/check-error-management.py` | Enforces AppError wrappers & no swallowed errors | Manual / Wrapping |
| `Legacy Refs Check` | `.github/scripts/check-legacy-refs.py .` | Prevents outdated repository references | Manual / Replace |
| `Spell Check (US locale)` | `.github/scripts/misspell-changed.py` | Validates American English spelling | Manual / Correction |
| `Go Compile Gate` | `go test -run=^$ ./... -count=1` | Verifies full compiler type-check | Code compilation |

### 3. Execution Control

1. **Automated Mode (`--fix` / `-f`)**: Executes all probes and applies automated repairs (`gofmt -w .`) without interactive prompts.
2. **Read-Only Check Mode (`--check` / `-c`)**: Executes all probes without applying file modifications.
3. **Interactive Mode**: Prompts the user:
   `Would you like to run internal CI/CD diagnostic & auto-repair scripts? [y/N]: `
   When confirmed or when `-y` / `--yes` is present, executes the suite and prints formatted results.
