// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import "fmt"

// PrintNodesCloneHelp outputs the modern CLI help for fleet clone commands.
func PrintNodesCloneHelp(kind NodesCloneKind) error {
	kindStr := string(kind)
	fmt.Println()
	fmt.Println("  ╔══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Printf("  ║ GITMAP NODES %-86s║\n", toUpperStr(kindStr)+" - ASYNC FLEET CLONE & MANIFEST STAGING")
	fmt.Println("  ╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	printUsageSection(kindStr)
	printDescriptionSection()
	printParametersSection(kindStr)
	printFlagsSection()
	printExamplesSection(kindStr)
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

func printUsageSection(kindStr string) {
	fmt.Println("  Usage:")
	fmt.Printf("    gitmap nodes %s [flags] <repo|url|file>\n", kindStr)
	fmt.Printf("    gitmap nodes %s [flags] <repo1,repo2,...>\n", kindStr)
	fmt.Printf("    gitmap nodes %s [flags]\n\n", kindStr)
}

func printDescriptionSection() {
	fmt.Println("  Description:")
	fmt.Println("    Executes repository cloning asynchronously across local host and remote nodes:")
	fmt.Println("      1. Runs directly on the current host machine in-process.")
	fmt.Println("      2. Dispatches concurrently to all enrolled remote fleet nodes over SSH.")
	fmt.Println("      3. Manifest File Staging: When a file like gitmap.json is provided from the")
	fmt.Println("         current directory, gitmap automatically copies it over SSH to the default")
	fmt.Println("         work directory (D:\\work on Windows, ~/work on Linux/Unix) on all remote")
	fmt.Println("         nodes, then executes the clone pipeline across every machine.")
	fmt.Println()
}

func printParametersSection(kindStr string) {
	fmt.Println("  Supported Targets & Input Formats:")
	fmt.Printf("    • Repository Name Only:      gitmap nodes %s owner/repo (e.g. ChrisTitusTech/winutil)\n", kindStr)
	fmt.Printf("    • Single Git URL:            gitmap nodes %s https://github.com/user/project.git\n", kindStr)
	fmt.Printf("    • Comma-Separated Targets:   gitmap nodes %s repo1,repo2,https://github.com/user/repo3\n", kindStr)
	fmt.Printf("    • Manifest File (.json):     gitmap nodes %s gitmap.json (staged to remote work dirs first)\n", kindStr)
	fmt.Printf("    • Auto-Detected Manifest:    gitmap nodes %s (auto-discovers gitmap.json in current dir)\n\n", kindStr)
}

func printFlagsSection() {
	fmt.Println("  Flags:")
	fmt.Println("    -t, --target string     Target specific node alias or IP (default: all nodes)")
	fmt.Println("        --exclude string    Exclude nodes by alias or IP (comma separated)")
	fmt.Println("        --skip-local        Execute only on remote nodes, skipping current machine")
	fmt.Println("        --dry-run           Simulate clone actions without modifying filesystem")
	fmt.Println("    -j, --json              Output machine-readable JSON telemetry")
	fmt.Println("    -h, --help              Show this help message")
	fmt.Println()
}

func printExamplesSection(kindStr string) {
	fmt.Println("  Examples:")
	fmt.Printf("    gitmap nodes %s ChrisTitusTech/winutil                 # Clone repo by name across fleet\n", kindStr)
	fmt.Printf("    gitmap nodes %s https://github.com/user/repo.git       # Clone full Git URL across fleet\n", kindStr)
	fmt.Printf("    gitmap nodes %s repo1,repo2,repo3                      # Clone multiple comma-separated repos\n", kindStr)
	fmt.Printf("    gitmap nodes %s gitmap.json                            # Copy local manifest to remote work dirs & clone\n", kindStr)
	fmt.Printf("    gitmap nodes %s                                        # Auto-detect local gitmap.json & clone fleet-wide\n", kindStr)
	fmt.Printf("    gitmap nodes %s --target w1 gitmap.json                # Dispatch only to node 'w1'\n", kindStr)
	fmt.Printf("    gitmap nodes %s --skip-local gitmap.json               # Execute on remote fleet nodes only\n", kindStr)
	fmt.Printf("    gitmap nodes %s --dry-run gitmap.json                  # Preview fleet execution without modifying disk\n", kindStr)
	fmt.Println()
	fmt.Println("  Aliases:")
	fmt.Printf("    gitmap node %s, gitmap nodes-%s, gitmap fleet-nodes %s\n\n", kindStr, kindStr, kindStr)
}
