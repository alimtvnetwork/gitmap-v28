package cmdssh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"golang.org/x/crypto/ssh"
)

// SSHImportAllOptions holds options for importing all settings from a node.
type SSHImportAllOptions struct {
	Target        string
	IsLocalBundle bool
	IsDryRun      bool
	IsForce       bool
	IsHelp        bool
}

// RunSSHImportAllCLI executes the import-all command from a node.
func RunSSHImportAllCLI(args []string) error {
	opts := parseImportAllOptions(args)
	if opts.IsHelp {
		printSSHImportAllHelp()
		return nil
	}
	if opts.IsLocalBundle {
		return executeLocalBundleImport(opts)
	}
	if opts.Target == "" {
		fmt.Println("Error: target node required for import-all (e.g. 'gitmap ssh import-all node devbox')")
		return apperror.NewValidationError("target node required for import-all")
	}
	return executeSSHImportAll(opts)
}

func parseImportAllOptions(args []string) SSHImportAllOptions {
	opts := SSHImportAllOptions{}
	for i := 0; i < len(args); i++ {
		processImportAllArg(args[i], args, &i, &opts)
	}
	return opts
}

func processImportAllArg(arg string, args []string, idx *int, opts *SSHImportAllOptions) {
	if processImportAllFlag(arg, args, idx, opts) {
		return
	}
	if arg == "node" && *idx+1 < len(args) {
		*idx++
		opts.Target = args[*idx]
		return
	}
	if !strings.HasPrefix(arg, "-") && opts.Target == "" {
		opts.Target = arg
	}
}

func processImportAllFlag(arg string, args []string, idx *int, opts *SSHImportAllOptions) bool {
	switch {
	case arg == "--local-bundle":
		opts.IsLocalBundle = true
		return true
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
		return processImportAllTargetFlag(arg, args, idx, opts)
	}
}

func processImportAllTargetFlag(arg string, args []string, idx *int, opts *SSHImportAllOptions) bool {
	if (arg == "-t" || arg == "--target") && *idx+1 < len(args) {
		*idx++
		opts.Target = args[*idx]
		return true
	}
	return false
}

func executeSSHImportAll(opts SSHImportAllOptions) error {
	conn, err := resolveExportConnection(opts.Target)
	if err != nil {
		return err
	}
	logImportProcessStart(conn.Alias, conn.IPAddress)
	client, isConnected := connectSSHClient(*conn, fmt.Sprintf("[%s]", conn.Alias))
	if isConnected == false {
		return apperror.NewExecutionError("failed to authenticate with " + conn.Alias)
	}
	defer client.Close()
	bundle, fetchErr := fetchRemoteExportBundle(client, *conn)
	if fetchErr != nil {
		return fetchErr
	}
	return applyImportBundleWithDetailedLogs(bundle, *conn, opts)
}

func fetchRemoteExportBundle(client *ssh.Client, c db.SSHConnection) (*GitmapExportBundle, error) {
	shell := determineFallbackShell(c.OS)
	cmd := buildRemoteReadBundleCmd(isWindowsOS(c.OS))
	out, err := crypto.RunCommand(client, cmd, shell)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read remote export bundle from "+c.Alias)
	}
	raw, decErr := decodeRemoteBase64Output(out)
	if decErr != nil {
		return nil, decErr
	}
	var bundle GitmapExportBundle
	if err := json.Unmarshal(raw, &bundle); err != nil {
		return nil, apperror.WrapSimple(err, "unmarshal remote export bundle")
	}
	return &bundle, nil
}

func applyImportBundleLocally(bundle *GitmapExportBundle, c db.SSHConnection, opts SSHImportAllOptions) error {
	return applyImportBundleWithDetailedLogs(bundle, c, opts)
}

func importConfigLocally(cfg map[string]any) {
	if len(cfg) == 0 {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir := filepath.Join(home, ".gitmap")
	_ = os.MkdirAll(dir, 0755)
	data, _ := json.MarshalIndent(cfg, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "config.json"), data, 0644)
}

func importConnectionsLocally(conns []db.SSHConnection) int {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return 0
	}
	defer dbConn.Close()
	count := 0
	for _, c := range conns {
		if appErr := db.InsertOrUpdateSSHConnection(dbConn.Context(), dbConn.SQL(), c); appErr == nil {
			count++
		}
	}
	return count
}

func importKnownHostsLocally(kh string) int {
	if kh == "" {
		return 0
	}
	khPath, err := DefaultKnownHostsPath()
	if err != nil {
		return 0
	}
	lines := strings.Split(strings.TrimSpace(kh), "\n")
	count := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && appendRawKnownHostEntry(khPath, trimmed) == nil {
			count++
		}
	}
	return count
}

func appendRawKnownHostEntry(path string, entry string) error {
	_ = os.MkdirAll(filepath.Dir(path), 0700)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(entry + "\n")
	return err
}

func executeLocalBundleImport(opts SSHImportAllOptions) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".gitmap", "export_bundle.json")
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return readErr
	}
	var bundle GitmapExportBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return err
	}
	c := db.SSHConnection{Alias: "local", IPAddress: "127.0.0.1"}
	return applyImportBundleLocally(&bundle, c, opts)
}

func printSSHImportAllHelp() {
	fmt.Println("Import all GitMap settings, macros, config, and SSH data from a remote node.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh import-all node \"<node name, ip, alias>\" [flags]")
	fmt.Println("  gitmap ssh import-all \"<node name, ip, alias>\" [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t, --target string     Target machine alias or IP")
	fmt.Println("  -f, --force             Overwrite existing local files")
	fmt.Println("      --dry-run           Simulate import without writing local files")
	fmt.Println("  -h, --help              Show help for import-all")
}
