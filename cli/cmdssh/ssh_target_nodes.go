package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func loadTargetNodes(target string) ([]db.SSHConnection, error) {
	conns, err := loadSSHConnectionsForTarget(target)
	if err != nil {
		return nil, err
	}

	if len(conns) == 0 {
		return nil, apperror.NewNotFoundError("no target SSH machines found for target: " + target)
	}

	return conns, nil
}

func checkRemoteNodeOnline(ip, header string) bool {
	isOnline, reason := CheckConnLiveness(nil, ip, 22, 0)
	if !isOnline {
		appErr := apperror.NewExecutionError(fmt.Sprintf("node %s is unreachable: %s", header, reason))
		fmt.Printf("  %s %sOffline:%s %v\n", header, constants.ColorRed, constants.ColorReset, appErr)
		return false
	}

	return true
}
