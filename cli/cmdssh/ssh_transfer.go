package cmdssh

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type sshTransferOptions struct {
	Target  string
	Except  string
	Exclude string
	Args    []string
	IsHelp  bool
}

func parseTransferFlags(cmdName string, args []string) sshTransferOptions {
	if hasHelpFlag(args) || (len(args) > 0 && args[0] == "help") {
		return sshTransferOptions{IsHelp: true}
	}

	fs := flag.NewFlagSet(cmdName, flag.ExitOnError)
	var opts sshTransferOptions
	fs.StringVar(&opts.Target, "target", "", "Target machine alias or IP")
	fs.StringVar(&opts.Target, "t", "", "Target machine alias or IP (shorthand)")
	fs.StringVar(&opts.Except, "except", "", "Exclude machines by alias, IP, or ID")
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude machines (comma separated)")
	fs.Parse(args)
	opts.Args = fs.Args()

	return opts
}

func resolveTransferExclusions(opts sshTransferOptions) string {
	if opts.Except != "" && opts.Exclude != "" {
		return opts.Except + "," + opts.Exclude
	}
	if opts.Except != "" {
		return opts.Except
	}

	return opts.Exclude
}

func runSSHCopyCLI(args []string) error {
	return runSSHTransferCLI(args, false)
}

func runSSHMvCLI(args []string) error {
	return runSSHTransferCLI(args, true)
}

func runSSHTransferCLI(args []string, isMove bool) error {
	cmdName := "ssh-copy"
	if isMove {
		cmdName = "ssh-mv"
	}

	opts := parseTransferFlags(cmdName, args)
	if opts.IsHelp || len(opts.Args) < 2 {
		printSSHTransferHelp(isMove)
		return nil
	}

	return executeSSHTransfer(opts, isMove)
}

func executeSSHTransfer(opts sshTransferOptions, isMove bool) error {
	srcPath := ExpandUniversalPath(opts.Args[0], runtime.GOOS)
	fileData, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Printf("  %sError: source file not found: %s (%v)%s\n", constants.ColorRed, srcPath, err, constants.ColorReset)
		return nil
	}

	conns := loadTransferConnections(opts)
	if len(conns) == 0 {
		fmt.Printf("  %sNo target machines available for transfer.%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	transferToAllNodes(conns, srcPath, opts.Args[1], fileData, isMove)

	return nil
}

func loadTransferConnections(opts sshTransferOptions) []db.SSHConnection {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer dbConn.Close()

	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		return nil
	}

	filtered := filterSSHConns(res.Data, resolveTransferExclusions(opts))
	if opts.Target != "" && opts.Target != "all" {
		return filterConnectionsByTarget(filtered, opts.Target)
	}

	return filtered
}

func transferToAllNodes(conns []db.SSHConnection, srcPath, rawDest string, data []byte, isMove bool) {
	online, offline := partitionOnlineOffline(conns)
	printTransferStartBanner(online, offline, srcPath, rawDest, isMove)
	if len(online) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, c := range online {
		wg.Add(1)
		go transferWorker(c, srcPath, rawDest, data, &wg)
	}
	wg.Wait()

	finalizeTransfer(srcPath, len(online), offline, isMove)
}

func transferWorker(c db.SSHConnection, srcPath, rawDest string, data []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	destPath := resolveRemoteDestPath(srcPath, rawDest, c.OS)
	client, isConnected := connectSSHClient(c)
	if !isConnected {
		printNodeResultOutput(c.Alias, c.IPAddress, "(authentication required - transfer skipped)", nil)
		return
	}
	defer client.Close()

	cmd := buildRemoteWriteCmd(destPath, data, isWindowsOS(c.OS))
	shell := determineFallbackShell(c.OS)
	out, err := crypto.RunCommand(client, cmd, shell)
	if err == nil {
		out = fmt.Sprintf("Transferred %d bytes -> %s", len(data), destPath)
	}
	printNodeResultOutput(c.Alias, c.IPAddress, out, err)
}

