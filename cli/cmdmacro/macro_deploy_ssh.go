package cmdmacro

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
	"golang.org/x/crypto/ssh"
)

// MacroDeployTarget represents an active cluster node for macro distribution.
type MacroDeployTarget struct {
	ID       string
	Alias    string
	IP       string
	Username string
	Port     int
	Password string
	KeyPath  string
	OS       string
}

// MacroDeployOptions contains parsed CLI flags for fleet deployment.
type MacroDeployOptions struct {
	Target   string
	Except   string
	IsDryRun bool
	IsForce  bool
}

// MacroDeployNodeResult stores the deployment status on a single node.
type MacroDeployNodeResult struct {
	Alias      string
	IP         string
	IsSuccess  bool
	Status     string
	MacroCount int
	DurationMs int64
	Details    string
	Error      error
}

// LoadClusterTargetsFn is a mockable provider for cluster target discovery.
var LoadClusterTargetsFn = loadClusterTargets

// TransferMacrosToTargetFn is a mockable transport worker for macro transfer.
var TransferMacrosToTargetFn = transferMacrosToTarget

// ExecuteMacroDeploySSH distributes macros across active cluster targets.
func ExecuteMacroDeploySSH(args []string) error {
	opts := parseMacroDeployFlags(args)
	macros, err := collectAllLocalMacros()
	if err != nil {
		return apperror.WrapSimple(err, "ExecuteMacroDeploySSH.collectAllLocalMacros")
	}

	targets, err := LoadClusterTargetsFn()
	if err != nil {
		return apperror.WrapSimple(err, "ExecuteMacroDeploySSH.LoadClusterTargets")
	}

	exclusionSet := parseExclusionSet(opts.Except)
	filteredTargets, excludedCount := filterDeployTargets(targets, opts.Target, exclusionSet)
	if len(filteredTargets) == 0 {
		printDeployNoTargetsBanner(excludedCount)
		return nil
	}

	results := executeFleetDeployment(filteredTargets, macros, opts)
	renderDeploySummary(results, excludedCount, len(macros))
	return nil
}

func parseMacroDeployFlags(args []string) MacroDeployOptions {
	cleanArgs := stripLeadingKeywords(args)
	opts := MacroDeployOptions{}
	for i := 0; i < len(cleanArgs); i++ {
		arg := cleanArgs[i]
		processDeployFlag(arg, cleanArgs, &i, &opts)
	}
	return opts
}

func stripLeadingKeywords(args []string) []string {
	tokens := args
	for len(tokens) > 0 {
		first := strings.ToLower(tokens[0])
		isKeyword := first == "deploy" || first == "ssh" || first == "peat" || first == "pea" || first == "macro"
		if !isKeyword {
			break
		}
		tokens = tokens[1:]
	}
	return tokens
}

func processDeployFlag(arg string, args []string, index *int, opts *MacroDeployOptions) {
	if isExceptParam(arg) {
		consumeExceptFlag(args, index, opts)
		return
	}
	if isTargetParam(arg) {
		consumeTargetFlag(args, index, opts)
		return
	}
	if arg == "--dry-run" {
		opts.IsDryRun = true
		return
	}
	if arg == "-f" || arg == "--force" {
		opts.IsForce = true
		return
	}
}

func isExceptParam(arg string) bool {
	return arg == "--except" || arg == "--exclude" || strings.HasPrefix(arg, "--except=") || strings.HasPrefix(arg, "--exclude=")
}

func isTargetParam(arg string) bool {
	return arg == "-t" || arg == "--target" || strings.HasPrefix(arg, "--target=")
}

func consumeExceptFlag(args []string, index *int, opts *MacroDeployOptions) {
	arg := args[*index]
	if strings.Contains(arg, "=") {
		opts.Except = strings.SplitN(arg, "=", 2)[1]
		return
	}
	var tokens []string
	for *index+1 < len(args) {
		next := args[*index+1]
		isFlag := strings.HasPrefix(next, "-")
		if isFlag {
			break
		}
		*index++
		tokens = append(tokens, next)
	}
	opts.Except = strings.Join(tokens, ",")
}

func consumeTargetFlag(args []string, index *int, opts *MacroDeployOptions) {
	arg := args[*index]
	if strings.Contains(arg, "=") {
		opts.Target = strings.SplitN(arg, "=", 2)[1]
		return
	}
	if *index+1 < len(args) {
		*index++
		opts.Target = args[*index]
	}
}

