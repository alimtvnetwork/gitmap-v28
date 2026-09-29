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

	if !strings.Contains(out, "1. alpha\nD:/work/alpha\n\n") && !strings.Contains(out, "1. alpha\r\nD:/work/alpha\r\n\r\n") {
		t.Errorf("unexpected preview output:\n%s", out)
	}
	if !strings.Contains(out, "2. beta\nD:/work/beta\n\n") && !strings.Contains(out, "2. beta\r\nD:/work/beta\r\n\r\n") {
		t.Errorf("unexpected preview output:\n%s", out)
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
