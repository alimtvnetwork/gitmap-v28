package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// listSUGProjects renders the comprehensive status, target list, rules, and commands.
func listSUGProjects() error {
	cfg := loadSUGConfig()
	st, isRunning := LoadSUGRuntimeStatus()

	printSUGHeader(len(cfg.ProjectTargets))
	printSUGWatchStatus(st, isRunning)
	printSUGTargetList(cfg.ProjectTargets)
	printSUGShutdownCondition(len(cfg.ProjectTargets))
	printSUGManagementExamples()
	return nil
}

func printSUGHeader(count int) {
	fmt.Println()
	fmt.Printf("  %s╔════ SHUTDOWN-UNTIL-GREEN WATCH LIST (%d projects) ════╗%s\n\n",
		constants.ColorCyan, count, constants.ColorReset)
}

func printSUGWatchStatus(st SUGRuntimeState, isRunning bool) {
	if isRunning {
		fmt.Printf("  %s● Watch Loop State:%s ACTIVE & RUNNING (PID: %d, interval: %ds, started: %s)\n",
			constants.ColorGreen, constants.ColorReset, st.PID, st.IntervalSec, st.StartedAt)
		fmt.Printf("  %s● Background Execution:%s ACTIVE (daemon process is actively monitoring pipelines)\n\n",
			constants.ColorGreen, constants.ColorReset)
		return
	}
	fmt.Printf("  %s○ Watch Loop State:%s IDLE (Not currently watching)\n",
		constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  %sℹ Background Execution:%s NO (Does NOT run automatically in background unless 'gitmap sug watch' or 'gitmap sug run' is started)\n\n",
		constants.ColorDim, constants.ColorReset)
}

func printSUGTargetList(targets []string) {
	if len(targets) == 0 {
		fmt.Printf("  %sMonitored Projects: (Watch list is empty. Add projects with 'add-projects <target>')%s\n\n",
			constants.ColorDim, constants.ColorReset)
		return
	}
	fmt.Println("  Monitored Projects:")
	for i, p := range targets {
		fmt.Printf("    [%d] %s\n", i+1, p)
	}
	fmt.Println()
}

func printSUGShutdownCondition(count int) {
	fmt.Println("  Shutdown Condition:")
	if count <= 1 {
		fmt.Println("    • OS shutdown will trigger ONLY when all monitored projects turn 100% green.")
		fmt.Println()
		return
	}
	fmt.Printf("    • OS shutdown will trigger ONLY when ALL %d projects are green (100%% passing) together at the same time.\n", count)
	fmt.Println("    • It will NEVER shut down if only 1 or a partial subset passes.")
	fmt.Println()
}

func printSUGManagementExamples() {
	fmt.Println("  Management Commands:")
	fmt.Println("    • Remove single target : gitmap sug remove <target>   (or 'gitmap sug rm <target>')")
	fmt.Println("    • Clear all targets    : gitmap sug clear             (or 'gitmap sug remove-all', 'gitmap sug rm-all')")
	fmt.Println("    • Add project target   : gitmap sug add <target>      (or 'gitmap sug add-projects <target>')")
	fmt.Println("    • Auto-add running AGY : gitmap sug agy-running-projects (or 'gitmap sug arp')")
	fmt.Println("    • Start watching loop  : gitmap sug watch             (or 'gitmap sug run', 'gitmap sug ui')")
	fmt.Println()
}
