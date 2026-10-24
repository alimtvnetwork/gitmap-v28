package cmdcargo

import (
	"fmt"
	"os/exec"
	"strings"

	cargoconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/constants"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(cargoconst.ColorDraculaPurple))
	keyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(cargoconst.ColorDraculaWhite)).Width(cargoconst.KeyColumnWidth)
	valStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(cargoconst.ColorDraculaGreen))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color(cargoconst.ColorDraculaRed))
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(cargoconst.ColorDraculaGray))
)

func runCargoStatus() error {
	fmt.Println(headerStyle.Render(cargoconst.StatusHeaderTitle))
	fmt.Println()

	printToolStatus(cargoconst.DisplayNameCargo, ResolveCargoBinary())
	printToolStatus(cargoconst.DisplayNameRustc, ResolveRustcBinary())
	printToolStatus(cargoconst.DisplayNameRustup, ResolveRustupBinary())
	fmt.Println()

	if ResolveCargoBinary() == "" {
		printInstallHint()
	}

	return nil
}

func printInstallHint() {
	cyanStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(cargoconst.ColorDraculaCyan))
	fmt.Println(dimStyle.Render(cargoconst.PromptToInstall) + cyanStyle.Render(cargoconst.SuggestGitMapInstall))
	fmt.Println()
}

func printToolStatus(name, path string) {
	if path == "" {
		fmt.Printf("  %s %s\n", keyStyle.Render(name+":"), warnStyle.Render(cargoconst.StatusNotInstalled))

		return
	}

	ver := queryBinaryVersion(path)
	fmt.Printf("  %s %s %s\n", keyStyle.Render(name+":"), valStyle.Render(ver), dimStyle.Render("("+path+")"))
}

func queryBinaryVersion(path string) string {
	out, err := exec.Command(path, cargoconst.FlagLongVersion).CombinedOutput()
	if err != nil {
		return cargoconst.StatusInstalled
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		return lines[0]
	}

	return cargoconst.StatusInstalled
}
