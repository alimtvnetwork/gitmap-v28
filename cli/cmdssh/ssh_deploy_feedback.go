// Package cmdssh — ssh_deploy_feedback.go renders post-execution diagnostic docket tables and suggestions.
package cmdssh

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// DeployNodeOutcome records the execution result for an individual target node.
type DeployNodeOutcome struct {
	NodeName   string
	NodeIP     string
	NodeOS     string
	Status     string
	SkipReason string
	Detail     string
	Version    string
	Duration   time.Duration
}

// DeployFeedbackDocket captures comprehensive execution metrics across all targets.
type DeployFeedbackDocket struct {
	BinaryPath   string
	BinaryFormat BinaryFormatType
	TotalTargets int
	Succeeded    []DeployNodeOutcome
	SkippedOS    []DeployNodeOutcome
	SkippedOff   []DeployNodeOutcome
	Failed       []DeployNodeOutcome
	TotalTime    time.Duration
}

// NewDeployFeedbackDocket creates an initialized diagnostic docket for binary deployment.
func NewDeployFeedbackDocket(binPath string, format BinaryFormatType, totalTargets int) *DeployFeedbackDocket {
	return &DeployFeedbackDocket{
		BinaryPath:   binPath,
		BinaryFormat: format,
		TotalTargets: totalTargets,
	}
}

// AddSuccess records a successfully deployed node outcome.
func (d *DeployFeedbackDocket) AddSuccess(outcome DeployNodeOutcome) {
	d.Succeeded = append(d.Succeeded, outcome)
}

// AddSkippedOS records a node skipped due to OS incompatibility.
func (d *DeployFeedbackDocket) AddSkippedOS(outcome DeployNodeOutcome) {
	d.SkippedOS = append(d.SkippedOS, outcome)
}

// AddSkippedOffline records a node skipped due to being offline or unreachable.
func (d *DeployFeedbackDocket) AddSkippedOffline(outcome DeployNodeOutcome) {
	d.SkippedOff = append(d.SkippedOff, outcome)
}

// AddFailure records a failed deployment outcome.
func (d *DeployFeedbackDocket) AddFailure(outcome DeployNodeOutcome) {
	d.Failed = append(d.Failed, outcome)
}

// RenderDeployFeedbackDocket prints a formatted diagnostic summary table and suggestions.
func RenderDeployFeedbackDocket(d *DeployFeedbackDocket) {
	renderDocketHeader(d)
	renderDocketStats(d)
	renderSucceededSection(d)
	renderSkippedOSSection(d)
	renderSkippedOffSection(d)
	renderFailedSection(d)
	renderDocketFooter()
	renderDocketSuggestions(d)
}

func renderDocketHeader(d *DeployFeedbackDocket) {
	fmt.Printf("\n%s%s📋 GitMap Fleet Binary Deployment Diagnostic Docket%s\n", constants.ColorCyan, constants.ColorBold, constants.ColorReset)
	fmt.Println("  ────────────────────────────────────────────────────────────")
	fmt.Printf("  • Local Binary:   %s (Format: %s)\n", d.BinaryPath, d.BinaryFormat)
	fmt.Printf("  • Total Targets:  %d node(s)\n", d.TotalTargets)
	fmt.Printf("  • Execution Time: %s\n\n", d.TotalTime)
}

func renderDocketStats(d *DeployFeedbackDocket) {
	fmt.Printf("  ┌─── Summary Statistics ─────────────────────────────────────┐\n")
	fmt.Printf("  │ Total: %-3d │ Succeeded: %-3d │ Skipped OS: %-3d │ Offline: %-3d │ Failed: %-3d │\n",
		d.TotalTargets, len(d.Succeeded), len(d.SkippedOS), len(d.SkippedOff), len(d.Failed))
	fmt.Printf("  └─── Result Breakdown ──────────────────────────────────────┘\n\n")
}

func renderSucceededSection(d *DeployFeedbackDocket) {
	hasItems := len(d.Succeeded) > 0
	if !hasItems {
		return
	}

	fmt.Printf("  %s✓ Succeeded (%d):%s\n", constants.ColorGreen, len(d.Succeeded), constants.ColorReset)
	for _, o := range d.Succeeded {
		renderSingleSuccess(o)
	}

	fmt.Println()
}

func renderSingleSuccess(o DeployNodeOutcome) {
	verInfo := ""
	if o.Version != "" {
		verInfo = fmt.Sprintf(" [%s]", o.Version)
	}

	fmt.Printf("    • %s (%s)%s in %s\n", o.NodeName, o.NodeIP, verInfo, o.Duration)
}

