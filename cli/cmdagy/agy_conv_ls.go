package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"github.com/spf13/cobra"
)

var (
	agyConvLsFile string
)

var agyConvCmd = &cobra.Command{
	Use:     "conv",
	Aliases: []string{"conversations", "con"},
	Short:   "Manage conversations",
}

var agyConvLsCmd = &cobra.Command{
	Use:   "ls [N]",
	Short: "List recent conversations",
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
	agyConvLsCmd.Flags().StringVarP(&agyConvLsFile, "file", "f", "", "Export JSON to file")
	agyConvCmd.AddCommand(agyConvLsCmd)
	AgyCmd.AddCommand(agyConvCmd)
}

func runAgyConvLs(n int) error {
	queues, _ := DiscoverAllWorkspaceQueues()
	
	if agyConvLsFile != "" {
		data, _ := json.MarshalIndent(queues, "", "  ")
		return os.WriteFile(agyConvLsFile, data, 0644)
	}
	
	fmt.Printf("Listing top %d conversations...\n", n)
	return nil
}
