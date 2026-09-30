// Package cmdclone implements repository cloning operations.
package cmdclone

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	isFleetActive bool
	fleetMu       sync.RWMutex
)

// SetFleetCloneActive marks whether fleet-wide clone execution is running.
func SetFleetCloneActive(active bool) {
	fleetMu.Lock()
	defer fleetMu.Unlock()
	isFleetActive = active
}

// IsFleetCloneActive returns true if fleet-wide clone is running.
func IsFleetCloneActive() bool {
	fleetMu.RLock()
	defer fleetMu.RUnlock()
	return isFleetActive
}

// MaybePrintFleetCloneSuggestion outputs actionable fleet clone tips.
func MaybePrintFleetCloneSuggestion(source string) {
	if IsFleetCloneActive() {
		return
	}
	src := strings.TrimSpace(source)
	if src == "" {
		src = "gitmap.json"
	}
	fmt.Println()
	fmt.Printf("  %s💡 Fleet Clone Suggestion:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • Want to clone on all fleet nodes?  %sgitmap nodes clone %s%s\n", constants.ColorGreen, src, constants.ColorReset)
	fmt.Printf("    • Want to clone & fix across fleet?  %sgitmap nodes cfr %s%s\n", constants.ColorGreen, src, constants.ColorReset)
	fmt.Println()
}

// MaybePrintFleetCFRSuggestion outputs actionable fleet CFR tips.
func MaybePrintFleetCFRSuggestion(source string, isPub bool) {
	if IsFleetCloneActive() {
		return
	}
	src := strings.TrimSpace(source)
	cmdName := "cfr"
	if isPub {
		cmdName = "cfrp"
	}
	if src == "" {
		src = "gitmap.json"
	}
	fmt.Println()
	fmt.Printf("  %s💡 Fleet CFR Suggestion:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    • Want to run CFR on all fleet nodes?  %sgitmap nodes %s %s%s\n", constants.ColorGreen, cmdName, src, constants.ColorReset)
	fmt.Println()
}
