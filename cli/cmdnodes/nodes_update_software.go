// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// SoftwareUpdateOptions specifies flags and target filters for fleet updates.
type SoftwareUpdateOptions struct {
	FilterOpts NodeFilterOptions `json:"filterOpts"`
	IsDryRun   bool              `json:"isDryRun"`
	IsJSON     bool              `json:"isJSON"`
	IsVerbose  bool              `json:"isVerbose"`
	IsDeploy   bool              `json:"isDeploy"`
}

// SoftwareUpdateNodeResult captures execution telemetry for a single node.
type SoftwareUpdateNodeResult struct {
	NodeAlias  string `json:"nodeAlias"`
	Host       string `json:"host"`
	OS         string `json:"os"`
	App        string `json:"app"`
	Status     string `json:"status"`
	IsSuccess  bool   `json:"isSuccess"`
	IsOffline  bool   `json:"isOffline"`
	Details    string `json:"details,omitempty"`
	DurationMs int64  `json:"durationMs"`
}

// SoftwareUpdateSummary aggregates results across all target nodes.
type SoftwareUpdateSummary struct {
	App         string                     `json:"app"`
	Total       int                        `json:"total"`
	Success     int                        `json:"success"`
	Failed      int                        `json:"failed"`
	Skipped     int                        `json:"skipped"`
	DurationMs  int64                      `json:"durationMs"`
	NodeResults []SoftwareUpdateNodeResult `json:"nodeResults"`
}

// RunFleetSoftwareUpdate coordinates parallel software updates across target fleet nodes.
func RunFleetSoftwareUpdate(appRaw string, args []string) error {
	app, isSupported := normalizeSoftwareApp(appRaw)
	if !isSupported {
		return apperror.NewValidation("app", "E5001", fmt.Sprintf("unsupported software %q (supported: cursor, vscode, agm, agy)", appRaw))
	}

	opts := parseSoftwareUpdateOptions(args)
	targets, err := resolveFleetSoftwareTargets(opts.FilterOpts)
	if err != nil || len(targets) == 0 {
		return err
	}

	return executeFleetSoftwareUpdatePlan(app, targets, opts)
}

func normalizeSoftwareApp(name string) (string, bool) {
	clean := strings.ToLower(strings.TrimSpace(name))
	switch clean {
	case "cursor", "cursor-ide":
		return "cursor", true
	case "vscode", "code", "vs-code":
		return "vscode", true
	case "agm", "ag-manager", "antigravity-manager":
		return "agm", true
	case "agy", "agy-tools", "antigravity", "antigravity-ide", "ide":
		return "agy", true
	default:
		return "", false
	}
}

func resolveFleetSoftwareTargets(filterOpts NodeFilterOptions) ([]db.SSHConnection, error) {
	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil {
		return nil, apperror.WrapSimple(err, "fetch ssh connections")
	}

	targets := FilterFleetNodes(conns, filterOpts)
	if len(targets) == 0 {
		fmt.Println("No matching fleet nodes found.")
		return nil, nil
	}

	return targets, nil
}

func executeFleetSoftwareUpdatePlan(app string, targets []db.SSHConnection, opts SoftwareUpdateOptions) error {
	start := time.Now()
	results := dispatchConcurrentSoftwareUpdate(app, targets, opts)
	summary := buildSoftwareUpdateSummary(app, results, time.Since(start))

	if opts.IsJSON {
		return printSoftwareUpdateJSON(summary)
	}

	renderSoftwareUpdateCard(summary)

	return nil
}

func dispatchConcurrentSoftwareUpdate(app string, targets []db.SSHConnection, opts SoftwareUpdateOptions) []SoftwareUpdateNodeResult {
	results := make([]SoftwareUpdateNodeResult, len(targets))
	var wg sync.WaitGroup

	for i, conn := range targets {
		wg.Add(1)
		go executeSoftwareWorker(&wg, results, i, conn, app, opts)
	}

	wg.Wait()

	return results
}

func executeSoftwareWorker(wg *sync.WaitGroup, res []SoftwareUpdateNodeResult, idx int, c db.SSHConnection, app string, opts SoftwareUpdateOptions) {
	defer wg.Done()

	res[idx] = executeSingleNodeSoftwareUpdate(c, app, opts)
}

func executeSingleNodeSoftwareUpdate(conn db.SSHConnection, app string, opts SoftwareUpdateOptions) SoftwareUpdateNodeResult {
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	isOnline := probeNodeTCPConnection(ctx, conn, 1200*time.Millisecond)
	cancel()

	if !isOnline {
		return buildOfflineSoftwareResult(conn, app)
	}

	if opts.IsDryRun {
		return buildDryRunSoftwareResult(conn, app)
	}

	return executeLiveSoftwareUpdate(conn, app, opts.IsDeploy)
}

