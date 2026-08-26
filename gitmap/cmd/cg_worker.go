package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

var cgSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
var cgErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555"))

func executeCGWorkers(repos []string) {
	fmt.Printf("Running coding guidelines installer on %d repositories...\n", len(repos))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for _, repo := range repos {
		wg.Add(1)
		sem <- struct{}{}
		go runCgWorker(repo, &wg, sem)
	}
	wg.Wait()
	fmt.Println("Done.")
}

func runCgWorker(repo string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }()

	if err := runCgScriptInRepo(repo); err != nil {
		fmt.Printf("%s [%s] %v\n", cgErrorStyle.Render("✖"), repo, err)
		return
	}
	fmt.Printf("%s [%s]\n", cgSuccessStyle.Render("✔"), repo)
}

func runCgScriptInRepo(repo string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		psCmd := "irm https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.ps1 | iex"
		cmd = exec.Command("powershell", "-NoProfile", "-Command", psCmd)
	} else {
		shCmd := "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/coding-guidelines-v24/main/install.sh | bash"
		cmd = exec.Command("bash", "-c", shCmd)
	}
	cmd.Dir = repo
	return cmd.Run()
}
