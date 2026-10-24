package cmdtoolchain

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/charmbracelet/lipgloss"
)

// RunRust is the main entry point for gitmap rust / rst / rustc.
func RunRust(args []string) error {
	checkToolchainHelp("rust", args)
	if len(args) == 0 {
		return RunRustStatus()
	}

	sub := strings.ToLower(args[0])
	if isStatusSubcommand(sub) {
		return RunRustStatus()
	}

	return dispatchRustSubcommand(sub, args)
}

func dispatchRustSubcommand(sub string, args []string) error {
	if sub == "c" {
		return forwardToRustc(args[1:])
	}

	if sub == "up" {
		return forwardToRustup(args[1:])
	}

	return forwardToCargo(args)
}

func forwardToRustc(args []string) error {
	bin := ResolveRustcBinary()
	if bin == "" {
		return apperror.NewNotFoundError("rustc is not installed or not found in PATH")
	}

	return ExecEnrichedCommand(bin, args)
}

func forwardToRustup(args []string) error {
	bin := ResolveRustupBinary()
	if bin == "" {
		return apperror.NewNotFoundError("rustup is not installed or not found in PATH")
	}

	return ExecEnrichedCommand(bin, args)
}

func forwardToCargo(args []string) error {
	bin := ResolveCargoBinary()
	if bin == "" {
		return apperror.NewNotFoundError("cargo is not installed or not found in PATH")
	}

	return ExecEnrichedCommand(bin, args)
}

// ResolveRustcBinary locates rustc on the system.
func ResolveRustcBinary() string {
	return resolveBinaryCandidate("rustc")
}

// ResolveRustupBinary locates rustup on the system.
func ResolveRustupBinary() string {
	return resolveBinaryCandidate("rustup")
}

// ResolveClangBinary locates clang on the system.
func ResolveClangBinary() string {
	if p, err := exec.LookPath("clang"); err == nil {
		return p
	}

	dir := ResolveLibclangDir()
	if dir == "" {
		return ""
	}

	target := filepath.Join(dir, "clang.exe")
	if isRegularFile(target) {
		return target
	}

	return ""
}

// RunRustStatus renders toolchain dashboard.
func RunRustStatus() error {
	renderStatusHeader()
	renderBinaryStatus("Rustc", ResolveRustcBinary())
	renderBinaryStatus("Cargo", ResolveCargoBinary())
	renderBinaryStatus("Rustup", ResolveRustupBinary())
	renderBinaryStatus("Clang", ResolveClangBinary())
	renderLibclangStatus()

	return nil
}

func renderStatusHeader() {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	fmt.Println(headerStyle.Render("Rust & LLVM Toolchain Status:"))
	fmt.Println()
}

func renderBinaryStatus(name, path string) {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")).Width(16)
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8"))
	if path == "" {
		fmt.Printf("  %s %s\n", keyStyle.Render(name+":"), warnStyle.Render("Not installed"))

		return
	}

	ver := queryBinaryVersion(path)
	fmt.Printf("  %s %s %s\n", keyStyle.Render(name+":"), valStyle.Render(ver), dimStyle.Render("("+path+")"))
}

func renderLibclangStatus() {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")).Width(16)
	valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8"))
	dir := ResolveLibclangDir()
	if dir == "" {
		fmt.Printf("  %s %s\n", keyStyle.Render("LIBCLANG_PATH:"), dimStyle.Render("Not found / not set"))

		return
	}

	fmt.Printf("  %s %s\n", keyStyle.Render("LIBCLANG_PATH:"), valStyle.Render(dir))
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
