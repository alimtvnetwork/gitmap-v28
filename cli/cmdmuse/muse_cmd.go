// Package cmdmuse provides CLI commands, onboarding, and help cards for the Muse AI agent persona.
package cmdmuse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	museTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	museCyanStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	museGreenStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	museMutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	// museTableStyle was deleted here on 31-Oct-26 - it was unused
	// (golangci-lint unused finding on the v6.643.0 baseline diff).
)

// MusePromptDoc holds metadata describing a specific Muse prompt template.
type MusePromptDoc struct {
	Name        string
	FileName    string
	Description string
	TriggerCmd  string
}

// AllMusePrompts returns the catalog of supported Muse prompts.
func AllMusePrompts() []MusePromptDoc {
	return []MusePromptDoc{
		{
			Name:        "Master Autonomous Prompt",
			FileName:    "01-muse-master-prompt.md",
			Description: "Master autonomous coordinator governing multi-agent tree swarm execution and state persistence.",
			TriggerCmd:  "gitmap muse prompt master",
		},
		{
			Name:        "Generic Engineering Prompt",
			FileName:    "02-muse-generic.md",
			Description: "Core software engineering, system architecture design, and guideline-compliant coding.",
			TriggerCmd:  "gitmap muse prompt generic",
		},
		{
			Name:        "Standalone Single-Repo Prompt",
			FileName:    "03-muse-standalone-master-prompt.md",
			Description: "Independent single-repository task executor operating without external swarm dependencies.",
			TriggerCmd:  "gitmap muse prompt standalone",
		},
		{
			Name:        "Universal Master Prompt",
			FileName:    "04-muse-universal-master-prompt.md",
			Description: "Polyglot multi-stack workflows, cross-language abstractions, and framework migrations.",
			TriggerCmd:  "gitmap muse prompt universal",
		},
		{
			Name:        "Atomic Single-Step Prompt",
			FileName:    "05-muse-execute-in-a-step.md",
			Description: "Fast atomic single-step execution, surgical bug fixes, and zero-overhead micro-tasks.",
			TriggerCmd:  "gitmap muse prompt step",
		},
		{
			Name:        "Verification & Ledger Prompt",
			FileName:    "09-execute-with-verification-letterly.md",
			Description: "Retrospective test verification, evidence ledger generation, and audit compliance.",
			TriggerCmd:  "gitmap muse prompt verify",
		},
		{
			Name:        "Automated Release Prompt",
			FileName:    "10-execute-with-release-letterly.md",
			Description: "Automated minor version bump release ceremony, changelog synchronization, and tag push.",
			TriggerCmd:  "gitmap muse prompt release",
		},
		{
			Name:        "Time-Restricted Recommit",
			FileName:    "14-execute-recommit-time-restricted-letterly.md",
			Description: "Time-window restricted recommits (< 08:00 AM weekdays) and historical commit replay.",
			TriggerCmd:  "gitmap muse prompt recommit",
		},
	}
}

// RunMuseCLI dispatches muse subcommands (help, onboarding, prompt).
func RunMuseCLI(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		RenderMuseHelp()
		return nil
	}

	subcommand := strings.ToLower(args[0])
	switch subcommand {
	case "onboarding", "onboard", "train", "setup":
		return RunMuseOnboarding()
	case "prompt", "pmt":
		return RunMusePrompt(args[1:])
	default:
		RenderMuseHelp()
		return nil
	}
}

