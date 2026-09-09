package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func isSensitiveFile(name string) bool {
	switch name {
	case "wp-config.php", ".env":

		return true
	}
	ext := filepath.Ext(name)

	return ext == ".key" || ext == ".pem"
}

func isWritablePath(relPath string, appType PermsAppType) bool {
	cleanRel := filepath.ToSlash(relPath)
	if appType == PermsAppTypeWordPress {
		return strings.HasPrefix(cleanRel, "wp-content/uploads")
	}
	if appType == PermsAppTypeLaravel {
		return strings.HasPrefix(cleanRel, "storage") || strings.HasPrefix(cleanRel, "bootstrap/cache")
	}

	return false
}

func resolveDirMode(relPath string, appType PermsAppType) os.FileMode {
	if isWritablePath(relPath, appType) {
		return 0775
	}

	return 0755
}

func resolveFileMode(relPath string, appType PermsAppType) os.FileMode {
	base := filepath.Base(relPath)
	if isSensitiveFile(base) {
		return 0600
	}
	if isWritablePath(relPath, appType) {
		return 0664
	}

	return 0644
}

func resolveExpectedUnixMode(relPath string, isDir bool, appType PermsAppType) os.FileMode {
	if isDir {
		return resolveDirMode(relPath, appType)
	}

	return resolveFileMode(relPath, appType)
}

func fixItemMode(path string, expected os.FileMode, isFix, isDryRun bool) (bool, *apperror.AppError) {
	if !isFix {
		return false, nil
	}
	if isDryRun || currentOS == "windows" {
		return true, nil
	}
	err := os.Chmod(path, expected)
	if err != nil {
		return false, apperror.WrapSimple(err, "os.Chmod")
	}

	return true, nil
}

func auditItem(path, relPath string, info os.FileInfo, opts PermsOptions, report *PermsReport) *apperror.AppError {
	report.ScannedCount++
	expected := resolveExpectedUnixMode(relPath, info.IsDir(), opts.AppType)
	actual := info.Mode().Perm()
	isMatch := actual == expected
	if isMatch {
		return nil
	}
	msg := fmt.Sprintf("%s: mode %04o != expected %04o", relPath, actual, expected)
	report.Violations = append(report.Violations, msg)
	wasFixed, fixErr := fixItemMode(path, expected, opts.IsFix, opts.IsDryRun)
	if fixErr != nil {
		return fixErr
	}
	if wasFixed {
		report.FixedCount++
	}

	return nil
}

func createPermsWalkFunc(targetDir string, opts PermsOptions, report *PermsReport) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(targetDir, path)
		if relErr != nil {
			return relErr
		}
		appErr := auditItem(path, rel, info, opts, report)
		if appErr != nil {
			return appErr
		}

		return nil
	}
}

func applyOwnership(targetDir, owner string) *apperror.AppError {
	hasOwner := strings.TrimSpace(owner) != ""
	if !hasOwner {
		return nil
	}
	_, lookErr := exec.LookPath("chown")
	if lookErr != nil {
		return nil
	}
	runErr := commandRunner("chown", "-R", owner, targetDir)
	if runErr != nil {
		return apperror.WrapSimple(runErr, "chown")
	}

	return nil
}

func handleUnixOwnership(opts PermsOptions) *apperror.AppError {
	if !opts.IsFix {
		return nil
	}

	return applyOwnership(opts.TargetDir, opts.Owner)
}

// ApplyUnixPermissions audits and fixes Linux/Unix directory and file permissions.
func ApplyUnixPermissions(opts PermsOptions) (*PermsReport, *apperror.AppError) {
	report := &PermsReport{Violations: make([]string, 0)}
	walkFn := createPermsWalkFunc(opts.TargetDir, opts, report)
	walkErr := filepath.Walk(opts.TargetDir, walkFn)
	if walkErr != nil {
		return report, apperror.WrapSimple(walkErr, "filepath.Walk")
	}
	ownErr := handleUnixOwnership(opts)
	if ownErr != nil {
		return report, ownErr
	}

	return report, nil
}
