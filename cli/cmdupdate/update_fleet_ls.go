package cmdupdate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/termtable"
)

// FleetInventoryItem represents an installed component on a remote node.
type FleetInventoryItem struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	InstallPath string `json:"path,omitempty"`
}

// FleetNodeInventory represents the full software inventory returned from a node.
type FleetNodeInventory struct {
	NodeID string               `json:"node_id,omitempty"`
	Alias  string               `json:"alias,omitempty"`
	IP     string               `json:"ip,omitempty"`
	OS     string               `json:"os,omitempty"`
	Items  []FleetInventoryItem `json:"items,omitempty"`
	Error  string               `json:"error,omitempty"`
}

// FleetLSNodeResult stores query outcome for a single node.
type FleetLSNodeResult struct {
	Alias      string
	IP         string
	IsSuccess  bool
	Inventory  FleetNodeInventory
	DurationMs int64
	Error      error
}

// ExecuteRemoteInventoryFn is a mockable remote inventory query executor.
var ExecuteRemoteInventoryFn = executeDefaultRemoteInventory

// ExecuteFleetUpdateLS queries cluster nodes in parallel for installed software inventory.
func ExecuteFleetUpdateLS(args []string) error {
	opts := parseFleetUpdateOptions(args)
	targets, err := LoadFleetTargetsFn()
	if err != nil {
		return apperror.WrapSimple(err, "ExecuteFleetUpdateLS.LoadFleetTargets")
	}

	exclusionSet := parseFleetExclusionSet(opts.Except)
	filteredTargets, excludedCount := filterFleetTargets(targets, opts.Target, exclusionSet)
	if len(filteredTargets) == 0 {
		printFleetNoTargetsBanner("inventory ls", excludedCount)
		return nil
	}

	results := executeParallelFleetLS(filteredTargets)
	renderFleetLSSummary(results, excludedCount)
	return nil
}

func executeParallelFleetLS(targets []FleetTarget) []FleetLSNodeResult {
	results := make([]FleetLSNodeResult, len(targets))
	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Printf("\n%s[FLEET UPDATE LS]%s Querying software inventory across %d cluster node(s) in parallel...\n\n",
		constants.ColorCyan, constants.ColorReset, len(targets))

	for idx, t := range targets {
		wg.Add(1)
		go func(i int, target FleetTarget) {
			defer wg.Done()
			res := runSingleFleetLS(target)
			mu.Lock()
			results[i] = res
			mu.Unlock()
		}(idx, t)
	}
	wg.Wait()
	return results
}

func runSingleFleetLS(target FleetTarget) FleetLSNodeResult {
	start := time.Now()
	rawOutput, err := ExecuteRemoteInventoryFn(target)
	dur := time.Since(start).Milliseconds()

	inventory := ParseFleetInventoryTelemetry(rawOutput, target, err)
	hasItems := len(inventory.Items) > 0
	isSuccess := err == nil && (hasItems || inventory.Error == "")

	res := FleetLSNodeResult{
		Alias:      target.Alias,
		IP:         target.IP,
		IsSuccess:  isSuccess,
		Inventory:  inventory,
		DurationMs: dur,
		Error:      err,
	}
	printSingleFleetLSProgress(res)
	return res
}

func printSingleFleetLSProgress(res FleetLSNodeResult) {
	if res.IsSuccess {
		fmt.Printf("  %s✓%s [%s|%s] Found %d installed package(s) (%dms)\n",
			constants.ColorGreen, constants.ColorReset, res.Alias, res.IP, len(res.Inventory.Items), res.DurationMs)
		return
	}
	errDesc := "query failed"
	if res.Inventory.Error != "" {
		errDesc = res.Inventory.Error
	}
	fmt.Printf("  %s✖%s [%s|%s] FAILED: %s (%dms)\n",
		constants.ColorRed, constants.ColorReset, res.Alias, res.IP, errDesc, res.DurationMs)
}

