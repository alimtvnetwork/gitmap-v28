// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import "fmt"

// PrintNodesCloneHelp outputs the modern CLI help for fleet clone commands.
func PrintNodesCloneHelp(kind NodesCloneKind) error {
	kindStr := string(kind)
	printHelpHeaderBox(kindStr)
	printUsageSection(kindStr)
	printDescriptionSection()
	printParametersSection(kindStr)
	printFlagsSection()
	printFlagsFilterSection()
	printExamplesSection(kindStr)
	printExamplesTelemetrySection(kindStr)
	printAliasesSection(kindStr)
	return nil
}

func toUpperStr(s string) string {
	switch s {
	case "cfr":
		return "CFR (CLONE-FIX-REPO)"
	case "cfrp":
		return "CFRP (CLONE-FIX-REPO-PUB)"
	default:
		return "CLONE"
	}
}

func printHelpHeaderBox(kindStr string) {
	fmt.Println()
	fmt.Println("  ┌──────────────────────────────────────────────────────────────────────────────────────────────────────┐")
	fmt.Printf("  │ GITMAP NODES %-86s│\n", toUpperStr(kindStr)+" - ASYNC FLEET CLONE & MANIFEST STAGING")
	fmt.Println("  └──────────────────────────────────────────────────────────────────────────────────────────────────────┘")
	fmt.Println()
}

func printUsageSection(kindStr string) {
	fmt.Println("  Usage:")
	fmt.Printf("    gitmap nodes %s [flags] <repo|url|file> [dest]\n", kindStr)
	fmt.Printf("    gitmap nodes %s except-self <repo|url|file> [dest]\n", kindStr)
	fmt.Printf("    gitmap nodes %s [flags] <repo1,repo2,...> [dest]\n", kindStr)
	fmt.Printf("    gitmap nodes %s --except <node> <repo|url|file> [dest]\n", kindStr)
	fmt.Printf("    gitmap nodes %s --accept <node> <repo|url|file> [dest]\n", kindStr)
	fmt.Printf("    gitmap nodes %s [flags]\n\n", kindStr)
}

func printDescriptionSection() {
	fmt.Println("  Description:")
	fmt.Println("    Executes repository cloning asynchronously across local host and remote nodes:")
	fmt.Println("      1. Subfolder Mirroring: When invoked inside a local work subfolder (e.g. D:\\work\\sub),")
	fmt.Println("         the relative path is replicated on remote nodes (<work-root>/sub).")
	fmt.Println("      2. Target Directory Pre-Creation: Remote directories are created automatically")
	fmt.Println("         prior to cloning (mkdir -p / New-Item -ItemType Directory).")
	fmt.Println("      3. Custom Destination: Passing [dest] or --dest <path> overrides target directory.")
	fmt.Println("      4. Except-Self Mode: Skips local execution to run strictly on remote workers.")
	fmt.Println("      5. Manifest Staging: Passing a .json manifest stages it to all nodes before cloning.")
	fmt.Println()
}

func printParametersSection(kindStr string) {
	fmt.Println("  Supported Targets & Input Formats:")
	fmt.Printf("    • Repository Name Only:      gitmap nodes %s owner/repo [dest]\n", kindStr)
	fmt.Printf("    • Single Git URL:            gitmap nodes %s https://github.com/user/project.git [dest]\n", kindStr)
	fmt.Printf("    • Except-Self Subcommand:    gitmap nodes %s except-self https://github.com/user/project.git\n", kindStr)
	fmt.Printf("    • Comma-Separated Targets:   gitmap nodes %s repo1,repo2 [dest]\n", kindStr)
	fmt.Printf("    • Manifest File (.json):     gitmap nodes %s gitmap.json\n\n", kindStr)
}

func printFlagsSection() {
	fmt.Println("  Flags:")
	fmt.Println("    -t, --target string       Target specific node alias or IP (default: all nodes)")
	fmt.Println("        --exclude string      Exclude nodes by alias or IP (comma separated)")
	fmt.Println("        --except string       Exclude specific node alias or IP from execution")
	fmt.Println("        --accept string       Accept and dispatch only to specified node alias or IP")
	fmt.Println("    -d, --dest, --dir string  Explicit destination directory for cloned repositories")
	fmt.Println("        --target-dir string   Alias for --dest")
}

func printFlagsFilterSection() {
	fmt.Println("        --except-self         Execute only on remote nodes, skipping local master host")
	fmt.Println("        --no-self             Alias for --except-self")
	fmt.Println("        --skip-local          Alias for --except-self")
	fmt.Println("        --dry-run             Simulate actions without modifying filesystem")
	fmt.Println("    -j, --json                Output machine-readable JSON telemetry")
	fmt.Println("    -h, --help                Show this help message")
	fmt.Println()
}

func printExamplesSection(kindStr string) {
	fmt.Println("  Examples:")
	fmt.Printf("    gitmap nodes %s ChrisTitusTech/winutil                 # Clone repo by name across fleet\n", kindStr)
	fmt.Printf("    gitmap nodes %s https://github.com/u/repo D:\\work\\sub   # Clone URL into custom destination path\n", kindStr)
	fmt.Printf("    gitmap nodes %s except-self ChrisTitusTech/winutil     # Clone strictly to remote fleet nodes\n", kindStr)
	fmt.Printf("    gitmap nodes %s --dest ~/work/special user/repo        # Clone with destination flag\n", kindStr)
	fmt.Printf("    gitmap nodes %s gitmap.json                            # Copy manifest file and clone fleet-wide\n", kindStr)
	fmt.Printf("    gitmap nodes %s --target w1 user/repo                  # Dispatch only to targeted node\n", kindStr)
	fmt.Printf("    gitmap nodes %s --exclude w3 user/repo                 # Exclude node from execution\n", kindStr)
}

func printExamplesTelemetrySection(kindStr string) {
	fmt.Printf("    gitmap nodes %s --accept w1,w2 user/repo               # Accept and target specific nodes\n", kindStr)
	fmt.Println("    gitmap machine --ssh                                   # Query machine telemetry and SSH identity")
	fmt.Println("    gitmap nodes ping                                      # Probe fleet nodes ping and inspect latency")
	fmt.Println("    gitmap ssh trust w1                                    # Accept and trust remote host SSH keys")
	fmt.Printf("    gitmap nodes %s --dry-run user/repo                    # Preview fleet execution without changes\n\n", kindStr)
}

func printAliasesSection(kindStr string) {
	fmt.Println("  Aliases:")
	fmt.Printf("    gitmap node %s, gitmap nodes-%s, gitmap fleet-nodes %s\n\n", kindStr, kindStr, kindStr)
}
