package cmdssh

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// RunSSHAuthKeyDeployCLI executes public key deployment to SSH fleet or node.
func RunSSHAuthKeyDeployCLI(args []string) error {
	if hasHelpFlag(args) {
		printAuthKeyHelp()
		return nil
	}
	return executeAuthKeyDeploy(args)
}

func printAuthKeyHelp() {
	fmt.Printf("\n%sDeploy SSH Public Key to Fleet or Node%s\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Usage:")
	fmt.Println("  gitmap ssh fix-auth [target...] [-i <pubkey>] [--unix]")
	fmt.Println("  gitmap ssh auth-key deploy [target] [-i <pubkey>]")
	fmt.Println("  gitmap ssh copy-id [target] [-i <pubkey>]")
	fmt.Println("  gitmap sj fix-auth [target] [-i <pubkey>]")
	fmt.Println("  gitmap cluster fix-auth [target] [-i <pubkey>]")
	fmt.Println("  gitmap sc fix-auth [target] [-i <pubkey>]")
	fmt.Println("\nExamples:")
	fmt.Println("  gitmap ssh fix-auth machineid, ip, id")
	fmt.Println("  gitmap ssh fix-auth devbox")
	fmt.Println("  gitmap ssh copy-id user@192.168.1.14")
	fmt.Println("  gitmap ssh fix-auth --all")
	fmt.Println("  gitmap ssh fix-auth node1 -i ~/.ssh/id_ed25519.pub")
}

func executeAuthKeyDeploy(args []string) error {
	target, identityPath, isForceUnix, appErr := parseAuthKeyArgs(args)
	if appErr != nil {
		return appErr
	}

	pubKey, keyPath, appErr := discoverPublicKey(identityPath)
	if appErr != nil {
		return appErr
	}

	return runFleetKeyDeployment(target, pubKey, keyPath, isForceUnix)
}

func parseAuthKeyArgs(args []string) (string, string, bool, *apperror.AppError) {
	filtered := stripSubcommandKeyword(args, "deploy")
	target, identityPath, isForceUnix := "all", "", false
	for i := 0; i < len(filtered); i++ {
		t, id, unix, nextIdx, err := stepAuthKeyArg(filtered, i, target, identityPath, isForceUnix)
		if err != nil {
			return "", "", false, err
		}
		target, identityPath, isForceUnix, i = t, id, unix, nextIdx
	}
	return target, identityPath, isForceUnix, nil
}

func isUnixFlag(arg string) bool {
	return arg == "--unix" || arg == "-u" || arg == "--linux"
}

func stepAuthKeyArg(args []string, i int, curTarget, curID string, curUnix bool) (string, string, bool, int, *apperror.AppError) {
	arg := args[i]
	if isUnixFlag(arg) {
		return curTarget, curID, true, i, nil
	}
	if !isIdentityFlag(arg) {
		return resolveTargetFromArg(arg, curTarget), curID, curUnix, i, nil
	}
	val, nextIdx, err := parseIdentityFlagVal(args, i)
	if err != nil {
		return "", "", curUnix, i, err
	}
	return curTarget, val, curUnix, nextIdx, nil
}

func isSubcommandKeyword(s string) bool {
	return strings.EqualFold(s, "deploy") ||
		strings.EqualFold(s, "fix-auth") ||
		strings.EqualFold(s, "auth-key") ||
		strings.EqualFold(s, "copy-id")
}

func stripSubcommandKeyword(args []string, keyword string) []string {
	if len(args) > 0 && (strings.EqualFold(args[0], keyword) || isSubcommandKeyword(args[0])) {
		return args[1:]
	}
	return args
}

func isIdentityFlag(arg string) bool {
	return arg == "-i" || arg == "--identity" || arg == "--identity-file" ||
		strings.HasPrefix(arg, "-i=") || strings.HasPrefix(arg, "--identity=")
}

func parseIdentityFlagVal(args []string, idx int) (string, int, *apperror.AppError) {
	arg := args[idx]
	if eqIdx := strings.Index(arg, "="); eqIdx != -1 {
		return arg[eqIdx+1:], idx, nil
	}
	hasValue := idx+1 < len(args)
	if hasValue {
		return args[idx+1], idx + 1, nil
	}
	return "", idx, apperror.NewValidationError("missing identity path for flag " + arg)
}

func resolveTargetFromArg(arg, currentTarget string) string {
	if isAllFlag(arg) {
		return "all"
	}
	isPositional := !strings.HasPrefix(arg, "-")
	if isPositional {
		if currentTarget == "all" {
			return arg
		}
		return currentTarget + " " + arg
	}
	return currentTarget
}

