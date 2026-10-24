# Component Spec: 319 - Package Constants & Enums DRY Restructuring

## Executive Summary
This document provides the exhaustive technical breakdown for organizing enums and constants across the 5 target packages based on deep AST scans and literal auditing by the 3 research subagents.

---

## 1. Package 1: `cli/cmdpull` Enums & Constants

### A. Dedicated Enum Definitions (`cli/cmdpull/enums/`)
1. **`PullStepType`** (`byte` with static label lookup):
   ```go
   type PullStepType byte

   const (
       PullStepInvalid PullStepType = iota
       PullStepPending
       PullStepInspecting
       PullStepFetching
       PullStepMerging
       PullStepUpToDate
       PullStepFastForward
       PullStepConflict
       PullStepError
       PullStepSkipped
   )
   ```
   - Methods: `IsValid() bool`, `String() string`, `IsComplete() bool`, `HasFailed() bool`, `Badge() string`.
2. **`PullOutcomeStatusType`** (`string`):
   - Variants: `StatusUpToDate = "up-to-date"`, `StatusUpdated = "updated"`, `StatusFastForward = "fast-forward"`, `StatusMerged = "auto-merged"`, `StatusDirty = "dirty"`, `StatusFailed = "failed"`, `StatusClean = "clean"`, `StatusLocalOnly = "local-only"`.
   - Methods: `IsValid() bool`, `String() string`, `IsDirty() bool`, `IsClean() bool`.
3. **`PullTaskStatusType`** (`string`):
   - Variants: `TaskPending = "pending"`, `TaskRunning = "running"`, `TaskCompleted = "completed"`, `TaskFailed = "failed"`, `TaskCanceled = "canceled"`.
4. **`ConflictSideType`** (`string`):
   - Variants: `ConflictSideTheirs = "theirs"`, `ConflictSideOurs = "ours"`.
5. **`RemediationModeType`** (`string`):
   - Variants: `RemediationAutoMerge = "auto-merge"`, `RemediationRebase = "rebase"`, `RemediationSendCommits = "send"`, `RemediationStashAll = "stash"`, `RemediationStep = "step"`, `RemediationSkip = "skip"`.

### B. Dedicated Constants (`cli/cmdpull/constants/`)
- **Flags (`constants_flags.go`):** `FlagParallel = "parallel"`, `FlagWorkers = "worker"`, `FlagHands = "hand"`, `FlagWWOH = "wwoh"`, `FlagAutoScale = "auto-scale"`, `FlagHighPerf = "high-perf"`, `FlagLowCPU = "low-cpu"`, `FlagStatus = "status"`, `FlagJSON = "json"`, `FlagAutoFix = "fix"`, `FlagYes = "yes"`, `FlagNoFix = "no-fix"`.
- **Table Headers (`constants_table.go`):** `TableColHeaderRepo = "REPO"`, `TableColHeaderBranch = "BRANCH"`, `TableColHeaderLatestBranch = "LATEST BRANCH"`, `TableColHeaderRelease = "RELEASE"`, `TableColHeaderCommit = "COMMIT"`, `TableColHeaderPR = "PR"`, `TableColHeaderStatus = "STATUS"`.
- **Tree Connectors (`constants_tree.go`):** `GlyphTreeBranch = "├── "`, `GlyphTreeTerminal = "└── "`, `GlyphTreeContinuation = "│   "`, `GlyphTreeIndent = "    "`.
- **Diagnostic Signatures (`constants_signatures.go`):** Detection substrings for Windows Credential Manager (`"wincredman"`), merge conflicts (`"E_MERGE_CONFLICT"`, `"conflict"`), and branch divergence (`"diverged"`).

---

## 2. Package 2: `cli/cloner` Enums & Constants

### A. Dedicated Enum Definitions (`cli/cloner/enums/`)
1. **`AuditActionType`** (`byte` with static label lookup):
   ```go
   type AuditActionType byte

   const (
       AuditActionInvalid AuditActionType = iota
       AuditActionClone    // "clone"
       AuditActionPull     // "pull"
       AuditActionCached   // "cached"
       AuditActionConflict // "conflict"
   )
   ```
   - Methods: `IsValid() bool`, `String() string`, `Marker() string` (`"+"`, `"~"`, `"="`, `"!"`).
2. **`GitUnmergedStatusType`** (`string`):
   - Variants: `StatusBothModified = "UU"`, `StatusBothAdded = "AA"`, `StatusUpdatedByUs = "UD"`, `StatusUpdatedByThem = "DU"`, `StatusBothDeleted = "DD"`.
