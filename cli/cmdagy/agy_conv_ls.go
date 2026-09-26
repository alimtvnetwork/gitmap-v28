package cmdagy

import (
	"fmt"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

var (
	agyConvLsFile string
)

var agyConvCmd = &cobra.Command{
	Use:     "conv",
	Aliases: []string{"conversations", "con"},
	Short:   "Manage and inspect conversations",
}

var agyConvLsCmd = &cobra.Command{
	Use:   "ls [N]",
	Short: "List recent conversations with status and queued counts",
	RunE: func(cmd *cobra.Command, args []string) error {
		n := 8
		n = parseAgyConvLsArgCount(args, n)
		return runAgyConvLs(n)
	},
}

func parseAgyConvLsArgCount(args []string, defaultN int) int {
	if len(args) == 0 {
		return defaultN
	}
	val, err := strconv.Atoi(args[0])
	if err == nil && val > 0 {
		return val
	}
	return defaultN
}

func init() {
	agyConvLsCmd.Flags().StringVarP(&agyConvLsFile, "file", "f", "", "Export JSON to file (default repo-conversations.json)")
	if flag := agyConvLsCmd.Flags().Lookup("file"); flag != nil {
		flag.NoOptDefVal = "repo-conversations.json"
	}
	agyConvCmd.AddCommand(agyConvLsCmd)
	AgyCmd.AddCommand(agyConvCmd)
}

func runAgyConvLs(n int) error {
	projectFilter := ""
	currentRepo, isInsideRepo := detectCurrentProjectContext()
	if isInsideRepo {
		projectFilter = currentRepo
	}

	rows, isDepthMode, err := FetchConversationInspectRows(n, projectFilter)
	if err != nil {
		return apperror.WrapSimple(err, "fetch conversation inspect rows")
	}

	CacheInspectRowSequences(rows, "conv")

	if agyConvLsFile != "" {
		return ExportInspectRowsToJSONFile(rows, agyConvLsFile)
	}

	if isDepthMode {
		fmt.Printf("\n  \033[36m[Depth Mode]\033[0m Showing last %d conversations for repo: %s\n", len(rows), currentRepo)
	}

	RenderInspectRowsTable(rows, isDepthMode)
	return nil
}