func runFleetKeyDeployment(target, pubKey, keyPath string, isForceUnix bool) error {
	conns, err := resolveAuthKeyTargets(target)
	if err != nil {
		return err
	}
	if len(conns) == 0 {
		return apperror.NewNotFoundError("no target SSH machines found for target: " + target)
	}
	executeFleetKeyDeployment(conns, target, pubKey, keyPath, isForceUnix)
	return nil
}

func executeFleetKeyDeployment(conns []db.SSHConnection, target, pubKey, keyPath string, isForceUnix bool) {
	fmt.Printf("\n%s Deploying authorized key (%s) to SSH fleet (%s):%s\n\n",
		constants.ColorCyan, filepath.Base(keyPath), target, constants.ColorReset)
	for _, c := range conns {
		deployAuthKeyToNode(c, pubKey, isForceUnix)
	}
	fmt.Printf("\nAuthorized key deployment complete.\n\n")
}

func resolveTargetOS(osType string, isForceUnix bool) string {
	if isForceUnix || osType == "" {
		return "linux"
	}
	return osType
}

func deployAuthKeyToNode(c db.SSHConnection, pubKey string, isForceUnix bool) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := connectAuthKeyNode(c, header)
	if !isConnected {
		return
	}
	defer client.Close()

	osType := resolveTargetOS(c.OS, isForceUnix)
	executeKeyInjection(client, header, osType, pubKey)
}

func connectAuthKeyNode(c db.SSHConnection, header string) (*ssh.Client, bool) {
	client, isConnected := connectSSHNode(c, header)
	if isConnected {
		return client, true
	}
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return promptAndConnectSSH(c, header)
	}
	return nil, false
}

func promptAndConnectSSH(c db.SSHConnection, header string) (*ssh.Client, bool) {
	prompt := fmt.Sprintf("%s Enter SSH password for %s@%s: ", header, c.Username, c.IPAddress)
	pass, err := PromptSSHPassword(context.Background(), prompt, int(os.Stdin.Fd()))
	if err != nil || pass == "" {
		return nil, false
	}
	client, connErr := crypto.ConnectWithPassword(c.IPAddress, c.Username, pass)
	if connErr != nil {
		printHeaderError(header, "Password connect failed", connErr)
		return nil, false
	}
	return client, true
}

func executeKeyInjection(client *ssh.Client, header, osType, pubKey string) {
	script := buildInjectAuthKeyScript(pubKey, osType)
	shellType := resolveRemoteShell(osType)
	out, err := crypto.RunCommand(client, script, shellType)
	reportRemoteExecution(header, "Auth key deployment", out, err)
}

func resolveAuthKeyTargets(target string) ([]db.SSHConnection, error) {
	conns, err := loadSSHConnectionsForTarget(target)
	if err != nil {
		return nil, apperror.WrapSimple(err, "loadSSHConnectionsForTarget")
	}
	hasMatches := len(conns) > 0 || target == "all" || target == ""
	if hasMatches {
		return conns, nil
	}
	return resolveFallbackTarget(target)
}

func resolveFallbackTarget(target string) ([]db.SSHConnection, error) {
	t, parseErr := ParseSSHTarget(target, "", 22)
	isParsed := parseErr == nil && t != nil && t.IP != ""
	if isParsed {
		conn := db.SSHConnection{
			Alias:     t.IP,
			IPAddress: t.IP,
			Username:  t.Username,
		}
		return []db.SSHConnection{conn}, nil
	}
	return nil, apperror.NewNotFoundError("no SSH machines found for target: " + target)
}

func discoverPublicKey(identityPath string) (string, string, *apperror.AppError) {
	if identityPath != "" {
		return loadSpecifiedPublicKey(identityPath)
	}
	return loadDefaultPublicKey()
}

func expandHomePath(p string) string {
	if !strings.HasPrefix(p, "~/") && p != "~" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	return filepath.Join(home, strings.TrimPrefix(p, "~/"))
}

func loadSpecifiedPublicKey(rawPath string) (string, string, *apperror.AppError) {
	path := expandHomePath(rawPath)
	resolvedPath := resolvePubKeyFile(path)
	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return "", "", apperror.Wrap(err, "readPublicKey", map[string]any{"path": resolvedPath})
	}
	content := strings.TrimSpace(string(data))
	isEmpty := content == ""
	if isEmpty {
		return "", "", apperror.NewValidationError("public key file is empty: " + resolvedPath)
	}
	return content, resolvedPath, nil
}

