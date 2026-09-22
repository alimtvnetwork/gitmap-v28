package cmdos

import (
	"fmt"
	"os/exec"
	"strings"
)

func discoverUpdateToolchains() []UpdateToolchain {
	candidates := []UpdateToolchain{
		{"winget", "winget", []string{"upgrade"}, []string{"upgrade", "--all", "--include-unknown"}},
		{"apt", "apt-get", []string{"update"}, []string{"dist-upgrade", "-y"}},
		{"dnf", "dnf", []string{"check-update"}, []string{"upgrade", "-y"}},
		{"pacman", "pacman", []string{"-Sy"}, []string{"-Syu", "--noconfirm"}},
		{"flatpak", "flatpak", []string{"update", "--appstream"}, []string{"update", "-y"}},
		{"snap", "snap", []string{"refresh", "--list"}, []string{"refresh"}},
		{"brew", "brew", []string{"update"}, []string{"upgrade"}},
	}

	var discovered []UpdateToolchain
	for _, tc := range candidates {
		if _, err := exec.LookPath(tc.Binary); err == nil {
			discovered = append(discovered, tc)
		}
	}

	return discovered
}

func executeToolchain(tc UpdateToolchain, isUpgrade bool, isDryRun bool) UpdateResult {
	args := tc.UpdateArgs
	if isUpgrade {
		args = tc.UpgradeArgs
	}

	if isDryRun {
		fmt.Printf("  • [dry-run] %s %s\n", tc.Binary, strings.Join(args, " "))

		return UpdateResult{Name: tc.Name, Success: true}
	}

	cmd := exec.Command(tc.Binary, args...)
	out, err := cmd.CombinedOutput()
	isSuccess := err == nil

	return UpdateResult{
		Name:    tc.Name,
		Success: isSuccess,
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