func parseExclusionSet(rawExcept string) map[string]bool {
	set := make(map[string]bool)
	if rawExcept == "" {
		return set
	}
	parts := strings.Split(rawExcept, ",")
	for _, p := range parts {
		cleaned := strings.ToLower(strings.TrimSpace(p))
		if cleaned != "" {
			set[cleaned] = true
		}
	}
	return set
}

func filterDeployTargets(targets []MacroDeployTarget, targetFilter string, exclusions map[string]bool) ([]MacroDeployTarget, int) {
	var filtered []MacroDeployTarget
	excludedCount := 0
	targetLower := strings.ToLower(targetFilter)
	for _, t := range targets {
		if isTargetMatch(t, targetLower) == false {
			continue
		}
		if isNodeExcluded(t, exclusions) {
			excludedCount++
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered, excludedCount
}

func isTargetMatch(t MacroDeployTarget, target string) bool {
	if target == "" || target == "all" || target == "all-nodes" {
		return true
	}
	if strings.EqualFold(t.Alias, target) {
		return true
	}
	if strings.EqualFold(t.IP, target) {
		return true
	}
	return strings.EqualFold(t.ID, target)
}

func isNodeExcluded(t MacroDeployTarget, exclusions map[string]bool) bool {
	if len(exclusions) == 0 {
		return false
	}
	if exclusions[strings.ToLower(t.ID)] {
		return true
	}
	if exclusions[strings.ToLower(t.Alias)] {
		return true
	}
	if exclusions[strings.ToLower(t.IP)] {
		return true
	}
	userHost := fmt.Sprintf("%s@%s", t.Username, t.IP)
	return exclusions[strings.ToLower(userHost)]
}

func collectAllLocalMacros() ([]macro.Macro, error) {
	res := macro.ListMacros()
	if res.IsFailure() {
		return nil, res.AppError()
	}
	return res.Data, nil
}

func loadClusterTargets() ([]MacroDeployTarget, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenDefault")
	}
	defer dbConn.Close()

	hosts, _ := store.ListSSHHosts(dbConn.Context(), dbConn.SQL())
	connsRes := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())

	var conns []db.SSHConnection
	if connsRes.IsSuccess() {
		conns = connsRes.Data
	}
	return mergeClusterTargets(hosts, conns), nil
}

func mergeClusterTargets(hosts []store.SSHHost, conns []db.SSHConnection) []MacroDeployTarget {
	seen := make(map[string]bool)
	var targets []MacroDeployTarget

	for _, h := range hosts {
		key := strings.ToLower(h.IP)
		seen[key] = true
		targets = append(targets, MacroDeployTarget{
			ID:       h.ID,
			Alias:    h.Alias,
			IP:       h.IP,
			Username: h.Username,
			Port:     resolvePortDefault(h.Port),
			Password: h.EncryptedPassword,
			OS:       "linux",
		})
	}

	for _, c := range conns {
		key := strings.ToLower(c.IPAddress)
		if seen[key] {
			continue
		}
		seen[key] = true
		targets = append(targets, MacroDeployTarget{
			ID:       fmt.Sprintf("host-%s", c.IPAddress),
			Alias:    c.Alias,
			IP:       c.IPAddress,
			Username: c.Username,
			Port:     22,
			Password: c.EncryptedPassword,
			KeyPath:  c.KeyPath,
			OS:       resolveOSDefault(c.OS),
		})
	}
	return targets
}

func resolvePortDefault(p int) int {
	if p > 0 {
		return p
	}
	return 22
}

func resolveOSDefault(osName string) string {
	if osName != "" {
		return osName
	}
	return "linux"
}

func executeFleetDeployment(targets []MacroDeployTarget, macros []macro.Macro, opts MacroDeployOptions) []MacroDeployNodeResult {
	results := make([]MacroDeployNodeResult, len(targets))
	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Printf("\n%s[MACRO FLEET DEPLOY]%s Dispatching %d macro(s) across %d target node(s)...\n\n",
		constants.ColorCyan, constants.ColorReset, len(macros), len(targets))

	for idx, t := range targets {
		wg.Add(1)
		go func(i int, target MacroDeployTarget) {
			defer wg.Done()
			res := runSingleTargetDeploy(target, macros, opts)
			mu.Lock()
			results[i] = res
			mu.Unlock()
		}(idx, t)
	}
	wg.Wait()
	return results
}

