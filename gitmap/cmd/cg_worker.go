package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

var cgSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
var cgErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555"))
var cgHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Bold(true)
var cgVersionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9"))

type cgUpdateResult struct {
	repo       string
	isFail     bool
	errorMsg   string
	oldVersion string
	newVersion string
	hasChanged bool
	stdout     string
}

func executeCGWorkers(repos []string) {
	fmt.Println(cgHeaderStyle.Render(fmt.Sprintf("Running coding guidelines installer on %d repositories...", len(repos))))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)
	results := make(chan cgUpdateResult, len(repos))
	for _, repo := range repos {
		wg.Add(1)
		sem <- struct{}{}
		go runCgWorker(repo, &wg, sem, results)
	}
	wg.Wait()
	close(results)
	printCGUpdateSummary(results)
}

func runCgWorker(repo string, wg *sync.WaitGroup, sem chan struct{}, results chan<- cgUpdateResult) error {
	defer wg.Done()
	defer func() { <-sem }()
	res := cgUpdateResult{repo: repo, isFail: false}
	if oldMeta, err := ReadCGMetadata(repo); err == nil {
		res.oldVersion = oldMeta.Version
	}
	out, runErr := runCgScriptInRepo(repo)
	if runErr != nil {
		res.isFail = true
		res.errorMsg = runErr.Error()
	}
	res.stdout = out
	if newMeta, err := ReadCGMetadata(repo); err == nil {
		res.newVersion = newMeta.Version
	}
	res.hasChanged = (res.oldVersion != res.newVersion)
	results <- res
	return nil
}

func getCgScriptCmd(repo string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		psCmd := "irm https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.ps1 | iex"
		return exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	}
	shCmd := "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.sh | bash"
	return exec.Command("bash", "-c", shCmd)
}

func runCgScriptInRepo(repo string) (string, error) {
	cmd := getCgScriptCmd(repo)
	cmd.Dir = repo
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func printCGUpdateSummary(results <-chan cgUpdateResult) {
	fmt.Println(cgHeaderStyle.Render("\n--- Update Summary ---"))
	for r := range results {
		switch {
		case r.isFail:
			fmt.Printf("%s [%s] Failed: %s\n", cgErrorStyle.Render("X"), r.repo, r.errorMsg)
		case !r.hasChanged:
			fmt.Printf("%s [%s] Up to date (%s)\n", cgSuccessStyle.Render("OK"), r.repo, cgVersionStyle.Render(r.newVersion))
		default:
			fmt.Printf("%s [%s] Updated: %s -> %s\n", cgSuccessStyle.Render("OK"), r.repo, cgErrorStyle.Render(r.oldVersion), cgSuccessStyle.Render(r.newVersion))
			printModifiedFiles(r.stdout)
		}
	}
	fmt.Println(cgHeaderStyle.Render("Done."))
}

func printModifiedFiles(stdout string) {
	lines := strings.Split(stdout, "\n")
	hasFiles := false
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "downloaded") || strings.Contains(lower, "created") || strings.Contains(lower, "updated") {
			fmt.Printf("    - %s\n", strings.TrimSpace(line))
			hasFiles = true
		}
	}
	if !hasFiles {
		fmt.Printf("    - version.json\n    - .lovable/coding-guidelines/\n")
	}
}
