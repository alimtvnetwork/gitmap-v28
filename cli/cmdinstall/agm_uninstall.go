package cmdinstall

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	isAgmUninstallAll    bool
	isAgmUninstallDryRun bool
	isAgmUninstallForce  bool
)

var agmTargetProcRegex = regexp.MustCompile(`(?i)^(antigravity-manager|agm)(\.exe)?$`)

// AgmUninstallCmd uninstalls Antigravity Manager.
var AgmUninstallCmd = &cobra.Command{
	Use:     "uninstall [flags]",
	Aliases: []string{"remove", "rm", "purge"},
	Short:   "Uninstall Antigravity Manager (AGM) desktop application and tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAGMUninstall(isAgmUninstallAll, isAgmUninstallForce, isAgmUninstallDryRun)
	},
}

// AgmUninstallAllCmd performs a full purge of Antigravity Manager.
var AgmUninstallAllCmd = &cobra.Command{
	Use:   "uninstall-all [flags]",
	Short: "Full purge of Antigravity Manager (AGM), configs, shortcuts, and registry entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAGMUninstall(true, isAgmUninstallForce, isAgmUninstallDryRun)
	},
}

func initAgmUninstallCmd() {
	bindAgmUninstallFlags()
	AgmCmd.AddCommand(AgmUninstallCmd)
	AgmCmd.AddCommand(AgmUninstallAllCmd)
}

func bindAgmUninstallFlags() {
	AgmUninstallCmd.Flags().BoolVarP(&isAgmUninstallAll, "all", "a", false, "Full purge including configs, shortcuts, and registry")
	AgmUninstallCmd.Flags().BoolVarP(&isAgmUninstallDryRun, "dry-run", "n", false, "Simulate uninstallation without deleting")
	AgmUninstallCmd.Flags().BoolVarP(&isAgmUninstallForce, "force", "f", false, "Force uninstallation without confirmation")
	AgmUninstallAllCmd.Flags().BoolVarP(&isAgmUninstallDryRun, "dry-run", "n", false, "Simulate uninstallation without deleting")
	AgmUninstallAllCmd.Flags().BoolVarP(&isAgmUninstallForce, "force", "f", false, "Force uninstallation without confirmation")
}

// RunAGMUninstall terminates processes, removes binaries, and purges configurations.
func RunAGMUninstall(isFullPurge, isForce, isDryRun bool) error {
	printAgmUninstallHeader(isFullPurge)
	candidates := collectAgmCandidatePaths(isFullPurge)
	validErr := validateSafeAgmCandidates(candidates)
	hasValidErr := validErr != nil
	if hasValidErr {
		return validErr
	}
	return executeAgmUninstallStages(candidates, isFullPurge, isForce, isDryRun)
}

func executeAgmUninstallStages(candidates []string, isFullPurge, isForce, isDryRun bool) error {
	terminateRunningAgmProcesses(isDryRun)
	isConfirmed := checkAgmConfirmation(candidates, isForce, isDryRun)
	if !isConfirmed {
		fmt.Println("Uninstallation cancelled by user.")
		return nil
	}
	return applyAgmRemoval(candidates, isFullPurge, isDryRun)
}

func printAgmUninstallHeader(isFullPurge bool) {
	if isFullPurge {
		fmt.Printf("%s[GitMap]%s Starting Antigravity Manager (AGM) FULL PURGE...\n", constants.ColorCyan, constants.ColorReset)
		return
	}
	fmt.Printf("%s[GitMap]%s Starting Antigravity Manager (AGM) Standard Uninstallation...\n", constants.ColorCyan, constants.ColorReset)
}

func validateSafeAgmCandidates(candidates []string) error {
	for _, p := range candidates {
		isSafe := isPathSafeToDelete(p)
		if !isSafe {
			return apperror.NewSimple(fmt.Sprintf("safety violation: protected path %s cannot be deleted", p), "E9001")
		}
	}
	return nil
}

func isPathSafeToDelete(path string) bool {
	clean := filepath.Clean(path)
	hasEmpty := clean == "" || clean == "." || clean == ".."
	if hasEmpty {
		return false
	}
	isRoot := filepath.VolumeName(clean) == clean || clean == "/" || clean == "\\"
	if isRoot {
		return false
	}
	isWorkOverlap := isWorkDirectoryOverlap(clean)
	if isWorkOverlap {
		return false
	}
	isGitRepo := isGitRepositoryPath(clean)
	if isGitRepo {
		return false
	}
	return true
}