func runSingleTargetDeploy(target MacroDeployTarget, macros []macro.Macro, opts MacroDeployOptions) MacroDeployNodeResult {
	start := time.Now()
	err := TransferMacrosToTargetFn(target, macros, opts)
	dur := time.Since(start).Milliseconds()

	isSuccess := err == nil
	status := "SUCCESS"
	details := fmt.Sprintf("Deployed %d macro(s)", len(macros))
	if err != nil {
		status = "FAILED"
		details = err.Error()
	}
	if opts.IsDryRun {
		details = fmt.Sprintf("[DRY-RUN] %s", details)
	}

	res := MacroDeployNodeResult{
		Alias:      target.Alias,
		IP:         target.IP,
		IsSuccess:  isSuccess,
		Status:     status,
		MacroCount: len(macros),
		DurationMs: dur,
		Details:    details,
		Error:      err,
	}
	printSingleDeployProgress(res)
	return res
}

func printSingleDeployProgress(res MacroDeployNodeResult) {
	if res.IsSuccess {
		fmt.Printf("  %s✓%s [%s|%s] %s (%dms)\n",
			constants.ColorGreen, constants.ColorReset, res.Alias, res.IP, res.Details, res.DurationMs)
		return
	}
	fmt.Printf("  %s✖%s [%s|%s] %s: %v (%dms)\n",
		constants.ColorRed, constants.ColorReset, res.Alias, res.IP, res.Status, res.Error, res.DurationMs)
}

func transferMacrosToTarget(target MacroDeployTarget, macros []macro.Macro, opts MacroDeployOptions) error {
	if opts.IsDryRun {
		return nil
	}
	if tryRestMacroDeploy(target, macros) {
		return nil
	}
	return deployMacrosViaSSH(target, macros, opts.IsForce)
}

func tryRestMacroDeploy(target MacroDeployTarget, macros []macro.Macro) bool {
	bundle, err := json.Marshal(macros)
	if err != nil {
		return false
	}
	url := fmt.Sprintf("http://%s:49152/api/v1/macro/import", target.IP)
	client := http.Client{Timeout: 1500 * time.Millisecond}
	resp, reqErr := client.Post(url, "application/json", strings.NewReader(string(bundle)))
	if reqErr != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated
}

func deployMacrosViaSSH(target MacroDeployTarget, macros []macro.Macro, isForce bool) error {
	client, err := dialTargetSSH(target)
	if err != nil {
		return err
	}
	defer client.Close()

	isWin := strings.EqualFold(target.OS, "windows")
	shellType := "sh"
	if isWin {
		shellType = ""
	}

	for _, m := range macros {
		if deployErr := writeMacroToRemoteNode(client, m, isWin, shellType, isForce); deployErr != nil {
			return deployErr
		}
	}
	return nil
}

func dialTargetSSH(target MacroDeployTarget) (*ssh.Client, error) {
	if c, ok := tryDialPassword(target); ok {
		return c, nil
	}
	if c, ok := tryDialKeyPath(target); ok {
		return c, nil
	}
	if c, ok := tryDialCandidateKeys(target); ok {
		return c, nil
	}
	return nil, fmt.Errorf("ssh connection failed for %s@%s", target.Username, target.IP)
}

func tryDialPassword(target MacroDeployTarget) (*ssh.Client, bool) {
	if target.Password == "" {
		return nil, false
	}
	c, err := crypto.ConnectWithPassword(target.IP, target.Username, target.Password)
	return c, err == nil
}

func tryDialKeyPath(target MacroDeployTarget) (*ssh.Client, bool) {
	if target.KeyPath == "" {
		return nil, false
	}
	c, err := crypto.ConnectWithKey(target.IP, target.Username, target.KeyPath)
	return c, err == nil
}

func tryDialCandidateKeys(target MacroDeployTarget) (*ssh.Client, bool) {
	home, _ := os.UserHomeDir()
	candidateKeys := []string{
		fmt.Sprintf("%s/.ssh/id_ed25519", home),
		fmt.Sprintf("%s/.ssh/id_rsa", home),
	}
	for _, k := range candidateKeys {
		c, err := crypto.ConnectWithKey(target.IP, target.Username, k)
		if err == nil {
			return c, true
		}
	}
	return nil, false
}

