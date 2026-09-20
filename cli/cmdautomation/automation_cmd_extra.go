package cmdautomation

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	newlinesCmd = &cobra.Command{
		Use:     "newlines [paths...]",
		Aliases: []string{"fix-newlines", "lf"},
		Short:   "Polyglot CRLF to LF and trailing whitespace normalizer",
		RunE:    runNewlinesCmd,
	}

	cacheCmd = &cobra.Command{
		Use:     "cache [status|read|warm|clear]",
		Aliases: []string{"mem-cache"},
		Short:   "Manage fast in-memory file and scan cache",
		RunE:    runCacheCmd,
	}

	locateCmd = &cobra.Command{
		Use:     "locate [tool-name]",
		Aliases: []string{"find-tool", "vcvars", "find-exe"},
		Short:   "Fast developer tool locator with vswhere and vcvarsall.bat fast paths",
		RunE:    runLocateCmd,
	}

	newlineOpts NewlineOptions
	locateOpts  LocateOptions
)

func runNewlinesCmd(cmd *cobra.Command, args []string) error {
	newlineOpts.Paths = args
	res, err := RunNormalizeNewlines(newlineOpts)
	if err != nil {
		return err
	}
	renderNewlineResult(res)
	return nil
}

func renderNewlineResult(res NewlineResult) {
	fmt.Printf("\n%s[Polyglot Newline Normalization]%s\n", constants.ColorBold, constants.ColorReset)
	fmt.Printf("  Scanned Files:   %d\n", res.ScannedFiles)
	fmt.Printf("  Modified Files:  %d\n", res.ModifiedFiles)
	fmt.Printf("  CRLF Converted:  %d\n", res.CrlfCount)
	fmt.Printf("  Duration:        %s\n\n", res.Duration)
}

func runCacheCmd(cmd *cobra.Command, args []string) error {
	action := "status"
	if len(args) > 0 {
		action = args[0]
	}
	switch action {
	case "read":
		path := "."
		if len(args) > 1 {
			path = args[1]
		}
		return toError(RunCacheRead(path))
	case "warm":
		dir := "."
		if len(args) > 1 {
			dir = args[1]
		}
		return toError(RunCacheWarm(dir))
	case "clear", "purge", "rm":
		return toError(RunCacheClear())
	default:
		return toError(RunCacheStatus())
	}
}

func renderSearchResults(res SearchResult) {
	fmt.Printf("\n%sSearch completed:%s %d hit(s) across %d file(s) in %s\n",
		constants.ColorGreen, constants.ColorReset, res.TotalHits, res.TotalFiles, res.Duration)
	for _, m := range res.Matches {
		fmt.Printf("  %s%s:%d%s %s\n", constants.ColorCyan, m.Path, m.LineNumber, constants.ColorReset, m.LineContent)
	}
	fmt.Println()
}

func runLocateCmd(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		locateOpts.Target = args[0]
	}
	res, err := RunLocate(locateOpts)
	renderLocateResult(res, locateOpts.IsJSON)
	if err != nil {
		return err
	}
	return nil
}

func initSubsystems() {
	AutomationCmd.AddCommand(newlinesCmd)
	AutomationCmd.AddCommand(cacheCmd)
	AutomationCmd.AddCommand(locateCmd)
	AutomationCmd.AddCommand(guardCmd)
	AutomationCmd.AddCommand(sequenceCmd)
	AutomationCmd.AddCommand(excludeCmd)
	AutomationCmd.AddCommand(relPathsCmd)
	AutomationCmd.AddCommand(namingCmd)
	AutomationCmd.AddCommand(resultWrapperCmd)
	AutomationCmd.AddCommand(paramsCmd)
	AutomationCmd.AddCommand(enumsCmd)

	newlinesCmd.Flags().BoolVarP(&newlineOpts.IsFixMode, "fix", "f", false, "Write normalized changes to disk")
	newlinesCmd.Flags().BoolVar(&newlineOpts.IsDryRun, "dry-run", false, "Preview modifications without writing")
	locateCmd.Flags().BoolVar(&locateOpts.IsJSON, "json", false, "Output results as JSON")
	locateCmd.Flags().BoolVar(&locateOpts.IsCmd, "cmd", false, "Output shell initialization command")
	initGuardFlags()
	initSequenceFlags()
	initRelPathsFlags()
	initNamingFlags()
	initRuleFlags()
}
