package cmdssh

import (
	"context"
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
	isOnline, reason := CheckConnLiveness(context.Background(), ip, 22, 0)
	if !isOnline {
		appErr := apperror.NewExecutionError(fmt.Sprintf("node %s is unreachable: %s", header, reason))
		printAppErrorWithStack(header, "Offline", appErr)
		return false
	}

	return true
}

func printAppErrorWithStack(header, prefix string, appErr *apperror.AppError) {
	if appErr == nil {
		return
	}

	fmt.Printf("  %s %s%s:%s %v\n", header, constants.ColorRed, prefix, constants.ColorReset, appErr)
	if appErr.Stack != "" {
		fmt.Printf("  %sStack Trace:%s%s\n", constants.ColorYellow, constants.ColorReset, appErr.Stack)
	}
}
