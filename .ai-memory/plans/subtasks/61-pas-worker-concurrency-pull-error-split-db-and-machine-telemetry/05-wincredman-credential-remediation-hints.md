# Subtask 05: Wincredman Credential Store Remediation Hints

> **Parent Plan:** [61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry](../../pending/61-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Tracking Spec:** [198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md](../../../../02-spec/21-app/198-pas-worker-concurrency-pull-error-split-db-and-machine-telemetry.md)  
> **Primary File Targets:** `cli/cmdpull/pull_remediation_hint.go`  

---

## 1. Objective

Provide automated, actionable recovery guidance when Git operations encounter Windows Credential Manager failures:
1. Detect error signature: `fatal: Unable to persist credentials with the 'wincredman' credential store`.
2. Format actionable recovery instructions:
   - Run `cmdkey /list` to inspect saved credentials.
   - Configure Git credential store fallback (`git config --global credential.helper manager` or `wincred`).
   - Run `gitmap pull-error` to view full diagnostic logs.

---

## 2. Implementation Scope

- **`cli/cmdpull/pull_remediation_hint.go`:**
  - Enhance `detectRemediationHint(errStr string)` to recognize `wincredman` and credential persistence failure signatures.
  - Return formatted remediation block with copy-pasteable recovery commands.

---

## 3. Verification

- Verify hint generator recognizes wincredman errors.
- Verify compilation with `go build ./...` in `cli/`.
