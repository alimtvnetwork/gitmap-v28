package cmdssh

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

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
	ExceptOS   string
	TargetOS   string
	Target     string
	IP         string
	IsJSON     bool
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
	fmt.Println("      --except-os string  Exclude machines by OS (e.g. unix, win, linux, ubuntu, darwin)")
	fmt.Println("      --os string         Target machines by OS (e.g. win, unix, linux, darwin)")
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
	fmt.Println("  gitmap ssh exec cmd1,cmd2,cmd3 --except-os unix")
	fmt.Println("  gitmap ssh exec cmd1,cmd2,cmd3 --except-os win")
	fmt.Println("  gitmap ssh exec cmd1,cmd2,cmd3 --except-os ubuntu")
}

func configureSEFlags(fs *flag.FlagSet, opts *seOptions) {
	fs.StringVar(&opts.Exclude, "exclude", "", "Exclude machines (comma separated)")
	fs.StringVar(&opts.Except, "except", "", "Exclude machines by alias, IP, or ID (comma separated)")
	fs.StringVar(&opts.ExceptOS, "except-os", "", "Exclude machines by OS (unix, win, ubuntu, etc.)")
	fs.StringVar(&opts.TargetOS, "os", "", "Target machines by OS")
	fs.StringVar(&opts.Target, "target", "", "Target machine alias or IP")
	fs.StringVar(&opts.Target, "t", "", "Target machine alias or IP (shorthand)")
	fs.StringVar(&opts.IP, "ip", "", "Target machine IP address")
	fs.BoolVar(&opts.IsJSON, "json", false, "Output results in JSON format")
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

func extractSECustomFlags(args []string) (seOptions, []string) {
	var opts seOptions
	var clean []string

	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--json" {
			opts.IsJSON = true
			continue
		}
		if a == "-h" || a == "--help" {
			opts.IsShowHelp = true
			continue
		}
		if (a == "--except-os" || a == "--exclude-os" || a == "--skip-os") && i+1 < len(args) {
			opts.ExceptOS = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--except-os=") {
			opts.ExceptOS = strings.TrimPrefix(a, "--except-os=")
			continue
		}
		if strings.HasPrefix(a, "--exclude-os=") {
			opts.ExceptOS = strings.TrimPrefix(a, "--exclude-os=")
			continue
		}
		if strings.HasPrefix(a, "--skip-os=") {
			opts.ExceptOS = strings.TrimPrefix(a, "--skip-os=")
			continue
		}

		if (a == "--os" || a == "--target-os") && i+1 < len(args) {
			opts.TargetOS = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--os=") {
			opts.TargetOS = strings.TrimPrefix(a, "--os=")
			continue
		}
		if strings.HasPrefix(a, "--target-os=") {
			opts.TargetOS = strings.TrimPrefix(a, "--target-os=")
			continue
		}

		if (a == "--except" || a == "--excep") && i+1 < len(args) {
			opts.Except = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--except=") || strings.HasPrefix(a, "--excep=") {
			opts.Except = strings.TrimPrefix(strings.TrimPrefix(a, "--except="), "--excep=")
			continue
		}

		if a == "--exclude" && i+1 < len(args) {
			opts.Exclude = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--exclude=") {
			opts.Exclude = strings.TrimPrefix(a, "--exclude=")
			continue
		}

		if (a == "-t" || a == "--target") && i+1 < len(args) {
			opts.Target = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--target=") {
			opts.Target = strings.TrimPrefix(a, "--target=")
			continue
		}

		if a == "--ip" && i+1 < len(args) {
			opts.IP = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--ip=") {
			opts.IP = strings.TrimPrefix(a, "--ip=")
			continue
		}

		clean = append(clean, a)
	}

	return opts, clean
}

func parseSEFlags(args []string) seOptions {
	if hasHelpFlag(args) {
		printSSHExecHelp()
		return seOptions{IsShowHelp: true}
	}
	opts, remaining := extractSECustomFlags(args)
	opts.Args = remaining
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
	if isInteractiveMacroAdd(opts.Args) {
		printInteractiveMacroAdvice()

		return apperror.NewValidationError("interactive macro creation cannot run over non-interactive SSH exec")
	}

	return executeSSHFromOptions(opts)
}

func executeSSHFromOptions(opts seOptions) error {
	conns, err := loadFilteredSSHConns(opts)
	if err != nil {
		return err
	}
	conns, execArgs := resolveExecTargetAndArgs(conns, opts)
	if len(conns) == 0 {
		return handleEmptySSHConns(opts.IsJSON)
	}
	return dispatchSSHExecIfAllowed(conns, execArgs, opts.IsJSON)
}

