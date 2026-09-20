package cmdssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"golang.org/x/crypto/ssh"
)

// SSHExportAllOptions holds options for exporting all settings to nodes.
type SSHExportAllOptions struct {
	Target   string
	Except   string
	Exclude  string
	IsDryRun bool
	IsForce  bool
	IsHelp   bool
}

// RunSSHExportAllCLI executes the export-all command across nodes.
func RunSSHExportAllCLI(args []string) error {
	opts := parseExportAllOptions(args)
	if opts.IsHelp {
		printSSHExportAllHelp()
		return nil
	}
	bundle, err := collectExportBundle()
	if err != nil {
		return err
	}
	return executeSSHExportAll(bundle, opts)
}

func parseExportAllOptions(args []string) SSHExportAllOptions {
	opts := SSHExportAllOptions{}
	for i := 0; i < len(args); i++ {
		processExportAllArg(args[i], args, &i, &opts)
	}
	if opts.Target == "" || opts.Target == "nodes" || opts.Target == "all" {
		opts.Target = "all"
	}
	return opts
}

func processExportAllArg(arg string, args []string, idx *int, opts *SSHExportAllOptions) {
	if processExportAllFlag(arg, args, idx, opts) {
		return
	}
	if !strings.HasPrefix(arg, "-") && opts.Target == "" {
		opts.Target = arg
	}
}

func processExportAllFlag(arg string, args []string, idx *int, opts *SSHExportAllOptions) bool {
	switch {
	case arg == "--dry-run":
		opts.IsDryRun = true
		return true
	case arg == "-f" || arg == "--force":
		opts.IsForce = true
		return true
	case arg == "-h" || arg == "--help" || arg == "help":
		opts.IsHelp = true
		return true
	default:
		return processExportAllTargetFlags(arg, args, idx, opts)
	}
}

func processExportAllTargetFlags(arg string, args []string, idx *int, opts *SSHExportAllOptions) bool {
	if (arg == "-t" || arg == "--target") && *idx+1 < len(args) {
		*idx++
		opts.Target = args[*idx]
		return true
	}
	return processExportAllFilterFlags(arg, args, idx, opts)
}

func processExportAllFilterFlags(arg string, args []string, idx *int, opts *SSHExportAllOptions) bool {
	if arg == "--except" && *idx+1 < len(args) {
		*idx++
		opts.Except = args[*idx]
		return true
	}
	if arg == "--exclude" && *idx+1 < len(args) {
		*idx++
		opts.Exclude = args[*idx]
		return true
	}
	return false
}

