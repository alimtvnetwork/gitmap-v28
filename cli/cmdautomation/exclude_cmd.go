package cmdautomation

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var excludeCmd = &cobra.Command{
	Use:     "exclude [list|add|rm|clear] [pattern]",
	Aliases: []string{"waiver", "exclusions"},
	Short:   "Manage persistent search and audit exclusion patterns",
	RunE:    runExcludeCmd,
}

func runExcludeCmd(cmd *cobra.Command, args []string) error {
	action := "list"
	if len(args) > 0 {
		action = args[0]
	}

	switch action {
	case "add":
		return handleExcludeAdd(args)
	case "rm", "remove", "delete":
		return handleExcludeRemove(args)
	case "clear", "purge":
		return handleExcludeClear()
	default:
		return handleExcludeList()
	}
}

func handleExcludeList() error {
	entries, err := ListExclusions()
	if err != nil {
		return err
	}
	renderExclusionsList(entries)
	return nil
}

func renderExclusionsList(entries []ExclusionEntry) {
	fmt.Printf("\n%s[Search & Audit Exclusion Registry]%s\n", constants.ColorBold, constants.ColorReset)
	if len(entries) == 0 {
		fmt.Printf("  %sNo custom exclusions configured.%s (Using built-in binary & large-file defaults)\n\n",
			constants.ColorDim, constants.ColorReset)
		return
	}
	fmt.Printf("%-32s | %-16s | %s\n", "Pattern / Path", "Reason", "Added At")
	fmt.Println("---------------------------------|------------------|---------------------")
	for _, e := range entries {
		fmt.Printf("%-32s | %-16s | %s\n", e.Pattern, e.Reason, e.CreatedAt)
	}
	fmt.Println()
}

func handleExcludeAdd(args []string) error {
	if len(args) < 2 {
		return apperror.NewValidationError("pattern is required: gitmap automation exclude add <pattern> [reason]")
	}
	reason := "user_excluded"
	if len(args) > 2 {
		reason = args[2]
	}
	appErr := AddExclusion(args[1], reason)
	if appErr != nil {
		return appErr
	}
	fmt.Printf("\n%s✔ Added exclusion pattern:%s %s (reason: %s)\n\n",
		constants.ColorGreen, constants.ColorReset, args[1], reason)
	return nil
}

func handleExcludeRemove(args []string) error {
	if len(args) < 2 {
		return apperror.NewValidationError("pattern is required: gitmap automation exclude rm <pattern>")
	}
	appErr := RemoveExclusion(args[1])
	if appErr != nil {
		return appErr
	}
	fmt.Printf("\n%s✔ Removed exclusion pattern:%s %s\n\n",
		constants.ColorGreen, constants.ColorReset, args[1])
	return nil
}

func handleExcludeClear() error {
	appErr := ClearExclusions()
	if appErr != nil {
		return appErr
	}
	fmt.Printf("\n%s✔ Cleared all custom exclusions.%s\n\n",
		constants.ColorGreen, constants.ColorReset)
	return nil
}
