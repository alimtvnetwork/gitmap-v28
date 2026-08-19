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

	ErrClusterServerProtected        = "lifecycle commands cannot target server nodes"
	ErrClusterLifecycleRequiresForce = "--force-lifecycle flag is required for this operation"
	ErrClusterPasswordRequired       = "password is required for this operation"
	MsgClusterCountdown              = "⚠ %s %d nodes in %ds… Press Ctrl+C to abort"
	MsgClusterAuditFooter            = "Cluster run completed"
)
