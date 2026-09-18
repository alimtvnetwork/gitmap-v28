package cmdssh

import (
	"fmt"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func runSSHCodeCLI(args []string) error {
	if len(args) == 0 {
		showSSHCodeUsage()
		return nil
	}

	target, codeArgs := parseRemoteTargetAndArgs(args)
	return executeCodeOnFleet(target, codeArgs)
}

func executeCodeOnFleet(target string, codeArgs []string) error {
	conns, err := loadTargetNodes(target)
	if err != nil {
		return err
	}

	for _, c := range conns {
		runCodeOnNode(c, codeArgs)
	}

	return nil
}

func showSSHCodeUsage() {
	fmt.Printf("\n%s Usage:%s gitmap ssh code <target> [open] <remote-path>\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Examples:")
	fmt.Println("    gitmap ssh code w1 /var/www/my-project")
	fmt.Println("    gitmap ssh code 192.168.1.8 open /home/user/workspace")
	fmt.Println()
}

func runCodeOnNode(c db.SSHConnection, codeArgs []string) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return
	}

	if tryLaunchLocalCodeRemote(header, c.Username, c.IPAddress, resolveTargetCodePath(codeArgs)) {
		return
	}

	runRemoteCodeBinary(c, header, codeArgs)
}

func resolveTargetCodePath(codeArgs []string) string {
	if len(codeArgs) == 0 {
		return "."
	}

	return codeArgs[len(codeArgs)-1]
}

func tryLaunchLocalCodeRemote(header, username, ip, path string) bool {
	remoteURI := fmt.Sprintf("ssh-remote+%s@%s", username, ip)
	cmd := exec.Command("code", "--remote", remoteURI, path)
	if err := cmd.Start(); err == nil {
		fmt.Printf("  %s %sOpening remote VS Code window for %s:%s%s\n",
			header, constants.ColorGreen, remoteURI, path, constants.ColorReset)
		return true
	}

	return false
}
