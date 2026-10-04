package cmdpushfix

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/glyphs"
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// PushFixBadgeType designates the status badge state.
type PushFixBadgeType string

const (
	BadgeDiagnosing PushFixBadgeType = "DIAGNOSING"
	BadgeAuthError  PushFixBadgeType = "AUTH_ERROR"
	BadgeDiverged   PushFixBadgeType = "DIVERGED"
	BadgeTransport  PushFixBadgeType = "TRANSPORT"
	BadgeLeaseForce PushFixBadgeType = "LEASE_FORCE"
	BadgeDryRun     PushFixBadgeType = "DRY_RUN"
	BadgeSuccess    PushFixBadgeType = "SUCCESS"
)

var safePushFixBadges = map[PushFixBadgeType]string{
	BadgeDiagnosing: "[DIAGNOSE]",
	BadgeAuthError:  "[AUTH-FAIL]",
	BadgeDiverged:   "[NON-FF]",
	BadgeTransport:  "[TRANSPORT]",
	BadgeLeaseForce: "[FORCE-LEASE]",
	BadgeDryRun:     "[DRY-RUN]",
	BadgeSuccess:    "[SUCCESS]",
}

var richPushFixBadges = map[PushFixBadgeType]string{
	BadgeDiagnosing: "[ 🔍 DIAGNOSING ]",
	BadgeAuthError:  "[ ✗ AUTH ERROR ]",
	BadgeDiverged:   "[ ⚠ DIVERGED ]",
	BadgeTransport:  "[ ⚡ PROTOCOL ]",
	BadgeLeaseForce: "[ 🛡️ FORCE-LEASE ]",
	BadgeDryRun:     "[ ℹ DRY-RUN ]",
	BadgeSuccess:    "[ ✓ RECOVERED ]",
}

var badgeColors = map[PushFixBadgeType]string{
	BadgeDiagnosing: constants.ColorCyan,
	BadgeAuthError:  constants.ColorRed,
	BadgeDiverged:   constants.ColorYellow,
	BadgeTransport:  constants.ColorMagenta,
	BadgeLeaseForce: constants.ColorYellow,
	BadgeDryRun:     constants.ColorBlue,
	BadgeSuccess:    constants.ColorGreen,
}

var phaseStepIcons = map[int]string{
	1: "🔍 ",
	2: "🔑 ",
	3: "🔄 ",
	4: "🚀 ",
}

// PushFixDiagnosticCard encapsulates details for rendering a diagnostic summary.
type PushFixDiagnosticCard struct {
	Category      string
	RemoteTarget  string
	ActiveBranch  string
	RootCause     string
	PlannedAction string
}

// FormatPushFixBadge returns plain badge string adhering to safe or rich glyph mode.
func FormatPushFixBadge(badge PushFixBadgeType, isSafe bool) string {
	if isSafe {
		if text, hasText := safePushFixBadges[badge]; hasText {
			return text
		}
		return "[STATUS]"
	}
	if text, hasText := richPushFixBadges[badge]; hasText {
		return text
	}
	return "[STATUS]"
}

// ResolveBadgeColor looks up ANSI color code for the badge type.
func ResolveBadgeColor(badge PushFixBadgeType) string {
	if color, hasColor := badgeColors[badge]; hasColor {
		return color
	}
	return constants.ColorCyan
}

// ColorizePushFixBadge decorates badge with ANSI color and resets terminal formatting.
func ColorizePushFixBadge(badge PushFixBadgeType, isSafe bool) string {
	color := ResolveBadgeColor(badge)
	text := FormatPushFixBadge(badge, isSafe)
	return fmt.Sprintf("%s%s%s", color, text, constants.ColorReset)
}

// RenderPushFixBadge dynamically resolves console mode and renders a colorized badge.
func RenderPushFixBadge(badge PushFixBadgeType) string {
	isSafe := glyphs.Resolve() == glyphs.ModeSafe
	return ColorizePushFixBadge(badge, isSafe)
}

func resolvePhaseStepIcon(step int, isSafe bool) string {
	if isSafe {
		return ""
	}
	return phaseStepIcons[step]
}

// PrintPushFixHeader displays the diagnosis and remediation HUD header.
func PrintPushFixHeader(repoDir string) {
	isSafe := glyphs.Resolve() == glyphs.ModeSafe
	icon := "⚡"
	if isSafe {
		icon = "*"
	}
	fmt.Printf("\n  %s%s%s %s[PUSH-FIX]%s %sAutonomous Push & Auth Recovery%s\n",
		constants.ColorBold, constants.ColorCyan, icon,
		constants.ColorBold, constants.ColorReset,
		constants.ColorDim, constants.ColorReset,
	)
	fmt.Printf("  Target Directory: %s%s%s\n\n", constants.ColorDim, repoDir, constants.ColorReset)
}

// PrintPhaseStep outputs a numbered execution stage card.
func PrintPhaseStep(step int, name, detail string) {
	isSafe := glyphs.Resolve() == glyphs.ModeSafe
	icon := resolvePhaseStepIcon(step, isSafe)
	fmt.Printf("  [%d/4] %s%s%-24s%s • %s\n",
		step, icon,
		constants.ColorCyan, name, constants.ColorReset,
		detail,
	)
}

// PrintPhaseStepWithStatus outputs a numbered execution stage card with status badge.
func PrintPhaseStepWithStatus(step int, name, status string) {
	isSafe := glyphs.Resolve() == glyphs.ModeSafe
	icon := resolvePhaseStepIcon(step, isSafe)
	fmt.Printf("  [%d/4] %s%-55s %s%s%s\n",
		step, icon, name+"...",
		constants.ColorGreen, status, constants.ColorReset,
	)
}

