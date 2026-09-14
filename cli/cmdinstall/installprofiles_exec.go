package cmdinstall

import (
	"errors"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runInstallProfile(profileName string, opts installOptions) error {
	p, found := FindInstallProfile(profileName)
	if !found {
		fmt.Printf("Unknown installation profile: %s\n", profileName)

		return nil
	}

	installed := loadInstalledLookup()
	if opts.Tree {
		renderProfileTree(p, installed)

		return nil
	}

	return handleProfileExecution(p, opts, installed)
}

func handleProfileExecution(p InstallProfile, opts installOptions, installed map[string]string) error {
	splitDB, _ := store.OpenInstallationSplitDB()
	if splitDB == nil {
		return executeProfileWorkflow(p, opts, installed, nil)
	}

	defer splitDB.Close()
	if isProfileAlreadyInstalled(p.Name, splitDB, opts) {
		handleAlreadyInstalled(p, splitDB, installed)

		return nil
	}

	return executeProfileWorkflow(p, opts, installed, splitDB)
}

func isProfileAlreadyInstalled(name string, splitDB *store.InstallationSplitDB, opts installOptions) bool {
	if opts.Force {
		return false
	}

	return splitDB.IsProfileInstalled(name)
}

func handleAlreadyInstalled(p InstallProfile, splitDB *store.InstallationSplitDB, installed map[string]string) {
	rec, _ := splitDB.GetProfileInstallation(p.Name)
	installedAt := resolveInstalledTime(rec)
	fmt.Printf("[INFO] Profile '%s' is already installed (installed at: %s)\n", p.Name, installedAt)
	renderProfileTree(p, installed)
}

func resolveInstalledTime(rec *store.ProfileInstallationRecord) string {
	if rec != nil && rec.InstalledAt != "" {
		return rec.InstalledAt
	}

	return "previously"
}

func executeProfileWorkflow(
	p InstallProfile, opts installOptions,
	installed map[string]string, splitDB *store.InstallationSplitDB,
) error {
	printProfileStartHeader(p, installed)
	profileId := startProfileAudit(p, splitDB)
	start := time.Now()

	installedCount, execErr := executeProfileToolsWithAudit(p, opts, profileId, splitDB)
	durationMs := time.Since(start).Milliseconds()

	if execErr != nil {
		handleProfileFailure(profileId, durationMs, installedCount, execErr, splitDB)

		return execErr
	}

	handleProfileSuccess(p, profileId, durationMs, installedCount, splitDB)

	return nil
}

func startProfileAudit(p InstallProfile, splitDB *store.InstallationSplitDB) int64 {
	if splitDB == nil {
		return 0
	}

	id, err := splitDB.RecordProfileStart(p.Name, p.Title, p.Description, len(p.Tools))
	if err != nil {
		return 0
	}

	return id
}

func handleProfileFailure(
	profileId, durationMs int64, installedCount int,
	execErr error, splitDB *store.InstallationSplitDB,
) {
	if splitDB == nil {
		return
	}

	stackTrace := extractAppStackTrace(execErr)
	_ = splitDB.RecordProfileCompletion(profileId, durationMs, false, installedCount, 1, stackTrace, execErr.Error(), "")
}

func handleProfileSuccess(
	p InstallProfile, profileId, durationMs int64,
	installedCount int, splitDB *store.InstallationSplitDB,
) {
	if splitDB != nil {
		_ = splitDB.RecordProfileCompletion(profileId, durationMs, true, installedCount, 0, "", "", "")
	}

	printProfileSummary(p, installedCount)
}

func extractAppStackTrace(err error) string {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return appErr.Stack
	}

	return ""
}

func printProfileStartHeader(p InstallProfile, installed map[string]string) {
	fmt.Printf("\n=== Installing Profile: %s (%s) ===\n", p.Name, p.Title)
	fmt.Printf("Description: %s\n", p.Description)
	fmt.Printf("Total Tools: %d\n\n", len(p.Tools))
	renderProfileTreeNodes(p, installed)
	fmt.Println()
}

func executeProfileToolsWithAudit(
	p InstallProfile, opts installOptions,
	profileId int64, splitDB *store.InstallationSplitDB,
) (int, error) {
	installed := loadInstalledLookup()
	successCount := 0
	for idx, tool := range p.Tools {
		stepNum := idx + 1
		if isToolInstalledAlready(tool, installed) {
			announceToolAlreadyInstalled(stepNum, len(p.Tools), tool, installed)
			recordPackageAudit(profileId, tool, true, 0, 0, "", "", splitDB)
			successCount++
			continue
		}

		stepErr, dur := runSingleProfileStep(stepNum, len(p.Tools), tool, opts)
		if stepErr != nil {
			recordPackageAudit(profileId, tool, false, 1, dur, stepErr.Error(), extractAppStackTrace(stepErr), splitDB)

			return successCount, stepErr
		}

		recordPackageAudit(profileId, tool, true, 0, dur, "", "", splitDB)
		successCount++
	}

	return successCount, nil
}

func runSingleProfileStep(step, total int, tool string, opts installOptions) (error, int64) {
	start := time.Now()
	err := installSingleProfileTool(step, total, tool, opts)
	dur := time.Since(start).Milliseconds()

	return err, dur
}

func isToolInstalledAlready(tool string, installed map[string]string) bool {
	status, _ := resolveToolStatus(tool, installed)

	return status == constants.StatusInstalled
}

func announceToolAlreadyInstalled(step, total int, tool string, installed map[string]string) {
	_, ver := resolveToolStatus(tool, installed)
	fmt.Printf("  ✓ [%d/%d] %s is already installed (%s)\n", step, total, tool, ver)
}

func installSingleProfileTool(step, total int, tool string, opts installOptions) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = apperror.NewWithDetails("cmdinstall.profile", "E5001", fmt.Sprintf("tool install failed: %v", r), "cmdinstall", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
		}
	}()

	fmt.Printf("\n  → [%d/%d] Installing %s...\n", step, total, tool)
	toolOpts := opts
	toolOpts.Tool = tool
	toolOpts.Yes = true
	executeInstall(toolOpts)

	return nil
}

func recordPackageAudit(
	profId int64, tool string, isSuccess bool,
	exitCode int, dur int64, stderr, stack string, splitDB *store.InstallationSplitDB,
) {
	if splitDB == nil {
		return
	}

	status := resolvePackageStatus(isSuccess)
	_ = splitDB.RecordPackageInstallation(store.PackageInstallationRecord{
		ProfileInstallationId: profId,
		PackageName:           tool,
		Status:                status,
		IsSuccess:             isSuccess,
		ExitCode:              exitCode,
		DurationMs:            dur,
		Stderr:                stderr,
		StackTrace:            stack,
	})
}

func resolvePackageStatus(isSuccess bool) string {
	if isSuccess {
		return "installed"
	}

	return "failed"
}

func printProfileSummary(p InstallProfile, count int) {
	fmt.Printf("\n" + constants.ColorGreen + "✓" + constants.ColorReset)
	fmt.Printf(" Profile '%s' setup complete! (%d/%d tools processed)\n\n", p.Name, count, len(p.Tools))
}
