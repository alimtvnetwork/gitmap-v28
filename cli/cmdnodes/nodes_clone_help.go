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
	fmt.Printf("    gitmap nodes %s [flags] <repo|url|file> [target_path]\n", kindStr)
	fmt.Printf("    gitmap nodes %s except-self <repo|url|file> [target_path]\n", kindStr)
	fmt.Printf("    gitmap nodes %s [flags] <repo1,repo2,...>\n", kindStr)
	fmt.Printf("    gitmap nodes %s [flags]\n\n", kindStr)
}

func printDescriptionSection() {
	fmt.Println("  Description:")
	fmt.Println("    Executes repository cloning asynchronously across local host and remote nodes:")
	fmt.Println("      1. Default Work Directory: Commands automatically navigate to the work directory")
	fmt.Println("         (D:\\work on Windows, ~/work on Linux/Unix) prior to cloning, ensuring clean")
	fmt.Println("         folder organization across every machine in the fleet.")
	fmt.Println("      2. Custom Destination Path: Passing [target_path] after the URL clones into that")
	fmt.Println("         explicit directory on all target nodes.")
	fmt.Println("      3. Except-Self Mode: Using 'except-self' or '--except-self' skips cloning on the")
	fmt.Println("         master/local host, running purely across remote fleet worker nodes.")
	fmt.Println("      4. Manifest Staging: Passing a file like gitmap.json automatically transfers it")
	fmt.Println("         to all remote work directories before parallel execution begins.")
	fmt.Println()
}

func printParametersSection(kindStr string) {
	fmt.Println("  Supported Targets & Input Formats:")
	fmt.Printf("    • Repository Name Only:      gitmap nodes %s owner/repo [dest] (e.g. ChrisTitusTech/winutil)\n", kindStr)
	fmt.Printf("    • Single Git URL:            gitmap nodes %s https://github.com/user/project.git [dest]\n", kindStr)
	fmt.Printf("    • Except-Self Subcommand:    gitmap nodes %s except-self https://github.com/user/project.git\n", kindStr)
	fmt.Printf("    • Comma-Separated Targets:   gitmap nodes %s repo1,repo2,https://github.com/user/repo3\n", kindStr)
	fmt.Printf("    • Manifest File (.json):     gitmap nodes %s gitmap.json (staged to remote work dirs first)\n", kindStr)
	fmt.Printf("    • Auto-Detected Manifest:    gitmap nodes %s (auto-discovers gitmap.json in current dir)\n\n", kindStr)
}

func printFlagsSection() {
	fmt.Println("  Flags:")
	fmt.Println("    -t, --target string     Target specific node alias or IP (default: all nodes)")
	fmt.Println("        --exclude string    Exclude nodes by alias or IP (comma separated)")
	fmt.Println("        --except-self       Execute only on remote nodes, skipping current local host")
	fmt.Println("        --no-self           Alias for --except-self")
	fmt.Println("        --skip-local        Alias for --except-self")
	fmt.Println("        --dry-run           Simulate clone actions without modifying filesystem")
	fmt.Println("    -j, --json              Output machine-readable JSON telemetry")
	fmt.Println("    -h, --help              Show this help message")
	fmt.Println()
}

func printExamplesSection(kindStr string) {
	fmt.Println("  Examples:")
	fmt.Printf("    gitmap nodes %s ChrisTitusTech/winutil                 # Clone repo by name across fleet\n", kindStr)
	fmt.Printf("    gitmap nodes %s https://github.com/user/repo.git       # Clone URL into default work dir\n", kindStr)
	fmt.Printf("    gitmap nodes %s https://github.com/user/repo.git D:\\w # Clone URL into custom target path\n", kindStr)
	fmt.Printf("    gitmap nodes %s except-self ChrisTitusTech/winutil     # Clone strictly to remote fleet nodes\n", kindStr)
	fmt.Printf("    gitmap nodes %s --except-self https://github.com/u/r   # Remote-only clone using flag\n", kindStr)
	fmt.Printf("    gitmap nodes %s gitmap.json                            # Copy local manifest & clone fleet-wide\n", kindStr)
	fmt.Printf("    gitmap nodes %s --target w1 gitmap.json                # Dispatch only to node 'w1'\n", kindStr)
	fmt.Printf("    gitmap nodes %s --dry-run gitmap.json                  # Preview fleet execution\n", kindStr)
	fmt.Println()
	fmt.Println("  Aliases:")
	fmt.Printf("    gitmap node %s, gitmap nodes-%s, gitmap fleet-nodes %s\n\n", kindStr, kindStr, kindStr)
}
