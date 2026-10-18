package cmdssh

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func startFleetPullHeartbeat(nodeCount int, isJSON bool) func() {
	if isJSON {
		return func() {}
	}
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	start := time.Now()
	isTTY := termout.IsStdoutTTY()

	go func() {
		defer close(doneCh)
		interval := 1 * time.Second
		if !isTTY {
			interval = 30 * time.Second
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return
			case now := <-ticker.C:
				elapsed := now.Sub(start)
				renderFleetHeartbeatTick(isTTY, elapsed, nodeCount)
			}
		}
	}()

	return func() {
		close(stopCh)
		<-doneCh
		if isTTY {
			fmt.Print("\r\033[K")
		}
	}
}

func renderFleetHeartbeatTick(isTTY bool, elapsed time.Duration, nodeCount int) {
	if isTTY && elapsed >= 3*time.Second {
		fmt.Printf("\r\033[K    %s[⏳ %ds elapsed] Pulling across SSH fleet (%d online nodes, 1 local)...%s",
			constants.ColorDim, int(elapsed.Seconds()), nodeCount, constants.ColorReset)
		return
	}
	if !isTTY && elapsed >= 30*time.Second {
		fmt.Printf("    %s[⏳ %ds elapsed] Pulling across SSH fleet (%d online nodes, 1 local)...%s\n",
			constants.ColorDim, int(elapsed.Seconds()), nodeCount, constants.ColorReset)
	}
}
