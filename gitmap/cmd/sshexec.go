package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/crypto"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

var seCommand = "se"

type seOptions struct {
	Exclude string
	Args    []string
}

func parseSEFlags(args []string) seOptions {
	fs := flag.NewFlagSet(seCommand, flag.ExitOnError)
	var opts seOptions
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude machines (comma separated)")
	fs.Parse(args)

	opts.Args = fs.Args()
	if len(opts.Args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap se [shell] <command> [--exclude m1,m2]\n")
		os.Exit(1)
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

	conns, err := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if err != nil {
		fmt.Printf("Failed to get connections: %v\n", err)
		return nil
	}

	conns = filterSSHConns(conns, opts.Exclude)
	if len(conns) == 0 {
		fmt.Println("No machines to execute on.")
		return nil
	}

	executeOnAllSSH(conns, opts.Args)
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

	var client *ssh.Client
	var err error

	if c.EncryptedPassword != "" {
		passBytes, decErr := crypto.Decrypt(c.EncryptedPassword, getEncryptionKey())
		if decErr != nil {
			fmt.Printf("%s Decrypt error: %v\n", header, decErr)
			return nil
		}
		client, err = crypto.ConnectWithPassword(c.IPAddress, c.Username, string(passBytes))
	} else if c.KeyPath != "" {
		client, err = crypto.ConnectWithKey(c.IPAddress, c.Username, c.KeyPath)
	} else {
		fmt.Printf("%s No password or key configured\n", header)
		return nil
	}

	if err != nil {
		fmt.Printf("%s Connect error: %v\n", header, err)
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

func determineSSHCommand(osType string, args []string) (string, string, bool) {
	if len(args) == 0 {
		return "", "", false
	}

	first := args[0]
	// Check if it's a known gitmap delegation command
	if first == "mkdir" || first == "cat" || first == "ssh" {
		return "", "", true
	}

	// Check if explicit shell
	if first == "ps" || first == "cmd" || first == "bash" || first == "sh" {
		if len(args) > 1 {
			return first, strings.Join(args[1:], " "), false
		}
		return first, "", false
	}

	// Default shell based on OS
	shell := "bash"
	if strings.ToLower(osType) == "windows" {
		shell = "ps"
	}
	return shell, strings.Join(args, " "), false
}

func ensureGitmapInstalled(client *ssh.Client, osType, header string) error {
	_, err := crypto.RunCommand(client, "gitmap --version", "")
	if err == nil {
		return nil // installed
	}

	fmt.Printf("%s gitmap not found, installing...\n", header)
	var installCmd string
	if strings.ToLower(osType) == "windows" {
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
	if strings.ToLower(osType) == "windows" {
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