func buildOfflineSoftwareResult(conn db.SSHConnection, app string) SoftwareUpdateNodeResult {
	return SoftwareUpdateNodeResult{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OS:        conn.OS,
		App:       app,
		Status:    "SKIPPED",
		IsOffline: true,
		Details:   "Node offline (1200ms probe timed out)",
	}
}

func buildDryRunSoftwareResult(conn db.SSHConnection, app string) SoftwareUpdateNodeResult {
	isWin := isWindowsNode(conn)
	cmd := resolveSoftwareUpdateCommand(app, isWin)

	return SoftwareUpdateNodeResult{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OS:        conn.OS,
		App:       app,
		Status:    "DRY-RUN",
		IsSuccess: true,
		Details:   "Would run: " + cmd,
	}
}

func executeLiveSoftwareUpdate(conn db.SSHConnection, app string, isDeploy bool) SoftwareUpdateNodeResult {
	isWin := isWindowsNode(conn)
	cmd := resolveSoftwareUpdateCommand(app, isWin)
	shell := resolveNodeShell(isWin)

	start := time.Now()
	out, errExec := currentSSHExecutor.Execute(conn, cmd, shell)
	duration := time.Since(start).Milliseconds()

	if errExec == nil && isDeploy && app == "agy" {
		_ = runSoftwareAgyDeploy(conn)
	}

	return populateSoftwareExecutionResult(conn, app, out, errExec, duration)
}

func runSoftwareAgyDeploy(conn db.SSHConnection) error {
	opts := cmdagy.AgyDeployOptions{
		TargetNode: conn.Alias,
		IsAll:      true,
	}
	_, err := cmdagy.ExecuteAgyDeploy(opts)

	return err
}

func populateSoftwareExecutionResult(conn db.SSHConnection, app, out string, errExec error, duration int64) SoftwareUpdateNodeResult {
	res := SoftwareUpdateNodeResult{
		NodeAlias:  conn.Alias,
		Host:       conn.IPAddress,
		OS:         conn.OS,
		App:        app,
		DurationMs: duration,
	}
	assignSoftwareOutcome(&res, out, errExec)

	return res
}

func assignSoftwareOutcome(res *SoftwareUpdateNodeResult, out string, errExec error) {
	if errExec != nil {
		res.Status = "FAILED"
		res.Details = errExec.Error()
		return
	}

	res.Status = "SUCCESS"
	res.IsSuccess = true
	res.Details = resolveSoftwareOutputSnippet(out)
}

func resolveSoftwareOutputSnippet(out string) string {
	trimmed := strings.TrimSpace(out)
	if len(trimmed) > 50 {
		return trimmed[:47] + "..."
	}

	return trimmed
}

func resolveSoftwareUpdateCommand(app string, isWin bool) string {
	if isWin {
		return resolveWindowsSoftwareCommand(app)
	}

	return resolveLinuxSoftwareCommand(app)
}

func resolveWindowsSoftwareCommand(app string) string {
	switch app {
	case "cursor":
		return "winget upgrade Anysphere.Cursor --silent --accept-source-agreements --accept-package-agreements"
	case "vscode":
		return "winget upgrade Microsoft.VisualStudioCode --silent --accept-source-agreements --accept-package-agreements"
	case "agm", "agy":
		return "irm https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1 | iex"
	default:
		return ""
	}
}

func resolveLinuxSoftwareCommand(app string) string {
	switch app {
	case "cursor":
		return "if [ -f repo-secrets/05-scripts/setup-cursor-ubuntu.py ]; then python3 repo-secrets/05-scripts/setup-cursor-ubuntu.py; elif [ -f ~/gitmap/repo-secrets/05-scripts/setup-cursor-ubuntu.py ]; then python3 ~/gitmap/repo-secrets/05-scripts/setup-cursor-ubuntu.py; else python3 repo-secrets/05-scripts/setup-cursor-ubuntu.py; fi"
	case "vscode":
		return "sudo apt-get update && sudo apt-get --only-upgrade install -y code"
	case "agm", "agy":
		return "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash"
	default:
		return ""
	}
}

func parseSoftwareUpdateOptions(args []string) SoftwareUpdateOptions {
	opts := SoftwareUpdateOptions{
		FilterOpts: ParseNodeFilterOptions(args),
	}
	for i := 0; i < len(args); i++ {
		step := parseSingleSoftwareFlag(&opts, args, i)
		if step > 0 {
			i += (step - 1)
			continue
		}
	}
	normalizeSoftwareFilterTarget(&opts.FilterOpts)

	return opts
}

