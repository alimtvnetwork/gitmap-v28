package cmdos

import (
	"testing"
)

func TestParseDevCleanOptions(t *testing.T) {
	args := []string{"--dry-run", "-y", "--json", "--verbose", "--only", "go,npm"}
	opts := parseDevCleanOptions(args)
	assertDevCleanFlags(t, opts)
	assertDevCleanCategories(t, opts)
}

func assertDevCleanFlags(t *testing.T, opts DevCleanOptions) {
	if !opts.IsDryRun || !opts.HasAutoYes || !opts.IsJSON || !opts.IsVerbose {
		t.Errorf("expected all boolean flags to be true: %+v", opts)
	}
}

func assertDevCleanCategories(t *testing.T, opts DevCleanOptions) {
	if len(opts.OnlyCategories) != 2 || opts.OnlyCategories[0] != "go" || opts.OnlyCategories[1] != "npm" {
		t.Errorf("expected OnlyCategories [go, npm], got %v", opts.OnlyCategories)
	}
}

func TestIsHelpDevClean(t *testing.T) {
	if !isHelpDevClean([]string{"--help"}) {
		t.Errorf("expected --help to be recognized")
	}
	if !isHelpDevClean([]string{"-h"}) {
		t.Errorf("expected -h to be recognized")
	}
	if !isHelpDevClean([]string{"help"}) {
		t.Errorf("expected help to be recognized")
	}
	if isHelpDevClean([]string{"--dry-run"}) {
		t.Errorf("expected --dry-run not to be recognized as help")
	}
}

func TestIsDevTarget(t *testing.T) {
	validTargets := []string{
		"dev", "clean-dev", "dev-clean",
		"devtools", "devtool", "dev-tools", "dev-tool",
		"devtools-cache", "dev-tools-cache",
	}
	for _, target := range validTargets {
		if !isDevTarget(target) {
			t.Errorf("expected isDevTarget(%q) to be true", target)
		}
	}
	if isDevTarget("temp") {
		t.Errorf("expected temp to be false")
	}
}

func TestRunOSDevSubcommand_CleanRouting(t *testing.T) {
	err := runOSDevSubcommand([]string{"clean", "--dry-run", "--json"})
	if err != nil {
		t.Errorf("expected no error running os dev clean dry-run, got %v", err)
	}
}

func TestNormalizeEcosystemID_RustAndCargo(t *testing.T) {
	rustAliases := []string{"rust", "rustc", "cargo-cache", "crates"}
	for _, alias := range rustAliases {
		if normalizeEcosystemID(alias) != "cargo" {
			t.Errorf("expected %q to normalize to 'cargo', got %q", alias, normalizeEcosystemID(alias))
		}
	}
	if normalizeEcosystemID("ccache") != "ccache" {
		t.Errorf("expected ccache to normalize to ccache")
	}
}

func getCargoSafeTestPaths() []string {
	return []string{
		`C:\Users\Admin\.cargo\registry\cache`,
		`C:\Users\Admin\.cargo\registry\src`,
		`C:\Users\Admin\.cargo\git\db`,
		`C:\Users\Admin\.cargo\git\checkouts`,
		`C:\Users\Admin\.rustup\downloads`,
		`C:\Users\Admin\AppData\Local\Mozilla\sccache`,
		`/home/user/.cache/sccache`,
		`/home/user/.cargo/git/checkouts`,
	}
}

func TestIsDevCachePathSafe_CargoAndRust(t *testing.T) {
	for _, p := range getCargoSafeTestPaths() {
		if !isDevCachePathSafe(p) {
			t.Errorf("expected path %q to be recognized as safe dev cache path", p)
		}
	}
}

func TestResolveEnvCachePaths_CargoHomeNotRoot(t *testing.T) {
	paths := resolveEnvCachePaths("cargo", "CARGO_HOME", `C:\Users\Admin\.cargo`)
	for _, p := range paths {
		if p.Path == `C:\Users\Admin\.cargo` {
			t.Errorf("CARGO_HOME root should never be targeted directly for wipe")
		}
	}
}

func TestRunOSDevClean_DryRunExecution(t *testing.T) {
	err := RunOSDevClean([]string{"--dry-run", "--json", "--only", "cargo,go"})
	if err != nil {
		t.Errorf("expected no error running os dev-clean dry-run with json, got %v", err)
	}
}
