package cmdssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/ghtoken"
)

// RunTokenDeployCLI coordinates deploying the GitHub token across remote SSH nodes.
func RunTokenDeployCLI(args []string) error {
	token, _, err := ghtoken.Resolve()
	hasToken := err == nil && len(token) > 0
	if !hasToken {
		return fmt.Errorf("no GitHub token available to deploy; set GH_TOKEN, login via gh auth, or gitmap token add")
	}

	targetNode, exceptNode, isDryRun := parseTokenDeployArgs(args)
	conns, err := fetchAllSSHConnections()
	hasNoConns := err != nil || len(conns) == 0
	if hasNoConns {
		printNoSSHNodesWarning()
		return nil
	}

	opts := FleetParallelOptions{
		Target:   targetNode,
		Except:   exceptNode,
		TaskName: "Git Token Deploy",
		IsDryRun: isDryRun,
	}

	RunParallelFleetExecution(conns, opts, func(c db.SSHConnection) (string, error) {
		return executeSingleNodeTokenDeploy(c, token, isDryRun)
	})

	return nil
}

func printNoSSHNodesWarning() {
	fmt.Printf("%s✖ No remote SSH nodes registered in cluster database.%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Println("  Run `gitmap ssh-join` or `gitmap ssh-join-common` to onboard nodes.")
}

func executeSingleNodeTokenDeploy(c db.SSHConnection, token string, isDryRun bool) (string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if isDryRun {
		return fmt.Sprintf("[DRY-RUN] Would deploy GitHub token to %s", c.Alias), nil
	}

	client, isConnected := connectSSHNode(c, header)
	if !isConnected {
		return "", fmt.Errorf("unable to connect to %s", header)
	}
	defer client.Close()

	osType := resolveNodeOS(client, c.OS)
	cmd := fmt.Sprintf("git config --global github.token %s", token)
	out, err := crypto.RunCommand(client, cmd, resolveRemoteShell(osType))
	reportRemoteExecution(header, "Token Deployed", out, err)

	return out, err
}

func resolveNodeOS(client *ssh.Client, fallbackOS string) string {
	probed := probeRemoteOSType(client)
	hasProbed := len(probed) > 0
	if hasProbed {
		return probed
	}

	return fallbackOS
}

func parseTokenDeployArgs(args []string) (string, string, bool) {
	var target, except string
	var isDryRun bool

	for i := 0; i < len(args); i++ {
		a := args[i]
		isTargetFlag := (a == "--target" || a == "-t" || a == "--nodes") && i+1 < len(args)
		if isTargetFlag {
			target = args[i+1]
			i++
			continue
		}
		isExceptFlag := (a == "--except" || a == "-e") && i+1 < len(args)
		if isExceptFlag {
			except = args[i+1]
			i++
			continue
		}
		if a == "--dry-run" || a == "-d" {
			isDryRun = true
		}
	}

	return target, except, isDryRun
}
