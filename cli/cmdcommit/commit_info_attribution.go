package cmdcommit

import (
	"fmt"
	"os"
	"strings"
)

// DetectCommitOrigin determines the provenance and origin category of a commit.
func DetectCommitOrigin(rep *CommitInfoReport) CommitProvenance {
	prov := rep.Provenance
	hasAiTraces := prov.IsAiCommitted || rep.CrossReferences.HasSplitDbRecord || rep.CrossReferences.HasHeatmapRecord
	if !hasAiTraces {
		return buildHumanProvenance(prov)
	}

	localHost, _ := os.Hostname()
	isLocal := isLocalMachineHost(prov.MachineHostname, localHost)
	if isLocal {
		return buildLocalAiProvenance(prov, localHost)
	}

	return buildForeignAiProvenance(prov)
}

func isLocalMachineHost(machineHost, localHost string) bool {
	if len(machineHost) == 0 {
		return true
	}

	return strings.EqualFold(machineHost, localHost)
}

func buildHumanProvenance(prov CommitProvenance) CommitProvenance {
	prov.OriginType = OriginHumanOrExternal
	prov.IsLocalMachine = false
	prov.IsAiCommitted = false
	prov.AttributionText = "The current machine's AI has not committed this (likely another worker/node/human)"

	return prov
}

func buildLocalAiProvenance(prov CommitProvenance, localHost string) CommitProvenance {
	prov.OriginType = OriginLocalAI
	prov.IsLocalMachine = true
	prov.IsAiCommitted = true
	if len(prov.MachineHostname) == 0 {
		prov.MachineHostname = localHost
	}
	model := resolveModelDisplay(prov)
	prov.AttributionText = fmt.Sprintf("Committed by local AI agent on this machine (Host: %s, AI: %s)", prov.MachineHostname, model)

	return prov
}

func buildForeignAiProvenance(prov CommitProvenance) CommitProvenance {
	prov.OriginType = OriginForeignAI
	prov.IsLocalMachine = false
	prov.IsAiCommitted = true
	model := resolveModelDisplay(prov)
	prov.AttributionText = fmt.Sprintf("Committed by remote cluster worker AI (Host: %s, AI: %s)", prov.MachineHostname, model)

	return prov
}

func resolveModelDisplay(prov CommitProvenance) string {
	if len(prov.AiModel) > 0 {
		return prov.AiModel
	}
	if len(prov.AiName) > 0 {
		return prov.AiName
	}

	return "Unknown AI"
}
