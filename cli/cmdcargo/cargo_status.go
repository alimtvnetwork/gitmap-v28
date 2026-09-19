package cmdcargo

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	keyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")).Width(14)
	valStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555"))
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
)

func runCargoStatus() error {
	fmt.Println(headerStyle.Render("Cargo & Rust Toolchain Status:"))
	fmt.Println()

	printToolStatus("Cargo", ResolveCargoBinary())
	printToolStatus("Rustc", ResolveRustcBinary())
	printToolStatus("Rustup", ResolveRustupBinary())
	fmt.Println()

	if ResolveCargoBinary() == "" {
		fmt.Println(dimStyle.Render("  To install Cargo, run: ") + lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd")).Render("gitmap install cargo"))
		fmt.Println()
	}

	return nil
}

func printToolStatus(name, path string) {
	if path == "" {
		fmt.Printf("  %s %s\n", keyStyle.Render(name+":"), warnStyle.Render("Not installed"))

		return
	}

	ver := queryBinaryVersion(path)
	fmt.Printf("  %s %s %s\n", keyStyle.Render(name+":"), valStyle.Render(ver), dimStyle.Render("("+path+")"))
}

func queryBinaryVersion(path string) string {
	out, err := exec.Command(path, "--version").CombinedOutput()
	if err != nil {
		return "installed"
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) > 0 {
		return lines[0]
	}

	return "installed"
}
