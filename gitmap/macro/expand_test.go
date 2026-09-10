package macro

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPathAndEnv_WindowsEnv(t *testing.T) {
	os.Setenv("TEST_MACRO_VAR", "my_custom_value")
	defer os.Unsetenv("TEST_MACRO_VAR")

	input := "echo %TEST_MACRO_VAR%/subfolder"
	got := ExpandPathAndEnv(input)
	want := "echo my_custom_value/subfolder"
	if got != want {
		t.Errorf("ExpandPathAndEnv(%q) = %q, want %q", input, got, want)
	}
}

func TestExpandPathAndEnv_UnixEnv(t *testing.T) {
	os.Setenv("TEST_UNIX_VAR", "hello_world")
	defer os.Unsetenv("TEST_UNIX_VAR")

	input := "cat $TEST_UNIX_VAR/file.txt"
	got := ExpandPathAndEnv(input)
	want := "cat hello_world/file.txt"
	if got != want {
		t.Errorf("ExpandPathAndEnv(%q) = %q, want %q", input, got, want)
	}
}

func TestExpandPathAndEnv_Tilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || len(home) == 0 {
		t.Skip("User home dir not available")
	}

	input := "cd ~/projects/app"
	got := ExpandPathAndEnv(input)
	want := "cd " + filepath.Join(home, "projects/app")
	if got != want {
		t.Errorf("ExpandPathAndEnv(%q) = %q, want %q", input, got, want)
	}
}

func TestExpandPathAndEnv_Temp(t *testing.T) {
	tempVal := os.Getenv("TEMP")
	if tempVal == "" {
		tempVal = os.Getenv("TMP")
	}
	if tempVal == "" {
		t.Skip("TEMP not set")
	}

	input := "cd %temp%"
	got := ExpandPathAndEnv(input)
	if got != "cd "+tempVal {
		t.Errorf("ExpandPathAndEnv(%q) = %q, want %q", input, got, "cd "+tempVal)
	}
}

func TestNormalizeTargetPath_TempAliases(t *testing.T) {
	expectedTemp := filepath.Clean(os.TempDir())
	aliases := []string{"//temp", "/temp", `\temp`, `\\temp`, "/tmp", "//tmp"}
	for _, alias := range aliases {
		got := NormalizeTargetPath(alias, "")
		if got != expectedTemp {
			t.Errorf("NormalizeTargetPath(%q) = %q, want %q", alias, got, expectedTemp)
		}
	}
}

func TestNormalizeTargetPath_TempSubdir(t *testing.T) {
	expected := filepath.Clean(filepath.Join(os.TempDir(), "subfolder"))
	aliases := []string{"//temp/subfolder", `\temp\subfolder`, "/tmp/subfolder"}
	for _, alias := range aliases {
		got := NormalizeTargetPath(alias, "")
		if got != expected {
			t.Errorf("NormalizeTargetPath(%q) = %q, want %q", alias, got, expected)
		}
	}
}

func TestNormalizeTargetPath_QuotesAndEnv(t *testing.T) {
	expectedTemp := filepath.Clean(os.TempDir())
	inputs := []string{`"%temp%"`, `'%temp%'`, "%temp%"}
	for _, in := range inputs {
		got := NormalizeTargetPath(in, "")
		if got != expectedTemp {
			t.Errorf("NormalizeTargetPath(%q) = %q, want %q", in, got, expectedTemp)
		}
	}
}

func TestNormalizeTargetPath_RelativeAndDash(t *testing.T) {
	if got := NormalizeTargetPath("-", "/base"); got != "-" {
		t.Errorf("NormalizeTargetPath('-') = %q, want '-'", got)
	}
	baseDir := filepath.Clean("/base/project")
	got := NormalizeTargetPath("src//sub", baseDir)
	want := filepath.Clean(filepath.Join(baseDir, "src/sub"))
	if got != want {
		t.Errorf("NormalizeTargetPath('src//sub') = %q, want %q", got, want)
	}
}
