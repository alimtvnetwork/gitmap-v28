package cmdhealth

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	cmdconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdhealth/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhealth/enums"
)

// hasHelpFlag reports whether help arguments are present.
func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == cmdconstants.FlagArgHelp || a == cmdconstants.FlagArgHelpShort {
			return true
		}
	}

	return false
}

// parseHealthWorkdir extracts and resolves the workdir from flag arguments.
func parseHealthWorkdir(args []string) (string, *appfault.AppError) {
	fs := flag.NewFlagSet(cmdconstants.CmdHealth, flag.ContinueOnError)
	workdirFlag := fs.String(cmdconstants.FlagWorkdir, "", cmdconstants.FlagDescWorkdir)
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(cmdconstants.HelpTextHealth)

		return "", apperror.WrapSimple(parseErr, cmdconstants.WrapHealthFlags)
	}

	return resolveHealthWorkdir(*workdirFlag), nil
}

// printRepoHealth formats and prints git repository status.
func printRepoHealth(facts healthFacts) {
	fmt.Printf(cmdconstants.FmtReposTotal, facts.repoTotal, facts.repoDirty)
	if len(facts.dirtyNames) > 0 {
		fmt.Printf(cmdconstants.FmtDirtyItem, strings.Join(facts.dirtyNames, cmdconstants.ItemSeparator))
	}
	fmt.Println(cmdconstants.ParenClose)
}

// printPipelineHealth formats and prints CI/CD pipeline status.
func printPipelineHealth(facts healthFacts) {
	if facts.pipelineNote != "" {
		fmt.Printf(cmdconstants.FmtPipelinesNote, facts.pipelineNote)

		return
	}
	fmt.Printf(cmdconstants.FmtPipelinesFailing, facts.pipelineRed)
	if len(facts.redNames) > 0 {
		fmt.Printf(cmdconstants.FmtDirtyItem, strings.Join(facts.redNames, cmdconstants.ItemSeparator))
	}
	fmt.Println()
}

// printFleetHealth formats and prints fleet node status.
func printFleetHealth(facts healthFacts) {
	if facts.fleetNote == cmdconstants.NoteNodeRegistryUnreachable {
		fmt.Println(cmdconstants.MsgFleetRegistryUnreachable)

		return
	}

	fmt.Printf(cmdconstants.FmtFleetRegistered, facts.fleetNodes, facts.fleetNote)
}

// printDbHealth formats and prints SQLite database stale count.
func printDbHealth(facts healthFacts) {
	fmt.Printf(cmdconstants.FmtDbsStale, facts.staleDBs)
	if len(facts.staleNames) > 0 {
		fmt.Printf(cmdconstants.FmtDirtyItem, strings.Join(facts.staleNames, cmdconstants.ItemSeparator))
	}
	fmt.Println()
}

// printWorstItem formats and prints the top priority action.
func printWorstItem(facts healthFacts) {
	fmt.Printf(cmdconstants.FmtWorstLine, worstItemPointer(facts))
}

// renderHealthReport outputs the complete health aggregation summary.
func renderHealthReport(facts healthFacts) {
	fmt.Println(cmdconstants.HeaderWorkstationHealth)
	printRepoHealth(facts)
	printPipelineHealth(facts)
	printFleetHealth(facts)
	printDbHealth(facts)
	printWorstItem(facts)
}

// RunHealth implements `gitmap health` (alias: hlth).
func RunHealth(args []string) *appfault.AppError {
	if hasHelpFlag(args) {
		fmt.Print(cmdconstants.HelpTextHealth)

		return nil
	}
	workdir, appErr := parseHealthWorkdir(args)
	if appErr != nil {
		return appErr
	}
	facts := collectHealthFacts(workdir)
	renderHealthReport(facts)

	return nil
}

// worstItemPointer picks the most urgent follow-up command.
func worstItemPointer(facts healthFacts) string {
	if facts.pipelineRed > 0 {
		return fmt.Sprintf(cmdconstants.WorstTriageFmt, facts.redNames[0], facts.pipelineRed)
	}
	if facts.staleDBs > 0 {
		return cmdconstants.WorstDbDoctorRepair
	}
	if facts.repoDirty > 0 {
		return fmt.Sprintf(cmdconstants.WorstDirtyWorkingTreesFmt, facts.repoDirty)
	}

	return cmdconstants.WorstNothingUrgent
}

// DetermineWorstSeverity returns the HealthSeverityType for the aggregated health.
func DetermineWorstSeverity(facts healthFacts) enums.HealthSeverityType {
	if facts.pipelineRed > 0 {
		return enums.HealthSeverityError
	}
	if facts.staleDBs > 0 || facts.repoDirty > 0 {
		return enums.HealthSeverityWarn
	}

	return enums.HealthSeverityOk
}

// DetermineWorstCategory returns the HealthCategoryType for the most urgent issue.
func DetermineWorstCategory(facts healthFacts) enums.HealthCategoryType {
	if facts.pipelineRed > 0 {
		return enums.HealthCategoryPipeline
	}
	if facts.staleDBs > 0 {
		return enums.HealthCategoryDb
	}
	if facts.repoDirty > 0 {
		return enums.HealthCategoryGit
	}

	return enums.HealthCategorySystem
}

// resolveHealthWorkdir mirrors the shared workdir resolution.
func resolveHealthWorkdir(flagVal string) string {
	if strings.TrimSpace(flagVal) != "" {
		return flagVal
	}
	if envVal := strings.TrimSpace(os.Getenv(cmdconstants.EnvGitmapWorkdir)); envVal != "" {
		return envVal
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return cmdconstants.FallbackWorkdir
	}

	return filepath.Join(home, cmdconstants.DefaultWorkspaceDir, cmdconstants.DefaultReposDir)
}
