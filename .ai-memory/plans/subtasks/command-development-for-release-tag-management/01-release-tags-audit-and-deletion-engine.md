# Engineering Subtask Plan: Release Tag Audit & Deletion Engine (`01-release-tags-audit-and-deletion-engine`)

**Subtask ID:** `01-release-tags-audit-and-deletion-engine`  
**Subtask Code:** `Task-01-ReleaseTagsEngine`  
**Parent Plan:** `command-development-for-release-tag-management` ([command-development-for-release-tag-management.md](../../command-development-for-release-tag-management.md))  
**Run Number:** 84  
**Assigned Agent Role:** Implementation Worker 01 (Downstream Execution)  
**Spec Reference:**  
- [01-architecture-and-audit-engine-spec.md](../../../../02-spec/21-app/command-development-for-release-tag-management/01-architecture-and-audit-engine-spec.md)  
**Status:** `READY FOR IMPLEMENTATION`  

---

## 1. Targeted File Inventory & Disjoint Boundaries

To guarantee multi-agent isolation and eliminate concurrency conflicts:

### 1.1 Owned Subtask & Spec Files (Author 01)
- `02-spec/21-app/command-development-for-release-tag-management/01-architecture-and-audit-engine-spec.md`
- `.ai-memory/plans/subtasks/command-development-for-release-tag-management/01-release-tags-audit-and-deletion-engine.md`

### 1.2 Downstream Owned Implementation Files (Worker 01)
- `cli/cmdfixreleasetags/types.go` (Domain models, reason enums, step statuses, positive booleans)
- `cli/cmdfixreleasetags/safety_guard.go` (Safety invariants: active version protection, latest healthy release immunity, 15-min grace window)
- `cli/cmdfixreleasetags/audit_engine.go` (Tag discovery, GitHub release inspection, asset verification, CI/CD run correlation)
- `cli/cmdfixreleasetags/delete_executor.go` (3-tier deletion execution flow + sidecar metadata cleanup)
- `cli/cmdfixreleasetags/audit_engine_test.go` (Unit tests for audit criteria detection and safety rules)
- `cli/cmdfixreleasetags/delete_executor_test.go` (Unit tests for deletion sequence and error isolation)

### 1.3 Disjoint Boundaries with Subtask 02 (Worker 02)
Worker 02 owns CLI routing and user interaction files:
- `cli/cmdfix/fix_cmd.go` (Routing `gitmap fix release tags`, `fix-release-tags`, `frt`)
- `cli/cmdfixreleasetags/fix_release_tags_cmd.go` (Cobra command setup, CLI entrypoint, flag parsing `--dry-run`, `-n`, `-y`, `--confirm`, `--json`)
- `cli/cmdfixreleasetags/ui_render.go` (Diagnostic table preview via `termout.PrintTable`, interactive confirmation prompt)
- `cli/constants/constants_cli.go` (Command registration and help strings)
- `cli/cmdroothelp/roothelp_clusters.go` (Help cluster integration)
- `cli/helpdoc/fix-release-tags.md` (Markdown documentation)

Worker 01 exposes a clean public API consumed by Worker 02:
```go
package cmdfixreleasetags

func AuditReleaseTags(repoPath string, opts AuditOptions) (*AuditReport, error)
func EvaluateTagSafety(tag string, activeVersion string, latestHealthyTag string, runs []CIWorkflowRunInfo, now time.Time) (ProtectionStatus, bool)
func ExecuteTagDeletion(repoPath string, record ReleaseTagAuditRecord) (*ReleaseTagDeletionResult, error)
func ExecuteBatchDeletions(repoPath string, records []ReleaseTagAuditRecord) ([]ReleaseTagDeletionResult, error)
```

