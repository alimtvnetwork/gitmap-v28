package cluster

import "github.com/alimtvnetwork/gitmap-v28/cli/result"

type AliasEntry struct {
	Alias string
	Path  string
}

func ParseSetPathAliasArg(raw string) result.ResultSlice[AliasEntry] {
	return result.OkSlice([]AliasEntry{})
}