// ParseFleetInventoryTelemetry parses JSON summary inventory returned from remote nodes.
func ParseFleetInventoryTelemetry(raw string, target FleetTarget, execErr error) FleetNodeInventory {
	if execErr != nil {
		return FleetNodeInventory{
			NodeID: target.ID,
			Alias:  target.Alias,
			IP:     target.IP,
			Error:  execErr.Error(),
		}
	}

	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return FleetNodeInventory{
			NodeID: target.ID,
			Alias:  target.Alias,
			IP:     target.IP,
			Error:  "Empty inventory response",
		}
	}

	if inv, ok := tryParseNodeInventoryObject(trimmed, target); ok {
		return inv
	}
	if inv, ok := tryParseInventoryItemsArray(trimmed, target); ok {
		return inv
	}
	if inv, ok := tryParseInventoryMap(trimmed, target); ok {
		return inv
	}
	return parsePlainTextInventory(trimmed, target)
}

func tryParseNodeInventoryObject(raw string, target FleetTarget) (FleetNodeInventory, bool) {
	var parsed FleetNodeInventory
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return FleetNodeInventory{}, false
	}
	if len(parsed.Items) == 0 && parsed.Error == "" {
		return FleetNodeInventory{}, false
	}
	populateNodeInventoryDefaults(&parsed, target)
	return parsed, true
}

func populateNodeInventoryDefaults(inv *FleetNodeInventory, target FleetTarget) {
	if inv.NodeID == "" {
		inv.NodeID = target.ID
	}
	if inv.Alias == "" {
		inv.Alias = target.Alias
	}
	if inv.IP == "" {
		inv.IP = target.IP
	}
	if inv.OS == "" {
		inv.OS = target.OS
	}
}

func tryParseInventoryItemsArray(raw string, target FleetTarget) (FleetNodeInventory, bool) {
	var items []FleetInventoryItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return FleetNodeInventory{}, false
	}
	return FleetNodeInventory{
		NodeID: target.ID,
		Alias:  target.Alias,
		IP:     target.IP,
		OS:     target.OS,
		Items:  items,
	}, true
}

func tryParseInventoryMap(raw string, target FleetTarget) (FleetNodeInventory, bool) {
	var kv map[string]string
	if err := json.Unmarshal([]byte(raw), &kv); err != nil {
		return FleetNodeInventory{}, false
	}
	var items []FleetInventoryItem
	for k, v := range kv {
		items = append(items, FleetInventoryItem{
			Name:    k,
			Version: v,
			Status:  "installed",
		})
	}
	return FleetNodeInventory{
		NodeID: target.ID,
		Alias:  target.Alias,
		IP:     target.IP,
		OS:     target.OS,
		Items:  items,
	}, true
}

func parsePlainTextInventory(raw string, target FleetTarget) FleetNodeInventory {
	lines := strings.Split(raw, "\n")
	var items []FleetInventoryItem
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if isIgnoredInventoryLine(l) {
			continue
		}
		item := parseSingleInventoryLine(l)
		if item.Name != "" {
			items = append(items, item)
		}
	}
	return FleetNodeInventory{
		NodeID: target.ID,
		Alias:  target.Alias,
		IP:     target.IP,
		OS:     target.OS,
		Items:  items,
	}
}

func isIgnoredInventoryLine(l string) bool {
	return l == "" || strings.HasPrefix(l, "#") || strings.HasPrefix(l, "===")
}

func parseSingleInventoryLine(l string) FleetInventoryItem {
	parts := strings.Fields(l)
	if len(parts) >= 2 {
		return FleetInventoryItem{
			Name:    parts[0],
			Version: parts[1],
			Status:  "active",
		}
	}
	return FleetInventoryItem{
		Name:    parts[0],
		Version: "unknown",
		Status:  "active",
	}
}

func executeDefaultRemoteInventory(target FleetTarget) (string, error) {
	if out, ok := tryRestRemoteInventory(target); ok {
		return out, nil
	}
	return executeSSHRemoteInventory(target)
}

func tryRestRemoteInventory(target FleetTarget) (string, bool) {
	url := fmt.Sprintf("http://%s:49152/api/v1/update/ls", target.IP)
	client := http.Client{Timeout: 2000 * time.Millisecond}
	resp, reqErr := client.Get(url)
	if reqErr != nil {
		return "", false
	}
	defer resp.Body.Close()
	isOk := resp.StatusCode == http.StatusOK
	if isOk {
		var buf strings.Builder
		_, _ = buf.ReadFrom(resp.Body)
		return buf.String(), true
	}
	return "", false
}

