package cmdcursor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func findWindowsCursor() (string, bool) {
	username := os.Getenv("USERNAME")
	sysDrive := os.Getenv("SystemDrive")
	if sysDrive == "" {
		sysDrive = "C:"
	}
	candidates := []string{
		filepath.Join(sysDrive, "Users", username, "AppData", "Local", "Programs", "cursor", "Cursor.exe"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "cursor", "Cursor.exe"),
		filepath.Join(os.Getenv("ProgramFiles"), "Cursor", "Cursor.exe"),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Programs", "cursor", "Cursor.exe"),
	}
	for _, cand := range candidates {
		if cand != "" && isFile(cand) {
			return cand, true
		}
	}
	if path, err := exec.LookPath("cursor.cmd"); err == nil {
		return path, true
	}
	if path, err := exec.LookPath("cursor.exe"); err == nil {
		return path, true
	}
	return "", false
}

func findLinuxCursor() (string, bool) {
	home, _ := os.UserHomeDir()
	candidates := []string{
		"/usr/bin/cursor",
		"/usr/local/bin/cursor",
		"/opt/cursor/cursor",
		"/opt/cursor/Cursor.AppImage",
		filepath.Join(home, ".local", "bin", "cursor"),
	}
	for _, cand := range candidates {
		if cand != "" && isFile(cand) {
			return cand, true
		}
	}
	if path, err := exec.LookPath("cursor"); err == nil {
		return path, true
	}
	return "", false
}

func findDarwinCursor() (string, bool) {
	home, _ := os.UserHomeDir()
	candidates := []string{
		"/Applications/Cursor.app/Contents/MacOS/Cursor",
		filepath.Join(home, "Applications", "Cursor.app", "Contents", "MacOS", "Cursor"),
	}
	for _, cand := range candidates {
		if cand != "" && isFile(cand) {
			return cand, true
		}
	}
	if path, err := exec.LookPath("cursor"); err == nil {
		return path, true
	}
	return "", false
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// FindCursorExecutable locates the Cursor IDE executable on the host system.
func FindCursorExecutable() (string, bool) {
	switch runtime.GOOS {
	case "windows":
		return findWindowsCursor()
	case "darwin":
		return findDarwinCursor()
	default:
		return findLinuxCursor()
	}
}

func isPathPresent(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func matchRepoByExactName(repos []model.ScanRecord, target string) (string, bool) {
	for _, r := range repos {
		if strings.EqualFold(r.RepoName, target) || strings.EqualFold(r.Slug, target) {
			return r.AbsolutePath, true
		}
	}
	return "", false
}

func matchRepoBySubstring(repos []model.ScanRecord, targetLower string) (string, bool) {
	for _, r := range repos {
		if strings.Contains(strings.ToLower(r.Slug), targetLower) || strings.Contains(strings.ToLower(r.RepoName), targetLower) {
			return r.AbsolutePath, true
		}
	}
	return "", false
}

func findRepoPathByTarget(target string) (string, bool) {
	s, err := store.OpenDefault()
	if err != nil {
		return "", false
	}
	defer s.Close()

	records, _ := s.FindBySlug(target)
	if len(records) > 0 {
		return records[0].AbsolutePath, true
	}
	repos, err := s.ListRepos()
	if err != nil {
		return "", false
	}
	if match, found := matchRepoByExactName(repos, target); found {
		return match, true
	}
	return matchRepoBySubstring(repos, strings.ToLower(target))
}

func resolveOpenTarget(args []string) (string, error) {
	target := "."
	if len(args) > 0 && args[0] != "" {
		target = args[0]
	}
	absPath, err := filepath.Abs(target)
	if err != nil {
		return "", apperror.WrapSimple(err, "resolve target path")
	}
	if isPathPresent(absPath) {
		return absPath, nil
	}
	if repoPath, found := findRepoPathByTarget(target); found {
		return repoPath, nil
	}
	return absPath, nil
}

func launchCursorProcess(exePath, absTarget string) error {
	cmd := exec.Command(exePath, absTarget)
	if err := cmd.Start(); err != nil {
		return apperror.WrapSimple(err, "launch cursor process")
	}
	return nil
}

// RunCursorOpen opens the specified target directory or file in Cursor.
func RunCursorOpen(args []string) error {
	target, err := resolveOpenTarget(args)
	if err != nil {
		return err
	}
	exePath, isFound := FindCursorExecutable()
	if !isFound {
		return apperror.NewWithDetails("cmd.cursor.open", "E1031", "Cursor executable not found on system. Run 'gitmap cursor install' to configure.", "cmdcursor", apperror.ErrorTypeNotFound, apperror.SeverityError, nil)
	}
	if launchErr := launchCursorProcess(exePath, target); launchErr != nil {
		return launchErr
	}
	fmt.Printf("%s✔ Opened in Cursor:%s %s\n", constants.ColorGreen, constants.ColorReset, target)
	return nil
}