func writeMacroToRemoteNode(client *ssh.Client, m macro.Macro, isWin bool, shellType string, isForce bool) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	cmd := buildAtomicMacroWriteScript(m.Name, b64, isWin, isForce)
	_, runErr := crypto.RunCommand(client, cmd, shellType)
	return runErr
}

func resolveWindowsForceFlag(isForce bool) string {
	if isForce {
		return "$true"
	}
	return "$false"
}

func resolveUnixForceCheck(name string, isForce bool) string {
	if isForce {
		return ""
	}
	return fmt.Sprintf(`[ -f "$MDIR/%s.json" ] && exit 0; `, name)
}

func buildAtomicMacroWriteScript(name, b64 string, isWin, isForce bool) string {
	if isWin {
		forceFlag := resolveWindowsForceFlag(isForce)
		return fmt.Sprintf(`powershell -NoProfile -Command "$dir=[IO.Path]::Combine($env:USERPROFILE, '.gitmap', 'macros'); if (-not (Test-Path $dir)) { [IO.Directory]::CreateDirectory($dir) | Out-Null }; $p=[IO.Path]::Combine($dir, '%s.json'); if (%s -or -not (Test-Path $p)) { $d=[Convert]::FromBase64String('%s'); $tmp=$p + '.tmp'; [IO.File]::WriteAllBytes($tmp, $d); Move-Item -Force $tmp $p }"`,
			name, forceFlag, b64)
	}
	forceCheck := resolveUnixForceCheck(name, isForce)
	return fmt.Sprintf(`MDIR=""; for d in "$HOME/.gitmap/macros" "$HOME/.config/gitmap/macros" "$HOME/.local/share/gitmap/macros" "/tmp/.gitmap-$USER/macros"; do if mkdir -p "$d" 2>/dev/null && [ -w "$d" ]; then MDIR="$d"; break; fi; done; if [ -z "$MDIR" ]; then exit 1; fi; %sprintf '%%s' '%s' | base64 -d > "$MDIR/%s.json.tmp" && mv -f "$MDIR/%s.json.tmp" "$MDIR/%s.json"`,
		forceCheck, b64, name, name, name)
}

func printDeployNoTargetsBanner(excludedCount int) {
	fmt.Printf("\n%s[MACRO FLEET DEPLOY]%s No matching active targets found (excluded: %d).\n\n",
		constants.ColorYellow, constants.ColorReset, excludedCount)
}

func renderDeploySummary(results []MacroDeployNodeResult, excludedCount, macroCount int) {
	if len(results) == 0 {
		return
	}
	successCount, failCount := calculateDeployMetrics(results)
	fmt.Printf("\n%s================================================================================%s\n",
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" %sMacro Fleet Deployment Summary:%s Total: %d | Succeeded: %d | Failed: %d | Excluded: %d\n",
		constants.ColorBold, constants.ColorReset, len(results), successCount, failCount, excludedCount)
	fmt.Printf("%s--------------------------------------------------------------------------------%s\n",
		constants.ColorDim, constants.ColorReset)

	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "ALIAS", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "IP", Align: termtable.AlignLeft, MinWidth: 16},
			{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 10},
			{Title: "MACROS", Align: termtable.AlignRight, MinWidth: 10},
			{Title: "DURATION", Align: termtable.AlignRight, MinWidth: 10},
			{Title: "DETAILS", Align: termtable.AlignLeft, MinWidth: 25},
		},
		Rows: buildDeployTableRows(results),
	}
	termtable.PrintTable(cfg)
	fmt.Printf("%s================================================================================%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func calculateDeployMetrics(results []MacroDeployNodeResult) (int, int) {
	succeeded := 0
	failed := 0
	for _, r := range results {
		if r.IsSuccess {
			succeeded++
			continue
		}
		failed++
	}
	return succeeded, failed
}

func buildDeployTableRows(results []MacroDeployNodeResult) []termtable.Row {
	rows := make([]termtable.Row, 0, len(results))
	for _, r := range results {
		statusStr := constants.ColorRed + "FAILED" + constants.ColorReset
		if r.IsSuccess {
			statusStr = constants.ColorGreen + "SUCCESS" + constants.ColorReset
		}
		rows = append(rows, termtable.Row{
			Cells: []string{
				r.Alias,
				r.IP,
				statusStr,
				fmt.Sprintf("%d", r.MacroCount),
				fmt.Sprintf("%dms", r.DurationMs),
				r.Details,
			},
		})
	}
	return rows
}
