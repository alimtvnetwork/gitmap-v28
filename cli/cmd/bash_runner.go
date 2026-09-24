package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runBash(args []string) error {
	gitPath, hasGit := findGitExecutable()
	if !hasGit {
		printGitMissingAdvice()
		return apperror.NewSimple("git is required to run bash", "E_GIT_NOT_FOUND")
	}

	recordGitPathInDB(gitPath)
	bashPath, hasBash := findBashExecutable(gitPath)
	if !hasBash {
		return apperror.NewSimple("bash not found", "E_BASH_NOT_FOUND")
	}

	return executeBashProcess(bashPath, args)
}

func runShell(args []string) error {
	return runBash(args)
}

func findGitExecutable() (string, bool) {
	if p, err := exec.LookPath("git"); err == nil && p != "" {
		return p, true
	}
	if runtime.GOOS == constants.OSWindows {
		candidates := []string{
			`C:\Program Files\Git\cmd\git.exe`,
			`C:\Program Files\Git\bin\git.exe`,
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				return c, true
			}
		}
	}
	return "", false
}

func recordGitPathInDB(gitPath string) {
	if gitPath == "" {
		return
	}
	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()
	_ = db.SetSetting("git_path", gitPath)
}

func findBashExecutable(gitPath string) (string, bool) {
	if runtime.GOOS == constants.OSWindows {
		return findWindowsGitBash(gitPath)
	}
	if p, err := exec.LookPath("bash"); err == nil && p != "" {
		return p, true
	}
	if p, err := exec.LookPath("sh"); err == nil && p != "" {
		return p, true
	}
	return "", false
}

func findWindowsGitBash(gitPath string) (string, bool) {
	candidates := []string{
		`C:\Program Files\Git\bin\bash.exe`,
		`C:\Program Files\Git\usr\bin\bash.exe`,
	}
	if gitPath != "" {
		candidates = append(candidates, deriveGitBashFromGit(gitPath)...)
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, true
		}
	}
	if p, err := exec.LookPath("bash"); err == nil && p != "" {
		return p, true
	}
	return "", false
}

func deriveGitBashFromGit(gitPath string) []string {
	normalized := strings.ReplaceAll(gitPath, "\\", "/")
	idx := strings.LastIndex(normalized, "/cmd/git.exe")
	if idx >= 0 {
		root := gitPath[:idx]
		return []string{
			root + `\bin\bash.exe`,
			root + `\usr\bin\bash.exe`,
		}
	}
	return nil
}

func executeBashProcess(bashPath string, args []string) error {
	cmd := buildBashCmd(bashPath, args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "execute bash")
	}
	return nil
}

func buildBashCmd(bashPath string, args []string) *exec.Cmd {
	if len(args) == 0 {
		return exec.Command(bashPath)
	}
	if args[0] == "-c" || strings.HasPrefix(args[0], "-") {
		return exec.Command(bashPath, args...)
	}
	return exec.Command(bashPath, "-c", strings.Join(args, " "))
}

func printGitMissingAdvice() {
	printGitMissingHeader()
	printGitInstallOptions()
}

func printGitMissingHeader() {
	fmt.Println()
	fmt.Println("✗ Git is not installed or could not be found.")
	fmt.Println()
	fmt.Println("GitMap requires Git to execute bash and manage repositories.")
	fmt.Println("You can install Git using one of the following commands:")
	fmt.Println()
}

func printGitInstallOptions() {
	fmt.Println("  Local Installation:")
	fmt.Println("    ● gitmap install git")
	fmt.Println("    ● gitmap compact")
	fmt.Println()
	fmt.Println("  Remote SSH Installation:")
	fmt.Println("    ● gitmap ssh install git --target <node-alias>")
	fmt.Println("    ● gitmap ssh exec <node-alias> \"sudo apt-get update && sudo apt-get install -y git\"")
	fmt.Println("    ● gitmap ssh exec <node-alias> \"winget install --id Git.Git -e --source winget\"")
	fmt.Println()
}