func handleEmptySSHConns(isJSON bool) error {
	if isJSON {
		fmt.Println("[]")
		return nil
	}
	fmt.Println("No machines to execute on.")
	return nil
}

func dispatchSSHExecIfAllowed(conns []db.SSHConnection, args []string, isJSON bool) error {
	if isInteractiveMacroAdd(args) {
		printInteractiveMacroAdvice()
		return apperror.NewValidationError("interactive macro creation cannot run over non-interactive SSH exec")
	}
	executeOnAllSSH(conns, args, isJSON)
	return nil
}

func loadFilteredSSHConns(opts seOptions) ([]db.SSHConnection, error) {
	data, err := fetchAllSSHConnections()
	if err != nil {
		return nil, nil
	}
	conns := filterSSHConns(data, resolveExcludeCSV(opts))
	if opts.Except != "" {
		conns = FilterSSHConnectionsByExcept(conns, opts.Except)
	}
	if opts.ExceptOS != "" || opts.TargetOS != "" {
		conns = FilterSSHConnectionsByOS(conns, opts.TargetOS, opts.ExceptOS)
	}
	return conns, nil
}

func fetchAllSSHConnections() ([]db.SSHConnection, error) {
	dbConn, err := store.OpenDefault()
	if err == nil {
		defer dbConn.Close()
		connsRes := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
		if !connsRes.IsFailure() && len(connsRes.Data) > 0 {
			return connsRes.Data, nil
		}
	}

	globalConn, gErr := store.OpenGlobalDefault()
	if gErr == nil {
		defer globalConn.Close()
		gRes := db.GetSSHConnections(globalConn.Context(), globalConn.SQL())
		if !gRes.IsFailure() && len(gRes.Data) > 0 {
			return gRes.Data, nil
		}
	}

	if err != nil {
		fmt.Printf("Failed to open DB: %v\n", err)
		return nil, err
	}
	return nil, nil
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

type NodeExecResult struct {
	Alias      string `json:"alias"`
	IP         string `json:"ip"`
	Status     string `json:"status"`
	ExitCode   int    `json:"exitCode"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMs int64  `json:"durationMs"`
	Error      string `json:"error,omitempty"`
}

func executeOnAllSSH(conns []db.SSHConnection, args []string, isJSON bool) {
	online, offline := partitionOnlineOffline(conns)
	if isJSON {
		executeOnAllSSHJSON(online, offline, args)
		return
	}
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

func executeOnAllSSHJSON(online, offline []db.SSHConnection, args []string) {
	var results []NodeExecResult
	for _, c := range offline {
		results = append(results, NodeExecResult{
			Alias:  c.Alias,
			IP:     c.IPAddress,
			Status: "offline",
		})
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, c := range online {
		wg.Add(1)
		go runSSHWorkerJSON(c, args, &results, &mu, &wg)
	}
	wg.Wait()
	_ = json.NewEncoder(os.Stdout).Encode(results)
}

func runSSHWorkerJSON(c db.SSHConnection, args []string, results *[]NodeExecResult, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	client, isConnected := connectSSHClient(c, "")
	if !isConnected {
		mu.Lock()
		*results = append(*results, NodeExecResult{
			Alias:  c.Alias,
			IP:     c.IPAddress,
			Status: "auth_failed",
			Error:  "authentication failed",
		})
		mu.Unlock()
		return
	}
	defer client.Close()
	c.OS = probeRemoteOSType(client)
	shellType, cmdStr, isDelegate := resolveWorkerCommand(c, args)
	if isDelegate && !ensureDelegateInstalledQuiet(client, c) {
		cmdStr = strings.Join(args, " ")
		shellType = determineFallbackShell(c.OS)
	}
	start := time.Now()
	out, err := crypto.RunCommand(client, cmdStr, shellType)
	dur := time.Since(start).Milliseconds()
	exitCode := resolveProcessExitCode(err)
	mu.Lock()
	*results = append(*results, NodeExecResult{
		Alias:      c.Alias,
		IP:         c.IPAddress,
		Status:     "online",
		ExitCode:   exitCode,
		Stdout:     out,
		DurationMs: dur,
	})
	mu.Unlock()
}

func ensureDelegateInstalledQuiet(client *ssh.Client, c db.SSHConnection) bool {
	err := ensureGitmapInstalled(client, c.OS, c.Alias)
	return err == nil
}

func runSSHWorker(c db.SSHConnection, args []string, wg *sync.WaitGroup) error {
	defer wg.Done()

	client, isConnected := connectSSHClient(c, fmt.Sprintf("[%s]", c.Alias))
	if isConnected == false {
		printNodeResultOutput(c.Alias, c.IPAddress, formatAuthFailedAdvice(c.Alias), nil)
		return nil
	}
	defer client.Close()

	if c.EncryptedPassword == "" {
		c.EncryptedPassword = queryHostPasswordFromDB(c.Alias, c.IPAddress)
	}

	return executeSSHPayload(client, c, args)
}

func formatAuthFailedAdvice(alias string) string {
	return fmt.Sprintf("(auth failed: run 'gitmap ssh fix-auth %s' or configure password)", alias)
}

func ensureDelegateInstalled(client *ssh.Client, c db.SSHConnection) error {
	err := ensureGitmapInstalled(client, c.OS, c.Alias)
	if err != nil {
		printNodeResultOutput(c.Alias, c.IPAddress, "", err)
		return err
	}
	return nil
}

func executeSSHPayload(client *ssh.Client, c db.SSHConnection, args []string) error {
	isHandled, _ := handleRemoteMacroAdd(client, c, args)
	if isHandled {
		return nil
	}

	c.OS = probeRemoteOSType(client)
	shellType, cmdStr, isDelegate := resolveWorkerCommand(c, args)
	if isDelegate && ensureDelegateInstalled(client, c) != nil {
		rawCmd := strings.Join(args, " ")
		fallbackShell := determineFallbackShell(c.OS)
		out, err := crypto.RunCommand(client, rawCmd, fallbackShell)
		printNodeResultOutput(c.Alias, c.IPAddress, out, err)
		return nil
	}
	if shellType == "ps" || shellType == "pwsh" {
		_ = ensurePowerShellInstalled(client, c.OS, c.Alias)
	}

	out, err := crypto.RunCommand(client, cmdStr, shellType)
	printNodeResultOutput(c.Alias, c.IPAddress, out, err)

	return nil
}

func resolveWorkerCommand(c db.SSHConnection, args []string) (string, string, bool) {
	if len(args) > 0 && isPowerCommand(extractFirstToken(args[0])) {
		shell, cmd := resolvePowerCommand(c, args)
		return shell, cmd, false
	}

	shellType, commandStr, delegateToGitmap := determineSSHCommand(c.OS, args)
	if delegateToGitmap {
		commandStr = resolveGitmapCommandString(args)
		shellType = ""
	}
	if isWindowsOS(c.OS) == false {
		commandStr = wrapUnixPath(commandStr)
	}

	return shellType, commandStr, delegateToGitmap
}

func connectSSHClient(c db.SSHConnection, headers ...string) (*ssh.Client, bool) {
	header := ""
	if len(headers) > 0 {
		header = headers[0]
	}
	if client, isPassOk := tryConnectWithPassword(c, header); isPassOk {
		return client, true
	}
	if client, isKeyOk := tryConnectWithKey(c, header); isKeyOk {
		return client, true
	}
	return connectClientDefaults(c, header)
}

func connectClientDefaults(c db.SSHConnection, header string) (*ssh.Client, bool) {
	if client, isDefaultOk := connectWithDefaultKey(c.IPAddress, c.Username, header); isDefaultOk {
		return client, true
	}
	return tryFallbackDBPassword(c, header)
}

func tryConnectWithPassword(c db.SSHConnection, header string) (*ssh.Client, bool) {
	if c.EncryptedPassword == "" {
		return nil, false
	}

	return connectWithEncryptedPassword(c, header)
}

func tryConnectWithKey(c db.SSHConnection, header string) (*ssh.Client, bool) {
	if c.KeyPath == "" {
		return nil, false
	}

	return connectWithKeyPath(c, header)
}

func tryFallbackDBPassword(c db.SSHConnection, header string) (*ssh.Client, bool) {
	if c.EncryptedPassword != "" {
		return nil, false
	}
	encPass := queryHostPasswordFromDB(c.Alias, c.IPAddress)
	if encPass == "" {
		return nil, false
	}
	c.EncryptedPassword = encPass
	return connectWithEncryptedPassword(c, header)
}

func queryHostPasswordFromDB(alias, ip string) string {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer dbConn.Close()

	var encPass string
	row := dbConn.SQL().QueryRow("SELECT COALESCE(encrypted_password, '') FROM ssh_hosts WHERE alias = ? OR ip = ? LIMIT 1", alias, ip)
	if scanErr := row.Scan(&encPass); scanErr == nil && encPass != "" {
		return encPass
	}
	return ""
}

func connectWithKeyPath(c db.SSHConnection, header string) (*ssh.Client, bool) {
	client, err := crypto.ConnectWithKey(c.IPAddress, c.Username, c.KeyPath)
	if err != nil {
		printHeaderError(header, "Connect error", err)
		return nil, false
	}

	return client, true
}

func decryptPasswordCandidate(enc string) (string, error) {
	plain, err := DecryptSSHPassword(enc)
	if err == nil && plain != "" {
		return plain, nil
	}
	passBytes, decErr := crypto.Decrypt(enc, getEncryptionKey())
	if decErr == nil {
		return string(passBytes), nil
	}
	if err != nil {
		return "", err
	}
	return "", decErr
}

func connectWithEncryptedPassword(c db.SSHConnection, header string) (*ssh.Client, bool) {
	plain, decErr := decryptPasswordCandidate(c.EncryptedPassword)
	if decErr != nil {
		printHeaderError(header, "Decrypt error", decErr)
		return nil, false
	}

	client, err := crypto.ConnectWithPassword(c.IPAddress, c.Username, plain)
	if err != nil {
		printHeaderError(header, "Connect error", err)
		return nil, false
	}

	return client, true
}

func printHeaderError(header, prefix string, err error) {
	if header != "" {
		fmt.Printf("%s %s: %v\n", header, prefix, err)
	}
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

type quoteTokenScanner struct {
	tokens    []string
	cur       strings.Builder
	inQuote   bool
	quoteChar byte
}

func (s *quoteTokenScanner) processChar(c byte) {
	if s.inQuote {
		s.processInQuote(c)

		return
	}
	s.processOutQuote(c)
}

func (s *quoteTokenScanner) processInQuote(c byte) {
	if c == s.quoteChar {
		s.inQuote = false

		return
	}
	s.cur.WriteByte(c)
}

func (s *quoteTokenScanner) processOutQuote(c byte) {
	if isQuoteChar(c) {
		s.inQuote = true
		s.quoteChar = c

		return
	}
	if isWhitespaceChar(c) {
		s.flushCurrentToken()

		return
	}
	s.cur.WriteByte(c)
}

func isQuoteChar(c byte) bool {
	return c == '"' || c == '\''
}

func isWhitespaceChar(c byte) bool {
	return c == ' ' || c == '\t'
}

func (s *quoteTokenScanner) flushCurrentToken() {
	if s.cur.Len() > 0 {
		s.tokens = append(s.tokens, s.cur.String())
		s.cur.Reset()
	}
}

func splitTokensWithQuotes(raw string) []string {
	scanner := &quoteTokenScanner{}
	for i := 0; i < len(raw); i++ {
		scanner.processChar(raw[i])
	}
	scanner.flushCurrentToken()

	return scanner.tokens
}

func isSplittableMacroArg(a string) bool {
	hasSpace := strings.Contains(a, " ")
	hasQuote := strings.Contains(a, "\"") || strings.Contains(a, "'")
	hasMacro := strings.Contains(a, "macro ") || strings.Contains(a, "gitmap ")

	return (hasSpace && hasQuote) || hasMacro
}

func tokenizeMacroArg(a string) []string {
	if isSplittableMacroArg(a) {
		return splitTokensWithQuotes(a)
	}

	return []string{a}
}

func extractMacroTokens(args []string) []string {
	var tokens []string
	for _, a := range args {
		tokens = append(tokens, tokenizeMacroArg(a)...)
	}

	return tokens
}

func findMacroAddTokenIndex(tokens []string) int {
	for i := 0; i < len(tokens)-1; i++ {
		if strings.EqualFold(tokens[i], "macro") && strings.EqualFold(tokens[i+1], "add") {
			return i + 1
		}
	}

	return -1
}

func isMacroHelpToken(token string) bool {
	return token == "help" || token == "--help" || token == "-h"
}

func isMacroFlagWithArg(a string) bool {
	return a == "--desc" || a == "--description" || a == "--tag"
}

func isFlagToken(a string) bool {
	return strings.HasPrefix(a, "-")
}

func hasMacroCommands(args []string) bool {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if isMacroFlagWithArg(a) {
			i++
			continue
		}
		if isFlagToken(a) {
			continue
		}
		return true
	}
	return false
}

func isInteractiveMacroAdd(args []string) bool {
	tokens := extractMacroTokens(args)
	addIdx := findMacroAddTokenIndex(tokens)
	if addIdx < 0 {
		return false
	}

	return isMacroAddPayloadInteractive(tokens, addIdx+1)
}

func isMacroAddPayloadInteractive(tokens []string, nameIdx int) bool {
	if nameIdx >= len(tokens) {
		return true
	}
	if isMacroHelpToken(tokens[nameIdx]) {
		return false
	}
	if hasMacroCommands(tokens[nameIdx+1:]) {
		return false
	}

	return true
}

func printInteractiveMacroAdvice() {
	fmt.Println("Interactive macro creation cannot run over non-interactive SSH exec. Create locally and sync:")
	fmt.Println("  1. gitmap macro add <name> <cmd1> [cmd2...] (non-interactive)")
	fmt.Println("  2. gitmap macro sync --all (sync to all SSH nodes)")
}
