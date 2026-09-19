package cmdssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"golang.org/x/crypto/ssh"
)

// RunSSHMacroSyncCLI executes the ssh macro-sync command.
func RunSSHMacroSyncCLI(args []string) error {
	return runSSHMacroSyncCLI(args)
}

func runSSHMacroSyncCLI(args []string) error {
	opts := parseSSHMacroSyncOptions(args)
	if opts.IsHelp {
		printSSHMacroSyncHelp()
		return nil
	}
	macros, appErr := collectLocalMacrosForSync(opts.Name)
	if appErr != nil {
		return appErr
	}
	return executeSSHMacroSync(macros, opts)
}

func parseSSHMacroSyncOptions(args []string) SSHMacroSyncOptions {
	opts := SSHMacroSyncOptions{}
	for i := 0; i < len(args); i++ {
		processSyncArg(args[i], args, &i, &opts)
	}
	return opts
}

func processSyncArg(arg string, args []string, index *int, opts *SSHMacroSyncOptions) {
	if processSyncFlag(arg, args, index, opts) {
		return
	}
	if !strings.HasPrefix(arg, "-") && opts.Name == "" {
		opts.Name = arg
	}
}

func processSyncFlag(arg string, args []string, index *int, opts *SSHMacroSyncOptions) bool {
	switch {
	case arg == "--force" || arg == "-f":
		opts.IsForce = true
		return true
	case arg == "--dry-run":
		opts.IsDryRun = true
		return true
	case arg == "-h" || arg == "--help" || arg == "help":
		opts.IsHelp = true
		return true
	default:
		return processSyncTargetFlag(arg, args, index, opts)
	}
}

func processSyncTargetFlag(arg string, args []string, index *int, opts *SSHMacroSyncOptions) bool {
	if (arg == "-t" || arg == "--target") && *index+1 < len(args) {
		*index++
		opts.Target = args[*index]
		return true
	}
	return processSyncFilterFlag(arg, args, index, opts)
}

func processSyncFilterFlag(arg string, args []string, index *int, opts *SSHMacroSyncOptions) bool {
	switch {
	case arg == "--except" && *index+1 < len(args):
		*index++
		opts.Except = args[*index]
		return true
	case arg == "--exclude" && *index+1 < len(args):
		*index++
		opts.Exclude = args[*index]
		return true
	default:
		return false
	}
}

func collectLocalMacrosForSync(name string) ([]macro.Macro, *apperror.AppError) {
	if name == "" || name == "all" || name == "*" {
		res := macro.ListMacros()
		if res.IsFailure() {
			return nil, res.AppError()
		}
		return res.Data, nil
	}
	m, err := macro.LoadMacro(name)
	if err != nil {
		return nil, apperror.WrapSimple(err, "macro.LoadMacro: "+name)
	}
	return []macro.Macro{*m}, nil
}

