package cmd

import "github.com/alimtvnetwork/gitmap-v28/gitmap/model"

// ResolvedRepo describes a repository found by the unified resolver.
type ResolvedRepo struct {
	Record         model.ScanRecord
	MatchedBy      string
	OriginalTarget string
}