func resolveRemoteDestPath(srcPath, rawDest, osType string) string {
	expanded := ExpandUniversalPath(rawDest, osType)
	base := filepath.Base(srcPath)
	if !isDirDestination(rawDest) {
		return expanded
	}

	return joinRemotePath(expanded, base, isWindowsOS(osType))
}

func isDirDestination(rawDest string) bool {
	return strings.HasSuffix(rawDest, "/") || strings.HasSuffix(rawDest, "\\") || filepath.Ext(rawDest) == ""
}

func joinRemotePath(dir, base string, isWin bool) string {
	if isWin {
		return filepath.Join(dir, base)
	}

	return strings.TrimRight(dir, "/") + "/" + base
}

func buildRemoteWriteCmd(destPath string, data []byte, isWin bool) string {
	b64 := base64.StdEncoding.EncodeToString(data)
	if isWin {
		return fmt.Sprintf(`powershell -NoProfile -Command "$d=[Convert]::FromBase64String('%s'); $p='%s'; $dir=[IO.Path]::GetDirectoryName($p); if ($dir -and -not (Test-Path $dir)) { [IO.Directory]::CreateDirectory($dir) | Out-Null }; [IO.File]::WriteAllBytes($p, $d)"`, b64, destPath)
	}

	return fmt.Sprintf(`mkdir -p "$(dirname '%s')" && printf '%%s' '%s' | base64 -d > '%s'`, destPath, b64, destPath)
}

func printTransferStartBanner(online, offline []db.SSHConnection, src, dest string, isMove bool) {
	action := "Copying"
	if isMove {
		action = "Moving"
	}

	fmt.Println()
	printOfflineAtStart(offline)
	fmt.Printf("  %s● %s '%s' -> '%s' across %d machine(s):%s\n",
		constants.ColorCyan, action, src, dest, len(online), constants.ColorReset)
	fmt.Printf("    %sProcessing remote transfer...%s\n\n", constants.ColorDim, constants.ColorReset)
}

func finalizeTransfer(srcPath string, onlineCount int, offline []db.SSHConnection, isMove bool) {
	if isMove && onlineCount > 0 {
		_ = os.Remove(srcPath)
		fmt.Printf("  %s✓ Local source file '%s' removed.%s\n", constants.ColorGreen, srcPath, constants.ColorReset)
	}

	printExecFinishSummary(onlineCount, offline)
}

func printSSHTransferHelp(isMove bool) {
	action := "copy"
	if isMove {
		action = "mv"
	}

	fmt.Println()
	fmt.Printf("Transfer files to remote SSH machines with macro expansion and node filtering.\n\n")
	fmt.Printf("Usage:\n")
	fmt.Printf("  gitmap ssh %s <from-path> <to-path> [flags]\n\n", action)
	fmt.Printf("Flags:\n")
	fmt.Printf("  -t, --target string     Target machine alias or IP (default: all online machines)\n")
	fmt.Printf("      --except string     Exclude machines by alias, IP, or ID (comma separated)\n")
	fmt.Printf("      --exclude string    Exclude machines (comma separated)\n\n")
	fmt.Printf("Path Macros:\n")
	fmt.Printf("  ~             User home directory\n")
	fmt.Printf("  %%win%%         Windows folder (C:\\Windows)\n")
	fmt.Printf("  %%win-drive%%   Windows system drive (C:)\n")
	fmt.Printf("  %%temp%%        Temporary directory\n")
	fmt.Printf("  %%appdata%%     Application data directory\n\n")
	fmt.Printf("Examples:\n")
	fmt.Printf("  gitmap ssh %s ./app.tar.gz /var/www/\n", action)
	fmt.Printf("  gitmap ssh %s ./config.json .\\test\n", action)
	fmt.Printf("  gitmap ssh %s ~/setup.sh ~\n", action)
	fmt.Printf("  gitmap ssh %s ./file.txt .\\test --except worker-2,192.168.1.20\n", action)
	fmt.Printf("  gitmap ssh %s ./migration.sql /tmp/ --except alias\n", action)
	fmt.Printf("  gitmap ssh %s ./app.zip %%win%%\\Temp\n\n", action)
}
