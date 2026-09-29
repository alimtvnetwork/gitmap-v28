package cmdfoldertree

import "fmt"

// PrintFolderTreeHelp prints the interactive help menu for folder-tree.
func PrintFolderTreeHelp() {
	printFolderTreeBanner()
	printFolderTreeUsage()
	printFolderTreeSubcommands()
	printFolderTreeDisplayModes()
	printFolderTreeFlags()
	printFolderTreeTips()
}

func printFolderTreeBanner() {
	fmt.Println()
	fmt.Println("  ╔═════════════════════════════════════════════════════════════╗")
	fmt.Println("  ║   Folder & Repo Tree Engine (gitmap folder-tree / ft)       ║")
	fmt.Println("  ╚═════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func printFolderTreeUsage() {
	fmt.Println("  Usage:")
	fmt.Println("    gitmap folder-tree [path] [flags]                  (default tree view)")
	fmt.Println("    gitmap folder-tree ls [path] [flags]               (list/tree view)")
	fmt.Println("    gitmap folder-tree view [path] [flags]             (alias for ls)")
	fmt.Println("    gitmap folder-tree export [path] [flags]           (export hierarchy)")
	fmt.Println("    gitmap folder-tree import <file> [target] [flags]  (recreate structure)")
	fmt.Println("    gitmap ft [command] [args]                         (short alias)")
	fmt.Println()
}

func printFolderTreeSubcommands() {
	fmt.Println("  Commands:")
	fmt.Println("    ls / view [path]           Render directory hierarchy with Git detection & sequence numbers")
	fmt.Println("    export [path]              Export hierarchy to JSON, YAML, tree text, or one-line gap preview")
	fmt.Println("    import <file> [target-dir] Recreate directory skeleton and empty placeholder files from file")
	fmt.Println("    help                       Show this comprehensive help menu")
	fmt.Println()
}

func printFolderTreeDisplayModes() {
	fmt.Println("  Display Modes:")
	fmt.Println("    --tree                     Emoji tree view with branch connectors (default)")
	fmt.Println("    --preview (-p, --gap)      Numbered sequence with one-line gap (name + path)")
	fmt.Println("    --dirs-only (-d)           Display or export directories only, omit files")
	fmt.Println("    --files (-f)               Include files in tree rendering (default: true)")
	fmt.Println("    --git-only                 Filter to display only Git repositories")
	fmt.Println()
}

func printFolderTreeFlags() {
	fmt.Println("  Flags:")
	fmt.Println("    --depth <N> (-L <N>)       Maximum scan recursion depth (default: 3, 0 for infinite)")
	fmt.Println("    --format <fmt> (-f <fmt>)  Export format: json (default), yaml, tree, preview")
	fmt.Println("    --output <file> (-o <f>)   Write exported output directly to specified file")
	fmt.Println("    --no-numbers               Suppress sequence numbering in output")
	fmt.Println("    --all (-a)                 Include hidden items (.git internals, node_modules)")
	fmt.Println("    --dry-run (-n)             Preview import actions without modifying disk")
	fmt.Println("    -h, --help                 Show this help menu")
	fmt.Println()
}

func printFolderTreeTips() {
	fmt.Println("  [tip] Tip: Works on ANY directory on disk — no prior database scan required.")
	fmt.Println("  [tip] Tip: Git repositories are automatically tagged with branch: [git: main].")
	fmt.Println("  [tip] Tip: Use 'gitmap ft --preview' for clean numbered prompt-ready list with one-line gap.")
	fmt.Println("  [tip] Tip: Use 'gitmap ft export -f json -o skeleton.json' to export directory trees.")
	fmt.Println("  [tip] Tip: Use 'gitmap ft import skeleton.json ./new-project' to replicate scaffolding.")
	fmt.Println()
}