// RunMuseOnboarding displays the Muse persona configuration and onboarding instructions.
func RunMuseOnboarding() error {
	fmt.Println()
	fmt.Println(museTitleStyle.Render("=== GitMap Muse Autonomous Persona Onboarding ==="))
	fmt.Println()

	fmt.Println(museGreenStyle.Render("● Bot Identity:"), "gitmap automation bot - muse - runs autonomously <info@rgitmap.com>")
	fmt.Println(museCyanStyle.Render("● Engine Architecture:"), "V7 Master Engine (Tree Swarm, Brain Roadmap, SQLite Task DB)")
	fmt.Println(museCyanStyle.Render("● Calendar Windows:"), "Weekdays < 08:00 AM (or --auto-window / --schedule); Weekends unconstrained")
	fmt.Println()

	fmt.Println(museTitleStyle.Render("■ Mandatory Autonomous Directives:"))
	fmt.Println("  1. " + museGreenStyle.Render("Mandatory Pre-Flight Pull:") + " Always run 'gitmap pull' before inspecting or editing code.")
	fmt.Println("  2. " + museGreenStyle.Render("Strict Relative Paths:") + " Reference all specs as 02-spec/..., plans as .ai-memory/...; NO absolute paths.")
	fmt.Println("  3. " + museGreenStyle.Render("Exclusively Use GitMap CLI:") + " gitmap aum search, gitmap find, gitmap cat, gitmap ps. (Banned: rg, grep, Select-String).")
	fmt.Println("  4. " + museGreenStyle.Render("Zero Scratch Scripts:") + " Never create throwaway Python/Node scripts in brain/.../scratch/ or sleep loops.")
	fmt.Println("  5. " + museGreenStyle.Render("Dedicated Semantic Commits:") + " Always use gitmap cpr, gitmap cpc, gitmap cbf, gitmap cfr, gitmap cpf with --ai-agent muse.")
	fmt.Println()

	return nil
}

// RenderMuseHelp displays all available Muse prompts and their descriptions in a Tokyo Night card.
func RenderMuseHelp() {
	prompts := AllMusePrompts()
	fmt.Println()
	fmt.Println(museTitleStyle.Render("=== GitMap Muse AI Prompt Catalog & Reference ==="))
	fmt.Println(museMutedStyle.Render("Authoritative prompts for the Muse autonomous persona (01-prompts/27-muse-prompts/ & 22-letterly/):"))
	fmt.Println()

	fmt.Printf("%-32s | %-65s\n", museCyanStyle.Render("MUSE PROMPT TYPE"), museCyanStyle.Render("PURPOSE & DESCRIPTION"))
	fmt.Println(strings.Repeat("-", 100))

	for _, p := range prompts {
		namePill := museGreenStyle.Render(p.Name)
		fmt.Printf("%-44s | %s\n", namePill, p.Description)
		fmt.Printf("   File: %-25s | Trigger: %s\n", museMutedStyle.Render(p.FileName), museCyanStyle.Render(p.TriggerCmd))
		fmt.Println(strings.Repeat(".", 100))
	}

	fmt.Println()
	fmt.Println(museMutedStyle.Render("Tip: Run 'gitmap muse onboarding' to inspect persona configuration and execution rules."))
	fmt.Println(museMutedStyle.Render("Tip: Run 'gitmap cg execute v7' to stream the master V7 autonomous execution standard."))
	fmt.Println()
}

// RunMusePrompt resolves and streams a specific Muse prompt by keyword.
func RunMusePrompt(args []string) error {
	if len(args) == 0 {
		RenderMuseHelp()
		return nil
	}

	targetKeyword := strings.ToLower(args[0])
	prompts := AllMusePrompts()
	for _, p := range prompts {
		if strings.Contains(strings.ToLower(p.Name), targetKeyword) || strings.Contains(strings.ToLower(p.TriggerCmd), targetKeyword) {
			return streamMusePromptFile(p)
		}
	}

	fmt.Fprintf(os.Stderr, "No Muse prompt matching %q found. Available prompts:\n", targetKeyword)
	RenderMuseHelp()

	return nil
}

func streamMusePromptFile(p MusePromptDoc) error {
	fmt.Println(museTitleStyle.Render(fmt.Sprintf("=== Muse Prompt: %s (%s) ===", p.Name, p.FileName)))
	candidates := []string{
		filepath.Join("01-prompts", "27-muse-prompts", p.FileName),
		filepath.Join("01-prompts", "22-letterly", p.FileName),
		filepath.Join("01-prompts", "24-muse-prompts", p.FileName),
	}

	for _, c := range candidates {
		if data, err := os.ReadFile(c); err == nil {
			fmt.Println(string(data))
			return nil
		}
	}

	fmt.Printf("(Prompt template %s not found in local working tree)\n", p.FileName)
	return nil
}
