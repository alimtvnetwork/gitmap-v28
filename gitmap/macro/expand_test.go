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
	tempDir := getNormalizedTempDir()

	tests := []struct {
		input string
		want  string
	}{
		{"cd %temp%", "cd " + tempDir},
		{"cd %TEMP%", "cd " + tempDir},
		{"cd $TEMP", "cd " + tempDir},
		{"cd $TMP", "cd " + tempDir},
		{"cd //temp", "cd " + tempDir},
		{"cd /temp", "cd " + tempDir},
		{"cd //tmp", "cd " + tempDir},
		{"cd /tmp", "cd " + tempDir},
		{"mkdir //temp/subproject", "mkdir " + filepath.Join(tempDir, "subproject")},
		{"mkdir /temp/subproject", "mkdir " + filepath.Join(tempDir, "subproject")},
	}

	for _, tc := range tests {
		got := ExpandPathAndEnv(tc.input)
		if got != tc.want {
			t.Errorf("ExpandPathAndEnv(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
