# Subtask 01: Universal Forward-Slash Path Normalization

- **Target Files:**
  - `cli/cmdpull/pull_remediation_hint.go`
  - `cli/cmdpull/pull_efficient_render.go`
  - `cli/cmdpull/pull_batch_remediation.go`
  - `cli/cmd/remediation_box.go`
  - `cli/cmd/reconcile_prompt.go`

- **Checklist:**
  - [x] Update `formatRepoGitCmd(repoDir, gitArgs)` to convert `repoDir` to forward slashes:
        `cleanDir := filepath.ToSlash(filepath.Clean(repoDir))`
        `fmt.Sprintf("git -C \"%s\" %s", cleanDir, gitArgs)`
  - [x] Update itemized dirty file rendering in `renderItemizedDirtyFiles` to call `filepath.ToSlash(f)`.
  - [x] Normalize `item.RepoPath` and dirty file entry formatting in `remediation_box.go` and `reconcile_prompt.go`.
  - [x] Verify unit tests in `cli/cmdpull/` reflect forward slash paths without `\\` escapes.
