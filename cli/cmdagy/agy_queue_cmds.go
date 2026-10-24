package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var agyQueueCmd = &cobra.Command{
	Use:     "queue",
	Aliases: []string{"q"},
	Short:   "Manage Antigravity prompt queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueueStatus()
	},
}

var agyQueueStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current prompt queue status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueueStatus()
	},
}

var agyQueueLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all queued prompts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueueLs()
	},
}

var agyQueueClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all prompts from the queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueueClear()
	},
}

var agyQueuePopCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop and stage the next queued prompt",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyQueuePop()
	},
}

func initAgyQueueCommands() {
	agyQueueCmd.AddCommand(agyQueueStatusCmd)
	agyQueueCmd.AddCommand(agyQueueLsCmd)
	agyQueueCmd.AddCommand(agyQueueClearCmd)
	agyQueueCmd.AddCommand(agyQueuePopCmd)
	AgyCmd.AddCommand(agyQueueCmd)
}

func initAgyPromptAndStatusCommands() {
	initAgyPromptSubcommands()
	agyStatusCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runAgyStatusWithQueue()
	}
}

func runAgyQueueStatus() error {
	q, err := GetQueueStatus()
	if err != nil {
		return apperror.WrapSimple(err, "get queue status")
	}

	fmt.Printf("%s● Antigravity Prompt Queue Status%s\n", constants.ColorCyan, constants.ColorReset)
	renderQueueActiveItem(q.Active)
	fmt.Printf("  Queued Prompts: %d\n", len(q.Queued))
	if len(q.UpdatedAt) > 0 {
		fmt.Printf("  Last Updated:   %s\n", q.UpdatedAt)
	}

	return nil
}

func renderQueueActiveItem(item *AgyPromptQueueEntry) {
	if item == nil {
		fmt.Println("  Active Prompt:  None")

		return
	}

	fmt.Printf("  Active Prompt:  #%d [%s] %s (%s)\n", item.ID, item.Type, item.Title, item.Status)
}

func runAgyQueueLs() error {
	q, err := GetQueueStatus()
	if err != nil {
		return apperror.WrapSimple(err, "list queue")
	}
	if len(q.Queued) == 0 {
		fmt.Println("No prompts currently queued.")

		return nil
	}

	fmt.Printf("%sQueued Prompts (%d):%s\n", constants.ColorCyan, len(q.Queued), constants.ColorReset)
	for _, item := range q.Queued {
		fmt.Printf("  #%d [%s] %s (Created: %s)\n", item.ID, item.Type, item.Title, item.CreatedAt)
	}

	return nil
}

func runAgyQueueClear() error {
	if err := ClearPromptQueue(); err != nil {
		return apperror.WrapSimple(err, "clear prompt queue")
	}
	fmt.Printf("%s✔ Cleared Antigravity prompt queue%s\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func runAgyQueuePop() error {
	popped, err := PopNextQueuedPrompt()
	if err != nil {
		return apperror.WrapSimple(err, "pop queued prompt")
	}
	if popped == nil {
		fmt.Println("Queue is empty. No prompts to pop.")

		return nil
	}
	fmt.Printf("%s✔ Popped prompt #%d [%s] %s%s\n", constants.ColorGreen, popped.ID, popped.Type, popped.Title, constants.ColorReset)
	fmt.Printf("  Staged to: %s\n", resolveActiveAgyPromptPath())
	fmt.Println("  Copied prompt content to clipboard.")

	return nil
}

func runAgyStatusWithQueue() error {
	cwd, _ := os.Getwd()
	ideProc := DetectRunningAntigravityIDE()
	convState := DetectConversationExecutionState(cwd)
	q, _ := GetQueueStatus()
	activeConvs, _ := FetchActiveRunningConversations()
	wsQueues, _ := DiscoverAllWorkspaceQueues()
	totalQueued := countTotalQueuedPrompts(wsQueues)

	renderStatusHeader(ideProc, convState, len(q.Queued), len(activeConvs), totalQueued)

	return runAgyLs(5)
}

func renderStatusHeader(ideProc AgyProcessInfoWrap, state AgyConversationExecutionState, localQ, activeCount, totalQ int) {
	fmt.Printf("%s● Antigravity System Status%s\n", constants.ColorCyan, constants.ColorReset)
	renderIDEProcessStatus(ideProc)
	renderConversationStatus(state)
	renderActiveCountStatus(activeCount)
	renderQueueCountStatus(localQ, totalQ)
	fmt.Println()
}

func renderActiveCountStatus(activeCount int) {
	if activeCount > 0 {
		fmt.Printf("  Active Prompts: %s%d RUNNING%s ('gitmap agy active')\n",
			constants.ColorYellow, activeCount, constants.ColorReset)
		return
	}
	fmt.Printf("  Active Prompts: %s0 running (idle)%s\n", constants.ColorGreen, constants.ColorReset)
}

func renderQueueCountStatus(localQ, totalQ int) {
	if totalQ > 0 {
		fmt.Printf("  Prompt Queues:  %s%d pending across workspaces%s (local: %d, 'gitmap agy queues')\n",
			constants.ColorCyan, totalQ, constants.ColorReset, localQ)
		return
	}
	fmt.Printf("  Prompt Queues:  0 pending\n")
}

func renderIDEProcessStatus(ideProc AgyProcessInfoWrap) {
	if ideProc.IsSuccess() {
		fmt.Printf("  IDE Process:   %sRUNNING%s (PID: %d)\n", constants.ColorGreen, constants.ColorReset, ideProc.Value.PID)

		return
	}
	fmt.Printf("  IDE Process:   %sNOT RUNNING%s\n", constants.ColorYellow, constants.ColorReset)
}

func renderConversationStatus(state AgyConversationExecutionState) {
	color := constants.ColorGreen
	if state.Status == enums.AgyConvStatusRunning {
		color = constants.ColorYellow
	}
	if len(state.ConvID) > 0 {
		fmt.Printf("  Conversation:  %s%s%s (ID: %s)\n", color, strings.ToUpper(state.Status.String()), constants.ColorReset, state.ConvID)

		return
	}
	fmt.Printf("  Conversation:  %s%s%s\n", color, strings.ToUpper(state.Status.String()), constants.ColorReset)
}

func getProjectsDirPath() (string, error) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", homeErr
	}

	return filepath.Join(homeDir, ".gemini", "config", "projects"), nil
}

func ensureDirExistsErr(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}
