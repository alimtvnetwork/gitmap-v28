package cluster

import "context"

func ExecCloneURL(
	ctx context.Context,
	node ClusterNode,
	url,
	workdir,
	subCmd string,
) (string, string, int, error) {
	return "", "", 0, nil
}
