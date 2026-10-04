// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// NodesAgyQueryResult wraps the aggregated Antigravity query response across fleet.
type NodesAgyQueryResult struct {
	Projects      []cmdagy.RunningProjectRecord `json:"projects"`
	TotalProjects int                           `json:"totalProjects"`
	HasActiveIDE  bool                          `json:"hasActiveIDE"`
	IDEPID        int                           `json:"idePid,omitempty"`
	IDEName       string                        `json:"ideName,omitempty"`
	FleetNodes    []cmdagy.NodeStatusRecord     `json:"fleetNodes,omitempty"`
}

// RunNodesAgyQuery queries running Antigravity projects and IDE instances across nodes.
func RunNodesAgyQuery(args []string) error {
	projects, err := cmdagy.DiscoverRunningProjects()
	if err != nil {
		return apperror.WrapSimple(err, "discover running projects")
	}

	isJSON := isJSONRequested(args)
	if isSSHRequested(args) {
		return cmdagy.AggregateSSHRunningProjects(projects, isJSON, "")
	}

	queryRes := buildAgyQueryResult(projects)
	if isJSON {
		return printQueryResultJSON(queryRes)
	}

	renderAgyQueryBoxedDashboard(queryRes)
	return nil
}

func isJSONRequested(args []string) bool {
	for _, a := range args {
		tok := strings.TrimSpace(a)
		if tok == "--json" || tok == "-j" || tok == "-json" {
			return true
		}
	}
	return false
}

func isSSHRequested(args []string) bool {
	for _, a := range args {
		tok := strings.TrimSpace(a)
		if tok == "--ssh" || tok == "-s" {
			return true
		}
	}
	return false
}

func buildAgyQueryResult(projects []cmdagy.RunningProjectRecord) NodesAgyQueryResult {
	ideProc := cmdagy.DetectRunningAntigravityIDE()
	pid, name := extractIDEProcessDetails(ideProc)
	fleetNodes := cmdagy.CollectFleetNodeStatuses(len(projects))

	return NodesAgyQueryResult{
		Projects:      projects,
		TotalProjects: len(projects),
		HasActiveIDE:  ideProc.IsSuccess(),
		IDEPID:        pid,
		IDEName:       name,
		FleetNodes:    fleetNodes,
	}
}

func extractIDEProcessDetails(proc result.Result[cmdagy.AgyProcessInfo]) (int, string) {
	if proc.IsSuccess() {
		return proc.Value.PID, proc.Value.Name
	}
	return 0, ""
}

func printQueryResultJSON(queryRes NodesAgyQueryResult) error {
	data, err := json.MarshalIndent(queryRes, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal nodes agy query json")
	}

	fmt.Println(string(data))
	return nil
}

func renderAgyQueryBoxedDashboard(res NodesAgyQueryResult) {
	printBoxedBanner("FLEET ANTIGRAVITY RUNNING PROJECTS & INSTANCE QUERY", 80)
	printIDEInstanceStatus(res)
	if len(res.FleetNodes) > 1 {
		printFleetNodesSummary(res.FleetNodes)
	}
	cmdagy.RenderRunningProjectsTable(res.Projects)
}

func printFleetNodesSummary(nodes []cmdagy.NodeStatusRecord) {
	fmt.Printf("  %s● Fleet Nodes:%s ", constants.ColorCyan, constants.ColorReset)
	var labels []string
	for _, n := range nodes {
		indicator := constants.ColorGreen + "●" + constants.ColorReset
		if !n.IsOnline {
			indicator = constants.ColorDim + "○" + constants.ColorReset
		}
		labels = append(labels, fmt.Sprintf("%s %s (%s)", indicator, n.Alias, n.OS))
	}
	fmt.Println(strings.Join(labels, "  "))
	fmt.Println()
}

func printBoxedBanner(title string, width int) {
	inner := width - 4
	padLeft := (inner - len(title)) / 2
	padRight := inner - len(title) - padLeft
	fmt.Println()
	fmt.Printf("  %s┌%s┐%s\n", constants.ColorCyan, strings.Repeat("─", width-2), constants.ColorReset)
	fmt.Printf("  %s│ %s%s%s │%s\n", constants.ColorCyan, strings.Repeat(" ", padLeft), title, strings.Repeat(" ", padRight), constants.ColorReset)
	fmt.Printf("  %s└%s┘%s\n", constants.ColorCyan, strings.Repeat("─", width-2), constants.ColorReset)
}

func printIDEInstanceStatus(res NodesAgyQueryResult) {
	if res.HasActiveIDE {
		fmt.Printf("  %s●%s IDE Instance: %sRunning%s (PID: %d, Binary: %s)\n",
			constants.ColorGreen, constants.ColorReset,
			constants.ColorGreen, constants.ColorReset, res.IDEPID, res.IDEName)
		return
	}

	fmt.Printf("  %s○%s IDE Instance: %sOffline%s (no active Antigravity process detected)\n",
		constants.ColorYellow, constants.ColorReset,
		constants.ColorDim, constants.ColorReset)
}
