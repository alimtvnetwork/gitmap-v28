package cmd

// CLI helpers for the --checkout flag, split out of clonefrom.go
// to keep that file under the project's 200-line cap.

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// applyCheckoutDefault stamps the global --checkout default onto
// every row whose own Checkout field is empty. Per-row values
// always win — this only fills in the holes.
func applyCheckoutDefault(plan *clonefrom.Plan, globalDefault string) {
	if len(globalDefault) == 0 {
		return
	}
	for i := range plan.Rows {
		if len(plan.Rows[i].Checkout) == 0 {
			plan.Rows[i].Checkout = globalDefault
		}
	}
}

// validateCheckoutFlag exits 2 when --checkout is set to anything
// other than the empty string or one of the three concrete modes.
func validateCheckoutFlag(v string) {
	if len(v) == 0 {
		return
	}
	switch v {
	case constants.CloneFromCheckoutAuto,
		constants.CloneFromCheckoutSkip,
		constants.CloneFromCheckoutForce:
		return
	}
	fmt.Fprintf(os.Stderr, constants.MsgCloneFromBadCheckoutFlag+"\n", v)
	os.Exit(2)
}
