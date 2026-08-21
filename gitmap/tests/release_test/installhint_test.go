package release_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

func TestShouldPrintInstallHint_GitmapRepos(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want bool
	}{
		{"HTTPS match", "https://github.com/alimtvnetwork/gitmap-v28.git", true},
		{"HTTPS without .git", "https://github.com/alimtvnetwork/gitmap-v28", true},
		{"SSH match", "git@github.com:alimtvnetwork/gitmap-v28.git", true},
		{"Mixed case", "https://GitHub.com/AlimTVNetwork/Gitmap-V2.git", true},
		{"Subpath match", "https://github.com/alimtvnetwork/gitmap-v28/tree/main", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := release.ShouldPrintInstallHint(tc.url)
			if got != tc.want {
				t.Errorf("ShouldPrintInstallHint(%q) = %v, want %v", tc.url, got, tc.want)
			}
		})
	}
}

func TestShouldPrintInstallHint_NonGitmapRepos(t *testing.T) {
	cases := []struct {
		name string
		url  string
	}{
		{"Different org", "https://github.com/otherorg/gitmap-v28.git"},
		{"Different repo", "https://github.com/alimtvnetwork/other-repo.git"},
		{"Unrelated repo", "https://github.com/user/myproject.git"},
		{"Empty string", ""},
		{"Partial prefix", "https://github.com/alimtvnetwork/gitmap.git"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := release.ShouldPrintInstallHint(tc.url)
			if got {
				t.Errorf("ShouldPrintInstallHint(%q) = true, want false", tc.url)
			}
		})
	}
}

func TestPrintInstallHint_OutputContent(t *testing.T) {
	v, _ := release.Parse("2.60.0")

	output := fmt.Sprintf(constants.MsgInstallHintHeader, v.String()) +
		constants.MsgInstallHintWindows +
		constants.MsgInstallHintUnix

	if len(output) == 0 {
		t.Fatal("expected install hint output, got empty string")
	}

	checks := []string{
		"v2.60.0",
		"install.ps1",
		"install.sh",
		"PowerShell",
		"Linux",
	}

	for _, check := range checks {
		if !bytes.Contains([]byte(output), []byte(check)) {
			t.Errorf("output missing %q", check)
		}
	}
}