func parseSingleSoftwareFlag(opts *SoftwareUpdateOptions, args []string, i int) int {
	arg := args[i]
	if slices.Contains([]string{"--json", "-j"}, arg) {
		opts.IsJSON = true
		return 1
	}
	if slices.Contains([]string{"--dry-run", "-n"}, arg) {
		opts.IsDryRun = true
		return 1
	}

	return parseSecondarySoftwareFlag(opts, arg)
}

func parseSecondarySoftwareFlag(opts *SoftwareUpdateOptions, arg string) int {
	if slices.Contains([]string{"--verbose", "-v"}, arg) {
		opts.IsVerbose = true
		return 1
	}
	if slices.Contains([]string{"--deploy", "--sync"}, arg) {
		opts.IsDeploy = true
		return 1
	}

	return 0
}

func normalizeSoftwareFilterTarget(opts *NodeFilterOptions) {
	isAll := strings.EqualFold(opts.Target, "all") || strings.EqualFold(opts.Target, "all-nodes") || opts.Target == "*"
	if isAll {
		opts.Target = ""
	}
}

func buildSoftwareUpdateSummary(app string, results []SoftwareUpdateNodeResult, duration time.Duration) SoftwareUpdateSummary {
	summary := SoftwareUpdateSummary{
		App:         app,
		Total:       len(results),
		DurationMs:  duration.Milliseconds(),
		NodeResults: results,
	}
	for _, r := range results {
		accumulateSoftwareSummaryCounts(&summary, r)
	}

	return summary
}

func accumulateSoftwareSummaryCounts(summary *SoftwareUpdateSummary, r SoftwareUpdateNodeResult) {
	if r.IsSuccess {
		summary.Success++
		return
	}
	if r.IsOffline {
		summary.Skipped++
		return
	}

	summary.Failed++
}

func renderSoftwareUpdateCard(summary SoftwareUpdateSummary) {
	renderSoftwareCardHeader(summary.App)
	renderSoftwareResultsTable(summary.NodeResults)
	renderSoftwareCardFooter(summary)
}

func renderSoftwareCardHeader(app string) {
	title := fmt.Sprintf("FLEET SOFTWARE UPDATE: %s", strings.ToUpper(app))
	fmt.Println()
	fmt.Println("  ╭──────────────────────────────────────────────────────────╮")
	fmt.Printf("  │ %-56s │\n", title)
	fmt.Println("  ╰──────────────────────────────────────────────────────────╯")
}

func renderSoftwareResultsTable(results []SoftwareUpdateNodeResult) {
	fmt.Println()
	fmt.Printf("  %s%-18s %-16s %-10s %-22s %-10s%s\n",
		constants.ColorCyan, "NODE ALIAS", "HOST", "OS", "STATUS", "LATENCY", constants.ColorReset)
	fmt.Println("  " + strings.Repeat("-", 80))
	for _, r := range results {
		statusBadge := formatSoftwareStatusBadge(r.Status)
		fmt.Printf("  %-18s %-16s %-10s %-32s %-10s\n",
			r.NodeAlias, r.Host, r.OS, statusBadge, fmt.Sprintf("%dms", r.DurationMs))
	}
	fmt.Println()
}

func formatSoftwareStatusBadge(status string) string {
	switch status {
	case "SUCCESS":
		return constants.ColorGreen + "● SUCCESS" + constants.ColorReset
	case "DRY-RUN":
		return constants.ColorYellow + "◌ DRY-RUN" + constants.ColorReset
	case "SKIPPED":
		return constants.ColorDim + "○ SKIPPED (OFFLINE)" + constants.ColorReset
	default:
		return constants.ColorRed + "▲ FAILED" + constants.ColorReset
	}
}

func renderSoftwareCardFooter(s SoftwareUpdateSummary) {
	durSec := float64(s.DurationMs) / 1000.0
	fmt.Println("  " + strings.Repeat("-", 80))
	fmt.Printf("  %sSummary: %d nodes | %s%d Succeeded%s | %s%d Failed%s | %s%d Skipped (Offline)%s (Duration: %.2fs)%s\n\n",
		constants.ColorBold, s.Total,
		constants.ColorGreen, s.Success, constants.ColorReset,
		constants.ColorRed, s.Failed, constants.ColorReset,
		constants.ColorDim, s.Skipped, constants.ColorReset,
		durSec, constants.ColorReset)
}

func printSoftwareUpdateJSON(summary SoftwareUpdateSummary) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(summary)
}
