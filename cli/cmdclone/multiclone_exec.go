package cmdclone

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

// BuildMultiCloneItems converts raw URLs to destination-mapped MultiCloneItem entries.
func BuildMultiCloneItems(urls []string, baseDir string) []MultiCloneItem {
	items := make([]MultiCloneItem, 0, len(urls))
	for _, url := range urls {
		repo := repoNameFromURL(url)
		destFolder := resolveCloneDestination(baseDir, repo)
		absPath, _ := filepath.Abs(destFolder)

		items = append(items, MultiCloneItem{
			URL:           url,
			RepoName:      repo,
			TargetDir:     destFolder,
			TargetAbsPath: absPath,
		})
	}
	return items
}

func resolveCloneDestination(baseDir, repoName string) string {
	if baseDir != "" {
		return filepath.Join(baseDir, repoName)
	}
	return repoName
}

// RenderMultiCloneDryRun outputs the preview table of repositories to clone.
func RenderMultiCloneDryRun(items []MultiCloneItem, targetDir string) {
	fmt.Printf("\n%s[GitMap MultiClone — Dry Run Preview]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Target Directory: %s\n", resolveTargetDisplay(targetDir))
	fmt.Printf("  Repositories to clone: %d\n\n", len(items))

	for i, item := range items {
		fmt.Printf("  %s[%d/%d]%s %-30s -> %s\n",
			constants.ColorCyan, i+1, len(items), constants.ColorReset,
			item.RepoName, item.TargetDir)
		fmt.Printf("         URL: %s\n", item.URL)
	}
	fmt.Printf("\n%sDry-run mode active. No repositories were cloned.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func resolveTargetDisplay(dir string) string {
	if dir == "" {
		return "Current Directory (.)"
	}
	return dir
}

// ExecuteMultiCloneBatch clones items sequentially, reporting progress and updating VS Code.
func ExecuteMultiCloneBatch(items []MultiCloneItem, opts MultiCloneOptions) MultiCloneSummary {
	fmt.Printf(constants.MsgCloneMultiBegin, len(items))
	summary := MultiCloneSummary{UniqueCount: len(items)}
	pmPairs := make([]vscodepm.Pair, 0, len(items))

	for i, item := range items {
		printCloneItemProgress(i+1, len(items), item)
		err := executeDirectCloneOne(item.URL, item.TargetDir, opts.GHDesktop, opts.NoReplace)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCloneMultiFailedFmt, i+1, len(items), item.URL, err)
			summary.Failed++
			continue
		}

		summary.Succeeded++
		pmPairs = append(pmPairs, buildClonePMPair(item.TargetAbsPath, item.RepoName))
	}

	syncClonedReposToVSCodePM(pmPairs, opts.NoVSCodeSync)
	printMultiCloneSummary(summary)
	return summary
}

func printCloneItemProgress(idx, total int, item MultiCloneItem) {
	fmt.Printf(constants.MsgCloneMultiItem, idx, total, item.URL)
	if item.TargetDir != item.RepoName {
		fmt.Printf("       ↳ Target destination: %s\n", item.TargetDir)
	}
}

func printMultiCloneSummary(summary MultiCloneSummary) {
	fmt.Printf(constants.MsgCloneSummaryMultiFmt, summary.Succeeded, summary.Failed, summary.UniqueCount)
}