func renderSkippedOSSection(d *DeployFeedbackDocket) {
	hasItems := len(d.SkippedOS) > 0
	if !hasItems {
		return
	}

	fmt.Printf("  %s⊘ Skipped (OS Incompatible) (%d):%s\n", constants.ColorYellow, len(d.SkippedOS), constants.ColorReset)
	for _, o := range d.SkippedOS {
		renderSingleSkippedOS(o)
	}

	fmt.Println()
}

func renderSingleSkippedOS(o DeployNodeOutcome) {
	fmt.Printf("    • %s (%s) [%s]\n", o.NodeName, o.NodeIP, o.NodeOS)
	fmt.Printf("      %sRationale: %s%s\n", constants.ColorDim, o.Detail, constants.ColorReset)
}

func renderSkippedOffSection(d *DeployFeedbackDocket) {
	hasItems := len(d.SkippedOff) > 0
	if !hasItems {
		return
	}

	fmt.Printf("  %s⊘ Skipped (Offline / Unreachable) (%d):%s\n", constants.ColorYellow, len(d.SkippedOff), constants.ColorReset)
	for _, o := range d.SkippedOff {
		renderSingleSkippedOff(o)
	}

	fmt.Println()
}

func renderSingleSkippedOff(o DeployNodeOutcome) {
	fmt.Printf("    • %s (%s)\n", o.NodeName, o.NodeIP)
	fmt.Printf("      %sRationale: %s%s\n", constants.ColorDim, o.Detail, constants.ColorReset)
}

func renderFailedSection(d *DeployFeedbackDocket) {
	hasItems := len(d.Failed) > 0
	if !hasItems {
		return
	}

	fmt.Printf("  %s✖ Failed (%d):%s\n", constants.ColorRed, len(d.Failed), constants.ColorReset)
	for _, o := range d.Failed {
		fmt.Printf("    • %s (%s): %s\n", o.NodeName, o.NodeIP, o.Detail)
	}

	fmt.Println()
}

func renderDocketFooter() {
	fmt.Println("  ────────────────────────────────────────────────────────────")
}

func renderDocketSuggestions(d *DeployFeedbackDocket) {
	hasSuggestions := len(d.SkippedOS) > 0 || len(d.SkippedOff) > 0
	if !hasSuggestions {
		return
	}

	fmt.Printf("\n  %s💡 Next-Step Actionable Suggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	renderOSIncompatSuggestions(d)
	renderOfflineSuggestions(d)
	renderForceOverrideSuggestion(d)

	fmt.Println()
}

func renderOSIncompatSuggestions(d *DeployFeedbackDocket) {
	hasSkippedOS := len(d.SkippedOS) > 0
	if !hasSkippedOS {
		return
	}

	renderCrossCompileSuggestion(d.BinaryFormat)
}

func renderCrossCompileSuggestion(format BinaryFormatType) {
	if format == BinaryFormatPE {
		fmt.Println("    • Cross-compile Linux binary for skipped Linux nodes:")
		fmt.Println("      GOOS=linux GOARCH=amd64 go build -o gitmap-linux ./cli")
		fmt.Println("      gitmap ssh deploy bin <linux-target> --file ./gitmap-linux")
	}

	if format == BinaryFormatELF {
		fmt.Println("    • Cross-compile Windows executable for skipped Windows nodes:")
		fmt.Println("      GOOS=windows GOARCH=amd64 go build -o gitmap.exe ./cli")
		fmt.Println("      gitmap ssh deploy bin <win-target> --file ./gitmap.exe")
	}
}

func renderOfflineSuggestions(d *DeployFeedbackDocket) {
	hasOffline := len(d.SkippedOff) > 0
	if !hasOffline {
		return
	}

	first := d.SkippedOff[0]
	fmt.Println("    • Verify network reachability and SSH daemon on offline nodes:")
	fmt.Printf("      Test-NetConnection %s -Port 22  (or: ping %s)\n", first.NodeIP, first.NodeIP)
	fmt.Printf("      gitmap ssh probe %s\n", first.NodeName)
}

func renderForceOverrideSuggestion(d *DeployFeedbackDocket) {
	hasSkippedOS := len(d.SkippedOS) > 0
	if !hasSkippedOS {
		return
	}

	fmt.Println("    • Override compatibility check if running under emulation:")
	fmt.Println("      gitmap ssh deploy bin <target> --force (or: --force-copy)")
}
