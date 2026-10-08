# Subtask 03: Interactive Remediation Prompt (All, Single, Skip)

- **Target Files:**
  - `cli/cmdpull/pull_execute.go`
  - `cli/cmd/remediation_box.go`
  - `cli/cmd/remediation_box_prompt_test.go`

- **Checklist:**
  - [x] In `cli/cmdpull/pull_execute.go:handlePullRemediation`:
        Removed `isConcisePullOutput(opts.all, opts.showStatus)` check so `gitmap pa` / `gitmap pull-all` does not silently suppress prompts when dirty/failed repos exist.
  - [x] In `cli/cmd/remediation_box.go:promptForRemediation`:
        Implemented interactive prompt menu:
        ```text
        Remediate dirty / failed repository(ies) now?
          [a/1] Fix all (stash & re-apply with pre-pull)
          [s/2] Fix single / step through repositories
          [k/q] Skip / Exit
        Choice [a/s/k]:
        ```
  - [x] Supported `[a/1]` for all, `[s/2]` for single/stepping, `[k/q]` for skip.
  - [x] Ensured non-interactive / JSON mode skips cleanly without hanging.
