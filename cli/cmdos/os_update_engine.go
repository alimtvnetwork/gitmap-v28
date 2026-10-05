package cmdos

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func discoverUpdateToolchains() []UpdateToolchain {
	candidates := []UpdateToolchain{
		{Name: "winget", Binary: "winget", UpdateArgs: []string{"upgrade"}, UpgradeArgs: []string{"upgrade", "--all", "--include-unknown"}},
		{Name: "apt", Binary: "apt-get", UpdateArgs: []string{"update"}, UpgradeArgs: []string{"dist-upgrade", "-y"}, NeedsSudo: true, RequiresSudo: true},
		{Name: "dnf", Binary: "dnf", UpdateArgs: []string{"check-update"}, UpgradeArgs: []string{"upgrade", "-y"}, NeedsSudo: true, RequiresSudo: true},
		{Name: "pacman", Binary: "pacman", UpdateArgs: []string{"-Sy"}, UpgradeArgs: []string{"-Syu", "--noconfirm"}, NeedsSudo: true, RequiresSudo: true},
		{Name: "flatpak", Binary: "flatpak", UpdateArgs: []string{"update", "--appstream"}, UpgradeArgs: []string{"update", "-y"}},
		{Name: "snap", Binary: "snap", UpdateArgs: []string{"refresh", "--list"}, UpgradeArgs: []string{"refresh"}},
		{Name: "brew", Binary: "brew", UpdateArgs: []string{"update"}, UpgradeArgs: []string{"upgrade"}},
	}

	var discovered []UpdateToolchain
	for _, tc := range candidates {
		if _, err := exec.LookPath(tc.Binary); err == nil {
			discovered = append(discovered, tc)
		}
	}

	return discovered
}

func isLinuxNonRoot() bool {
	return runtime.GOOS == "linux" && os.Geteuid() != 0
}

func hasSudoBinary() bool {
	_, err := exec.LookPath("sudo")
	return err == nil
}

func isElevationRequired(tc UpdateToolchain) bool {
	return (tc.NeedsSudo || tc.RequiresSudo) && isLinuxNonRoot() && hasSudoBinary()
}

func resolveUpdateArgs(tc UpdateToolchain, isUpgrade bool) []string {
	if isUpgrade {
		return tc.UpgradeArgs
	}
	return tc.UpdateArgs
}

func printDryRunToolchain(tc UpdateToolchain, args []string, isElevated bool) UpdateResult {
	prefix := ""
	if isElevated {
		prefix = "sudo -n "
	}
	fmt.Printf("  • [dry-run] %s%s %s\n", prefix, tc.Binary, strings.Join(args, " "))

	return UpdateResult{Name: tc.Name, Success: true}
}

func buildUpdateCmd(binary string, args []string, isElevated bool) *exec.Cmd {
	if isElevated {
		return exec.Command("sudo", append([]string{"-n", binary}, args...)...)
	}
	return exec.Command(binary, args...)
}

func executeToolchain(tc UpdateToolchain, isUpgrade bool, isDryRun bool) UpdateResult {
	args := resolveUpdateArgs(tc, isUpgrade)
	isElevated := isElevationRequired(tc)
	if isDryRun {
		return printDryRunToolchain(tc, args, isElevated)
	}

	cmd := buildUpdateCmd(tc.Binary, args, isElevated)
	out, err := cmd.CombinedOutput()

	return UpdateResult{
		Name:    tc.Name,
		Success: err == nil,
		Output:  string(out),
		Error:   err,
	}
}

func runAllUpdates(isUpgrade bool, isDryRun bool) []UpdateResult {
	toolchains := discoverUpdateToolchains()
	var results []UpdateResult

	for _, tc := range toolchains {
		res := executeToolchain(tc, isUpgrade, isDryRun)
		results = append(results, res)
	}

	return results
}
