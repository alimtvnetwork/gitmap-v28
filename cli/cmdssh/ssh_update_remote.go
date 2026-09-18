package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func runSSHUpdateCLI(args []string) error {
	target := "all"
	if len(args) > 0 && args[0] != "gitmap" {
		target = args[0]
	} else if len(args) > 1 {
		target = args[1]
	}

	fmt.Printf("\n%s Updating GitMap across SSH fleet (%s):%s\n\n",
		constants.ColorCyan, target, constants.ColorReset)

	conns, err := loadSSHConnectionsForTarget(target)
	if err != nil {
		return err
	}

	for _, c := range conns {
		updateSingleSSHNode(c)
	}

	fmt.Println("\nSSH Fleet Update complete.")

	return nil
}

func updateSingleSSHNode(c db.SSHConnection) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)

	isOnline, reason := CheckConnLiveness(nil, c.IPAddress, 22, 0)
	if !isOnline {
		fmt.Printf("  %s %sOFFLINE (skipped: %s)%s\n", header, constants.ColorYellow, reason, constants.ColorReset)
		return
	}

	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return
	}
	defer client.Close()

	shell := "bash"
	if isWindowsOS(c.OS) {
		shell = "ps"
	}

	out, err := crypto.RunCommand(client, "gitmap update", shell)
	if err != nil {
		fmt.Printf("  %s %sUpdate failed:%s %v\n", header, constants.ColorRed, constants.ColorReset, err)
		return
	}

	fmt.Printf("  %s %sUpdated:%s\n%s\n", header, constants.ColorGreen, constants.ColorReset, strings.TrimSpace(out))
}
