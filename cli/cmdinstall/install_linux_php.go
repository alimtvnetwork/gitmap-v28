package cmdinstall

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PHP on Linux: the distro `php` package is NOT the 8.5 line
// (Ubuntu noble ships PHP 8.3), and `apt install php=8.5.12` fails
// because no such apt version exists. This handler installs the
// 8.5 line from the ondrej/php PPA instead, and pins the exact
// build when a version is requested (for example 8.5.12).
//
// Other managers keep the generic path:
//   - Windows winget: PHP.PHP.8.5 with --version (see WingetPkgPHP)
//   - Windows choco: php with --version
//   - macOS brew: php (latest 8.5.x; Homebrew cannot pin a patch
//     version, so brew never claims an exact 8.5.12)

const (
	phpAptLine    = "8.5"
	phpAptPackage = "php8.5"
	phpPpa        = "ondrej/php"
)

var phpAptExtensions = []string{
	"php8.5-cli", "php8.5-common", "php8.5-xml", "php8.5-mbstring",
	"php8.5-curl", "php8.5-zip", "php8.5-sqlite3", "php8.5-intl",
}

// runInstallPhpLinux installs PHP 8.5 on Ubuntu/Debian via the
// ondrej/php PPA. A forced non-apt manager falls back to the
// generic installer so --manager brew/choco/winget keeps working.
func runInstallPhpLinux(opts installOptions) error {
	if opts.Manager != "" && opts.Manager != constants.PkgMgrApt {
		executeGenericInstall(opts)

		return nil
	}

	printPhpInstallPlan(opts)
	if opts.DryRun {
		printPhpDryRun(opts)

		return nil
	}

	if installed := detectInstalledVersion(constants.ToolPHP); installed != "" {
		if parsed := parseVersionFromOutput(installed); phpVersionMatches(parsed, opts.Version) {
			fmt.Printf(constants.MsgInstallFound, constants.ToolPHP, parsed)

			return nil
		}
	}

	if err := ensurePhpPpa(opts.Verbose); err != nil {
		return err
	}

	exact, err := resolvePhpAptVersion(opts.Version)
	if err != nil {
		return err
	}

	installArgs := phpAptInstallArgs(exact)
	if err := runPhaseWithAudit(constants.ToolPHP, constants.PkgMgrApt, opts.Version, installArgs, "php-apt-install", opts.Verbose); err != nil {
		return err
	}

	return verifyPhpVersion(opts.Version)
}

func printPhpInstallPlan(opts installOptions) {
	version := opts.Version
	if version == "" {
		version = phpAptLine + " (PPA latest)"
	}

	fmt.Printf("\n  +-- Install Plan ---------------------\n")
	fmt.Printf("  | Tool:    %s\n", constants.ToolPHP)
	fmt.Printf("  | Version: %s\n", version)
	fmt.Printf("  | Manager: %s (ppa:%s)\n", constants.PkgMgrApt, phpPpa)
	fmt.Printf("  | Command: %s\n", strings.Join(phpAptInstallArgs(""), " "))
	fmt.Printf("  +--------------------------------------\n\n")
}

func printPhpDryRun(opts installOptions) {
	fmt.Printf(constants.MsgInstallDryCmd, "sudo add-apt-repository -y ppa:"+phpPpa)
	fmt.Printf(constants.MsgInstallDryCmd, "sudo apt-get update")
	fmt.Printf(constants.MsgInstallDryCmd, strings.Join(phpAptInstallArgs(""), " "))
}

