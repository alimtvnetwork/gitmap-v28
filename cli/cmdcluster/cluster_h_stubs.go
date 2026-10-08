package cmdcluster

import "github.com/alimtvnetwork/gitmap-v28/cli/cluster"

func RunClusterLS(selector cluster.TargetSelectorType, args []string) error             { return nil }
func RunClusterCat(selector cluster.TargetSelectorType, args []string) error            { return nil }
func RunClusterWrite(selector cluster.TargetSelectorType, args []string) error          { return nil }
func RunClusterSetDefaultPath(selector cluster.TargetSelectorType, args []string) error { return nil }
func RunClusterSetPathAlias(selector cluster.TargetSelectorType, args []string) error   { return nil }
func RunClusterUpdate(selector cluster.TargetSelectorType, isAll bool, args []string) error {
	return nil
}

func RunClusterClone(selector cluster.TargetSelectorType, subCmd string, args []string) error {
	return nil
}
