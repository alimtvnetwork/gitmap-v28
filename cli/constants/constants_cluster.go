package constants

const (
	MsgClusterPreflight = "Cluster preflight check..."
	ErrClusterNoNodes   = "no nodes available in cluster"
	ErrFilterExclusive  = "--except cannot be used with --ip or --id"
	ErrInvalidRange     = "invalid range format"

	IPv4OctetCount         = 4
	IPv4TrailingOctetIndex = 3
	RangeSeparator         = "-"
	DotSeparator           = "."

	EmptySliceLength    = 0
	RangePartsCount     = 2
	RangeStartPartIndex = 0
	RangeEndPartIndex   = 1

	NodeRoleServer = "server"
	NodeRoleAgent  = "agent"
	NodeRoleNode   = "node"

	ErrClusterServerProtected        = "lifecycle commands cannot target server nodes"
	ErrClusterLifecycleRequiresForce = "--force-lifecycle flag is required for this operation"
	ErrClusterPasswordRequired       = "password is required for this operation"
	ErrClusterInvalidPassword        = "invalid password"
	MsgClusterCountdown              = "⚠ %s %d nodes in %ds… Press Ctrl+C to abort"
	MsgClusterAbortedByUser          = "Aborted by user."
	MsgClusterAuditFooter            = "Cluster run completed"

	LifecycleCmdShutdown       = "shutdown"
	LifecycleCmdReboot         = "reboot"
	LifecycleCmdLogoff         = "logoff"
	ArgRestart                 = "/r"
	ArgShutdownWin             = "/s"
	ArgTimeout                 = "/t"
	ArgZero                    = "0"
	ArgHalt                    = "-h"
	ArgNow                     = "now"
	LifecycleCmdUnixLogoffArgs = "pkill -KILL -u whoami"

	ClusterDefaultHistoryLimit = 50
	ClusterBcryptCost          = 12

	ClusterHeaderRunRef         = "RunRef"
	ClusterHeaderCommandKind    = "CommandKind"
	ClusterHeaderTargetSelector = "TargetSelector"
	ClusterHeaderNodes          = "Nodes"
	ClusterHeaderOK             = "OK"
	ClusterHeaderFAIL           = "FAIL"
	ClusterHeaderStartedAt      = "StartedAt"
	ClusterHeaderNode           = "Node"
	ClusterHeaderSubCommand     = "SubCommand"
	ClusterHeaderResult         = "Result"
	ClusterHeaderExitCode       = "ExitCode"
	ClusterHeaderDurationMs     = "DurationMs"

	ClusterStatusOffline     = "offline"
	ClusterStatusUnreachable = "unreachable"

	FlagClusterFormat  = "--format"
	FlagClusterOutput  = "--output"
	FlagClusterJSON    = "--json"
	FlagClusterID      = "--id"
	FlagClusterConfirm = "--confirm"
	FlagClusterBefore  = "--before"
)
