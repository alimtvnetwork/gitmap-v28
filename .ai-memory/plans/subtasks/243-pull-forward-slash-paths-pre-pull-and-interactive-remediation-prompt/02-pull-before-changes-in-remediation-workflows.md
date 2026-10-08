# Subtask 02: Pull Before Changes in Remediation Workflows

- **Target Files:**
  - `cli/gitutil/remediation_stash.go`
  - `cli/gitutil/remediation_commit.go`
  - `cli/gitutil/remediation_discard.go`
  - `cli/cmdpull/pull_remediation_hint.go`

- **Checklist:**
  - [x] Audit remediation recipes in `cli/gitutil/remediation_*.go` to ensure pulling is sequenced prior to finalizing state (`stash -u && pull && stash pop`, `add -A && commit && pull --rebase`, `reset --hard && clean -fd && pull`).
  - [x] Update `resolveDirtyTreeDualHints` to ensure `commitCmd` and `stashCmd` include remote pulls before finishing.
  - [x] Verify non-destructive invariants.