func executeSSHMacroSync(macros []macro.Macro, opts SSHMacroSyncOptions) error {
	if len(macros) == 0 {
		fmt.Printf("  %sNo local macros found to sync.%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}
	conns := loadTransferConnections(sshTransferOptions{
		Target: opts.Target, Except: opts.Except, Exclude: opts.Exclude,
	})
	if len(conns) == 0 {
		fmt.Printf("  %sNo target machines available for macro sync.%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}
	return syncMacrosToNodes(macros, conns, opts)
}

func syncMacrosToNodes(macros []macro.Macro, conns []db.SSHConnection, opts SSHMacroSyncOptions) error {
	online, offline := partitionOnlineOffline(conns)
	printSyncStartBanner(len(macros), len(online), len(offline))
	if len(online) == 0 {
		return nil
	}
	var wg sync.WaitGroup
	for _, c := range online {
		wg.Add(1)
		go syncNodeWorker(c, macros, opts, &wg)
	}
	wg.Wait()
	printExecFinishSummary(len(online), offline)
	return nil
}

func syncNodeWorker(c db.SSHConnection, macros []macro.Macro, opts SSHMacroSyncOptions, wg *sync.WaitGroup) {
	defer wg.Done()
	if opts.IsDryRun {
		msg := fmt.Sprintf("(dry-run) would sync %d macro(s) to ~/.gitmap/macros/", len(macros))
		printNodeResultOutput(c.Alias, c.IPAddress, msg, nil)
		return
	}
	client, isConnected := connectSSHClient(c)
	if !isConnected {
		msg := "(authentication required - macro sync skipped)"
		printNodeResultOutput(c.Alias, c.IPAddress, msg, nil)
		return
	}
	defer client.Close()
	syncMacrosViaClient(client, c, macros, opts.IsForce)
}

func syncMacrosViaClient(client *ssh.Client, c db.SSHConnection, macros []macro.Macro, isForce bool) {
	syncedCount := syncAllMacrosToClient(client, c, macros, isForce)
	msg := fmt.Sprintf("Synced %d/%d macro(s) -> ~/.gitmap/macros/", syncedCount, len(macros))
	printNodeResultOutput(c.Alias, c.IPAddress, msg, nil)
}

func syncAllMacrosToClient(client *ssh.Client, c db.SSHConnection, macros []macro.Macro, isForce bool) int {
	syncedCount := 0
	shell := determineFallbackShell(c.OS)
	for _, m := range macros {
		if syncSingleMacroToClient(client, c, m, shell, isForce) {
			syncedCount++
		}
	}
	return syncedCount
}

func syncSingleMacroToClient(client *ssh.Client, c db.SSHConnection, m macro.Macro, shell string, isForce bool) bool {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return false
	}
	cmd := buildRemoteAtomicMacroWriteCmd(m.Name, data, isWindowsOS(c.OS), isForce)
	_, runErr := crypto.RunCommand(client, cmd, shell)
	return runErr == nil
}

func buildRemoteAtomicMacroWriteCmd(name string, data []byte, isWin bool, isForce bool) string {
	b64 := base64.StdEncoding.EncodeToString(data)
	if isWin {
		return buildWindowsMacroWriteCmd(name, b64, isForce)
	}
	return buildUnixMacroWriteCmd(name, b64, isForce)
}

func buildUnixMacroWriteCmd(name string, b64 string, isForce bool) string {
	forceCheck := ""
	if !isForce {
		forceCheck = fmt.Sprintf(`[ -f "$HOME/.gitmap/macros/%s.json" ] && exit 0; `, name)
	}
	return fmt.Sprintf(`mkdir -p "$HOME/.gitmap/macros" && %sprintf '%%s' '%s' | base64 -d > "$HOME/.gitmap/macros/%s.json.tmp" && mv -f "$HOME/.gitmap/macros/%s.json.tmp" "$HOME/.gitmap/macros/%s.json"`,
		forceCheck, b64, name, name, name)
}

func buildWindowsMacroWriteCmd(name string, b64 string, isForce bool) string {
	return fmt.Sprintf(`powershell -NoProfile -Command "$dir=[IO.Path]::Combine($env:USERPROFILE, '.gitmap', 'macros'); if (-not (Test-Path $dir)) { [IO.Directory]::CreateDirectory($dir) | Out-Null }; $p=[IO.Path]::Combine($dir, '%s.json'); if (%t -or -not (Test-Path $p)) { $d=[Convert]::FromBase64String('%s'); $tmp=$p + '.tmp'; [IO.File]::WriteAllBytes($tmp, $d); Move-Item -Force $tmp $p }"`,
		name, isForce, b64)
}

func printSyncStartBanner(macroCount, onlineCount, offlineCount int) {
	fmt.Println()
	fmt.Printf("  %s● Syncing %d macro(s) across %d machine(s)%s\n",
		constants.ColorCyan, macroCount, onlineCount, constants.ColorReset)
	if offlineCount > 0 {
		fmt.Printf("    %s(%d offline machine(s) skipped)%s\n",
			constants.ColorYellow, offlineCount, constants.ColorReset)
	}
	fmt.Println()
}

func printSSHMacroSyncHelp() {
	fmt.Println("Push local macro(s) to remote SSH machines.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh macro sync [name|all] [flags]")
	fmt.Println("  gitmap ssh macro-sync [name|all] [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: all online)")
	fmt.Println("      --except string     Exclude machines by alias, IP, or ID")
	fmt.Println("      --exclude string    Exclude machines (comma separated)")
	fmt.Println("  -f, --force             Overwrite existing remote macros")
	fmt.Println("      --dry-run           Simulate sync without writing files")
	fmt.Println("  -h, --help              Show help for macro sync")
}
