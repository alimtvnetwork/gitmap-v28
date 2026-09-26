package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func runMigrateWizard(args []string) error {
	scriptPath := findMigrateScript()
	if scriptPath == "" {
		return apperror.NewSimple("Migration wizard script not found.", "E1001")
	}

	psArgs := append([]string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath}, args...)
	cmd := exec.Command("powershell", psArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func findMigrateScript() string {
	candidates := defaultCandidatePaths()
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append([]string{
			filepath.Join(exeDir, "migrate.ps1"),
			filepath.Join(exeDir, "scripts", "migrate-gitmap-v28.ps1"),
		}, candidates...)
	}

	return searchCandidatePaths(candidates)
}

func defaultCandidatePaths() []string {
	return []string{
		"migrate.ps1",
		filepath.Join("scripts", "migrate-gitmap-v28.ps1"),
		filepath.Join(os.Getenv("LOCALAPPDATA"), "gitmap-cli", "migrate.ps1"),
		filepath.Join("D:\\work\\gitmap", "scripts", "migrate-gitmap-v28.ps1"),
		filepath.Join("D:\\work\\gitmap", "migrate.ps1"),
	}
}

func searchCandidatePaths(candidates []string) string {
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	return ""
}