### 1.4 Strictly Prohibited Actions
- **TOTAL BAN ON GIT COMMANDS:** Subagents MUST NOT execute `git add`, `git commit`, `git status`, `git push`, `git diff`, or `git checkout`.
- **Strictly Relative Paths Only:** All file references in documentation, code comments, and test fixtures must use forward-slash paths relative to repository root (`02-spec/...`, `.ai-memory/...`, `cli/...`). Zero absolute paths or `file:///` URIs.
- **Search Exclusively via GitMap:** Use `gitmap aum search`, `gitmap find`, and `gitmap cat`. Strict ban on `grep`, `ripgrep`, `rg`, `Select-String`, `findstr`.
- **No Build / No Test in Spec Phase:** Spec authoring requires zero build or test execution.
- **Positive Booleans Only:** All fields and predicate functions must be positive booleans (`IsValid`, `HasAssets`, `IsProtected`, `IsDraft`, `IsCompleted`). Double negatives or inverted prefixes are strictly forbidden.

---

## 2. Objective & Detailed Scope

Worker 01 will implement the audit engine and safe deletion executor in package `cli/cmdfixreleasetags`:

### 2.1 Audit Engine Criteria Detection
Implement logic to detect the five specified audit failure modes:
1. `ORPHAN_TAG`: Git tag exists in local or remote tracking refs, but no corresponding GitHub release is published.
2. `DRAFT_RELEASE`: A GitHub release exists with `isDraft == true`.
3. `MISSING_ASSETS`: A published GitHub release exists with 0 uploaded binary assets or missing checksums.
4. `CORRUPTED_ASSETS`: An uploaded binary asset has 0 bytes or is truncated.
5. `CICD_FAILED`: GitHub Actions workflow run for the tag/commit failed (`conclusion != "success"` and `conclusion != "neutral"`).

### 2.2 Safety Invariant Enforcement
Enforce non-negotiable protections before flagging any tag for deletion:
1. **Active Binary Version:** Protect running version (`constants.Version`). If tag matches `constants.Version` (ignoring leading `v`), mark `PROTECTED_ACTIVE_VERSION` and prevent deletion.
2. **Latest Healthy Release:** The highest semantic version release that is healthy (`TagReasonHealthy`) is marked `PROTECTED_LATEST_HEALTHY` and immune from deletion.
3. **In-Progress Grace Window (15 Minutes):** Protect queued, waiting, or in-progress workflow runs, or runs created within 15 minutes of execution time (`PROTECTED_GRACE_WINDOW`).

### 2.3 3-Tier Deletion Execution Flow + Sidecar Metadata
Execute surgical removal in exact 4-step sequence:
- **Step 1:** Delete GitHub release via API / `gh release delete <tag> -y`.
- **Step 2:** Delete remote git tag via `git push origin :refs/tags/<tag>` / API.
- **Step 3:** Delete local git tag via `git tag -d <tag>`.
- **Step 4:** Clean local sidecar metadata in `.gitmap/release/<tag>.json` and `.gitmap/release/temp/`.

---

## 3. Operational Step-by-Step Implementation Instructions

