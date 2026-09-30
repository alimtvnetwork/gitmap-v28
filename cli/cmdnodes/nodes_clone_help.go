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
	fmt.Println("    Executes repository cloning simultaneously across the entire fleet:")
	fmt.Println("      1. Runs directly on the current host machine in-process.")
	fmt.Println("      2. Dispatches concurrently to all enrolled remote fleet nodes over SSH.")
	fmt.Println("      3. When a local manifest file (e.g. gitmap.json) is specified, gitmap")
	fmt.Println("         automatically stages the file to the default work directory")
	fmt.Println("         (D:\\work on Windows, ~/work on Linux/Unix) on all remote nodes")
	fmt.Println("         prior to triggering the clone pipeline in parallel.")
	fmt.Println()
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
	fmt.Printf("    gitmap nodes %s ChrisTitusTech/winutil\n", kindStr)
	fmt.Printf("    gitmap nodes %s https://github.com/user/project.git\n", kindStr)
	fmt.Printf("    gitmap nodes %s repo1,repo2,repo3\n", kindStr)
	fmt.Printf("    gitmap nodes %s gitmap.json\n", kindStr)
	fmt.Printf("    gitmap nodes %s                           # Auto-detects local gitmap.json\n", kindStr)
	fmt.Printf("    gitmap nodes %s --target w1 gitmap.json    # Target single remote node\n", kindStr)
	fmt.Printf("    gitmap nodes %s --dry-run gitmap.json     # Preview fleet execution\n", kindStr)
	fmt.Println()
	fmt.Println("  Aliases:")
	fmt.Printf("    gitmap node %s, gitmap nodes-%s, gitmap fleet-nodes %s\n\n", kindStr, kindStr, kindStr)
}
