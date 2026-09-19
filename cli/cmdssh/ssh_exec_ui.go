package cmdssh

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/termpad"
)

func partitionOnlineOffline(conns []db.SSHConnection) ([]db.SSHConnection, []db.SSHConnection) {
	var online, offline []db.SSHConnection
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, c := range conns {
		wg.Add(1)
		go checkAndAppendConn(c, &online, &offline, &mu, &wg)
	}
	wg.Wait()

	return online, offline
}

func checkAndAppendConn(c db.SSHConnection, on, off *[]db.SSHConnection, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	isOnline, _ := CheckConnLiveness(context.Background(), c.IPAddress, 22, 0)
	mu.Lock()
	defer mu.Unlock()
	if isOnline {
		*on = append(*on, c)
		return
	}
	*off = append(*off, c)
}

func printExecStartBanner(online, offline []db.SSHConnection, cmdStr string) {
	fmt.Println()
	printOfflineAtStart(offline)
	if len(online) == 0 {
		fmt.Printf("  %sAll target machines are currently off. No nodes available to execute.%s\n\n",
			constants.ColorYellow, constants.ColorReset)
		return
	}
	printInjectedBanner(online, cmdStr)
}

func printOfflineAtStart(offline []db.SSHConnection) {
	if len(offline) == 0 {
		return
	}
	fmt.Printf("  %sNotice: The following machine(s) are currently OFF or unreachable:%s\n",
		constants.ColorYellow, constants.ColorReset)
	for _, c := range offline {
		fmt.Printf("    • [%s | %s] (machine is off)\n", c.Alias, c.IPAddress)
	}
	fmt.Println()
}

func printInjectedBanner(online []db.SSHConnection, cmdStr string) {
	var targets []string
	for _, c := range online {
		targets = append(targets, fmt.Sprintf("%s:%s", c.Alias, c.IPAddress))
	}
	fmt.Printf("  %s● Commands injected for %d machine(s) [%s]:%s\n",
		constants.ColorCyan, len(online), strings.Join(targets, ", "), constants.ColorReset)
	fmt.Printf("    %sCommand(s):%s %s\n", constants.ColorWhite, constants.ColorReset, cmdStr)
	fmt.Printf("    %sProcessing execution across active node(s)...%s\n\n",
		constants.ColorDim, constants.ColorReset)
}

func printNodeResultOutput(alias, ip, out string, err error) {
	fmt.Printf("  %s─── [%s | %s] ───%s\n", constants.ColorCyan, alias, ip, constants.ColorReset)
	if err != nil {
		fmt.Printf("    %sExecute error: %v%s\n", constants.ColorRed, err, constants.ColorReset)
		if strings.TrimSpace(out) != "" {
			fmt.Printf("    %s\n", strings.ReplaceAll(strings.TrimSpace(out), "\n", "\n    "))
		}
		fmt.Println()
		return
	}
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		fmt.Printf("    %s(no output)%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}
	fmt.Printf("    %s\n\n", strings.ReplaceAll(trimmed, "\n", "\n    "))
}

func printExecFinishSummary(onlineCount int, offline []db.SSHConnection) {
	if len(offline) > 0 {
		fmt.Printf("  %sSummary of offline machines (%d machine(s) off):%s\n",
			constants.ColorYellow, len(offline), constants.ColorReset)
		for _, c := range offline {
			fmt.Printf("    • [%s | %s] (off)\n", c.Alias, c.IPAddress)
		}
		fmt.Println()
	}
	msg := fmt.Sprintf("  %s✓ SSH Execution completed across %d active node(s).%s\n",
		constants.ColorGreen, onlineCount, constants.ColorReset)
	fmt.Print(msg)
	termpad.EnsureBottomPadding(msg)
}
