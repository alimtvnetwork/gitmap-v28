package constants

// CLI flag constants for cmdpull.
const (
	FlagParallel          = "parallel"
	FlagParallelShort     = "p"
	FlagConcurrency       = "concurrency"
	FlagWorkers           = "workers"
	FlagWorker            = "worker"
	FlagWorkerShort       = "w"
	FlagHands             = "hands"
	FlagHand              = "hand"
	FlagHandShort         = "h"
	FlagWWOH              = "wwoh"
	FlagWorkerWithOneHand = "worker-with-one-hand"
	FlagAutoScale         = "auto-scale"
	FlagHighPerf          = "high-perf"
	FlagTurbo             = "turbo"
	FlagLowCPU            = "low-cpu"
	FlagConservative      = "conservative"
	FlagAutoFix           = "auto-fix"
	FlagFix               = "fix"
	FlagYes               = "yes"
	FlagYesShort          = "y"
	FlagNoFix             = "no-fix"
	FlagJSON              = "json"
	FlagStatus            = "status"
	FlagTable             = "table"
	FlagStatusTable       = "status-table"
	FlagProbe             = "probe"
	FlagProbeRepos        = "probe-repos"
	FlagOnlyAvailable     = "only-available"
	FlagRaw               = "raw"
	FlagSSH               = "ssh"
	FlagHTTPS             = "https"
	FlagVerbose           = "verbose"
	FlagVerboseShort      = "v"
	FlagGroup             = "group"
	FlagGroupShort        = "g"
	FlagAll               = "all"
	FlagAllShort          = "a"
	FlagStopOnFail        = "stop-on-fail"
	FlagStopOnFailShort   = "s"
)

// Table column header constants for cmdpull table layout rendering.
const (
	TableColHeaderRepo         = "REPO"
	TableColHeaderBranch       = "BRANCH"
	TableColHeaderLatestBranch = "LATEST BRANCH"
	TableColHeaderRelease      = "RELEASE"
	TableColHeaderCommit       = "COMMIT"
	TableColHeaderPR           = "PR"
	TableColHeaderStatus       = "STATUS"
)

// Tree glyph constants for cmdpull tree formatting.
const (
	GlyphTreeBranch       = "├── "
	GlyphTreeTerminal     = "└── "
	GlyphTreeContinuation = "│   "
	GlyphTreeIndent       = "    "
)
