package cmdagy

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	isAgyUninstallAll    bool
	isAgyUninstallDryRun bool
	isAgyUninstallForce  bool
	agyUninstallBackup   string
)

var agyUninstallCmd = &cobra.Command{
	Use:     "uninstall [flags]",
	Aliases: []string{"remove", "rm", "purge"},
	Short:   "Uninstall Antigravity (AGY) CLI and application",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAGYUninstall(isAgyUninstallAll, isAgyUninstallForce, isAgyUninstallDryRun, agyUninstallBackup)
	},
}

var agyUninstallAllCmd = &cobra.Command{
	Use:   "uninstall-all [flags]",
	Short: "Full purge of Antigravity (AGY), caches, brain, and configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAGYUninstall(true, isAgyUninstallForce, isAgyUninstallDryRun, agyUninstallBackup)
	},
}

func initAgyUninstallCmd() {
	bindAgyUninstallFlags()
	AgyCmd.AddCommand(agyUninstallCmd)
	AgyCmd.AddCommand(agyUninstallAllCmd)
}

func bindAgyUninstallFlags() {
	agyUninstallCmd.Flags().BoolVarP(&isAgyUninstallAll, "all", "a", false, "Full purge including .gemini, brain, and caches")
	agyUninstallCmd.Flags().BoolVarP(&isAgyUninstallDryRun, "dry-run", "n", false, "Simulate uninstallation without deleting")
	agyUninstallCmd.Flags().BoolVarP(&isAgyUninstallForce, "force", "f", false, "Force uninstallation without confirmation")
	agyUninstallCmd.Flags().StringVarP(&agyUninstallBackup, "backup", "b", "", "Custom path for state snapshot JSON")
	agyUninstallAllCmd.Flags().BoolVarP(&isAgyUninstallDryRun, "dry-run", "n", false, "Simulate uninstallation without deleting")
	agyUninstallAllCmd.Flags().BoolVarP(&isAgyUninstallForce, "force", "f", false, "Force uninstallation without confirmation")
	agyUninstallAllCmd.Flags().StringVarP(&agyUninstallBackup, "backup", "b", "", "Custom path for state snapshot JSON")
}

// RunAGYUninstall removes Antigravity binaries, caches, and state according to purge tier.
func RunAGYUninstall(isFullPurge, isForce, isDryRun bool, backupPath string) error {
	printUninstallHeader(isFullPurge)
	candidates, collectErr := buildUninstallCandidates(isFullPurge)
	hasCollectErr := collectErr != nil
	if hasCollectErr {
		return collectErr
	}
	validErr := validateAllCandidates(candidates)
	hasValidErr := validErr != nil
	if hasValidErr {
		return validErr
	}
	return executeUninstallStages(candidates, isFullPurge, isForce, isDryRun, backupPath)
}

func executeUninstallStages(candidates []string, isFullPurge, isForce, isDryRun bool, backupPath string) error {
	snapErr := handlePreUninstallSnapshot(isFullPurge, isDryRun, backupPath)
	hasSnapErr := snapErr != nil
	if hasSnapErr {
		return snapErr
	}
	terminateRunningProcesses(isDryRun)
	isConfirmed := checkUninstallConfirmation(candidates, isForce, isDryRun)
	if !isConfirmed {
		fmt.Println("Uninstallation cancelled by user.")
		return nil
	}
	return applyUninstallRemoval(candidates, isFullPurge, isDryRun)
}

func printUninstallHeader(isFullPurge bool) {
	if isFullPurge {
		fmt.Printf("%s[GitMap]%s Starting Antigravity (AGY) FULL PURGE...\n", constants.ColorCyan, constants.ColorReset)
		return
	}
	fmt.Printf("%s[GitMap]%s Starting Antigravity (AGY) Standard Uninstallation...\n", constants.ColorCyan, constants.ColorReset)
}

func buildUninstallCandidates(isFullPurge bool) ([]string, error) {
	var list []string
	list = append(list, collectAGYStandardPaths()...)
	if isFullPurge {
		list = append(list, collectAGYPurgePaths()...)
	}
	return list, nil
}

