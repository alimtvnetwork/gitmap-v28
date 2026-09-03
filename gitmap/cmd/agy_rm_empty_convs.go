package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var (
	agyRmEmptyExcept string
	agyRmEmptyDryRun bool
	agyRmEmptyYes    bool
)

var agyRemoveEmptyConvsCmd = &cobra.Command{
	Use:     "remove-projects-with-empty-conversations",
	Aliases: []string{"rm-empty-conversations", "clean-empty-conversations", "prune-empty-conversations"},
	Short:   "Remove Antigravity projects that have no active conversations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgyRemoveEmptyConvs(args)
	},
}

func init() {
	agyRemoveEmptyConvsCmd.Flags().StringVarP(&agyRmEmptyExcept, "except", "e", "", "Exclude projects matching id, name, path, or from a CSV/text file")
	agyRemoveEmptyConvsCmd.Flags().BoolVarP(&agyRmEmptyDryRun, "dry-run", "d", false, "Preview projects to be removed without deleting")
	agyRemoveEmptyConvsCmd.Flags().BoolVarP(&agyRmEmptyYes, "yes", "y", false, "Confirm removal without prompt")
}

func runAgyRemoveEmptyConvs(args []string) error {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return apperror.WrapSimple(err, "path error")
	}
	projects, err := loadAllAgyProjects(dirPath)
	if err != nil {
		return apperror.WrapSimple(err, "load projects")
	}
	convs, err := scanAllConversations()
	if err != nil {
		return apperror.WrapSimple(err, "scan conversations")
	}
	return processAgyEmptyRemoval(dirPath, projects, convs)
}

func processAgyEmptyRemoval(dirPath string, projects []AgyProject, convs []AgyConvInfo) error {
	mapped := mapProjectsToConversations(projects, convs)
	candidates := filterEmptyProjectConvs(mapped)
	exceptTokens := parseAgyExceptTokens(agyRmEmptyExcept)
	targets := filterExceptedCandidates(candidates, exceptTokens)

	if len(targets) == 0 {
		fmt.Printf("%s No Antigravity projects with empty conversations found to remove.\n",
			constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}
	return executeAgyEmptyRemoval(dirPath, targets, len(projects)-len(targets))
}

func filterExceptedCandidates(candidates []AgyProjectConvs, tokens []string) []AgyProjectConvs {
	var targets []AgyProjectConvs
	for _, c := range candidates {
		if !isAgyProjectExceptedWithTokens(c.Project, tokens) {
			targets = append(targets, c)
		}
	}
	return targets
}

func executeAgyEmptyRemoval(dirPath string, targets []AgyProjectConvs, remaining int) error {
	printEmptyRemovalPreview(targets)
	if agyRmEmptyDryRun {
		fmt.Printf("\n%s [dry-run] %d project(s) would be removed. Remaining: %d\n",
			constants.ColorYellow+"ℹ"+constants.ColorReset, len(targets), remaining)
		return nil
	}
	if !agyRmEmptyYes && !confirmEmptyRemoval(len(targets)) {
		fmt.Println("Removal canceled. No changes made.")
		return nil
	}
	return deleteTargetProjectFiles(dirPath, targets, remaining)
}

func confirmEmptyRemoval(count int) bool {
	msg := fmt.Sprintf("%sAre you sure you want to remove %d Antigravity project(s) with empty conversations? [y/N]: %s",
		constants.ColorYellow, count, constants.ColorReset)
	ok, err := promptConfirm(msg)
	return err == nil && ok
}

func deleteTargetProjectFiles(dirPath string, targets []AgyProjectConvs, remaining int) error {
	deleted := 0
	for _, t := range targets {
		p := filepath.Join(dirPath, t.Project.ID+".json")
		if err := os.Remove(p); err == nil {
			deleted++
		}
	}
	fmt.Printf("\n%s Successfully removed %d Antigravity project(s) with empty conversations. Remaining active: %d\n",
		constants.ColorGreen+"✓"+constants.ColorReset, deleted, remaining)
	return nil
}

func printEmptyRemovalPreview(targets []AgyProjectConvs) {
	fmt.Printf("\n  %sTargeting %d project(s) with empty conversations for removal:%s\n\n",
		constants.ColorYellow, len(targets), constants.ColorReset)
	for _, t := range targets {
		fmt.Printf("    • %-36s %-22s %s\n", t.Project.ID, truncateStr(t.Project.Name, 21), truncateStr(t.Project.GetPath(), 35))
	}
}
