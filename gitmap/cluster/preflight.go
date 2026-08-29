package cluster

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// PrintPreflight renders the preflight summary box using existing terminal formatting helpers.
// It returns true if the user confirmed or if autoConfirm is true.
func PrintPreflight(selector TargetSelectorType, effective []ClusterNode, command string, runRef string, autoConfirm bool) (bool, error) {
	fmt.Print(constants.ColorCyan + "==================================================================\n" + constants.ColorReset)
	fmt.Printf("%sCluster Command Preflight%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Print(constants.ColorCyan + "==================================================================\n" + constants.ColorReset)

	selectorStr := "Servers & Clients"
	if selector == ClientsOnly {
		selectorStr = "Clients Only"
	}

	fmt.Printf("  %sRun Ref:%s    %s\n", constants.ColorDim, constants.ColorReset, runRef)
	fmt.Printf("  %sCommand:%s    %s\n", constants.ColorDim, constants.ColorReset, command)
	fmt.Printf("  %sSelector:%s   %s\n", constants.ColorDim, constants.ColorReset, selectorStr)
	fmt.Printf("  %sNodes:%s      %d effective targets\n", constants.ColorDim, constants.ColorReset, len(effective))

	fmt.Print(constants.ColorCyan + "------------------------------------------------------------------\n" + constants.ColorReset)

	for _, n := range effective {
		fmt.Printf("  - %s (ID: %d, IP: %s)\n", n.ID, n.DisplayId, n.IP)
	}

	fmt.Print(constants.ColorCyan + "==================================================================\n" + constants.ColorReset)

	if autoConfirm {
		return true, nil
	}

	fmt.Printf("%sExecute command on these nodes? [y/N]: %s", constants.ColorYellow, constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	ans, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	ans = strings.TrimSpace(strings.ToLower(ans))
	isConfirmed := ans == "y" || ans == "yes"
	return isConfirmed, nil
}