// ensurePhpPpa adds the ondrej/php PPA once, then refreshes apt.
func ensurePhpPpa(verbose bool) error {
	if !isCommandAvailable("add-apt-repository") {
		if err := runPhaseWithAudit(constants.ToolPHP, constants.PkgMgrApt, "", []string{"sudo", "apt-get", "update"}, "php-apt-update", verbose); err != nil {
			return err
		}
		if err := runPhaseWithAudit(constants.ToolPHP, constants.PkgMgrApt, "", []string{"sudo", "apt", "install", "-y", "software-properties-common"}, "php-ppa-prereq", verbose); err != nil {
			return err
		}
	}

	if !isPhpPpaPresent() {
		if err := runPhaseWithAudit(constants.ToolPHP, constants.PkgMgrApt, "", []string{"sudo", "add-apt-repository", "-y", "ppa:" + phpPpa}, "php-add-ppa", verbose); err != nil {
			return err
		}
	}

	return runPhaseWithAudit(constants.ToolPHP, constants.PkgMgrApt, "", []string{"sudo", "apt-get", "update"}, "php-apt-update", verbose)
}

func isPhpPpaPresent() bool {
	cmd := exec.Command("sh", "-c", "grep -rq "+phpPpa+" /etc/apt/sources.list.d/ 2>/dev/null")

	return cmd.Run() == nil
}

// resolvePhpAptVersion maps a requested version (for example
// 8.5.12) to the exact apt version string the PPA publishes
// (for example 8.5.12-1+ubuntu24.04.1+deb.sury.org+1). An empty
// request means "latest 8.5 from the PPA" and returns "".
func resolvePhpAptVersion(requested string) (string, error) {
	if requested == "" {
		return "", nil
	}

	out, err := exec.Command("apt-cache", "madison", phpAptPackage).Output()
	if err != nil {
		return "", fmt.Errorf("could not list %s versions from apt (is the ondrej/php PPA added?): %w", phpAptPackage, err)
	}

	exact, available := parseAptMadisonVersion(string(out), requested)
	if exact == "" {
		return "", fmt.Errorf("PHP %s is not available from the ondrej/php PPA on this system (available: %s)", requested, strings.Join(available, ", "))
	}

	return exact, nil
}

// parseAptMadisonVersion parses `apt-cache madison` output. Each
// line is "pkg | version | source". It returns the first version
// that equals the request or starts with request + "-".
func parseAptMadisonVersion(madison, requested string) (string, []string) {
	var available []string

	for _, line := range strings.Split(madison, "\n") {
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			continue
		}

		version := strings.TrimSpace(parts[1])
		if version == "" {
			continue
		}

		available = append(available, version)
		if version == requested || strings.HasPrefix(version, requested+"-") {
			return version, available
		}
	}

	return "", available
}

// phpAptInstallArgs builds the apt install command. The main
// package is pinned to the exact PPA version when known; the
// extensions resolve to the matching build automatically.
func phpAptInstallArgs(exact string) []string {
	args := []string{"sudo", "apt", "install", "-y"}
	target := phpAptPackage
	if exact != "" {
		target = phpAptPackage + "=" + exact
	}

	args = append(args, target)
	args = append(args, phpAptExtensions...)

	return args
}

// phpVersionMatches reports whether an installed version satisfies
// the request: an empty request means any 8.5.x, otherwise the
// installed version must equal the request exactly.
func phpVersionMatches(installed, requested string) bool {
	if installed == "" {
		return false
	}

	if requested == "" {
		return phpMajorMinor(installed) == phpAptLine
	}

	return installed == requested
}

func phpMajorMinor(version string) string {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return version
	}

	return parts[0] + "." + parts[1]
}

// verifyPhpVersion checks `php -v` after install and fails
// honestly instead of reporting a wrong version as success.
func verifyPhpVersion(requested string) error {
	out, err := exec.Command("php", "-v").Output()
	if err != nil {
		return fmt.Errorf("php binary not found after install: %w", err)
	}

	installed := parseVersionFromOutput(string(out))
	if !phpVersionMatches(installed, requested) {
		if requested == "" {
			return fmt.Errorf("installed PHP version %s is not in the 8.5 line", installed)
		}

		return fmt.Errorf("installed PHP version %s does not match requested %s", installed, requested)
	}

	fmt.Printf(constants.MsgInstallSuccess, constants.ToolPHP+" "+installed)
	verifyInstallation(constants.ToolPHP)
	recordInstallation(constants.ToolPHP, constants.PkgMgrApt)

	return nil
}
