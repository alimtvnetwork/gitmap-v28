package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/crypto"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh"
)

var seCommand = "se"

type seOptions struct {
	Exclude string
	Shell   string
	Command string
}

func parseSEFlags(args []string) seOptions {
	fs := flag.NewFlagSet(seCommand, flag.ExitOnError)
	var opts seOptions
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude machines (comma separated)")
	fs.Parse(reorderFlagsBeforeArgs(args))

	argsAfterParse := fs.Args()
	if len(argsAfterParse) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gitmap se <shell> <command> [--exclude m1,m2]\n")
		os.Exit(1)
	}

	opts.Shell = argsAfterParse[0]
	opts.Command = argsAfterParse[1]
	return opts
}

func runSSHExec(args []string) {
	opts := parseSEFlags(args)

	dbConn, err := store.OpenDefault()
	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
		return
	}
	defer dbConn.Close()

	conns, err := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if err != nil {
		fmt.Printf("Failed to get connections: %v\n", err)
		return
	}

	conns = filterSSHConns(conns, opts.Exclude)
	if len(conns) == 0 {
		fmt.Println("No machines to execute on.")
		return
	}

	executeOnAllSSH(conns, opts.Shell, opts.Command)
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

func executeOnAllSSH(conns []db.SSHConnection, shellType, command string) {
	var wg sync.WaitGroup
	for _, c := range conns {
		wg.Add(1)
		go runSSHWorker(c, shellType, command, &wg)
	}
	wg.Wait()
	fmt.Println("SSH Execution Done.")
}

func runSSHWorker(c db.SSHConnection, shellType, command string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd")).Render(fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress))

	var client *ssh.Client
	var err error

	if c.EncryptedPassword != "" {
		passBytes, decErr := crypto.Decrypt(c.EncryptedPassword, getEncryptionKey())
		if decErr != nil {
			fmt.Printf("%s Decrypt error: %v\n", header, decErr)
			return
		}
		client, err = crypto.ConnectWithPassword(c.IPAddress, c.Username, string(passBytes))
	} else if c.KeyPath != "" {
		client, err = crypto.ConnectWithKey(c.IPAddress, c.Username, c.KeyPath)
	} else {
		fmt.Printf("%s No password or key configured\n", header)
		return
	}

	if err != nil {
		fmt.Printf("%s Connect error: %v\n", header, err)
		return
	}
	defer client.Close()

	out, err := crypto.RunCommand(client, command, shellType)
	if err != nil {
		fmt.Printf("%s Execute error: %v\n", header, err)
		return
	}

	fmt.Printf("%s\n%s\n", header, strings.TrimSpace(out))
}
