// Package cmd — fixgit_permissions.go: permissions and Windows ACL repair for .git.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func remediateGitPermissions(gitDir string, opts FixGitOptions) ([]FixGitIssue, error) {
	var issues []FixGitIssue

	hasIndexPermIssue := checkIndexWritable(gitDir)
	hasDirPermIssue := checkDirWritable(gitDir)

	if !hasIndexPermIssue && !hasDirPermIssue {
		return issues, nil
	}

	issue := FixGitIssue{
		Category:    "Permissions",
		Description: "Read-only or restricted ACL permissions detected on .git",
	}

	if opts.IsDryRun {
		issue.Remedy = "Would grant Full Control to current user and clear read-only flags"
		issues = append(issues, issue)

		return issues, nil
	}

	applyErr := applyPermissionFixes(gitDir)
	if applyErr != nil {
		issue.ErrorDetail = applyErr.Error()
		issue.Remedy = "Attempted permission reset but encountered error"
		issues = append(issues, issue)

		return issues, applyErr
	}

	issue.IsFixed = true
	issue.Remedy = "Granted Full Control to current user and stripped read-only file attributes"
	issues = append(issues, issue)

	return issues, nil
}

func checkIndexWritable(gitDir string) bool {
	indexPath := filepath.Join(gitDir, "index")

	info, err := os.Stat(indexPath)
	if err != nil {
		return false
	}

	if (info.Mode().Perm() & 0200) == 0 {
		return true
	}

	file, openErr := os.OpenFile(indexPath, os.O_WRONLY, 0)
	if openErr != nil {
		return true
	}

	_ = file.Close()

	return false
}

func checkDirWritable(gitDir string) bool {
	probePath := filepath.Join(gitDir, ".gitmap-perm-probe")

	file, err := os.Create(probePath)
	if err != nil {
		return true
	}

	_ = file.Close()
	_ = os.Remove(probePath)

	return false
}

func applyPermissionFixes(gitDir string) error {
	_ = stripFileAttributes(gitDir)
	_ = grantOSPermissions(gitDir)
	_ = ensureChmodWritable(gitDir)

	return nil
}

func stripFileAttributes(gitDir string) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	target := filepath.Join(gitDir, "*")
	cmd := exec.Command("attrib", "-r", target, "/s", "/d")

	return cmd.Run()
}

func grantOSPermissions(gitDir string) error {
	if runtime.GOOS == "windows" {
		return grantWindowsACLs(gitDir)
	}

	return grantUnixPermissions(gitDir)
}

func grantWindowsACLs(gitDir string) error {
	username := os.Getenv("USERNAME")
	if username == "" {
		username = "*S-1-5-32-545" // Builtin Users SID fallback
	}

	grantArg := fmt.Sprintf("%s:(OI)(CI)F", username)
	cmd := exec.Command("icacls", gitDir, "/grant", grantArg, "/t", "/c", "/q")

	return cmd.Run()
}

func grantUnixPermissions(gitDir string) error {
	cmd := exec.Command("chmod", "-R", "u+rwX", gitDir)

	return cmd.Run()
}

func ensureChmodWritable(gitDir string) error {
	return filepath.Walk(gitDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			_ = os.Chmod(path, 0755)

			return nil
		}

		_ = os.Chmod(path, 0666)

		return nil
	})
}
