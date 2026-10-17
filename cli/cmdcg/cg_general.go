package cmdcg

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/charmbracelet/lipgloss"
)

// RunCGGeneral executes the general repository guidelines installation across targets.
func RunCGGeneral(args []string) error {
	var (
		isDryRun bool
		isAll    bool
		repos    []string
	)

	for _, arg := range args {
		switch arg {
		case "--dry-run", "-n":
			isDryRun = true
		case "--all", "-a":
			isAll = true
		case "general", "general-install", "install":
			// skip command noise tokens
		default:
			if !strings.HasPrefix(arg, "-") {
				repos = append(repos, arg)
			}
		}
	}

	targets := resolveGeneralTargets(repos, isAll)
	if len(targets) == 0 {
		fmt.Println("No repositories to process for general guidelines installation.")

		return nil
	}

	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Bold(true)
	fmt.Println(titleStyle.Render(fmt.Sprintf("Running general repository guidelines installer on %d repositories...", len(targets))))

	cgPath := findLocalCodingGuidelinesPath()

	for _, target := range targets {
		runGeneralInstallForRepo(target, cgPath, isDryRun)
	}

	fmt.Println(titleStyle.Render("General installation complete."))

	return nil
}

func resolveGeneralTargets(specified []string, isAll bool) []string {
	if len(specified) > 0 {
		return specified
	}

	if isAll {
		return resolveAllCGRepos()
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	if childRepos, findErr := fsutil.DiscoverChildGitRepos(cwd); findErr == nil && len(childRepos) > 0 {
		return childRepos
	}

	return []string{cwd}
}

func findLocalCodingGuidelinesPath() string {
	candidates := []string{
		"d:\\work\\coding-guidelines",
		"d:/work/coding-guidelines",
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr == nil {
		candidates = append(candidates,
			filepath.Join(cwd, "coding-guidelines"),
			filepath.Join(filepath.Dir(cwd), "coding-guidelines"),
		)
	}

	for _, cand := range candidates {
		configPath := filepath.Join(cand, "general-install-config.json")
		if _, statErr := os.Stat(configPath); statErr == nil {
			return cand
		}
	}

	return ""
}

func runGeneralInstallForRepo(repo string, localGuidelinesPath string, isDryRun bool) {
	fmt.Printf("--> Installing general guidelines into: %s\n", repo)

	var cmd *exec.Cmd

	hasLocal := len(localGuidelinesPath) > 0

	if hasLocal {
		cmd = buildLocalGeneralCmd(repo, localGuidelinesPath, isDryRun)
	}

	if !hasLocal {
		cmd = buildRemoteGeneralCmd(repo, isDryRun)
	}

	cmd.Dir = repo
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	execErr := cmd.Run()
	if execErr != nil {
		fmt.Printf("    [X] Failed: %v\n    %s\n", execErr, stderr.String())

		return
	}

	outStr := strings.TrimSpace(stdout.String())
	if len(outStr) > 0 {
		lines := strings.Split(outStr, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if len(trimmed) > 0 {
				fmt.Printf("    %s\n", trimmed)
			}
		}
	}

	fmt.Printf("    [OK] Finished %s\n", repo)
}

func buildLocalGeneralCmd(targetRepo, guidelinesDir string, isDryRun bool) *exec.Cmd {
	if runtime.GOOS == "windows" {
		scriptPath := filepath.Join(guidelinesDir, "general-install.ps1")
		args := []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath, "-TargetDir", targetRepo}
		if isDryRun {
			args = append(args, "-DryRun")
		}

		return exec.Command("powershell", args...)
	}

	scriptPath := filepath.Join(guidelinesDir, "general-install.sh")
	args := []string{scriptPath, "--target", targetRepo}
	if isDryRun {
		args = append(args, "--dry-run")
	}

	return exec.Command("bash", args...)
}

func buildRemoteGeneralCmd(targetRepo string, isDryRun bool) *exec.Cmd {
	if runtime.GOOS == "windows" {
		flag := ""
		if isDryRun {
			flag = " -DryRun"
		}

		psCmd := fmt.Sprintf("& ([scriptblock]::Create((irm https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/general-install.ps1))) -TargetDir '%s'%s", targetRepo, flag)

		return exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	}

	flag := ""
	if isDryRun {
		flag = " --dry-run"
	}

	shCmd := fmt.Sprintf("curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/general-install.sh | bash -s -- --target '%s'%s", targetRepo, flag)

	return exec.Command("bash", "-c", shCmd)
}
