package cluster

import "context"

func ExecUpdate(
	ctx context.Context,
	node ClusterNode,
	isAll bool,
	packages ...string,
) (string, string, int, error) {
	return "", "", 0, nil
}