### Step 1: Create `cli/cmdfixreleasetags/types.go`
- **Objective:** Establish domain models, reason enums, step statuses, and configuration options.
- **Key Definitions:**
  ```go
  package cmdfixreleasetags

  import "time"

  type AuditReason string

  const (
      ReasonHealthy         AuditReason = "HEALTHY"
      ReasonOrphanTag       AuditReason = "ORPHAN_TAG"
      ReasonDraftRelease    AuditReason = "DRAFT_RELEASE"
      ReasonMissingAssets   AuditReason = "MISSING_ASSETS"
      ReasonCorruptedAssets AuditReason = "CORRUPTED_ASSETS"
      ReasonCICDFailed      AuditReason = "CICD_FAILED"
  )

  type ProtectionStatus string

  const (
      StatusEligible        ProtectionStatus = "ELIGIBLE"
      StatusProtectedActive ProtectionStatus = "PROTECTED_ACTIVE_VERSION"
      StatusProtectedLatest ProtectionStatus = "PROTECTED_LATEST_HEALTHY"
      StatusProtectedGrace  ProtectionStatus = "PROTECTED_GRACE_WINDOW"
  )

  type StepStatus string

  const (
      StepPending StepStatus = "PENDING"
      StepSuccess StepStatus = "SUCCESS"
      StepSkipped StepStatus = "SKIPPED"
      StepFailed  StepStatus = "FAILED"
  )

  type ReleaseAssetInfo struct {
      Name    string `json:"name"`
      Size    int64  `json:"size"`
      State   string `json:"state"`
      IsValid bool   `json:"isValid"`
  }

  type CIWorkflowRunInfo struct {
      RunId        uint64        `json:"runId"`
      WorkflowName string        `json:"workflowName"`
      Status       string        `json:"status"`
      Conclusion   string        `json:"conclusion"`
      CreatedAt    time.Time     `json:"createdAt"`
      Duration     time.Duration `json:"duration"`
      IsSuccessful bool          `json:"isSuccessful"`
      IsInProgress bool          `json:"isInProgress"`
  }

  type ReleaseTagAuditRecord struct {
      Tag                   string              `json:"tag"`
      CommitSha             string              `json:"commitSha"`
      HasLocalTag           bool                `json:"hasLocalTag"`
      HasRemoteTag          bool                `json:"hasRemoteTag"`
      HasGitHubRelease      bool                `json:"hasGitHubRelease"`
      IsDraft               bool                `json:"isDraft"`
      AssetCount            int                 `json:"assetCount"`
      Assets                []ReleaseAssetInfo  `json:"assets,omitempty"`
      HasChecksumFile       bool                `json:"hasChecksumFile"`
      WorkflowRuns          []CIWorkflowRunInfo `json:"workflowRuns,omitempty"`
      AuditReason           AuditReason         `json:"auditReason"`
      ProtectionStatus      ProtectionStatus    `json:"protectionStatus"`
      IsEligibleForDeletion bool                `json:"isEligibleForDeletion"`
  }

  type StepExecutionResult struct {
      StepName string     `json:"stepName"`
      Status   StepStatus `json:"status"`
      Detail   string     `json:"detail,omitempty"`
  }

  type ReleaseTagDeletionResult struct {
      Tag          string                `json:"tag"`
      AuditReason  AuditReason           `json:"auditReason"`
      Steps        []StepExecutionResult `json:"steps"`
      IsCompleted  bool                  `json:"isCompleted"`
      ErrorMessage string                `json:"errorMessage,omitempty"`
  }

  type AuditSummary struct {
      TotalTagsChecked       int `json:"totalTagsChecked"`
      HealthyReleasesCount   int `json:"healthyReleasesCount"`
      OrphanTagsCount        int `json:"orphanTagsCount"`
      DraftReleasesCount     int `json:"draftReleasesCount"`
      MissingAssetsCount     int `json:"missingAssetsCount"`
      CorruptedAssetsCount   int `json:"corruptedAssetsCount"`
      CICDFailedCount        int `json:"cicdFailedCount"`
      ProtectedTagsCount     int `json:"protectedTagsCount"`
      EligibleDeletionsCount int `json:"eligibleDeletionsCount"`
  }

  type AuditReport struct {
      RepoPath         string                     `json:"repoPath"`
      ActiveVersion    string                     `json:"activeVersion"`
      LatestHealthyTag string                     `json:"latestHealthyTag"`
      Summary          AuditSummary               `json:"summary"`
      Records          []ReleaseTagAuditRecord    `json:"records"`
      DeletionResults  []ReleaseTagDeletionResult `json:"deletionResults,omitempty"`
  }

  type AuditOptions struct {
      GracePeriod time.Duration
      ActiveVersionOverride string
      MockExec bool
  }
  ```

---

