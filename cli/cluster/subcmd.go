package cluster

import "github.com/alimtvnetwork/gitmap-v28/cli/db"

type ClusterSubCommand struct {
	Kind   db.CommandKindType
	RawArg string
}
