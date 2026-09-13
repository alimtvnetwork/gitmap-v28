package cmdinstall

import (
	"runtime"
	"testing"
)

func TestIsInstallTarCommand_IdentifiesArchiveCommands(t *testing.T) {
	validCases := [][]string{
		{"tar", "app.tar.gz"},
		{"archive", "pkg.zip"},
		{"myapp.tar.gz"},
		{"tool.zip"},
		{"utility.tgz"},
		{"standalone.gz"},
		{"package.tar.xz"},
		{"package.tar.bz2"},
	}
	for _, tc := range validCases {
		if !isInstallTarCommand(tc) {
			t.Errorf("expected isInstallTarCommand(%v) to be true", tc)
		}
	}
}

func TestIsInstallTarCommand_RejectsNonArchiveCommands(t *testing.T) {
	invalidCases := [][]string{
		{"git"},
		{"agy"},
		{"node"},
		{"profile", "dev"},
		{},
	}
	for _, tc := range invalidCases {
		if isInstallTarCommand(tc) {
			t.Errorf("expected isInstallTarCommand(%v) to be false", tc)
		}
	}
}

func TestExtractInstallTarArgs_ExtractsPositionalArgs(t *testing.T) {
	args1 := extractInstallTarArgs([]string{"tar", "my-app.tar.gz", "--name", "app"})
	if len(args1) != 3 || args1[0] != "my-app.tar.gz" {
		t.Errorf("unexpected extracted args: %v", args1)
	}

	args2 := extractInstallTarArgs([]string{"pkg.zip", "--verbose"})
	if len(args2) != 2 || args2[0] != "pkg.zip" {
		t.Errorf("unexpected extracted args: %v", args2)
	}
}

func TestDeriveArchiveAppName_StripsKnownExtensions(t *testing.T) {
	tests := map[string]string{
		"myapp.tar.gz":        "myapp",
		"tool-v2.1.tgz":       "tool-v2.1",
		"package.zip":         "package",
		"single.gz":           "single",
		"complex.tar.xz":      "complex",
		"archive.tar.bz2":     "archive",
		"plain.tar":           "plain",
		"/path/to/nested.zip": "nested",
	}
	for input, expected := range tests {
		actual := deriveArchiveAppName(input)
		if actual != expected {
			t.Errorf("deriveArchiveAppName(%q) = %q; expected %q", input, actual, expected)
		}
	}
}

func TestParseInstallTarArgs_ParsesFlagsAndAppName(t *testing.T) {
	opts, err := parseInstallTarArgs([]string{"--name", "custom-app", "--dry-run", "bundle.tar.gz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.AppName != "custom-app" || !opts.DryRun || opts.ArchivePath != "bundle.tar.gz" {
		t.Errorf("unexpected parsed options: %+v", opts)
	}
}

func TestParseInstallTarArgs_RejectsMissingPath(t *testing.T) {
	_, err := parseInstallTarArgs([]string{"--name", "custom-app"})
	if err == nil {
		t.Errorf("expected error on missing archive path")
	}
}

func TestVerifyLinuxArchivePlatform_ValidatesCurrentOS(t *testing.T) {
	err := verifyLinuxArchivePlatform()
	if runtime.GOOS == "linux" && err != nil {
		t.Errorf("expected nil on Linux, got %v", err)
	}
	if runtime.GOOS != "linux" && err == nil {
		t.Errorf("expected E_LINUX_ONLY error on non-Linux, got nil")
	}
}

func TestResolvePackageStrategy_ResolvesBinaryApp(t *testing.T) {
	res := ArchiveInspectionResult{BinaryPath: "/tmp/bin/mytool"}
	strategy := resolvePackageStrategy(res)
	if strategy != StrategyBinaryApp {
		t.Errorf("expected StrategyBinaryApp, got %s", strategy)
	}
}

func TestResolvePackageStrategy_ResolvesScriptAndSource(t *testing.T) {
	scriptRes := ArchiveInspectionResult{ScriptPath: "/tmp/install.sh"}
	if resolvePackageStrategy(scriptRes) != StrategyScript {
		t.Errorf("expected StrategyScript")
	}

	sourceRes := ArchiveInspectionResult{HasMakefile: true}
	if resolvePackageStrategy(sourceRes) != StrategySource {
		t.Errorf("expected StrategySource")
	}

	unknownRes := ArchiveInspectionResult{}
	if resolvePackageStrategy(unknownRes) != StrategyUnknown {
		t.Errorf("expected StrategyUnknown")
	}
}
