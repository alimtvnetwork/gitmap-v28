package cmdssh

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// ExecuteFleetSendCommitsAction sends commits across all affected nodes.
func ExecuteFleetSendCommitsAction(issues []FleetNodeIssue, conns []db.SSHConnection) {
	fmt.Printf("\n  %s⚡ Sending commits for dirty repositories across affected nodes...%s\n",
		constants.ColorCyan, constants.ColorReset)
	cmd := "gitmap sends cp all \"wip: save changes\""
	results := DispatchCommandToAffectedNodes(cmd, issues, conns)
	RenderNodeExecResults(results, "Send Commits")
}

// ExecuteFleetSendAllAction sends commits across all online fleet nodes.
func ExecuteFleetSendAllAction(conns []db.SSHConnection) {
	fmt.Printf("\n  %s⚡ Sending commits across all online fleet nodes...%s\n",
		constants.ColorCyan, constants.ColorReset)
	cmd := "gitmap sends cp all \"wip: save changes\""
	var allIssues []FleetNodeIssue
	for _, c := range conns {
		allIssues = append(allIssues, FleetNodeIssue{NodeName: c.Alias, IP: c.IPAddress, IsLocal: false})
	}
	allIssues = append(allIssues, FleetNodeIssue{NodeName: "Local VM", IP: "127.0.0.1", IsLocal: true})
	results := DispatchCommandToAffectedNodes(cmd, allIssues, conns)
	RenderNodeExecResults(results, "Send All Commits")
}

// ExecuteFleetStashAllAction runs git stash across all affected nodes.
func ExecuteFleetStashAllAction(issues []FleetNodeIssue, conns []db.SSHConnection) {
	fmt.Printf("\n  %s📦 Stashing uncommitted changes across affected nodes...%s\n",
		constants.ColorYellow, constants.ColorReset)
	cmd := "gitmap stash"
	results := DispatchCommandToAffectedNodes(cmd, issues, conns)
	RenderNodeExecResults(results, "Stash All")
}

// ExecuteFleetFixAllAction runs gitmap fix --all across all affected nodes.
func ExecuteFleetFixAllAction(issues []FleetNodeIssue, conns []db.SSHConnection) {
	fmt.Printf("\n  %s🔧 Running gitmap fix --all across affected nodes...%s\n",
		constants.ColorCyan, constants.ColorReset)
	cmd := "gitmap fix --all"
	results := DispatchCommandToAffectedNodes(cmd, issues, conns)
	RenderNodeExecResults(results, "Fix All")
}

// ExecuteFleetRetryPullAction re-executes gitmap pa across affected nodes and returns new issues.
func ExecuteFleetRetryPullAction(issues []FleetNodeIssue, conns []db.SSHConnection) []FleetNodeIssue {
	fmt.Printf("\n  %s🔄 Retrying pull-all across affected nodes...%s\n",
		constants.ColorMagenta, constants.ColorReset)
	cmd := "gitmap pa --json"
	results := DispatchCommandToAffectedNodes(cmd, issues, conns)
	var newOutcomes []FleetNodePullOutcome
	for _, r := range results {
		newOutcomes = append(newOutcomes, parseFleetPullOutcome(r.NodeName, r.IP, r.IsLocal, r.Output))
	}
	renderFleetPullAllOutcomes(newOutcomes, false)

	return ExtractFleetIssues(newOutcomes)
}

// ExecuteFleetStepThroughAction steps through affected nodes individually.
func ExecuteFleetStepThroughAction(issues []FleetNodeIssue, conns []db.SSHConnection) {
	for i, issue := range issues {
		fmt.Printf("\n  [%d/%d] Inspecting Node: %s%s%s (%s)\n",
			i+1, len(issues), constants.ColorBold, issue.NodeName, constants.ColorReset, issue.IP)
		renderSingleNodeIssueDetails(issue)
		handleSingleNodeStepChoice(issue, conns)
	}
}

func handleSingleNodeStepChoice(issue FleetNodeIssue, conns []db.SSHConnection) {
	fmt.Printf("    Actions: [1] Send commits  [2] Stash  [3] Fix all  [4] Skip\n")
	fmt.Printf("    %sAction for [%s] [1/2/3/4]:%s ", constants.ColorCyan, issue.NodeName, constants.ColorReset)
	choice := readSingleNodeStepChoice()
	conn, hasConn := resolveConnectionForNode(issue.NodeName, issue.IP, conns)
	if !issue.IsLocal && !hasConn {
		fmt.Println("    ✖ Connection not found, skipping.")
		return
	}
	dispatchSingleNodeStepChoice(choice, conn, issue.IsLocal)
}

func readSingleNodeStepChoice() string {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "4"
	}

	return strings.TrimSpace(line)
}

func dispatchSingleNodeStepChoice(choice string, conn db.SSHConnection, isLocal bool) {
	cmd := ""
	switch choice {
	case "1", "s", "send":
		cmd = "gitmap sends cp all \"wip: save changes\""
	case "2", "t", "stash":
		cmd = "gitmap stash"
	case "3", "f", "fix":
		cmd = "gitmap fix --all"
	default:
		return
	}
	res := ExecuteCommandOnFleetNode(conn, isLocal, cmd)
	renderSingleRemediationResult(res)
}