func isWorkDirectoryOverlap(path string) bool {
	norm := strings.ToLower(filepath.ToSlash(path))
	hasWork := strings.HasPrefix(norm, "d:/work") || norm == "d:/work"
	if hasWork {
		return true
	}
	isWorkParent := strings.HasPrefix("d:/work", norm+"/")
	if isWorkParent {
		return true
	}
	return false
}

func isGitRepositoryPath(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	hasGit := err == nil && info.IsDir()
	if hasGit {
		return true
	}
	hasGitInPath := strings.Contains(strings.ToLower(path), ".git")
	if hasGitInPath {
		return true
	}
	return false
}

func terminateRunningAgmProcesses(isDryRun bool) {
	pids := findRunningAgmProcessPIDs()
	hasPids := len(pids) > 0
	if !hasPids {
		return
	}
	if isDryRun {
		fmt.Printf("  [dry-run] Would terminate %d running AGM process(es)\n", len(pids))
		return
	}
	killRunningAgmPIDs(pids)
}

func findRunningAgmProcessPIDs() []int {
	hasWindows := runtime.GOOS == "windows"
	if hasWindows {
		return findAgmPIDsWindows()
	}
	return findAgmPIDsUnix()
}

func findAgmPIDsWindows() []int {
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	setHiddenProcessAttr(cmd)
	out, err := cmd.Output()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	return parseAgmTasklistCSV(string(out))
}

func parseAgmTasklistCSV(output string) []int {
	var pids []int
	currentPid := os.Getpid()
	reader := csv.NewReader(strings.NewReader(output))
	records, err := reader.ReadAll()
	hasErr := err != nil
	if hasErr {
		return pids
	}
	for _, rec := range records {
		pid, isMatch := matchAgmCSVRecord(rec, currentPid)
		if isMatch {
			pids = append(pids, pid)
		}
	}
	return pids
}

func matchAgmCSVRecord(rec []string, currentPid int) (int, bool) {
	hasValidLen := len(rec) >= 2
	if !hasValidLen {
		return 0, false
	}
	name := strings.TrimSpace(rec[0])
	pid, convErr := strconv.Atoi(strings.TrimSpace(rec[1]))
	hasConvErr := convErr != nil || pid == currentPid
	if hasConvErr {
		return 0, false
	}
	isTarget := agmTargetProcRegex.MatchString(name)
	return pid, isTarget
}

func findAgmPIDsUnix() []int {
	cmd := exec.Command("ps", "-eo", "pid,comm")
	out, err := cmd.Output()
	hasErr := err != nil
	if hasErr {
		return nil
	}
	return parseAgmPsOutput(string(out))
}

func parseAgmPsOutput(output string) []int {
	var pids []int
	currentPid := os.Getpid()
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		pid, isMatch := matchAgmPsLine(line, currentPid)
		if isMatch {
			pids = append(pids, pid)
		}
	}
	return pids
}

func matchAgmPsLine(line string, currentPid int) (int, bool) {
	fields := strings.Fields(line)
	hasFields := len(fields) >= 2
	if !hasFields {
		return 0, false
	}
	pid, convErr := strconv.Atoi(fields[0])
	hasConvErr := convErr != nil || pid == currentPid
	if hasConvErr {
		return 0, false
	}
	name := filepath.Base(fields[1])
	isTarget := agmTargetProcRegex.MatchString(name)
	return pid, isTarget
}

func killRunningAgmPIDs(pids []int) {
	killed := 0
	for _, pid := range pids {
		isKilled := killSingleAgmPID(pid)
		if isKilled {
			killed++
		}
	}
	fmt.Printf("  %s✓%s Terminated %d running AGM process(es)\n", constants.ColorGreen, constants.ColorReset, killed)
}

func killSingleAgmPID(pid int) bool {
	hasWindows := runtime.GOOS == "windows"
	if hasWindows {
		cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
		setHiddenProcessAttr(cmd)
		return cmd.Run() == nil
	}
	proc, err := os.FindProcess(pid)
	hasProcErr := err != nil
	if hasProcErr {
		return false
	}
	return proc.Kill() == nil
}

func checkAgmConfirmation(candidates []string, isForce, isDryRun bool) bool {
	if isForce || isDryRun {
		return true
	}
	fmt.Printf("  Ready to remove %d detected AGM item(s). Confirm? [y/N]: ", len(candidates))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	clean := strings.TrimSpace(strings.ToLower(input))
	return clean == "y" || clean == "yes"
}

