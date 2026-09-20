package cmdssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

// RunSSHMacroExportCLI executes the ssh macro-export command.
func RunSSHMacroExportCLI(args []string) error {
	return runSSHMacroExportCLI(args)
}

// RunSSHMacroImportCLI executes the ssh macro-import command.
func RunSSHMacroImportCLI(args []string) error {
	return runSSHMacroImportCLI(args)
}

func runSSHMacroExportCLI(args []string) error {
	opts := parseSSHMacroExportOptions(args)
	if opts.IsHelp {
		printSSHMacroExportHelp()
		return nil
	}
	if opts.Name == "" {
		fmt.Println("Error: macro name required for export")
		return apperror.NewValidationError("macro name required for export")
	}
	return executeSSHMacroExport(opts)
}

func parseSSHMacroExportOptions(args []string) SSHMacroIOOptions {
	opts := SSHMacroIOOptions{}
	for i := 0; i < len(args); i++ {
		processExportArg(args[i], args, &i, &opts)
	}
	return opts
}

func processExportArg(arg string, args []string, index *int, opts *SSHMacroIOOptions) {
	if processExportFlag(arg, args, index, opts) {
		return
	}
	if !strings.HasPrefix(arg, "-") && opts.Name == "" {
		opts.Name = arg
	}
}

func processExportFlag(arg string, args []string, index *int, opts *SSHMacroIOOptions) bool {
	switch {
	case (arg == "-f" || arg == "--file") && *index+1 < len(args):
		*index++
		opts.FilePath = args[*index]
		return true
	case (arg == "-t" || arg == "--target") && *index+1 < len(args):
		*index++
		opts.Target = args[*index]
		return true
	case arg == "-h" || arg == "--help" || arg == "help":
		opts.IsHelp = true
		return true
	default:
		return false
	}
}

func executeSSHMacroExport(opts SSHMacroIOOptions) error {
	conn, err := resolveExportConnection(opts.Target)
	if err != nil {
		return err
	}
	data, err := readRemoteMacroData(*conn, opts.Name)
	if err != nil {
		return err
	}
	return saveExportedMacro(data, opts, conn.Alias)
}

func resolveExportConnection(target string) (*db.SSHConnection, error) {
	conns := loadTransferConnections(sshTransferOptions{Target: target})
	online, _ := partitionOnlineOffline(conns)
	if len(online) == 0 {
		fmt.Printf("  %sNo online SSH machine found for target: %s%s\n",
			constants.ColorYellow, target, constants.ColorReset)
		return nil, apperror.NewNotFoundError("no online SSH machine for target")
	}
	return &online[0], nil
}

func readRemoteMacroData(c db.SSHConnection, name string) ([]byte, error) {
	client, isConnected := connectSSHClient(c)
	if !isConnected {
		return nil, apperror.NewExecutionError("failed to connect to " + c.Alias)
	}
	defer client.Close()
	cmd := buildRemoteReadMacroCmd(name, isWindowsOS(c.OS))
	shell := determineFallbackShell(c.OS)
	out, err := crypto.RunCommand(client, cmd, shell)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read remote macro: "+name)
	}
	return decodeRemoteBase64Output(out)
}

func buildRemoteReadMacroCmd(name string, isWin bool) string {
	if isWin {
		return fmt.Sprintf(`powershell -NoProfile -Command "$p=[IO.Path]::Combine($env:USERPROFILE, '.gitmap', 'macros', '%s.json'); if (Test-Path $p) { [Convert]::ToBase64String([IO.File]::ReadAllBytes($p)) } else { exit 1 }"`, name)
	}
	return fmt.Sprintf(`cat "$HOME/.gitmap/macros/%s.json" | base64`, name)
}

func decodeRemoteBase64Output(out string) ([]byte, error) {
	trimmed := strings.TrimSpace(out)
	data, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, apperror.WrapSimple(err, "decode base64 macro")
	}
	return data, nil
}

func saveExportedMacro(data []byte, opts SSHMacroIOOptions, alias string) error {
	if opts.FilePath != "" {
		err := os.WriteFile(opts.FilePath, data, 0644)
		if err != nil {
			return apperror.WrapSimple(err, "write exported file")
		}
		printExportSuccessBanner(opts.Name, opts.FilePath, alias)
		return nil
	}
	var m macro.Macro
	if err := json.Unmarshal(data, &m); err != nil {
		return apperror.WrapSimple(err, "unmarshal exported macro")
	}
	if err := macro.SaveMacro(&m); err != nil {
		return err
	}
	printExportSuccessBanner(opts.Name, "local store", alias)
	return nil
}

