package cmdssh

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var seCommand = "se"

type seOptions struct {
	Exclude    string
	Except     string
	Target     string
	IP         string
	Args       []string
	IsShowHelp bool
}

func printSSHExecHelp() {
	fmt.Println("Execute remote commands across SSH machines with automatic liveness checks.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh exec [target] \"<command>\" [flags]")
	fmt.Println("  gitmap se [target] \"<command>\" [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: all online machines)")
	fmt.Println("      --exclude string    Exclude machines by alias or IP (comma separated)")
	fmt.Println("      --except string     Exclude machines by alias, IP, or ID (comma separated)")
	fmt.Println("      --ip string         Target machine IP address")
	fmt.Println("  -h, --help              Show help for ssh exec")
	printSSHExecExamples()
}

func printSSHExecExamples() {
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap ssh exec \"uptime\"")
	fmt.Println("  gitmap ssh exec devbox \"uname -a && df -h\"")
	fmt.Println("  gitmap ssh exec devbox \"cd /var/www && git status; ls -la\"")
	fmt.Println("  gitmap ssh exec devbox,worker-1 \"free -m\"")
	fmt.Println("  gitmap ssh exec devbox gitmap status")
	fmt.Println("  gitmap ssh exec devbox \"gitmap status && gitmap pipeline\"")
	fmt.Println("  gitmap ssh exec all gitmap --version")
	fmt.Println("  gitmap ssh exec --target devbox \"docker ps\"")
	fmt.Println("  gitmap ssh exec --exclude worker-1,192.168.1.20 \"free -m\"")
	fmt.Println("  gitmap ssh exec cmd1,cmd2,cmd3 --except worker-1,2")
}

func configureSEFlags(fs *flag.FlagSet, opts *seOptions) {
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude machines (comma separated)")
	fs.StringVar(&opts.Except, "except", "", "Exclude machines by alias, IP, or ID (comma separated)")
	fs.StringVar(&opts.Target, "target", "", "Target machine alias or IP")
	fs.StringVar(&opts.Target, "t", "", "Target machine alias or IP (shorthand)")
	fs.StringVar(&opts.IP, "ip", "", "Target machine IP address")
}

func validateSEArgs(args []string) {
	if len(args) == 0 {
		printSSHExecHelp()
		err := apperror.NewWithDetails(
			"cmd.sshexec.parseFlags",
			"E1151",
			"missing command to execute: gitmap se [target] \"<command>\"",
			"cmd.sshexec",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(err, 1)
	}
}

func parseSEFlags(args []string) seOptions {
	if hasHelpFlag(args) {
		printSSHExecHelp()
		return seOptions{IsShowHelp: true}
	}
	fs := flag.NewFlagSet(seCommand, flag.ExitOnError)
	var opts seOptions
	configureSEFlags(fs, &opts)
	fs.Parse(args)
	opts.Args = fs.Args()
	validateSEArgs(opts.Args)

	return opts
}

func resolveExcludeCSV(opts seOptions) string {
	if opts.Except != "" && opts.Exclude != "" {
		return opts.Except + "," + opts.Exclude
	}
	if opts.Except != "" {
		return opts.Except
	}
	return opts.Exclude
}

func runSSHExec(args []string) error {
	opts := parseSEFlags(args)
	if opts.IsShowHelp {
		return nil
	}

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

	conns := filterSSHConns(connsRes.Data, resolveExcludeCSV(opts))
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
		if !isConnExcluded(c, excludeList) {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

func isConnExcluded(c db.SSHConnection, excludeList []string) bool {
	userHost := fmt.Sprintf("%s@%s", c.Username, c.IPAddress)
	for _, ex := range excludeList {
		ex = strings.TrimSpace(ex)
		if ex == "" {
			continue
		}
		if strings.EqualFold(c.Alias, ex) || c.IPAddress == ex || strings.EqualFold(userHost, ex) {
			return true
		}
	}

	return false
}

func executeOnAllSSH(conns []db.SSHConnection, args []string) {
	online, offline := partitionOnlineOffline(conns)
	cmdStr := strings.Join(args, " ")
	printExecStartBanner(online, offline, cmdStr)
	if len(online) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, c := range online {
		wg.Add(1)
		go runSSHWorker(c, args, &wg)
	}

	wg.Wait()
	printExecFinishSummary(len(online), offline)
}

func runSSHWorker(c db.SSHConnection, args []string, wg *sync.WaitGroup) error {
	defer wg.Done()

	client, isConnected := connectSSHClient(c)
	if !isConnected {
		printNodeResultOutput(c.Alias, c.IPAddress, "(authentication required: configure key or password)", nil)
		return nil
	}
	defer client.Close()

	if err := ensureGitmapInstalled(client, c.OS, c.Alias); err != nil {
		printNodeResultOutput(c.Alias, c.IPAddress, "", err)
		return nil
	}

	shellType, commandStr, delegateToGitmap := determineSSHCommand(c.OS, args)
	if shellType == "ps" || shellType == "pwsh" {
		_ = ensurePowerShellInstalled(client, c.OS, c.Alias)
	}
	if delegateToGitmap {
		commandStr = resolveGitmapCommandString(args)
		shellType = ""
	}

	out, err := crypto.RunCommand(client, commandStr, shellType)
	printNodeResultOutput(c.Alias, c.IPAddress, out, err)

	return nil
}

func connectSSHClient(c db.SSHConnection, headers ...string) (*ssh.Client, bool) {
	header := ""
	if len(headers) > 0 {
		header = headers[0]
	}

	if c.EncryptedPassword != "" {
		return connectWithEncryptedPassword(c, header)
	}

	if c.KeyPath != "" {
		return connectWithKeyPath(c, header)
	}

	if client, isDefaultOk := connectWithDefaultKey(c.IPAddress, c.Username, header); isDefaultOk {
		return client, true
	}

	return nil, false
}

func connectWithKeyPath(c db.SSHConnection, header string) (*ssh.Client, bool) {
	client, err := crypto.ConnectWithKey(c.IPAddress, c.Username, c.KeyPath)
	if err != nil {
		if header != "" {
			fmt.Printf("%s Connect error: %v\n", header, err)
		}

		return nil, false
	}

	return client, true
}

func connectWithEncryptedPassword(c db.SSHConnection, header string) (*ssh.Client, bool) {
	passBytes, decErr := crypto.Decrypt(c.EncryptedPassword, getEncryptionKey())
	if decErr != nil {
		if header != "" {
			fmt.Printf("%s Decrypt error: %v\n", header, decErr)
		}

		return nil, false
	}

	client, err := crypto.ConnectWithPassword(c.IPAddress, c.Username, string(passBytes))
	if err != nil {
		if header != "" {
			fmt.Printf("%s Connect error: %v\n", header, err)
		}

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