// PrintRemediationApplied prints an auto-fix action badge.
func PrintRemediationApplied(action string) {
	isSafe := glyphs.Resolve() == glyphs.ModeSafe
	badge := "✔ [FIX APPLIED]"
	if isSafe {
		badge = "* [FIX-APPLIED]"
	}
	fmt.Printf("        %s%s%s %s\n",
		constants.ColorGreen, badge, constants.ColorReset,
		action,
	)
}

// PrintPushFixSuccess renders the final completion summary card.
func PrintPushFixSuccess(remote, branch, sha string, isDryRun bool) {
	badge := RenderPushFixBadge(BadgeSuccess)
	status := "Pushed to GitHub successfully"
	if isDryRun {
		badge = RenderPushFixBadge(BadgeDryRun)
		status = "Simulated push (dry-run mode) — ready for live push"
	}
	fmt.Printf("\n  %s %s\n", badge, status)
	fmt.Printf("    • Remote: %s%s%s\n", constants.ColorCyan, remote, constants.ColorReset)
	fmt.Printf("    • Branch: %s%s%s\n", constants.ColorCyan, branch, constants.ColorReset)
	if len(sha) > 0 {
		fmt.Printf("    • Commit: %s%s%s\n", constants.ColorDim, sha, constants.ColorReset)
	}
	fmt.Println()
}

// RenderPushFixDiagnosticCard prints a boxed diagnostic summary card.
func RenderPushFixDiagnosticCard(card PushFixDiagnosticCard) {
	isSafe := glyphs.Resolve() == glyphs.ModeSafe
	renderCardBorders(card, isSafe)
}

func renderCardBorders(card PushFixDiagnosticCard, isSafe bool) {
	fmt.Println()
	renderCardTopBorder(isSafe)
	renderCardFields(card, isSafe)
	renderCardBottomBorder(isSafe)
	fmt.Println()
}

func renderCardTopBorder(isSafe bool) {
	color := constants.ColorCyan
	reset := constants.ColorReset
	if isSafe {
		fmt.Printf("  %s+-- Push Failure Diagnosis ------------------------------------------+%s\n", color, reset)
		return
	}
	fmt.Printf("  %s┌── Push Failure Diagnosis ──────────────────────────────────────────┐%s\n", color, reset)
}

func renderCardBottomBorder(isSafe bool) {
	color := constants.ColorCyan
	reset := constants.ColorReset
	if isSafe {
		fmt.Printf("  %s+--------------------------------------------------------------------+%s\n", color, reset)
		return
	}
	fmt.Printf("  %s└────────────────────────────────────────────────────────────────────┘%s\n", color, reset)
}

func renderCardFields(card PushFixDiagnosticCard, isSafe bool) {
	renderCardFieldLine("Issue Category", card.Category, isSafe)
	renderCardFieldLine("Remote Target", card.RemoteTarget, isSafe)
	renderCardFieldLine("Active Branch", card.ActiveBranch, isSafe)
	renderCardFieldLine("Root Cause", card.RootCause, isSafe)
	renderCardFieldLine("Planned Action", card.PlannedAction, isSafe)
}

func renderCardFieldLine(label, value string, isSafe bool) {
	border := "│"
	if isSafe {
		border = "|"
	}
	color := constants.ColorCyan
	reset := constants.ColorReset
	fmt.Printf("  %s%s%s %-14s : %s\n", color, border, reset, label, value)
}

// RenderPushFixHelpMenu renders the rich CLI help menu for push-fix.
func RenderPushFixHelpMenu() {
	termhelp.RenderMenu(BuildPushFixHelpMenu())
}

func buildPushFixUsage() []string {
	return []string{
		"gitmap push-fix [flags]",
		"gitmap pf [flags]",
		"gitmap push fix [flags]",
	}
}

// BuildPushFixHelpMenu constructs structured help configuration conforming to spec §6.3.
func BuildPushFixHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title:      "GITMAP PUSH-FIX & AUTH RECOVERY",
		UsageLines: buildPushFixUsage(),
		Sections:   []termhelp.HelpSection{buildPushFixRemediationSection(), buildPushFixOptionsSection()},
		Tips:       []string{"Run 'gitmap pf -n' to safely preview the repair plan before modifying remotes."},
	}
}

func buildPushFixRemediationSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Autonomous Remediation",
		Color: constants.ColorCyan,
		Entries: []termhelp.CommandEntry{
			{Command: "gitmap pf", Description: "Diagnose and auto-repair push failures"},
			{Command: "gitmap pf --dry-run", Description: "Simulate recovery actions safely"},
			{Command: "gitmap pf --ssh", Description: "Convert remote to SSH and push"},
			{Command: "gitmap pf --https", Description: "Convert remote to HTTPS and push"},
			{Command: "gitmap pf --force", Description: "Safe push with --force-with-lease"},
		},
	}
}

func buildPushFixOptionsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Options",
		Color: constants.ColorWhite,
		Entries: []termhelp.CommandEntry{
			{Command: "-n, --dry-run", Description: "Preview fixes without making changes"},
			{Command: "-r, --remote <name>", Description: "Target remote (default: origin)"},
			{Command: "-b, --branch <name>", Description: "Target branch (default: current)"},
			{Command: "-f, --force", Description: "Push with --force-with-lease"},
			{Command: "-y, --yes", Description: "Non-interactive bypass for prompts"},
		},
	}
}