func validateAllCandidates(candidates []string) error {
	for _, p := range candidates {
		isSafe := IsPathSafeToDelete(p)
		if !isSafe {
			return apperror.NewSimple(fmt.Sprintf("safety violation: protected path %s cannot be deleted", p), "E9001")
		}
	}
	return nil
}

// IsPathSafeToDelete verifies path is not a system root, git repo, or d:\work directory.
func IsPathSafeToDelete(path string) bool {
	clean := filepath.Clean(path)
	hasEmpty := clean == "" || clean == "." || clean == ".."
	if hasEmpty {
		return false
	}
	isRoot := isRootDirectory(clean)
	if isRoot {
		return false
	}
	isWorkOverlap := IsWorkDirectoryOverlap(clean)
	if isWorkOverlap {
		return false
	}
	isGitRepo := isGitRepositoryPath(clean)
	if isGitRepo {
		return false
	}
	return true
}

func isRootDirectory(clean string) bool {
	vol := filepath.VolumeName(clean)
	isVolRoot := clean == vol || clean == vol+"\\" || clean == vol+"/"
	if isVolRoot {
		return true
	}
	return clean == "/" || clean == "\\"
}

// IsWorkDirectoryOverlap validates if a path overlaps with or is a parent of d:\work.
func IsWorkDirectoryOverlap(path string) bool {
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

func handlePreUninstallSnapshot(isFullPurge, isDryRun bool, backupPath string) error {
	if !isFullPurge {
		return nil
	}
	if isDryRun {
		fmt.Printf("  [dry-run] Would export AGY restore snapshot before full purge (backup: %s)\n", backupPath)
		return nil
	}
	snap, err := ExportAGYRestoreSnapshot(backupPath)
	hasErr := err != nil
	if hasErr {
		return fmt.Errorf("snapshot export failed prior to purge: %w", err)
	}
	fmt.Printf("  %s✓%s Snapshot saved (%d projects, %d conversations)\n",
		constants.ColorGreen, constants.ColorReset, snap.TotalProjects, snap.TotalConvs)
	return nil
}

func terminateRunningProcesses(isDryRun bool) {
	procs, err := DiscoverTargetProcesses()
	hasProcs := err == nil && len(procs) > 0
	if !hasProcs {
		return
	}
	if isDryRun {
		fmt.Printf("  [dry-run] Would terminate %d running Antigravity processes\n", len(procs))
		return
	}
	killed, _ := TerminateProcesses(procs)
	fmt.Printf("  %s✓%s Terminated %d running processes\n", constants.ColorGreen, constants.ColorReset, killed)
}

func checkUninstallConfirmation(candidates []string, isForce, isDryRun bool) bool {
	if isForce || isDryRun {
		return true
	}
	fmt.Printf("  Ready to remove %d detected Antigravity item(s). Confirm? [y/N]: ", len(candidates))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	clean := strings.TrimSpace(strings.ToLower(input))
	return clean == "y" || clean == "yes"
}

func applyUninstallRemoval(candidates []string, isFullPurge, isDryRun bool) error {
	if isDryRun {
		renderDryRunCandidates(candidates)
		return nil
	}
	for _, p := range candidates {
		removePathSafe(p)
	}
	if isFullPurge {
		cleanEmptyGeminiDir()
	}
	purgeAgyFromStores()
	fmt.Printf("  %s✓%s Antigravity (AGY) uninstallation completed.\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func purgeAgyFromStores() {
	purgeAgySplitDB()
	purgeAgyRootDB()
}

func purgeAgySplitDB() {
	splitDB, err := store.OpenInstallationSplitDB()
	hasErr := err != nil
	if hasErr {
		return
	}
	defer splitDB.Close()
	_ = splitDB.RemoveInstalledTool(constants.ToolAntigravity)
	_ = splitDB.RemoveInstalledTool(constants.ToolAgy)
}

func purgeAgyRootDB() {
	rootDB, err := store.OpenDefault()
	hasErr := err != nil
	if hasErr {
		return
	}
	defer rootDB.Close()
	conn := rootDB.Conn()
	hasConn := conn != nil
	if hasConn {
		_, _ = conn.Exec(constants.SQLDeleteInstalledTool, constants.ToolAntigravity)
		_, _ = conn.Exec(constants.SQLDeleteInstalledTool, constants.ToolAgy)
	}
	_ = rootDB.SyncKnownSplitDatabases()
}

func renderDryRunCandidates(candidates []string) {
	for _, p := range candidates {
		fmt.Printf("  [dry-run] Would remove: %s\n", p)
	}
}

func removePathSafe(p string) {
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

func cleanEmptyGeminiDir() {
	home, err := os.UserHomeDir()
	hasErr := err != nil
	if hasErr {
		return
	}
	geminiDir := filepath.Join(home, ".gemini")
	entries, readErr := os.ReadDir(geminiDir)
	hasReadErr := readErr != nil
	if hasReadErr {
		return
	}
	hasEmpty := len(entries) == 0
	if hasEmpty {
		_ = os.Remove(geminiDir)
	}
}

func collectAGYStandardPaths() []string {
	var paths []string
	home, _ := os.UserHomeDir()
	localApp := os.Getenv("LOCALAPPDATA")
	progFiles := os.Getenv("ProgramFiles")
	hasWindows := runtime.GOOS == "windows"
	if hasWindows {
		paths = append(paths, getWindowsAGYBinPaths(home, localApp, progFiles)...)
	}
	hasNonWindows := runtime.GOOS != "windows"
	if hasNonWindows {
		paths = append(paths, getUnixAGYBinPaths(home)...)
	}
	return filterExistingPaths(paths)
}

func getWindowsAGYBinPaths(home, localApp, progFiles string) []string {
	binDir := filepath.Join(localApp, "agy", "bin")
	return []string{
		filepath.Join(binDir, "agy.exe"),
		filepath.Join(binDir, "agy.cmd"),
		filepath.Join(binDir, "agy.ps1"),
		filepath.Join(binDir, "antigravity.cmd"),
		filepath.Join(binDir, "antigravity.ps1"),
		filepath.Join(binDir, "antigravity.exe"),
		filepath.Join(home, ".antigravity", "bin", "agy.exe"),
		filepath.Join(localApp, "Programs", "Antigravity"),
		filepath.Join(localApp, "Programs", "antigravity"),
		filepath.Join(progFiles, "Antigravity"),
	}
}

func getUnixAGYBinPaths(home string) []string {
	return []string{
		filepath.Join(home, ".local", "bin", "agy"),
		filepath.Join(home, ".local", "bin", "antigravity"),
		filepath.Join(home, ".local", "agy", "bin", "agy"),
		filepath.Join(home, ".antigravity", "bin", "agy"),
		filepath.Join(home, ".agy", "bin", "agy"),
		"/usr/local/bin/agy",
		"/usr/local/bin/antigravity",
		"/Applications/Antigravity.app",
		filepath.Join(home, "Applications", "Antigravity.app"),
		"/opt/antigravity",
		filepath.Join(home, ".local", "share", "antigravity"),
		filepath.Join(home, ".local", "share", "applications", "antigravity.desktop"),
		"/usr/share/applications/antigravity.desktop",
	}
}

func collectAGYPurgePaths() []string {
	var paths []string
	home, _ := os.UserHomeDir()
	localApp := os.Getenv("LOCALAPPDATA")
	appData := os.Getenv("APPDATA")
	hasHome := home != ""
	if hasHome {
		paths = append(paths, getHomeGeminiPaths(home)...)
	}
	hasWindows := runtime.GOOS == "windows"
	if hasWindows {
		paths = append(paths, getWindowsPurgeDirs(localApp, appData)...)
	}
	return filterExistingPaths(paths)
}

func getHomeGeminiPaths(home string) []string {
	return []string{
		filepath.Join(home, ".gemini", "antigravity"),
		filepath.Join(home, ".gemini", "config", "projects"),
		filepath.Join(home, ".gemini", "antigravity-cli"),
		filepath.Join(home, ".cache", "antigravity"),
	}
}

func getWindowsPurgeDirs(localApp, appData string) []string {
	return []string{
		filepath.Join(localApp, "antigravity"),
		filepath.Join(localApp, "agy"),
		filepath.Join(localApp, "antigravity-updater"),
		filepath.Join(appData, "Antigravity"),
	}
}

func filterExistingPaths(paths []string) []string {
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
