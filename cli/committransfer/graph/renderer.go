package graph

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

// GraphEvent captures a commit or merge event to render in the ASCII execution graph.
type GraphEvent struct {
	CommitSha  string `json:"commitSha"`
	BranchName string `json:"branchName"`
	IsMerge    bool   `json:"isMerge"`
	PRNumber   int    `json:"prNumber,omitempty"`
	ReleaseTag string `json:"releaseTag,omitempty"`
	Message    string `json:"message"`
}

// RenderExecutionGraph renders an ASCII/Unicode visual execution graph depicting
// mainline commits, feature branch forks, PR merge nodes, and release tags.
func RenderExecutionGraph(events []GraphEvent) string {
	if len(events) == 0 {
		return ""
	}

	var sb strings.Builder
	hasFeatureBranch := detectFeatureBranch(events)
	if hasFeatureBranch {
		renderBranchedGraph(&sb, events)
	} else {
		renderLinearGraph(&sb, events)
	}

	renderEventDetails(&sb, events)

	return sb.String()
}

func detectFeatureBranch(events []GraphEvent) bool {
	for _, ev := range events {
		if ev.IsMerge || ev.PRNumber > 0 {
			return true
		}
	}

	return false
}

func renderLinearGraph(sb *strings.Builder, events []GraphEvent) {
	sb.WriteString(pterm.LightCyan("  main: "))
	for i, ev := range events {
		renderCommitNode(sb, ev)
		if i < len(events)-1 {
			sb.WriteString(pterm.Gray("───"))
		}
	}
	sb.WriteString(pterm.LightGreen(" (HEAD)\n"))
}

func renderCommitNode(sb *strings.Builder, ev GraphEvent) {
	if ev.ReleaseTag != "" {
		sb.WriteString(pterm.Yellow("◆[" + ev.ReleaseTag + "]"))

		return
	}
	sb.WriteString(pterm.Cyan("●"))
}

func renderBranchedGraph(sb *strings.Builder, events []GraphEvent) {
	sb.WriteString(pterm.LightCyan("  main:   "))
	renderMainlineTrack(sb, events)
	sb.WriteString("\n")
	sb.WriteString(pterm.LightYellow("  branch: "))
	renderFeatureTrack(sb, events)
	sb.WriteString("\n")
}

func renderMainlineTrack(sb *strings.Builder, events []GraphEvent) {
	for i, ev := range events {
		if ev.IsMerge {
			sb.WriteString(pterm.Green("M[PR#" + fmt.Sprintf("%d", ev.PRNumber) + "]"))
		} else {
			renderCommitNode(sb, ev)
		}
		if i < len(events)-1 {
			sb.WriteString(pterm.Gray("───"))
		}
	}
	sb.WriteString(pterm.LightGreen(" (HEAD)"))
}

func renderFeatureTrack(sb *strings.Builder, events []GraphEvent) {
	for i, ev := range events {
		switch {
		case ev.IsMerge:
			sb.WriteString(pterm.Yellow("▲───────┘"))
		case ev.PRNumber > 0:
			sb.WriteString(pterm.Yellow("●───────"))
		default:
			sb.WriteString(pterm.Gray("        "))
		}
		if i < len(events)-1 {
			sb.WriteString("   ")
		}
	}
}

func renderEventDetails(sb *strings.Builder, events []GraphEvent) {
	sb.WriteString("\n" + pterm.Gray("  Node Details:") + "\n")
	for idx, ev := range events {
		shortSha := ev.CommitSha
		if len(shortSha) > 7 {
			shortSha = shortSha[:7]
		}
		prefix := pterm.Cyan(fmt.Sprintf("    [%d] %s", idx+1, shortSha))
		desc := formatEventDescription(ev)
		sb.WriteString(fmt.Sprintf("%s %s\n", prefix, desc))
	}
}

func formatEventDescription(ev GraphEvent) string {
	msg := ev.Message
	if len(msg) > 60 {
		msg = msg[:57] + "..."
	}
	if ev.IsMerge {
		return pterm.Green(fmt.Sprintf("[MERGE PR #%d] %s", ev.PRNumber, msg))
	}
	if ev.ReleaseTag != "" {
		return pterm.Yellow(fmt.Sprintf("[TAG %s] %s", ev.ReleaseTag, msg))
	}

	return pterm.White(msg)
}