3. **`CloneDepthType`** (`byte`):
   - Variants: `CloneDepthFull`, `CloneDepthShallow`, `CloneDepthBlobless`, `CloneDepthTreeless`.
4. **`GitTransportType`** (`string`):
   - Variants: `TransportSSH = "ssh"`, `TransportHTTPS = "https"`, `TransportGit = "git"`.

### B. Dedicated Constants (`cli/cloner/constants/`)
- **Predefined Git CLI Arguments (`constants_git_args.go`):** Immutable package-level slices eliminating heap allocations during batch loops:
  - `GitPullFastForwardArgs = []string{"pull", "--progress", "--ff-only", "--autostash"}`
  - `GitPullAutoMergeArgs = []string{"pull", "--progress", "--no-rebase", "--no-edit", "-X", "diff-algorithm=histogram", "--autostash"}`
  - `GitStatusPorcelainArgs = []string{"status", "--porcelain"}`
  - `GitAbortMergeArgs = []string{"merge", "--abort"}`
  - `GitAbortRebaseArgs = []string{"rebase", "--abort"}`
- **Signatures & Environment (`constants_signatures.go`):**
  - `EnvGitLFSNoSmudge = "GIT_LFS_SKIP_SMUDGE=1"`
  - `EnvGitSSHCommand = "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=accept-new"`
  - LFS failure patterns: `"smudge filter lfs failed"`, `"external filter 'git-lfs"`

---

## 3. Package 3: `cli/cmdclone` Enums & Constants

### A. Dedicated Enum Definitions (`cli/cmdclone/enums/`)
1. **`AuthMethodType`** (`int`):
   - Variants: `AuthMethodToken = 1`, `AuthMethodBrowser = 2`, `AuthMethodSkip = 3`.
   - Methods: `IsValid() bool`, `String() string`, `IsToken() bool`, `IsBrowser() bool`.
2. **`CloneStrategyType`** (`string`):
   - Variants: `StrategyBatch = "batch"`, `StrategySequential = "sequential"`, `StrategyInteractive = "interactive"`, `StrategyParallel = "parallel"`.
3. **`CloneCheckoutStrategyType`** (`string`):
   - Variants: `CheckoutAuto = "auto"`, `CheckoutSkip = "skip"`, `CheckoutForce = "force"`.
4. **`CloneNowOnExistsType`** (`string`):
   - Variants: `OnExistsSkip = "skip"`, `OnExistsUpdate = "update"`, `OnExistsForce = "force"`.
5. **`BatchStatusType`** (`string`):
   - Variants: `BatchStatusOK = "ok"`, `BatchStatusSkipped = "skipped"`, `BatchStatusFailed = "failed"`.
6. **`RCMenuOptionType`** (`string`):
   - Variants: `RCMenuOptionAll = "A"`, `RCMenuOptionSelect = "S"`, `RCMenuOptionRotate = "R"`, `RCMenuOptionQuit = "Q"`.

### B. Dedicated Constants (`cli/cmdclone/constants/`)
- **Flags (`constants_flags.go`):** `FlagMaxConcurrency`, `FlagWorkers`, `FlagParallel`, `FlagTargetDir`, `FlagOnly`, `FlagExclude`, `FlagToken`, `FlagSSHKey`, `FlagClean`, `FlagForce`, `FlagReclone`, `FlagSafePull`, `FlagRepoCache`.
- **UI Banners & Prompts (`constants_messages.go`):** `MsgMultiCloneNoURLs`, `MsgMultiCloneDryRunTitle`, `MsgRCOptionHeader`, `MsgRCOptionAll`, `MsgRCOptionSelect`, `MsgAuthRequiredFmt`.

---

## 4. Package 4: `cli/store` Enums & Constants

### A. Dedicated Enum Definitions (`cli/store/enums/`)
1. **`StoreTableType`** (`string`):
   - Variants: `TableRepo = "Repo"`, `TableGroup = "Group"`, `TablePipelineRun = "PipelineRun"`, `TableSplitDatabaseRegistry = "SplitDatabaseRegistry"`, `TableInstallerPins = "installer_pins"`, `TableSiteRegistry = "SiteRegistry"`, `TableSSHHosts = "ssh_hosts"`, `TableParentTask = "ParentTask"`, `TableSubtask = "Subtask"`, `TableCompletedTask = "CompletedTask"`, `TablePendingTask = "PendingTask"`.
   - Methods: `IsValid() bool`, `String() string`, `TableName() string`.
