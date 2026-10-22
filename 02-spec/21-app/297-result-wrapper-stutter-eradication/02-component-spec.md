# Component Specification 297: Result Wrapper Component Signatures & Types

## 1. Subsystem Specifications

### 1.1 `cli/cmdagy/`
- **`cli/cmdagy/types.go` (new):**
  ```go
  package cmdagy

  import (
  	"github.com/alimtvnetwork/gitmap-v28/cli/result"
  )

  type (
  	AgyProcessInfoWrap = result.Wrap[AgyProcessInfo]
  	AgyConvIDWrap      = result.Wrap[string]
  )
  ```
- **`cli/cmdagy/agy_agentapi_ops.go`:**
  - `AgentAPISendMessage(...) result.StringResult`
  - `parseSendMessageResponse(...) result.StringResult`
  - `AgentAPINewConversation(...) result.StringResult`
  - `AgentAPINewConversationWithOptions(...) result.StringResult`
  - `AgentAPINewConversationInDir(...) result.StringResult`
  - `parseNewConversationResponse(...) result.StringResult`
- **`cli/cmdagy/agy_open_paths.go`:**
  - `ResolveAntigravityIDE() result.StringResult`
  - `ResolveAntigravityCLI() result.StringResult`
  - `ResolveAntigravityBinary() result.StringResult`
  - `DetectRunningAntigravityIDE() AgyProcessInfoWrap`
  - `checkCandidatePaths(candidates []string) result.StringResult`
  - `findBinaryInPath(names []string) result.StringResult`

### 1.2 `cli/prdb/`
- **`cli/prdb/pr_split_conn.go`:**
  - `OpenPrSplitDb(repoSlug, repoRoot string) PrSplitDbWrap`
  - `OpenPrSplitDbAt(dbPath string) PrSplitDbWrap`
  - `openPrSplitDbWithSlug(dbPath, slug string) PrSplitDbWrap`
  - `initPrSplitConn(conn *sql.DB, repoSlug, dbPath string) PrSplitDbWrap`
- **`cli/prdb/pr_split_ops.go`:**
  - `CreatePullRequest(rec PullRequestRecord) result.Int64Result`
  - `scanPullRequestRecord(row *sql.Row) PullRequestWrap`
  - `GetPullRequestByNumber(prNum int) PullRequestWrap`
  - `UpdatePullRequestStatus(prNum int, status, mergeSha string) result.BoolResult`
  - `AddPrRelease(rec PrReleaseRecord) result.Int64Result`
  - `UpsertPrBranch(rec PrBranchRecord) result.BoolResult`
  - `scanSinglePrBranch(rows *sql.Rows) PrBranchRecordWrap`
  - `MarkPrBranchDeleted(branchName string) result.BoolResult`

### 1.3 `cli/osuser/` & `cli/osclean/`
- **`cli/osuser/group_portable.go`:**
  - `ExportGroup(name string) PortableGroupWrap`
  - `ImportGroup(pg PortableGroup) result.BoolResult`
  - `ImportAllGroups(groups []PortableGroup) result.IntResult`
  - `EditGroup(oldName, newName string) result.BoolResult`
- **`cli/osuser/portable.go`:**
  - `ExportUser(username string) PortableUserWrap`
  - `ImportUser(pu PortableUser) result.BoolResult`
  - `ImportAllUsers(users []PortableUser) result.IntResult`
  - `EditUser(username, newShell, newHome string) result.BoolResult`
  - `applyUserEditsLinux(...) result.BoolResult`
- **`cli/osclean/dev_clean.go`:**
  - `CleanDevCaches(opts DevCleanOptions) DevCleanSummaryWrap`
- **`cli/osclean/terminal_clean.go`:**
  - `CleanTerminalHistory(opts TerminalCleanOptions) TerminalCleanSummaryWrap`

### 1.4 `cli/cmdautomation/`
- Replace `result.Result[T]` with `result.Wrap[T]` across `phase1_types.go`, `phase2_types.go`, `phase3_types.go`, `phase4_types.go`, `phase5_types.go`, `phase6_types.go`, and `types.go`.
- Replace `OpenAutomationDB` and `GetCachedRuntime` in `runtime_db.go` with canonical `result.Wrap[*sql.DB]` and `RuntimeRecordWrap`.

### 1.5 `cli/committransfer/`
- **`cli/committransfer/pr_engine.go`:**
  - `runGit(...) result.StringResult`
  - `replaySingleIngress(...) PRCommitInfoWrap`
  - `recordPRInDatabase(...) PrSplitDbWrap`
  - `mergeFeatureBranch(...) result.StringResult`
  - `executePRReplayAndMerge(...) result.StringResult`
  - `ProcessPR(...) result.StringResult`
