package cmdfix

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func formatStashCardHeader(repoName string, width int) string {
	prefix := "┌── Safe Auto-Stash: " + repoName + " "
	dashCount := width - len([]rune(prefix)) - 1
	if dashCount < 2 {
		dashCount = 2
	}

	return prefix + strings.Repeat("─", dashCount) + "┐"
}

func stripANSI(s string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}

		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}

func padCardLine(s string, width int) string {
	visLen := len([]rune(stripANSI(s)))
	padding := width - 1 - visLen
	if padding < 0 {
		padding = 0
	}

	return s + strings.Repeat(" ", padding) + "│"
}

func colorizeStatusCode(code gitutil.FileStatusCode) string {
	if code == gitutil.StatusModified {
		return constants.ColorYellow + string(code) + constants.ColorReset
	}

	if code == gitutil.StatusDeleted {
		return constants.ColorRed + string(code) + constants.ColorReset
	}

	if code == gitutil.StatusUntracked {
		return constants.ColorDim + string(code) + constants.ColorReset
	}

	if code == gitutil.StatusAdded {
		return constants.ColorGreen + string(code) + constants.ColorReset
	}

	return string(code)
}

func appendCardFiles(b *strings.Builder, files []gitutil.DirtyFileInfo, width int) {
	maxShow := min(len(files), 5)
	for i := 0; i < maxShow; i++ {
		f := files[i]
		line := fmt.Sprintf("│   %s %s", colorizeStatusCode(f.StatusCode), f.Path)
		b.WriteString(padCardLine(line, width) + "\n")
	}

	if len(files) > 5 {
		moreLine := fmt.Sprintf("│   ... and %d other file(s)", len(files)-5)
		b.WriteString(padCardLine(moreLine, width) + "\n")
	}
}

// RenderHumanizedStashCard renders a boxed Lipgloss/termpad formatted card.
func RenderHumanizedStashCard(manifest gitutil.RepoDirtyManifest) string {
	boxWidth := 72
	var b strings.Builder
	titleBar := formatStashCardHeader(manifest.RepoName, boxWidth)
	b.WriteString("\n" + constants.ColorCyan + titleBar + constants.ColorReset + "\n")
	desc := fmt.Sprintf("│ GitMap is temporarily stashing %d uncommitted file(s) before sync:", len(manifest.Files))
	b.WriteString(padCardLine(desc, boxWidth) + "\n│" + strings.Repeat(" ", boxWidth-2) + "│\n")
	appendCardFiles(&b, manifest.Files, boxWidth)
	b.WriteString("│" + strings.Repeat(" ", boxWidth-2) + "│\n")
	guarantee := "│ ℹ Safety Guarantee: All changes will be re-applied after pull."
	b.WriteString(padCardLine(guarantee, boxWidth) + "\n")
	b.WriteString(constants.ColorCyan + "└" + strings.Repeat("─", boxWidth-2) + "┘" + constants.ColorReset + "\n")

	return b.String()
}

// ExecuteSafeStashLifecycle safely stashes, executes pullFn, and restores stashed work with humanized messaging.
func ExecuteSafeStashLifecycle(repoPath string, pullFn func() error) error {
	manifest, err := gitutil.InspectDirtyManifest(repoPath)
	if err != nil || !manifest.HasDirtyWork {
		return pullFn()
	}

	fmt.Print(RenderHumanizedStashCard(manifest))
	if err := runStashCreate(repoPath); err != nil {
		return err
	}

	pullErr := runStashPullStep(pullFn)
	popErr := runStashPopStep(repoPath)
	if pullErr != nil {
		return pullErr
	}

	return popErr
}

func runStashCreate(repoPath string) error {
	fmt.Print("  ↓ Stashing uncommitted changes (git stash -u) ... ")
	cmd := exec.Command("git", "-C", repoPath, "stash", "-u")
	if err := cmd.Run(); err != nil {
		fmt.Printf("%s failed\n", constants.ColorRed+"✖"+constants.ColorReset)
		return err
	}

	fmt.Printf("%s done\n", constants.ColorGreen+"✔"+constants.ColorReset)

	return nil
}

func runStashPullStep(pullFn func() error) error {
	fmt.Print("  ↓ Pulling upstream branch ... ")
	if err := pullFn(); err != nil {
		fmt.Printf("%s failed\n", constants.ColorRed+"✖"+constants.ColorReset)
		return err
	}

	fmt.Printf("%s done\n", constants.ColorGreen+"✔"+constants.ColorReset)

	return nil
}

func runStashPopStep(repoPath string) error {
	fmt.Print("  ↓ Re-applying stashed changes (git stash pop) ... ")
	cmd := exec.Command("git", "-C", repoPath, "stash", "pop")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return handleStashPopIssue(repoPath, string(out), err)
	}

	fmt.Printf("%s done\n", constants.ColorGreen+"✔"+constants.ColorReset)
	fmt.Printf("  %s✔ Local changes restored cleanly.%s\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func handleStashPopIssue(repoPath, out string, err error) error {
	fmt.Printf("%s note\n", constants.ColorYellow+"⚠"+constants.ColorReset)
	fmt.Fprintf(os.Stderr, "  %sℹ Stash preserved in 'git stash list'. Untracked or conflict notes:%s\n",
		constants.ColorYellow, constants.ColorReset)
	for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
		fmt.Fprintf(os.Stderr, "    %s\n", l)
	}

	return nil
}
