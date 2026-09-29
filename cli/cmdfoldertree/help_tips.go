package cmdfoldertree

import "fmt"

func printFolderTreeTips() {
	fmt.Println("  [tip] Tip: Works on ANY directory on disk — no prior database scan required.")
	fmt.Println("  [tip] Tip: Git repositories are automatically tagged with branch: [git: main].")
	fmt.Println("  [tip] Tip: Use 'gitmap ft --preview' for clean numbered prompt-ready list with one-line gap.")
	fmt.Println("  [tip] Tip: Use 'gitmap ft export folder' or '--dirs-only' to export folder structure only.")
	fmt.Println("  [tip] Tip: Use 'gitmap ft export -f json -o skeleton.json' to export directory trees.")
	fmt.Println("  [tip] Tip: Use 'gitmap ft import skeleton.json ./new-project' to replicate scaffolding.")
	fmt.Println("  [tip] Tip: Use 'gitmap ft import folder skeleton.json ./target' to create only folders.")
	fmt.Println()
}

func printFolderTreeExamples() {
	fmt.Println("  Examples:")
	fmt.Println("    gitmap ft                                 # Tree view of current directory")
	fmt.Println("    gitmap ft ls --preview                    # Numbered list with one-line gap")
	fmt.Println("    gitmap ft export folder -o dirs.txt       # Export folder structure only")
	fmt.Println("    gitmap ft export -f yaml -o schema.yaml   # Export full tree to YAML")
	fmt.Println("    gitmap ft import schema.yaml ./target     # Recreate folders & placeholder files")
	fmt.Println("    gitmap ft import folder dirs.txt ./target # Recreate directory structure only")
	fmt.Println()
}
