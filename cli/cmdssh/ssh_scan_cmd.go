package cmdssh

import (
	"context"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runSSHScanCLI(args []string) error {
	ctx := context.Background()
	InvalidateLivenessCache()

	hosts, appErr := loadAllRegisteredHosts(ctx)
	if appErr != nil {
		return appErr
	}

	if len(hosts) == 0 {
		fmt.Printf("\n  %sNo SSH machines registered. Run 'gitmap ssh join <user@ip> [alias]' to add one.%s\n\n",
			constants.ColorYellow, constants.ColorReset)
		return nil
	}

	results := probeAllHosts(ctx, hosts)
	renderScanTable(results)

	return nil
}

func probeAllHosts(ctx context.Context, hosts []store.SSHHost) []SSHHealthResult {
	var results []SSHHealthResult
	timeout := 1200 * time.Millisecond

	for _, h := range hosts {
		res := probeHostHealth(ctx, h, resolveHealthPort(h.Port), timeout)
		setCachedLiveness(h.IP+":22", res.IsOnline, res.Details)
		results = append(results, res)
	}

	return results
}
