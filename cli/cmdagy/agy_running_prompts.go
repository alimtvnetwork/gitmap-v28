package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/spf13/cobra"
)

var (
	runningPromptsInstanceFlag    string
	runningPromptsAllInstanceFlag bool
)

func init() {
	setupMultiInstanceRunningPromptsFlags()
	setupMultiInstanceRunningProjectsHook()
}

func setupMultiInstanceRunningPromptsFlags() {
	for _, cmd := range RunningPromptsCmd.Commands() {
		if cmd.Name() == "ls" {
			cmd.Flags().StringVarP(&runningPromptsInstanceFlag, "instance", "i", "", "Filter by Antigravity instance ID or alias (default: primary)")
			cmd.Flags().BoolVarP(&runningPromptsAllInstanceFlag, "all-instances", "a", false, "Aggregate prompts across all active Antigravity instances")

			origRunE := cmd.RunE
			cmd.RunE = func(c *cobra.Command, args []string) error {
				instFlag, _ := c.Flags().GetString("instance")
				allFlag, _ := c.Flags().GetBool("all-instances")
				limit, _ := c.Flags().GetInt("limit")
				wc := resolveWordCountFlag(c)
				isFull, _ := c.Flags().GetBool("full")
				isJSON, _ := c.Flags().GetBool("json")

				if instFlag != "" || allFlag {
					return RunMultiInstancePromptsLs(instFlag, allFlag, limit, wc, isFull, isJSON)
				}
				if origRunE != nil {
					return origRunE(c, args)
				}
				return nil
			}
		}
	}
}

func setupMultiInstanceRunningProjectsHook() {
	AgyRunningProjectsCmd.Flags().StringVarP(&runningPromptsInstanceFlag, "instance", "i", "", "Filter by Antigravity instance ID or alias")
	AgyRunningProjectsCmd.Flags().BoolVarP(&runningPromptsAllInstanceFlag, "all-instances", "a", false, "Aggregate projects across all active Antigravity instances")

	origRunE := AgyRunningProjectsCmd.RunE
	AgyRunningProjectsCmd.RunE = func(cmd *cobra.Command, args []string) error {
		instID, isAll := extractInstanceArgs(args)
		if instID != "" || isAll {
			isJSON := hasArgFlag(args, "--json") || runningProjectsJSON
			return RunMultiInstanceProjectsCLI(instID, isAll, isJSON, args)
		}
		if origRunE != nil {
			return origRunE(cmd, args)
		}
		return nil
	}
}

func extractInstanceArgs(args []string) (string, bool) {
	instID := ""
	isAll := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-a" || arg == "--all-instances" {
			isAll = true
			continue
		}
		if arg == "-i" || arg == "--instance" {
			if i+1 < len(args) {
				instID = args[i+1]
				i++
			}
			continue
		}
		if strings.HasPrefix(arg, "--instance=") {
			instID = strings.TrimPrefix(arg, "--instance=")
			continue
		}
	}
	return instID, isAll
}

// RunMultiInstancePromptsLs retrieves and formats prompts across instances.
func RunMultiInstancePromptsLs(instanceID string, isAll bool, limit, wordCount int, isFull, isJSON bool) error {
	opts := AgyInstancePromptQueryOptions{
		InstanceID:   instanceID,
		IsAll:        isAll,
		Limit:        limit,
		MaxWords:     wordCount,
		IncludeConvs: true,
	}

	resp, err := QueryInstancePrompts(opts)
	if err != nil {
		return err
	}

	if isJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(resp)
	}

	renderMultiInstancePromptsView(resp, isFull)
	return nil
}

func renderMultiInstancePromptsView(resp *AgyMultiInstancePromptResponse, isFull bool) {
	if resp == nil || len(resp.Instances) == 0 {
		fmt.Printf("%sNo Antigravity instances found.%s\n", constants.ColorYellow, constants.ColorReset)
		return
	}

	fmt.Printf("\n%s⚡ Antigravity Multi-Instance Prompt Inspection%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("Total Instances: %d | Total Prompts: %d\n", resp.TotalInstances, resp.TotalPrompts)
	fmt.Println(strings.Repeat("─", 70))

	for _, inst := range resp.Instances {
		fmt.Printf("\n%s● Instance: %s [%s]%s\n", constants.ColorBold, inst.InstanceName, inst.InstanceID, constants.ColorReset)
		fmt.Printf("  Running Prompts: %d | Queued Prompts: %d\n", inst.RunningCount, inst.QueuedCount)

		if len(inst.Running) == 0 && len(inst.Queued) == 0 {
			fmt.Printf("  %s(No active or queued prompts)%s\n", constants.ColorDim, constants.ColorReset)
			continue
		}

		for idx, item := range inst.Running {
			fmt.Printf("  %s[%d] [RUNNING] %s (%s)%s\n", constants.ColorGreen, idx+1, item.Title, item.ProjectName, constants.ColorReset)
			if item.PromptPreview != "" {
				preview := item.PromptPreview
				if !isFull && len(preview) > 120 {
					preview = preview[:120] + "..."
				}
				fmt.Printf("      Preview: %s\n", preview)
			}
		}

		for idx, item := range inst.Queued {
			fmt.Printf("  %s[%d] [QUEUED] %s%s\n", constants.ColorYellow, idx+1, item.Title, constants.ColorReset)
			if item.Prompt != "" {
				preview := item.Prompt
				if !isFull && len(preview) > 120 {
					preview = preview[:120] + "..."
				}
				fmt.Printf("      Prompt: %s\n", preview)
			}
		}
	}
	fmt.Println()
}

// RunMultiInstanceProjectsCLI filters running projects by target instance or aggregates across instances.
func RunMultiInstanceProjectsCLI(instanceID string, isAll bool, isJSON bool, args []string) error {
	allProjects, err := DiscoverRunningProjects()
	if err != nil {
		return err
	}

	instances, iErr := DiscoverAllAgyInstances()
	if iErr != nil {
		return iErr
	}

	allowedWorkspaces := make(map[string]bool)
	targetID := strings.TrimSpace(instanceID)

	if isAll || strings.EqualFold(targetID, "all") || targetID == "" {
		for _, inst := range instances {
			for _, ws := range inst.ActiveWorkspaces {
				allowedWorkspaces[strings.ToLower(ws)] = true
			}
		}
	} else {
		resolved, rErr := ResolveInstance(targetID)
		if rErr != nil {
			return rErr
		}
		for _, ws := range resolved.ActiveWorkspaces {
			allowedWorkspaces[strings.ToLower(ws)] = true
		}
	}

	var filtered []RunningProjectRecord
	for _, p := range allProjects {
		pClean := strings.ToLower(filepath.Clean(p.ProjectPath))
		if len(allowedWorkspaces) == 0 || allowedWorkspaces[pClean] {
			filtered = append(filtered, p)
		}
	}

	if isJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(filtered)
	}

	return outputRunningProjects(filtered, false, "")
}
