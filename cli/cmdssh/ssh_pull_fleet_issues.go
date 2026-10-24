package cmdssh

import (
	"strings"
)

// FleetRepoIssueItem represents a single repository that has an issue on a node.
type FleetRepoIssueItem struct {
	RepoName string
	Status   string
	Reason   string
	IsDirty  bool
	IsFailed bool
}

// FleetNodeIssue encapsulates node-level or repo-level failures for a single machine.
type FleetNodeIssue struct {
	NodeName     string
	IP           string
	IsLocal      bool
	NodeError    string
	HasNodeError bool
	RepoIssues   []FleetRepoIssueItem
}

// ExtractFleetIssues parses all node outcomes and collects machines with failures or dirty repos.
func ExtractFleetIssues(outcomes []FleetNodePullOutcome) []FleetNodeIssue {
	var issues []FleetNodeIssue
	for _, o := range outcomes {
		if issue, hasIssue := extractSingleNodeIssue(o); hasIssue {
			issues = append(issues, issue)
		}
	}

	return issues
}

func extractSingleNodeIssue(o FleetNodePullOutcome) (FleetNodeIssue, bool) {
	if o.IsSkipped {
		return FleetNodeIssue{}, false
	}
	if !o.Success {
		return buildNodeErrorIssue(o), true
	}
	repos := extractRepoIssuesFromStates(o.Summary.States)
	if len(repos) == 0 && len(o.ActiveStates) > 0 {
		repos = extractRepoIssuesFromStates(o.ActiveStates)
	}
	if len(repos) > 0 {
		return buildRepoIssuesNode(o, repos), true
	}

	return FleetNodeIssue{}, false
}

func buildNodeErrorIssue(o FleetNodePullOutcome) FleetNodeIssue {
	return FleetNodeIssue{
		NodeName:     o.NodeName,
		IP:           o.IP,
		IsLocal:      o.IsLocal,
		NodeError:    o.ErrorMsg,
		HasNodeError: true,
	}
}

func buildRepoIssuesNode(o FleetNodePullOutcome, repos []FleetRepoIssueItem) FleetNodeIssue {
	return FleetNodeIssue{
		NodeName:     o.NodeName,
		IP:           o.IP,
		IsLocal:      o.IsLocal,
		HasNodeError: false,
		RepoIssues:   repos,
	}
}

func extractRepoIssuesFromStates(states []fleetPullRepoItem) []FleetRepoIssueItem {
	var repos []FleetRepoIssueItem
	for _, item := range states {
		if isRepoStateProblematic(item) {
			repos = append(repos, buildFleetRepoIssue(item))
		}
	}

	return repos
}

func isRepoStateProblematic(item fleetPullRepoItem) bool {
	status := strings.ToLower(item.Status)
	changes := strings.ToLower(item.Changes)
	if status == "dirty" || changes == "dirty" || strings.Contains(changes, "dirty") {
		return true
	}
	if status == "failed" || status == "conflict" || status == "error" {
		return true
	}

	return item.ErrorDetails != "" || item.RemediationHint != ""
}

func buildFleetRepoIssue(item fleetPullRepoItem) FleetRepoIssueItem {
	status := strings.ToLower(item.Status)
	isDirty := status == "dirty" || strings.Contains(strings.ToLower(item.Changes), "dirty")
	isFailed := status == "failed" || status == "error" || item.ErrorDetails != ""

	return FleetRepoIssueItem{
		RepoName: item.RepoName,
		Status:   item.Status,
		Reason:   resolveRepoIssueReason(item),
		IsDirty:  isDirty,
		IsFailed: isFailed,
	}
}

func resolveRepoIssueReason(item fleetPullRepoItem) string {
	if item.ErrorDetails != "" {
		return item.ErrorDetails
	}
	if item.RemediationHint != "" {
		return item.RemediationHint
	}
	status := strings.ToLower(item.Status)
	if status == "dirty" || strings.Contains(strings.ToLower(item.Changes), "dirty") {
		return "working tree dirty / autostash conflict"
	}
	if status == "failed" {
		return "pull operation failed"
	}

	return "requires attention after pull"
}

// CountTotalFleetIssues returns the total count of node and repository issues.
func CountTotalFleetIssues(issues []FleetNodeIssue) int {
	count := 0
	for _, item := range issues {
		if item.HasNodeError {
			count++
		}
		count += len(item.RepoIssues)
	}

	return count
}

// HasAnyFleetIssues returns true if there is at least one issue in the slice.
func HasAnyFleetIssues(issues []FleetNodeIssue) bool {
	return len(issues) > 0
}
