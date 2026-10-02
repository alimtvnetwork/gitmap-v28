# Subtask 05: Repository Canonical Deduplication & Structured Sub-Node Remediation Solutions

> **Parent Plan:** [.ai-memory/plans/pending/65-gitmap-update-all-zip-and-fixes.md](../../pending/65-gitmap-update-all-zip-and-fixes.md)  
> **Spec Reference:** [02-spec/21-app/65-gitmap-update-all-zip-and-fixes/01-architecture-spec.md](../../../../02-spec/21-app/65-gitmap-update-all-zip-and-fixes/01-architecture-spec.md)  
> **Status:** `PENDING`  
> **Target Files:**
> - `cli/cmdpull/pull_efficient.go`
> - `cli/cmdpull/pull_efficient_render.go`
> - `cli/cmdpull/pull_remediation_hint.go`

---

## 1. Technical Objective

1. **Path Canonical Deduplication**:
   - Prevent the same physical repository from appearing multiple times in `gitmap pull-all`, `gitmap pas`, `gitmap pae`, and `gitmap status`.
   - Normalize paths on Windows using `filepath.Clean(strings.ToLower(r.AbsolutePath))` and filter out duplicates in `resolveAllTrackedRecords()`.
2. **Ambiguous Repository Disambiguation**:
   - When different repositories in different directories share identical names (e.g., `D:\work\frontend` and `D:\external\frontend`), display the path alongside the repository name so the user can distinguish them.
3. **Structured Sub-Node Remediation Solutions (Option 1 vs Option 2)**:
   - For pull and status failures, upgrade the single-line remediation hint to structured sub-nodes offering up to two distinct operational solutions:
     - **Diverged Branch**: Option 1 (Preserve Local / Rebase) vs Option 2 (Discard Local / Hard Reset).
     - **Dirty Tree**: Option 1 (Commit WIP via `cpar`) vs Option 2 (Stash Changes via `stash`).
     - **Auth Failure**: Option 1 (Fix Credential Store via `fc`) vs Option 2 (Deploy SSH Keys).

---

## 2. Implementation Scope & File Edits

### Target 1: `cli/cmdpull/pull_efficient.go`
- **Deduplicate Tracked Repository Records**:
  In `resolveAllTrackedRecords() []model.ScanRecord`:
  ```go
  func resolveAllTrackedRecords() []model.ScanRecord {
      db, err := store.OpenDefault()
      if err != nil {
          return nil
      }
      defer db.Close()

      records, err := db.ListRepos()
      if err != nil {
          return nil
      }

      seen := make(map[string]bool, len(records))
      unique := make([]model.ScanRecord, 0, len(records))
      for _, r := range records {
          if r.AbsolutePath == "" {
              continue
          }
          canonical := filepath.Clean(strings.ToLower(r.AbsolutePath))
          if seen[canonical] {
              continue
          }
          seen[canonical] = true
          unique = append(unique, r)
      }
      return unique
  }
  ```

### Target 2: `cli/cmdpull/pull_remediation_hint.go`
- **Introduce Structured Remediation Models**:
  ```go
  type RemediationOption struct {
      OptionNumber int    `json:"option_number"`
      Title        string `json:"title"`
      Command      string `json:"command"`
  }

  type StructuredRemediation struct {
      Reason  string              `json:"reason"`
      Options []RemediationOption `json:"options"`
  }
  ```
- **Implement `ResolveStructuredRemediation(s *PullRepoState) StructuredRemediation`**:
  - **Diverged Branch**:
    - Option 1 (Preserve Local / Rebase): `gitmap pull --rebase <repo>` (or `git -C "<path>" pull --rebase`)
    - Option 2 (Discard Local / Hard Reset): `git -C "<path>" reset --hard origin/<branch>`
  - **Dirty Working Tree**:
    - Option 1 (Commit WIP): `gitmap cpar "wip: save changes"`
    - Option 2 (Stash Changes): `gitmap stash` (or `git -C "<path>" stash`)
  - **Authentication Rejection**:
    - Option 1 (Fix Windows Credentials): `gitmap fix-credential`
    - Option 2 (Deploy SSH Keys): `gitmap ssh deploy-keys`

### Target 3: `cli/cmdpull/pull_efficient_render.go`
- **Render Distinct Paths for Colliding Names**:
  - In `renderUpdatedGroup`, `renderDirtyGroup`, `renderFailedGroup`:
  - Calculate frequency of each `RepoName` across all records.
  - If `count > 1`, render format as `fmt.Sprintf("%s (%s)", s.RepoName, s.RepoPath)`.
- **Render Structured Sub-Nodes in Failed & Dirty Groups**:
  ```go
  for _, opt := range structured.Options {
      fmt.Fprintf(w, "        %s↳ Sub-node Option %d (%s):%s\n",
          constants.ColorCyan, opt.OptionNumber, opt.Title, constants.ColorReset)
      fmt.Fprintf(w, "          %s%s%s\n",
          constants.ColorDim, opt.Command, constants.ColorReset)
  }
  ```

---

## 3. Verification Protocol

- Run pull package unit tests:
  ```bash
  go test -v ./cli/cmdpull -run TestPull
  ```
- Verify compilation across CLI packages:
  ```bash
  go build ./cli/...
  ```
- CLI Manual Check:
  - Simulate a diverged or dirty repo in workspace.
  - Execute `gitmap pull-all` or `gitmap pae`.
  - Verify deduplication of Windows case variants and appearance of Option 1 vs Option 2 sub-nodes.
