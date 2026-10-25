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

func TestCleanFlexiblePath_AtBracketFormat(t *testing.T) {
	raw := `@[C:\Users\Administrator\...\media_1792894009084_45b4893f.png]`
	want := filepath.Clean(`C:\Users\Administrator\...\media_1792894009084_45b4893f.png`)

	assertCleanAndExpand(t, raw, want)
}

func TestCleanFlexiblePath_MarkdownImageFormat(t *testing.T) {
	raw := `![Attached Image 1](` + `file:/` + `//C:/Users/Administrator/.../media_1792894009084_45b4893f.png)`
	want := filepath.Clean(`C:/Users/Administrator/.../media_1792894009084_45b4893f.png`)

	assertCleanAndExpand(t, raw, want)
}

func TestCleanFlexiblePath_NumberedRepoRelative(t *testing.T) {
	raw := `1. @[assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo`
	want := filepath.Clean(`assets/media_1792894009084_45b4893f.png`)

	assertCleanAndExpand(t, raw, want)
}

func TestCleanFlexiblePath_NumberedRepoCacheRelative(t *testing.T) {
	raw := `1. @[../repo-cache/assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo outside of the repo`
	want := filepath.Clean(`../repo-cache/assets/media_1792894009084_45b4893f.png`)

	assertCleanAndExpand(t, raw, want)
}

func assertCleanAndExpand(t *testing.T, raw, want string) {
	t.Helper()
	gotClean := CleanFlexiblePath(raw)
	if gotClean != want {
		t.Errorf("CleanFlexiblePath(%q) = %q, want %q", raw, gotClean, want)
	}

	gotExpand := ExpandPathAndEnv(raw)
	if gotExpand != want {
		t.Errorf("ExpandPathAndEnv(%q) = %q, want %q", raw, gotExpand, want)
	}
}

func TestCleanFlexiblePath_RepoRelativeExisting(t *testing.T) {
	repoRoot := getTestRepoRoot(t)
	tempAsset := createTempAssetFile(t, repoRoot, "assets/test_expand_flexible_asset.png")
	defer os.Remove(tempAsset)

	raw := `1. @[assets/test_expand_flexible_asset.png] # consider the file is relative from the repo`
	assertCleanAndExpand(t, raw, tempAsset)
}

func TestCleanFlexiblePath_RepoCacheExisting(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Skip("Cannot get cwd")
	}

	cwdParent := filepath.Dir(cwd)
	tempAsset := createTempAssetFile(t, cwdParent, "repo-cache/assets/test_expand_flexible_rc.png")
	defer os.Remove(tempAsset)

	raw := `1. @[../repo-cache/assets/test_expand_flexible_rc.png] # consider the file is relative from the repo outside of the repo`
	assertCleanAndExpand(t, raw, tempAsset)
}

func getTestRepoRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Skip("Cannot get cwd")
	}

	root := resolveRepoRoot(cwd)
	if len(root) == 0 {
		t.Skip("Cannot resolve repo root")
	}

	return root
}

func createTempAssetFile(t *testing.T, baseDir, rel string) string {
	t.Helper()
	target := filepath.Join(baseDir, filepath.FromSlash(rel))
	os.MkdirAll(filepath.Dir(target), 0755)
	err := os.WriteFile(target, []byte("fake"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp asset: %v", err)
	}

	return target
}
