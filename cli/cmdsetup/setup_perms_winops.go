package cmdsetup

import (
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

var commandRunner = runCommandDiscrete

func runCommandDiscrete(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	return cmd.Run()
}

func runAttribClear(targetDir string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}

	pattern := filepath.Join(targetDir, "*.*")
	err := commandRunner("attrib", "-r", pattern, "/s", "/d")
	if err != nil {
		return apperror.WrapSimple(err, "attrib")
	}

	return nil
}

func runIcaclsGrant(targetDir string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}

	err := commandRunner("icacls", targetDir, "/grant", "Users:(OI)(CI)M", "/T", "/C", "/Q")
	if err != nil {
		return apperror.WrapSimple(err, "icacls.grant")
	}

	return nil
}

func runIcaclsCredentials(credPath string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}

	err := commandRunner("icacls", credPath, "/inheritance:r", "/grant:r", "Administrators:(F)", "SYSTEM:(F)", "/Q")
	if err != nil {
		return apperror.WrapSimple(err, "icacls.credentials")
	}

	return nil
}

func secureSingleCredential(filePath string, isDryRun bool) *apperror.AppError {
	hasFile := fileExists(filePath)
	if !hasFile {
		return nil
	}

	return runIcaclsCredentials(filePath, isDryRun)
}

func secureWindowsCredentials(targetDir string, isDryRun bool) *apperror.AppError {
	wpErr := secureSingleCredential(filepath.Join(targetDir, "wp-config.php"), isDryRun)
	if wpErr != nil {
		return wpErr
	}

	envErr := secureSingleCredential(filepath.Join(targetDir, ".env"), isDryRun)
	if envErr != nil {
		return envErr
	}

	return nil
}

func executeWindowsFixSteps(targetDir string, isDryRun bool) *apperror.AppError {
	attribErr := runAttribClear(targetDir, isDryRun)
	if attribErr != nil {
		return attribErr
	}

	grantErr := runIcaclsGrant(targetDir, isDryRun)
	if grantErr != nil {
		return grantErr
	}

	return secureWindowsCredentials(targetDir, isDryRun)
}

// ApplyWindowsPermissions audits and fixes Windows NTFS directory and file ACLs.
func ApplyWindowsPermissions(opts PermsOptions) (*PermsReport, *apperror.AppError) {
	report := &PermsReport{Violations: make([]string, 0)}
	if !opts.IsFix {
		report.Violations = append(report.Violations, "Windows permission audit requires --fix to apply NTFS ACLs")

		return report, nil
	}

	execErr := executeWindowsFixSteps(opts.TargetDir, opts.IsDryRun)
	if execErr != nil {
		return report, execErr
	}

	report.FixedCount++

	return report, nil
}
