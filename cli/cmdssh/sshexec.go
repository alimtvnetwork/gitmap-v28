package cmdssh

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var seCommand = "se"

type seOptions struct {
	Exclude string
	Target  string
	IP      string
	Args    []string
}

func parseSEFlags(args []string) seOptions {
	fs := flag.NewFlagSet(seCommand, flag.ExitOnError)
	var opts seOptions
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude machines (comma separated)")
	fs.StringVar(&opts.Target, "target", "", "Target machine alias or IP")
	fs.StringVar(&opts.Target, "t", "", "Target machine alias or IP (shorthand)")
	fs.StringVar(&opts.IP, "ip", "", "Target machine IP address")
	fs.Parse(args)

	opts.Args = fs.Args()
	if len(opts.Args) == 0 {
		err := apperror.NewWithDetails(
			"cmd.sshexec.parseFlags",
			"E1151",
			"Usage: gitmap se [shell] <command> [--target <alias|ip>] [--exclude m1,m2]",
			"cmd.sshexec",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
	}

	return opts
}

func runSSHExec(args []string) error {
	opts := parseSEFlags(args)

	dbConn, err := store.OpenDefault()
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)

		return nil
	}

	defer dbConn.Close()

	connsRes := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if connsRes.IsFailure() {
		fmt.Printf("Failed to get connections: %v\n", connsRes.AppError())

		return nil
	}

	conns := filterSSHConns(connsRes.Data, opts.Exclude)
	conns, execArgs := resolveExecTargetAndArgs(conns, opts)
	if len(conns) == 0 {
		fmt.Println("No machines to execute on.")

		return nil
	}

	executeOnAllSSH(conns, execArgs)

	return nil
}

func filterSSHConns(conns []db.SSHConnection, excludeCSV string) []db.SSHConnection {
	if excludeCSV == "" {
		return conns
	}

	excludeList := strings.Split(excludeCSV, ",")
	var filtered []db.SSHConnection
	for _, c := range conns {
		excluded := false
		for _, ex := range excludeList {
			ex = strings.TrimSpace(ex)
			if c.Alias == ex || c.IPAddress == ex {
				excluded = true
				break
			}
		}

		if !excluded {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

func executeOnAllSSH(conns []db.SSHConnection, args []string) {
	var wg sync.WaitGroup
	for _, c := range conns {
		wg.Add(1)
		go runSSHWorker(c, args, &wg)
	}

	wg.Wait()
	fmt.Println("SSH Execution Done.")
}

func runSSHWorker(c db.SSHConnection, args []string, wg *sync.WaitGroup) error {
	defer wg.Done()

	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd")).Render(fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress))

	isOnline, reason := CheckConnLiveness(context.Background(), c.IPAddress, 22, 0)
	if !isOnline {
		fmt.Printf("%s OFFLINE (skipped: %s)\n", header, reason)

		return nil
	}

	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return nil
	}

	defer client.Close()

	if err := ensureGitmapInstalled(client, c.OS, header); err != nil {
		fmt.Printf("%s Failed to ensure gitmap: %v\n", header, err)

		return nil
	}

	shellType, commandStr, delegateToGitmap := determineSSHCommand(c.OS, args)

	if shellType == "ps" || shellType == "pwsh" {
		_ = ensurePowerShellInstalled(client, c.OS, header)
	}

	if delegateToGitmap {
		commandStr = "gitmap " + strings.Join(args, " ")
		shellType = "" // default shell
	}

	out, err := crypto.RunCommand(client, commandStr, shellType)
	if err != nil {
		fmt.Printf("%s Execute error: %v\n%s\n", header, err, strings.TrimSpace(out))

		return nil
	}

	fmt.Printf("%s\n%s\n", header, strings.TrimSpace(out))

	return nil
}

func connectSSHClient(c db.SSHConnection, header string) (*ssh.Client, bool) {
	if c.EncryptedPassword != "" {
		return connectWithEncryptedPassword(c, header)
	}

	if c.KeyPath != "" {
		return connectWithKeyPath(c, header)
	}

	if client, isDefaultOk := connectWithDefaultKey(c.IPAddress, c.Username, header); isDefaultOk {
		return client, true
	}

	fmt.Printf("%s %s\n", header, formatMissingAuthAdvice(c.Alias, c.IPAddress, c.Username))

	return nil, false
}

func connectWithKeyPath(c db.SSHConnection, header string) (*ssh.Client, bool) {
	client, err := crypto.ConnectWithKey(c.IPAddress, c.Username, c.KeyPath)
	if err != nil {
		fmt.Printf("%s Connect error: %v\n", header, err)

		return nil, false
	}

	return client, true
}

func connectWithEncryptedPassword(c db.SSHConnection, header string) (*ssh.Client, bool) {
	passBytes, decErr := crypto.Decrypt(c.EncryptedPassword, getEncryptionKey())
	if decErr != nil {
		fmt.Printf("%s Decrypt error: %v\n", header, decErr)

		return nil, false
	}

	client, err := crypto.ConnectWithPassword(c.IPAddress, c.Username, string(passBytes))
	if err != nil {
		fmt.Printf("%s Connect error: %v\n", header, err)

		return nil, false
	}

	return client, true
}

func ensureGitmapInstalled(client *ssh.Client, osType, header string) error {
	_, err := crypto.RunCommand(client, "gitmap --version", "")
	if err == nil {
		return nil // installed
	}

	fmt.Printf("%s gitmap not found, installing...\n", header)
	var installCmd string
	if strings.EqualFold(osType, "windows") {
		installCmd = "irm https://gitmap.dev/install.ps1 | iex"
		_, err = crypto.RunCommand(client, installCmd, "ps")
	} else {
		installCmd = "curl -fsSL https://gitmap.dev/install.sh | bash"
		_, err = crypto.RunCommand(client, installCmd, "bash")
	}

	if err != nil {
		return fmt.Errorf("auto-install failed: %w", err)
	}

	return nil
}

func ensurePowerShellInstalled(client *ssh.Client, osType, header string) error {
	if strings.EqualFold(osType, "windows") {
		return nil
	}

	_, err := crypto.RunCommand(client, "pwsh --version", "bash")
	if err == nil {
		return nil
	}

	fmt.Printf("%s PowerShell not found, installing via package manager...\n", header)
	installCmd := `if command -v apt-get &> /dev/null; then sudo apt-get update && sudo apt-get install -y powershell; elif command -v yum &> /dev/null; then sudo yum install -y powershell; elif command -v brew &> /dev/null; then brew install --cask powershell; fi`
	_, err = crypto.RunCommand(client, installCmd, "bash")
	if err != nil {
		fmt.Printf("%s Note: auto-installing PowerShell failed. It may require manual setup.\n", header)
	}

	return nil
}
