package cmdclone

import (
	"strings"
	"testing"
)

func TestParseMultiCloneText_UserSample(t *testing.T) {
	input := "```\n" +
		"ChrisTitusTech/ChrisTitusTech\n" +
		"https://github.com/ChrisTitusTech/ChrisTitusTech\n\n" +
		"ChrisTitusTech/linutil: Chris Titus Tech's Linux Toolbox - Linutil is a distro-agnostic toolbox designed to simplify everyday Linux tasks.\n" +
		"https://github.com/ChrisTitusTech/linutil\n\n" +
		"ChrisTitusTech/winutil: Chris Titus Tech's Windows Utility - Install Programs, Tweaks, Fixes, and Updates\n" +
		"https://github.com/ChrisTitusTech/winutil\n\n" +
		"ChrisTitusTech/image-tools: CLI/Thunar Image Resize, Upscale, and Webp conversion\n" +
		"https://github.com/ChrisTitusTech/image-tools\n\n" +
		"ChrisTitusTech/winutil: Chris Titus Tech's Windows Utility - Install Programs, Tweaks, Fixes, and Updates\n" +
		"https://github.com/ChrisTitusTech/winutil\n\n" +
		"ChrisTitusTech/dwm-titus: My Linux Desktop - Fedora/X11 Desktop Environment\n" +
		"https://github.com/ChrisTitusTech/dwm-titus\n\n" +
		"ChrisTitusTech/warframe-linux: Planning a Linux-first Warframe companion: inventory, market prices, relic rewards, and crafting wi\n" +
		"https://github.com/ChrisTitusTech/warframe-linux\n" +
		"```"

	urls := ParseMultiCloneText(input)
	if len(urls) != 6 {
		t.Fatalf("expected 6 unique URLs from user sample, got %d: %v", len(urls), urls)
	}

	expectedRepos := []string{
		"ChrisTitusTech",
		"linutil",
		"winutil",
		"image-tools",
		"dwm-titus",
		"warframe-linux",
	}

	for i, expected := range expectedRepos {
		isMatching := strings.Contains(urls[i], expected)
		if !isMatching {
			t.Errorf("item %d: expected %s, got %s", i, expected, urls[i])
		}
	}
}

func TestBuildMultiCloneItems_TargetDir(t *testing.T) {
	urls := []string{
		"https://github.com/owner/repo1.git",
		"https://github.com/owner/repo2",
	}

	items := BuildMultiCloneItems(urls, "custom-dir")
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}

	if items[0].RepoName != "repo1" || !strings.Contains(items[0].TargetDir, "custom-dir") {
		t.Errorf("unexpected item 0 target dir: %s", items[0].TargetDir)
	}
}

func TestResolveMultiCloneOptions_Flags(t *testing.T) {
	args := []string{"-d", "workspace", "--dry-run", "owner/testrepo"}
	opts, err := ResolveMultiCloneOptions(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.TargetDir != "workspace" {
		t.Errorf("expected targetDir workspace, got %s", opts.TargetDir)
	}
	if !opts.IsDryRun {
		t.Errorf("expected dry-run true")
	}
	if !strings.Contains(opts.RawInput, "owner/testrepo") {
		t.Errorf("expected raw input to contain owner/testrepo, got %s", opts.RawInput)
	}
}