### Step 2: Implement `cli/cmdfixreleasetags/safety_guard.go`
- **Objective:** Implement invariant evaluation logic to safeguard active versions, latest healthy releases, and in-flight CI/CD runs.
- **Key Logic:**
  1. `NormalizeVersionTag(v string) string`: Strips leading `v` for comparison against `constants.Version`.
  2. `IsActiveVersion(tag, activeVersion string) bool`: Compares normalized tag with normalized `constants.Version`.
  3. `EvaluateTagSafety(tag string, activeVersion string, latestHealthyTag string, runs []CIWorkflowRunInfo, now time.Time, graceWindow time.Duration) (ProtectionStatus, bool)`:
     - Check 1: If `IsActiveVersion(tag, activeVersion)` -> return `(StatusProtectedActive, false)`.
     - Check 2: If `tag == latestHealthyTag` -> return `(StatusProtectedLatest, false)`.
     - Check 3: Check workflow runs. If any run has `IsInProgress == true` or `now.Sub(run.CreatedAt) < graceWindow` -> return `(StatusProtectedGrace, false)`.
     - Check 4: Return `(StatusEligible, true)`.

---

### Step 3: Implement `cli/cmdfixreleasetags/audit_engine.go`
- **Objective:** Collect tags from local repository, remote origin, GitHub release API, and GitHub Actions runs; evaluate audit criteria.
- **Key Functions:**
  1. `CollectLocalTags(repoPath string) ([]string, error)`: Runs `git -C <repo> tag -l` and parses non-empty lines.
  2. `CollectRemoteTags(repoPath string) (map[string]bool, error)`: Runs `git -C <repo> ls-remote --tags origin` and maps tag names.
  3. `FetchGitHubReleases(repo string) ([]GitHubReleaseItem, error)`: Queries `gh release list --repo <repo> --limit 100 --json tagName,isDraft,createdAt,assets`. Supports pluggable command executor for unit testing.
  4. `FetchTagWorkflows(repo, tag, commitSha string) ([]CIWorkflowRunInfo, error)`: Queries `gh run list --repo <repo> --commit <commitSha> --limit 10 --json databaseId,name,status,conclusion,createdAt`.
  5. `AuditReleaseTags(repoPath string, opts AuditOptions) (*AuditReport, error)`:
     - Resolves `repo` owner/slug via `gitutil.ResolveRemoteOrigin(repoPath)`.
     - Collects all unique tag names across local, remote, and GitHub releases.
     - Identifies the `LatestHealthyTag` by parsing semantic versions in descending order and verifying healthy status.
     - For each tag, evaluates the 5 audit criteria in order:
       * If no GitHub release -> `ReasonOrphanTag`
       * Else if `release.IsDraft` -> `ReasonDraftRelease`
       * Else if `len(release.Assets) == 0` or missing checksum -> `ReasonMissingAssets`
       * Else if any asset has size `== 0` -> `ReasonCorruptedAssets`
       * Else if latest workflow run has `conclusion != "success"` -> `ReasonCICDFailed`
       * Else -> `ReasonHealthy`
     - Applies `EvaluateTagSafety(...)` to compute `ProtectionStatus` and `IsEligibleForDeletion`.
     - Aggregates report and summary statistics.

---

### Step 4: Implement `cli/cmdfixreleasetags/delete_executor.go`
- **Objective:** Implement the safe 4-step sequential deletion flow with error isolation and sidecar cleanup.
- **Key Functions:**
  1. `ExecuteTagDeletion(repoPath string, record ReleaseTagAuditRecord) (*ReleaseTagDeletionResult, error)`:
     - Initializes `ReleaseTagDeletionResult` with steps: GitHub Release, Remote Git Tag, Local Git Tag, Sidecar Metadata.
     - **Step 1 (GitHub Release):**
       * If `record.HasGitHubRelease`: executes `gh release delete <record.Tag> -y` via command executor.
       * If success: records `StepSuccess`.
       * If `!record.HasGitHubRelease`: records `StepSkipped` (e.g. for `ORPHAN_TAG`).
     - **Step 2 (Remote Git Tag):**
       * If `record.HasRemoteTag`: executes `git -C <repo> push origin :refs/tags/<record.Tag>`.
       * If success: records `StepSuccess`.
       * If `!record.HasRemoteTag`: records `StepSkipped`.
     - **Step 3 (Local Git Tag):**
       * If `record.HasLocalTag`: executes `git -C <repo> tag -d <record.Tag>`.
       * If success: records `StepSuccess`.
       * If `!record.HasLocalTag`: records `StepSkipped`.
     - **Step 4 (Sidecar Metadata):**
       * Checks `.gitmap/release/<tag>.json`. If file exists, removes it via `os.Remove`.
       * Checks `.gitmap/release/temp/` for matching staging files and purges them.
       * Records `StepSuccess`.
     - Sets `result.IsCompleted = true`.
  2. `ExecuteBatchDeletions(repoPath string, records []ReleaseTagAuditRecord) ([]ReleaseTagDeletionResult, error)`:
     - Iterates through all records where `record.IsEligibleForDeletion == true`.
     - Invokes `ExecuteTagDeletion` for each and returns aggregated results.