2. **`SQLitePragmaType`** (`string`):
   - Variants: `PragmaJournalWAL = "PRAGMA journal_mode = WAL;"`, `PragmaSynchronousNormal = "PRAGMA synchronous = NORMAL;"`, `PragmaBusyTimeout5s = "PRAGMA busy_timeout = 5000;"`, `PragmaForeignKeysOn = "PRAGMA foreign_keys = ON;"`.
3. **`SplitDatabaseType`** (`string`):
   - Variants: `SplitDBTypeRepodb = "repodb"`, `SplitDBTypeSearch = "search"`, `SplitDBTypeTasks = "tasks"`, `SplitDBTypeAiInstruction = "ai-instruction"`, `SplitDBTypeAiAnalysis = "ai_analysis"`, `SplitDBTypeStartup = "startup"`.
4. **`SubtaskStatusType`** (`string`):
   - Variants: `SubtaskPending = "PENDING"`, `SubtaskInProgress = "IN_PROGRESS"`, `SubtaskCompleted = "COMPLETED"`, `SubtaskFailed = "FAILED"`.
5. **`TxOperationType`** (`string`):
   - Variants: `TxOpBegin = "BEGIN"`, `TxOpCommit = "COMMIT"`, `TxOpRollback = "ROLLBACK"`.

### B. Dedicated Constants (`cli/store/constants/`)
- **Driver & DSN:** `DriverSQLite = "sqlite"`, `PragmaBusyTimeoutDSN = "?_pragma=busy_timeout(2000)"`, `DefaultBusyTimeoutMs = 5000`, `DefaultSSHPort = 22`.
- **Pragma Statements:** Predefined immutable array of default pragmas eliminating 6 duplicate slice allocations on every database open.

---

## 5. Rollout Package 5: `cli/cmdcommit` Enums & Constants

### A. Dedicated Enum Definitions (`cli/cmdcommit/enums/`)
1. **`CommitVerbType`** (`string`):
   - Variants: `CommitVerbFeature = "cpf"`, `CommitVerbBug = "cbf"`, `CommitVerbBugAlt = "cpb"`, `CommitVerbRefactor = "cfr"`, `CommitVerbChore = "cpc"`, `CommitVerbRelease = "cpr"`, `CommitVerbFlat = "cm"`, `CommitVerbPullPush = "pcp"`.
   - Methods: `IsValid() bool`, `String() string`, `Prefix() string`, `UppercaseToken() string`, `IsFeature() bool`, `IsBug() bool`.
2. **`ConventionalVerbType`** (`string`):
   - Variants: `ConventionalFeat = "FEAT"`, `ConventionalFix = "FIX"`, `ConventionalChore = "CHORE"`, `ConventionalRefactor = "REFACTOR"`, `ConventionalRelease = "RELEASE"`, `ConventionalTest = "TEST"`, `ConventionalDocs = "DOCS"`.
3. **`CommitFormatType`** (`byte`):
   - Variants: `FormatStandardBracket`, `FormatHyphenSeparated`, `FormatConventional`.
4. **`OriginType`** (`string`):
   - Variants: `OriginLocalAI = "LocalAI"`, `OriginForeignAI = "ForeignAI"`, `OriginHumanOrExternal = "HumanOrExternal"`.
5. **`GitStepType`** (`string`):
   - Variants: `GitStepAdd = "add"`, `GitStepCommit = "commit"`, `GitStepPush = "push"`, `GitStepStatus = "status"`, `GitStepReset = "reset"`.

### B. Dedicated Constants (`cli/cmdcommit/constants/`)
- **Flags:** `FlagAgentID = "--agent-id"`, `FlagHeatmapID = "--heatmap-id"`, `FlagAI = "--ai"`, `FlagSpec = "--spec"`, `FlagTask = "--task"`, `FlagPush = "--push"`, `FlagDryRun = "--dry-run"`.
- **Environment Variables:** `EnvGitmapAgentID = "GITMAP_AGENT_ID"`, `EnvAntigravityAgentID = "ANTIGRAVITY_AGENT_ID"`, `EnvGitmapHeatmapID = "GITMAP_HEATMAP_ID"`, `EnvAIModel = "AI_MODEL"`.
- **Patterns:** `PatternHyphen = `^([A-Za-z0-9_\-\./]+)\s*[:-]\s*(.+)$``, `PatternConventional = `^([A-Za-z]+)\(([A-Za-z0-9_\-\./]+)\):\s*(.+)$``.