func executeSSHRemoteInventory(target FleetTarget) (string, error) {
	client, err := dialFleetSSH(target)
	if err != nil {
		return "", err
	}
	defer client.Close()

	cmd := resolveFleetInventoryCommand(target.OS)
	shell := resolveFleetShell(target.OS)
	return crypto.RunCommand(client, cmd, shell)
}

func resolveFleetInventoryCommand(osType string) string {
	isWin := strings.EqualFold(osType, "windows")
	if isWin {
		return "powershell -NoProfile -Command \"gitmap update ls --json 2>$null || gitmap list --json 2>$null || gitmap --version\""
	}
	return "gitmap update ls --json 2>/dev/null || gitmap list --json 2>/dev/null || gitmap --version"
}

func renderFleetLSSummary(results []FleetLSNodeResult, excludedCount int) {
	if len(results) == 0 {
		return
	}
	totalItems := countTotalInventoryItems(results)
	fmt.Printf("\n%s================================================================================%s\n",
		constants.ColorCyan, constants.ColorReset)
	fmt.Printf(" %sSSH Fleet Installed Software Inventory:%s Nodes: %d | Items: %d | Excluded: %d\n",
		constants.ColorBold, constants.ColorReset, len(results), totalItems, excludedCount)
	fmt.Printf("%s--------------------------------------------------------------------------------%s\n",
		constants.ColorDim, constants.ColorReset)

	cfg := termtable.TableConfig{
		Columns: []termtable.Column{
			{Title: "ALIAS", Align: termtable.AlignLeft, MinWidth: 15},
			{Title: "IP", Align: termtable.AlignLeft, MinWidth: 16},
			{Title: "PACKAGE / APP", Align: termtable.AlignLeft, MinWidth: 18},
			{Title: "VERSION", Align: termtable.AlignLeft, MinWidth: 12},
			{Title: "STATUS", Align: termtable.AlignLeft, MinWidth: 10},
			{Title: "DETAILS", Align: termtable.AlignLeft, MinWidth: 20},
		},
		Rows: buildFleetLSTableRows(results),
	}
	termtable.PrintTable(cfg)
	fmt.Printf("%s================================================================================%s\n\n",
		constants.ColorCyan, constants.ColorReset)
}

func countTotalInventoryItems(results []FleetLSNodeResult) int {
	total := 0
	for _, r := range results {
		total += len(r.Inventory.Items)
	}
	return total
}

func buildFleetLSTableRows(results []FleetLSNodeResult) []termtable.Row {
	var rows []termtable.Row
	for _, r := range results {
		rows = append(rows, buildSingleNodeLSRows(r)...)
	}
	return rows
}

func buildSingleNodeLSRows(r FleetLSNodeResult) []termtable.Row {
	hasItems := len(r.Inventory.Items) > 0
	isOperational := r.IsSuccess && hasItems
	if isOperational {
		return buildOperationalLSRows(r)
	}
	return buildFailedLSRow(r)
}

func buildFailedLSRow(r FleetLSNodeResult) []termtable.Row {
	errDesc := "unreachable"
	if r.Inventory.Error != "" {
		errDesc = r.Inventory.Error
	}
	statusStr := constants.ColorRed + "FAILED" + constants.ColorReset
	return []termtable.Row{{
		Cells: []string{
			r.Alias,
			r.IP,
			"-",
			"-",
			statusStr,
			errDesc,
		},
	}}
}

func buildOperationalLSRows(r FleetLSNodeResult) []termtable.Row {
	var rows []termtable.Row
	for _, item := range r.Inventory.Items {
		statusStr := constants.ColorGreen + item.Status + constants.ColorReset
		detail := item.InstallPath
		if detail == "" {
			detail = fmt.Sprintf("active (%dms)", r.DurationMs)
		}
		rows = append(rows, termtable.Row{
			Cells: []string{
				r.Alias,
				r.IP,
				item.Name,
				item.Version,
				statusStr,
				detail,
			},
		})
	}
	return rows
}