func applyAgmRemoval(candidates []string, isFullPurge, isDryRun bool) error {
	if isDryRun {
		renderDryRunAgmCandidates(candidates)
		return nil
	}
	for _, p := range candidates {
		removeAgmPathSafe(p)
	}
	if isFullPurge {
		cleanAgmRegistryAndShortcuts()
	}
	cleanAgmInstalledDB()
	fmt.Printf("  %s✓%s Antigravity Manager (AGM) uninstallation completed.\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func renderDryRunAgmCandidates(candidates []string) {
	for _, p := range candidates {
		fmt.Printf("  [dry-run] Would remove: %s\n", p)
	}
}

func removeAgmPathSafe(p string) {
	info, err := os.Stat(p)
	hasErr := err != nil
	if hasErr {
		return
	}
	hasDir := info.IsDir()
	if hasDir {
		_ = os.RemoveAll(p)
		return
	}
	_ = os.Remove(p)
}

func cleanAgmRegistryAndShortcuts() {
	hasWindows := runtime.GOOS == "windows"
	if hasWindows {
		purgeAgmRegistryEntries()
	}
}

func purgeAgmRegistryEntries() {
	keys := []string{"Antigravity Manager", "AntigravityManager", "AGM", "agm", "antigravity-manager"}
	for _, k := range keys {
		cmd := exec.Command("reg", "delete", `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`, "/v", k, "/f")
		setHiddenProcessAttr(cmd)
		_ = cmd.Run()
	}
}

func cleanAgmInstalledDB() {
	splitDB, err := store.OpenInstallationSplitDB()
	hasErr := err != nil
	if hasErr {
		return
	}
	defer splitDB.Close()
	_ = splitDB.RemoveInstalledTool("ag-manager")
	_ = splitDB.RemoveInstalledTool("agm")
	_ = splitDB.RemoveInstalledTool("antigravity-manager")
}

func collectAgmCandidatePaths(isFullPurge bool) []string {
	var paths []string
	home, _ := os.UserHomeDir()
	localApp := os.Getenv("LOCALAPPDATA")
	appData := os.Getenv("APPDATA")
	paths = append(paths, getAgmBinaryPaths(home, localApp)...)
	if isFullPurge {
		paths = append(paths, getAgmPurgePaths(home, localApp, appData)...)
	}
	return filterExistingAgmPaths(paths)
}

func getAgmBinaryPaths(home, localApp string) []string {
	binDir := filepath.Join(localApp, "agm", "bin")
	return []string{
		filepath.Join(binDir, "agm.cmd"),
		filepath.Join(binDir, "agm.ps1"),
		filepath.Join(binDir, "agm.exe"),
		filepath.Join(binDir, "antigravity-manager.cmd"),
		filepath.Join(binDir, "antigravity-manager.exe"),
		filepath.Join(localApp, "Programs", "antigravity-manager"),
		filepath.Join(localApp, "Programs", "agm"),
		filepath.Join(home, ".local", "bin", "agm"),
		filepath.Join(home, ".local", "bin", "agm.cmd"),
		filepath.Join(home, ".local", "bin", "ag-manager"),
		filepath.Join(home, ".local", "bin", "antigravity-manager"),
		filepath.Join(home, "bin", "agm"),
		filepath.Join(home, "bin", "agm.cmd"),
		filepath.Join(home, "bin", "antigravity-manager"),
		"/usr/local/bin/agm",
		"/usr/local/bin/ag-manager",
		"/usr/local/bin/antigravity-manager",
		"/Applications/Antigravity Manager.app",
		"/opt/antigravity-manager",
	}
}

func getAgmPurgePaths(home, localApp, appData string) []string {
	userProf := os.Getenv("USERPROFILE")
	return []string{
		filepath.Join(home, ".agm"),
		filepath.Join(home, ".config", "agm"),
		filepath.Join(home, ".config", "antigravity-manager"),
		filepath.Join(home, ".cache", "agm"),
		filepath.Join(home, ".cache", "antigravity-manager"),
		filepath.Join(appData, "agm"),
		filepath.Join(appData, "antigravity-manager"),
		filepath.Join(localApp, "agm"),
		filepath.Join(localApp, "antigravity-manager"),
		filepath.Join(userProf, "Desktop", "Antigravity Manager.lnk"),
		filepath.Join(userProf, "Desktop", "agm.lnk"),
		filepath.Join(userProf, "Desktop", "antigravity-manager.lnk"),
		filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Antigravity Manager.lnk"),
		filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "agm.lnk"),
		filepath.Join(home, ".config", "autostart", "antigravity-manager.desktop"),
		filepath.Join(home, ".local", "share", "applications", "antigravity-manager.desktop"),
	}
}

func filterExistingAgmPaths(paths []string) []string {
	var existing []string
	for _, p := range paths {
		hasPath := p != ""
		if !hasPath {
			continue
		}
		_, err := os.Stat(p)
		hasExist := err == nil
		if hasExist {
			existing = append(existing, p)
		}
	}
	return existing
}