func resolvePubKeyFile(path string) string {
	fi, err := os.Stat(path)
	isFile := err == nil && !fi.IsDir()
	if isFile {
		return path
	}
	pubVariant := path + ".pub"
	if fiPub, errPub := os.Stat(pubVariant); errPub == nil && !fiPub.IsDir() {
		return pubVariant
	}
	return path
}

func loadDefaultPublicKey() (string, string, *apperror.AppError) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", apperror.WrapSimple(err, "os.UserHomeDir")
	}

	candidates := []string{
		filepath.Join(home, ".ssh", "id_ed25519.pub"),
		filepath.Join(home, ".ssh", "id_rsa.pub"),
		filepath.Join(home, ".ssh", "id_ecdsa.pub"),
	}

	return scanKeyCandidates(candidates)
}

func scanKeyCandidates(candidates []string) (string, string, *apperror.AppError) {
	for _, c := range candidates {
		key, isFound := readKeyCandidate(c)
		if isFound {
			return key, c, nil
		}
	}
	return "", "", apperror.NewNotFoundError("no SSH public key found (~/.ssh/id_ed25519.pub, ~/.ssh/id_rsa.pub, ~/.ssh/id_ecdsa.pub); specify with -i <path>")
}

func readKeyCandidate(path string) (string, bool) {
	fi, statErr := os.Stat(path)
	if statErr != nil || fi.IsDir() {
		return "", false
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return "", false
	}
	content := strings.TrimSpace(string(data))
	return content, len(content) > 0
}

func buildInjectAuthKeyScript(pubKey, osType string) string {
	escapedKey := escapeSingleQuotes(pubKey)
	if strings.EqualFold(osType, "windows") {
		return buildWindowsAuthKeyScript(escapedKey)
	}

	return buildUnixAuthKeyScript(escapedKey)
}

func buildUnixAuthKeyScript(escapedKey string) string {
	return fmt.Sprintf("mkdir -p ~/.ssh && chmod 700 ~/.ssh && touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && (grep -qF '%s' ~/.ssh/authorized_keys || echo '%s' >> ~/.ssh/authorized_keys)", escapedKey, escapedKey)
}

const winAuthKeyDeployScript = `powershell -NoProfile -Command "$k = '%s'.Trim(); $tokens = $k -split '\s+'; if ($tokens.Length -lt 2) { exit 1 }; $keyBody = $tokens[1]; $isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator); if (-not $isAdmin) { $members = net localgroup administrators 2>$null; if ($members -match $env:USERNAME) { $isAdmin = $true } }; if ($isAdmin) { $keysFile = Join-Path $env:ProgramData 'ssh\administrators_authorized_keys'; $sshDir = Join-Path $env:ProgramData 'ssh'; if (!(Test-Path $sshDir)) { New-Item -ItemType Directory -Path $sshDir -Force | Out-Null }; if (!(Test-Path $keysFile)) { New-Item -ItemType File -Path $keysFile -Force | Out-Null }; icacls $keysFile /inheritance:r /grant 'Administrators:F' /grant 'SYSTEM:F' | Out-Null } else { $sshDir = Join-Path $env:USERPROFILE '.ssh'; if (!(Test-Path $sshDir)) { New-Item -ItemType Directory -Path $sshDir -Force | Out-Null }; $keysFile = Join-Path $sshDir 'authorized_keys'; if (!(Test-Path $keysFile)) { New-Item -ItemType File -Path $keysFile -Force | Out-Null } }; $lines = Get-Content $keysFile -ErrorAction SilentlyContinue; $hasKey = $false; foreach ($line in $lines) { $t = $line.Trim() -split '\s+'; if ($t.Length -ge 2 -and $t[1] -eq $keyBody) { $hasKey = $true; break } }; if (-not $hasKey) { Add-Content -Path $keysFile -Value $k }; $svc = Get-Service sshd -ErrorAction SilentlyContinue; if ($svc) { if ($svc.StartType -ne 'Automatic') { Set-Service sshd -StartupType Automatic -ErrorAction SilentlyContinue }; if ($svc.Status -ne 'Running') { Start-Service sshd -ErrorAction SilentlyContinue } }"`

func buildWindowsAuthKeyScript(escapedKey string) string {
	return fmt.Sprintf(winAuthKeyDeployScript, escapedKey)
}
