package cmdssh

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func parseFleetPullOutcome(name, ip string, isLocal bool, rawOutput string) FleetNodePullOutcome {
	jsonStr := extractJSONPayload(rawOutput)
	var summary fleetPullJSONSummary
	if err := json.Unmarshal([]byte(jsonStr), &summary); err != nil {
		return buildFailedFleetOutcome(name, ip, isLocal, "invalid JSON: "+strings.TrimSpace(rawOutput))
	}
	active := filterActiveFleetStates(summary.States)

	return FleetNodePullOutcome{
		NodeName:     name,
		IP:           ip,
		IsLocal:      isLocal,
		Success:      true,
		Summary:      summary,
		ActiveStates: active,
	}
}

func extractJSONPayload(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end <= start {
		return ""
	}

	return strings.TrimSpace(s[start : end+1])
}

func filterActiveFleetStates(states []fleetPullRepoItem) []fleetPullRepoItem {
	var active []fleetPullRepoItem
	for _, s := range states {
		if isFleetRepoItemActive(s) {
			active = append(active, s)
		}
	}

	return active
}

func isFleetRepoItemActive(item fleetPullRepoItem) bool {
	status := strings.ToLower(item.Status)
	changes := strings.ToLower(item.Changes)
	if status == "up-to-date" || changes == "up-to-date" || status == "synced" || changes == "synced" {
		return false
	}
	if status == "failed" || status == "dirty" {
		return true
	}
	if strings.HasPrefix(item.Changes, "+") {
		return true
	}

	return status != "" && status != "up-to-date"
}

func renderFleetPullAllOutcomes(outcomes []FleetNodePullOutcome, isJSON bool) error {
	if isJSON {
		data, _ := json.MarshalIndent(map[string]any{"fleet": outcomes}, "", "  ")
		fmt.Println(string(data))

		return nil
	}
	for _, o := range outcomes {
		renderSingleFleetOutcome(o)
	}
	fmt.Println()

	return nil
}

func renderSingleFleetOutcome(o FleetNodePullOutcome) {
	if o.IsSkipped {
		return
	}
	if !o.Success {
		fmt.Printf("  %s▶ Node [%s] (IP: %s): %sfailed%s\n",
			constants.ColorYellow, o.NodeName, o.IP, constants.ColorRed, constants.ColorReset)
		fmt.Printf("      %s(%s)%s\n", constants.ColorDim, o.ErrorMsg, constants.ColorReset)

		return
	}
	activeCount := len(o.ActiveStates)
	upToDateCount := o.Summary.PulledCount - activeCount
	if upToDateCount < 0 {
		upToDateCount = 0
	}
	printNodeOutcomeHeader(o, activeCount, upToDateCount)
	if activeCount == 0 {
		fmt.Printf("      %s(all repositories are up-to-date)%s\n", constants.ColorDim, constants.ColorReset)

		return
	}
	renderActiveStatesList(o.ActiveStates)
}

func printNodeOutcomeHeader(o FleetNodePullOutcome, activeCount, upToDateCount int) {
	nodeTitle := fmt.Sprintf("[%s] (IP: %s)", o.NodeName, o.IP)
	if o.IsLocal {
		nodeTitle = fmt.Sprintf("Local VM (%s - localhost)", o.IP)
	}
	if activeCount > 0 {
		fmt.Printf("\n  %s• %s:%s %d pulled (%d active, %d up-to-date)\n",
			constants.ColorCyan, nodeTitle, constants.ColorReset,
			o.Summary.PulledCount, activeCount, upToDateCount)

		return
	}
	fmt.Printf("\n  %s• %s:%s %d pulled (all up-to-date)\n",
		constants.ColorCyan, nodeTitle, constants.ColorReset, o.Summary.PulledCount)
}

func renderActiveStatesList(states []fleetPullRepoItem) {
	colWidth := resolveFleetStatesColWidth(states)
	for i, s := range states {
		renderSingleFleetItem(s, colWidth, i == len(states)-1)
	}
}

func renderSingleFleetItem(s fleetPullRepoItem, colWidth int, isLast bool) {
	conn, subPrefix := resolveFleetItemConnectors(isLast)
	styled := styleFleetStatus(s.Status, s.Changes)
	fmt.Printf("    %s %-*s  %s\n", conn, colWidth, s.RepoName, styled)
	if s.ErrorDetails != "" {
		fmt.Printf("    %s ↳ Reason: %s%s\n", subPrefix, s.ErrorDetails, constants.ColorReset)
	}
	if s.RemediationHint != "" {
		fmt.Printf("    %s ↳ Next Step: %s%s\n", subPrefix, s.RemediationHint, constants.ColorReset)
	}
}

func resolveFleetItemConnectors(isLast bool) (string, string) {
	if isLast {
		return "└──", "   "
	}

	return "├──", "│  "
}

func resolveFleetStatesColWidth(states []fleetPullRepoItem) int {
	colWidth := 26
	for _, s := range states {
		if len(s.RepoName) > colWidth {
			colWidth = len(s.RepoName)
		}
	}
	return colWidth
}

func styleFleetStatus(status, changes string) string {
	if status == "failed" || changes == "failed" {
		return constants.ColorRed + "failed" + constants.ColorReset
	}
	if status == "dirty" || changes == "dirty" {
		return constants.ColorYellow + "dirty" + constants.ColorReset
	}
	if strings.HasPrefix(changes, "+") {
		return formatDiffStats(changes)
	}

	return changes
}

func formatDiffStats(s string) string {
	parts := strings.SplitN(s, " ", 2)
	diffPart := parts[0]
	slashIdx := strings.Index(diffPart, "/")
	if slashIdx == -1 {
		return constants.ColorGreen + s + constants.ColorReset
	}
	ins := diffPart[:slashIdx]
	del := diffPart[slashIdx:]
	colored := constants.ColorGreen + ins + constants.ColorReset + constants.ColorRed + del + constants.ColorReset
	if len(parts) > 1 {
		return colored + " " + constants.ColorDim + parts[1] + constants.ColorReset
	}

	return colored
}
