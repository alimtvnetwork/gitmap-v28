package cmdtoolchain

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestIsValidLibclangDir(t *testing.T) {
	if isValidLibclangDir("non-existent-dir-for-test-12345") {
		t.Errorf("expected non-existent dir to be invalid")
	}

	tempDir := t.TempDir()
	if isValidLibclangDir(tempDir) {
		t.Errorf("expected empty dir to be invalid")
	}

	createMockLibclangFile(t, tempDir)
	if !isValidLibclangDir(tempDir) {
		t.Errorf("expected dir with libclang file to be valid")
	}
}

func createMockLibclangFile(t *testing.T, dir string) {
	name := "libclang.so"
	if runtime.GOOS == "windows" {
		name = "libclang.dll"
	}
	if runtime.GOOS == "darwin" {
		name = "libclang.dylib"
	}

	filePath := filepath.Join(dir, name)
	err := os.WriteFile(filePath, []byte("mock"), 0644)
	if err != nil {
		t.Fatalf("failed to create mock libclang: %v", err)
	}
}

func TestResolveLibclangDir_Env(t *testing.T) {
	tempDir := t.TempDir()
	createMockLibclangFile(t, tempDir)

	orig := os.Getenv("LIBCLANG_PATH")
	defer os.Setenv("LIBCLANG_PATH", orig)

	os.Setenv("LIBCLANG_PATH", tempDir)
	resolved := ResolveLibclangDir()
	if resolved != tempDir {
		t.Errorf("expected %q, got %q", tempDir, resolved)
	}
}

func TestEnrichEnv(t *testing.T) {
	baseEnv := []string{"PATH=/usr/bin", "FOO=BAR"}
	baseBin := filepath.Join("opt", "cargo", "bin", "cargo")
	libclangDir := filepath.Join("opt", "llvm", "lib")

	enriched := EnrichEnv(baseEnv, baseBin, libclangDir)
	m := sliceToEnvMap(enriched)

	verifyEnrichedEnv(t, m, baseBin, libclangDir)
}

func verifyEnrichedEnv(t *testing.T, m map[string]string, baseBin, libclangDir string) {
	pathVal := m["PATH"]
	binDir := filepath.Dir(baseBin)
	if !strings.Contains(pathVal, binDir) {
		t.Errorf("PATH does not contain binDir: %s", pathVal)
	}

	if !strings.Contains(pathVal, libclangDir) {
		t.Errorf("PATH does not contain libclangDir: %s", pathVal)
	}

	if m["LIBCLANG_PATH"] != libclangDir {
		t.Errorf("expected LIBCLANG_PATH %s, got %s", libclangDir, m["LIBCLANG_PATH"])
	}
}

func TestKillOptions_Parse(t *testing.T) {
	opts1, err1 := parseKillArgs([]string{"--force", "--tree", "12345"})
	if err1 != nil || !opts1.HasForce || !opts1.HasTree || opts1.Target != "12345" {
		t.Errorf("unexpected parse result for opts1: %+v, err: %v", opts1, err1)
	}

	opts2, err2 := parseKillArgs([]string{"-p", "3000"})
	if err2 != nil || opts2.Target != ":3000" {
		t.Errorf("unexpected parse result for opts2: %+v, err: %v", opts2, err2)
	}

	opts3, err3 := parseKillArgs([]string{":8080"})
	if err3 != nil || opts3.Target != ":8080" {
		t.Errorf("unexpected parse result for opts3: %+v, err: %v", opts3, err3)
	}
}

func TestSafetyGuards(t *testing.T) {
	verifyProtectedPIDs(t)
	verifyProtectedNames(t)
}

func verifyProtectedPIDs(t *testing.T) {
	if !isProtectedPID(0) || !isProtectedPID(1) || !isProtectedPID(4) {
		t.Errorf("expected system PIDs to be protected")
	}

	if !isProtectedPID(os.Getpid()) {
		t.Errorf("expected self PID to be protected")
	}

	if isProtectedPID(999999) {
		t.Errorf("expected normal PID not to be protected")
	}
}

func verifyProtectedNames(t *testing.T) {
	if !isProtectedProcessName("antigravity") || !isProtectedProcessName("code.exe") {
		t.Errorf("expected IDEs to be protected")
	}

	if !isProtectedProcessName("pwsh") || !isProtectedProcessName("bash.exe") {
		t.Errorf("expected shells to be protected")
	}

	if isProtectedProcessName("my-custom-service") {
		t.Errorf("expected custom service not to be protected")
	}
}

func getMockNetstatData() string {
	return "Active Connections\n" +
		"  TCP    0.0.0.0:3000           0.0.0.0:0              LISTENING       14920\n" +
		"  TCP    127.0.0.1:8080         0.0.0.0:0              LISTENING       9988\n"
}

func TestParseNetstatOutput(t *testing.T) {
	mockOutput := getMockNetstatData()
	pid1, isFound1 := parseNetstatOutput(mockOutput, "3000")
	if !isFound1 || pid1 != 14920 {
		t.Errorf("expected 14920, got %d (found=%v)", pid1, isFound1)
	}

	pid2, isFound2 := parseNetstatOutput(mockOutput, "8080")
	if !isFound2 || pid2 != 9988 {
		t.Errorf("expected 9988, got %d (found=%v)", pid2, isFound2)
	}

	_, isFound3 := parseNetstatOutput(mockOutput, "9000")
	if isFound3 {
		t.Errorf("expected port 9000 not to be found")
	}
}

func TestSubcommands(t *testing.T) {
	if !isStatusSubcommand("status") || !isStatusSubcommand("st") || !isStatusSubcommand("info") {
		t.Errorf("expected status subcommands to be true")
	}

	if isStatusSubcommand("build") || isStatusSubcommand("test") {
		t.Errorf("expected build/test not to be status subcommands")
	}
}