func printExportSuccessBanner(name, dest, alias string) {
	fmt.Printf("\n  %s✔ Exported macro '%s' from %s to %s%s\n\n",
		constants.ColorGreen, name, alias, dest, constants.ColorReset)
}

func printSSHMacroExportHelp() {
	fmt.Println("Pull macro from remote SSH node into local store or file.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh macro export <name> [flags]")
	fmt.Println("  gitmap ssh macro-export <name> [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t, --target string     Target machine alias or IP")
	fmt.Println("  -f, --file string       Local destination file path")
	fmt.Println("  -h, --help              Show help for macro export")
}

func runSSHMacroImportCLI(args []string) error {
	opts := parseSSHMacroImportOptions(args)
	if opts.IsHelp {
		printSSHMacroImportHelp()
		return nil
	}
	if opts.FilePath == "" && opts.Name == "" {
		fmt.Println("Error: macro file or name required for import")
		return apperror.NewValidationError("macro file or name required for import")
	}
	return executeSSHMacroImport(opts)
}

func parseSSHMacroImportOptions(args []string) SSHMacroIOOptions {
	opts := SSHMacroIOOptions{}
	for i := 0; i < len(args); i++ {
		processImportArg(args[i], args, &i, &opts)
	}
	return opts
}

func processImportArg(arg string, args []string, index *int, opts *SSHMacroIOOptions) {
	if processImportFlag(arg, args, index, opts) {
		return
	}
	if !strings.HasPrefix(arg, "-") && opts.FilePath == "" {
		opts.FilePath = arg
	}
}

func processImportFlag(arg string, args []string, index *int, opts *SSHMacroIOOptions) bool {
	switch {
	case (arg == "-t" || arg == "--target") && *index+1 < len(args):
		*index++
		opts.Target = args[*index]
		return true
	case arg == "--as" && *index+1 < len(args):
		*index++
		opts.RenameAs = args[*index]
		return true
	case arg == "-f" || arg == "--force":
		opts.IsForce = true
		return true
	default:
		return processImportAltFlag(arg, args, index, opts)
	}
}

func processImportAltFlag(arg string, args []string, index *int, opts *SSHMacroIOOptions) bool {
	switch {
	case arg == "--dry-run":
		opts.IsDryRun = true
		return true
	case arg == "-h" || arg == "--help" || arg == "help":
		opts.IsHelp = true
		return true
	default:
		return false
	}
}

func executeSSHMacroImport(opts SSHMacroIOOptions) error {
	m, appErr := loadMacroForImport(opts)
	if appErr != nil {
		return appErr
	}
	syncOpts := SSHMacroSyncOptions{
		Target: opts.Target, Except: opts.Except, Exclude: opts.Exclude,
		IsForce: opts.IsForce, IsDryRun: opts.IsDryRun,
	}
	return executeSSHMacroSync([]macro.Macro{*m}, syncOpts)
}

func loadMacroForImport(opts SSHMacroIOOptions) (*macro.Macro, *apperror.AppError) {
	m, err := readMacroFromFileOrStore(opts.FilePath, opts.Name)
	if err != nil {
		return nil, apperror.WrapSimple(err, "load macro for import")
	}
	if opts.RenameAs != "" {
		m.Name = opts.RenameAs
	}
	return m, nil
}

func readMacroFromFileOrStore(filePath, name string) (*macro.Macro, error) {
	target := filePath
	if target == "" {
		target = name
	}
	data, err := os.ReadFile(target)
	if err == nil {
		var m macro.Macro
		jsonErr := json.Unmarshal(data, &m)
		if jsonErr == nil {
			return &m, nil
		}
	}
	return macro.LoadMacro(target)
}

func printSSHMacroImportHelp() {
	fmt.Println("Push macro file or saved macro to remote SSH nodes.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh macro import <file|name> [flags]")
	fmt.Println("  gitmap ssh macro-import <file|name> [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: all online)")
	fmt.Println("      --as string         Rename macro on remote nodes")
	fmt.Println("  -f, --force             Overwrite existing remote macros")
	fmt.Println("      --dry-run           Simulate import without writing files")
	fmt.Println("  -h, --help              Show help for macro import")
}
