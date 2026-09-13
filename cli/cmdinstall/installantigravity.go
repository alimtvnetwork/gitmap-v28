package cmdinstall

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func announceAntigravityPlatform(opts installOptions) (AntigravityPlatformInfo, *apperror.AppError) {
	info, appErr := detectAntigravityPlatform(opts.Version, "", opts.Prefix)

	if appErr != nil {
		return AntigravityPlatformInfo{}, appErr
	}

	announcePlatform(info, resolveBuildVersion(opts.Version))

	return info, nil
}

func validateAntigravityPrereqs(info AntigravityPlatformInfo) *apperror.AppError {
	if appErr := verifyAntigravityPrerequisites(info); appErr != nil {
		return appErr
	}

	return nil
}

func hasExistingAntigravity(isForce bool) bool {
	if isForce {
		return false
	}

	path, isFound := findInstalledAntigravityDesktopPath()

	if isFound {
		fmt.Printf("  ✓ Google Antigravity Desktop IDE is already installed (%s)\n", path)

		return true
	}

	return false
}

func handleAntigravityDryRun(info AntigravityPlatformInfo, isDryRun bool) bool {
	if !isDryRun {
		return false
	}

	fmt.Println("  [dry-run] Would download and install Google Antigravity Desktop IDE")
	fmt.Printf("  [dry-run] URL: %s\n", info.DownloadURL)
	fmt.Printf("  [dry-run] Destination: %s\n", info.InstallDir)

	return true
}

func checkAntigravityPreconditions(info AntigravityPlatformInfo, opts installOptions) bool {
	if hasExisting := hasExistingAntigravity(opts.Force); hasExisting {
		return true
	}

	return handleAntigravityDryRun(info, opts.DryRun)
}

func verifyAntigravityPostInstall() *apperror.AppError {
	_, isFound := findInstalledAntigravityDesktopPath()

	if !isFound {
		return apperror.NewSimple("desktop application binary not found", ErrPrerequisiteFailed)
	}

	return nil
}

func finalizeAntigravityInstall() {
	path, _ := findInstalledAntigravityDesktopPath()
	fmt.Printf(constants.ColorGreen+"✓"+constants.ColorReset+" Google Antigravity Desktop IDE installed: %s\n", path)
	recordAntigravityDesktopInstalled(path)
}

func executeAntigravityDeployment(opts installOptions) error {
	fmt.Println("Installing Google Antigravity Desktop IDE...")

	if err := installAntigravityDesktopPlatform(opts); err != nil {
		return apperror.WrapSimple(err, "installAntigravityDesktopPlatform")
	}

	if appErr := verifyAntigravityPostInstall(); appErr != nil {
		return appErr
	}

	finalizeAntigravityInstall()

	return nil
}

func runInstallAntigravityWithOpts(opts installOptions) error {
	info, errAnnounce := announceAntigravityPlatform(opts)

	if errAnnounce != nil {
		return errAnnounce
	}

	if errPrereq := validateAntigravityPrereqs(info); errPrereq != nil {
		return errPrereq
	}

	if isHandled := checkAntigravityPreconditions(info, opts); isHandled {
		return nil
	}

	return executeAntigravityDeployment(opts)
}
