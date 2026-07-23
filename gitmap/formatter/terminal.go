// Package formatter renders ScanRecords to terminal output.
package formatter

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/render"
)

// Terminal writes a professional colored output to the given writer.
// When quiet is true, the clone help section is suppressed (useful for CI/scripts).
//
// The "Per-Repo Summary" section emits one standardized block per repo
// via render.RenderRepoTermBlocks. The block is intentionally rendered
// without color so the output is grep-friendly across all four
// commands that share it (scan, clone-from, clone-next, probe).
func Terminal(w io.Writer, records []model.ScanRecord, outputDir string, quiet bool) error {
	printBanner(w, len(records))
	printRepoList(w, records)
	printRepoSummaryBlocks(w, records)
	printFolderTree(w, records)
	printOutputFiles(w, outputDir)
	if quiet {
		return nil
	}
	printCloneHelp(w)

	return nil
}

// printRepoSummaryBlocks writes the standardized per-repo summary
// section. Empty input is a no-op so quiet/edge runs stay clean.
func printRepoSummaryBlocks(w io.Writer, records []model.ScanRecord) {
	if len(records) == 0 {
		return
	}
	fmt.Fprintf(w, constants.ColorYellow+"  ■ Per-Repo Summary"+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	_ = render.RenderRepoTermBlocks(w, render.FromScanRecords(records))
	fmt.Fprintln(w)
}

// printBanner writes the header section.
func printBanner(w io.Writer, count int) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, constants.ColorCyan+constants.TermBannerTop+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorCyan+constants.TermBannerTitle+constants.ColorReset+"\n", constants.Version)
	fmt.Fprintf(w, constants.ColorCyan+constants.TermBannerBottom+constants.ColorReset+"\n")
	fmt.Fprintln(w)
	fmt.Fprintf(w, constants.ColorGreen+constants.TermFoundFmt+constants.ColorReset+"\n", count)
	fmt.Fprintln(w)
}

// printRepoList writes each repo with folder name and clone instruction.
func printRepoList(w io.Writer, records []model.ScanRecord) {
	fmt.Fprintf(w, constants.ColorYellow+constants.TermReposHeader+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	for i, r := range records {
		printOneRepo(w, r, i+1, len(records))
	}
	fmt.Fprintln(w)
}

// printOneRepo writes a single repo entry with index.
func printOneRepo(w io.Writer, r model.ScanRecord, idx, total int) {
	fmt.Fprintf(w, constants.ColorDim+"  %d/%d "+constants.ColorReset, idx, total)
	fmt.Fprintf(w, constants.ColorGreen+"📦 %s"+constants.ColorReset, r.RepoName)
	fmt.Fprintf(w, constants.ColorDim+" (%s)"+constants.ColorReset+"\n", r.Branch)
	fmt.Fprintf(w, constants.ColorDim+"       └─ "+constants.ColorReset)
	fmt.Fprintf(w, constants.ColorWhite+"%s"+constants.ColorReset+"\n", r.CloneInstruction)
}

// printFolderTree writes the folder structure to terminal.
func printFolderTree(w io.Writer, records []model.ScanRecord) {
	fmt.Fprintf(w, constants.ColorYellow+constants.TermTreeHeader+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	paths := collectTermPaths(records)
	tree := buildTermTree(paths)
	renderTermTree(w, tree, "  ")
	fmt.Fprintln(w)
}

// printOutputFiles shows the generated output files.
func printOutputFiles(w io.Writer, outputDir string) {
	fmt.Fprintf(w, constants.ColorYellow+"  ■ Output Files"+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+"  📁 %s/"+constants.ColorReset+"\n", outputDir)
	printOutputFile(w, outputDir, constants.DefaultCSVFile, "Repo data in CSV")
	printOutputFile(w, outputDir, constants.DefaultJSONFile, "Repo data in JSON")
	printOutputFile(w, outputDir, constants.DefaultStructureFile, "Folder tree")
	printOutputFile(w, outputDir, constants.DefaultCloneScript, "PowerShell clone script")
	printOutputFile(w, outputDir, constants.DefaultDirectCloneScript, "Plain git clone commands (HTTPS)")
	printOutputFile(w, outputDir, constants.DefaultDirectCloneSSHScript, "Plain git clone commands (SSH)")
	printOutputFile(w, outputDir, constants.DefaultDesktopScript, "GitHub Desktop registration")
	fmt.Fprintln(w)
}

// printOutputFile shows one output file entry.
func printOutputFile(w io.Writer, dir, name, desc string) {
	fullPath := filepath.Join(dir, name)
	fmt.Fprintf(w, constants.ColorDim+"  ├── "+constants.ColorReset)
	fmt.Fprintf(w, constants.ColorCyan+"📄 %s"+constants.ColorReset, name)
	fmt.Fprintf(w, constants.ColorDim+"  %s"+constants.ColorReset+"\n", desc)
	_ = fullPath
}

// printCloneHelp writes instructions for cloning on another machine.
func printCloneHelp(w io.Writer) {
	fmt.Fprintf(w, constants.ColorYellow+constants.TermCloneHeader+constants.ColorReset+"\n")
	fmt.Fprintf(w, constants.ColorDim+constants.TermSeparator+constants.ColorReset+"\n")
	printCloneStep(w, constants.TermCloneStep1, constants.TermCloneCmd1)
	printCloneStepMulti(w, constants.TermCloneStep2, constants.TermCloneCmd2, constants.TermCloneCmd2Alt)
	printCloneStepMulti(w, constants.TermCloneStep3, constants.TermCloneCmd3, constants.TermCloneCmd3Alt)
	printCloneStepMulti(w, constants.TermCloneStep3t, constants.TermCloneCmd3t, constants.TermCloneCmd3tAlt)
	printCloneStep(w, constants.TermCloneStep3b, constants.TermCloneCmd3b)
	printCloneStepMulti(w, constants.TermCloneStep4, constants.TermCloneCmd4HTTPS, constants.TermCloneCmd4SSH)
	printCloneStep(w, constants.TermCloneStep5, constants.TermCloneCmd5)
	printCloneStep(w, constants.TermCloneStep6, constants.TermCloneCmd6)
	fmt.Fprintf(w, constants.ColorDim+constants.TermCloneNote+constants.ColorReset+"\n")
	fmt.Fprintln(w)
}

// printCloneStep writes a single step with one command.
func printCloneStep(w io.Writer, step, cmd string) {
	fmt.Fprintf(w, "%s%s%s\n", constants.ColorWhite, step, constants.ColorReset)
	fmt.Fprintf(w, "%s%s%s\n", constants.ColorCyan, cmd, constants.ColorReset)
	fmt.Fprintln(w)
}

// printCloneStepMulti writes a step with multiple command lines.
func printCloneStepMulti(w io.Writer, step string, cmds ...string) {
	fmt.Fprintf(w, "%s%s%s\n", constants.ColorWhite, step, constants.ColorReset)
	for _, cmd := range cmds {
		fmt.Fprintf(w, "%s%s%s\n", constants.ColorCyan, cmd, constants.ColorReset)
	}
	fmt.Fprintln(w)
}
