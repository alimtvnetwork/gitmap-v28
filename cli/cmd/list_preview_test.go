package cmd

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestPrintListPreview(t *testing.T) {
	records := []model.ScanRecord{
		{
			RepoName:     "alpha",
			AbsolutePath: "D:/work/alpha",
		},
		{
			RepoName:     "beta",
			AbsolutePath: "D:/work/beta",
		},
	}

	out, _ := captureStdout(t, func() int {
		PrintListPreview(records)
		return 0
	})

	if !strings.Contains(out, "alpha") || !strings.Contains(out, "D:/work/alpha") {
		t.Errorf("unexpected preview output:\n%s", out)
	}
	if !strings.Contains(out, "beta") || !strings.Contains(out, "D:/work/beta") {
		t.Errorf("unexpected preview output:\n%s", out)
	}
	if !strings.Contains(out, "[git]") {
		t.Errorf("expected [git] in preview output:\n%s", out)
	}
}

func TestPrintListPreviewWithOptions(t *testing.T) {
	records := []model.ScanRecord{
		{
			RepoName:     "alpha",
			Branch:       "feature-1",
			AbsolutePath: "D:/work/alpha",
		},
	}

	opts := ListPreviewOptions{
		ShowNumbers: false,
		ShowEmoji:   false,
	}

	out, _ := captureStdout(t, func() int {
		PrintListPreviewWithOptions(records, opts)
		return 0
	})

	if strings.Contains(out, "1. ") {
		t.Errorf("expected no sequence numbers when ShowNumbers is false:\n%s", out)
	}
	if !strings.Contains(out, "[git: feature-1]") {
		t.Errorf("expected [git: feature-1] in output:\n%s", out)
	}
}

func TestPrintListTree(t *testing.T) {
	records := []model.ScanRecord{
		{
			RepoName:     "alpha",
			AbsolutePath: "D:/work/repos/alpha",
		},
	}

	out, _ := captureStdout(t, func() int {
		PrintListTree(records)
		return 0
	})

	if !strings.Contains(out, "📁") {
		t.Errorf("expected emoji folder icon in tree output:\n%s", out)
	}
	if !strings.Contains(out, "alpha") {
		t.Errorf("expected repo name in tree output:\n%s", out)
	}
}