func executeSSHExportAll(bundle *GitmapExportBundle, opts SSHExportAllOptions) error {
	conns := loadTransferConnections(sshTransferOptions{
		Target: opts.Target, Except: opts.Except, Exclude: opts.Exclude,
	})
	if len(conns) == 0 {
		fmt.Printf("  %sNo target machines found for export-all.%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}
	return exportToOnlineNodes(bundle, conns, opts)
}

func exportToOnlineNodes(bundle *GitmapExportBundle, conns []db.SSHConnection, opts SSHExportAllOptions) error {
	online, offline := partitionOnlineOffline(conns)
	printExportStartBanner(bundle, len(online), len(offline))
	if len(online) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	for _, c := range online {
		wg.Add(1)
		go exportNodeWorker(c, bundle, opts, &wg)
	}
	wg.Wait()
	printExecFinishSummary(len(online), offline)
	return nil
}

func exportNodeWorker(c db.SSHConnection, bundle *GitmapExportBundle, opts SSHExportAllOptions, wg *sync.WaitGroup) {
	defer wg.Done()
	if opts.IsDryRun {
		msg := fmt.Sprintf("(dry-run) would export %d macro(s), config, known_hosts to node", len(bundle.Macros))
		printNodeResultOutput(c.Alias, c.IPAddress, msg, nil)
		return
	}
	client, isConnected := connectSSHClient(c)
	if !isConnected {
		printNodeResultOutput(c.Alias, c.IPAddress, "(auth failed: export skipped)", nil)
		return
	}
	defer client.Close()
	deployBundleToNode(client, c, bundle, opts)
}

func deployBundleToNode(client *ssh.Client, c db.SSHConnection, bundle *GitmapExportBundle, opts SSHExportAllOptions) {
	shell := determineFallbackShell(c.OS)
	deployBundleFile(client, c, bundle, shell)
	deployConfigFile(client, c, bundle, shell)
	deployMacroFiles(client, c, bundle.Macros, shell, opts.IsForce)
	deployKnownHosts(client, c, bundle.KnownHosts, shell)
	triggerRemoteLocalBundleImport(client, shell)
	msg := fmt.Sprintf("Exported config, %d macro(s), known_hosts, and DB registry", len(bundle.Macros))
	printNodeResultOutput(c.Alias, c.IPAddress, msg, nil)
}

func deployBundleFile(client *ssh.Client, c db.SSHConnection, bundle *GitmapExportBundle, shell string) {
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	cmd := buildRemoteExportWriteCmd(".gitmap", "export_bundle.json", b64, isWindowsOS(c.OS))
	_, _ = crypto.RunCommand(client, cmd, shell)
}

func deployConfigFile(client *ssh.Client, c db.SSHConnection, bundle *GitmapExportBundle, shell string) {
	if len(bundle.Config) == 0 {
		return
	}
	data, err := json.MarshalIndent(bundle.Config, "", "  ")
	if err != nil {
		return
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	cmd := buildRemoteExportWriteCmd(".gitmap", "config.json", b64, isWindowsOS(c.OS))
	_, _ = crypto.RunCommand(client, cmd, shell)
}

func deployMacroFiles(client *ssh.Client, c db.SSHConnection, macros []macro.Macro, shell string, isForce bool) {
	for _, m := range macros {
		_ = syncSingleMacroToClient(client, c, m, shell, isForce)
	}
}

func deployKnownHosts(client *ssh.Client, c db.SSHConnection, kh string, shell string) {
	if kh == "" {
		return
	}
	b64 := base64.StdEncoding.EncodeToString([]byte(kh))
	cmd := buildRemoteExportAppendCmd(".ssh", "known_hosts", b64, isWindowsOS(c.OS))
	_, _ = crypto.RunCommand(client, cmd, shell)
}

func triggerRemoteLocalBundleImport(client *ssh.Client, shell string) {
	cmd := "gitmap ssh import-all --local-bundle 2>/dev/null || true"
	_, _ = crypto.RunCommand(client, cmd, shell)
}

func printExportStartBanner(bundle *GitmapExportBundle, onlineCount, offlineCount int) {
	fmt.Println()
	fmt.Printf("  %s● Exporting all settings and configurations to %d machine(s)%s\n",
		constants.ColorCyan, onlineCount, constants.ColorReset)
	fmt.Printf("    Payload: Config, %d Macro(s), %d Connection(s), Known Hosts\n",
		len(bundle.Macros), len(bundle.Connections))
	if offlineCount > 0 {
		fmt.Printf("    %s(%d offline machine(s) skipped)%s\n",
			constants.ColorYellow, offlineCount, constants.ColorReset)
	}
	fmt.Println()
}

func printSSHExportAllHelp() {
	fmt.Println("Export all GitMap settings, macros, config, and SSH data to remote nodes.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh export-all [nodes|all] [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: all online)")
	fmt.Println("      --except string     Exclude machines by alias, IP, or ID")
	fmt.Println("      --exclude string    Exclude machines (comma separated)")
	fmt.Println("  -f, --force             Overwrite existing remote files")
	fmt.Println("      --dry-run           Simulate export without writing remote files")
	fmt.Println("  -h, --help              Show help for export-all")
}
