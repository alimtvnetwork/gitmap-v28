package cmdssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func runSSHInstallCLI(args []string) error {
	target := "all"
	if len(args) > 0 && args[0] != "gitmap" {
		target = args[0]
	} else if len(args) > 1 {
		target = args[1]
	}

	fmt.Printf("\n%s Installing / Updating GitMap across SSH fleet (%s):%s\n\n",
		constants.ColorCyan, target, constants.ColorReset)

	return executeSSHInstallOnTarget(target)
}

func executeSSHInstallOnTarget(target string) error {
	conns, err := loadSSHConnectionsForTarget(target)
	if err != nil {
		return err
	}

	for _, c := range conns {
		installOrUpdateSingleNode(c)
	}

	fmt.Println("\nSSH Install & Update complete.")

	return nil
}

func installOrUpdateSingleNode(c db.SSHConnection) {
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

	if isGitmapPresent(client) {
		updateExistingGitmap(client, header, c.OS)
		return
	}

	installFreshGitmap(client, header, c.OS)
}

func isGitmapPresent(client *ssh.Client) bool {
	out, err := crypto.RunCommand(client, "gitmap --version", "")

	return err == nil && strings.Contains(out, "gitmap")
}
