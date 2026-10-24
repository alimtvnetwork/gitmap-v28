package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type fleetPullRepoItem struct {
	RepoName        string `json:"repoName"`
	Status          string `json:"status"`
	Changes         string `json:"changes"`
	ErrorDetails    string `json:"errorDetails,omitempty"`
	RemediationHint string `json:"remediationHint,omitempty"`
}

type fleetPullJSONSummary struct {
	Total        int                 `json:"total"`
	PulledCount  int                 `json:"pulledCount"`
	SuccessCount int                 `json:"successCount"`
	FailedCount  int                 `json:"failedCount"`
	States       []fleetPullRepoItem `json:"states"`
	DurationMs   int64               `json:"durationMs"`
}

// FleetNodePullOutcome represents the pull result for a local or remote node.
type FleetNodePullOutcome struct {
	NodeName     string               `json:"node"`
	IP           string               `json:"ip"`
	IsLocal      bool                 `json:"isLocal"`
	Success      bool                 `json:"success"`
	IsSkipped    bool                 `json:"isSkipped,omitempty"`
	ErrorMsg     string               `json:"errorMsg,omitempty"`
	Summary      fleetPullJSONSummary `json:"summary,omitempty"`
	ActiveStates []fleetPullRepoItem  `json:"activeStates,omitempty"`
}

type fleetNodeLiveness struct {
	conn        db.SSHConnection
	isOnline    bool
	isLocalHost bool
	version     string
}

// RunSSHPullAllFleet executes pull-all across SSH fleet nodes and Local VM concurrently.
func RunSSHPullAllFleet(cleanArgs []string) error {
	isJSON := hasFleetJSONFlag(cleanArgs)
	outcomes, conns := probeAndExecuteFleetPull(cleanArgs, isJSON)
	if renderErr := renderFleetPullAllOutcomes(outcomes, isJSON); renderErr != nil {
		return renderErr
	}
	if tryRemediateFleetIssues(outcomes, conns, isJSON) {
		return nil
	}

	return checkFleetFailures(outcomes)
}

func probeAndExecuteFleetPull(cleanArgs []string, isJSON bool) ([]FleetNodePullOutcome, []db.SSHConnection) {
	conns, _ := fetchAllSSHConnections()
	probes := probeFleetLiveness(conns)
	if !isJSON {
		printFleetEnqueueBanner(probes)
	}

	return executeFleetPullAll(probes, cleanArgs), conns
}

func tryRemediateFleetIssues(outcomes []FleetNodePullOutcome, conns []db.SSHConnection, isJSON bool) bool {
	if !hasFleetRemediationTrigger(outcomes, isJSON) {
		return false
	}
	remErr := PromptFleetPostPullRemediation(outcomes, conns)

	return remErr == nil
}

func checkFleetFailures(outcomes []FleetNodePullOutcome) error {
	failures := countFleetFailures(outcomes)
	if failures > 0 {
		return apperror.NewSimple(fmt.Sprintf("SSH fleet pull-all had %d node failure(s)", failures), "E_SSH_FLEET_FAIL")
	}

	return nil
}

func hasFleetRemediationTrigger(outcomes []FleetNodePullOutcome, isJSON bool) bool {
	if isJSON || !IsInteractiveFleetPromptEligible() {
		return false
	}
	issues := ExtractFleetIssues(outcomes)

	return HasAnyFleetIssues(issues)
}

func countFleetFailures(outcomes []FleetNodePullOutcome) int {
	failures := 0
	for _, o := range outcomes {
		if hasOutcomeFailed(o) {
			failures++
		}
	}
	return failures
}

func hasOutcomeFailed(o FleetNodePullOutcome) bool {
	if o.IsSkipped {
		return false
	}
	return !o.Success
}

func hasFleetJSONFlag(args []string) bool {
	for _, a := range args {
		if strings.EqualFold(a, "--json") || strings.EqualFold(a, "-json") {
			return true
		}
	}

	return false
}

func buildSkippedFleetOutcome(name, ip string) FleetNodePullOutcome {
	return FleetNodePullOutcome{
		NodeName:  name,
		IP:        ip,
		IsLocal:   false,
		Success:   false,
		IsSkipped: true,
		ErrorMsg:  "offline (skipped, no task enqueued)",
	}
}

func buildLocalSkippedFleetOutcome(name, ip string) FleetNodePullOutcome {
	return FleetNodePullOutcome{
		NodeName:  name,
		IP:        ip,
		IsLocal:   true,
		Success:   true,
		IsSkipped: true,
		ErrorMsg:  "local machine (covered by local execution)",
	}
}

func buildFailedFleetOutcome(name, ip string, isLocal bool, msg string) FleetNodePullOutcome {
	logFleetFailure(isLocal, msg, ip)
	return FleetNodePullOutcome{
		NodeName: name,
		IP:       ip,
		IsLocal:  isLocal,
		Success:  false,
		ErrorMsg: msg,
	}
}

func logFleetFailure(isLocal bool, msg, ip string) {
	errType := "SSH_DELEGATION_ERROR"
	if isLocal {
		errType = "LOCAL_PULL_ERROR"
	}
	store.LogInternalError("PULL_ALL_SSH", errType, msg, ip, "")
}