---

### Step 5: Implement `cli/cmdfixreleasetags/audit_engine_test.go`
- **Objective:** Unit test all 5 audit criteria and all 3 safety invariants using mocked command outputs.
- **Test Inventory:**
  1. `TestAuditCriteria_OrphanTag`: Verify tag present in git but missing from GitHub release is flagged `ReasonOrphanTag`.
  2. `TestAuditCriteria_DraftRelease`: Verify release with `isDraft: true` is flagged `ReasonDraftRelease`.
  3. `TestAuditCriteria_MissingAssets`: Verify release with 0 assets is flagged `ReasonMissingAssets`.
  4. `TestAuditCriteria_CorruptedAssets`: Verify release with 0-byte asset is flagged `ReasonCorruptedAssets`.
  5. `TestAuditCriteria_CICDFailed`: Verify workflow run with `conclusion: "failure"` is flagged `ReasonCICDFailed`.
  6. `TestSafetyGuard_ActiveVersion`: Verify tag matching `constants.Version` receives `StatusProtectedActive` and `IsEligibleForDeletion == false`.
  7. `TestSafetyGuard_LatestHealthy`: Verify latest valid release receives `StatusProtectedLatest` and `IsEligibleForDeletion == false`.
  8. `TestSafetyGuard_GraceWindow`: Verify workflow in progress or created 5 minutes ago receives `StatusProtectedGrace` and `IsEligibleForDeletion == false`.

---

### Step 6: Implement `cli/cmdfixreleasetags/delete_executor_test.go`
- **Objective:** Unit test 4-tier deletion flow, step skipping, and sidecar cleanup.
- **Test Inventory:**
  1. `TestDeleteExecutor_FullFourTierSuccess`: Verify all 4 steps execute in order and record `StepSuccess`.
  2. `TestDeleteExecutor_OrphanTagSkipsGHRelease`: Verify GitHub release step is cleanly skipped for orphan tags.
  3. `TestDeleteExecutor_SidecarMetadataClean`: Create mock file in `.gitmap/release/v1.0.0.json`, execute Step 4, and assert file no longer exists.
  4. `TestDeleteExecutor_ErrorIsolation`: Verify that if remote tag deletion fails, local tag deletion and sidecar cleanup still attempt execution and error is reported in `ErrorMessage`.

---

## 4. Definition of Done & Downstream Verification Gates

Worker 01 task completion requires satisfying all verification criteria:

- [ ] **Data Model Integrity:** `cli/cmdfixreleasetags/types.go` compiles cleanly with positive booleans and JSON tags.
- [ ] **All 5 Failure Reasons Audited:** Unit tests verify detection for `ORPHAN_TAG`, `DRAFT_RELEASE`, `MISSING_ASSETS`, `CORRUPTED_ASSETS`, `CICD_FAILED`.
- [ ] **All 3 Safety Invariants Guarded:** Unit tests verify protection for Active Version, Latest Healthy Release, and 15-Minute Grace Window.
- [ ] **3-Tier Deletion + Sidecar Cleanup:** Step 1 through Step 4 execute in exact sequence with proper error capture.
- [ ] **Relative Path Hygiene:** Zero absolute paths or `file:///` URIs in any authored files or comments.
- [ ] **Zero Git Invocations:** Subagent executes zero git state-changing commands (`git add`, `git commit`, `git push`, etc.).
