package cmdssh

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunSSHPullFleetRemediate executes standalone multi-machine post-pull remediation.
func RunSSHPullFleetRemediate(args []string) error {
	fmt.Printf("\n  %sScanning multi-machine fleet for unresolved issues...%s\n",
		constants.ColorCyan, constants.ColorReset)
	conns, err := fetchAllSSHConnections()
	if err != nil {
		return err
	}
	probes := probeFleetLiveness(conns)
	cleanArgs := []string{"--json"}
	outcomes := executeFleetPullAll(probes, cleanArgs)
	issues := ExtractFleetIssues(outcomes)
	if !HasAnyFleetIssues(issues) {
		fmt.Printf("  %s✓ No unresolved issues found across multi-machine fleet.%s\n\n",
			constants.ColorGreen, constants.ColorReset)
		return nil
	}

	return runFleetRemediationCircle(issues, conns)
}
