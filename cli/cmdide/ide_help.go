package cmdide

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcursor"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

// RunHelp displays a formatted reference for gitmap ide commands.
func RunHelp() error {
	fmt.Printf("%sGitMap IDE Suite%s — Cross-IDE Repository Synchronization\n\nUsage:\n  gitmap ide <command> [flags]\n\n", constants.ColorCyan, constants.ColorReset)
	printHelpDetails()
	return nil
}

func printHelpDetails() {
	printHelpCommands()
	printHelpFlags()
	printHelpExamples()
}

func printHelpCommands() {
	fmt.Print("Commands:\n")
	fmt.Print("  setup [target]   Configure AI IDEs (cline, continue, zed, all) (alias: cursor-setup, setup-ai, configure)\n")
	fmt.Print("  add <path>       Add repository to IDEs (--target <ide>)\n")
	fmt.Print("  sync [dir]       Synchronize repositories across IDEs\n")
	fmt.Print("  remove <path>    Remove repository from IDEs (alias: rm)\n")
	fmt.Print("  list             List repositories and status (alias: ls)\n")
	fmt.Print("  status           Display overview of installed IDEs\n")
	fmt.Print("  help             Show this help reference\n\n")
}

func printHelpFlags() {
	fmt.Print("Flags:\n")
	fmt.Print("  -p, --provider <name> AI provider preset (kimi, openrouter, deepseek, anthropic, ollama, custom)\n")
	fmt.Print("  -m, --model <id>      Model identifier override\n")
	fmt.Print("  -b, --base-url <url>  Custom endpoint base URL override\n")
	fmt.Print("  -k, --api-key <key>   API key for the provider\n")
	fmt.Print("      --test            Test connectivity before configuring\n")
	fmt.Print("  -t, --target <ide>    Target specific IDE (vscode, cursor, agy, desktop, all)\n")
	fmt.Print("  -n, --dry-run         Simulate operations without writing changes\n")
	fmt.Print("  -j, --json            Emit machine-readable JSON output\n")
	fmt.Print("  -q, --quiet           Suppress informational messages\n\n")
}

func printHelpExamples() {
	fmt.Print("Examples:\n")
	fmt.Print("  gitmap ide setup --provider kimi --api-key sk-xxxxxx\n")
	fmt.Print("  gitmap ide setup all --provider ollama --model deepseek-r1\n")
	fmt.Print("  gitmap ide setup cline --dry-run\n")
	fmt.Print("  gitmap ide setup --test --provider deepseek -k sk-xxxxxx\n\n")
}

func runIDEStatus(args []string) error {
	opts, _ := parseIDEOptions(args)
	items := collectIDEStatusItems()
	if opts.IsJSON {
		b, _ := json.MarshalIndent(items, "", "  ")
		fmt.Println(string(b))
		return nil
	}
	printStatusCards(items)
	return nil
}

func collectIDEStatusItems() []IDEStatusItem {
	vCfg, _ := vscodepm.ProjectsJSONPath()
	vEntries, _ := vscodepm.ListEntries()
	cCfg, _ := cmdcursor.GetCursorProjectsJSONPath()
	cEntries, _ := vscodepm.ReadEntries(cCfg)
	cli := desktop.ResolveCLI()
	return []IDEStatusItem{
		probeEditorStatus("VS Code", "code", vCfg, len(vEntries)),
		probeEditorStatus("Cursor", "cursor", cCfg, len(cEntries)),
		probeAntigravityStatus(),
		{Name: "GitHub Desktop", IsInstalled: cli != "", ExecutablePath: cli},
	}
}

func probeEditorStatus(name, exeName, cfg string, count int) IDEStatusItem {
	exe, _ := exec.LookPath(exeName)
	return IDEStatusItem{
		Name:            name,
		IsInstalled:     exe != "" || cfg != "",
		ExecutablePath:  exe,
		ConfigPath:      cfg,
		RegisteredCount: count,
	}
}

func probeAntigravityStatus() IDEStatusItem {
	home, _ := os.UserHomeDir()
	cfg := filepath.Join(home, ".gemini", "config", "projects")
	entries, _ := os.ReadDir(cfg)
	count := 0
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			count++
		}
	}
	return IDEStatusItem{Name: "Antigravity", IsInstalled: count > 0, ConfigPath: cfg, RegisteredCount: count}
}

func printStatusCards(items []IDEStatusItem) {
	fmt.Printf("\n%s● Installed Developer IDEs & Status:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, item := range items {
		statusStr := constants.ColorGreen + "Installed" + constants.ColorReset
		if !item.IsInstalled {
			statusStr = constants.ColorDim + "Not Detected" + constants.ColorReset
		}
		fmt.Printf("  • %-16s %s (repos: %d)\n", item.Name, statusStr, item.RegisteredCount)
		if item.ConfigPath != "" {
			fmt.Printf("    config: %s\n", item.ConfigPath)
		}
	}
	fmt.Println()
}
