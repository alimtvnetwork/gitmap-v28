package cluster

import "github.com/alimtvnetwork/gitmap-v28/gitmap/db"

type ClusterSubCommand struct {
	Kind   db.CommandKindType
	RawArg string
}
